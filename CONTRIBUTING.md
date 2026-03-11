# Contributing to CV Builder | المساهمة في نظام السيرة الذاتية

Thank you for your interest in contributing! | شكراً لاهتمامك بالمساهمة!

## How to Contribute | كيفية المساهمة

### Reporting Bugs | الإبلاغ عن الأخطاء
- Use GitHub Issues to report bugs
- Include steps to reproduce the issue
- Add screenshots if applicable

### Feature Requests | طلب مميزات جديدة
- Open an issue with the "enhancement" label
- Describe the feature and its use case
- Explain how it benefits academics/students

### Pull Requests | طلبات السحب
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Test your changes thoroughly
5. Commit: `git commit -m "Add: my feature"`
6. Push: `git push origin feature/my-feature`
7. Open a Pull Request

### Development Setup | إعداد بيئة التطوير

```bash
# Backend
cd backend
go mod download
go run main.go

# Frontend
cd frontend
npm install
npm run dev
```

### Code Style | أسلوب الكتابة
- Go: Follow standard Go conventions (`gofmt`)
- Vue/TypeScript: Follow Vue.js 3 Composition API style
- Use meaningful variable and function names
- Add comments for complex logic

### Translation | الترجمة
- Translations are in `frontend/src/i18n/`
- Arabic: `ar.ts`, English: `en.ts`
- Ensure both languages are updated when adding UI text

## Code of Conduct | قواعد السلوك
- Be respectful and constructive
- Focus on the issue, not the person
- Welcome newcomers
