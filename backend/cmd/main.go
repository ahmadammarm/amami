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
	"github.com/gin-contrib/cors"
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

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Fix Gin warning by setting trusted proxies
	_ = r.SetTrustedProxies(nil)

	// API V1 Group
	v1 := r.Group("/api/v1")
	{
		// Auth Module (Public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimit(5, time.Minute), config.AuthHandler.Login)
			auth.GET("/me", middleware.RequireAuth(), config.AuthHandler.GetMe)
		}

		// User Management Module (Protected)
		users := v1.Group("/users", middleware.RequireAuth())
		{
			users.POST("/invite", middleware.RequirePermission(config.RoleRepo, "user:manage"), config.UserHandler.Handler.InviteUser)
			users.GET("", middleware.RequirePermission(config.RoleRepo, "user:manage"), config.UserHandler.Handler.GetAllUsers)
			users.PATCH("/:id/status", middleware.RequirePermission(config.RoleRepo, "user:manage"), config.UserHandler.Handler.UpdateUserStatus)
		}

		// System Settings Module
		settings := v1.Group("/settings")
		{
			settings.GET("/profile", config.SettingsHandler.Handler.GetMosqueProfile)
			
			// Protected settings
			protected := settings.Group("", middleware.RequireAuth())
			{
				protected.PUT("/profile", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.UpdateMosqueProfile)
				protected.GET("/smtp", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.GetSMTPConfig)
				protected.PUT("/smtp", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.UpdateSMTPConfig)
				protected.POST("/smtp/test", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.TestSMTPConnection)
				protected.GET("/health", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.GetSystemHealth)
				protected.GET("/audit", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.GetAuditLogs)
			}
		}
	}

	appAddr := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	logger.Info("Starting server on " + appAddr)
	if err := r.Run(appAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
