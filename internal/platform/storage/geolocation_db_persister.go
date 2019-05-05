package storage

import (
	"context"
	"fmt"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/pkg/log"
	"github.com/jmoiron/sqlx"
)

// geolocalationDBPersister to hold necessary dependencies to persist the Geolocalation entity into the DB
type geolocalationDBPersister struct {
	db    *sqlx.DB
	table string
}

// NewGeolocalationDBPersister creates a geolocalation db storage instance
func NewGeolocalationDBPersister(db *sqlx.DB, table string) domain.Persister {
	return &geolocalationDBPersister{
		db:    db,
		table: table,
	}
}

func (p *geolocalationDBPersister) Persist(ctx context.Context, g *domain.Geolocation) error {
	logger := log.FromContext(ctx)

	query := `INSERT INTO %[1]s (id, ip_address, country_code, country, city, latitude, longitude, mystery_value) 
					VALUES (:id, :ip_address, :country_code, :country, :city, :latitude, :longitude, :mystery_value)`
	query = fmt.Sprintf(query, p.table)

	if logger != nil {
		logger.Debugf("exec in transaction sql %s, values %+v", query, g)
	}

	return execInTransaction(p.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExecContext(ctx, query, g)
		if err != nil {
			return err
		}

		return nil
	})
}
