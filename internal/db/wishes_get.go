package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// GetWishByID gets a wish by its id, returns nil if the wish is not found
func (c *Connector) GetWishByID(tx *sqlx.Tx, wishID int) (wish *Wish, err error) {
	var conn *sqlx.DB
	if tx == nil {
		conn, err = c.getDbConn()
		if err != nil {
			return
		}
	}

	query := `SELECT 
			wish.id AS wishId, text, author, zipCode, createdAt, x, y,
			category.id AS categoryId, category.name AS categoryName,
			category.description AS categoryDescription
	  FROM wish INNER JOIN category ON (category.id = wish.category_id)
		WHERE wish.id = ?`

	var rows *sqlx.Rows

	if tx == nil {
		rows, err = conn.Queryx(query, wishID)
	} else {
		rows, err = tx.Queryx(query, wishID)
	}

	if err != nil {
		err = fmt.Errorf("could not perform GetWishByID db query >> %w", err)
		return
	}
	defer rows.Close()

	var foundWish Wish
	for rows.Next() {
		err = rows.StructScan(&foundWish)
		if err != nil {
			err = fmt.Errorf("error scanning wish row >> %w", err)
		}

		wish = &foundWish
		return
	}

	return nil, nil
}

// GetRandomFreeWish gets a wish that hasn't been taken by a user
func GetRandomFreeWish(tx *sqlx.Tx) (wish Wish, err error) {
	query := `SELECT TOP 1
			wish.id AS wishId, text, author, zipCode, createdAt, x, y,
			category.id AS categoryId, category.name AS categoryName,
			category.description AS categoryDescription
		FROM wish	LEFT JOIN category ON (category.id = wish.category_id)
		WHERE text IS NULL AND isArchived=0
		ORDER BY wish.id DESC`

	var rows *sqlx.Rows
	rows, err = tx.Queryx(query)
	if err != nil {
		err = fmt.Errorf("could not perform GetRandomFreeWish db query >> %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&wish)
		if err != nil {
			err = fmt.Errorf("error scanning wish row in GetRandomFreeWish >> %w", err)
			return
		}

		return
	}

	err = errors.New("could not find any free wish")
	return
}

// GetAllWishes gets all existing wishes
func (c *Connector) GetAllWishes() (wishes []Wish, err error) {
	wishes = []Wish{}

	var conn *sqlx.DB
	conn, err = c.getDbConn()
	if err != nil {
		return
	}

	query := `SELECT 
			wish.id AS wishId, text, author, zipCode, createdAt, x, y,
			category.id AS categoryId, category.name AS categoryName,
			category.description AS categoryDescription
		FROM wish	LEFT JOIN category ON (category.id = wish.category_id) 
		WHERE isArchived = 0 AND isHidden = 0
		ORDER BY wish.createdAt DESC, wish.id DESC`

	var rows *sqlx.Rows
	rows, err = conn.Queryx(query)
	if err != nil {
		err = fmt.Errorf("could not perform GetAllWishes db query >> %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var wish Wish
		err = rows.StructScan(&wish)
		if err != nil {
			err = fmt.Errorf("error scanning wish row in GetAllWishes >> %w", err)
			return
		}

		wishes = append(wishes, wish)
	}

	return
}

// GetNumberOfTakenWishes gets number of taken wishes not archived
func GetNumberOfTakenWishes(tx sqlx.Tx) (numberOfWishes int, err error) {
	var row *sqlx.Row
	row = tx.QueryRowx(`
		SELECT COUNT(*)
		FROM wish
		WHERE isArchived = 0 AND text IS NOT NULL AND isHidden = 0
	`)

	err = row.Err()
	if err != nil {
		err = fmt.Errorf("could not perform GetNumberOfTakenWishes query >> %w", err)
		return
	}

	err = row.Scan(&numberOfWishes)
	if err != nil {
		err = fmt.Errorf("could not scan count result >> %w", err)
		return
	}

	return
}

// GetAllTreeWishCounts gets the numbers of shown wishes, total wishes, and archived wishes
func (c *Connector) GetAllTreeWishCounts() (wishCounts TreeWishCounts, err error) {
	var conn *sqlx.DB
	conn, err = c.getDbConn()
	if err != nil {
		return
	}

	shownTakenWishesQuery := `SELECT COUNT(*) FROM wish
		WHERE isArchived = 0 AND text IS NOT NULL AND isHidden = 0`
	totalWishesQuery := `SELECT COUNT(*) FROM wish
		WHERE isArchived = 0 AND isHidden = 0`
	archivedWishesQuery := `SELECT COUNT(*) FROM wish
		WHERE isArchived = 1`

	shownTakenWishes, err := queryCount(*conn, shownTakenWishesQuery)
	if err != nil {
		return
	}
	totalWishes, err := queryCount(*conn, totalWishesQuery)
	if err != nil {
		return
	}
	archivedWishes, err := queryCount(*conn, archivedWishesQuery)
	if err != nil {
		return
	}

	wishCounts = TreeWishCounts{
		ShownTakenWishes: shownTakenWishes,
		TotalWishes:      totalWishes,
		ArchivedWishes:   archivedWishes,
	}

	return
}

func queryCount(conn sqlx.DB, query string) (result int, err error) {
	var row *sqlx.Row
	row = conn.QueryRowx(query)

	err = row.Err()
	if err != nil {
		err = fmt.Errorf("could not perform row count query >> %w", err)
		return
	}

	err = row.Scan(&result)
	if err != nil {
		err = fmt.Errorf("could not scan row count result >> %w", err)
		return
	}

	return
}
