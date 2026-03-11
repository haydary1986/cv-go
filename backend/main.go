package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cv-go/config"
	"cv-go/handlers"
	"cv-go/middleware"
	"cv-go/models"
	"cv-go/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting CV Builder Server...")

	cfg := config.Load()
	log.Printf("Config loaded. Port: %s, DB: %s", cfg.Port, cfg.DBPath)
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Database
	log.Println("Connecting to database...")
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Auto migrate
	log.Println("Running migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.CV{},
		&models.Faculty{},
		&models.Department{},
		&models.Notification{},
		&models.ActivityLog{},
		&models.SystemSetting{},
		&models.AdSetting{},
		&models.AIUsageLog{},
		&models.BrandingSetting{},
		&models.AISetting{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migrations completed successfully")

	// Seed admin user
	seedAdmin(db)

	// Initialize handlers
	loginRL := middleware.NewLoginRateLimiter()
	authHandler := &handlers.AuthHandler{
		DB:                 db,
		LoginRL:            loginRL,
		GoogleClientID:     cfg.GoogleClientID,
		GoogleClientSecret: cfg.GoogleClientSecret,
		GoogleRedirectURL:  cfg.GoogleRedirectURL,
		FrontendURL:        cfg.FrontendURL,
	}
	cvHandler := &handlers.CVHandler{DB: db}
	adminHandler := &handlers.AdminHandler{DB: db, AESKey: cfg.AESKey}
	aiHandler := &handlers.AIHandler{DB: db, AESKey: cfg.AESKey}
	notifHandler := &handlers.NotificationHandler{DB: db}

	// Router
	r := gin.Default()

	// CORS - allow same-origin and configured frontend URL
	allowedOrigins := []string{cfg.FrontendURL}
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, o := range allowedOrigins {
				if o == origin {
					return true
				}
			}
			// Allow same-origin requests (no Origin header or empty)
			return origin == ""
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiter
	rl := middleware.NewRateLimiter(100, time.Minute)
	r.Use(rl.Middleware())

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)

			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware())
			{
				authProtected.GET("/profile", authHandler.GetProfile)
				authProtected.PUT("/profile", authHandler.UpdateProfile)
				authProtected.PUT("/change-password", authHandler.ChangePassword)
			}
		}

		// CV routes
		cvRoutes := api.Group("/cvs")
		cvRoutes.Use(middleware.AuthMiddleware())
		{
			cvRoutes.GET("", cvHandler.ListCVs)
			cvRoutes.POST("", cvHandler.CreateCV)
			cvRoutes.GET("/:id", cvHandler.GetCV)
			cvRoutes.PUT("/:id", cvHandler.UpdateCV)
			cvRoutes.DELETE("/:id", cvHandler.DeleteCV)
			cvRoutes.POST("/:id/toggle-share", cvHandler.ToggleShare)
			cvRoutes.GET("/export/json", cvHandler.ExportCVsJSON)
			cvRoutes.GET("/export/csv", cvHandler.ExportCVsCSV)
		}

		// Shared CV (public)
		api.GET("/shared/:token", cvHandler.GetSharedCV)

		// AI routes
		aiRoutes := api.Group("/ai")
		aiRoutes.Use(middleware.AuthMiddleware())
		{
			aiRoutes.POST("/improve-text", aiHandler.ImproveText)
			aiRoutes.POST("/analyze/:id", aiHandler.AnalyzeCV)
			aiRoutes.POST("/cover-letter", aiHandler.GenerateCoverLetter)
			aiRoutes.POST("/suggest-jobs/:id", aiHandler.SuggestJobs)
			aiRoutes.POST("/evaluate-research/:id", aiHandler.EvaluateResearch)
			aiRoutes.POST("/tips/:id", aiHandler.GetAITips)
		}

		// Notification routes
		notifRoutes := api.Group("/notifications")
		notifRoutes.Use(middleware.AuthMiddleware())
		{
			notifRoutes.GET("", notifHandler.GetNotifications)
			notifRoutes.PUT("/:id/read", notifHandler.MarkAsRead)
			notifRoutes.PUT("/read-all", notifHandler.MarkAllAsRead)
			notifRoutes.GET("/unread-count", notifHandler.GetUnreadCount)
			notifRoutes.DELETE("/:id", notifHandler.DeleteNotification)
		}

		// Admin routes
		adminRoutes := api.Group("/admin")
		adminRoutes.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			adminRoutes.GET("/dashboard", adminHandler.GetDashboardStats)
			adminRoutes.GET("/cvs", adminHandler.ListAllCVs)
			adminRoutes.POST("/cvs/:id/approve", adminHandler.ApproveCV)
			adminRoutes.POST("/cvs/:id/reject", adminHandler.RejectCV)
			adminRoutes.POST("/cvs/:id/revision", adminHandler.RequestRevision)
			adminRoutes.GET("/users", adminHandler.ListUsers)
			adminRoutes.PUT("/users/:id/credits", adminHandler.UpdateUserCredits)
			adminRoutes.GET("/faculties", adminHandler.ListFaculties)
			adminRoutes.POST("/faculties", adminHandler.CreateFaculty)
			adminRoutes.PUT("/faculties/:id", adminHandler.UpdateFaculty)
			adminRoutes.DELETE("/faculties/:id", adminHandler.DeleteFaculty)
			adminRoutes.GET("/departments", adminHandler.ListDepartments)
			adminRoutes.POST("/departments", adminHandler.CreateDepartment)
			adminRoutes.PUT("/departments/:id", adminHandler.UpdateDepartment)
			adminRoutes.DELETE("/departments/:id", adminHandler.DeleteDepartment)
			adminRoutes.GET("/branding", adminHandler.GetBranding)
			adminRoutes.PUT("/branding", adminHandler.UpdateBranding)
			adminRoutes.GET("/ai-settings", adminHandler.GetAISettings)
			adminRoutes.PUT("/ai-settings", adminHandler.UpdateAISettings)
			adminRoutes.GET("/ad-settings", adminHandler.GetAdSettings)
			adminRoutes.PUT("/ad-settings", adminHandler.UpdateAdSettings)
			adminRoutes.POST("/notifications", adminHandler.SendNotification)
			adminRoutes.GET("/activity-logs", adminHandler.GetActivityLogs)
			adminRoutes.GET("/stats", adminHandler.GetStats)
			adminRoutes.GET("/users/export/csv", adminHandler.ExportUsersCSV)
			adminRoutes.POST("/users/import/csv", adminHandler.ImportUsersCSV)
		}

		// Public API v1
		v1 := api.Group("/v1")
		{
			v1.GET("/faculties", adminHandler.ListFaculties)
			v1.GET("/departments", adminHandler.ListDepartments)
			v1.GET("/stats", func(c *gin.Context) {
				var userCount, cvCount int64
				db.Model(&models.User{}).Count(&userCount)
				db.Model(&models.CV{}).Count(&cvCount)
				c.JSON(200, gin.H{
					"total_users": userCount,
					"total_cvs":   cvCount,
				})
			})
		}
	}

	// Serve frontend static files
	if _, err := os.Stat("./static"); err == nil {
		r.Static("/assets", "./static/assets")
		r.StaticFile("/vite.svg", "./static/vite.svg")
		r.NoRoute(func(c *gin.Context) {
			c.File("./static/index.html")
		})
	} else {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		})
	}

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(fmt.Sprintf(":%s", cfg.Port))
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, _ := utils.HashPassword("admin123")
		admin := models.User{
			Email:      "admin@cvbuilder.com",
			Password:   hashedPassword,
			FullNameAr: "مدير النظام",
			FullNameEn: "System Admin",
			Role:       "admin",
			AICredits:  999,
			IsActive:   true,
		}
		db.Create(&admin)
		log.Println("Admin user created: admin@cvbuilder.com / admin123")
	}
}
