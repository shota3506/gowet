package gotools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoVet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		path := "github.com/tenntenn/greeting"
		dir := t.TempDir()

		err := GoModInit(dir)
		require.NoError(t, err)

		err = GoGet(path, dir)
		require.NoError(t, err)

		module, err := GoList(path, dir)
		require.NoError(t, err)

		out, err := GoVet(module.Dir)
		require.NoError(t, err)
		require.NotNil(t, out)
	})
}
