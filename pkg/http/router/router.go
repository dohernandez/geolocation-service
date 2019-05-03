package router

import (
	"net/http"

	"github.com/dohernandez/geolocation-service/pkg/app"
	httpHandler "github.com/dohernandez/geolocation-service/pkg/http/router/handler"
	"github.com/dohernandez/geolocation-service/pkg/version"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	health "github.com/hellofresh/health-go"
)

// NewRouter creates base instrumented router
func NewRouter(c *app.Container) chi.Router {
	var r chi.Router = chi.NewRouter()

	AddMiddleWares(r, c)
	AddStandardHandlers(r, c)

	return r
}

// AddMiddleWares register middleware
func AddMiddleWares(r chi.Router, c *app.Container) {
	r.Use(app.PanicRecoverer(c.PanicCatcher()))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(chiMiddleware.RealIP)
}

// AddStandardHandlers registers standard handlers
func AddStandardHandlers(r chi.Router, c *app.Container) {
	logger := c.Logger()

	// HelloWorld
	r.Get("/", func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Set("content-type", "text/html")
		_, err := rw.Write([]byte("Welcome to " + c.Cfg().ServiceName +
			`. Please read API <a href="/docs/api.html">documentation</a>.`))
		if err != nil {
			logger.WithError(err).Error("failed to write response")
		}
	})
	logger.Debug("added `/` route")

	// Endpoint shows the version of the api
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, version.Info())
	})
	logger.Debug("added `/version` route")

	// Endpoint check the health of the api
	r.Get("/status", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	logger.Debug("added `/status` route")

	// Endpoint shows current status
	r.Method(http.MethodGet, "/health", health.Handler())
	logger.Debug("added `/health` route")

	// Endpoint shows API documentation
	r.Method(http.MethodGet, "/docs/*", httpHandler.NewDocsHandler("/docs", "/resources/docs"))
	logger.Debug("added `/docs` route")
}
