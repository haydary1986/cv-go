<template>
  <div :dir="locale === 'ar' ? 'rtl' : 'ltr'" :class="{ 'rtl': locale === 'ar' }">
    <!-- Navbar -->
    <nav class="navbar navbar-expand-lg navbar-dark sticky-top shadow-sm"
         :style="{ backgroundColor: brandingStore.branding.primary_color || '#0d6efd' }">
      <div class="container">
        <router-link class="navbar-brand d-flex align-items-center fw-bold" to="/">
          <img v-if="brandingStore.logoUrl" :src="brandingStore.logoUrl" alt="Logo" class="navbar-logo me-2" />
          <i v-else class="fas fa-file-alt me-2"></i>
          {{ brandingStore.systemName }}
        </router-link>

        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
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

          <ul class="navbar-nav">
            <!-- Language Switcher -->
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" data-bs-toggle="dropdown">
                <i class="fas fa-globe me-1"></i>{{ t('app.language') }}
              </a>
              <ul class="dropdown-menu dropdown-menu-end">
                <li><a class="dropdown-item" href="#" @click.prevent="switchLocale('ar')">{{ t('app.arabic') }}</a></li>
                <li><a class="dropdown-item" href="#" @click.prevent="switchLocale('en')">{{ t('app.english') }}</a></li>
              </ul>
            </li>

            <template v-if="authStore.isAuthenticated">
              <!-- Notifications -->
              <li class="nav-item">
                <router-link to="/notifications" class="nav-link position-relative">
                  <i class="fas fa-bell"></i>
                  <span v-if="notifStore.unreadCount > 0"
                    class="position-absolute top-0 start-100 translate-middle badge rounded-pill bg-danger"
                    style="font-size: 0.6rem">
                    {{ notifStore.unreadCount }}
                  </span>
                </router-link>
              </li>

              <!-- User Menu -->
              <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" href="#" data-bs-toggle="dropdown">
                  <i class="fas fa-user me-1"></i>
                  {{ authStore.user?.full_name_ar || authStore.user?.full_name_en || authStore.user?.email }}
                </a>
                <ul class="dropdown-menu dropdown-menu-end">
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
                <router-link to="/register" class="nav-link">{{ t('app.register') }}</router-link>
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
    <footer v-if="route.name !== 'home'" class="bg-dark text-white text-center py-3 mt-5">
      <div class="container">
        <p class="mb-0">&copy; {{ new Date().getFullYear() }} {{ brandingStore.systemName }}</p>
      </div>
    </footer>

    <ToastNotification />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'
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
}

body {
  font-family: var(--font-en);
  background-color: #f8f9fa;
}

[dir="rtl"] body,
.rtl {
  font-family: var(--font-ar);
}

.rtl .me-1 { margin-right: 0 !important; margin-left: 0.25rem !important; }
.rtl .me-2 { margin-right: 0 !important; margin-left: 0.5rem !important; }
.rtl .ms-2 { margin-left: 0 !important; margin-right: 0.5rem !important; }
.rtl .dropdown-menu-end { left: 0; right: auto; }

.navbar-logo {
  height: 32px;
  width: auto;
  object-fit: contain;
}

@media print {
  .navbar, footer, .btn { display: none !important; }
  .cv-preview-container { box-shadow: none !important; margin: 0 !important; }
}
</style>
