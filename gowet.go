package gowet

import (
	"bytes"
	"context"
	"os"

	"github.com/shota3506/gowet/gotools"
)

func Run(ctx context.Context, modulePath string) (*bytes.Buffer, error) {
	workingDir, err := os.MkdirTemp("", "example")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workingDir)

	// go mod init
	err = gotools.GoModInit(workingDir)
	if err != nil {
		return nil, err
	}

	// go get
	err = gotools.GoGet(modulePath, workingDir)
	if err != nil {
		return nil, err
	}

	// go link
	dir, err := gotools.GoListDir(modulePath, workingDir)
	if err != nil {
		return nil, err
	}

	// go vet
	out, err := gotools.GoVet(dir)
	if err != nil {
		return nil, err
	}
	return out, nil
}
