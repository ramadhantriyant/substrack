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

	// Auth (public)
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.RefreshToken)

	// Protected routes — all handlers here require a valid Bearer JWT
	protectedMux := http.NewServeMux()

	// Category
	protectedMux.HandleFunc("GET /api/category", h.ListCategory)
	protectedMux.HandleFunc("GET /api/category/id/{id}", h.GetCategoryByID)
	protectedMux.HandleFunc("GET /api/category/name/{name}", h.GetCategoryByName)
	protectedMux.HandleFunc("POST /api/category", h.CreateCategory)
	protectedMux.HandleFunc("PUT /api/category/{id}", h.UpdateCategory)
	protectedMux.HandleFunc("PUT /api/category/{id}/name", h.UpdateCategoryName)
	protectedMux.HandleFunc("PUT /api/category/{id}/description", h.UpdateCategoryDescription)
	protectedMux.HandleFunc("DELETE /api/category/{id}", h.DeleteCategory)

	// Subscriptions
	protectedMux.HandleFunc("GET /api/subscription", h.ListSubscription)
	protectedMux.HandleFunc("GET /api/subscription/active", h.ListActiveSubscription)
	protectedMux.HandleFunc("GET /api/subscription/expired", h.ListExpiredSubscription)
	protectedMux.HandleFunc("GET /api/subscription/cycle/{billCycle}", h.ListSubscriptionsByBillingCycle)
	protectedMux.HandleFunc("GET /api/subscription/{id}", h.GetSubscription)
	protectedMux.HandleFunc("POST /api/subscription", h.CreateSubscription)
	protectedMux.HandleFunc("PUT /api/subscription/{id}", h.UpdateSubscription)
	protectedMux.HandleFunc("PUT /api/subscription/{id}/status", h.UpdateSubscriptionStatus)
	protectedMux.HandleFunc("PUT /api/subscription/{id}/cost", h.UpdateSubscriptionCost)
	protectedMux.HandleFunc("PATCH /api/subscription/{id}/pause", h.PauseSubscription)
	protectedMux.HandleFunc("DELETE /api/subscription/{id}", h.DeleteSubscription)

	// User
	protectedMux.HandleFunc("POST /auth/logout", h.Logout)
	protectedMux.HandleFunc("GET /api/user/me", h.GetMe)
	protectedMux.HandleFunc("PUT /api/user/me", h.UpdateMe)
	protectedMux.HandleFunc("PUT /api/user/me/password", h.UpdatePassword)
	protectedMux.HandleFunc("DELETE /api/user/me", h.DeleteMe)
	protectedMux.HandleFunc("GET /api/user/me/subscription", h.ListUserSubscriptions)
	protectedMux.HandleFunc("POST /api/user/me/subscription/{id}", h.AddUserSubscription)
	protectedMux.HandleFunc("DELETE /api/user/me/subscription/{id}", h.RemoveUserSubscription)
	protectedMux.HandleFunc("GET /api/user/me/category", h.ListUserCategories)
	protectedMux.HandleFunc("POST /api/user/me/category/{id}", h.AddUserCategory)
	protectedMux.HandleFunc("DELETE /api/user/me/category/{id}", h.RemoveUserCategory)

	// Mount protected mux last so explicit public routes take priority
	mux.Handle("/", middlewares.RequireAuth(config.JWTSecret)(protectedMux))

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
