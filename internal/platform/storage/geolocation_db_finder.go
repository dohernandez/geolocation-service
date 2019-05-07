package storage

import (
	"context"
	"database/sql"

	"fmt"

	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/pkg/log"
	"github.com/jmoiron/sqlx"
)

// geolocalationDBFinder to hold necessary dependencies to persist the Geolocalation entity into the DB
type geolocalationDBFinder struct {
	db    *sqlx.DB
	table string
}

// NewGeolocalationDBFinder creates a geolocalation db storage instance
func NewGeolocalationDBFinder(db *sqlx.DB, table string) domain.Finder {
	return &geolocalationDBFinder{
		db:    db,
		table: table,
	}
}

// ByIPAddress finds geolocations that contains the given ip
func (p *geolocalationDBFinder) ByIPAddress(ctx context.Context, ip string) (*domain.Geolocation, error) {
	logger := log.FromContext(ctx)

	query := "SELECT * FROM %[1]s WHERE ip_address = $1"
	query = fmt.Sprintf(query, p.table)

	if logger != nil {
		logger.
			WithField("sql", query).
			WithField("ip", ip).
			Debugf("exec in transaction")
	}

	var g domain.Geolocation

	err := p.db.GetContext(ctx, &g, query, ip)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return &g, nil
}
