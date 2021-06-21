package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func render(w http.ResponseWriter, in []byte, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
