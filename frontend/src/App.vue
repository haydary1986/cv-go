<template>
  <div :dir="locale === 'ar' ? 'rtl' : 'ltr'" :class="{ 'rtl': locale === 'ar' }">
    <!-- Navbar -->
    <nav class="app-navbar navbar navbar-expand-lg navbar-dark sticky-top"
         :class="{ 'app-navbar--scrolled': isScrolled }">
      <div class="container">
        <router-link class="navbar-brand d-flex align-items-center gap-2 fw-bold" to="/">
          <img :src="brandingStore.logoUrl || '/logo.svg'" alt="Logo" class="navbar-logo" />
          <span class="navbar-brand-text">{{ brandingStore.systemName }}</span>
        </router-link>

        <button class="navbar-toggler border-0" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
          <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav me-auto">
            <li class="nav-item">
              <router-link to="/" class="nav-link">{{ t('app.home') }}</router-link>
            </li>
            <template v-if="authStore.isAuthenticated">
              <li class="nav-item">
                <router-link to="/dashboard" class="nav-link">{{ t('app.dashboard') }}</router-link>
              </li>
              <li class="nav-item">
                <router-link to="/cv/create" class="nav-link">{{ t('app.createCV') }}</router-link>
              </li>
              <li class="nav-item" v-if="authStore.isAdmin">
                <router-link to="/admin" class="nav-link">{{ t('app.admin') }}</router-link>
              </li>
            </template>
          </ul>

          <ul class="navbar-nav align-items-lg-center">
            <!-- Language Switcher -->
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" data-bs-toggle="dropdown">
                <i class="fas fa-globe me-1"></i>{{ t('app.language') }}
              </a>
              <ul class="dropdown-menu dropdown-menu-end app-dropdown">
                <li><a class="dropdown-item" href="#" @click.prevent="switchLocale('ar')">{{ t('app.arabic') }}</a></li>
                <li><a class="dropdown-item" href="#" @click.prevent="switchLocale('en')">{{ t('app.english') }}</a></li>
              </ul>
            </li>

            <template v-if="authStore.isAuthenticated">
              <!-- Notifications -->
              <li class="nav-item">
                <router-link to="/notifications" class="nav-link position-relative notif-link">
                  <i class="fas fa-bell"></i>
                  <span v-if="notifStore.unreadCount > 0"
                    class="notif-badge">
                    {{ notifStore.unreadCount }}
                  </span>
                </router-link>
              </li>

              <!-- User Menu -->
              <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle d-flex align-items-center gap-2" href="#" data-bs-toggle="dropdown">
                  <span class="user-avatar">
                    {{ userInitials }}
                  </span>
                  <span class="d-none d-lg-inline user-name-text">
                    {{ authStore.user?.full_name_ar || authStore.user?.full_name_en || authStore.user?.email }}
                  </span>
                </a>
                <ul class="dropdown-menu dropdown-menu-end app-dropdown user-dropdown">
                  <li class="dropdown-header user-dropdown-header">
                    <div class="user-dropdown-avatar">{{ userInitials }}</div>
                    <div>
                      <div class="fw-bold">{{ authStore.user?.full_name_ar || authStore.user?.full_name_en }}</div>
                      <small class="text-muted">{{ authStore.user?.email }}</small>
                    </div>
                  </li>
                  <li><hr class="dropdown-divider" /></li>
                  <li>
                    <router-link to="/profile" class="dropdown-item">
                      <i class="fas fa-user-cog me-2"></i>{{ t('app.profile') }}
                    </router-link>
                  </li>
                  <li><hr class="dropdown-divider" /></li>
                  <li>
                    <a class="dropdown-item text-danger" href="#" @click.prevent="handleLogout">
                      <i class="fas fa-sign-out-alt me-2"></i>{{ t('app.logout') }}
                    </a>
                  </li>
                </ul>
              </li>
            </template>

            <template v-else>
              <li class="nav-item">
                <router-link to="/login" class="nav-link">{{ t('app.login') }}</router-link>
              </li>
              <li class="nav-item">
                <router-link to="/register" class="nav-link nav-register-btn">{{ t('app.register') }}</router-link>
              </li>
            </template>
          </ul>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main>
      <router-view />
    </main>

    <!-- Footer (hidden on home page which has its own footer) -->
    <footer v-if="route.name !== 'home'" class="app-footer">
      <div class="container">
        <div class="app-footer-inner">
          <div class="app-footer-brand">
            <img :src="brandingStore.logoUrl || '/logo.svg'" alt="Logo" class="app-footer-logo" />
            <span>{{ brandingStore.systemName }}</span>
          </div>
          <p class="app-footer-copy mb-0">&copy; {{ new Date().getFullYear() }} {{ brandingStore.systemName }}</p>
        </div>
      </div>
    </footer>

    <ToastNotification />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from './stores/auth'
import { useNotificationStore } from './stores/notification'
import { useBrandingStore } from './stores/branding'
import ToastNotification from './components/ToastNotification.vue'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const notifStore = useNotificationStore()
const brandingStore = useBrandingStore()

const isScrolled = ref(false)

const userInitials = computed(() => {
  const name = authStore.user?.full_name_ar || authStore.user?.full_name_en || authStore.user?.email || ''
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.substring(0, 2).toUpperCase()
})

function handleScroll() {
  isScrolled.value = window.scrollY > 10
}

let notifInterval: ReturnType<typeof setInterval> | null = null

