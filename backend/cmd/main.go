package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ahmadammarm/amami/backend/database/seed"
	"github.com/ahmadammarm/amami/backend/internal/di"
	"github.com/ahmadammarm/amami/backend/internal/middleware"
	"github.com/ahmadammarm/amami/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Logger
	logger.Initialize()
	defer logger.Log.Sync()

	// Flags
	seedFlag := flag.Bool("seed", false, "Seed the database with initial data")
	flag.Parse()

	// Initialize all dependencies
	config, err := di.Initialize()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Handle seeding
	if *seedFlag {
		if err := seed.Seed(config.DB); err != nil {
			log.Fatalf("failed to seed database: %v", err)
		}
		os.Exit(0)
	}

	r := gin.Default()

	// Fix Gin warning by setting trusted proxies
	// In production, set this to your proxy IPs (e.g., Nginx)
	_ = r.SetTrustedProxies(nil)

	// Public routes
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			// Apply rate limit to login
			auth.POST("/login", middleware.RateLimit(5, time.Minute), config.AuthHandler.Login)
		}
	}

	appAddr := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	if err := r.Run(appAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
