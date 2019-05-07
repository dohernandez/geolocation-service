package app

import (
	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/pkg/app"
	"github.com/jmoiron/sqlx"
)

// Container contains application resources
type Container struct {
	*app.Container

	cfg Config
	db  *sqlx.DB

	geolocationPersister domain.Persister
	geolocationFinder    domain.Finder
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

// WithDB sets sqlx.DB instance
func (c *Container) WithDB(db *sqlx.DB) *Container {
	c.db = db

	return c
}

// DB returns app-level sqlx.DB  instance
// nolint:unused
func (c *Container) DB() *sqlx.DB {
	return c.db
}

// WithGeolocationPersister sets domain.Persister instance
func (c *Container) WithGeolocationPersister(geolocationPersister domain.Persister) {
	c.geolocationPersister = geolocationPersister
}

// GeolocationPersister returns service-level domain.Persister instance
func (c *Container) GeolocationPersister() domain.Persister {
	return c.geolocationPersister
}

// WithGeolocationFinder sets domain.Finder instance
func (c *Container) WithGeolocationFinder(geolocationFinder domain.Finder) {
	c.geolocationFinder = geolocationFinder
}

// GeolocationFinder returns service-level domain.Finder instance
func (c *Container) GeolocationFinder() domain.Finder {
	return c.geolocationFinder
}
