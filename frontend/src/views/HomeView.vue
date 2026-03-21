<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-bg-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>
      <div class="container position-relative">
        <div class="row align-items-center min-vh-75">
          <div class="col-lg-7 text-white py-5">
            <div class="hero-badge mb-3 fade-in">
              <span class="badge bg-white bg-opacity-25 px-3 py-2 rounded-pill fs-6">
                <i class="fas fa-star me-1"></i> {{ t('home.heroBadge') }}
              </span>
            </div>
            <h1 class="display-3 fw-bold mb-3 hero-title fade-in">
              {{ t('home.heroTitle') }}
            </h1>
            <p class="lead mb-4 opacity-90 hero-subtitle fade-in" style="max-width: 540px;">
              {{ t('home.heroSubtitle') }}
            </p>
            <div class="d-flex gap-3 flex-wrap fade-in">
              <template v-if="authStore.isAuthenticated">
                <router-link to="/cv/create" class="btn btn-light btn-lg px-4 fw-semibold shadow-sm hero-btn">
                  <i class="fas fa-plus-circle me-2"></i>{{ t('home.createCVBtn') }}
                </router-link>
                <router-link to="/dashboard" class="btn btn-outline-light btn-lg px-4 fw-semibold hero-btn">
                  <i class="fas fa-th-large me-2"></i>{{ t('home.dashboardBtn') }}
                </router-link>
              </template>
              <template v-else>
                <router-link to="/register" class="btn btn-light btn-lg px-4 fw-semibold shadow-sm hero-btn">
                  <i class="fas fa-rocket me-2"></i>{{ t('home.getStartedBtn') }}
                </router-link>
                <router-link to="/cv/guest" class="btn btn-outline-light btn-lg px-4 fw-semibold hero-btn">
                  <i class="fas fa-eye me-2"></i>{{ t('home.tryDemoBtn') }}
                </router-link>
              </template>
            </div>
            <div class="mt-4 d-flex align-items-center gap-4 hero-trust fade-in">
              <div class="d-flex align-items-center">
                <div class="d-flex">
                  <div class="trust-avatar" style="background: #4fc3f7;">A</div>
                  <div class="trust-avatar" style="background: #81c784; margin-left: -8px;">M</div>
                  <div class="trust-avatar" style="background: #ffb74d; margin-left: -8px;">S</div>
                </div>
                <span class="text-white-50 ms-2 small">{{ t('home.trustedBy') }}</span>
              </div>
            </div>
          </div>
          <div class="col-lg-5 d-none d-lg-block">
            <div class="hero-visual fade-in-right">
              <div class="cv-preview-card">
                <div class="cv-preview-header">
                  <div class="cv-dots">
                    <span class="dot dot-red"></span>
                    <span class="dot dot-yellow"></span>
                    <span class="dot dot-green"></span>
                  </div>
                </div>
                <div class="cv-preview-body">
                  <div class="cv-line cv-line-title"></div>
                  <div class="cv-line cv-line-subtitle"></div>
                  <div class="cv-section-divider"></div>
                  <div class="cv-line cv-line-heading"></div>
                  <div class="cv-line cv-line-text"></div>
                  <div class="cv-line cv-line-text short"></div>
                  <div class="cv-section-divider"></div>
                  <div class="cv-line cv-line-heading"></div>
                  <div class="cv-line cv-line-text"></div>
                  <div class="cv-line cv-line-text medium"></div>
                  <div class="cv-line cv-line-text short"></div>
                  <div class="cv-section-divider"></div>
                  <div class="cv-skills-row">
                    <span class="cv-skill-chip"></span>
                    <span class="cv-skill-chip"></span>
                    <span class="cv-skill-chip"></span>
                  </div>
                </div>
              </div>
              <div class="floating-badge floating-badge-ai">
                <i class="fas fa-magic"></i> AI
              </div>
              <div class="floating-badge floating-badge-pdf">
                <i class="fas fa-file-pdf"></i> PDF
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features-section py-5" id="features">
      <div class="container py-4">
        <div class="text-center mb-5">
          <span class="section-label">{{ t('home.featuresLabel') }}</span>
          <h2 class="display-5 fw-bold mt-2">{{ t('home.featuresTitle') }}</h2>
          <p class="text-muted lead mx-auto" style="max-width: 600px;">{{ t('home.featuresSubtitle') }}</p>
        </div>
        <div class="row g-4">
          <div class="col-md-6 col-lg-3" v-for="(feature, index) in features" :key="index">
            <div class="feature-card h-100 p-4 text-center">
              <div class="feature-icon-wrapper mb-3">
                <div class="feature-icon" :style="{ background: feature.gradient }">
                  <i :class="feature.icon" class="fa-lg text-white"></i>
                </div>
              </div>
              <h5 class="fw-bold mb-2">{{ t(feature.titleKey) }}</h5>
              <p class="text-muted small mb-0">{{ t(feature.descKey) }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="how-it-works-section py-5 bg-light">
      <div class="container py-4">
        <div class="text-center mb-5">
          <span class="section-label">{{ t('home.howItWorksLabel') }}</span>
          <h2 class="display-5 fw-bold mt-2">{{ t('home.howItWorksTitle') }}</h2>
          <p class="text-muted lead mx-auto" style="max-width: 600px;">{{ t('home.howItWorksSubtitle') }}</p>
        </div>
        <div class="row g-4 justify-content-center">
          <div class="col-md-4" v-for="(step, index) in steps" :key="index">
            <div class="step-card text-center p-4 h-100">
              <div class="step-number-wrapper mb-3">
                <div class="step-number">{{ index + 1 }}</div>
                <div class="step-icon-bg">
                  <i :class="step.icon" class="fa-2x text-primary"></i>
                </div>
              </div>
              <h5 class="fw-bold mb-2">{{ t(step.titleKey) }}</h5>
              <p class="text-muted small mb-0">{{ t(step.descKey) }}</p>
              <div v-if="index < steps.length - 1" class="step-arrow d-none d-md-block">
                <i class="fas fa-chevron-right text-primary opacity-25 fa-2x"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Templates Preview Section -->
    <section class="templates-section py-5">
      <div class="container py-4">
        <div class="text-center mb-5">
          <span class="section-label">{{ t('home.templatesLabel') }}</span>
          <h2 class="display-5 fw-bold mt-2">{{ t('home.templatesTitle') }}</h2>
          <p class="text-muted lead mx-auto" style="max-width: 600px;">{{ t('home.templatesSubtitle') }}</p>
        </div>
        <div class="row g-3">
          <div class="col-6 col-sm-4 col-md-3 col-lg-2" v-for="tmpl in templates" :key="tmpl.name">
            <div class="template-card">
              <div class="template-preview" :style="{ background: tmpl.color }">
                <div class="template-lines">
                  <div class="tl tl-header"></div>
                  <div class="tl tl-sub"></div>
                  <div class="tl-divider"></div>
                  <div class="tl tl-text"></div>
                  <div class="tl tl-text short"></div>
                  <div class="tl tl-text"></div>
                </div>
                <div class="template-overlay">
                  <i class="fas fa-eye"></i>
                </div>
              </div>
              <div class="template-name">
                <small class="fw-semibold">{{ t(`templates.${tmpl.name}`) }}</small>
              </div>
            </div>
          </div>
        </div>
        <div class="text-center mt-4">
          <router-link
            :to="authStore.isAuthenticated ? '/cv/create' : '/register'"
            class="btn btn-primary btn-lg px-5 rounded-pill shadow-sm">
            <i class="fas fa-palette me-2"></i>{{ t('home.browseTemplatesBtn') }}
          </router-link>
        </div>
      </div>
    </section>

    <!-- Statistics Section -->
    <section class="stats-section py-5">
      <div class="container py-4">
        <div class="row g-4 justify-content-center">
          <div class="col-6 col-md-3" v-for="stat in stats" :key="stat.labelKey">
            <div class="stat-card text-center p-4">
              <div class="stat-icon mb-2">
                <i :class="stat.icon" class="fa-2x"></i>
              </div>
              <div class="stat-value display-5 fw-bold">{{ stat.value }}</div>
              <div class="stat-label text-muted small">{{ t(stat.labelKey) }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Guest CTA Section -->
    <section v-if="!authStore.isAuthenticated" class="guest-cta-section py-5 bg-light">
      <div class="container py-4">
        <div class="row justify-content-center">
          <div class="col-lg-8">
            <div class="guest-card p-5 text-center">
              <div class="guest-icon mb-3">
                <i class="fas fa-bolt fa-3x text-warning"></i>
              </div>
              <h3 class="fw-bold mb-3">{{ t('home.guestTitle') }}</h3>
              <p class="text-muted mb-4 mx-auto" style="max-width: 500px;">{{ t('home.guestSubtitle') }}</p>
              <router-link to="/cv/guest" class="btn btn-warning btn-lg px-5 shadow-sm rounded-pill fw-semibold mb-3">
                <i class="fas fa-bolt me-2"></i>{{ t('cv.createWithoutAccount') }}
              </router-link>
              <div class="mt-3">
                <div class="d-flex align-items-center justify-content-center gap-2 text-muted small">
                  <i class="fas fa-exclamation-triangle text-warning"></i>
                  <span>{{ t('home.guestWarningShort') }}</span>
                </div>
              </div>
              <div class="mt-4 pt-3 border-top">
                <p class="mb-2 text-muted small">{{ t('home.guestOrRegister') }}</p>
                <router-link to="/register" class="btn btn-outline-primary rounded-pill px-4">
                  <i class="fas fa-user-plus me-2"></i>{{ t('home.createAccountBtn') }}
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Final CTA Section (for authenticated users) -->
    <section v-if="authStore.isAuthenticated" class="final-cta-section py-5 bg-light">
      <div class="container py-4 text-center">
        <h3 class="fw-bold mb-3">{{ t('home.readyTitle') }}</h3>
        <p class="text-muted mb-4">{{ t('home.readySubtitle') }}</p>
        <router-link to="/cv/create" class="btn btn-primary btn-lg px-5 rounded-pill shadow-sm fw-semibold">
          <i class="fas fa-plus-circle me-2"></i>{{ t('home.createCVBtn') }}
        </router-link>
      </div>
    </section>

    <!-- Footer -->
    <footer class="home-footer py-5">
      <div class="container">
        <div class="row g-4">
          <div class="col-lg-4 mb-3">
            <div class="d-flex align-items-center mb-3">
              <img v-if="brandingStore.logoUrl" :src="brandingStore.logoUrl" alt="Logo" class="footer-logo me-2" />
              <i v-else class="fas fa-file-alt fa-2x text-primary me-2"></i>
              <h5 class="fw-bold mb-0">{{ brandingStore.systemName }}</h5>
            </div>
            <p class="text-muted small">{{ t('home.footerDesc') }}</p>
          </div>
          <div class="col-6 col-lg-2">
            <h6 class="fw-bold mb-3">{{ t('home.footerProduct') }}</h6>
            <ul class="list-unstyled footer-links">
              <li><router-link to="/">{{ t('app.home') }}</router-link></li>
              <li>
                <router-link :to="authStore.isAuthenticated ? '/cv/create' : '/register'">
                  {{ t('app.createCV') }}
                </router-link>
              </li>
              <li><router-link to="/cv/guest">{{ t('home.tryDemoBtn') }}</router-link></li>
            </ul>
          </div>
          <div class="col-6 col-lg-2">
            <h6 class="fw-bold mb-3">{{ t('home.footerAccount') }}</h6>
            <ul class="list-unstyled footer-links">
              <template v-if="authStore.isAuthenticated">
                <li><router-link to="/dashboard">{{ t('app.dashboard') }}</router-link></li>
                <li><router-link to="/profile">{{ t('app.profile') }}</router-link></li>
              </template>
              <template v-else>
                <li><router-link to="/login">{{ t('app.login') }}</router-link></li>
                <li><router-link to="/register">{{ t('app.register') }}</router-link></li>
              </template>
            </ul>
          </div>
          <div class="col-lg-4">
            <h6 class="fw-bold mb-3">{{ t('home.footerFeatures') }}</h6>
            <div class="d-flex flex-wrap gap-2">
              <span class="badge bg-primary bg-opacity-10 text-primary px-3 py-2 rounded-pill">
                <i class="fas fa-magic me-1"></i>{{ t('home.featureAI') }}
              </span>
              <span class="badge bg-success bg-opacity-10 text-success px-3 py-2 rounded-pill">
                <i class="fas fa-file-pdf me-1"></i>PDF
              </span>
              <span class="badge bg-info bg-opacity-10 text-info px-3 py-2 rounded-pill">
                <i class="fas fa-qrcode me-1"></i>QR
              </span>
              <span class="badge bg-warning bg-opacity-10 text-warning px-3 py-2 rounded-pill">
                <i class="fas fa-language me-1"></i>{{ t('home.featureRTL') }}
              </span>
            </div>
          </div>
        </div>
        <hr class="my-4 opacity-25" />
        <div class="d-flex flex-column flex-md-row justify-content-between align-items-center">
          <p class="text-muted small mb-2 mb-md-0">
            &copy; {{ new Date().getFullYear() }} {{ brandingStore.systemName }}. {{ t('home.footerRights') }}
          </p>
          <div class="d-flex gap-3">
            <a href="#" class="footer-social-link"><i class="fab fa-github"></i></a>
            <a href="#" class="footer-social-link"><i class="fab fa-linkedin"></i></a>
            <a href="#" class="footer-social-link"><i class="fab fa-twitter"></i></a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { useBrandingStore } from '../stores/branding'
import { onMounted } from 'vue'

const { t } = useI18n()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()

const features = [
  {
    icon: 'fas fa-magic',
    titleKey: 'home.featureAITitle',
    descKey: 'home.featureAIDesc',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  },
  {
    icon: 'fas fa-palette',
    titleKey: 'home.featureTemplatesTitle',
    descKey: 'home.featureTemplatesDesc',
    gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
  },
  {
    icon: 'fas fa-language',
    titleKey: 'home.featureLanguageTitle',
    descKey: 'home.featureLanguageDesc',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
  },
  {
    icon: 'fas fa-share-alt',
    titleKey: 'home.featureShareTitle',
    descKey: 'home.featureShareDesc',
    gradient: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
  },
]

const steps = [
  {
    icon: 'fas fa-palette',
    titleKey: 'home.step1Title',
    descKey: 'home.step1Desc',
  },
  {
    icon: 'fas fa-edit',
    titleKey: 'home.step2Title',
    descKey: 'home.step2Desc',
  },
  {
    icon: 'fas fa-download',
    titleKey: 'home.step3Title',
    descKey: 'home.step3Desc',
  },
]

const templates = [
  { name: 'modern', color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { name: 'professional', color: 'linear-gradient(135deg, #2c3e50 0%, #3498db 100%)' },
  { name: 'creative', color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { name: 'elegant', color: 'linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)' },
  { name: 'minimalist', color: 'linear-gradient(135deg, #bdc3c7 0%, #2c3e50 100%)' },
  { name: 'executive', color: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 100%)' },
  { name: 'tech', color: 'linear-gradient(135deg, #00c9ff 0%, #92fe9d 100%)' },
  { name: 'ats', color: 'linear-gradient(135deg, #3a7bd5 0%, #00d2ff 100%)' },
  { name: 'academic', color: 'linear-gradient(135deg, #834d9b 0%, #d04ed6 100%)' },
  { name: 'designer', color: 'linear-gradient(135deg, #fc5c7d 0%, #6a82fb 100%)' },
  { name: 'engineer', color: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)' },
  { name: 'medical', color: 'linear-gradient(135deg, #56ab2f 0%, #a8e063 100%)' },
  { name: 'compact', color: 'linear-gradient(135deg, #614385 0%, #516395 100%)' },
  { name: 'lawyer', color: 'linear-gradient(135deg, #232526 0%, #414345 100%)' },
  { name: 'teacher', color: 'linear-gradient(135deg, #ff6a00 0%, #ee0979 100%)' },
  { name: 'standard', color: 'linear-gradient(135deg, #536976 0%, #292E49 100%)' },
]

const stats = [
  { icon: 'fas fa-file-alt text-primary', value: '10,000+', labelKey: 'home.statCVs' },
  { icon: 'fas fa-palette text-success', value: '16', labelKey: 'home.statTemplates' },
  { icon: 'fas fa-language text-info', value: '2', labelKey: 'home.statLanguages' },
  { icon: 'fas fa-robot text-warning', value: '6', labelKey: 'home.statAITools' },
]

// Intersection Observer for scroll animations
onMounted(() => {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
        }
      })
    },
    { threshold: 0.1, rootMargin: '0px 0px -50px 0px' }
  )

  document.querySelectorAll('.feature-card, .step-card, .template-card, .stat-card, .guest-card').forEach((el) => {
    el.classList.add('animate-on-scroll')
    observer.observe(el)
  })
})
</script>

