package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// GetIsTreeOpen returns true if the tree is open, false otherwise
func (c *Connector) GetIsTreeOpen() (isOpen bool, err error) {
	var conn *sqlx.DB
	conn, err = c.getDbConn()
	if err != nil {
		return
	}

	var row *sqlx.Row
	row = conn.QueryRowx(`SELECT isOpen FROM tree_status`)

	err = row.Err()
	if err != nil {
		err = fmt.Errorf("could not perform GetTreeStatus query >> %w", err)
		return
	}

	var isOpenInt int
	err = row.Scan(&isOpenInt)
	if err != nil {
		err = fmt.Errorf("could not scan tree status result >> %w", err)
		return
	}

	return isOpenInt == 1, nil
}

// UpdateTreeStatus updates the tree's isOpen row
func (c *Connector) UpdateTreeStatus(isOpen bool) (err error) {
	var conn *sqlx.DB
	conn, err = c.getDbConn()
	if err != nil {
		return
	}

	_, err = conn.Exec(`UPDATE tree_status SET isOpen = ?`, isOpen)

	if err != nil {
		err = fmt.Errorf("could not perform UpdateTreeStatus query >> %w", err)
		return
	}

	return
}
