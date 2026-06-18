//go:build go1.8
// +build go1.8

package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sid/service/feature"
	"sid/service/repository"
	"sid/service/routes"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Using pgdriver (recommended)
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(os.Getenv("DSN")),
	))
	// Configure connection pool
	sqldb.SetMaxOpenConns(25)                 // Maximum open connections
	sqldb.SetMaxIdleConns(10)                 // Maximum idle connections
	sqldb.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime
	sqldb.SetConnMaxIdleTime(5 * time.Minute) // Idle connection timeout

	// Test the connection
	if err := sqldb.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := gin.New()

	// Trust only your known proxy addresses
	// router.SetTrustedProxies([]string{"10.0.0.1", "192.168.1.0/24"})

	repositories := repository.New(db)
	features := feature.New(repositories)

	routes.Setup(router, logger, features)

	srv := &http.Server{
		Addr:    os.Getenv("HTTP_ADDR"),
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	if err := sqldb.Close(); err != nil {
		log.Println("Database Close:", err)
	}
}
