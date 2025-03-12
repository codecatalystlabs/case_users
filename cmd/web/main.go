package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"

	_ "github.com/lib/pq"

	"case/internal/handlers"
)

var store = session.New() // Session store

func main() {

	config := getConfig()

	mlogger := initLogger(config.LogFile)
	// Initialize Fiber app
	app := fiber.New()

	// Serve static files
	app.Static("/static", "../../ui/static")

	// Add Logger middleware
	app.Use(logger.New())

	db := getDB(config, mlogger)
	defer db.Close()

	// Add CSRF Middleware
	/*
		app.Use(csrf.New(csrf.Config{
			KeyLookup:      "header:X-CSRF-Token",
			CookieName:     "csrf_token",
			CookieSecure:   true,
			CookieHTTPOnly: true,
		}))
	*/
	// Set up routes
	SetRoute(app, db, store, mlogger, config) // Pass the appropriate *sql.DB instance here
	fmt.Println("Starting server on:", config.Address)
	// Start the app
	app.Listen(config.Address)
}

// connect to database
func getDB(config handlers.Config, sl *slog.Logger) *sql.DB {
	connStr := "host=localhost user=" + config.Ux + " password='" + config.Px + "' dbname=" + config.Dx + " sslmode=disable"
	println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Cannot reach db:", err.Error())
		log.Fatal(err)
		sl.Error("Failed to connect to database: " + err.Error())
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Cannot reach db:", err.Error())
	}

	return db
}

// get config details
func getConfig() (config handlers.Config) {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return config
	}

	// Decode the JSON data into a Config struct
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return config
	}

	return config
}

func initLogger(logFile string) *slog.Logger {

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a handler for writing logs to the file
	fileHandler := slog.NewTextHandler(file, nil)
	logger := slog.New(fileHandler)

	// Set this logger as the default
	slog.SetDefault(logger)

	return logger
}
