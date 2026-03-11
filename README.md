# CV Builder System | نظام إنشاء السيرة الذاتية

<div align="center">

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go)
![Vue.js](https://img.shields.io/badge/Vue.js-3-4FC08D?style=for-the-badge&logo=vuedotjs)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?style=for-the-badge&logo=typescript)
![Bootstrap](https://img.shields.io/badge/Bootstrap-5.3-7952B3?style=for-the-badge&logo=bootstrap)
![SQLite](https://img.shields.io/badge/SQLite-3-003B57?style=for-the-badge&logo=sqlite)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**نظام شامل لإنشاء السير الذاتية الاحترافية مع دعم الذكاء الاصطناعي**

**A comprehensive CV builder system with AI support for academics and students**

[العربية](#arabic) | [English](#english)

</div>

---

<a name="arabic"></a>
## 🇸🇦 نبذة عن المشروع

نظام متكامل لإنشاء السيرة الذاتية مصمم خصيصاً لمساعدة **الأكاديميين** و**الطلاب** في بناء سيرهم الذاتية بسهولة واحترافية. يهدف النظام إلى تمكين غير التقنيين من إنشاء سير ذاتية احترافية، وتهيئة الطلاب لسوق العمل قبل التخرج.

### الفئة المستهدفة
- **الأكاديميون**: أعضاء هيئة التدريس والباحثون الذين يحتاجون سيراً ذاتية أكاديمية
- **الطلاب**: لتجهيز أنفسهم لسوق العمل قبل وبعد التخرج
- **الجامعات**: كنظام مؤسسي لإدارة السير الذاتية للمنتسبين

---

<a name="english"></a>
## About

A full-stack CV builder system built with **Go** (backend) and **Vue.js 3** (frontend), designed specifically to help **academics and students** create professional CVs easily. The system empowers non-technical users to build professional CVs and helps students prepare for the job market before graduation.

### Target Audience
- **Academics**: Faculty members and researchers who need academic CVs
- **Students**: To prepare themselves for the job market before and after graduation
- **Universities**: As an institutional system for managing CVs

---

## Features | المميزات

### CV Creation | إنشاء السيرة الذاتية
- Comprehensive form covering: personal data, experience, education, skills, languages, academic links
- **16 professional templates**: Academic, ATS, Compact, Creative, Designer, Elegant, Engineer, Executive, Lawyer, Medical, Minimalist, Modern, Professional, Standard, Teacher, Tech
- Personal photo upload with Base64 conversion
- Auto-save across sessions
- Edit and delete CVs from dashboard

### Multi-Language | دعم اللغات
- Arabic (RTL) and English (LTR) with dynamic switching
- Custom fonts: Cairo for Arabic, Inter for English
- Fully translated interface

### PDF Export | تصدير PDF
- High-quality PDF export using html2canvas + jsPDF
- Automatic multi-page support
- Custom filename (username + CV)

### AI Features | الذكاء الاصطناعي
**3 Providers:** OpenAI (GPT-4o, 4o-mini, 4-Turbo, 3.5) | Google Gemini (1.5 Pro, Flash) | DeepSeek (Chat, Coder)

- AI-powered text improvement
- CV analysis with ATS scoring
- Professional cover letter generation (Arabic/English)
- Job suggestions
- Research capabilities evaluation
- Personalized improvement tips
- Credit system for AI operations

### CV Sharing | مشاركة السيرة
- Unique share links (64-character token)
- Auto-generated QR codes
- View count tracking
- Enable/disable sharing
- Professional public display pages

### Authentication & Security | المصادقة والأمان
- Email/password login
- Google OAuth 2.0
- Password encryption (bcrypt)
- Brute force protection (5 attempts → 15-minute lock)
- AES-256-GCM for sensitive data encryption
- RBAC permissions (student, teacher, admin)

### Admin Panel | لوحة الإدارة
- Comprehensive statistics with Chart.js
- CV management: view, search, approve/reject
- User management with credit system
- Faculty & department management
- Branding customization
- AI provider configuration
- Activity logging
- CSV/JSON import/export

### REST API
- `GET /api/v1/faculties` - List faculties
- `GET /api/v1/departments` - List departments
- `GET /api/v1/stats` - System statistics
- Full CORS support

---

## Tech Stack | التقنيات المستخدمة

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.23, Gin Framework, GORM |
| Frontend | Vue.js 3, TypeScript, Composition API |
| UI | Bootstrap 5.3, Font Awesome 6 |
| Database | SQLite |
| Auth | JWT, bcrypt, Google OAuth 2.0 |
| Charts | Chart.js |
| PDF | html2canvas, jsPDF |
| i18n | vue-i18n |
| Deployment | Docker, Docker Compose, Coolify |

---

## Quick Start | البدء السريع

### Prerequisites
- Go 1.23+
- Node.js 20+
- npm or yarn

### Local Development

```bash
# Clone the repository
git clone https://github.com/haydary1986/cv-go.git
cd cv-go

# Start Backend
cd backend
go mod download
go run main.go
# Server runs on http://localhost:8080

# Start Frontend (in another terminal)
cd frontend
npm install
npm run dev
# Frontend runs on http://localhost:5173
```

### Docker Deployment

```bash
# Build and run
docker-compose up -d

# Or build manually
docker build -t cv-builder .
docker run -p 8080:8080 -v cv-data:/app/data cv-builder
```

### Coolify Deployment

1. Connect your GitHub repository in Coolify
2. Select **Docker Compose** as the build method
3. Set the following environment variables:
   - `JWT_SECRET` - A strong random secret key
   - `AES_KEY` - 32-character encryption key
   - `FRONTEND_URL` - Your domain URL (e.g., https://cv.yourdomain.com)
   - `GOOGLE_CLIENT_ID` - (Optional) Google OAuth Client ID
   - `GOOGLE_CLIENT_SECRET` - (Optional) Google OAuth Secret
   - `GOOGLE_REDIRECT_URL` - (Optional) OAuth redirect URL
4. Deploy!

---

## Environment Variables | متغيرات البيئة

| Variable | Description | Default |
|----------|------------|---------|
| `PORT` | Server port | `8080` |
| `DB_PATH` | SQLite database path | `cvbuilder.db` |
| `JWT_SECRET` | JWT signing secret | `change-in-production` |
| `AES_KEY` | 32-char AES encryption key | `0123456789abcdef...` |
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Secret | - |
| `GOOGLE_REDIRECT_URL` | Google OAuth Redirect URL | - |
| `FRONTEND_URL` | Frontend URL for CORS | `http://localhost:5173` |
| `GIN_MODE` | Gin framework mode | `debug` |

---

## Default Admin Account | حساب المدير الافتراضي

| Field | Value |
|-------|-------|
| Email | `admin@cvbuilder.com` |
| Password | `admin123` |

> **Important:** Change the admin password after first login!

---

## Project Structure | هيكل المشروع

```
cv-go/
├── backend/
│   ├── config/        # Application configuration
│   ├── handlers/      # HTTP request handlers
│   │   ├── auth.go    # Authentication (login, register, OAuth)
│   │   ├── cv.go      # CV CRUD operations
│   │   ├── admin.go   # Admin panel operations
│   │   ├── ai.go      # AI features (OpenAI, Gemini, DeepSeek)
│   │   └── notification.go
│   ├── middleware/     # Auth, rate limiting, RBAC
│   ├── models/        # Database models (GORM)
│   ├── utils/         # Crypto utilities
│   ├── main.go        # Entry point & routes
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/  # Reusable components
│   │   │   └── CVTemplates.vue  # 16 CV templates
│   │   ├── views/       # Page views
│   │   │   ├── admin/   # Admin panel views
│   │   │   ├── HomeView.vue
│   │   │   ├── DashboardView.vue
│   │   │   ├── CVFormView.vue
│   │   │   ├── CVPreviewView.vue
│   │   │   └── ...
│   │   ├── stores/      # Pinia stores
│   │   ├── services/    # API service layer
│   │   ├── i18n/        # Arabic & English translations
│   │   ├── router/      # Vue Router configuration
│   │   ├── App.vue
│   │   └── main.ts
│   ├── package.json
│   └── vite.config.ts
├── Dockerfile           # Multi-stage Docker build
├── docker-compose.yml   # Docker Compose for deployment
└── README.md
```

---

## Summary | ملخص بالأرقام

| Item | Count |
|------|-------|
| CV Templates | 16 |
| AI Providers | 3 (OpenAI, Gemini, DeepSeek) |
| AI Models | 8+ |
| Database Tables | 11 |
| Security Systems | 6 (encryption, protection, RBAC, logging...) |
| API Endpoints | 30+ |
| User Roles | 3 (student, teacher, admin) |
| Languages | 2 (Arabic RTL, English LTR) |

---

## Contributing | المساهمة

Contributions are welcome! Please feel free to submit a Pull Request.

## License | الرخصة

This project is licensed under the MIT License.

---

<div align="center">

**Made with ❤️ for academics and students**

**صُنع بحب للأكاديميين والطلاب**

</div>
