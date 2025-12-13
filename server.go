package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/handlers"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/middlewares"
	"git.ramadhantriyant.id/ramadhantriyant/substrack/internal/models"
)

func createServer(config *models.AppConfig, port string) *http.Server {
	mux := http.NewServeMux()
	h := handlers.New(config)

	// Category
	mux.HandleFunc("GET /api/category", h.ListCategory)
	mux.HandleFunc("GET /api/category/id/{id}", h.GetCategoryByID)
	mux.HandleFunc("GET /api/category/name/{name}", h.GetCategoryByName)
	mux.HandleFunc("POST /api/category", h.CreateCategory)
	mux.HandleFunc("PUT /api/category/{id}", h.UpdateCategory)
	mux.HandleFunc("PUT /api/category/{id}/name", h.UpdateCategoryName)
	mux.HandleFunc("PUT /api/category/{id}/description", h.UpdateCategoryDescription)
	mux.HandleFunc("DELETE /api/category/{id}", h.DeleteCategory)

	// Subscriptions
	mux.HandleFunc("GET /api/subscription", h.ListSubscription)
	mux.HandleFunc("GET /api/subscription/active", h.ListActiveSubscription)
	mux.HandleFunc("GET /api/subscription/expired", h.ListExpiredSubscription)
	mux.HandleFunc("GET /api/subscribtion/cycle/{billCycle}", h.ListSubscriptionsByBillingCycle)
	mux.HandleFunc("GET /api/subscription/{id}", h.GetSubscription)
	mux.HandleFunc("POST /api/subscription", h.CreateSubscription)
	mux.HandleFunc("PUT /api/subscription/{id}", h.UpdateSubscription)
	mux.HandleFunc("PUT /api/subscription/{id}/status", h.UpdateSubscriptionStatus)
	mux.HandleFunc("PUT /api/subscription/{id}/cost", h.UpdateSubscriptionCost)
	mux.HandleFunc("PATCH /api/subscription/{id}/pause", h.PauseSubscription)
	mux.HandleFunc("DELETE /api/subscription/{id}", h.DeleteSubscription)

	handler := middlewares.Chain(mux, middlewares.Logger, middlewares.CORS, middlewares.ShouldJSON)
	log.Printf("listening to port 0.0.0.0%s", port)
	return &http.Server{
		Addr:    port,
		Handler: handler,
	}
}

func runServer(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	serverErr := make(chan error, 1)

	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("Shutting down...")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}

	log.Println("Shutdown complete")
	return nil
}
