package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type client struct {
	c redis.Conn
}

func NewClient(host string, port int) (*client, error) {
	c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	return &client{c: c}, nil
}

func (r *client) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := redis.Bytes(r.c.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *client) Set(ctx context.Context, key, value string) error {
	_, err := r.c.Do("SET", key, value)
	return err
}
