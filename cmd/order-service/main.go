/*
	DESCRIPTION
	===========
	Setting up project environment aka pre-nitialization of logger, config, etc.

*/

package main

import (
	"context"
	"log/slog"
	"os"

	// "ls-0/arti/order/internal/config"
	"ls-0/arti/order/internal/config"
	"ls-0/arti/order/internal/kafka"
	"ls-0/arti/order/internal/lib/logger/handlers/slogdiscard"
	"ls-0/arti/order/internal/lib/logger/handlers/slogpretty"
	"ls-0/arti/order/internal/storage/postgres"
	"ls-0/arti/order/internal/storage/safer"
	"ls-0/arti/order/internal/web/server"
)

// TODO: Add testing profile
const (
	envLocal = "local"
	envTest  = "test"
)

func main() {

	// load config from config file
	cfg := config.MustLoad()

	// set up logger by env
	log := setupLogger(cfg.Env)

	// init context
	ctx, cancel := context.WithCancel(context.Background())

	// set up storage implementation
	storage := postgres.New(ctx, cfg)

	// create a map to hold cash data
	msgMap := safer.NewSafeMap()

	// setting up a kafka consumer
	go kafka.SetUpNewConsumer(ctx, storage, msgMap, log, cfg)

	defer cancel()

	// TODO: init web-server
	server.MustStart(ctx, cfg, log, msgMap) // start web-server

	log.Info("Service stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envTest:
		log = slogdiscard.NewDiscardLogger()
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
