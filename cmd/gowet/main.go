package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shota3506/gowet/gotools"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)

	http.ListenAndServe(":8080", mux)
}

func handle(w http.ResponseWriter, r *http.Request) {
	modulePath := r.URL.Path[1:]

	workingDir, err := os.MkdirTemp("", "example")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(workingDir)

	// go mod init
	err = gotools.GoModInit(workingDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// go get
	err = gotools.GoGet(modulePath, workingDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// go link
	dir, err := gotools.GoListDir(modulePath, workingDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// go vet
	out, err := gotools.GoVet(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, out.String())
}
