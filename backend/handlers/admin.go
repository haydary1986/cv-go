package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"cv-go/models"
	"cv-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	DB     *gorm.DB
	AESKey string
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	var totalUsers, totalCVs, pendingCVs, activeToday, guestCVs int64

	h.DB.Model(&models.User{}).Count(&totalUsers)
	h.DB.Model(&models.CV{}).Count(&totalCVs)
	h.DB.Model(&models.CV{}).Where("status = ?", "pending").Count(&pendingCVs)
	h.DB.Model(&models.CV{}).Where("is_guest = ?", true).Count(&guestCVs)

	today := time.Now().Truncate(24 * time.Hour)
	h.DB.Model(&models.ActivityLog{}).Where("created_at >= ?", today).
		Distinct("user_id").Count(&activeToday)

	// Status counts
	var statusCounts []struct {
		Status string
		Count  int64
	}
	h.DB.Model(&models.CV{}).Select("status, count(*) as count").Group("status").Find(&statusCounts)

	c.JSON(http.StatusOK, gin.H{
		"total_users":   totalUsers,
		"total_cvs":     totalCVs,
		"pending_cvs":   pendingCVs,
		"active_today":  activeToday,
		"guest_cvs":     guestCVs,
		"status_counts": statusCounts,
	})
}

func (h *AdminHandler) ListAllCVs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	search := c.Query("search")
	language := c.Query("language")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := h.DB.Model(&models.CV{}).Preload("User")
	if status != "" {
		query = query.Where("cvs.status = ?", status)
	}
	if language != "" {
		query = query.Where("cvs.language = ?", language)
	}
	if search != "" {
		query = query.Joins("JOIN users ON users.id = cvs.user_id").
			Where("cvs.title LIKE ? OR users.full_name_ar LIKE ? OR users.full_name_en LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var cvs []models.CV
	query.Order("cvs.created_at DESC").Offset(offset).Limit(limit).Find(&cvs)

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"cvs":         cvs,
		"total":       total,
		"page":        page,
		"total_pages": totalPages,
	})
}

func (h *AdminHandler) ApproveCV(c *gin.Context) {
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	h.DB.Model(&cv).Update("status", "approved")

	// Send notification
	h.DB.Create(&models.Notification{
		UserID:    cv.UserID,
		TitleAr:   "تمت الموافقة على سيرتك الذاتية",
		TitleEn:   "Your CV has been approved",
		MessageAr: fmt.Sprintf("تمت الموافقة على السيرة الذاتية: %s", cv.Title),
		MessageEn: fmt.Sprintf("Your CV '%s' has been approved", cv.Title),
		Type:      "approval",
		CVID:      &cv.ID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "CV approved"})
}

func (h *AdminHandler) RejectCV(c *gin.Context) {
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	h.DB.Model(&cv).Updates(map[string]interface{}{
		"status":      "rejected",
		"reject_note": input.Reason,
	})

	h.DB.Create(&models.Notification{
		UserID:    cv.UserID,
		TitleAr:   "تم رفض سيرتك الذاتية",
		TitleEn:   "Your CV has been rejected",
		MessageAr: fmt.Sprintf("تم رفض السيرة الذاتية: %s. السبب: %s", cv.Title, input.Reason),
		MessageEn: fmt.Sprintf("Your CV '%s' has been rejected. Reason: %s", cv.Title, input.Reason),
		Type:      "rejection",
		CVID:      &cv.ID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "CV rejected"})
}

