package main

import (
	"flag"
	"fmt"
	"github.com/cliffdoyle/go-auth-app/internal/api"
	"github.com/cliffdoyle/go-auth-app/internal/database"
	"github.com/cliffdoyle/go-auth-app/internal/repository"
	"github.com/cliffdoyle/go-auth-app/internal/service"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// --- Define our 'create-admin' subcommand ---
	createAdminCmd := flag.NewFlagSet("create-admin", flag.ExitOnError)
	adminName := createAdminCmd.String("name", "", "Admin's name")
	adminEmail := createAdminCmd.String("email", "", "Admin's email address")
	adminPassword := createAdminCmd.String("password", "", "Admin's password (min 8 chars)")

	// Check if a subcommand has been provided
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "create-admin":
			// Parse the flags for this subcommand
			createAdminCmd.Parse(os.Args[2:])

			if *adminEmail == "" || *adminPassword == "" || *adminName == "" {
				slog.Error("Error: --name, --email and --password are required flags")
				createAdminCmd.Usage()
				os.Exit(1)
			}
			if len(*adminPassword) < 8 {
				slog.Error("Error: password must be at least 8 characters long")
				os.Exit(1)
			}

			// --- Run the Admin Creation Logic ---
			runCreateAdmin(*adminName, *adminEmail, *adminPassword)
			return // Exit after the command is done

		default:
			slog.Error(fmt.Sprintf("Unknown command: %s", os.Args[1]))
			os.Exit(1)
		}
	}

	// --- If no subcommand, run the Web Server (the original main logic) ---
	slog.Info("No command provided. Starting web server...")
	runWebServer()
}

// Extracted the web server logic into its own function
func runWebServer() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiryStr := os.Getenv("JWT_EXPIRATION_HOURS")
	jwtExpiry, _ := strconv.Atoi(jwtExpiryStr)
	if jwtExpiry == 0 {
		jwtExpiry = 72
	}

	// --- Dependency Injection for Web Server ---
	db := setupDatabase()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := api.NewAuthHandler(userService, jwtSecret, jwtExpiry)
	router := api.NewRouter(authHandler, jwtSecret)

	// --- Start Server ---
	serverAddr := fmt.Sprintf(":%s", port)
	slog.Info(fmt.Sprintf("Starting server on %s", serverAddr))
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

// NEW function to handle the admin creation logic
func runCreateAdmin(name, email, password string) {
	slog.Info("Running create-admin command...")
	// --- Dependency Injection for CLI command ---
	db := setupDatabase()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	payload := service.SignupPayload{
		Name:     name,
		Email:    email,
		Password: password,
	}

	admin, err := userService.CreateAdmin(payload)
	if err != nil {
		slog.Error("Failed to create admin user", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully created admin user", "email", admin.Email, "id", admin.ID)
}

// NEW helper function to avoid code duplication
func setupDatabase() *gorm.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL is not set")
		os.Exit(1)
	}
	db, err := database.Connect(dbURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connection successful")
	return db
}
