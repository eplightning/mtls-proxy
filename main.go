package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := loadConfig()
	srv := configureServer(conf)

	shutdownCh := shutdownOnSignal(srv)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen error: %v", err)
	}
	<-shutdownCh
}

func shutdownOnSignal(srv *http.Server) chan struct{} {
	shutdownCh := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("shutdown error, forcefully closing: %v", err)
		}

		close(shutdownCh)
	}()

	return shutdownCh
}
