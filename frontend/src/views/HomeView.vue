<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <div class="container">
        <div class="row align-items-center hero-row">
          <div class="col-lg-7 py-5">
            <!-- University Logo + Name -->
            <div class="d-flex align-items-center mb-4 fade-in">
              <img
                :src="brandingStore.logoUrl || '/logo.png'"
                alt="AL-Turath University Logo"
                class="hero-university-logo"
              />
              <div class="ms-3">
                <h4 class="mb-0 hero-uni-name">{{ brandingStore.systemName || 'جامعة التراث' }}</h4>
                <small class="hero-uni-subtitle">AL-Turath University</small>
              </div>
            </div>

            <div class="hero-badge mb-3 fade-in">
              <span class="hero-badge-pill">
                <i class="fas fa-star me-1"></i> {{ t('home.heroBadge') }}
              </span>
            </div>
            <h1 class="hero-title fade-in">
              {{ t('home.heroTitle') }}
            </h1>
            <p class="hero-subtitle fade-in">
              {{ t('home.heroSubtitle') }}
            </p>
            <div class="d-flex gap-3 flex-wrap fade-in">
              <template v-if="authStore.isAuthenticated">
                <router-link to="/cv/create" class="btn btn-gold btn-lg hero-btn">
                  <i class="fas fa-plus-circle me-2"></i>{{ t('home.createCVBtn') }}
                </router-link>
                <router-link to="/dashboard" class="btn btn-outline-dark btn-lg hero-btn">
                  <i class="fas fa-th-large me-2"></i>{{ t('home.dashboardBtn') }}
                </router-link>
              </template>
              <template v-else>
                <router-link to="/register" class="btn btn-gold btn-lg hero-btn">
                  <i class="fas fa-rocket me-2"></i>{{ t('home.getStartedBtn') }}
                </router-link>
                <router-link to="/cv/guest" class="btn btn-outline-dark btn-lg hero-btn">
                  <i class="fas fa-eye me-2"></i>{{ t('home.tryDemoBtn') }}
                </router-link>
              </template>
            </div>
            <div class="mt-5 d-flex align-items-center gap-4 hero-trust fade-in">
              <div class="d-flex align-items-center gap-3">
                <div class="hero-stat-item">
                  <i class="fas fa-graduation-cap me-1"></i> 14 {{ t('home.faculties') || 'كلية' }}
                </div>
                <div class="hero-stat-divider"></div>
                <div class="hero-stat-item">
                  <i class="fas fa-palette me-1"></i> 16 {{ t('home.statTemplates') || 'قالب' }}
                </div>
                <div class="hero-stat-divider"></div>
                <div class="hero-stat-item">
                  <i class="fas fa-robot me-1"></i> AI
                </div>
              </div>
            </div>
          </div>
          <div class="col-lg-5 d-none d-lg-block">
            <div class="hero-visual fade-in-right">
              <!-- Large University Logo -->
              <div class="hero-logo-showcase">
                <img :src="brandingStore.logoUrl || '/logo.png'" alt="University Logo" class="hero-logo-large" />
              </div>
              <div class="floating-badge floating-badge-ai">
                <i class="fas fa-magic"></i> AI
              </div>
              <div class="floating-badge floating-badge-pdf">
                <i class="fas fa-file-pdf"></i> PDF
              </div>
              <div class="floating-badge floating-badge-qr">
                <i class="fas fa-qrcode"></i> QR
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features-section" id="features">
      <div class="container">
        <div class="text-center mb-5">
          <span class="section-label">{{ t('home.featuresLabel') }}</span>
          <h2 class="section-title mt-2">{{ t('home.featuresTitle') }}</h2>
          <p class="section-subtitle mx-auto">{{ t('home.featuresSubtitle') }}</p>
        </div>
        <div class="row g-4">
          <div class="col-md-6 col-lg-3" v-for="(feature, index) in features" :key="index">
            <div class="feature-card h-100 p-4 text-center">
              <div class="feature-icon-wrapper mb-3">
                <div class="feature-icon" :style="{ background: feature.gradient }">
                  <i :class="feature.icon" class="fa-lg text-white"></i>
                </div>
              </div>
              <h5 class="feature-title mb-2">{{ t(feature.titleKey) }}</h5>
              <p class="feature-desc mb-0">{{ t(feature.descKey) }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- How It Works Section -->
    <section class="how-it-works-section">
      <div class="container">
        <div class="text-center mb-5">
          <span class="section-label">{{ t('home.howItWorksLabel') }}</span>
          <h2 class="section-title mt-2">{{ t('home.howItWorksTitle') }}</h2>
          <p class="section-subtitle mx-auto">{{ t('home.howItWorksSubtitle') }}</p>
        </div>
        <div class="row g-4 justify-content-center">
          <div class="col-md-4" v-for="(step, index) in steps" :key="index">
            <div class="step-card text-center p-4 h-100">
              <div class="step-number-wrapper mb-3">
                <div class="step-number">{{ index + 1 }}</div>
                <div class="step-icon-bg">
                  <i :class="step.icon" class="fa-2x step-icon-inner"></i>
                </div>
              </div>
              <h5 class="step-title mb-2">{{ t(step.titleKey) }}</h5>
              <p class="step-desc mb-0">{{ t(step.descKey) }}</p>
              <div v-if="index < steps.length - 1" class="step-connector d-none d-md-block">
                <div class="step-connector-line"></div>
                <i class="fas fa-chevron-right step-connector-arrow"></i>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Templates Preview Section -->
    <section class="templates-section">
      <div class="container">
        <div class="text-center mb-5">
          <span class="section-label">{{ t('home.templatesLabel') }}</span>
          <h2 class="section-title mt-2">{{ t('home.templatesTitle') }}</h2>
          <p class="section-subtitle mx-auto">{{ t('home.templatesSubtitle') }}</p>
        </div>
        <div class="row g-3">
          <div class="col-6 col-sm-4 col-md-3 col-lg-2" v-for="tmpl in templates" :key="tmpl.name">
            <div class="template-card">
              <div class="template-preview">
                <img :src="getTemplatePreviewDataUrl(tmpl.name)" :alt="t(`templates.${tmpl.name}`)" loading="lazy" />
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
        <div class="text-center mt-5">
          <router-link
            :to="authStore.isAuthenticated ? '/cv/create' : '/register'"
            class="btn btn-primary-uni btn-lg px-5">
            <i class="fas fa-palette me-2"></i>{{ t('home.browseTemplatesBtn') }}
          </router-link>
        </div>
      </div>
    </section>

    <!-- Statistics Section -->
    <section class="stats-section">
      <div class="container">
        <div class="row g-4 justify-content-center">
          <div class="col-6 col-md-3" v-for="stat in stats" :key="stat.labelKey">
            <div class="stat-card text-center p-4">
              <div class="stat-icon mb-3">
                <i :class="stat.icon" class="fa-2x"></i>
              </div>
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ t(stat.labelKey) }}</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Guest CTA Section -->
    <section v-if="!authStore.isAuthenticated" class="guest-cta-section">
      <div class="container">
        <div class="row justify-content-center">
          <div class="col-lg-8">
            <div class="guest-card p-5 text-center">
              <div class="guest-icon mb-4">
                <i class="fas fa-bolt fa-2x"></i>
              </div>
              <h3 class="guest-title mb-3">{{ t('home.guestTitle') }}</h3>
              <p class="guest-desc mb-4 mx-auto">{{ t('home.guestSubtitle') }}</p>
              <router-link to="/cv/guest" class="btn btn-gold btn-lg px-5 mb-3">
                <i class="fas fa-bolt me-2"></i>{{ t('cv.createWithoutAccount') }}
              </router-link>
              <div class="mt-3">
                <div class="d-flex align-items-center justify-content-center gap-2 guest-warning">
                  <i class="fas fa-exclamation-triangle"></i>
                  <span>{{ t('home.guestWarningShort') }}</span>
                </div>
              </div>
              <div class="mt-4 pt-3 guest-divider">
                <p class="mb-2 guest-or-text">{{ t('home.guestOrRegister') }}</p>
                <router-link to="/register" class="btn btn-outline-primary-uni px-4">
                  <i class="fas fa-user-plus me-2"></i>{{ t('home.createAccountBtn') }}
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Final CTA Section (for authenticated users) -->
    <section v-if="authStore.isAuthenticated" class="final-cta-section">
      <div class="container text-center">
        <h3 class="cta-title mb-3">{{ t('home.readyTitle') }}</h3>
        <p class="cta-subtitle mb-4">{{ t('home.readySubtitle') }}</p>
        <router-link to="/cv/create" class="btn btn-gold btn-lg px-5">
          <i class="fas fa-plus-circle me-2"></i>{{ t('home.createCVBtn') }}
        </router-link>
      </div>
    </section>

    <!-- Footer -->
    <footer class="home-footer">
      <div class="container">
        <div class="row g-4">
          <div class="col-lg-4 mb-3">
            <div class="d-flex align-items-center mb-3">
              <img :src="brandingStore.logoUrl || '/logo.png'" alt="Logo" class="footer-logo me-2" />
              <h5 class="footer-brand-name mb-0">{{ brandingStore.systemName }}</h5>
            </div>
            <p class="footer-desc">{{ t('home.footerDesc') }}</p>
          </div>
          <div class="col-6 col-lg-2">
            <h6 class="footer-heading mb-3">{{ t('home.footerProduct') }}</h6>
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
            <h6 class="footer-heading mb-3">{{ t('home.footerAccount') }}</h6>
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
            <h6 class="footer-heading mb-3">{{ t('home.footerFeatures') }}</h6>
            <div class="d-flex flex-wrap gap-2">
              <span class="footer-feature-badge footer-feature-ai">
                <i class="fas fa-magic me-1"></i>{{ t('home.featureAI') }}
              </span>
              <span class="footer-feature-badge footer-feature-pdf">
                <i class="fas fa-file-pdf me-1"></i>PDF
              </span>
              <span class="footer-feature-badge footer-feature-qr">
                <i class="fas fa-qrcode me-1"></i>QR
              </span>
              <span class="footer-feature-badge footer-feature-rtl">
                <i class="fas fa-language me-1"></i>{{ t('home.featureRTL') }}
              </span>
            </div>
          </div>
        </div>
        <hr class="footer-divider" />
        <div class="d-flex flex-column flex-md-row justify-content-between align-items-center">
          <p class="footer-copyright mb-2 mb-md-0">
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
import { getTemplatePreviewDataUrl } from '../data/templatePreviews'
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
  { icon: 'fas fa-file-alt', value: '10,000+', labelKey: 'home.statCVs' },
  { icon: 'fas fa-palette', value: '16', labelKey: 'home.statTemplates' },
  { icon: 'fas fa-language', value: '2', labelKey: 'home.statLanguages' },
  { icon: 'fas fa-robot', value: '6', labelKey: 'home.statAITools' },
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
/* ===== Design Tokens ===== */
:root {
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-text: #222222;
  --uni-text-secondary: #6a6a6a;
  --uni-bg: #ffffff;
  --uni-bg-light: #f7f7f7;
  --shadow-card: rgba(0,0,0,0.02) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 6px, rgba(0,0,0,0.1) 0px 4px 8px;
  --shadow-hover: rgba(0,0,0,0.08) 0px 4px 12px;
  --radius-btn: 8px;
  --radius-card: 12px;
  --radius-large: 20px;
}

