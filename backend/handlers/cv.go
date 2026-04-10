package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cv-go/audit"
	"cv-go/models"
	"cv-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CVHandler struct {
	DB *gorm.DB
}

type CreateCVInput struct {
	Title              string        `json:"title" binding:"required"`
	Language           string        `json:"language" binding:"required"`
	Template           string        `json:"template" binding:"required"`
	Data               models.CVData `json:"data"`
	IsUniversityMember bool          `json:"is_university_member"`
	FacultyID          *uint         `json:"faculty_id"`
	DepartmentID       *uint         `json:"department_id"`
}

type CreateGuestCVInput struct {
	Title              string        `json:"title" binding:"required"`
	Language           string        `json:"language" binding:"required"`
	Template           string        `json:"template" binding:"required"`
	Data               models.CVData `json:"data"`
	GuestName          string        `json:"guest_name" binding:"required"`
	GuestEmail         string        `json:"guest_email"`
	IsUniversityMember bool          `json:"is_university_member"`
	FacultyID          *uint         `json:"faculty_id"`
	DepartmentID       *uint         `json:"department_id"`
}

func (h *CVHandler) CreateGuestCV(c *gin.Context) {
	var input CreateGuestCVInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shareToken, err := utils.GenerateToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share token"})
		return
	}

	shareURL := fmt.Sprintf("%s/shared/%s", c.Request.Host, shareToken)

	cv := models.CV{
		UserID:             0,
		Title:              input.Title,
		Language:           input.Language,
		Template:           input.Template,
		Data:               input.Data,
		Status:             "draft",
		ShareToken:         shareToken,
		IsShared:           true,
		QRCodeData:         shareURL,
		IsGuest:            true,
		GuestName:          input.GuestName,
		GuestEmail:         input.GuestEmail,
		GuestIP:            c.ClientIP(),
		IsUniversityMember: input.IsUniversityMember,
		FacultyID:          input.FacultyID,
		DepartmentID:       input.DepartmentID,
	}

	if err := h.DB.Create(&cv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CV", "detail": err.Error()})
		return
	}

	h.DB.Create(&models.ActivityLog{
		UserID: 0, Action: "guest_create_cv",
		Details:   fmt.Sprintf("CV ID: %d, Guest: %s (%s)", cv.ID, input.GuestName, input.GuestEmail),
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})

	c.JSON(http.StatusCreated, gin.H{"cv": cv})
}

// CreateCV godoc
// @Summary Create a new CV
// @Description Create a new CV for the authenticated user
// @Tags CV
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CreateCVInput true "CV data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/cvs [post]
func (h *CVHandler) CreateCV(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input CreateCVInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate secure share token
	shareToken, err := utils.GenerateToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share token"})
		return
	}

	shareURL := fmt.Sprintf("%s/shared/%s", c.Request.Host, shareToken)

	cv := models.CV{
		UserID:             userID,
		Title:              input.Title,
		Language:           input.Language,
		Template:           input.Template,
		Data:               input.Data,
		Status:             "draft",
		ShareToken:         shareToken,
		QRCodeData:         shareURL,
		IsUniversityMember: input.IsUniversityMember,
		FacultyID:          input.FacultyID,
		DepartmentID:       input.DepartmentID,
	}

	if err := h.DB.Create(&cv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CV"})
		return
	}

	h.DB.Create(&models.ActivityLog{
		UserID: userID, Action: "create_cv",
		Details: fmt.Sprintf("CV ID: %d", cv.ID),
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})
	audit.Log(h.DB, cv.ID, userID, "create", "CV created", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusCreated, gin.H{"cv": cv})
}

// ListCVs godoc
// @Summary List user's CVs
// @Description Retrieve a paginated list of CVs for the authenticated user
// @Tags CV
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(12)
// @Success 200 {object} map[string]interface{}
// @Router /v1/cvs [get]
func (h *CVHandler) ListCVs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}
	offset := (page - 1) * limit

	var cvs []models.CV
	var total int64

	h.DB.Model(&models.CV{}).Where("user_id = ?", userID).Count(&total)
	h.DB.Preload("Faculty").Preload("Department").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&cvs)

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

