package database

import (
	"context"
	"errors"
)

func NewFakeClient() *fakeClient {
	return &fakeClient{}
}

var _ DB = (*fakeClient)(nil)

type fakeClient struct{}

func (f *fakeClient) Get(ctx context.Context, key string) (string, error) {
	return "", errors.New("nil returned")
}

func (f *fakeClient) Set(ctx context.Context, key, value string) error {
	return nil
}
