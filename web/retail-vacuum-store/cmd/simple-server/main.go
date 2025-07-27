package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/go-go-golems/go-go-labs/web/retail-vacuum-store/pkg/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Parse command line flags
	var (
		port     = flag.String("port", "8080", "Server port")
		logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Setup logging
	setupLogging(*logLevel)

	// Setup routes
	mux := handlers.SetupRoutes()

	// Create server
	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Info().
		Str("port", *port).
		Str("address", "http://localhost:"+*port).
		Msg("Starting VacuumMart server")

	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}

func setupLogging(level string) {
	// Setup structured logging
	zerolog.TimeFieldFormat = time.RFC3339

	// Set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Use console writer for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
