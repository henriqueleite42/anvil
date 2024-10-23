package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	logger.Trace().Msg("start")

	// ----------------------------
	//
	// Secrets
	//
	// ----------------------------

	// ----------------------------
	//
	// Databases
	//
	// ----------------------------

	logger.Trace().Msg("connecting to database")

	logger.Info().Msg("connected to database")

	// ----------------------------
	//
	// Adapters
	//
	// ----------------------------

	logger.Trace().Msg("initializing adapters")

	logger.Info().Msg("adapters initialized")

	// ----------------------------
	//
	// Repositories
	//
	// ----------------------------

	logger.Trace().Msg("initializing repositories")

	logger.Info().Msg("repositories initialized")

	// ----------------------------
	//
	// Services
	//
	// ----------------------------

	logger.Trace().Msg("initializing services")

	logger.Info().Msg("services initialized")

	// ----------------------------
	//
	// Deliveries
	//
	// ----------------------------

	logger.Trace().Msg("initializing deliveries")


	logger.Info().Msg("deliveries initialized")

	// ----------------------------
	//
	// Gracefully shutdown
	//
	// ----------------------------

	logger.Trace().Msg("setup gracefully shutdown")

	logger.Info().Msg("shutdown completed")
}
