package cmd

import (
	"context"
	"flag"
	"github.com/efrengarcial/framework/users/pkg/service"
	"github.com/go-kit/kit/log"
	"os"
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
		rsurl= envString("ROUTINGSERVICE_URL", defaultRoutingServiceURL)
		dburl= envString("MONGODB_URL", defaultMongoDBURL)
		dbname= envString("DB_NAME", defaultDBName)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		routingServiceURL = flag.String("service.routing", rsurl, "routing service URL")
		mongoDBURL = flag.String("db.url", dburl, "MongoDB URL")
		databaseName = flag.String("db.name", dbname, "MongoDB database name")
		inmemory = flag.Bool("inmem", false, "use in-memory repositories")

		ctx = context.Background()
	)

	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)


	// Setup repositories
	var service service.Repository


}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
