package feature

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ValueWhereBuilder is a builder function used to define condition in select
//
// Example:
// fmt.Printf("%v, %v", ValueWhereBuilder("id", 12, 1))
// Output:
// {"id = $1", 12}
type ValueWhereBuilder = func(col, v string, position int) (string, interface{})

// DBContext struct to hold necessary dependencies
type DBContext struct {
	DB     *sqlx.DB
	Tables []string
}

// CleanUpDB cleans up the tables
func (c *DBContext) CleanUpDB() {
	for _, table := range c.Tables {
		// #nosec G201
		_, err := c.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			log.Fatal(fmt.Sprintf("here %s", err))
		}
	}
}

// RunExistData check data exist thro select from db table the values from the DataTable.
// The DataTable could expand dynamically as many column contain the table, resulting in some cases a DataTable
// with 3 columns or 5 columns depends on the background data require
func (c *DBContext) RunExistData(selectCol, table string, data *gherkin.DataTable, valueWhereBuilder ValueWhereBuilder, skipColumnsFromCondition []string) (ids []string, err error) {
	columns, err := c.columns(data)
	if err != nil {
		return nil, err
	}

	for _, row := range data.Rows[1:] {
		var id string

		query, args, err := c.buildSelectIDWhereQuery(selectCol, table, columns, row.Cells, valueWhereBuilder, skipColumnsFromCondition)
		if err != nil {
			return nil, err
		}

		err = c.DB.Get(&id, query, args...)
		if err != nil {
			return nil, errors.Wrapf(err, "query [%s] with args [%+v]", query, args)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (c *DBContext) columns(data *gherkin.DataTable) ([]string, error) {
	var tColumns = make([]string, len(data.Rows[0].Cells))

	for k, cell := range data.Rows[0].Cells {
		tColumns[k] = cell.Value
	}

	if len(tColumns) == 0 {
		return nil, fmt.Errorf("there is no column defined")
	}

	return tColumns, nil
}

func (c *DBContext) buildSelectIDWhereQuery(selectCol, table string, columns []string, cells []*gherkin.TableCell, valueWhereBuilder ValueWhereBuilder, skipColumnsFromCondition []string) (query string, args []interface{}, err error) {

	// #nosec G201
	query = fmt.Sprintf("SELECT %s FROM %s", selectCol, table)

	where := make([]string, 0)
	var i int

	for k, cell := range cells {
		col := columns[k]

		if skipColumnsFromCondition != nil {
			var skipCol bool

			for i := range skipColumnsFromCondition {
				if skipColumnsFromCondition[i] == col {
					skipCol = true

					break
				}
			}

			if skipCol {
				continue
			}
		}

		condition, value := valueWhereBuilder(col, cell.Value, i+1)
		where = append(where, condition)

		if value == nil {
			continue
		}

		args = append(args, value)

		i++
	}

	if len(where) == 0 {
		return "", nil, fmt.Errorf("no criteria found")
	}

	// #nosec G201
	query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))

	return query, args, nil
}

// RunNoExistData check data does not exist thro select from db table the values from the DataTable.
// The DataTable could expand dynamically as many column contain the table, resulting in some cases a DataTable
// with 3 columns or 5 columns depends on the background data require
//
// Returns (false, error) if data exists, otherwise return (true, nil)
func (c *DBContext) RunNoExistData(selectCol, table string, data *gherkin.DataTable, valueWhereBuilder ValueWhereBuilder, skipColumnsFromCondition []string) (ok bool, err error) {
	columns, err := c.columns(data)
	if err != nil {
		return false, err
	}

	for _, row := range data.Rows[1:] {
		var id string

		query, args, err := c.buildSelectIDWhereQuery(selectCol, table, columns, row.Cells, valueWhereBuilder, skipColumnsFromCondition)
		if err != nil {
			return false, err
		}

		err = c.DB.Get(&id, query, args...)
		if err != nil {
			if err != sql.ErrNoRows {
				return false, errors.Wrapf(err, "query [%s] with args [%+v]", query, args)
			}

			return true, nil
		}
	}

	return false, errors.Errorf("data exists")
}

// RunStoreData inserts into the db table the values from the DataTable.
// The DataTable could expand dynamically as many column contain the table, resulting in some cases a DataTable
// with 3 columns or 5 columns depends on the background data require
func (c *DBContext) RunStoreData(table string, data *gherkin.DataTable, value func(col, v string) interface{}) error {
	tColumns := make([]string, 0, len(data.Rows[0].Cells))
	qValue := make([]string, 0, len(data.Rows[0].Cells))

	for i, cell := range data.Rows[0].Cells {
		tColumns = append(tColumns, cell.Value)
		qValue = append(qValue, fmt.Sprintf("$%d", i+1))
	}

	if len(tColumns) == 0 {
		return fmt.Errorf("there is no column defined INSERT INTO wallet")
	}

	// #nosec G201
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		table,
		strings.Join(tColumns, ", "),
		strings.Join(qValue, ", "),
	)

	for _, row := range data.Rows[1:] {
		var args []interface{}

		for k, cell := range row.Cells {
			if value != nil {
				args = append(args, value(data.Rows[0].Cells[k].Value, cell.Value))

				continue
			}

			args = append(args, cell.Value)
		}

		_, err := c.DB.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
