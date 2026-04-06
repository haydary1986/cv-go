// @title CV Builder API
// @version 1.0
// @description API for the University of Heritage CV Management System
// @host cv.erticaz.com
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cv-go/backup"
	"cv-go/config"
	"cv-go/email"
	"cv-go/handlers"
	"cv-go/logger"
	"cv-go/middleware"
	"cv-go/models"
	"cv-go/utils"
	"cv-go/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	logger.Init()
	slog.Info("Starting CV Builder Server...")

	cfg := config.Load()
	slog.Info("Config loaded", "port", cfg.Port, "db", cfg.DBPath)
	middleware.SetJWTSecret(cfg.JWTSecret)

	// Database
	slog.Info("Connecting to database...")
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connected successfully")

	// Auto migrate
	slog.Info("Running migrations...")
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
		&models.AuditEntry{},
	)
	if err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Migrations completed successfully")

	// Seed admin user and initial data
	seedAdmin(db)
	seedFacultiesAndDepartments(db)
	if os.Getenv("SEED_DEMO_DATA") == "true" {
		seedDemoData(db)
	}

	// Backup scheduler
	backupDir := os.Getenv("BACKUP_DIR")
	if backupDir == "" {
		backupDir = "/app/data/backups"
	}
	backupScheduler := backup.NewScheduler(db, cfg.DBPath, backupDir, 24*time.Hour, 7)
	backupScheduler.Start()

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
	twoFAHandler := &handlers.TwoFAHandler{DB: db, AESKey: cfg.AESKey}
	emailSender := email.NewSender()
	if emailSender.IsConfigured() {
		slog.Info("Email notifications enabled", "from", emailSender.From)
	} else {
		slog.Info("Email notifications disabled (set SMTP_ENABLED=true to enable)")
	}

	// WebSocket hub
	wsHub := ws.NewHub()
	go wsHub.Run()

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

	// Security Headers
	r.Use(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob: https:; font-src 'self' https:; connect-src 'self' wss: ws:; frame-ancestors 'none'")
		if c.Request.TLS != nil || os.Getenv("GIN_MODE") == "release" {
			c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		c.Next()
	})

	// Request body size limit (5MB)
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 5<<20)
		c.Next()
	})

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
				authProtected.POST("/2fa/setup", twoFAHandler.Setup2FA)
				authProtected.POST("/2fa/verify-setup", twoFAHandler.VerifySetup2FA)
				authProtected.POST("/2fa/disable", twoFAHandler.Disable2FA)
			}
			auth.POST("/2fa/validate", twoFAHandler.Validate2FA)
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
			cvRoutes.POST("/import/linkedin", handlers.ImportLinkedIn)
		}

		// Shared CV (public)
		api.GET("/shared/:token", cvHandler.GetSharedCV)

		// Guest CV creation (public - no auth required)
		api.POST("/guest/cv", cvHandler.CreateGuestCV)

		// Public branding (no auth required)
		api.GET("/branding", adminHandler.GetPublicBranding)

		// WebSocket endpoint
		api.GET("/ws", middleware.AuthMiddleware(), func(c *gin.Context) {
			ws.ServeWS(wsHub, c)
		})

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
			adminRoutes.POST("/users", adminHandler.CreateUser)
			adminRoutes.PUT("/users/:id", adminHandler.UpdateUser)
			adminRoutes.DELETE("/users/:id", adminHandler.DeleteUser)
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
			adminRoutes.POST("/branding/logo", adminHandler.UploadLogo)
			adminRoutes.GET("/ai-settings", adminHandler.GetAISettings)
			adminRoutes.PUT("/ai-settings", adminHandler.UpdateAISettings)
			adminRoutes.GET("/ad-settings", adminHandler.GetAdSettings)
			adminRoutes.PUT("/ad-settings", adminHandler.UpdateAdSettings)
			adminRoutes.POST("/notifications", adminHandler.SendNotification)
			adminRoutes.GET("/activity-logs", adminHandler.GetActivityLogs)
			adminRoutes.GET("/audit-trail", adminHandler.GetAuditTrail)
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

	// Serve uploaded files (logos, etc.)
	r.Static("/uploads", "/app/data/uploads")

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

	// Graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		slog.Info("Server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	// Close all background goroutines
	wsHub.Close()
	backupScheduler.Close()
	rl.Close()
	loginRL.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}
	slog.Info("Server exited gracefully")
}

