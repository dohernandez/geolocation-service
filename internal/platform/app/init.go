package app

import (
	"github.com/dohernandez/geolocation-service/internal/platform/storage"
	"github.com/dohernandez/geolocation-service/pkg/app"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // enable postgres driver
)

const geolocationTable = "geolocation"

// NewAppContainer initializes application container
func NewAppContainer(cfg Config) (*Container, error) {
	bc, err := app.NewAppContainer(cfg.Config)
	if err != nil {
		return nil, err
	}

	c := newContainer(cfg, bc)

	// Init db
	db, err := sqlx.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	c.WithDB(db)

	geolocationPersister := storage.NewGeolocalationDBPersister(db, geolocationTable)
	c.WithGeolocationPersister(geolocationPersister)

	geolocationFinder := storage.NewGeolocalationDBFinder(db, geolocationTable)
	c.WithGeolocationFinder(geolocationFinder)

	return c, nil
}
