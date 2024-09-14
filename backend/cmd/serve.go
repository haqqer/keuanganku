package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/haqqer/keuanganku/database"
	"github.com/haqqer/keuanganku/middleware"
	"github.com/haqqer/keuanganku/routes"
	"github.com/haqqer/keuanganku/utils/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Serve() {
	ctx := context.Background()

	// Create a new serve mux
	// mux := http.NewServeMux()

	pool, err := pgxpool.NewWithConfig(ctx, database.Config())

	if err != nil {
		log.Fatalf("error db connection")
	}

	database.DB = pool
	router := routes.Router()

	auth.AuthConfig = auth.GoogleConfig()

	middlewareLoader := middleware.Loader(
		middleware.Cors,
		middleware.Logging,
	)

	// mux.Handle("/", router.Router())
	// userHandler := handler.UsersHandler{}
	// // Define routes
	// mux.HandleFunc("GET /", handler.HomeHandler)
	// mux.HandleFunc("/api/users", userHandler)

	// Create a new server
	server := &http.Server{
		Addr:    ":8080",
		Handler: middlewareLoader(router),
	}

	// Start the server in a goroutine
	go func() {
		fmt.Println("Server is starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting")
}
