package main

import (
	"flag"
	"fmt"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/repository"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/efrengarcial/framework/users/pkg/transport"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort       = "8080"
	defaultDBHost     = "127.0.0.1"
	defaultDBNme      = "db"
	defaultDBUser     = "postgres"
	defaultDBPassword = "dbpwd"
)

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

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)


	// Creates a database connection and handles
	// closing it again before exit.
	db, err := CreateConnection()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	// Automatically migrates the user struct
	// into database columns/types etc. This will
	// check for changes and migrate them each time
	// this service is restarted.
	db.AutoMigrate(&model.User{}, &model.Authority{}, &model.Privilege{})

	// Setup repositories
	repo := repository.NewGormRepository(db)

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
	userRepository := repository.NewUserGormRepository(db)
	as := service.NewAuthService(userRepository, ts, log.With(logger, "component", "auth"))

	srv := transport.New(us, as, log.With(logger, "component", "http"))

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

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func CreateConnection() (*gorm.DB, error) {

	// Get database details from environment variables
	host := envString("DB_HOST", defaultDBHost)
	user :=  envString("DB_USER", defaultDBUser)
	DBName := envString("DB_NAME", defaultDBNme)
	password := envString("DB_PASSWORD", defaultDBPassword)

	return gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s sslmode=disable password=%s",
			host, user, DBName, password,
		),
	)
}
