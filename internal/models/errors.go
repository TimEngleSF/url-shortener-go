package models

import "errors"

var (
	UniqueViolationCode = "23505"
	ErrNoRecord         = errors.New("models: no matching record found")
	ErrUniqueViolation  = errors.New("models: unique constraint violation")
	ErrInsertRow        = errors.New("models: error inserting row")
	ErrEmptySuffix      = errors.New("models: empty suffix")
)
