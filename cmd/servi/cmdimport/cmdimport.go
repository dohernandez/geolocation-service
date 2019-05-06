package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/internal/platform/app"
	"github.com/dohernandez/geolocation-service/pkg/version"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.Info().String())

		return
	}

	ctx, cancelCtx := context.WithCancel(context.TODO())
	defer cancelCtx()

	cfg, err := app.LoadEnv()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	c, err := app.NewAppContainer(cfg)
	if err != nil {
		panic("failed init application container: " + err.Error())
	}

	app := cli.NewApp()
	app.Version = version.Info().Version
	app.Name = cfg.CliImport

	app.Usage = "To import data from a csv file."
	app.UsageText = fmt.Sprintf("%s [arguments]", cfg.CliImport)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "csv file",
		},
	}

	app.Action = func(ctxCli *cli.Context) error {
		fp := ctxCli.String("file")

		if fp == "" {
			return errors.New("file must be defined")
		}

		f, err := os.Open(filepath.Clean(fp))
		if err != nil {
			return err
		}
		// nolint:errcheck
		defer f.Close()

		// create context with logger
		//lCtx := logger.ToContext(ctx, c.Logger())

		uc := domain.NewImportGeolocationFromCSVFileToDBUseCase(c.GeolocationPersister())

		start := time.Now()
		processed, accepted, discarded, err := uc.Do(ctx, f)
		elapsed := time.Since(start)

		if err != nil {
			return err
		}

		fmt.Println("Import statistics")
		fmt.Printf("time elapsed: %s\n", elapsed)
		fmt.Printf("processed: %d\n", processed)
		fmt.Printf("accepted: %d\n", accepted)
		fmt.Printf("discarded: %d\n", discarded)

		return nil
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
