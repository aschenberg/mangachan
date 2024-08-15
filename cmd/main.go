package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"manga/config"
	"manga/internal/app"
	"manga/pkg/httpserver"
	"manga/pkg/logging"
)

func main() {
	cfg := config.NewConfig()
	log := logging.NewLogger(cfg)
	httpServer := app.InitServer(cfg, log)
	err := waitForSignals(log, httpServer)
	shutdown(err, httpServer, log)

}

func waitForSignals(log logging.Logger, httpServer *httpserver.Server) error {
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		log.Info(logging.General, logging.Startup, "app - Run - signal: "+s.String(), nil)
	case err = <-httpServer.Notify():
		log.Error(logging.General, logging.Startup, fmt.Sprint("app - Run - httpServer.Notify: %w", err), nil)

	}
	return err
}

func shutdown(err error, httpServer *httpserver.Server, log logging.Logger) {
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(logging.General, logging.Startup, fmt.Sprint("app - Run - httpServer.Shutdown: %w", err), nil)
	}

}
