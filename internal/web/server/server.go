package server

import (
	"context"
	"log/slog"
	"ls-0/arti/order/internal/config"
	"ls-0/arti/order/internal/storage/safer"
	"ls-0/arti/order/internal/web/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func MustStart(ctx context.Context, cfg *config.Config, log *slog.Logger, sfm *safer.SafeMap) {

	log.Info("Starting server...")
	// initialize router
	router := mux.NewRouter()

	// Handler registry
	// home := handlers.NewHomeHandler()
	ordersHandler := handlers.NewOrderHandler(sfm)

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("./front-end"))
	router.Handle("/", fs)

	router.HandleFunc("/order/{order_uid}", ordersHandler.GetOrder).Methods("GET")

	log.Info("Server is running on port 8087")

	// Set a port via configuration
	err := http.ListenAndServe(":8087", router)
	if err != nil {
		log.Error("Failed to start server", slog.String("error", err.Error()))
	}
}
