package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/shota3506/gowet/database"
	"github.com/shota3506/gowet/database/redis"
	"github.com/shota3506/gowet/server"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	// load environment variables
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGTERM, syscall.SIGINT)

	var d database.DB
	var err error
	if redisHost == "fake" {
		d = database.NewFakeClient()
	} else {
		d, err = redis.NewClient(fmt.Sprintf("%s:%s", redisHost, redisPort))
	}
	if err != nil {
		fmt.Printf("failed to connect database: %v\n", err)
		return 1
	}

	s := server.NewServer(port, d)

	errCh := make(chan error, 1)

	go func() {
		errCh <- s.Start()
	}()

	select {
	case <-termCh:
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err := s.Stop(ctx)
		if err != nil {
			fmt.Printf("failed to gracefully shutdown: %v\n", err)
			return 1
		}
		return 0
	case <-errCh:
		return 1
	}
}