/* ===== Hero Section ===== */
.hero-section {
  background: linear-gradient(180deg, #f0f6fb 0%, #ffffff 100%);
  position: relative;
  overflow: hidden;
  padding: 40px 0 0;
}

.hero-row {
  min-height: 75vh;
}

.hero-university-logo {
  width: 56px;
  height: 56px;
  object-fit: contain;
}

.hero-uni-name {
  font-weight: 700;
  color: var(--uni-text);
  letter-spacing: 0.02em;
}

.hero-uni-subtitle {
  color: var(--uni-text-secondary);
  font-size: 0.85rem;
}

.hero-badge-pill {
  display: inline-block;
  background: rgba(26, 82, 118, 0.08);
  color: var(--uni-primary);
  padding: 6px 16px;
  border-radius: 999px;
  font-size: 0.875rem;
  font-weight: 600;
}

.hero-title {
  font-size: 3.25rem;
  font-weight: 700;
  color: var(--uni-text);
  line-height: 1.15;
  letter-spacing: -0.02em;
  margin-bottom: 16px;
}

.hero-subtitle {
  font-size: 1.15rem;
  color: var(--uni-text-secondary);
  max-width: 540px;
  line-height: 1.6;
  margin-bottom: 32px;
}

.btn-gold {
  background: var(--uni-accent);
  color: #ffffff;
  border: none;
  font-weight: 600;
  border-radius: var(--radius-btn);
  padding: 12px 28px;
  transition: all 0.2s ease;
}

.btn-gold:hover {
  background: #d4a82f;
  color: #ffffff;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(192, 152, 43, 0.3);
}

.btn-outline-dark {
  border: 1.5px solid var(--uni-text);
  color: var(--uni-text);
  font-weight: 600;
  border-radius: var(--radius-btn);
  padding: 12px 28px;
  background: transparent;
  transition: all 0.2s ease;
}

.btn-outline-dark:hover {
  background: var(--uni-text);
  color: #ffffff;
  transform: translateY(-2px);
}

.hero-btn {
  border-radius: var(--radius-btn);
  font-size: 1rem;
}

.hero-stat-item {
  color: var(--uni-text-secondary);
  font-size: 0.9rem;
  font-weight: 500;
}

.hero-stat-divider {
  width: 1px;
  height: 16px;
  background: rgba(0, 0, 0, 0.12);
}

.hero-logo-showcase {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px;
}

.hero-logo-large {
  width: 280px;
  height: 280px;
  object-fit: contain;
  filter: drop-shadow(0 16px 32px rgba(0, 0, 0, 0.08));
  animation: logoFloat 6s ease-in-out infinite;
}

@keyframes logoFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
}

