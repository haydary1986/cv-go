<template>
  <div class="login-page">
    <div class="login-split">
      <!-- Left Side: Form -->
      <div class="login-form-side">
        <div class="login-form-container">
          <!-- Welcome Text -->
          <h2 class="login-greeting">{{ t('auth.welcomeBack') }}</h2>
          <p class="login-subtitle">{{ t('auth.signInToAccount') }}</p>

          <!-- Error Alert -->
          <div v-if="error" class="alert alert-danger d-flex align-items-center border-0 login-alert" role="alert">
            <i class="fas fa-exclamation-circle me-2 flex-shrink-0"></i>
            <div>{{ error }}</div>
          </div>

          <form @submit.prevent="handleLogin">
            <!-- Email -->
            <div class="mb-4">
              <label class="form-label login-label">{{ t('auth.email') }}</label>
              <div class="login-input-wrap">
                <i class="fas fa-envelope login-input-icon"></i>
                <input
                  type="email"
                  class="form-control login-input"
                  v-model="form.email"
                  :placeholder="t('auth.email')"
                  required
                  autocomplete="email"
                />
              </div>
            </div>

            <!-- Password -->
            <div class="mb-4">
              <label class="form-label login-label">{{ t('auth.password') }}</label>
              <div class="login-input-wrap">
                <i class="fas fa-lock login-input-icon"></i>
                <input
                  :type="showPassword ? 'text' : 'password'"
                  class="form-control login-input login-input-password"
                  v-model="form.password"
                  :placeholder="t('auth.password')"
                  required
                  autocomplete="current-password"
                />
                <button
                  type="button"
                  class="login-toggle-pw"
                  @click="showPassword = !showPassword"
                  :title="showPassword ? t('auth.hidePassword') : t('auth.showPassword')"
                >
                  <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                </button>
              </div>
            </div>

            <!-- Submit Button -->
            <button type="submit" class="btn login-submit-btn w-100 fw-semibold" :disabled="loading">
              <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              <i v-else class="fas fa-sign-in-alt me-2"></i>
              {{ loading ? t('app.loading') : t('auth.loginBtn') }}
            </button>
          </form>

          <!-- Divider -->
          <div class="login-divider">
            <span>{{ t('auth.orContinueWith') }}</span>
          </div>

          <!-- Google Login -->
          <button @click="handleGoogleLogin" class="btn login-google-btn w-100">
            <i class="fab fa-google me-2"></i>{{ t('auth.googleLogin') }}
          </button>

          <!-- Register Link -->
          <div class="login-register-link">
            <span>{{ t('auth.noAccount') }}</span>
            <router-link to="/register" class="fw-semibold">
              {{ t('auth.registerBtn') }}
            </router-link>
          </div>
        </div>
      </div>

      <!-- Right Side: University Branding -->
      <div class="login-brand-side">
        <div class="login-brand-content">
          <div class="login-brand-logo">
            <img :src="brandingStore.logoUrl || '/logo.svg'" alt="Logo" />
          </div>
          <h1 class="login-brand-name">{{ brandingStore.systemName }}</h1>
          <p class="login-brand-tagline">{{ t('auth.loginSubtitle') }}</p>
          <div class="login-brand-decoration">
            <div class="login-brand-line"></div>
            <i class="fas fa-graduation-cap login-brand-icon-accent"></i>
            <div class="login-brand-line"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { useBrandingStore } from '../stores/branding'
import { authAPI } from '../services/api'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()

const loading = ref(false)
const error = ref('')
const showPassword = ref(false)
const form = reactive({ email: '', password: '' })

async function handleLogin() {
  loading.value = true
  error.value = ''
  try {
    await authStore.login(form.email, form.password)
    const rawRedirect = route.query.redirect as string
    const redirect = (rawRedirect && rawRedirect.startsWith('/') && !rawRedirect.startsWith('//')) ? rawRedirect : '/dashboard'
    router.push(redirect)
  } catch (err: any) {
    error.value = err.response?.data?.error || t('auth.invalidCredentials')
  } finally {
    loading.value = false
  }
}

function handleGoogleLogin() {
  authAPI.googleLogin()
}
</script>

<style scoped>
/* ── Layout ── */
.login-page {
  min-height: 100vh;
}

.login-split {
  display: flex;
  min-height: 100vh;
}

/* ── Left: Form Side ── */
.login-form-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  background: #fff;
  animation: loginFadeIn 0.7s ease-out both;
}

@keyframes loginFadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-form-container {
  width: 100%;
  max-width: 420px;
}

.login-greeting {
  font-size: 2rem;
  font-weight: 800;
  color: #1a5276;
  margin-bottom: 0.25rem;
  letter-spacing: -0.5px;
}

.login-subtitle {
  color: #6c757d;
  font-size: 1rem;
  margin-bottom: 2rem;
}

