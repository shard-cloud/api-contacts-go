package main

import (
	"log"
	"os"

	"api-contacts-go/internal/config"
	"api-contacts-go/internal/database"
	"api-contacts-go/internal/handlers"
	"api-contacts-go/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logrus.SetLevel(logrus.InfoLevel)
	if cfg.Environment == "development" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		logrus.Warn("Failed to run migrations:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"message":   "API is healthy",
			"version":   "1.0.0",
			"timestamp": fiber.Map{"time": "now"},
		})
	})

	// API routes
	api := app.Group("/api/v1")
	handlers.SetupRoutes(api, db)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	logrus.Infof("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
