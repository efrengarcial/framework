package web

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
	"net/http"
	"os"
	"syscall"
	"time"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*gin.Engine
	och      *ochttp.Handler
	shutdown chan os.Signal
	logger   *logrus.Logger
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, logger *logrus.Logger) *App {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(cors.Default())
	engine.Use(Trace())
	app := App{
		Engine:   engine,
		shutdown: shutdown,
		logger:   logger,
	}

	// Create an OpenCensus HTTP Handler which wraps the router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	app.och = &ochttp.Handler{
		Handler:     app.Engine,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &app
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.logger.Error("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGKILL
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := trace.StartSpan(c.Request.Context(), "internal.platform.web")
		defer span.End()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Now:     time.Now(),
		}
		context.WithValue(ctx, KeyValues, &v)
	}
}