.login-alert {
  border-radius: 10px;
  background: #fef2f2;
  color: #b91c1c;
  margin-bottom: 1.5rem;
}

/* ── Inputs ── */
.login-label {
  font-weight: 600;
  font-size: 0.875rem;
  color: #2c3e50;
  margin-bottom: 0.4rem;
}

.login-input-wrap {
  position: relative;
}

.login-input-icon {
  position: absolute;
  top: 50%;
  inset-inline-start: 16px;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 0.95rem;
  pointer-events: none;
  z-index: 2;
}

.login-input {
  height: 52px;
  padding-inline-start: 46px;
  padding-inline-end: 16px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  font-size: 1rem;
  background: #f9fafb;
  transition: border-color 0.2s, background-color 0.2s, box-shadow 0.2s;
}

.login-input-password {
  padding-inline-end: 46px;
}

.login-input:focus {
  border-color: #1a5276;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
  outline: none;
}

.login-toggle-pw {
  position: absolute;
  top: 50%;
  inset-inline-end: 14px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
  z-index: 2;
  transition: color 0.15s;
}

.login-toggle-pw:hover {
  color: #1a5276;
}

/* ── Submit ── */
.login-submit-btn {
  height: 52px;
  font-size: 1.05rem;
  border-radius: 12px;
  background: linear-gradient(135deg, #c0982b 0%, #d4a933 100%);
  border: none;
  color: #fff;
  letter-spacing: 0.3px;
  transition: transform 0.15s, box-shadow 0.2s;
  margin-top: 0.5rem;
}

.login-submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(192, 152, 43, 0.35);
  background: linear-gradient(135deg, #b08a24 0%, #c0982b 100%);
  color: #fff;
}

.login-submit-btn:disabled {
  opacity: 0.65;
  color: #fff;
}

/* ── Divider ── */
.login-divider {
  position: relative;
  text-align: center;
  margin: 1.75rem 0;
}

.login-divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  border-top: 1px solid #e5e7eb;
}

.login-divider span {
  position: relative;
  padding: 0 16px;
  background: #fff;
  color: #9ca3af;
  font-size: 0.85rem;
}

/* ── Google ── */
.login-google-btn {
  height: 48px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  background: #fff;
  color: #374151;
  font-weight: 500;
  transition: all 0.2s;
}

.login-google-btn:hover {
  border-color: #1a5276;
  background: #f0f7fc;
  color: #1a5276;
}

/* ── Register Link ── */
.login-register-link {
  text-align: center;
  margin-top: 2rem;
  color: #6c757d;
  font-size: 0.95rem;
}

.login-register-link a {
  color: #1a5276;
  text-decoration: none;
  margin-inline-start: 6px;
  transition: color 0.15s;
}

.login-register-link a:hover {
  color: #c0982b;
}

/* ── Right: Branding Side ── */
.login-brand-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(160deg, #1a5276 0%, #0d6efd 100%);
  position: relative;
  overflow: hidden;
}

.login-brand-side::before {
  content: '';
  position: absolute;
  top: -30%;
  right: -20%;
  width: 500px;
  height: 500px;
  border-radius: 50%;
  background: rgba(192, 152, 43, 0.08);
}

.login-brand-side::after {
  content: '';
  position: absolute;
  bottom: -20%;
  left: -15%;
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
}

.login-brand-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 3rem;
  animation: loginBrandFade 0.9s ease-out 0.3s both;
}

@keyframes loginBrandFade {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.login-brand-logo {
  margin-bottom: 2rem;
}

.login-brand-logo img {
  height: 120px;
  width: auto;
  object-fit: contain;
  filter: drop-shadow(0 4px 20px rgba(0, 0, 0, 0.15));
}

.login-brand-name {
  font-size: 2rem;
  font-weight: 800;
  color: #fff;
  margin-bottom: 0.75rem;
  letter-spacing: -0.3px;
}

.login-brand-tagline {
  color: rgba(255, 255, 255, 0.75);
  font-size: 1.1rem;
  max-width: 320px;
  margin: 0 auto 2rem;
  line-height: 1.6;
}

.login-brand-decoration {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.login-brand-line {
  width: 48px;
  height: 2px;
  background: rgba(192, 152, 43, 0.6);
  border-radius: 1px;
}

.login-brand-icon-accent {
  color: #c0982b;
  font-size: 1.25rem;
}

/* ── Responsive ── */
@media (max-width: 991.98px) {
  .login-brand-side {
    display: none;
  }

  .login-form-side {
    padding: 2rem 1.5rem;
  }
}

@media (max-width: 576px) {
  .login-greeting {
    font-size: 1.5rem;
  }

  .login-form-container {
    max-width: 100%;
  }
}

/* ── RTL ── */
[dir="rtl"] .login-split {
  flex-direction: row-reverse;
}
</style>