.floating-badge {
  position: absolute;
  padding: 10px 18px;
  border-radius: var(--radius-card);
  font-weight: 700;
  font-size: 14px;
  box-shadow: var(--shadow-card);
  animation: floatBadge 3s ease-in-out infinite;
  background: #ffffff;
}

.floating-badge-ai {
  top: 10px;
  right: 0;
  color: #667eea;
  animation-delay: 0s;
}

.floating-badge-pdf {
  bottom: 30px;
  left: 0;
  color: #f5576c;
  animation-delay: 1.5s;
}

.floating-badge-qr {
  position: absolute;
  bottom: 20%;
  left: 5%;
  color: #009688;
  padding: 8px 14px;
  border-radius: var(--radius-large);
  font-weight: 700;
  font-size: 0.85rem;
  box-shadow: var(--shadow-card);
  animation: floatBadge 4s ease-in-out infinite 2s;
  background: #ffffff;
}

@keyframes floatBadge {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

/* Hero Visual */
.hero-visual {
  position: relative;
  padding: 20px;
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

/* ===== Section Common ===== */
.section-label {
  display: inline-block;
  color: var(--uni-primary);
  font-weight: 700;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 2px;
}

.section-title {
  font-size: 2.25rem;
  font-weight: 700;
  color: var(--uni-text);
  letter-spacing: -0.01em;
}

.section-subtitle {
  color: var(--uni-text-secondary);
  font-size: 1.05rem;
  max-width: 600px;
  line-height: 1.6;
}

/* ===== Features Section ===== */
.features-section {
  background: var(--uni-bg);
  padding: 96px 0;
}

.feature-card {
  background: var(--uni-bg);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  cursor: default;
  border: none;
}

.feature-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-card), var(--shadow-hover);
}

.feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.feature-card:hover .feature-icon {
  transform: scale(1.1);
}

.feature-title {
  font-weight: 600;
  color: var(--uni-text);
  font-size: 1rem;
}

.feature-desc {
  color: var(--uni-text-secondary);
  font-size: 0.875rem;
  line-height: 1.5;
}

/* ===== How It Works ===== */
.how-it-works-section {
  background: var(--uni-bg-light);
  padding: 96px 0;
}

.step-card {
  background: var(--uni-bg);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  position: relative;
  transition: all 0.3s ease;
  border: none;
}

.step-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-card), var(--shadow-hover);
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
  background: var(--uni-primary);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 13px;
  z-index: 1;
  box-shadow: 0 2px 8px rgba(26, 82, 118, 0.3);
}

.step-icon-bg {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: rgba(26, 82, 118, 0.06);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  transition: background 0.3s ease;
}

.step-icon-inner {
  color: var(--uni-primary);
}

.step-card:hover .step-icon-bg {
  background: rgba(26, 82, 118, 0.1);
}

.step-title {
  font-weight: 600;
  color: var(--uni-text);
}

.step-desc {
  color: var(--uni-text-secondary);
  font-size: 0.875rem;
}

.step-connector {
  position: absolute;
  top: 50%;
  right: -25px;
  transform: translateY(-50%);
  z-index: 2;
  display: flex;
  align-items: center;
}