function startNotificationPolling() {
  stopNotificationPolling()
  notifStore.fetchUnreadCount()
  notifInterval = setInterval(() => {
    notifStore.fetchUnreadCount()
  }, 60000)
}

function stopNotificationPolling() {
  if (notifInterval) {
    clearInterval(notifInterval)
    notifInterval = null
  }
}

function switchLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('locale', lang)
  document.documentElement.dir = lang === 'ar' ? 'rtl' : 'ltr'
  document.documentElement.lang = lang
}

function handleLogout() {
  authStore.logout()
  router.push('/')
}

onMounted(async () => {
  window.addEventListener('scroll', handleScroll, { passive: true })
  const savedLocale = localStorage.getItem('locale') || 'ar'
  switchLocale(savedLocale)

  await brandingStore.fetchBranding()
  document.title = brandingStore.systemName

  if (authStore.isAuthenticated) {
    await authStore.fetchProfile()
    startNotificationPolling()
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  stopNotificationPolling()
})

watch(() => authStore.isAuthenticated, (isAuth) => {
  if (isAuth) {
    startNotificationPolling()
  } else {
    stopNotificationPolling()
  }
})
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Cairo:wght@300;400;500;600;700&family=Inter:wght@300;400;500;600;700&display=swap');

:root {
  --font-ar: 'Cairo', sans-serif;
  --font-en: 'Inter', sans-serif;
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-secondary: #2c3e50;
}

body {
  font-family: var(--font-en);
  background-color: #f4f6f9;
}

[dir="rtl"] body,
.rtl {
  font-family: var(--font-ar);
}

.rtl .me-1 { margin-right: 0 !important; margin-left: 0.25rem !important; }
.rtl .me-2 { margin-right: 0 !important; margin-left: 0.5rem !important; }
.rtl .ms-2 { margin-left: 0 !important; margin-right: 0.5rem !important; }
.rtl .dropdown-menu-end { left: 0; right: auto; }

/* ====== App Navbar ====== */
.app-navbar {
  background: rgba(26, 82, 118, 0.92) !important;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 2px solid var(--uni-accent);
  transition: background 0.3s ease, box-shadow 0.3s ease;
  padding-top: 0.55rem;
  padding-bottom: 0.55rem;
}

.app-navbar--scrolled {
  background: var(--uni-primary) !important;
  box-shadow: 0 4px 20px rgba(26, 82, 118, 0.25);
}

.navbar-logo {
  height: 34px;
  width: auto;
  object-fit: contain;
  filter: brightness(0) invert(1);
}

.navbar-brand-text {
  font-size: 1.05rem;
  letter-spacing: 0.01em;
}

.app-navbar .nav-link {
  font-weight: 500;
  font-size: 0.9rem;
  padding: 0.5rem 0.85rem !important;
  border-radius: 8px;
  transition: background 0.2s;
}

.app-navbar .nav-link:hover {
  background: rgba(255, 255, 255, 0.1);
}

.app-navbar .nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.15);
}

.nav-register-btn {
  background: var(--uni-accent) !important;
  color: #fff !important;
  border-radius: 999px !important;
  padding: 0.4rem 1.2rem !important;
  font-weight: 600 !important;
}

.nav-register-btn:hover {
  background: #d4a82f !important;
}

/* Notification Bell */
.notif-link {
  position: relative;
}

.notif-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 18px;
  height: 18px;
  padding: 0 4px;
  border-radius: 999px;
  background: #dc3545;
  color: #fff;
  font-size: 0.6rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: notif-pulse 2s ease-in-out infinite;
  line-height: 1;
}

@keyframes notif-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(220, 53, 69, 0.5); }
  50% { box-shadow: 0 0 0 6px rgba(220, 53, 69, 0); }
}

/* User Avatar */
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--uni-accent);
  color: #fff;
  font-size: 0.72rem;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  letter-spacing: 0.03em;
}

.user-name-text {
  font-size: 0.85rem;
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Dropdown */
.app-dropdown {
  border: none;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  padding: 0.5rem;
  margin-top: 0.5rem;
}

.app-dropdown .dropdown-item {
  border-radius: 8px;
  padding: 0.5rem 0.75rem;
  font-size: 0.88rem;
  transition: background 0.15s;
}

.app-dropdown .dropdown-item:hover {
  background: rgba(26, 82, 118, 0.06);
}

.user-dropdown {
  min-width: 260px;
}

.user-dropdown-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
}

.user-dropdown-avatar {
  width: 42px;
  height: 42px;
  min-width: 42px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--uni-primary), #1e6291);
  color: var(--uni-accent);
  font-size: 0.85rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ====== App Footer ====== */
.app-footer {
  background: var(--uni-primary);
  color: rgba(255, 255, 255, 0.8);
  margin-top: 3rem;
  border-top: 3px solid var(--uni-accent);
}

.app-footer-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 0;
}

.app-footer-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 0.9rem;
}

.app-footer-logo {
  height: 24px;
  width: auto;
  object-fit: contain;
  filter: brightness(0) invert(1);
  opacity: 0.8;
}

.app-footer-copy {
  font-size: 0.8rem;
  opacity: 0.7;
}

@media (max-width: 576px) {
  .app-footer-inner {
    flex-direction: column;
    gap: 0.5rem;
    text-align: center;
  }
}

@media print {
  .app-navbar, .app-footer, .btn { display: none !important; }
  .cv-preview-container { box-shadow: none !important; margin: 0 !important; }
}
</style>
