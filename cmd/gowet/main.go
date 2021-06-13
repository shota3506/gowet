package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/database/redis"
	"github.com/shota3506/gowet/gotools"
)

var db database.DB

func main() {
	// load environment variables
	port := os.Getenv("PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	err := setup(redisHost, redisPort)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)

	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}

func setup(redisHost, redisPort string) error {
	// setup redis
	var err error
	if redisHost == "fake" {
		db = database.NewFakeClient()
	} else {
		db, err = redis.NewClient(fmt.Sprintf("%s:%s", redisHost, redisPort))
	}

	return err
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Path[1:]

	res, err := db.Get(ctx, path)
	if err == nil {
		fmt.Fprint(w, res) // return cached result
		return
	}

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
	err = gotools.GoGet(path, workingDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// go link
	module, err := gotools.GoList(path, workingDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pathVer := fmt.Sprintf("%s@%s", module.Path, module.Version)

	res, err = db.Get(ctx, pathVer)
	if err == nil {
		fmt.Fprint(w, res) // return cached result
		return
	}

	// go vet
	out, err := gotools.GoVet(module.Dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.Set(ctx, pathVer, out.String())
	if err != nil {
		log.Printf("failed to cache result: %v", err)
	}

	fmt.Fprint(w, out.String())
}
