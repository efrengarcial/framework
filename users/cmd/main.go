package main

import (
	"flag"
	"fmt"
	"github.com/efrengarcial/framework/users/pkg/model"
	"github.com/efrengarcial/framework/users/pkg/repository"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/efrengarcial/framework/users/pkg/transport"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort              = "8080"
	defaultRoutingServiceURL = "http://localhost:7878"
	defaultMongoDBURL        = "127.0.0.1"
	defaultDBName            = "dddsample"
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

	service := service.NewService(repo, log.With(logger, "component", "users"))

	srv := transport.New(service, log.With(logger, "component", "http"))

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
