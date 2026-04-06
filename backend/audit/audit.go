package audit

import (
	"cv-go/models"

	"gorm.io/gorm"
)

func Log(db *gorm.DB, cvID, userID uint, action, details, ip, userAgent string) {
	db.Create(&models.AuditEntry{
		CVID: cvID, UserID: userID, Action: action,
		Details: details, IP: ip, UserAgent: userAgent,
	})
}
