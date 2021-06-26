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

func (err *NotFoundError) Error() string {
	return err.err.Error()
}

func (err *NotFoundError) Unwrap() error {
	return err.err
}

func IsNotFoundError(err error) bool {
	var nErr *NotFoundError
	if errors.As(err, &nErr) {
		return true
	}
	return false
}
