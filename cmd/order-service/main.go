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
	"ls-0/arti/order/internal/lib/logger/handlers/slogdiscard"
	"ls-0/arti/order/internal/lib/logger/handlers/slogpretty"
	"ls-0/arti/order/internal/web/server"
)

// TODO: Add testing profile
const (
	envLocal = "local"
	envTest  = "test"
)

func main() {
	cfg := config.MustLoad()                                // load config from config file
	log := setupLogger(cfg.Env)                             // setup logger by env
	ctx, cancel := context.WithCancel(context.Background()) // init context
	// storage := postgres.New(ctx, cfg) // setup storage implementation
	// TODO: init kafka
	// go consumer.setUpConsumer()
	defer cancel()

	// TODO: init web-server
	server.MustStart(ctx, cfg, log) // start web-server

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
