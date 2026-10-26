package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ahmadammarm/amami/backend/database/seed"
	"github.com/ahmadammarm/amami/backend/internal/di"
	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/internal/middleware"
	"github.com/ahmadammarm/amami/backend/pkg/logger"
	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// maxDocumentationFiles mirrors the per-request file limit in the agenda service;
// together with MaxUploadBytes it bounds the size of an upload request.
const maxDocumentationFiles = 10

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

	utils.UseJSONFieldNames()

	r := gin.Default()
	r.MaxMultipartMemory = config.AppConfig.MaxUploadBytes

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

	// Uploaded documentation files; stored under unguessable UUID names.
	r.Static(di.UploadURLPrefix, config.AppConfig.UploadDir)

	// perm builds the permission check for a route.
	perm := func(code domain.PermissionCode) gin.HandlerFunc {
		return middleware.RequirePermission(config.RoleRepo, code)
	}
	// protected is the middleware every authenticated, non-auth route group uses.
	protected := []gin.HandlerFunc{middleware.RequireAuth(), middleware.RejectScope(utils.ScopePasswordReset)}

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

		// User Management Module
		users := v1.Group("/users", protected...)
		users.Use(perm(domain.PermUserManage))
		{
			users.GET("", config.UserHandler.Handler.GetAllUsers)
			users.GET("/roles", config.UserHandler.Handler.GetRoles)
			users.POST("/invite", config.UserHandler.Handler.InviteUser)
			users.PATCH("/:id/status", config.UserHandler.Handler.UpdateUserStatus)
			users.DELETE("/:id", config.UserHandler.Handler.DeleteUser)
		}

		// System Settings Module
		settings := v1.Group("/settings")
		{
			settings.GET("/profile", config.SettingsHandler.Handler.GetMosqueProfile)

			managed := settings.Group("", protected...)
			managed.Use(perm(domain.PermSettingsManage))
			{
				managed.PUT("/profile", config.SettingsHandler.Handler.UpdateMosqueProfile)
				managed.GET("/smtp", config.SettingsHandler.Handler.GetSMTPConfig)
				managed.PUT("/smtp", config.SettingsHandler.Handler.UpdateSMTPConfig)
				managed.POST("/smtp/test", config.SettingsHandler.Handler.TestSMTPConnection)
				managed.GET("/health", config.SettingsHandler.Handler.GetSystemHealth)
				managed.GET("/audit", config.SettingsHandler.Handler.GetAuditLogs)
			}
		}

		// Finance Module
		finance := v1.Group("/finance", protected...)
		{
			finance.GET("/funds", perm(domain.PermLedgerRead), config.FinanceHandler.Handler.GetFunds)
			finance.GET("/transactions", perm(domain.PermLedgerRead), config.FinanceHandler.Handler.GetTransactions)

			finance.POST("/funds", perm(domain.PermLedgerWrite), config.FinanceHandler.Handler.CreateFund)
			finance.POST("/transactions", perm(domain.PermLedgerWrite), config.FinanceHandler.Handler.CreateTransaction)
			finance.POST("/transactions/:id/reverse", perm(domain.PermLedgerWrite), config.FinanceHandler.Handler.ReverseTransaction)
		}

		// Dashboard Module
		dashboard := v1.Group("/dashboard", protected...)
		{
			dashboard.GET("/metrics", config.DashboardHandler.Handler.GetMetrics)
		}

		// Jamaah Module
		jamaah := v1.Group("/jamaah", protected...)
		{
			jamaah.GET("", perm(domain.PermJamaahManage), config.JamaahHandler.Handler.GetJamaahs)
			jamaah.GET("/:id", perm(domain.PermJamaahManage), config.JamaahHandler.Handler.GetJamaahByID)
			jamaah.POST("", perm(domain.PermJamaahManage), config.JamaahHandler.Handler.CreateJamaah)
			jamaah.PUT("/:id", perm(domain.PermJamaahManage), config.JamaahHandler.Handler.UpdateJamaah)

			jamaah.GET("/:id/mustahik-survey", perm(domain.PermJamaahManage), config.JamaahHandler.Handler.GetSurvey)
			jamaah.PUT("/:id/mustahik-survey", perm(domain.PermMustahikSurvey), config.JamaahHandler.Handler.SubmitSurvey)
		}

		// Zakat Module
		zakat := v1.Group("/zakat", protected...)
		{
			zakat.GET("/donations", perm(domain.PermZakatRead), config.ZakatHandler.Handler.GetDonations)
			zakat.GET("/distributions", perm(domain.PermZakatRead), config.ZakatHandler.Handler.GetDistributions)
			zakat.POST("/donations", perm(domain.PermZakatManage), config.ZakatHandler.Handler.CollectZakat)
			zakat.POST("/distributions", perm(domain.PermZakatManage), config.ZakatHandler.Handler.DistributeZakat)
		}

		// Qurban Module
		qurban := v1.Group("/qurban", protected...)
		{
			qurban.GET("/packages", perm(domain.PermQurbanRead), config.QurbanHandler.Handler.GetAllPackages)
			qurban.POST("/packages", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.CreatePackage)
			qurban.PUT("/packages/:id", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.UpdatePackage)
			qurban.DELETE("/packages/:id", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.DeletePackage)

			qurban.GET("/bookings", perm(domain.PermQurbanRead), config.QurbanHandler.Handler.GetBookings)
			qurban.POST("/bookings", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.CreateBooking)

			qurban.GET("/animals", perm(domain.PermQurbanRead), config.QurbanHandler.Handler.GetAnimals)
			qurban.POST("/animals", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.CreateAnimal)
			qurban.PUT("/animals/:id", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.UpdateAnimal)
			qurban.DELETE("/animals/:id", perm(domain.PermQurbanManage), config.QurbanHandler.Handler.DeleteAnimal)
		}

		// Inventory Module
		inventory := v1.Group("/inventory", protected...)
		{
			inventory.GET("/assets", perm(domain.PermInventoryRead), config.InventoryHandler.Handler.GetAllAssets)
			inventory.POST("/assets", perm(domain.PermInventoryManage), config.InventoryHandler.Handler.CreateAsset)
			inventory.PUT("/assets/:id", perm(domain.PermInventoryManage), config.InventoryHandler.Handler.UpdateAsset)
			inventory.DELETE("/assets/:id", perm(domain.PermInventoryManage), config.InventoryHandler.Handler.DeleteAsset)

			inventory.POST("/assets/:id/loans", perm(domain.PermInventoryManage), config.InventoryHandler.Handler.CreateLoan)
			inventory.PUT("/assets/loans/:loanId/return", perm(domain.PermInventoryManage), config.InventoryHandler.Handler.ReturnLoan)
		}

		// Agenda Module
		agenda := v1.Group("/agenda", protected...)
		{
			agenda.GET("", perm(domain.PermAgendaRead), config.AgendaHandler.Handler.GetAgendas)
			agenda.POST("", perm(domain.PermAgendaManage), config.AgendaHandler.Handler.CreateAgenda)
			agenda.PUT("/:id", perm(domain.PermAgendaManage), config.AgendaHandler.Handler.UpdateAgenda)
			agenda.DELETE("/:id", perm(domain.PermAgendaManage), config.AgendaHandler.Handler.DeleteAgenda)

			agenda.GET("/:id/documentation", perm(domain.PermAgendaRead), config.AgendaHandler.Handler.GetDocumentation)
			agenda.POST("/:id/documentation",
				perm(domain.PermAgendaManage),
				middleware.LimitBody(config.AppConfig.MaxUploadBytes*maxDocumentationFiles+1<<20),
				config.AgendaHandler.Handler.UploadDocumentation,
			)
			agenda.DELETE("/:id/documentation/:docId", perm(domain.PermAgendaManage), config.AgendaHandler.Handler.DeleteDocumentation)
		}
	}

	appAddr := fmt.Sprintf(":%s", config.AppConfig.AppPort)
	logger.Info("Starting server on " + appAddr)
	if err := r.Run(appAddr); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
