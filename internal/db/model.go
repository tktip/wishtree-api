package db

import "time"

// Wish represents a wish from a person
type Wish struct {
	ID                  int        `json:"id" db:"wishId"`
	X                   float64    `json:"x" db:"x"`
	Y                   float64    `json:"y" db:"y"`
	Text                *string    `json:"text" db:"text"`
	Author              *string    `json:"author" db:"author"`
	ZipCode             *string    `json:"zipCode" db:"zipCode"`
	CreatedAt           *time.Time `json:"createdAt" db:"createdAt"`
	CategoryID          *int       `json:"categoryId" db:"categoryId"`
	CategoryName        *string    `json:"categoryName" db:"categoryName"`
	CategoryDescription *string    `json:"categoryDescription" db:"categoryDescription"`
	IsArchived          bool       `json:"isArchived" db:"isArchived"`
}

// Category represents a category for a wish
type Category struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// TreeWishCounts represents numbers of wishes with various states in the tree
type TreeWishCounts struct {
	ShownTakenWishes int `json:"shownWishes"`
	TotalWishes      int `json:"totalWishes"`
	ArchivedWishes   int `json:"archivedWishes"`
}