func seedAdmin(db *gorm.DB) {
	adminPassword := os.Getenv("ADMIN_INITIAL_PASSWORD")
	if adminPassword == "" {
		adminPassword = "Sakina1990"
	}
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "haydary1986@gmail.com"
	}

	hashedPassword, err := utils.HashPassword(adminPassword)
	if err != nil {
		slog.Warn("Failed to hash admin password", "error", err)
		return
	}

	var admin models.User
	result := db.Where("role = ?", "admin").First(&admin)
	if result.Error != nil {
		// No admin exists, create one
		admin = models.User{
			Email:      adminEmail,
			Password:   hashedPassword,
			FullNameAr: "مدير النظام",
			FullNameEn: "System Admin",
			Role:       "admin",
			AICredits:  999,
			IsActive:   true,
		}
		if err := db.Create(&admin).Error; err != nil {
			slog.Warn("Failed to create admin user", "error", err)
			return
		}
		slog.Info("Admin user created successfully")
	} else {
		// Admin exists, update email and password
		if err := db.Model(&admin).Updates(map[string]interface{}{
			"email":    adminEmail,
			"password": hashedPassword,
		}).Error; err != nil {
			slog.Warn("Failed to update admin", "error", err)
		} else {
			slog.Info("Admin user updated", "email", adminEmail)
		}
	}
}

func seedFacultiesAndDepartments(db *gorm.DB) {
	var count int64
	db.Model(&models.Faculty{}).Count(&count)
	if count > 0 {
		return // Already seeded
	}

	type deptInfo struct {
		NameAr string
		NameEn string
	}
	type facultyInfo struct {
		NameAr      string
		NameEn      string
		Departments []deptInfo
	}

	faculties := []facultyInfo{
		{"كلية الصيدلة", "College of Pharmacy", []deptInfo{
			{"كلية الصيدلة", "College of Pharmacy"},
		}},
		{"كلية طب الاسنان", "College of Dentistry", []deptInfo{
			{"كلية طب الاسنان", "College of Dentistry"},
		}},
		{"كلية التمريض", "College of Nursing", []deptInfo{
			{"كلية التمريض", "College of Nursing"},
		}},
		{"كلية التربية البدنية وعلوم الرياضة", "College of Physical Education and Sport Sciences", []deptInfo{
			{"كلية التربية البدنية وعلوم الرياضة", "College of Physical Education and Sport Sciences"},
		}},
		{"كلية الادارة والاقتصاد", "College of Administration and Economics", []deptInfo{
			{"المحاسبة", "Accounting"},
			{"العلوم المالية والمصرفية", "Financial and Banking Sciences"},
			{"ادارة الاعمال", "Business Administration"},
			{"علوم السياحة", "Tourism Sciences"},
		}},
		{"كلية التقنيات الصحية والطبية", "College of Health and Medical Technologies", []deptInfo{
			{"تقنيات الاشعة", "Radiology Technologies"},
			{"تقنيات التخدير", "Anesthesia Technologies"},
			{"تقنيات المختبرات الطبية", "Medical Laboratory Technologies"},
			{"تقنيات صناعة الاسنان", "Dental Technologies"},
			{"تقنيات البصريات", "Optics Technologies"},
		}},
		{"كلية القانون", "College of Law", []deptInfo{
			{"كلية القانون", "College of Law"},
		}},
		{"كلية الاداب", "College of Arts", []deptInfo{
			{"اللغة الانكليزية", "English Language"},
			{"علم النفس السريري", "Clinical Psychology"},
		}},
		{"كلية العلوم", "College of Science", []deptInfo{
			{"علوم الحاسبات", "Computer Science"},
			{"علوم الادلة الجنائية", "Forensic Science"},
			{"علوم الامن السيبراني", "Cybersecurity"},
			{"علوم الحياة", "Life Sciences"},
			{"علوم الذكاء الاصطناعي", "Artificial Intelligence"},
			{"علوم الفيزياء الطبية", "Medical Physics"},
		}},
		{"كلية التربية", "College of Education", []deptInfo{
			{"التربية الاسلامية", "Islamic Education"},
			{"علوم التربوية والنفسية", "Educational and Psychological Sciences"},
		}},
		{"كلية الهندسة", "College of Engineering", []deptInfo{
			{"هندسة الطب الحياتي", "Biomedical Engineering"},
			{"هندسة تكرير النفط والغاز", "Oil and Gas Refining Engineering"},
			{"هندسة الطائرات", "Aircraft Engineering"},
			{"الهندسة المدنية", "Civil Engineering"},
			{"الهندسة المعمارية", "Architectural Engineering"},
		}},
		{"الكلية التقنية الهندسية", "Technical Engineering College", []deptInfo{
			{"هندسة تقنيات الحاسوب", "Computer Technologies Engineering"},
		}},
		{"كلية الاعلام", "College of Media", []deptInfo{
			{"كلية الاعلام", "College of Media"},
		}},
		{"كلية الفنون", "College of Arts (Fine Arts)", []deptInfo{
			{"التصميم", "Design"},
		}},
	}

	for _, f := range faculties {
		faculty := models.Faculty{
			NameAr: f.NameAr,
			NameEn: f.NameEn,
		}
		db.Create(&faculty)

		for _, d := range f.Departments {
			dept := models.Department{
				FacultyID: faculty.ID,
				NameAr:    d.NameAr,
				NameEn:    d.NameEn,
			}
			db.Create(&dept)
		}
	}
	slog.Info("Seeded faculties with departments", "count", len(faculties))
}