func (h *AdminHandler) RequestRevision(c *gin.Context) {
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input struct {
		Note string `json:"note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	h.DB.Model(&cv).Updates(map[string]interface{}{
		"status":      "draft",
		"reject_note": input.Note,
	})

	h.DB.Create(&models.Notification{
		UserID:    cv.UserID,
		TitleAr:   "مطلوب تعديل على سيرتك الذاتية",
		TitleEn:   "Revision requested for your CV",
		MessageAr: fmt.Sprintf("مطلوب تعديل على السيرة الذاتية: %s. الملاحظات: %s", cv.Title, input.Note),
		MessageEn: fmt.Sprintf("Revision requested for CV '%s'. Notes: %s", cv.Title, input.Note),
		Type:      "revision",
		CVID:      &cv.ID,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Revision requested"})
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	role := c.Query("role")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := h.DB.Model(&models.User{}).Preload("Faculty").Preload("Department")
	if search != "" {
		query = query.Where("email LIKE ? OR full_name_ar LIKE ? OR full_name_en LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users)

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total":       total,
		"page":        page,
		"total_pages": totalPages,
	})
}

func (h *AdminHandler) UpdateUserCredits(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input struct {
		Credits int `json:"credits" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Model(&models.User{}).Where("id = ?", userID).Update("ai_credits", input.Credits)
	c.JSON(http.StatusOK, gin.H{"message": "Credits updated"})
}

// Faculty CRUD
func (h *AdminHandler) ListFaculties(c *gin.Context) {
	var faculties []models.Faculty
	h.DB.Preload("Departments").Find(&faculties)
	c.JSON(http.StatusOK, gin.H{"faculties": faculties})
}

func (h *AdminHandler) CreateFaculty(c *gin.Context) {
	var input struct {
		NameAr string `json:"name_ar" binding:"required"`
		NameEn string `json:"name_en" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	faculty := models.Faculty{NameAr: input.NameAr, NameEn: input.NameEn}
	h.DB.Create(&faculty)
	c.JSON(http.StatusCreated, gin.H{"faculty": faculty})
}

func (h *AdminHandler) UpdateFaculty(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input struct {
		NameAr string `json:"name_ar"`
		NameEn string `json:"name_en"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.DB.Model(&models.Faculty{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name_ar": input.NameAr, "name_en": input.NameEn,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Faculty updated"})
}

func (h *AdminHandler) DeleteFaculty(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.DB.Delete(&models.Faculty{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Faculty deleted"})
}

// Department CRUD
func (h *AdminHandler) ListDepartments(c *gin.Context) {
	facultyID := c.Query("faculty_id")
	query := h.DB.Model(&models.Department{}).Preload("Faculty")
	if facultyID != "" {
		query = query.Where("faculty_id = ?", facultyID)
	}
	var departments []models.Department
	query.Find(&departments)
	c.JSON(http.StatusOK, gin.H{"departments": departments})
}

func (h *AdminHandler) CreateDepartment(c *gin.Context) {
	var input struct {
		FacultyID uint   `json:"faculty_id" binding:"required"`
		NameAr    string `json:"name_ar" binding:"required"`
		NameEn    string `json:"name_en" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dept := models.Department{FacultyID: input.FacultyID, NameAr: input.NameAr, NameEn: input.NameEn}
	h.DB.Create(&dept)
	c.JSON(http.StatusCreated, gin.H{"department": dept})
}

func (h *AdminHandler) UpdateDepartment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var input struct {
		NameAr string `json:"name_ar"`
		NameEn string `json:"name_en"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.DB.Model(&models.Department{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name_ar": input.NameAr, "name_en": input.NameEn,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Department updated"})
}

func (h *AdminHandler) DeleteDepartment(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.DB.Delete(&models.Department{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Department deleted"})
}

// Branding
func (h *AdminHandler) GetBranding(c *gin.Context) {
	var branding models.BrandingSetting
	h.DB.FirstOrCreate(&branding, models.BrandingSetting{
		Name: "CV Builder", PrimaryColor: "#0d6efd",
		SecondaryColor: "#6c757d", AccentColor: "#198754",
	})
	c.JSON(http.StatusOK, gin.H{"branding": branding})
}

func (h *AdminHandler) UpdateBranding(c *gin.Context) {
	var input models.BrandingSetting
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var branding models.BrandingSetting
	h.DB.First(&branding)
	h.DB.Model(&branding).Updates(input)
	c.JSON(http.StatusOK, gin.H{"branding": branding})
}

// AI Settings
func (h *AdminHandler) GetAISettings(c *gin.Context) {
	var settings []models.AISetting
	h.DB.Find(&settings)
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func (h *AdminHandler) UpdateAISettings(c *gin.Context) {
	var input struct {
		Provider  string `json:"provider" binding:"required"`
		Model     string `json:"model" binding:"required"`
		APIKey    string `json:"api_key"`
		MaxTokens int    `json:"max_tokens"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate provider
	validProviders := map[string]bool{"openai": true, "gemini": true, "deepseek": true}
	if !validProviders[input.Provider] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider. Must be openai, gemini, or deepseek"})
		return
	}

	if input.MaxTokens < 0 || input.MaxTokens > 16000 {
		input.MaxTokens = 2000
	}

	var setting models.AISetting
	h.DB.Where("provider = ?", input.Provider).FirstOrCreate(&setting, models.AISetting{Provider: input.Provider})

	updates := map[string]interface{}{
		"model":      input.Model,
		"max_tokens": input.MaxTokens,
		"is_active":  input.IsActive,
	}

	if input.APIKey != "" {
		encrypted, err := utils.EncryptAES(input.APIKey, h.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt API key"})
			return
		}
		updates["api_key_enc"] = encrypted
	}

	h.DB.Model(&setting).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"message": "AI settings updated"})
}

// Ad Settings
func (h *AdminHandler) GetAdSettings(c *gin.Context) {
	var settings models.AdSetting
	h.DB.FirstOrCreate(&settings, models.AdSetting{Type: "video", Duration: 5, SkipAfter: 3})
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

func (h *AdminHandler) UpdateAdSettings(c *gin.Context) {
	var input models.AdSetting
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var settings models.AdSetting
	h.DB.First(&settings)
	h.DB.Model(&settings).Updates(input)
	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// Notifications
func (h *AdminHandler) SendNotification(c *gin.Context) {
	var input struct {
		TitleAr      string `json:"title_ar" binding:"required"`
		TitleEn      string `json:"title_en" binding:"required"`
		MessageAr    string `json:"message_ar" binding:"required"`
		MessageEn    string `json:"message_en" binding:"required"`
		Target       string `json:"target"` // all, faculty, department
		FacultyID    *uint  `json:"faculty_id"`
		DepartmentID *uint  `json:"department_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users []models.User
	query := h.DB.Model(&models.User{})
	if input.Target == "faculty" && input.FacultyID != nil {
		query = query.Where("faculty_id = ?", *input.FacultyID)
	} else if input.Target == "department" && input.DepartmentID != nil {
		query = query.Where("department_id = ?", *input.DepartmentID)
	}
	query.Find(&users)

	for _, user := range users {
		h.DB.Create(&models.Notification{
			UserID:    user.ID,
			TitleAr:   input.TitleAr,
			TitleEn:   input.TitleEn,
			MessageAr: input.MessageAr,
			MessageEn: input.MessageEn,
			Type:      "announcement",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Notification sent to %d users", len(users))})
}

// Activity Logs
func (h *AdminHandler) GetActivityLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	action := c.Query("action")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	query := h.DB.Model(&models.ActivityLog{}).Preload("User")
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if search != "" {
		query = query.Joins("JOIN users ON users.id = activity_logs.user_id").
			Where("users.email LIKE ? OR activity_logs.details LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var logs []models.ActivityLog
	query.Order("activity_logs.created_at DESC").Offset(offset).Limit(limit).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
	})
}

// Statistics
func (h *AdminHandler) GetStats(c *gin.Context) {
	// By faculty
	type FacultyStat struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var byFaculty []FacultyStat
	h.DB.Raw(`SELECT f.name_en as name, COUNT(c.id) as count
		FROM cvs c JOIN users u ON u.id = c.user_id
		LEFT JOIN faculties f ON f.id = u.faculty_id
		WHERE c.deleted_at IS NULL GROUP BY f.name_en`).Scan(&byFaculty)

	// By language
	type LangStat struct {
		Language string `json:"language"`
		Count    int64  `json:"count"`
	}
	var byLanguage []LangStat
	h.DB.Model(&models.CV{}).Select("language, count(*) as count").Group("language").Find(&byLanguage)

	// By status
	type StatusStat struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var byStatus []StatusStat
	h.DB.Model(&models.CV{}).Select("status, count(*) as count").Group("status").Find(&byStatus)

	// Monthly trends (last 12 months)
	type MonthStat struct {
		Month string `json:"month"`
		Count int64  `json:"count"`
	}
	var monthlyTrends []MonthStat
	h.DB.Raw(`SELECT strftime('%Y-%m', created_at) as month, count(*) as count
		FROM cvs WHERE deleted_at IS NULL AND created_at >= datetime('now', '-12 months')
		GROUP BY month ORDER BY month`).Scan(&monthlyTrends)

	// Top users
	type TopUser struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var topUsers []TopUser
	h.DB.Raw(`SELECT u.email, u.full_name_en as name, count(c.id) as count
		FROM users u JOIN cvs c ON c.user_id = u.id
		WHERE c.deleted_at IS NULL GROUP BY u.id ORDER BY count DESC LIMIT 10`).Scan(&topUsers)

	c.JSON(http.StatusOK, gin.H{
		"by_faculty":     byFaculty,
		"by_language":    byLanguage,
		"by_status":      byStatus,
		"monthly_trends": monthlyTrends,
		"top_users":      topUsers,
	})
}

// Export Users CSV
func (h *AdminHandler) ExportUsersCSV(c *gin.Context) {
	var users []models.User
	h.DB.Preload("Faculty").Preload("Department").Find(&users)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=users.csv")

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"ID", "Email", "Name (AR)", "Name (EN)", "Role", "Faculty", "Credits", "Created At"})

	for _, u := range users {
		facultyName := ""
		if u.Faculty != nil {
			facultyName = u.Faculty.NameEn
		}
		writer.Write([]string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Email, u.FullNameAr, u.FullNameEn, u.Role,
			facultyName,
			strconv.Itoa(u.AICredits),
			u.CreatedAt.Format("2006-01-02"),
		})
	}
	writer.Flush()
}

// Import Users CSV
func (h *AdminHandler) ImportUsersCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	reader.Read()

	imported := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 4 {
			continue
		}

		// Generate unique random password for each imported user
		randomPw, _ := utils.GenerateToken(12)
		hashedPw, _ := utils.HashPassword(randomPw)
		role := record[3]
		if role != "student" && role != "teacher" && role != "admin" {
			role = "student"
		}
		user := models.User{
			Email:      record[0],
			FullNameAr: record[1],
			FullNameEn: record[2],
			Role:       role,
			Password:   hashedPw,
			AICredits:  10,
			IsActive:   true,
		}

		if err := h.DB.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err == nil {
			imported++
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Imported %d users", imported)})
}
