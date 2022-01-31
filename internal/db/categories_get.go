package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// GetAllCategories gets all categories
func (c *Connector) GetAllCategories() (categories []Category, err error) {
	var conn *sqlx.DB
	conn, err = c.getDbConn()
	if err != nil {
		return
	}

	query := `SELECT id, name, description FROM category`

	var rows *sqlx.Rows
	rows, err = conn.Queryx(query)
	if err != nil {
		err = fmt.Errorf("could not perform GetAllCategories db query >> %w", err)
		return
	}
	defer rows.Close()

	categories = []Category{}
	for rows.Next() {
		var category Category
		err = rows.StructScan(&category)
		if err != nil {
			err = fmt.Errorf("error scanning row in GetAllCategories >> %w", err)
			return
		}

		categories = append(categories, category)
	}

	return
}
