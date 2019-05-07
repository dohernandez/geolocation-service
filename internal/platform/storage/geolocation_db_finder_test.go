package storage_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/dohernandez/geolocation-service/internal/domain"
	"github.com/dohernandez/geolocation-service/internal/platform/storage"
	"github.com/dohernandez/geolocation-service/pkg/test/mocked"
	"github.com/stretchr/testify/assert"
)

func TestFinder(t *testing.T) {
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
		scenario           string
		ip                 string
		setMockExpectation func(ip string, g *domain.Geolocation, mock sqlmock.Sqlmock)
		err                error
	}{
		{
			scenario: "Successfully find goelocation by ip address",
			ip:       "160.103.7.140",
			setMockExpectation: func(ip string, g *domain.Geolocation, mock sqlmock.Sqlmock) {
				// #nosec G201
				mock.ExpectQuery(fmt.Sprintf(`^SELECT .+ FROM %[1]s WHERE .+$`, table)).WithArgs(
					ip,
				).WillReturnError(
					nil,
				).WillReturnRows(sqlmock.NewRows([]string{
					"id",
					"ip_address",
					"country_code",
					"country",
					"city",
					"latitude",
					"longitude",
					"mystery_value",
				}).AddRow(
					g.ID.String(),
					g.IPAddress,
					g.CountryCode,
					g.Country,
					g.City,
					g.Latitude,
					g.Longitude,
					g.MysteryValue,
				))
			},
		},
		{
			scenario: "Find goelocation by ip address fails, error happens",
			ip:       "160.103.7.140",
			setMockExpectation: func(ip string, g *domain.Geolocation, mock sqlmock.Sqlmock) {
				// #nosec G201
				mock.ExpectQuery(fmt.Sprintf(`^SELECT .+ FROM %[1]s WHERE .+$`, table)).WithArgs(
					ip,
				).WillReturnError(
					sql.ErrTxDone,
				)
			},
			err: sql.ErrTxDone,
		},
		{
			scenario: "Find goelocation by ip address fails, not found",
			ip:       "160.103.7.140",
			setMockExpectation: func(ip string, g *domain.Geolocation, mock sqlmock.Sqlmock) {
				// #nosec G201
				mock.ExpectQuery(fmt.Sprintf(`^SELECT .+ FROM %[1]s WHERE .+$`, table)).WithArgs(
					ip,
				).WillReturnError(
					sql.ErrNoRows,
				)
			},
			err: domain.ErrNotFound,
		},
	}

	dbMock := mocked.NewDBMock(t)
	defer dbMock.Close()

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(tc.scenario, func(t *testing.T) {
			tc.setMockExpectation(tc.ip, &g, dbMock.Sqlmock)

			f := storage.NewGeolocalationDBFinder(dbMock.SqlxDB, table)
			gf, err := f.ByIPAddress(context.TODO(), tc.ip)

			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, &g, gf)
			}

			err = dbMock.Sqlmock.ExpectationsWereMet()
			assert.NoErrorf(t, err, "there were unfulfilled expectations: %s", err)
		})
	}
}
