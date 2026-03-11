<template>
  <div class="home">
    <!-- Hero Section -->
    <section class="hero-section text-center py-5 bg-primary text-white">
      <div class="container">
        <h1 class="display-4 fw-bold mb-3">{{ t('app.name') }}</h1>
        <p class="lead mb-4" v-if="locale === 'ar'">أنشئ سيرتك الذاتية الاحترافية بسهولة مع 14 قالب مميز ودعم الذكاء الاصطناعي</p>
        <p class="lead mb-4" v-else>Create your professional CV easily with 14 stunning templates and AI support</p>
        <div class="d-flex gap-3 justify-content-center flex-wrap">
          <router-link v-if="!authStore.isAuthenticated" to="/register" class="btn btn-light btn-lg px-4">
            {{ t('auth.registerBtn') }}
          </router-link>
          <router-link v-if="!authStore.isAuthenticated" to="/login" class="btn btn-outline-light btn-lg px-4">
            {{ t('auth.loginBtn') }}
          </router-link>
          <router-link v-if="authStore.isAuthenticated" to="/cv/create" class="btn btn-light btn-lg px-4">
            {{ t('app.createCV') }}
          </router-link>
          <router-link v-if="authStore.isAuthenticated" to="/dashboard" class="btn btn-outline-light btn-lg px-4">
            {{ t('app.dashboard') }}
          </router-link>
        </div>
        <!-- Guest CTA -->
        <div v-if="!authStore.isAuthenticated" class="mt-4">
          <router-link to="/cv/guest" class="btn btn-warning btn-lg px-4 shadow">
            <i class="fas fa-bolt me-2"></i>{{ t('cv.createWithoutAccount') }}
          </router-link>
          <p class="mt-2 small opacity-75">{{ t('cv.createWithoutAccountDesc') }}</p>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-5">
      <div class="container">
        <h2 class="text-center mb-5" v-if="locale === 'ar'">المميزات الرئيسية</h2>
        <h2 class="text-center mb-5" v-else>Key Features</h2>
        <div class="row g-4">
          <div class="col-md-4" v-for="feature in features" :key="feature.icon">
            <div class="card h-100 border-0 shadow-sm text-center p-4">
              <div class="card-body">
                <i :class="feature.icon" class="fa-3x text-primary mb-3"></i>
                <h5 class="card-title">{{ locale === 'ar' ? feature.titleAr : feature.titleEn }}</h5>
                <p class="card-text text-muted">{{ locale === 'ar' ? feature.descAr : feature.descEn }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Templates Preview -->
    <section class="py-5 bg-light">
      <div class="container">
        <h2 class="text-center mb-5" v-if="locale === 'ar'">14 قالب احترافي</h2>
        <h2 class="text-center mb-5" v-else>14 Professional Templates</h2>
        <div class="row g-3">
          <div class="col-6 col-md-3" v-for="tmpl in templates" :key="tmpl">
            <div class="card border-0 shadow-sm text-center p-3">
              <i class="fas fa-file-alt fa-2x text-primary mb-2"></i>
              <small class="fw-bold">{{ t(`templates.${tmpl}`) }}</small>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'

const { t, locale } = useI18n()
const authStore = useAuthStore()

const templates = [
  'academic', 'ats', 'compact', 'creative', 'designer', 'elegant', 'engineer',
  'executive', 'lawyer', 'medical', 'minimalist', 'modern', 'professional', 'standard',
  'teacher', 'tech'
]

const features = [
  {
    icon: 'fas fa-magic',
    titleAr: 'ذكاء اصطناعي متقدم',
    titleEn: 'Advanced AI',
    descAr: 'تحسين النصوص وتحليل السيرة وإنشاء رسائل التغطية',
    descEn: 'Improve texts, analyze CVs, and generate cover letters',
  },
  {
    icon: 'fas fa-palette',
    titleAr: '14 قالب احترافي',
    titleEn: '14 Professional Templates',
    descAr: 'قوالب متنوعة تناسب جميع التخصصات والمجالات',
    descEn: 'Diverse templates suitable for all specializations',
  },
  {
    icon: 'fas fa-language',
    titleAr: 'دعم ثنائي اللغة',
    titleEn: 'Bilingual Support',
    descAr: 'دعم كامل للعربية والإنجليزية مع RTL',
    descEn: 'Full Arabic and English support with RTL',
  },
  {
    icon: 'fas fa-file-pdf',
    titleAr: 'تصدير PDF احترافي',
    titleEn: 'Professional PDF Export',
    descAr: 'تصدير عالي الجودة مع دعم صفحات متعددة',
    descEn: 'High-quality export with multi-page support',
  },
  {
    icon: 'fas fa-share-alt',
    titleAr: 'مشاركة سهلة',
    titleEn: 'Easy Sharing',
    descAr: 'روابط مشاركة فريدة مع رمز QR وتتبع المشاهدات',
    descEn: 'Unique share links with QR codes and view tracking',
  },
  {
    icon: 'fas fa-shield-alt',
    titleAr: 'أمان متقدم',
    titleEn: 'Advanced Security',
    descAr: 'تشفير البيانات وحماية من الهجمات ونظام صلاحيات',
    descEn: 'Data encryption, attack protection, and RBAC',
  },
]
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, #0d6efd 0%, #0a58ca 100%);
  min-height: 400px;
  display: flex;
  align-items: center;
}
</style>
