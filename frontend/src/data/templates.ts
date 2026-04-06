export interface TemplateMeta {
  id: string
  nameAr: string
  nameEn: string
  description: string
  icon: string
  color: string
}

export const templates: TemplateMeta[] = [
  { id: 'modern', nameAr: 'عصري', nameEn: 'Modern', description: 'Clean and modern design', icon: 'fa-laptop-code', color: '#2196F3' },
  { id: 'classic', nameAr: 'كلاسيكي', nameEn: 'Classic', description: 'Traditional professional layout', icon: 'fa-file-alt', color: '#607D8B' },
  { id: 'professional', nameAr: 'احترافي', nameEn: 'Professional', description: 'Corporate professional style', icon: 'fa-briefcase', color: '#1a5276' },
  { id: 'creative', nameAr: 'إبداعي', nameEn: 'Creative', description: 'Creative and colorful design', icon: 'fa-palette', color: '#E91E63' },
  { id: 'elegant', nameAr: 'أنيق', nameEn: 'Elegant', description: 'Sophisticated elegant layout', icon: 'fa-gem', color: '#9C27B0' },
  { id: 'academic', nameAr: 'أكاديمي', nameEn: 'Academic', description: 'Designed for academic CVs', icon: 'fa-graduation-cap', color: '#4CAF50' },
  { id: 'minimal', nameAr: 'بسيط', nameEn: 'Minimal', description: 'Simple and minimalist', icon: 'fa-minus-circle', color: '#795548' },
  { id: 'executive', nameAr: 'تنفيذي', nameEn: 'Executive', description: 'For senior executives', icon: 'fa-crown', color: '#FF9800' },
  { id: 'tech', nameAr: 'تقني', nameEn: 'Tech', description: 'For tech professionals', icon: 'fa-code', color: '#00BCD4' },
  { id: 'compact', nameAr: 'مضغوط', nameEn: 'Compact', description: 'Compact single-page layout', icon: 'fa-compress-alt', color: '#3F51B5' },
  { id: 'bold', nameAr: 'جريء', nameEn: 'Bold', description: 'Bold and impactful design', icon: 'fa-bold', color: '#F44336' },
  { id: 'clean', nameAr: 'نظيف', nameEn: 'Clean', description: 'Clean whitespace layout', icon: 'fa-broom', color: '#009688' },
  { id: 'premium', nameAr: 'مميز', nameEn: 'Premium', description: 'Premium quality template', icon: 'fa-star', color: '#c0982b' },
  { id: 'formal', nameAr: 'رسمي', nameEn: 'Formal', description: 'Formal document style', icon: 'fa-stamp', color: '#455A64' },
  { id: 'timeline', nameAr: 'زمني', nameEn: 'Timeline', description: 'Timeline-based layout', icon: 'fa-clock', color: '#FF5722' },
  { id: 'infographic', nameAr: 'إنفوجرافيك', nameEn: 'Infographic', description: 'Visual infographic style', icon: 'fa-chart-bar', color: '#8BC34A' },
]
