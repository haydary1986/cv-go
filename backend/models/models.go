package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Email             string         `gorm:"uniqueIndex;not null" json:"email"`
	Password          string         `gorm:"not null" json:"-"`
	FullNameAr        string         `json:"full_name_ar"`
	FullNameEn        string         `json:"full_name_en"`
	Phone             string         `json:"phone"`
	Role              string         `gorm:"default:student" json:"role"` // student, teacher, admin
	FacultyID         *uint          `json:"faculty_id"`
	DepartmentID      *uint          `json:"department_id"`
	Faculty           *Faculty       `gorm:"foreignKey:FacultyID" json:"faculty,omitempty"`
	Department        *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	PhotoURL          string         `json:"photo_url"`
	GoogleID          string         `gorm:"uniqueIndex" json:"-"`
	AICredits         int            `gorm:"default:10" json:"ai_credits"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	LoginAttempts     int            `gorm:"default:0" json:"-"`
	LockedUntil       *time.Time     `json:"-"`
	CVs               []CV           `gorm:"foreignKey:UserID" json:"cvs,omitempty"`
	Notifications     []Notification `gorm:"foreignKey:UserID" json:"-"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// CV represents a curriculum vitae
type CV struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Language    string         `gorm:"default:ar" json:"language"` // ar, en
	Template    string         `gorm:"default:modern" json:"template"`
	Status      string         `gorm:"default:draft" json:"status"` // draft, pending, approved, rejected
	RejectNote  string         `json:"reject_note"`
	Title       string         `json:"title"`
	Data        CVData         `gorm:"serializer:json" json:"data"`
	ShareToken  string         `gorm:"uniqueIndex" json:"share_token"`
	IsShared    bool           `gorm:"default:false" json:"is_shared"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	QRCodeData  string         `json:"qr_code_data"`
	IsGuest            bool           `gorm:"default:false" json:"is_guest"`
	GuestName          string         `json:"guest_name"`
	GuestEmail         string         `json:"guest_email"`
	GuestIP            string         `json:"guest_ip"`
	IsUniversityMember bool           `gorm:"default:false" json:"is_university_member"`
	FacultyID          *uint          `json:"faculty_id"`
	DepartmentID       *uint          `json:"department_id"`
	Faculty            Faculty        `gorm:"foreignKey:FacultyID" json:"faculty,omitempty"`
	Department         Department     `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// CVData holds all CV content
type CVData struct {
	PersonalInfo PersonalInfo  `json:"personal_info"`
	Objective    string        `json:"objective"`
	Experiences  []Experience  `json:"experiences"`
	Education    []Education   `json:"education"`
	Skills       []Skill       `json:"skills"`
	Languages    []LangSkill   `json:"languages"`
	Links        []Link        `json:"links"`
	Projects     []Project     `json:"projects"`
	Certificates []Certificate `json:"certificates"`
	References   []Reference   `json:"references"`
	Photo        string        `json:"photo"` // Base64
}

type PersonalInfo struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	DateOfBirth string `json:"date_of_birth"`
	Nationality string `json:"nationality"`
	LinkedIn    string `json:"linkedin"`
	Website     string `json:"website"`
	GitHub      string `json:"github"`
	JobTitle    string `json:"job_title"`
}

type Experience struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Current     bool   `json:"current"`
	Description string `json:"description"`
}

type Education struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	Location    string `json:"location"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	GPA         string `json:"gpa"`
	Description string `json:"description"`
}

type Skill struct {
	Name  string `json:"name"`
	Level string `json:"level"` // beginner, intermediate, advanced, expert
}

type LangSkill struct {
	Name  string `json:"name"`
	Level string `json:"level"` // basic, conversational, fluent, native
}

type Link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"` // academic, social, portfolio
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type Certificate struct {
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	Date         string `json:"date"`
	ExpiryDate   string `json:"expiry_date"`
	CredentialID string `json:"credential_id"`
	URL          string `json:"url"`
}

type Reference struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// Faculty represents an academic faculty
type Faculty struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	NameAr      string         `gorm:"not null" json:"name_ar"`
	NameEn      string         `gorm:"not null" json:"name_en"`
	Departments []Department   `gorm:"foreignKey:FacultyID" json:"departments,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Department represents an academic department
type Department struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	FacultyID uint           `gorm:"not null;index" json:"faculty_id"`
	Faculty   Faculty        `gorm:"foreignKey:FacultyID" json:"faculty,omitempty"`
	NameAr    string         `gorm:"not null" json:"name_ar"`
	NameEn    string         `gorm:"not null" json:"name_en"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Notification represents a user notification
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	TitleAr   string         `json:"title_ar"`
	TitleEn   string         `json:"title_en"`
	MessageAr string         `json:"message_ar"`
	MessageEn string         `json:"message_en"`
	Type      string         `json:"type"` // approval, rejection, revision, announcement
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CVID      *uint          `json:"cv_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ActivityLog records user activities
type ActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action    string    `gorm:"not null" json:"action"` // login, create_cv, edit_cv, delete_cv, view_cv, download_cv
	Details   string    `json:"details"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// SystemSetting stores system configuration
type SystemSetting struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;not null" json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdSetting stores advertisement configuration
type AdSetting struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Type          string `json:"type"` // video, adsense
	VideoURL      string `json:"video_url"`
	AdsenseCode   string `json:"adsense_code"`
	Duration      int    `json:"duration"`       // seconds
	SkipAfter     int    `json:"skip_after"`     // seconds before skip button
	IsActive      bool   `gorm:"default:false" json:"is_active"`
}

// AIUsageLog tracks AI feature usage
type AIUsageLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Feature   string    `json:"feature"` // improve_text, analyze_cv, cover_letter, job_suggest, research_eval
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	Credits   int       `json:"credits"`
	CreatedAt time.Time `json:"created_at"`
}

// BrandingSetting stores branding configuration
type BrandingSetting struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Name         string `json:"name"`
	LogoURL      string `json:"logo_url"`
	PrimaryColor string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	AccentColor  string `json:"accent_color"`
}

// AISetting stores AI provider configuration
type AISetting struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Provider     string `json:"provider"` // openai, gemini, deepseek
	Model        string `json:"model"`
	APIKeyEnc    string `json:"-"`         // AES-256 encrypted
	IsActive     bool   `gorm:"default:false" json:"is_active"`
	MaxTokens    int    `gorm:"default:2000" json:"max_tokens"`
}
