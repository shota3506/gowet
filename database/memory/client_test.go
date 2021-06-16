package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	c := NewClient()

	err := c.Set(ctx, "key1", "value1")
	require.NoError(t, err)

	err = c.Set(ctx, "key2", "value2")
	require.NoError(t, err)

	res, err := c.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "value1", string(res))

	res, err = c.Get(ctx, "key2")
	require.NoError(t, err)
	assert.Equal(t, "value2", string(res))

	res, err = c.Get(ctx, "key3")
	require.Error(t, err)
	require.Nil(t, res)

	err = c.Set(ctx, "key1", "new")
	require.NoError(t, err)

	res, err = c.Get(ctx, "key1")
	require.NoError(t, err)
	assert.Equal(t, "new", string(res))
}
