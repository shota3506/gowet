package database

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotFonundError(t *testing.T) {
	testcases := []struct {
		err      error
		expected bool
	}{
		{
			err:      NewNotFoundError(errors.New("some error")),
			expected: true,
		},
		{
			err:      fmt.Errorf("error: %w", NewNotFoundError(errors.New("some error"))),
			expected: true,
		},
		{
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, testcase := range testcases {
		actual := IsNotFoundError(testcase.err)
		assert.Equal(t, testcase.expected, actual)
	}
}
