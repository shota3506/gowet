package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shota3506/gowet"
	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/database/redis"
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
	modulePath := r.URL.Path[1:]

	res, err := db.Get(ctx, modulePath)
	if err == nil {
		fmt.Fprint(w, res)
		return
	}

	out, err := gowet.Run(ctx, modulePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.Set(ctx, modulePath, out.String())
	if err != nil {
		log.Printf("failed to cache result: %v", err)
	}

	fmt.Fprint(w, out.String())
}
