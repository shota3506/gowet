package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/gotools"
)

const (
	timeout = 30 * time.Second
)

type handler struct {
	db database.DB
}

func NewHandler(db database.DB) *handler {
	return &handler{
		db: db,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	path := r.URL.Path[1:]

	workingDir, err := os.MkdirTemp("", "example")
	if err != nil {
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(workingDir)

	res, err := h.handle(ctx, path, workingDir)
	if err != nil {
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}

	render(w, res, http.StatusOK)
}

func (h *handler) handle(ctx context.Context, path, workingDir string) ([]byte, error) {
	res, err := h.db.Get(ctx, path)
	if err == nil {
		return nil, err
	}

	module, err := h.getModule(ctx, path, workingDir)
	if err != nil {
		return nil, err
	}

	pathVer := fmt.Sprintf("%s@%s", module.Path, module.Version)
	res, err = h.db.Get(ctx, pathVer)
	if err == nil {
		return res, nil
	}

	res, err = gotools.GoVet(module.Dir)
	if err != nil {
		return nil, err
	}

	res, ok := marshalVet(res)
	if !ok {
		return nil, errors.New("failed to format vet output in JSON")
	}

	res, err = marshal(pathVer, res)
	if err != nil {
		return nil, err
	}

	err = h.db.Set(ctx, pathVer, string(res))
	if err != nil {
		log.Printf("failed to cache result: %v", err)
	}

	return res, nil
}

func (h *handler) getModule(ctx context.Context, path string, workingDir string) (*gotools.Module, error) {
	// go mod init
	err := gotools.GoModInit(workingDir)
	if err != nil {
		return nil, err
	}

	// go get
	err = gotools.GoGet(path, workingDir)
	if err != nil {
		return nil, err
	}

	// go link
	module, err := gotools.GoList(path, workingDir)
	if err != nil {
		return nil, err
	}

	return module, nil
}
