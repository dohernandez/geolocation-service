package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dohernandez/geolocation-service/internal/platform/app"
	"github.com/dohernandez/geolocation-service/pkg/version"
	"github.com/urfave/cli"
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

	app := cli.NewApp()
	app.Version = version.Info().Version
	app.Name = cfg.CliImport
	app.Usage = "To import data from a csv file."
	app.UsageText = fmt.Sprintf("%s [arguments]", cfg.CliImport)

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
