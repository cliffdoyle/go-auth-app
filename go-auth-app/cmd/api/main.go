package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"github.com/joho/godotenv"
	"github.com/cliffdoyle/go-auth-app/internal/api"
	"github.com/cliffdoyle/go-auth-app/internal/database"
	"github.com/cliffdoyle/go-auth-app/internal/repository"
	"github.com/cliffdoyle/go-auth-app/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Error("JWT_SECRET is not set")
		os.Exit(1)
	}
	jwtExpiryStr := os.Getenv("JWT_EXPIRATION_HOURS")
	jwtExpiry, err := strconv.Atoi(jwtExpiryStr)
	if err != nil {
		jwtExpiry = 72
	}

	slog.Info("Initializing application dependencies...")
	db, err := database.Connect(dbURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connection successful")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := api.NewAuthHandler(userService, jwtSecret, jwtExpiry)
	router := api.NewRouter(authHandler, jwtSecret)

	serverAddr := fmt.Sprintf(":%s", port)
	slog.Info(fmt.Sprintf("Starting server on %s", serverAddr))
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}