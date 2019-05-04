package app

import (
	"github.com/dohernandez/geolocation-service/pkg/app"
)

// Container contains application resources
type Container struct {
	*app.Container

	cfg Config
}

func newContainer(cfg Config, upstream *app.Container) *Container {
	return &Container{
		Container: upstream,
		cfg:       cfg,
	}
}

// Cfg returns app-level application configuration
// nolint:unused
func (c *Container) Cfg() Config {
	return c.cfg
}
