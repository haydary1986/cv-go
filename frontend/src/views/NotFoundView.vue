<template>
  <div class="notfound-page">
    <div class="notfound-content text-center">
      <!-- University Logo -->
      <div class="notfound-logo mb-3">
        <i class="fas fa-university"></i>
      </div>

      <!-- 404 Number -->
      <h1 class="notfound-number">404</h1>

      <!-- Messages -->
      <h3 class="notfound-title mb-2">{{ t('app.pageNotFound') }}</h3>
      <p class="notfound-subtitle mb-1" dir="rtl">الصفحة التي تبحث عنها غير موجودة</p>
      <p class="notfound-subtitle mb-4">The page you are looking for does not exist</p>

      <!-- Go Home Button -->
      <router-link to="/" class="notfound-btn">
        <i class="fas fa-home me-2"></i>{{ t('app.goHome') }}
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<style scoped>
.notfound-page {
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-secondary: #2c3e50;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    linear-gradient(135deg, rgba(26, 82, 118, 0.03) 25%, transparent 25%),
    linear-gradient(225deg, rgba(26, 82, 118, 0.03) 25%, transparent 25%),
    linear-gradient(315deg, rgba(26, 82, 118, 0.03) 25%, transparent 25%),
    linear-gradient(45deg, rgba(26, 82, 118, 0.03) 25%, transparent 25%),
    linear-gradient(180deg, #f4f6f9, #e8ecf1);
  background-size: 60px 60px, 60px 60px, 60px 60px, 60px 60px, 100% 100%;
}

.notfound-content {
  max-width: 480px;
  padding: 2rem;
}

.notfound-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--uni-primary);
  color: var(--uni-accent);
  font-size: 2rem;
  box-shadow: 0 6px 24px rgba(26, 82, 118, 0.2);
  margin-bottom: 0.5rem;
}

.notfound-number {
  font-size: 7rem;
  font-weight: 800;
  line-height: 1;
  color: var(--uni-accent);
  text-shadow: 0 4px 20px rgba(192, 152, 43, 0.15);
  margin-bottom: 0.5rem;
  letter-spacing: -0.04em;
}

.notfound-title {
  color: var(--uni-primary);
  font-weight: 700;
  font-size: 1.5rem;
}

.notfound-subtitle {
  color: #6c757d;
  font-size: 0.95rem;
}

.notfound-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.65rem 2rem;
  border-radius: 999px;
  background: var(--uni-primary);
  color: #fff;
  font-weight: 600;
  font-size: 1rem;
  text-decoration: none;
  border: none;
  box-shadow: 0 4px 14px rgba(26, 82, 118, 0.25);
  transition: all 0.25s ease;
}

.notfound-btn:hover {
  background: #1e6291;
  box-shadow: 0 6px 20px rgba(26, 82, 118, 0.35);
  transform: translateY(-1px);
  color: #fff;
}
</style>
