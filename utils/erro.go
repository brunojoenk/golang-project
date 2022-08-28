package utils

import "github.com/pkg/errors"

var (
	ErrAuthorIdNotFound = errors.New("Author ID not found")
	ErrBookIdNotFound   = errors.New("Book ID not found")
)
