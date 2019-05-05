package storage_test

import (
	"testing"

	"fmt"

	"context"

	"database/sql"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/internal/platform/storage"
	"github.com/dohernandez/geolocation-service/pkg/test/mocked"
	"github.com/stretchr/testify/assert"
)

const (
	table = "table"
)

func TestPersist(t *testing.T) {
	g := domain.Geolocation{
		IPAddress:    "160.103.7.140",
		CountryCode:  "CZ",
		Country:      "Nicaragua",
		City:         "New Neva",
		Latitude:     "-68.31023296602508",
		Longitude:    "-37.62435199624531",
		MysteryValue: 7301823115,
	}

	testCases := []struct {
		scenario string

		setMockExpectation func(g *domain.Geolocation, mock sqlmock.Sqlmock)
		err                error
	}{
		{
			scenario: "Successfully persist",
			setMockExpectation: func(g *domain.Geolocation, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				// #nosec G201
				mock.ExpectExec(fmt.Sprintf(`^INSERT INTO %[1]s \(.+\) VALUES \(.+\)$`, table)).WithArgs(
					g.ID,
					g.IPAddress,
					g.CountryCode,
					g.Country,
					g.City,
					g.Latitude,
					g.Longitude,
					g.MysteryValue,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			scenario: "Persist failed",
			setMockExpectation: func(g *domain.Geolocation, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				// #nosec G201
				mock.ExpectExec(fmt.Sprintf(`^INSERT INTO %[1]s \(.+\) VALUES \(.+\)$`, table)).WithArgs(
					g.ID,
					g.IPAddress,
					g.CountryCode,
					g.Country,
					g.City,
					g.Latitude,
					g.Longitude,
					g.MysteryValue,
				).WillReturnError(sql.ErrTxDone)

				mock.ExpectRollback()
			},
			err: sql.ErrTxDone,
		},
	}

	dbMock := mocked.NewDBMock(t)
	defer dbMock.Close()

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.scenario, func(t *testing.T) {
			tc.setMockExpectation(&g, dbMock.Sqlmock)

			p := storage.NewGeolocalationDBPersister(dbMock.SqlxDB, table)
			err := p.Persist(context.TODO(), &g)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}

			err = dbMock.Sqlmock.ExpectationsWereMet()
			assert.NoErrorf(t, err, "there were unfulfilled expectations: %s", err)
		})
	}
}
