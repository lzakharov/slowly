package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lzakharov/slowly/slowly"
)

var (
	addr    string
	timeout time.Duration
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "server address")
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "server request timeout")
	flag.Parse()

	log.Println("initializing...")

	app := slowly.NewApp(addr, timeout)

	log.Println("initialized")

	shutdown := make(chan error)

	go func() {
		if err := app.Run(); err != nil {
			shutdown <- err
		}
	}()

	log.Println("running...")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		log.Printf("received shutdown signal: %s", x)
	case err := <-shutdown:
		log.Printf("received shutdown message: %v", err)
	}

	log.Println("stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	err := app.Shutdown(ctx)
	if err != nil {
		log.Panicf("failed to shut down the application: %v", err)
	}

	cancel()

	log.Println("stopped")
}
