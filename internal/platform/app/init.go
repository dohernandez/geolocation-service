package app

import "github.com/dohernandez/geolocation-service/pkg/app"

// NewAppContainer initializes application container
func NewAppContainer(cfg Config) (*Container, error) {
	bc, err := app.NewAppContainer(cfg.Config)
	if err != nil {
		return nil, err
	}

	c := newContainer(cfg, bc)

	return c, nil
}
