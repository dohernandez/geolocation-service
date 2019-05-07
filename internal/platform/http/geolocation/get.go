package geolocation

import (
	"net/http"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/pkg/http/rest"
	"github.com/dohernandez/geolocation-service/pkg/log"
	"github.com/dohernandez/geolocation-service/pkg/must"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// NewGetIpAddressHandler creates Handler
func NewGetIpAddressHandler(c interface {
	GeolocationFinder() domain.Finder
	Logger() logrus.FieldLogger
}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ipAddress := chi.URLParam(r, "ip_address")
		if ipAddress == "" {
			must.NotFail(render.Render(w, r, rest.ErrBadRequest(errors.New("Ip address missing"))))

			return
		}

		ctx = log.ToContext(ctx, c.Logger())
		g, err := c.GeolocationFinder().ByIpAddress(ctx, ipAddress)
		if err != nil {
			if err == domain.ErrNotFound {
				must.NotFail(render.Render(w, r, rest.ErrNotFound(err)))

				return
			}

			must.NotFail(render.Render(w, r, rest.ErrInternal(err)))

			return
		}

		render.JSON(w, r, g)
	}
}
