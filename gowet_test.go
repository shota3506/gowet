package gowet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		path := "github.com/tenntenn/greeting"

		out, err := Run(ctx, path)
		require.NoError(t, err)
		require.NotNil(t, out)
	})

	t.Run("fail", func(t *testing.T) {
		ctx := context.Background()
		path := "example.com"

		out, err := Run(ctx, path)
		require.Error(t, err)
		require.Nil(t, out)
	})
}
