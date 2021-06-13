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

type httpHandler struct {
	db database.DB
}

func NewHTTPHandler(db database.DB) *httpHandler {
	return &httpHandler{db: db}
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path[1:]

	res, err := h.Run(ctx, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, res)
}

func (h *httpHandler) Run(ctx context.Context, path string) (string, error) {
	res, err := h.db.Get(ctx, path)
	if err == nil {
		return res, nil // return cached result
	}

	workingDir, err := os.MkdirTemp("", "example")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workingDir)

	// go mod init
	err = gotools.GoModInit(workingDir)
	if err != nil {
		return "", err
	}

	// go get
	err = gotools.GoGet(path, workingDir)
	if err != nil {
		return "", err
	}

	// go link
	module, err := gotools.GoList(path, workingDir)
	if err != nil {
		return "", err
	}

	pathVer := fmt.Sprintf("%s@%s", module.Path, module.Version)

	res, err = h.db.Get(ctx, pathVer)
	if err == nil {
		return res, nil // return cached result
	}

	// go vet
	out, err := gotools.GoVet(module.Dir)
	if err != nil {
		return "", err
	}

	err = h.db.Set(ctx, pathVer, out.String())
	if err != nil {
		log.Printf("failed to cache result: %v", err)
	}

	return out.String(), nil
}
