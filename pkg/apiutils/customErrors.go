package apiutils

import "errors"

var (
	// ErrWishIDTaken - when someone tries to claim a wish whose ID is taken
	ErrWishIDTaken = errors.New("the wish with this ID is already taken")
	// ErrBadWishQuery - when the wish id does not exist
	ErrBadWishQuery = errors.New("there is no wish with this ID")
)
