package gotools

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		path := "github.com/tenntenn/greeting"
		dir := t.TempDir()

		err := GoModInit(dir)
		require.NoError(t, err)

		err = GoGet(path, dir)
		require.NoError(t, err)

		resp, err := GoList(path, dir)
		require.NoError(t, err)
		assert.Contains(t, resp.Dir, path)
	})

	t.Run("fail", func(t *testing.T) {
		path := "example"
		dir := t.TempDir()

		err := GoModInit(dir)
		require.NoError(t, err)

		_, err = GoList(path, dir)
		require.Error(t, err)
	})
}
