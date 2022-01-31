package db

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	// driver import
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/jmoiron/sqlx"
)

// Connector - provides db functionality
type Connector struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	conn     *sqlx.DB
}

func (c *Connector) getDbConn() (*sqlx.DB, error) {
	if c.conn == nil {
		connString := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",
			c.Username, strings.TrimSpace(url.QueryEscape(c.Password)), c.URL, c.DBName)
		log.Debugf("Opening new connection to database '%s", connString)
		var err error
		c.conn, err = sqlx.Connect("mssql", connString)
		if err != nil {
			return nil, err
		}
		log.Debug("Database connection successfully opened.")
	}

	return c.conn, nil
}

// NewTransaction creates a db transaction object
func (c *Connector) NewTransaction() (tx *sqlx.Tx, err error) {
	var conn *sqlx.DB
	conn, err = c.getDbConn()
	if err != nil {
		err = fmt.Errorf("failed to open db connection: %w", err)
		return nil, err
	}

	return conn.BeginTxx(context.Background(), nil)
}

// TxRollbackIfErr rolls back the changes in the transaction if error becomes != nil
func TxRollbackIfErr(tx *sqlx.Tx, err *error) {
	if err == nil {
		panic("txRollbackIfErr used with non-pointer error")
	}

	if *err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			log.Errorf("transaction failed, and so did rollback: %v", rbErr)
		}
	}
}
