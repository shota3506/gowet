//go:generate mockgen -package=$GOPACKAGE -destination=mock_$GOFILE . DB
package database

import (
	"context"
	"errors"
)

type DB interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key, value string) error
}

type NotFoundError struct {
	err error
}

func NewNotFoundError(err error) error {
	return &NotFoundError{err: err}
}

func (e *NotFoundError) Error() string {
	return e.err.Error()
}

func (e *NotFoundError) Unwrap() error {
	return e.err
}

func IsNotFoundError(err error) bool {
	var nErr *NotFoundError
	if errors.As(err, &nErr) {
		return true
	}
	return false
}
