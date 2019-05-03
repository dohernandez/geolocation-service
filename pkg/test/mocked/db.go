package mocked

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// DBMock struct to hold necessary dependencies
type DBMock struct {
	db      *sql.DB
	SqlxDB  *sqlx.DB
	Sqlmock sqlmock.Sqlmock
}

// NewDBMock creates a new database mock instance
func NewDBMock(t *testing.T) *DBMock {
	db, mock, err := sqlmock.New()
	if err != nil {
		if t != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		panic(err)
	}

	return &DBMock{
		db:      db,
		SqlxDB:  sqlx.NewDb(db, "sqlmock"),
		Sqlmock: mock,
	}
}

// Close closes the database, releasing any open resources.
func (m *DBMock) Close() {
	// #nosec G104
	// nolint:errcheck
	m.db.Close()
}
