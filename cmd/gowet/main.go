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
	"github.com/shota3506/gowet/database/memory"
	"github.com/shota3506/gowet/database/redis"
	"github.com/shota3506/gowet/server"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	// load environment variables
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	databaseType := os.Getenv("DATABASE_TYPE")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))

	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGTERM, syscall.SIGINT)

	var d database.DB
	var err error
	switch databaseType {
	case "fake":
		d = database.NewFakeClient()
	case "memory":
		d = memory.NewClient()
	case "redis":
		d, err = redis.NewClient(redisHost, redisPort)
	default:
		err = fmt.Errorf("invalid database type: %s", databaseType)
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
