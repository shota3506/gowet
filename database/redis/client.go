package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

func NewClient(addr string) (*client, error) {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &client{c: c}, nil
}

type client struct {
	c redis.Conn
}

func (r *client) Get(ctx context.Context, key string) (string, error) {
	res, err := redis.String(r.c.Do("GET", key))
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *client) Set(ctx context.Context, key, value string) error {
	_, err := r.c.Do("SET", key, value)
	return err
}
