package storage

import "errors"

var (
	// ErrNotFound not found entity in Storage or ValuesStorage
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists entity is already exist in Storage
	ErrAlreadyExists = errors.New("already exists")
)
