package internal

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/dohernandez/geolocation-service/pkg/test/feature"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// DBContext struct to hold necessary dependencies
type DBContext struct {
	feature.DBContext
}

// RegisterDBContext is the place to truncate database, run background givens and check then should be stored
func RegisterDBContext(s *godog.Suite, db *sqlx.DB) *DBContext {
	c := DBContext{
		DBContext: feature.DBContext{
			DB: db,
			Tables: []string{
				"geolocation",
			},
		},
	}

	s.BeforeScenario(func(_ interface{}) {
		c.CleanUpDB()
	})

	s.Step(`^the following geolocation\(s\) should be stored in the table "([^"]*)"$`, c.theFollowingGeolocationsShouldBeStoredInTheTable)
	s.Step(`^that the following geolocation\(s\) are stored in the table "([^"]*)"$`, c.thatTheFollowingGeolocationsAreStoredInTheTable)
	s.Step(`^there should be "([^"]*)" geolocation\(s\) stored in the table "([^"]*)"$`, c.thereShouldBeGeolocationsStoredInTheTable)

	return &c
}

func (c *DBContext) theFollowingGeolocationsShouldBeStoredInTheTable(table string, data *gherkin.DataTable) error {
	_, err := c.RunExistData("id", table, data, func(col, v string, position int) (string, interface{}) {
		return fmt.Sprintf("%s = $%d", col, position), v
	}, nil)

	if err != nil {
		return err
	}

	return nil
}

func (c *DBContext) thatTheFollowingGeolocationsAreStoredInTheTable(table string, data *gherkin.DataTable) error {
	err := c.RunStoreData(table, data, nil)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	return nil
}

func (c *DBContext) thereShouldBeGeolocationsStoredInTheTable(amount int, table string) error {
	var a int

	// #nosec G201
	query := fmt.Sprintf("SELECT count(id) FROM %s", table)

	err := c.DB.Get(&a, query)
	if err != nil {
		return errors.Wrapf(err, "query [%s]", query)
	}

	if a != amount {
		return errors.Errorf("there should be %d, but there is %d", amount, a)
	}

	return nil
}
