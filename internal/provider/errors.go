package provider

import "errors"

var (
	// ErrNotFound not found entity
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists entity is already exist
	ErrAlreadyExists = errors.New("already exists")

	// ErrNotValid validation error
	ErrNotValid = errors.New("not valid")
)
