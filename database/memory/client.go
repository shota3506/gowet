package memory

import (
	"context"
	"errors"
	"sync"
)

var ErrNil = errors.New("redigo: nil returned")

type memoryDB map[string]string

// singleton pattern
var md memoryDB = memoryDB{}

type client struct {
	m *sync.RWMutex
}

func NewClient() *client {
	return &client{
		m: &sync.RWMutex{},
	}
}

func (r *client) Get(ctx context.Context, key string) ([]byte, error) {
	r.m.RLock()
	defer r.m.RUnlock()

	value, ok := md[key]
	if !ok {
		return nil, ErrNil
	}
	return []byte(value), nil
}

func (r *client) Set(ctx context.Context, key, value string) error {
	r.m.Lock()
	defer r.m.Unlock()

	md[key] = value
	return nil
}
