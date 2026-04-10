<template>
  <div class="login-page">
    <div class="login-card">
      <!-- Logo -->
      <div class="login-logo">
        <img :src="brandingStore.logoUrl || '/logo.svg'" alt="Logo" />
      </div>

      <!-- Welcome Text -->
      <h2 class="login-greeting">{{ t('auth.welcomeBack') }}</h2>
      <p class="login-subtitle">{{ t('auth.signInToAccount') }}</p>

      <!-- Error Alert -->
      <div v-if="error" class="login-alert" role="alert">
        <i class="fas fa-exclamation-circle"></i>
        <span>{{ error }}</span>
      </div>

      <form @submit.prevent="handleLogin">
        <!-- Email -->
        <div class="login-field">
          <label class="login-label">{{ t('auth.email') }}</label>
          <div class="login-input-wrap">
            <input
              type="email"
              class="login-input"
              v-model="form.email"
              :placeholder="t('auth.email')"
              required
              autocomplete="email"
            />
          </div>
        </div>

        <!-- Password -->
        <div class="login-field">
          <label class="login-label">{{ t('auth.password') }}</label>
          <div class="login-input-wrap">
            <input
              :type="showPassword ? 'text' : 'password'"
              class="login-input login-input-password"
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
        <button type="submit" class="login-submit-btn" :disabled="loading">
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
      <button @click="handleGoogleLogin" class="login-google-btn">
        <i class="fab fa-google me-2"></i>{{ t('auth.googleLogin') }}
      </button>

      <!-- Register Link -->
      <div class="login-register-link">
        <span>{{ t('auth.noAccount') }}</span>
        <router-link to="/register">
          {{ t('auth.registerBtn') }}
        </router-link>
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
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 480px;
  background: #ffffff;
  border-radius: 20px;
  padding: 48px 40px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
  animation: loginFadeIn 0.5s ease-out both;
}

@keyframes loginFadeIn {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Logo */
.login-logo {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo img {
  height: 48px;
  width: auto;
  object-fit: contain;
}

/* Heading */
.login-greeting {
  font-size: 1.75rem;
  font-weight: 700;
  color: #222222;
  margin-bottom: 4px;
  text-align: center;
}

.login-subtitle {
  color: #6a6a6a;
  font-size: 0.95rem;
  margin-bottom: 32px;
  text-align: center;
}

/* Alert */
.login-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 12px;
  background: #fef2f2;
  color: #b91c1c;
  font-size: 0.9rem;
  margin-bottom: 24px;
}

.login-alert i {
  flex-shrink: 0;
}

/* Fields */
.login-field {
  margin-bottom: 20px;
}

.login-label {
  display: block;
  font-weight: 500;
  font-size: 0.875rem;
  color: #222222;
  margin-bottom: 6px;
}

.login-input-wrap {
  position: relative;
}

.login-input {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  border: 1px solid #c1c1c1;
  border-radius: 12px;
  font-size: 0.95rem;
  color: #222222;
  background: #ffffff;
  transition: border-color 0.2s, box-shadow 0.2s;
  outline: none;
}

.login-input-password {
  padding-inline-end: 44px;
}

.login-input::placeholder {
  color: #b0b0b0;
}

.login-input:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 2px rgba(26, 82, 118, 0.15);
}

.login-toggle-pw {
  position: absolute;
  top: 50%;
  inset-inline-end: 14px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #b0b0b0;
  cursor: pointer;
  padding: 4px;
  z-index: 2;
  transition: color 0.15s;
}

.login-toggle-pw:hover {
  color: #1a5276;
}

/* Submit */
.login-submit-btn {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 8px;
  background: #1a5276;
  border: none;
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.2s, box-shadow 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-submit-btn:hover:not(:disabled) {
  background: #164666;
  box-shadow:
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.12) 0px 4px 12px;
}

.login-submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Divider */
.login-divider {
  position: relative;
  text-align: center;
  margin: 28px 0;
}

.login-divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  border-top: 1px solid #e8e8e8;
}

.login-divider span {
  position: relative;
  padding: 0 16px;
  background: #ffffff;
  color: #6a6a6a;
  font-size: 0.85rem;
}

/* Google */
.login-google-btn {
  width: 100%;
  height: 48px;
  border: 1px solid #c1c1c1;
  border-radius: 8px;
  background: #ffffff;
  color: #222222;
  font-weight: 500;
  font-size: 0.95rem;
  cursor: pointer;
  transition: border-color 0.2s, box-shadow 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-google-btn:hover {
  border-color: #222222;
  box-shadow:
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.08) 0px 4px 8px;
}

/* Register Link */
.login-register-link {
  text-align: center;
  margin-top: 28px;
  color: #6a6a6a;
  font-size: 0.9rem;
}

.login-register-link a {
  color: #1a5276;
  text-decoration: none;
  font-weight: 600;
  margin-inline-start: 4px;
  transition: color 0.15s;
}

.login-register-link a:hover {
  color: #c0982b;
}

/* Responsive */
@media (max-width: 576px) {
  .login-card {
    padding: 32px 24px;
    border-radius: 16px;
  }

  .login-greeting {
    font-size: 1.5rem;
  }
}
</style>
