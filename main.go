package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"luytbq/com.github/server"
)

// create function main
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := 1025

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	log.Printf("listening on: %d", port)
	if err != nil {
		return err
	}

	gs := server.NewGameServer()
	s := &http.Server{
		Handler:      gs,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
