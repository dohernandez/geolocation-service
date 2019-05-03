package main

import (
	"fmt"
	"os"

	"github.com/dohernandez/geolocation-service/internal/platform/app"
	"github.com/dohernandez/geolocation-service/internal/platform/http"
	"github.com/dohernandez/geolocation-service/pkg/http/server"
	"github.com/dohernandez/geolocation-service/pkg/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.Info().String())

		return
	}

	cfg, err := app.LoadEnv()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	c, err := app.NewAppContainer(cfg)
	if err != nil {
		panic("failed init application container: " + err.Error())
	}

	c.Logger().Debug("Creating routers")
	router := http.NewRouter(c)

	c.Logger().Infof("Starting server at port http://0.0.0.0:%d", cfg.Port)
	(&server.Instance{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
		Logger:  c.Logger(),
		Closer:  c,
	}).Start()
}
