<template>
  <div :dir="locale === 'ar' ? 'rtl' : 'ltr'" :class="{ 'rtl': locale === 'ar' }">
    <!-- Navbar -->
    <nav class="app-navbar navbar navbar-expand-lg sticky-top"
         :class="{ 'app-navbar--scrolled': isScrolled }">
      <div class="container">
        <router-link class="navbar-brand d-flex align-items-center gap-2" to="/">
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
@import url('https://fonts.googleapis.com/css2?family=Cairo:wght@400;500;600;700&family=Inter:wght@400;500;600;700&display=swap');

:root {
  --font-ar: 'Cairo', sans-serif;
  --font-en: 'Inter', -apple-system, system-ui, sans-serif;
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

body {
  font-family: var(--font-en);
  background-color: var(--uni-bg);
  color: var(--uni-text);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
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
  background: var(--uni-bg) !important;
  border-bottom: none;
  box-shadow: rgba(0,0,0,0.02) 0px 0px 0px 1px;
  transition: box-shadow 0.3s ease;
  padding-top: 0.6rem;
  padding-bottom: 0.6rem;
}

.app-navbar--scrolled {
  box-shadow: rgba(0,0,0,0.02) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 6px, rgba(0,0,0,0.08) 0px 4px 8px;
}

.navbar-logo {
  height: 32px;
  width: auto;
  object-fit: contain;
}

.navbar-brand-text {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--uni-text);
  letter-spacing: 0.01em;
}

.app-navbar .nav-link {
  color: var(--uni-text) !important;
  font-weight: 500;
  font-size: 0.9rem;
  padding: 0.5rem 0.85rem !important;
  border-radius: var(--radius-btn);
  transition: color 0.2s ease, background 0.2s ease;
  text-decoration: none;
}

.app-navbar .nav-link:hover {
  color: var(--uni-primary) !important;
  background: rgba(26, 82, 118, 0.04);
}

.app-navbar .nav-link.router-link-active {
  color: var(--uni-primary) !important;
  background: rgba(26, 82, 118, 0.06);
}

.navbar-toggler {
  border: none !important;
  padding: 0.4rem;
}

.navbar-toggler-icon {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 30 30'%3e%3cpath stroke='%23222222' stroke-linecap='round' stroke-miterlimit='10' stroke-width='2' d='M4 7h22M4 15h22M4 23h22'/%3e%3c/svg%3e") !important;
}

.nav-register-btn {
  background: var(--uni-accent) !important;
  color: #ffffff !important;
  border-radius: 999px !important;
  padding: 0.4rem 1.5rem !important;
  font-weight: 600 !important;
  transition: all 0.2s ease !important;
}

.nav-register-btn:hover {
  background: #d4a82f !important;
  color: #ffffff !important;
  box-shadow: var(--shadow-hover);
  transform: translateY(-1px);
}

/* Notification Bell */
.notif-link {
  position: relative;
}

.notif-link .fas.fa-bell {
  color: var(--uni-text);
  font-size: 1.05rem;
}

.notif-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 18px;
  height: 18px;
  padding: 0 4px;
  border-radius: 999px;
  background: #e53e3e;
  color: #fff;
  font-size: 0.6rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: notif-pulse 2.5s ease-in-out infinite;
  line-height: 1;
}

@keyframes notif-pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(229, 62, 62, 0.4); }
  50% { box-shadow: 0 0 0 6px rgba(229, 62, 62, 0); }
}

/* User Avatar */
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--uni-accent), #d4a82f);
  border: 2px solid transparent;
  background-clip: padding-box;
  position: relative;
  color: #fff;
  font-size: 0.72rem;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  letter-spacing: 0.03em;
  box-shadow: 0 0 0 2px var(--uni-accent);
}

.user-name-text {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--uni-text);
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Dropdown */
.app-dropdown {
  border: none;
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-card);
  padding: 0.5rem;
  margin-top: 0.5rem;
  background: var(--uni-bg);
}

.app-dropdown .dropdown-item {
  border-radius: var(--radius-btn);
  padding: 0.5rem 0.75rem;
  font-size: 0.88rem;
  font-weight: 500;
  color: var(--uni-text);
  transition: background 0.15s;
}

.app-dropdown .dropdown-item:hover {
  background: rgba(26, 82, 118, 0.04);
  color: var(--uni-primary);
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
  color: #ffffff;
  font-size: 0.85rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ====== App Footer ====== */
.app-footer {
  background: var(--uni-bg-light);
  color: var(--uni-text-secondary);
  margin-top: 3rem;
  border-top: 1px solid rgba(0,0,0,0.06);
}

.app-footer-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 0;
}

.app-footer-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--uni-text);
}

.app-footer-logo {
  height: 24px;
  width: auto;
  object-fit: contain;
  opacity: 0.7;
}

.app-footer-copy {
  font-size: 0.8rem;
  color: var(--uni-text-secondary);
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