.step-connector-line {
  width: 20px;
  height: 2px;
  background: rgba(26, 82, 118, 0.15);
}

.step-connector-arrow {
  color: var(--uni-primary);
  opacity: 0.3;
  font-size: 0.8rem;
}

/* ===== Templates Section ===== */
.templates-section {
  background: var(--uni-bg);
  padding: 96px 0;
}

.template-card {
  border-radius: var(--radius-card);
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
  box-shadow: var(--shadow-card);
  background: var(--uni-bg);
}

.template-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-card), var(--shadow-hover);
}

.template-preview {
  aspect-ratio: 160 / 220;
  position: relative;
  overflow: hidden;
  background: #f8f9fa;
}

.template-preview img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.template-overlay {
  position: absolute;
  inset: 0;
  background: rgba(26, 82, 118, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  color: #ffffff;
  font-size: 1.5rem;
}

.template-card:hover .template-overlay {
  opacity: 1;
}

.template-name {
  padding: 10px;
  text-align: center;
  background: var(--uni-bg);
  color: var(--uni-text);
  font-size: 0.85rem;
}

.btn-primary-uni {
  background: var(--uni-primary);
  color: #ffffff;
  border: none;
  border-radius: var(--radius-btn);
  font-weight: 600;
  transition: all 0.2s ease;
}

.btn-primary-uni:hover {
  background: #154360;
  color: #ffffff;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(26, 82, 118, 0.25);
}

/* ===== Stats Section ===== */
.stats-section {
  background: var(--uni-bg-light);
  padding: 96px 0;
}

.stat-card {
  border-radius: var(--radius-card);
  background: var(--uni-bg);
  box-shadow: var(--shadow-card);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-card), var(--shadow-hover);
}

.stat-icon {
  color: var(--uni-primary);
  opacity: 0.6;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--uni-primary);
  line-height: 1.2;
}

.stat-label {
  color: var(--uni-text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  margin-top: 4px;
}

/* ===== Guest CTA Section ===== */
.guest-cta-section {
  background: var(--uni-bg);
  padding: 96px 0;
}

.guest-card {
  background: var(--uni-bg);
  border-radius: var(--radius-large);
  box-shadow: var(--shadow-card);
}

.guest-icon {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: rgba(192, 152, 43, 0.08);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--uni-accent);
}

