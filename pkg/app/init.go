package app

import (
	"github.com/dohernandez/geolocation-service/pkg/log"
	"github.com/pkg/errors"
)

// NewAppContainer initializes application container
func NewAppContainer(cfg Config) (*Container, error) {
	c := newContainer(cfg)

	logger, err := log.NewLog(cfg.Log)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init logger")
	}
	c.WithLogger(logger)

	var panicCatcher []PanicCatcher
	if cfg.Environment == EnvDev {
		panicCatcher = append(panicCatcher, PanicPrinter)
	} else {
		panicCatcher = append(
			panicCatcher,
			PanicLogger(logger),
			PanicResponse(logger, false),
		)
	}
	c.WithPanicCatcher(panicCatcher)

	return c, nil
}
