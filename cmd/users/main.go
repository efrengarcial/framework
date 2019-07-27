package main

import (
	"flag"
	"fmt"
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/users/repository"
	"github.com/efrengarcial/framework/internal/users/service"
	"github.com/efrengarcial/framework/internal/users/handlers"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort       = "8282"
	defaultDBHost     = "127.0.0.1"
	defaultDBNme      = "db"
	defaultDBUser     = "postgres"
	defaultDBPassword = "password"
)

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}


func main() {
	var (
		addr= envString("PORT", defaultPort)
		//rsurl= envString("ROUTINGSERVICE_URL", defaultRoutingServiceURL)
		//dburl= envString("MONGODB_URL", defaultMongoDBURL)
		//dbname= envString("DB_NAME", defaultDBName)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		//routingServiceURL = flag.String("service.routing", rsurl, "routing service URL")
		//mongoDBURL = flag.String("db.url", dburl, "MongoDB URL")
		//databaseName = flag.String("db.name", dbname, "MongoDB database name")
		//inmemory = flag.Bool("inmem", false, "use in-memory repositories")
	)

	flag.Parse()

	//https://medium.com/google-cloud/hidden-super-powers-of-stackdriver-logging-ca110dae7e74
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	host := envString("DB_HOST", defaultDBHost)
	user :=  envString("DB_USER", defaultDBUser)
	DBName := envString("DB_NAME", defaultDBNme)
	password := envString("DB_PASSWORD", defaultDBPassword)
	// Creates a database connection and handles
	// closing it again before exit.
	db := database.Initialize("postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s sslmode=disable password=%s",
			host, user, DBName, password))

	defer db.Close()

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&service.User{}, &service.Authority{}, &service.Privilege{})

	// Setup repositories
	repo := repository.NewUserGormRepository(db)

	fieldKeys := []string{"method"}
	us := service.NewService(repo, log.With(logger, "component", "users"))
	us = service.NewInstrumentingService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		us,
	)
	ts := service.NewTokenService()
	as := service.NewAuthService(repo, ts, log.With(logger, "component", "auth"))

	srv := handlers.New(us, as, log.With(logger, "component", "http"))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, srv)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)



}