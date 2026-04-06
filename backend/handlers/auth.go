package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"cv-go/middleware"
	"cv-go/models"
	"cv-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB                 *gorm.DB
	LoginRL            *middleware.LoginRateLimiter
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL        string
}

type RegisterInput struct {
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FullNameAr   string `json:"full_name_ar" binding:"required"`
	FullNameEn   string `json:"full_name_en" binding:"required"`
	Phone        string `json:"phone"`
	FacultyID    *uint  `json:"faculty_id"`
	DepartmentID *uint  `json:"department_id"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existing models.User
	if err := h.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Email:        input.Email,
		Password:     hashedPassword,
		FullNameAr:   input.FullNameAr,
		FullNameEn:   input.FullNameEn,
		Phone:        input.Phone,
		FacultyID:    input.FacultyID,
		DepartmentID: input.DepartmentID,
		Role:         "student",
		AICredits:    10,
		IsActive:     true,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log activity
	h.DB.Create(&models.ActivityLog{
		UserID: user.ID, Action: "register",
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check rate limit
	if !h.LoginRL.Check(input.Email) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Account locked. Try again after 15 minutes"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		h.LoginRL.RecordFailure(input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is deactivated"})
		return
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		h.LoginRL.RecordFailure(input.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Reset login attempts on success
	h.LoginRL.Reset(input.Email)

	token, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Log activity
	h.DB.Create(&models.ActivityLog{
		UserID: user.ID, Action: "login",
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := utils.GenerateToken(16)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	// Store state in a secure cookie
	isSecure := os.Getenv("GIN_MODE") == "release"
	c.SetCookie("oauth_state", state, 300, "/", "", isSecure, true)
	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email+profile&access_type=offline&state=%s",
		url.QueryEscape(h.GoogleClientID),
		url.QueryEscape(h.GoogleRedirectURL),
		url.QueryEscape(state),
	)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=no_code")
		return
	}

	// Verify state parameter
	receivedState := c.Query("state")
	expectedState, err := c.Cookie("oauth_state")
	if err != nil || receivedState == "" || receivedState != expectedState {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=invalid_state")
		return
	}
	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", os.Getenv("GIN_MODE") == "release", true)

	// Exchange code for token
	data := url.Values{
		"code":          {code},
		"client_id":     {h.GoogleClientID},
		"client_secret": {h.GoogleClientSecret},
		"redirect_uri":  {h.GoogleRedirectURL},
		"grant_type":    {"authorization_code"},
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=token_exchange")
		return
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := (&http.Client{}).Do(tokenReq)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=token_exchange")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=token_exchange")
		return
	}
	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=token_exchange")
		return
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=no_access_token")
		return
	}

	// Get user info
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{Timeout: 10 * time.Second}
	userResp, err := client.Do(req)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=user_info")
		return
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=user_info")
		return
	}
	var googleUser map[string]interface{}
	if err := json.Unmarshal(userBody, &googleUser); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=user_info")
		return
	}

	googleID, _ := googleUser["id"].(string)
	email, _ := googleUser["email"].(string)
	name, _ := googleUser["name"].(string)

	if email == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=no_email")
		return
	}

	// Find or create user
	var user models.User
	result := h.DB.Where("google_id = ?", googleID).First(&user)
	if result.Error != nil {
		// Try finding by email
		result = h.DB.Where("email = ?", email).First(&user)
		if result.Error != nil {
			// Create new user
			randomPw, err := utils.GenerateToken(16)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=account_creation")
				return
			}
			hashedPw, err := utils.HashPassword(randomPw)
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=account_creation")
				return
			}
			user = models.User{
				Email:      email,
				Password:   hashedPw,
				FullNameEn: name,
				FullNameAr: name,
				GoogleID:   googleID,
				Role:       "student",
				AICredits:  10,
				IsActive:   true,
			}
			if err := h.DB.Create(&user).Error; err != nil {
				c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=account_creation")
				return
			}
		} else {
			// Link Google account
			h.DB.Model(&user).Update("google_id", googleID)
		}
	}

	jwtToken, err := middleware.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=token_generation")
		return
	}

	// Log activity
	h.DB.Create(&models.ActivityLog{
		UserID: user.ID, Action: "login_google",
		IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
	})

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback#token=%s", h.FrontendURL, jwtToken))
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := h.DB.Preload("Faculty").Preload("Department").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		FullNameAr string `json:"full_name_ar"`
		FullNameEn string `json:"full_name_en"`
		Phone      string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	h.DB.First(&user, userID)

	updates := map[string]interface{}{
		"full_name_ar": input.FullNameAr,
		"full_name_en": input.FullNameEn,
		"phone":        input.Phone,
	}
	h.DB.Model(&user).Updates(updates)

	h.DB.Preload("Faculty").Preload("Department").First(&user, userID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	h.DB.First(&user, userID)

	if !utils.CheckPassword(user.Password, input.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	h.DB.Model(&user).Update("password", hashedPassword)

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
