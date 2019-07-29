package web

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/kit/log"
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

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*gin.Engine
	och      *ochttp.Handler
	shutdown chan os.Signal
	logger   kitlog.Logger
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, logger kitlog.Logger, mw ...Middleware) *App {
	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.Use(cors.Default())
	app := App{
		Engine:   engine,
		shutdown: shutdown,
		logger:   logger,
		mw:       mw,
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
	a.logger.Log("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGSTOP
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r, params); err != nil {
			a.logger.Log("*****> critical shutdown error: %v", err)
			a.SignalShutdown()
			return
		}
	}

	// Add this handler for the specified verb and route.
	a.Engine..TreeMux.Handle(verb, path, h)
}


func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.ServeHTTP(w, r)
}