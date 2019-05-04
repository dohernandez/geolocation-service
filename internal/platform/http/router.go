package http

import (
	"github.com/dohernandez/geolocation-service/internal/platform/app"
	"github.com/dohernandez/geolocation-service/pkg/http/router"
	"github.com/go-chi/chi"
)

// NewRouter register the roles to the service
func NewRouter(c *app.Container) chi.Router {
	r := router.NewRouter(c.Container)

	return r
}