<style scoped>
/* ===== Hero Section ===== */
.hero-section {
  background: linear-gradient(135deg, #0a1628 0%, #0d3b66 30%, #0d6efd 70%, #4da3ff 100%);
  position: relative;
  overflow: hidden;
  min-height: 85vh;
  display: flex;
  align-items: center;
}

.min-vh-75 {
  min-height: 75vh;
}

.hero-bg-shapes .shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.06;
  background: white;
}

.shape-1 {
  width: 500px;
  height: 500px;
  top: -150px;
  right: -100px;
  animation: floatShape 20s ease-in-out infinite;
}

.shape-2 {
  width: 300px;
  height: 300px;
  bottom: -80px;
  left: -60px;
  animation: floatShape 15s ease-in-out infinite reverse;
}

.shape-3 {
  width: 200px;
  height: 200px;
  top: 40%;
  left: 20%;
  animation: floatShape 18s ease-in-out infinite;
}

.shape-4 {
  width: 150px;
  height: 150px;
  top: 20%;
  right: 30%;
  animation: floatShape 12s ease-in-out infinite reverse;
}

@keyframes floatShape {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(30px, -30px) scale(1.05); }
  50% { transform: translate(-20px, 20px) scale(0.95); }
  75% { transform: translate(10px, -10px) scale(1.02); }
}

