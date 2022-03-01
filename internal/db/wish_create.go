package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tktip/wishtree-api/pkg/apiutils"
)

// CreateWish populates a wish with data
func (c *Connector) CreateWish(
	tx *sqlx.Tx, wishID int, text string, author *string, zipCode string, categoryID int) (
	newWish Wish, err error,
) {
	err = verifyWishIsFree(tx, wishID)
	if err != nil {
		return
	}

	var queryParams []interface{}
	queryParams = append(queryParams, text, zipCode, categoryID, time.Now())

	query := "UPDATE wish SET text=?, zipCode=?, category_id=?, createdAt=?"
	if author != nil {
		query += ", author=? "
		queryParams = append(queryParams, *author)
	}
	query += " WHERE id=?"
	queryParams = append(queryParams, wishID)

	_, err = tx.Exec(query, queryParams...)

	if err != nil {
		err = fmt.Errorf("failed to perform CreateWish db query >> %w", err)
		return
	}

	var wish *Wish
	wish, err = c.GetWishByID(tx, wishID)
	newWish = *wish

	return
}

func verifyWishIsFree(tx *sqlx.Tx, wishID int) (err error) {
	checkExistingQuery := "SELECT text FROM wish WHERE id = ?"
	var rows *sqlx.Rows
	rows, err = tx.Queryx(checkExistingQuery, wishID)

	if err != nil {
		err = fmt.Errorf("could not perform verifyWishIsFree db query >> %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var wishText *string
		err = rows.Scan(&wishText)
		if err != nil {
			err = fmt.Errorf("could not scan wishText >> %w", err)
			return
		}

		if wishText != nil {
			err = apiutils.ErrWishIDTaken
		}

		return nil
	}

	return apiutils.ErrBadWishQuery
}

// ArchiveAndClearOldestActiveWish creates a copy of the wish with the oldest createdAt date,
// with isArchived=1, and frees the selected wish - and returns the id of this wish.
func ArchiveAndClearOldestActiveWish(tx sqlx.Tx) (archivedWishID int, err error) {
	var row *sqlx.Row
	row = tx.QueryRowx(`
		SELECT TOP 1 id AS wishId, text, author, zipCode, createdAt, x, y
		FROM wish
		WHERE isArchived = 0 AND text IS NOT NULL
		ORDER BY wish.createdAt ASC
	`)

	err = row.Err()
	if err != nil {
		err = fmt.Errorf("Could not perform ArchiveOldestActiveWish query >> %w", err)
		return
	}

	var oldestWish Wish
	err = row.Scan(&oldestWish)
	if err != nil {
		err = fmt.Errorf("Could not scan ArchiveOldestActiveWish result >> %w", err)
		return
	}

	createWishCopyQuery := `INSERT INTO wish (
		x, y, text, author, zipCode, createdAt, category_id, isArchived
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)`

	_, err = tx.Exec(createWishCopyQuery,
		oldestWish.X,
		oldestWish.Y,
		oldestWish.Text,
		oldestWish.Author,
		oldestWish.ZipCode,
		oldestWish.CreatedAt,
		oldestWish.CategoryID,
	)

	err = ClearWishByID(tx, oldestWish.ID)

	return oldestWish.ID, err
}
