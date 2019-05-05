package bootstrap

import (
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/dohernandez/geolocation-service/features/bootstrap/internal"
	"github.com/dohernandez/geolocation-service/internal/platform/app"
	"github.com/dohernandez/geolocation-service/internal/platform/http"
	"github.com/dohernandez/geolocation-service/pkg/http/server"
	"github.com/dohernandez/geolocation-service/pkg/test/feature"
)

func TestIntegration(t *testing.T) {
	if !feature.RunGoDogTests {
		t.Skip("Skipping integration tests")
	}

	cfg, err := app.LoadEnv()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	c, err := app.NewAppContainer(cfg)
	if err != nil {
		panic("failed init application container: " + err.Error())
	}

	router := http.NewRouter(c)

	srv := &server.Instance{
		Handler:      router,
		AddrAssigned: make(chan string),
		Logger:       c.Logger(),
	}

	go srv.Start()

	baseURL := <-srv.AddrAssigned

	feature.RunSuite("..", func(_ *testing.T, s *godog.Suite) {
		feature.RegisterRestContext(s, baseURL)
		internal.RegisterCommandContext(s)
		internal.RegisterDBContext(s, c.DB())

		RegisterFileContext(s)
	}, t)
}
