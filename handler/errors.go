package handler

import (
	"errors"
	"fmt"
)

type BadRequestError struct {
	err error
}

func NewBadRequestError(err error) error {
	return &BadRequestError{err: err}
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("bad request: %s", e.err.Error())
}

func (e *BadRequestError) Unwrap() error {
	return e.err
}

func IsBadRequestError(err error) bool {
	var nErr *BadRequestError
	if errors.As(err, &nErr) {
		return true
	}
	return false
}

type InternalServerError struct {
	err error
}

func NewInternalServerError(err error) error {
	return &InternalServerError{err: err}
}

func (e *InternalServerError) Error() string {
	return fmt.Sprintf("internal server error: %s", e.err.Error())
}

func (e *InternalServerError) Unwrap() error {
	return e.err
}

func IsInternalServerError(err error) bool {
	var nErr *InternalServerError
	if errors.As(err, &nErr) {
		return true
	}
	return false
}