// GetCV godoc
// @Summary Get a specific CV
// @Description Retrieve a CV by ID (owner or admin only)
// @Tags CV
// @Produce json
// @Security BearerAuth
// @Param id path int true "CV ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /v1/cvs/{id} [get]
func (h *CVHandler) GetCV(c *gin.Context) {
	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	cvID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CV ID"})
		return
	}

	var cv models.CV
	if err := h.DB.Preload("Faculty").Preload("Department").First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	roleStr, _ := role.(string)
	if cv.UserID != userID && roleStr != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	h.DB.Create(&models.ActivityLog{
		UserID: userID, Action: "view_cv",
		Details: fmt.Sprintf("CV ID: %d", cv.ID),
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{"cv": cv})
}

// UpdateCV godoc
// @Summary Update a CV
// @Description Update an existing CV (owner only)
// @Tags CV
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "CV ID"
// @Param input body object true "CV update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /v1/cvs/{id} [put]
func (h *CVHandler) UpdateCV(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CV ID"})
		return
	}

	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	if cv.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var input struct {
		Title    string        `json:"title"`
		Language string        `json:"language"`
		Template string        `json:"template"`
		Data     models.CVData `json:"data"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Language != "" {
		updates["language"] = input.Language
	}
	if input.Template != "" {
		updates["template"] = input.Template
	}
	updates["data"] = input.Data

	h.DB.Model(&cv).Updates(updates)
	h.DB.First(&cv, cvID)

	h.DB.Create(&models.ActivityLog{
		UserID: userID, Action: "edit_cv",
		Details: fmt.Sprintf("CV ID: %d", cv.ID),
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})
	audit.Log(h.DB, cv.ID, userID, "update", "CV updated", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"cv": cv})
}

// DeleteCV godoc
// @Summary Delete a CV
// @Description Delete a CV by ID (owner only)
// @Tags CV
// @Produce json
// @Security BearerAuth
// @Param id path int true "CV ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /v1/cvs/{id} [delete]
func (h *CVHandler) DeleteCV(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CV ID"})
		return
	}

	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	if cv.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	h.DB.Delete(&cv)

	h.DB.Create(&models.ActivityLog{
		UserID: userID, Action: "delete_cv",
		Details: fmt.Sprintf("CV ID: %d", cv.ID),
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})
	audit.Log(h.DB, uint(cvID), userID, "delete", "CV deleted", c.ClientIP(), c.Request.UserAgent())

	c.JSON(http.StatusOK, gin.H{"message": "CV deleted successfully"})
}

func (h *CVHandler) ToggleShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CV ID"})
		return
	}

	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
		return
	}

	if cv.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	h.DB.Model(&cv).Update("is_shared", !cv.IsShared)
	h.DB.First(&cv, cvID)

	c.JSON(http.StatusOK, gin.H{"cv": cv})
}

func (h *CVHandler) GetSharedCV(c *gin.Context) {
	token := c.Param("token")

	var cv models.CV
	if err := h.DB.Where("share_token = ? AND is_shared = ?", token, true).First(&cv).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CV not found or sharing is disabled"})
		return
	}

	// Atomic increment of view count
	h.DB.Model(&models.CV{}).Where("id = ?", cv.ID).Update("view_count", gorm.Expr("view_count + 1"))

	c.JSON(http.StatusOK, gin.H{"cv": cv})
}

func (h *CVHandler) ExportCVsJSON(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cvs []models.CV
	h.DB.Where("user_id = ?", userID).Find(&cvs)

	data, _ := json.MarshalIndent(cvs, "", "  ")

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=my-cvs.json")
	c.Data(http.StatusOK, "application/json", data)
}

func (h *CVHandler) ExportCVsCSV(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cvs []models.CV
	h.DB.Where("user_id = ?", userID).Find(&cvs)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=my-cvs.csv")

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"ID", "Title", "Language", "Template", "Status", "Created At"})

	for _, cv := range cvs {
		writer.Write([]string{
			strconv.FormatUint(uint64(cv.ID), 10),
			sanitizeCSVField(cv.Title),
			cv.Language,
			cv.Template,
			cv.Status,
			cv.CreatedAt.Format("2006-01-02"),
		})
	}
	writer.Flush()
}
