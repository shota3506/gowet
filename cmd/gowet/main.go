package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/database/redis"
	"github.com/shota3506/gowet/handler"
)

var db database.DB

func main() {
	// load environment variables
	port := os.Getenv("PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	var db database.DB
	var err error
	if redisHost == "fake" {
		db = database.NewFakeClient()
	} else {
		db, err = redis.NewClient(fmt.Sprintf("%s:%s", redisHost, redisPort))
	}
	if err != nil {
		log.Fatal(err)
	}

	handler := handler.NewHTTPHandler(db)
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
