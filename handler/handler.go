package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/gotools"
)

type handler struct {
	db database.DB
}

func NewHandler(db database.DB) *handler {
	return &handler{db: db}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path[1:]

	res, err := h.Run(ctx, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}

func (h *handler) Run(ctx context.Context, path string) ([]byte, error) {
	res, err := h.db.Get(ctx, path)
	if err == nil {
		return res, nil // return cached result
	}

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
	err = gotools.GoGet(path, workingDir)
	if err != nil {
		return nil, err
	}

	// go link
	module, err := gotools.GoList(path, workingDir)
	if err != nil {
		return nil, err
	}

	pathVer := fmt.Sprintf("%s@%s", module.Path, module.Version)

	res, err = h.db.Get(ctx, pathVer)
	if err == nil {
		return res, nil // return cached result
	}

	// go vet
	out, err := gotools.GoVet(module.Dir)
	if err != nil {
		return nil, err
	}

	err = h.db.Set(ctx, pathVer, string(out))
	if err != nil {
		log.Printf("failed to cache result: %v", err)
	}

	return out, nil
}
