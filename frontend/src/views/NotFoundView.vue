<template>
  <div class="notfound-page">
    <div class="notfound-content">
      <!-- 404 Number -->
      <h1 class="notfound-number">404</h1>

      <!-- Messages -->
      <h3 class="notfound-title">{{ t('app.pageNotFound') }}</h3>
      <p class="notfound-subtitle" dir="rtl">الصفحة التي تبحث عنها غير موجودة</p>
      <p class="notfound-subtitle">The page you are looking for does not exist</p>

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
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  padding: 24px;
}

.notfound-content {
  text-align: center;
  max-width: 480px;
}

.notfound-number {
  font-size: 8rem;
  font-weight: 700;
  line-height: 1;
  color: #c0982b;
  margin-bottom: 8px;
  letter-spacing: -0.04em;
}

.notfound-title {
  color: #222222;
  font-weight: 700;
  font-size: 1.5rem;
  margin-bottom: 8px;
}

.notfound-subtitle {
  color: #6a6a6a;
  font-size: 0.95rem;
  margin-bottom: 4px;
}

.notfound-subtitle:last-of-type {
  margin-bottom: 32px;
}

.notfound-btn {
  display: inline-flex;
  align-items: center;
  padding: 12px 28px;
  border-radius: 8px;
  background: #1a5276;
  color: #ffffff;
  font-weight: 600;
  font-size: 1rem;
  text-decoration: none;
  border: none;
  transition: background-color 0.2s, box-shadow 0.2s;
}

.notfound-btn:hover {
  background: #164666;
  box-shadow:
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.12) 0px 4px 12px;
  color: #ffffff;
}
</style>