.hero-title {
  line-height: 1.15;
  letter-spacing: -0.02em;
}

.hero-btn {
  transition: all 0.3s ease;
  border-radius: 50px;
}

.hero-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.trust-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

/* Hero Visual - CV Preview Card */
.hero-visual {
  position: relative;
  padding: 20px;
}

.cv-preview-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 25px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  transform: perspective(800px) rotateY(-5deg) rotateX(3deg);
  transition: transform 0.5s ease;
}

.cv-preview-card:hover {
  transform: perspective(800px) rotateY(0deg) rotateX(0deg);
}

.cv-preview-header {
  background: #f1f3f5;
  padding: 12px 16px;
  display: flex;
  align-items: center;
}

.cv-dots {
  display: flex;
  gap: 6px;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.dot-red { background: #ff5f57; }
.dot-yellow { background: #febc2e; }
.dot-green { background: #28c840; }

.cv-preview-body {
  padding: 24px;
}

.cv-line {
  height: 8px;
  border-radius: 4px;
  margin-bottom: 8px;
  animation: shimmerLine 2s ease-in-out infinite;
}

.cv-line-title {
  width: 60%;
  height: 14px;
  background: linear-gradient(90deg, #0d6efd, #4da3ff);
}

.cv-line-subtitle {
  width: 40%;
  height: 10px;
  background: #dee2e6;
}

.cv-line-heading {
  width: 35%;
  height: 10px;
  background: #495057;
}

.cv-line-text {
  width: 90%;
  background: #e9ecef;
}

.cv-line-text.short {
  width: 55%;
}

.cv-line-text.medium {
  width: 75%;
}

.cv-section-divider {
  height: 1px;
  background: #e9ecef;
  margin: 12px 0;
}

.cv-skills-row {
  display: flex;
  gap: 6px;
  margin-top: 8px;
}

.cv-skill-chip {
  width: 50px;
  height: 20px;
  border-radius: 10px;
  background: linear-gradient(135deg, #e3f2fd, #bbdefb);
}

@keyframes shimmerLine {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.floating-badge {
  position: absolute;
  padding: 10px 18px;
  border-radius: 12px;
  font-weight: 700;
  font-size: 14px;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  animation: floatBadge 3s ease-in-out infinite;
}

.floating-badge-ai {
  top: 10px;
  right: 0;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  animation-delay: 0s;
}

.floating-badge-pdf {
  bottom: 30px;
  left: 0;
  background: linear-gradient(135deg, #f5576c, #ff6b6b);
  color: white;
  animation-delay: 1.5s;
}

@keyframes floatBadge {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

/* ===== Animations ===== */
.fade-in {
  animation: fadeInUp 0.8s ease forwards;
  opacity: 0;
}

.fade-in:nth-child(1) { animation-delay: 0.1s; }
.fade-in:nth-child(2) { animation-delay: 0.2s; }
.fade-in:nth-child(3) { animation-delay: 0.3s; }
.fade-in:nth-child(4) { animation-delay: 0.4s; }
.fade-in:nth-child(5) { animation-delay: 0.5s; }

.fade-in-right {
  animation: fadeInRight 1s ease forwards;
  animation-delay: 0.5s;
  opacity: 0;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fadeInRight {
  from { opacity: 0; transform: translateX(40px); }
  to { opacity: 1; transform: translateX(0); }
}

/* ===== Section Labels ===== */
.section-label {
  display: inline-block;
  background: linear-gradient(135deg, #0d6efd, #6610f2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  font-weight: 700;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 2px;
}

/* ===== Features Section ===== */
.features-section {
  background: white;
}

.feature-card {
  background: white;
  border-radius: 16px;
  border: 1px solid #f0f0f0;
  transition: all 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  cursor: default;
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.08);
  border-color: transparent;
}

.feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.feature-card:hover .feature-icon {
  transform: scale(1.1) rotate(5deg);
}

/* ===== How It Works ===== */
.step-card {
  background: white;
  border-radius: 16px;
  border: 1px solid #f0f0f0;
  position: relative;
  transition: all 0.4s ease;
}

.step-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.06);
}

.step-number-wrapper {
  position: relative;
  display: inline-block;
}

.step-number {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #0d6efd, #6610f2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  z-index: 1;
  box-shadow: 0 4px 10px rgba(13, 110, 253, 0.3);
}

.step-icon-bg {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #f0f4ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  transition: background 0.3s ease;
}

.step-card:hover .step-icon-bg {
  background: #e0e8ff;
}

.step-arrow {
  position: absolute;
  top: 50%;
  right: -25px;
  transform: translateY(-50%);
  z-index: 2;
}

/* ===== Templates Section ===== */
.template-card {
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.4s ease;
  cursor: pointer;
  border: 2px solid transparent;
}

.template-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
  border-color: #0d6efd;
}

.template-preview {
  height: 140px;
  position: relative;
  padding: 16px;
  display: flex;
  align-items: flex-start;
  overflow: hidden;
}

.template-lines {
  width: 100%;
}

.tl {
  height: 5px;
  border-radius: 3px;
  background: rgba(255, 255, 255, 0.4);
  margin-bottom: 5px;
}

.tl-header {
  width: 55%;
  height: 8px;
  background: rgba(255, 255, 255, 0.7);
}

.tl-sub {
  width: 35%;
  background: rgba(255, 255, 255, 0.5);
}

.tl-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.2);
  margin: 6px 0;
}

.tl-text {
  width: 80%;
}

.tl-text.short {
  width: 50%;
}

.template-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  color: white;
  font-size: 1.5rem;
}

.template-card:hover .template-overlay {
  opacity: 1;
}

.template-name {
  padding: 8px;
  text-align: center;
  background: white;
  border-top: 1px solid #f0f0f0;
}

/* ===== Stats Section ===== */
.stats-section {
  background: linear-gradient(135deg, #0a1628 0%, #0d3b66 50%, #0d6efd 100%);
}

.stat-card {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.stat-card:hover {
  background: rgba(255, 255, 255, 0.12);
  transform: translateY(-4px);
}

.stat-icon {
  color: rgba(255, 255, 255, 0.6);
}

.stat-value {
  color: white;
}

.stat-label {
  color: rgba(255, 255, 255, 0.6) !important;
}

/* ===== Guest CTA Section ===== */
.guest-card {
  background: white;
  border-radius: 20px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.06);
  border: 1px solid #f0f0f0;
}

.guest-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #fff8e1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

/* ===== Footer ===== */
.home-footer {
  background: #0a1628;
  color: #a0aec0;
}

.home-footer h5,
.home-footer h6 {
  color: white;
}

.footer-links li {
  margin-bottom: 8px;
}

.footer-links a {
  color: #a0aec0;
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.footer-links a:hover {
  color: #4da3ff;
}

.footer-social-link {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #a0aec0;
  text-decoration: none;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.footer-social-link:hover {
  background: #0d6efd;
  color: white;
  transform: translateY(-2px);
}

/* ===== Scroll Animations ===== */
.animate-on-scroll {
  opacity: 0;
  transform: translateY(30px);
  transition: all 0.6s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.animate-on-scroll.visible {
  opacity: 1;
  transform: translateY(0);
}

.feature-card.animate-on-scroll:nth-child(1) { transition-delay: 0s; }
.feature-card.animate-on-scroll:nth-child(2) { transition-delay: 0.1s; }
.feature-card.animate-on-scroll:nth-child(3) { transition-delay: 0.2s; }
.feature-card.animate-on-scroll:nth-child(4) { transition-delay: 0.3s; }

.step-card.animate-on-scroll:nth-child(1) { transition-delay: 0s; }
.step-card.animate-on-scroll:nth-child(2) { transition-delay: 0.15s; }
.step-card.animate-on-scroll:nth-child(3) { transition-delay: 0.3s; }

/* ===== Responsive ===== */
@media (max-width: 991.98px) {
  .hero-section {
    min-height: auto;
    padding: 60px 0;
  }

  .min-vh-75 {
    min-height: auto;
  }

  .hero-title {
    font-size: 2.5rem;
  }
}

@media (max-width: 767.98px) {
  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .template-preview {
    height: 100px;
  }

  .stat-value {
    font-size: 1.5rem !important;
  }
}

/* ===== RTL Support ===== */
:deep([dir="rtl"]) .step-arrow {
  right: auto;
  left: -25px;
  transform: translateY(-50%) scaleX(-1);
}

:deep([dir="rtl"]) .cv-preview-card {
  transform: perspective(800px) rotateY(5deg) rotateX(3deg);
}

:deep([dir="rtl"]) .cv-preview-card:hover {
  transform: perspective(800px) rotateY(0deg) rotateX(0deg);
}

:deep([dir="rtl"]) .floating-badge-ai {
  right: auto;
  left: 0;
}

:deep([dir="rtl"]) .floating-badge-pdf {
  left: auto;
  right: 0;
}

.footer-logo {
  height: 32px;
  width: auto;
  object-fit: contain;
}
</style>
