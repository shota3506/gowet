package handler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBadRequestError(t *testing.T) {
	testcases := []struct {
		err      error
		expected bool
	}{
		{
			err:      NewBadRequestError(errors.New("some error")),
			expected: true,
		},
		{
			err:      fmt.Errorf("error: %w", NewBadRequestError(errors.New("some error"))),
			expected: true,
		},
		{
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := IsBadRequestError(testcase.err)
		assert.Equal(t, testcase.expected, actual)
	}
}

func TestInternalServerError(t *testing.T) {
	testcases := []struct {
		err      error
		expected bool
	}{
		{
			err:      NewInternalServerError(errors.New("some error")),
			expected: true,
		},
		{
			err:      fmt.Errorf("error: %w", NewInternalServerError(errors.New("some error"))),
			expected: true,
		},
		{
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := IsInternalServerError(testcase.err)
		assert.Equal(t, testcase.expected, actual)
	}
}