func seedDemoData(db *gorm.DB) {
	var cvCount int64
	db.Model(&models.CV{}).Count(&cvCount)
	if cvCount > 0 {
		return // Already has data
	}

	// Get faculty/department IDs for assigning to users
	var faculties []models.Faculty
	db.Find(&faculties)
	if len(faculties) == 0 {
		return
	}

	var departments []models.Department
	db.Find(&departments)

	// Helper to get faculty/dept IDs
	fid := func(idx int) *uint {
		if idx < len(faculties) {
			return &faculties[idx].ID
		}
		return &faculties[0].ID
	}
	did := func(idx int) *uint {
		if idx < len(departments) {
			return &departments[idx].ID
		}
		return &departments[0].ID
	}

	// --- Demo Users ---
	pw, _ := utils.HashPassword("demo123")

	users := []models.User{
		{
			Email: "ahmed@demo.com", Password: pw,
			FullNameAr: "أحمد محمد علي", FullNameEn: "Ahmed Mohammed Ali",
			Phone: "07701234567", Role: "student", AICredits: 10, IsActive: true,
			FacultyID: fid(8), DepartmentID: did(14), // كلية العلوم - علوم الحاسبات
		},
		{
			Email: "fatima@demo.com", Password: pw,
			FullNameAr: "فاطمة حسين كاظم", FullNameEn: "Fatima Hussein Kadhim",
			Phone: "07709876543", Role: "student", AICredits: 5, IsActive: true,
			FacultyID: fid(0), DepartmentID: did(0), // كلية الصيدلة
		},
		{
			Email: "dr.ali@demo.com", Password: pw,
			FullNameAr: "د. علي عبدالله جاسم", FullNameEn: "Dr. Ali Abdullah Jasim",
			Phone: "07705551234", Role: "teacher", AICredits: 20, IsActive: true,
			FacultyID: fid(10), DepartmentID: did(22), // كلية الهندسة - هندسة الطب الحياتي
		},
		{
			Email: "noor@demo.com", Password: pw,
			FullNameAr: "نور سعد محمود", FullNameEn: "Noor Saad Mahmoud",
			Phone: "07708884321", Role: "student", AICredits: 8, IsActive: true,
			FacultyID: fid(7), DepartmentID: did(12), // كلية الاداب - اللغة الانكليزية
		},
		{
			Email: "dr.sarah@demo.com", Password: pw,
			FullNameAr: "د. سارة خالد عمر", FullNameEn: "Dr. Sarah Khalid Omar",
			Phone: "07702223344", Role: "teacher", AICredits: 15, IsActive: true,
			FacultyID: fid(4), DepartmentID: did(5), // كلية الادارة والاقتصاد - المحاسبة
		},
		{
			Email: "omar@demo.com", Password: pw,
			FullNameAr: "عمر حيدر رشيد", FullNameEn: "Omar Haider Rashid",
			Phone: "07706667788", Role: "student", AICredits: 0, IsActive: true,
			FacultyID: fid(8), DepartmentID: did(17), // كلية العلوم - علوم الامن السيبراني
		},
	}

	for i := range users {
		db.Create(&users[i])
	}
	slog.Info("Seeded demo users", "count", len(users))

	// --- Demo CVs ---
	genToken := func() string {
		t, _ := utils.GenerateToken(32)
		return t
	}

	cvs := []models.CV{
		// Ahmed - Arabic Modern CV (approved, shared)
		{
			UserID: users[0].ID, Title: "سيرتي الذاتية - تقنية المعلومات", Language: "ar", Template: "modern",
			Status: "approved", ShareToken: genToken(), IsShared: true, ViewCount: 42,
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "أحمد محمد علي", Email: "ahmed@demo.com", Phone: "07701234567",
					Address: "بغداد، العراق", DateOfBirth: "1999-03-15", Nationality: "عراقي",
					JobTitle: "مطور برمجيات",
				},
				Objective: "مطور برمجيات طموح أبحث عن فرصة عمل في مجال تطوير تطبيقات الويب والذكاء الاصطناعي",
				Experiences: []models.Experience{
					{Title: "مطور ويب متدرب", Company: "شركة التقنية العراقية", Location: "بغداد", StartDate: "2023-06", EndDate: "2023-12", Description: "تطوير واجهات المستخدم باستخدام Vue.js وتطوير واجهات برمجية REST"},
					{Title: "مساعد مختبر", Company: "جامعة الكوفة", Location: "النجف", StartDate: "2022-09", EndDate: "2023-05", Description: "مساعدة الطلاب في مختبرات البرمجة وقواعد البيانات"},
				},
				Education: []models.Education{
					{Degree: "بكالوريوس علوم الحاسبات", Institution: "جامعة الكوفة", Location: "النجف", StartDate: "2019-09", EndDate: "2023-06", GPA: "3.6/4.0"},
				},
				Skills: []models.Skill{
					{Name: "Go", Level: "intermediate"}, {Name: "Vue.js", Level: "advanced"},
					{Name: "Python", Level: "intermediate"}, {Name: "SQL", Level: "advanced"},
					{Name: "Docker", Level: "beginner"}, {Name: "Git", Level: "advanced"},
				},
				Languages: []models.LangSkill{
					{Name: "العربية", Level: "native"}, {Name: "English", Level: "fluent"},
				},
				Projects: []models.Project{
					{Name: "نظام إدارة المكتبة", Description: "نظام متكامل لإدارة المكتبة الجامعية باستخدام Go و Vue.js", StartDate: "2023-01", EndDate: "2023-05"},
					{Name: "تطبيق الطقس الذكي", Description: "تطبيق موبايل يستخدم الذكاء الاصطناعي للتنبؤ بالطقس", StartDate: "2022-06", EndDate: "2022-09"},
				},
				Certificates: []models.Certificate{
					{Name: "AWS Cloud Practitioner", Issuer: "Amazon Web Services", Date: "2023-08"},
					{Name: "شهادة Google IT Support", Issuer: "Google", Date: "2022-12"},
				},
				References: []models.Reference{
					{Name: "د. محمد حسن", Position: "أستاذ مساعد", Company: "جامعة الكوفة", Email: "m.hassan@uokufa.edu.iq"},
				},
			},
		},
		// Ahmed - English Academic CV (draft)
		{
			UserID: users[0].ID, Title: "Academic CV - Computer Science", Language: "en", Template: "academic",
			Status: "draft", ShareToken: genToken(),
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "Ahmed Mohammed Ali", Email: "ahmed@demo.com", Phone: "+964-770-123-4567",
					Address: "Baghdad, Iraq", DateOfBirth: "1999-03-15", Nationality: "Iraqi",
					JobTitle: "Computer Science Graduate",
				},
				Objective: "Seeking a graduate research position in artificial intelligence and machine learning",
				Education: []models.Education{
					{Degree: "B.Sc. Computer Science", Institution: "University of Kufa", Location: "Najaf, Iraq", StartDate: "2019-09", EndDate: "2023-06", GPA: "3.6/4.0", Description: "Focus on AI and web development"},
				},
				Skills: []models.Skill{
					{Name: "Python", Level: "advanced"}, {Name: "Machine Learning", Level: "intermediate"},
					{Name: "TensorFlow", Level: "beginner"}, {Name: "Research Writing", Level: "intermediate"},
				},
				Languages: []models.LangSkill{
					{Name: "Arabic", Level: "native"}, {Name: "English", Level: "fluent"},
				},
			},
		},
		// Fatima - Arabic Professional CV (pending)
		{
			UserID: users[1].ID, Title: "السيرة الذاتية - صيدلانية", Language: "ar", Template: "professional",
			Status: "pending", ShareToken: genToken(),
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "فاطمة حسين كاظم", Email: "fatima@demo.com", Phone: "07709876543",
					Address: "النجف، العراق", DateOfBirth: "2000-07-22", Nationality: "عراقية",
					JobTitle: "صيدلانية",
				},
				Objective: "صيدلانية متخرجة حديثاً أبحث عن فرصة عمل في مجال الصيدلة السريرية",
				Experiences: []models.Experience{
					{Title: "صيدلانية متدربة", Company: "مستشفى الصدر التعليمي", Location: "النجف", StartDate: "2023-01", EndDate: "2023-06", Description: "تدريب سريري في قسم الصيدلة السريرية"},
				},
				Education: []models.Education{
					{Degree: "بكالوريوس صيدلة", Institution: "جامعة الكوفة", Location: "النجف", StartDate: "2018-09", EndDate: "2023-06", GPA: "3.8/4.0"},
				},
				Skills: []models.Skill{
					{Name: "الصيدلة السريرية", Level: "advanced"}, {Name: "التحليل الدوائي", Level: "intermediate"},
					{Name: "Microsoft Office", Level: "advanced"}, {Name: "البحث العلمي", Level: "intermediate"},
				},
				Languages: []models.LangSkill{
					{Name: "العربية", Level: "native"}, {Name: "English", Level: "conversational"},
				},
				Certificates: []models.Certificate{
					{Name: "شهادة الإسعافات الأولية", Issuer: "الهلال الأحمر العراقي", Date: "2022-05"},
				},
			},
		},
		// Dr. Ali - English Academic CV (approved, shared)
		{
			UserID: users[2].ID, Title: "Academic CV - Biomedical Engineering", Language: "en", Template: "academic",
			Status: "approved", ShareToken: genToken(), IsShared: true, ViewCount: 156,
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "Dr. Ali Abdullah Jasim", Email: "dr.ali@demo.com", Phone: "+964-770-555-1234",
					Address: "Najaf, Iraq", Nationality: "Iraqi",
					JobTitle: "Assistant Professor - Biomedical Engineering",
				},
				Objective: "Experienced academic and researcher in biomedical engineering with focus on medical imaging and AI-assisted diagnostics",
				Experiences: []models.Experience{
					{Title: "Assistant Professor", Company: "University of Kufa", Location: "Najaf", StartDate: "2018-09", Current: true, Description: "Teaching biomedical engineering courses, supervising graduate students, conducting research in medical imaging"},
					{Title: "Research Fellow", Company: "Imperial College London", Location: "London, UK", StartDate: "2016-01", EndDate: "2018-06", Description: "Post-doctoral research in AI-assisted medical image analysis"},
					{Title: "Lecturer", Company: "University of Kufa", Location: "Najaf", StartDate: "2012-09", EndDate: "2015-12", Description: "Teaching undergraduate biomedical engineering courses"},
				},
				Education: []models.Education{
					{Degree: "Ph.D. Biomedical Engineering", Institution: "University of Manchester", Location: "Manchester, UK", StartDate: "2012-09", EndDate: "2016-06", Description: "Thesis: AI-Based Medical Image Segmentation for Cancer Detection"},
					{Degree: "M.Sc. Biomedical Engineering", Institution: "University of Baghdad", Location: "Baghdad, Iraq", StartDate: "2009-09", EndDate: "2011-06", GPA: "3.9/4.0"},
					{Degree: "B.Sc. Biomedical Engineering", Institution: "University of Baghdad", Location: "Baghdad, Iraq", StartDate: "2005-09", EndDate: "2009-06", GPA: "3.7/4.0"},
				},
				Skills: []models.Skill{
					{Name: "MATLAB", Level: "expert"}, {Name: "Python", Level: "advanced"},
					{Name: "Medical Imaging", Level: "expert"}, {Name: "Deep Learning", Level: "advanced"},
					{Name: "Research Methodology", Level: "expert"}, {Name: "Academic Writing", Level: "expert"},
				},
				Languages: []models.LangSkill{
					{Name: "Arabic", Level: "native"}, {Name: "English", Level: "fluent"},
				},
				Projects: []models.Project{
					{Name: "AI-Cancer Detection System", Description: "Deep learning system for early cancer detection from MRI scans", StartDate: "2019-01", EndDate: "2021-12"},
					{Name: "Smart Prosthetics Lab", Description: "Establishing a lab for smart prosthetic limb research", StartDate: "2020-06"},
				},
				Certificates: []models.Certificate{
					{Name: "IEEE Senior Member", Issuer: "IEEE", Date: "2020-01"},
					{Name: "Certified LabVIEW Developer", Issuer: "National Instruments", Date: "2017-06"},
				},
				References: []models.Reference{
					{Name: "Prof. James Smith", Position: "Professor", Company: "University of Manchester", Email: "j.smith@manchester.ac.uk"},
					{Name: "Prof. Hassan Al-Kindi", Position: "Professor", Company: "University of Baghdad", Email: "h.alkindi@uobaghdad.edu.iq"},
				},
				Links: []models.Link{
					{Title: "Google Scholar", URL: "https://scholar.google.com", Type: "academic"},
					{Title: "ResearchGate", URL: "https://researchgate.net", Type: "academic"},
				},
			},
		},
		// Noor - Arabic Modern CV (rejected)
		{
			UserID: users[3].ID, Title: "السيرة الذاتية - مترجمة", Language: "ar", Template: "modern",
			Status: "rejected", RejectNote: "يرجى إضافة المزيد من التفاصيل حول الخبرات العملية والشهادات", ShareToken: genToken(),
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "نور سعد محمود", Email: "noor@demo.com", Phone: "07708884321",
					Address: "النجف، العراق", DateOfBirth: "2001-11-08", Nationality: "عراقية",
					JobTitle: "مترجمة",
				},
				Objective: "مترجمة محترفة في اللغتين العربية والإنكليزية",
				Education: []models.Education{
					{Degree: "بكالوريوس لغة إنكليزية", Institution: "جامعة الكوفة", Location: "النجف", StartDate: "2019-09", EndDate: "2023-06", GPA: "3.5/4.0"},
				},
				Skills: []models.Skill{
					{Name: "الترجمة التحريرية", Level: "advanced"}, {Name: "الترجمة الفورية", Level: "intermediate"},
				},
				Languages: []models.LangSkill{
					{Name: "العربية", Level: "native"}, {Name: "English", Level: "fluent"}, {Name: "Français", Level: "basic"},
				},
			},
		},
		// Noor - English Professional CV (draft)
		{
			UserID: users[3].ID, Title: "Translator CV", Language: "en", Template: "professional",
			Status: "draft", ShareToken: genToken(),
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "Noor Saad Mahmoud", Email: "noor@demo.com", Phone: "+964-770-888-4321",
					Address: "Najaf, Iraq", JobTitle: "Translator & Interpreter",
				},
				Objective: "Professional Arabic-English translator with academic background in English literature",
				Education: []models.Education{
					{Degree: "B.A. English Language", Institution: "University of Kufa", Location: "Najaf, Iraq", StartDate: "2019-09", EndDate: "2023-06", GPA: "3.5/4.0"},
				},
				Skills: []models.Skill{
					{Name: "Translation", Level: "advanced"}, {Name: "Interpretation", Level: "intermediate"},
					{Name: "Proofreading", Level: "advanced"}, {Name: "CAT Tools", Level: "intermediate"},
				},
				Languages: []models.LangSkill{
					{Name: "Arabic", Level: "native"}, {Name: "English", Level: "fluent"}, {Name: "French", Level: "basic"},
				},
			},
		},
		// Dr. Sarah - Arabic Academic CV (approved)
		{
			UserID: users[4].ID, Title: "السيرة الذاتية الأكاديمية", Language: "ar", Template: "academic",
			Status: "approved", ShareToken: genToken(), IsShared: true, ViewCount: 89,
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "د. سارة خالد عمر", Email: "dr.sarah@demo.com", Phone: "07702223344",
					Address: "النجف، العراق", Nationality: "عراقية",
					JobTitle: "أستاذ مساعد - المحاسبة",
				},
				Objective: "أكاديمية وباحثة في مجال المحاسبة والتدقيق المالي مع خبرة تدريسية تزيد عن 8 سنوات",
				Experiences: []models.Experience{
					{Title: "أستاذ مساعد", Company: "جامعة الكوفة", Location: "النجف", StartDate: "2019-09", Current: true, Description: "تدريس مواد المحاسبة المالية والتدقيق، إشراف على رسائل الماجستير"},
					{Title: "مدقق مالي", Company: "ديوان الرقابة المالية", Location: "بغداد", StartDate: "2014-01", EndDate: "2016-08", Description: "تدقيق الحسابات الحكومية وإعداد التقارير المالية"},
				},
				Education: []models.Education{
					{Degree: "دكتوراه محاسبة", Institution: "جامعة بغداد", Location: "بغداد", StartDate: "2016-09", EndDate: "2019-06"},
					{Degree: "ماجستير محاسبة", Institution: "جامعة الكوفة", Location: "النجف", StartDate: "2011-09", EndDate: "2013-06", GPA: "3.85/4.0"},
					{Degree: "بكالوريوس محاسبة", Institution: "جامعة الكوفة", Location: "النجف", StartDate: "2007-09", EndDate: "2011-06", GPA: "3.7/4.0"},
				},
				Skills: []models.Skill{
					{Name: "التدقيق المالي", Level: "expert"}, {Name: "المحاسبة الدولية IFRS", Level: "advanced"},
					{Name: "SPSS", Level: "advanced"}, {Name: "البحث العلمي", Level: "expert"},
				},
				Languages: []models.LangSkill{
					{Name: "العربية", Level: "native"}, {Name: "English", Level: "conversational"},
				},
				Links: []models.Link{
					{Title: "Google Scholar", URL: "https://scholar.google.com", Type: "academic"},
				},
			},
		},
		// Omar - Arabic Modern CV (pending)
		{
			UserID: users[5].ID, Title: "سيرتي الذاتية - أمن سيبراني", Language: "ar", Template: "modern",
			Status: "pending", ShareToken: genToken(),
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "عمر حيدر رشيد", Email: "omar@demo.com", Phone: "07706667788",
					Address: "النجف، العراق", DateOfBirth: "2001-05-20", Nationality: "عراقي",
					JobTitle: "متخصص أمن سيبراني",
				},
				Objective: "طالب متخصص في الأمن السيبراني مع شغف بأمن الشبكات واختبار الاختراق",
				Education: []models.Education{
					{Degree: "بكالوريوس علوم الأمن السيبراني", Institution: "جامعة الكوفة", Location: "النجف", StartDate: "2020-09", EndDate: "2024-06", GPA: "3.4/4.0"},
				},
				Skills: []models.Skill{
					{Name: "Linux", Level: "advanced"}, {Name: "Network Security", Level: "intermediate"},
					{Name: "Python", Level: "intermediate"}, {Name: "Penetration Testing", Level: "beginner"},
					{Name: "Wireshark", Level: "intermediate"}, {Name: "Firewall Configuration", Level: "intermediate"},
				},
				Languages: []models.LangSkill{
					{Name: "العربية", Level: "native"}, {Name: "English", Level: "conversational"},
				},
				Projects: []models.Project{
					{Name: "نظام كشف التسلل", Description: "تطوير نظام كشف تسلل باستخدام التعلم الآلي", StartDate: "2023-09", EndDate: "2024-04"},
				},
				Certificates: []models.Certificate{
					{Name: "CompTIA Security+", Issuer: "CompTIA", Date: "2024-01"},
				},
			},
		},
		// Omar - English CV (draft)
		{
			UserID: users[5].ID, Title: "Cybersecurity Specialist CV", Language: "en", Template: "professional",
			Status: "draft", ShareToken: genToken(),
			Data: models.CVData{
				PersonalInfo: models.PersonalInfo{
					FullName: "Omar Haider Rashid", Email: "omar@demo.com", Phone: "+964-770-666-7788",
					Address: "Najaf, Iraq", JobTitle: "Cybersecurity Specialist",
				},
				Objective: "Cybersecurity student passionate about network security and ethical hacking",
				Education: []models.Education{
					{Degree: "B.Sc. Cybersecurity", Institution: "University of Kufa", Location: "Najaf, Iraq", StartDate: "2020-09", EndDate: "2024-06", GPA: "3.4/4.0"},
				},
				Skills: []models.Skill{
					{Name: "Linux Administration", Level: "advanced"}, {Name: "Network Security", Level: "intermediate"},
					{Name: "Python Scripting", Level: "intermediate"}, {Name: "Penetration Testing", Level: "beginner"},
				},
				Languages: []models.LangSkill{
					{Name: "Arabic", Level: "native"}, {Name: "English", Level: "conversational"},
				},
			},
		},
	}

	for i := range cvs {
		cvs[i].QRCodeData = fmt.Sprintf("https://cvbuilder.example.com/shared/%s", cvs[i].ShareToken)
		db.Create(&cvs[i])
	}
	slog.Info("Seeded demo CVs", "count", len(cvs))

	// --- Notifications ---
	cvID1 := cvs[0].ID
	cvID4 := cvs[4].ID
	notifications := []models.Notification{
		{UserID: users[0].ID, TitleAr: "تمت الموافقة على سيرتك الذاتية", TitleEn: "Your CV has been approved",
			MessageAr: "تمت الموافقة على سيرتك الذاتية 'سيرتي الذاتية - تقنية المعلومات' من قبل المشرف",
			MessageEn: "Your CV 'IT CV' has been approved by the supervisor", Type: "approval", IsRead: true, CVID: &cvID1},
		{UserID: users[1].ID, TitleAr: "سيرتك الذاتية قيد المراجعة", TitleEn: "Your CV is under review",
			MessageAr: "تم إرسال سيرتك الذاتية للمراجعة بنجاح", MessageEn: "Your CV has been submitted for review",
			Type: "announcement", IsRead: false},
		{UserID: users[3].ID, TitleAr: "تم رفض سيرتك الذاتية", TitleEn: "Your CV has been rejected",
			MessageAr: "يرجى إضافة المزيد من التفاصيل حول الخبرات العملية والشهادات",
			MessageEn: "Please add more details about work experience and certificates",
			Type: "rejection", IsRead: false, CVID: &cvID4},
		{UserID: users[0].ID, TitleAr: "إعلان: تحديث النظام", TitleEn: "Announcement: System Update",
			MessageAr: "تم تحديث النظام وإضافة ميزات جديدة تشمل قوالب إضافية ودعم الذكاء الاصطناعي",
			MessageEn: "System updated with new features including additional templates and AI support",
			Type: "announcement", IsRead: false},
		{UserID: users[2].ID, TitleAr: "تمت الموافقة على سيرتك الذاتية", TitleEn: "Your CV has been approved",
			MessageAr: "تمت الموافقة على سيرتك الذاتية الأكاديمية", MessageEn: "Your academic CV has been approved",
			Type: "approval", IsRead: true},
		{UserID: users[4].ID, TitleAr: "تمت الموافقة على سيرتك الذاتية", TitleEn: "Your CV has been approved",
			MessageAr: "تمت الموافقة على سيرتك الذاتية الأكاديمية", MessageEn: "Your academic CV has been approved",
			Type: "approval", IsRead: false},
		{UserID: users[5].ID, TitleAr: "سيرتك الذاتية قيد المراجعة", TitleEn: "Your CV is under review",
			MessageAr: "تم إرسال سيرتك الذاتية للمراجعة", MessageEn: "Your CV has been submitted for review",
			Type: "announcement", IsRead: false},
	}

	for i := range notifications {
		db.Create(&notifications[i])
	}
	slog.Info("Seeded notifications", "count", len(notifications))

	// --- Activity Logs ---
	logs := []models.ActivityLog{
		{UserID: users[0].ID, Action: "register", IP: "192.168.1.10", UserAgent: "Mozilla/5.0"},
		{UserID: users[0].ID, Action: "login", IP: "192.168.1.10", UserAgent: "Mozilla/5.0"},
		{UserID: users[0].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[0].ID), IP: "192.168.1.10"},
		{UserID: users[0].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[1].ID), IP: "192.168.1.10"},
		{UserID: users[1].ID, Action: "register", IP: "192.168.1.20", UserAgent: "Mozilla/5.0"},
		{UserID: users[1].ID, Action: "login", IP: "192.168.1.20", UserAgent: "Mozilla/5.0"},
		{UserID: users[1].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[2].ID), IP: "192.168.1.20"},
		{UserID: users[2].ID, Action: "register", IP: "192.168.1.30", UserAgent: "Mozilla/5.0"},
		{UserID: users[2].ID, Action: "login", IP: "192.168.1.30", UserAgent: "Mozilla/5.0"},
		{UserID: users[2].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[3].ID), IP: "192.168.1.30"},
		{UserID: users[3].ID, Action: "register", IP: "192.168.1.40", UserAgent: "Mozilla/5.0"},
		{UserID: users[3].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[4].ID), IP: "192.168.1.40"},
		{UserID: users[4].ID, Action: "login", IP: "192.168.1.50", UserAgent: "Mozilla/5.0"},
		{UserID: users[4].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[6].ID), IP: "192.168.1.50"},
		{UserID: users[5].ID, Action: "register", IP: "192.168.1.60", UserAgent: "Mozilla/5.0"},
		{UserID: users[5].ID, Action: "create_cv", Details: fmt.Sprintf("CV ID: %d", cvs[7].ID), IP: "192.168.1.60"},
		{UserID: users[5].ID, Action: "view_cv", Details: fmt.Sprintf("CV ID: %d", cvs[0].ID), IP: "192.168.1.60"},
	}

	for i := range logs {
		db.Create(&logs[i])
	}
	slog.Info("Seeded activity logs", "count", len(logs))

	// --- Branding Settings ---
	var brandCount int64
	db.Model(&models.BrandingSetting{}).Count(&brandCount)
	if brandCount == 0 {
		db.Create(&models.BrandingSetting{
			Name:           "جامعة التراث - نظام السيرة الذاتية",
			PrimaryColor:   "#1a5276",
			SecondaryColor: "#2c3e50",
			AccentColor:    "#c0982b",
		})
		slog.Info("Seeded branding settings")
	}

	slog.Info("Demo data seeding completed successfully!")
}
