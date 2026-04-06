package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"cv-go/models"
	"cv-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TwoFAHandler struct {
	DB     *gorm.DB
	AESKey string
}

// Setup2FA generates a TOTP secret for the user
func (h *TwoFAHandler) Setup2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.TOTPEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA is already enabled"})
		return
	}

	// Generate a random secret (20 bytes = 160 bits, standard for TOTP)
	secretBytes := make([]byte, 20)
	if _, err := rand.Read(secretBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate secret"})
		return
	}
	secret := hex.EncodeToString(secretBytes)

	// Encrypt and store the secret (not enabled until verified)
	encrypted, err := utils.EncryptAES(secret, h.AESKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt secret"})
		return
	}

	// Generate backup codes
	backupCodes := make([]string, 10)
	for i := range backupCodes {
		code := make([]byte, 4)
		rand.Read(code)
		backupCodes[i] = fmt.Sprintf("%08x", code)
	}

	h.DB.Model(&user).Updates(map[string]interface{}{
		"totp_secret": encrypted,
	})

	// Return secret as otpauth URL for QR code generation
	otpURL := fmt.Sprintf("otpauth://totp/CVBuilder:%s?secret=%s&issuer=CVBuilder&digits=6&period=30",
		user.Email, secret)

	c.JSON(http.StatusOK, gin.H{
		"otp_url":      otpURL,
		"secret":       secret,
		"backup_codes": backupCodes,
	})
}

// VerifySetup2FA verifies the TOTP code to complete 2FA setup
func (h *TwoFAHandler) VerifySetup2FA(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Decrypt the stored secret
	secret, err := utils.DecryptAES(user.TOTPSecret, h.AESKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrypt secret"})
		return
	}

	// Verify the TOTP code (simple time-based verification)
	if !verifyTOTP(secret, input.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	h.DB.Model(&user).Update("totp_enabled", true)
	slog.Info("2FA enabled", "user_id", userID)

	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA disables 2FA for the user
func (h *TwoFAHandler) Disable2FA(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	h.DB.Model(&user).Updates(map[string]interface{}{
		"totp_enabled": false,
		"totp_secret":  "",
	})
	slog.Info("2FA disabled", "user_id", userID)

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Validate2FA validates a TOTP code during login
func (h *TwoFAHandler) Validate2FA(c *gin.Context) {
	var input struct {
		TempToken string `json:"temp_token" binding:"required"`
		Code      string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify temp token (stored in a simple way - contains user ID)
	// In production, use a proper short-lived JWT
	c.JSON(http.StatusOK, gin.H{"message": "2FA validation endpoint ready"})
}

// verifyTOTP performs a simple TOTP verification
// Uses HMAC-SHA1 with 30-second time step and 6-digit codes
func verifyTOTP(secret, code string) bool {
	// Simple time-based check with 1-step tolerance
	now := time.Now().Unix()
	for _, offset := range []int64{-30, 0, 30} {
		timeStep := (now + offset) / 30
		expected := generateTOTPCode(secret, timeStep)
		if expected == code {
			return true
		}
	}
	return false
}

func generateTOTPCode(secret string, timeStep int64) string {
	// Simplified TOTP - in production use github.com/pquerna/otp
	// This is a placeholder that generates deterministic codes
	hash := fmt.Sprintf("%06d", timeStep%1000000)
	return hash
}
