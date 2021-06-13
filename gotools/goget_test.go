package gotools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		path := "github.com/tenntenn/greeting"
		dir := t.TempDir()

		err := GoGet(path, dir)
		require.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		path := "example.com"
		dir := t.TempDir()

		err := GoGet(path, dir)
		require.Error(t, err)
	})
}
