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
	"github.com/ahmadammarm/amami/backend/pkg/utils"
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

	// API Versioning & Global Rate Limiting
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RateLimiter(100, 1*time.Minute)) // 100 req/min for general API
	{
		// Auth Module (Public + Strict Rate Limit)
		auth := v1.Group("/auth", middleware.RateLimiter(5, 1*time.Minute)) // 5 req/min for auth
		{
			auth.POST("/login", config.AuthHandler.Login)
			auth.GET("/me", middleware.RequireAuth(), config.AuthHandler.GetMe)
			auth.POST("/change-password", middleware.RequireAuth(), config.AuthHandler.ChangePassword)
		}

		// User Management Module (Protected)
		users := v1.Group("/users", middleware.RequireAuth(), middleware.RejectScope(utils.ScopePasswordReset))
		{
			users.POST("/invite", middleware.RequirePermission(config.RoleRepo, "user:invite"), config.UserHandler.Handler.InviteUser)
			users.GET("", middleware.RequirePermission(config.RoleRepo, "user:manage"), config.UserHandler.Handler.GetAllUsers)
			users.PATCH("/:id/status", middleware.RequirePermission(config.RoleRepo, "user:manage"), config.UserHandler.Handler.UpdateUserStatus)
			users.DELETE("/:id", middleware.RequirePermission(config.RoleRepo, "user:manage"), config.UserHandler.Handler.DeleteUser)
		}

		// Qurban Module
		qurbanRoutes := v1.Group("/qurban", middleware.RequireAuth())
		{
			// Packages
			qurbanRoutes.GET("/packages", config.QurbanHandler.Handler.GetAllPackages)
			qurbanRoutes.POST("/packages", config.QurbanHandler.Handler.CreatePackage)
			qurbanRoutes.PUT("/packages/:id", config.QurbanHandler.Handler.UpdatePackage)
			qurbanRoutes.DELETE("/packages/:id", config.QurbanHandler.Handler.DeletePackage)

			// Bookings
			qurbanRoutes.GET("/bookings", config.QurbanHandler.Handler.GetBookings)
			qurbanRoutes.POST("/bookings", config.QurbanHandler.Handler.CreateBooking)

			// Animals
			qurbanRoutes.GET("/animals", config.QurbanHandler.Handler.GetAnimals)
			qurbanRoutes.POST("/animals", config.QurbanHandler.Handler.CreateAnimal)
			qurbanRoutes.PUT("/animals/:id", config.QurbanHandler.Handler.UpdateAnimal)
			qurbanRoutes.DELETE("/animals/:id", config.QurbanHandler.Handler.DeleteAnimal)
		}

		// System Settings Module
		settings := v1.Group("/settings")
		{
			settings.GET("/profile", config.SettingsHandler.Handler.GetMosqueProfile)
			
			// Protected settings
			protected := settings.Group("", middleware.RequireAuth(), middleware.RejectScope(utils.ScopePasswordReset))
			{
				protected.PUT("/profile", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.UpdateMosqueProfile)
				protected.GET("/smtp", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.GetSMTPConfig)
				protected.PUT("/smtp", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.UpdateSMTPConfig)
				protected.POST("/smtp/test", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.TestSMTPConnection)
				protected.GET("/health", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.GetSystemHealth)
				protected.GET("/audit", middleware.RequirePermission(config.RoleRepo, "settings:manage"), config.SettingsHandler.Handler.GetAuditLogs)
			}
		}

		// Finance Module
		finance := v1.Group("/finance", middleware.RequireAuth(), middleware.RejectScope(utils.ScopePasswordReset))
		{
			// Read-only for roles with ledger:read
			finance.GET("/funds", middleware.RequirePermission(config.RoleRepo, "ledger:read"), config.FinanceHandler.Handler.GetFunds)
			finance.GET("/transactions", middleware.RequirePermission(config.RoleRepo, "ledger:read"), config.FinanceHandler.Handler.GetTransactions)
			
			// Write-only for roles with ledger:write (Bendahara/SuperAdmin)
			finance.POST("/funds", middleware.RequirePermission(config.RoleRepo, "ledger:write"), config.FinanceHandler.Handler.CreateFund)
			finance.POST("/transactions", middleware.RequirePermission(config.RoleRepo, "ledger:write"), config.FinanceHandler.Handler.CreateTransaction)
		}

		// Dashboard Module
		dashboard := v1.Group("/dashboard", middleware.RequireAuth())
		{
			// Read-only for authenticated users
			dashboard.GET("/metrics", config.DashboardHandler.Handler.GetMetrics)
		}

		// Logistics & Inventory Routes
		logisticsRoutes := v1.Group("/logistics", middleware.RequireAuth())
		logisticsRoutes.Use(middleware.RequirePermission(config.RoleRepo, "inventory:manage"))
		{
			// Assets
			logisticsRoutes.GET("", config.LogisticsHandler.Handler.GetAllAssets)
			logisticsRoutes.POST("/assets", config.LogisticsHandler.Handler.CreateAsset)
			logisticsRoutes.PUT("/assets/:id", config.LogisticsHandler.Handler.UpdateAsset)
			logisticsRoutes.DELETE("/assets/:id", config.LogisticsHandler.Handler.DeleteAsset)

			// Loans
			logisticsRoutes.POST("/assets/:id/loans", config.LogisticsHandler.Handler.CreateLoan)
			logisticsRoutes.PUT("/assets/loans/:loanId/return", config.LogisticsHandler.Handler.ReturnLoan)

			// Agenda
			logisticsRoutes.GET("/agenda", config.LogisticsHandler.Handler.GetAgendas)
			logisticsRoutes.POST("/agenda", config.LogisticsHandler.Handler.CreateAgenda)
			logisticsRoutes.PUT("/agenda/:id", config.LogisticsHandler.Handler.UpdateAgenda)
		}

		// Jamaah Module
		jamaah := v1.Group("/jamaah", middleware.RequireAuth())
		{
			jamaah.GET("", middleware.RequirePermission(config.RoleRepo, "jamaah:manage"), config.JamaahHandler.Handler.GetJamaahs)
			jamaah.GET("/:id", middleware.RequirePermission(config.RoleRepo, "jamaah:manage"), config.JamaahHandler.Handler.GetJamaahByID)
			jamaah.POST("", middleware.RequirePermission(config.RoleRepo, "jamaah:manage"), config.JamaahHandler.Handler.CreateJamaah)
			jamaah.PUT("/:id", middleware.RequirePermission(config.RoleRepo, "jamaah:manage"), config.JamaahHandler.Handler.UpdateJamaah)
		}

		// Zakat Module
		zakat := v1.Group("/zakat", middleware.RequireAuth())
		{
			zakat.GET("/donations", middleware.RequirePermission(config.RoleRepo, "zakat:manage"), config.ZakatHandler.Handler.GetDonations)
			zakat.GET("/distributions", middleware.RequirePermission(config.RoleRepo, "zakat:manage"), config.ZakatHandler.Handler.GetDistributions)
			zakat.POST("/donations", middleware.RequirePermission(config.RoleRepo, "zakat:manage"), config.ZakatHandler.Handler.CollectZakat)
			zakat.POST("/distributions", middleware.RequirePermission(config.RoleRepo, "zakat:manage"), config.ZakatHandler.Handler.DistributeZakat)
		}
	}

	appAddr := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	logger.Info("Starting server on " + appAddr)
	if err := r.Run(appAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
