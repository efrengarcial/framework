package main

import (
	"context"
	"expvar"
	"fmt"
	"github.com/efrengarcial/framework/internal/platform/auth"
	"github.com/efrengarcial/framework/internal/platform/cache"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/caarlos0/env/v6"
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/user/delivery"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sagikazarmark/go-gin-gorm-opencensus/pkg/ocgorm"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

type config struct {
	Web struct {
		APIHost         string        `envDefault:"0.0.0.0:8181"`
		DebugHost       string        `envDefault:"0.0.0.0:4000"`
		ReadTimeout     time.Duration `envDefault:"60s"`
		WriteTimeout    time.Duration `envDefault:"60s"`
		ShutdownTimeout time.Duration `envDefault:"5s"`
	}
	DB struct {
		User       string `envDefault:"postgres"`
		Password   string `envDefault:"password"`
		Host       string `env:"DB_HOST" envDefault:"0.0.0.0"`
		Name       string `envDefault:"postgres"`
		DisableTLS bool   `envDefault:"true"`
	}
	Auth struct {
		KeyID          string `envDefault:"1"`
		PrivateKeyFile string `envDefault:"/app/private.pem"`
		SecretKey	   string `envDefault:"mySuperSecretKeyLol"`
		Algorithm      string `envDefault:"HS512"`
	}
	Zipkin struct {
		LocalEndpoint string  `envDefault:"0.0.0.0:3000"`
		ReporterURI   string  `envDefault:"http://zipkin:9411/api/v2/spans"`
		ServiceName   string  `envDefault:"users-api"`
		Probability   float64 `envDefault:"0.05"`
	}
}

// Create a new instance of the logger. You can have any number of instances.
var (
	logger = log.New()
	build = "develop"
)

func main() {
	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	//https://medium.com/google-cloud/hidden-super-powers-of-stackdriver-logging-ca110dae7e74
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	logger.Formatter = &log.JSONFormatter{}

	if err := run(); err != nil {
		logger.Error("error :", err)
		os.Exit(1)
	}
}

func run()  error {

	// =========================================================================
	// Configuration
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		return errors.Wrap(err, "parsing config")
	}
	// =========================================================================
	// App Starting

	// Print the build version for our logs. Also expose it under /debug/vars.
	expvar.NewString("build").Set(build)
	logger.Info("main : Started : Application Initializing version ", build)
	defer logger.Info("main : Completed")

	fmt.Printf("main : Config : %+v\n", cfg)

	// =========================================================================
	// Initialize authentication support

	f := auth.NewSimpleKeyLookupFunc(cfg.Auth.KeyID,  []byte(cfg.Auth.SecretKey))
	authenticator, err := auth.NewAuthenticator([]byte(cfg.Auth.SecretKey), cfg.Auth.KeyID, cfg.Auth.Algorithm, f)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}


	// =========================================================================
	// Start Database

	logger.Info("main : Started : Initializing database support")

	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		logger.Info("main : Database Stopping : %s", cfg.DB.Host)
		db.Close()
	}()

	// Register instrumentation callbacks
	//https://github.com/sagikazarmark/go-gin-gorm-opencensus
	ocgorm.RegisterCallbacks(db)

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	//db.AutoMigrate(&domain.User{}, &domain.Authority{}, &domain.Privilege{})

	// =========================================================================
	// Start Tracing Support

	logger.Info("main : Started : Initializing zipkin tracing support")

	localEndpoint, err := openzipkin.NewEndpoint("users-api", cfg.Zipkin.LocalEndpoint)
	if err != nil {
		return err
	}

	reporter := zipkinHTTP.NewReporter(cfg.Zipkin.ReporterURI)
	ze := zipkin.NewExporter(reporter, localEndpoint)

	trace.RegisterExporter(ze)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	/*trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(cfg.Zipkin.Probability),
	})*/

	defer func() {
		log.Printf("main : Tracing Stopping : %s", cfg.Zipkin.LocalEndpoint)
		reporter.Close()
	}()
	// =========================================================================
	//Start Metrics Support
	//https://github.com/zsais/go-gin-prometheus
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Registry: prom.DefaultGatherer.(*prom.Registry),
	})

	if err != nil {
		return err
	}

	// Register stat views
	err = view.Register(
		// Gin (HTTP) stats
		ochttp.ServerRequestCountView,
		ochttp.ServerRequestBytesView,
		ochttp.ServerResponseBytesView,
		ochttp.ServerLatencyView,
		ochttp.ServerRequestCountByMethod,
		ochttp.ServerResponseCountByStatusCode,

		// Gorm stats
		ocgorm.QueryCountView,
	)
	if err != nil {
		return err
	}

	// Register prometheus as a stats exporter
	view.RegisterExporter(exporter)

	// =========================================================================
	// Start API Service
	logger.Info("main : Started : Initializing API support")
	cache := cache.NewRedisCache(cache.RedisOpts{
		Host:       "",
		Expiration: time.Hour,
	})

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	//https://gobyexample.com/signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      delivery.New(shutdown, db, logger /*, exporter*/, authenticator, cache),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}


	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.WithFields(log.Fields{
			"transport": "http",
			"address":    api.Addr,
		}).Info("msg", "listening")
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "starting server")

	case sig := <-shutdown:
		logger.Info("main : %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			logger.Info("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGKILL: // SIGSTOP (linux)
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
