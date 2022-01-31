package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ClearWishByID clears the fields of a wish, making it available to others again
func ClearWishByID(tx sqlx.Tx, wishID int) (err error) {
	query := `UPDATE wish SET 
			text = NULL,
			zipCode = NULL,
			author = NULL,
			category_id = NULL,
			createdAt = NULL,
			isPhysical = 0,
			isHidden = 0,
			isArchived = 0
		WHERE id = ?`

	_, err = tx.Exec(query, wishID)

	if err != nil {
		err = fmt.Errorf("failed to perform ClearWishByID db query >> %w", err)
	}

	return
}