.guest-title {
  font-weight: 700;
  color: var(--uni-text);
}

.guest-desc {
  color: var(--uni-text-secondary);
  max-width: 500px;
  line-height: 1.6;
}

.guest-warning {
  color: var(--uni-text-secondary);
  font-size: 0.85rem;
}

.guest-warning i {
  color: var(--uni-accent);
}

.guest-divider {
  border-top: 1px solid rgba(0,0,0,0.06);
}

.guest-or-text {
  color: var(--uni-text-secondary);
  font-size: 0.875rem;
}

.btn-outline-primary-uni {
  border: 1.5px solid var(--uni-primary);
  color: var(--uni-primary);
  border-radius: var(--radius-btn);
  font-weight: 600;
  background: transparent;
  transition: all 0.2s ease;
}

.btn-outline-primary-uni:hover {
  background: var(--uni-primary);
  color: #ffffff;
}

/* ===== Final CTA ===== */
.final-cta-section {
  background: var(--uni-bg-light);
  padding: 96px 0;
}

.cta-title {
  font-weight: 700;
  color: var(--uni-text);
}

.cta-subtitle {
  color: var(--uni-text-secondary);
}

/* ===== Footer ===== */
.home-footer {
  background: var(--uni-bg-light);
  padding: 64px 0 32px;
  border-top: 1px solid rgba(0,0,0,0.06);
}

.footer-logo {
  width: 36px;
  height: 36px;
  object-fit: contain;
}

.footer-brand-name {
  font-weight: 700;
  color: var(--uni-text);
  font-size: 1rem;
}

.footer-desc {
  color: var(--uni-text-secondary);
  font-size: 0.875rem;
  line-height: 1.6;
}

.footer-heading {
  font-weight: 600;
  color: var(--uni-text);
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.footer-links li {
  margin-bottom: 8px;
}

.footer-links a {
  color: var(--uni-text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: color 0.2s ease;
}

.footer-links a:hover {
  color: var(--uni-primary);
}

.footer-feature-badge {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
}

.footer-feature-ai {
  background: rgba(26, 82, 118, 0.06);
  color: var(--uni-primary);
}

.footer-feature-pdf {
  background: rgba(229, 62, 62, 0.06);
  color: #c53030;
}

.footer-feature-qr {
  background: rgba(0, 150, 136, 0.06);
  color: #009688;
}

.footer-feature-rtl {
  background: rgba(192, 152, 43, 0.06);
  color: var(--uni-accent);
}

.footer-divider {
  border-color: rgba(0, 0, 0, 0.06);
  margin: 32px 0;
}

.footer-copyright {
  color: var(--uni-text-secondary);
  font-size: 0.8rem;
}

.footer-social-link {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--uni-text-secondary);
  text-decoration: none;
  transition: all 0.2s ease;
  font-size: 1rem;
}

.footer-social-link:hover {
  background: var(--uni-primary);
  color: #ffffff;
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
    padding: 40px 0;
  }

  .hero-row {
    min-height: auto;
  }

  .hero-title {
    font-size: 2.5rem;
  }

  .features-section,
  .how-it-works-section,
  .templates-section,
  .stats-section,
  .guest-cta-section,
  .final-cta-section {
    padding: 72px 0;
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
    font-size: 1.75rem;
  }

  .section-title {
    font-size: 1.75rem;
  }

  .features-section,
  .how-it-works-section,
  .templates-section,
  .stats-section,
  .guest-cta-section,
  .final-cta-section {
    padding: 56px 0;
  }
}

/* ===== RTL Support ===== */
:deep([dir="rtl"]) .step-connector {
  right: auto;
  left: -25px;
  transform: translateY(-50%) scaleX(-1);
}

:deep([dir="rtl"]) .floating-badge-ai {
  right: auto;
  left: 0;
}

:deep([dir="rtl"]) .floating-badge-pdf {
  left: auto;
  right: 0;
}
</style>
