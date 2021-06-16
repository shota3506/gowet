package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/gotools"
)

func render(w http.ResponseWriter, in []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(in)
}

func marshal(path string, in []byte) ([]byte, error) {
	if !json.Valid(in) {
		return nil, errors.New("invalid json")
	}

	const tmpl = "{\"path\": \"%s\", \"results\": %s}"
	out := fmt.Sprintf(tmpl, path, in)

	var bf bytes.Buffer
	if err := json.Indent(&bf, []byte(out), "", "\t"); err != nil {
		return nil, err
	}

	return bf.Bytes(), nil
}

func marshalError(err error) []byte {
	message := ""
	if err != nil {
		message = err.Error()
	}

	s, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	return s
}

type handler struct {
	db database.DB
}

func NewHandler(db database.DB) *handler {
	return &handler{db: db}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path[1:]

	workingDir, err := os.MkdirTemp("", "example")
	if err != nil {
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(workingDir)

	res, err := h.db.Get(ctx, path)
	if err == nil {
		render(w, res, http.StatusOK)
		return
	}

	module, err := h.runGet(ctx, path, workingDir)
	if err != nil {
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}

	pathVer := fmt.Sprintf("%s@%s", module.Path, module.Version)
	res, err = h.db.Get(ctx, pathVer)
	if err == nil {
		render(w, res, http.StatusOK)
		return
	}

	res, err = gotools.GoVet(module.Dir)
	if err != nil {
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}

	res, ok := marshalVet(res)
	if !ok {
		err := errors.New("failed to format vet output in JSON")
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}

	res, err = marshal(pathVer, res)
	if err != nil {
		render(w, marshalError(err), http.StatusInternalServerError)
		return
	}

	err = h.db.Set(ctx, pathVer, string(res))
	if err != nil {
		log.Printf("failed to cache result: %v", err)
	}

	render(w, res, http.StatusOK)
}

func (h *handler) runGet(ctx context.Context, path string, workingDir string) (*gotools.Module, error) {
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
