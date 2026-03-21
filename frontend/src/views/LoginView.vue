<template>
  <div class="login-page">
    <div class="container">
      <div class="row justify-content-center align-items-center min-vh-100 py-4">
        <div class="col-sm-10 col-md-7 col-lg-5 col-xl-4">
          <!-- Branding -->
          <div class="text-center mb-4">
            <div class="brand-icon mb-3">
              <img v-if="brandingStore.logoUrl" :src="brandingStore.logoUrl" alt="Logo" class="brand-logo" />
              <i v-else class="fas fa-file-alt"></i>
            </div>
            <h2 class="fw-bold text-white mb-1">{{ brandingStore.systemName }}</h2>
            <p class="text-white-50">{{ t('auth.loginSubtitle') }}</p>
          </div>

          <!-- Login Card -->
          <div class="card login-card border-0 shadow-lg">
            <div class="card-body p-4 p-md-5">
              <h4 class="text-center fw-semibold mb-1">{{ t('auth.welcomeBack') }}</h4>
              <p class="text-center text-muted mb-4">{{ t('auth.signInToAccount') }}</p>

              <!-- Error Alert -->
              <div v-if="error" class="alert alert-danger d-flex align-items-center border-0 shadow-sm" role="alert">
                <i class="fas fa-exclamation-circle me-2 flex-shrink-0"></i>
                <div>{{ error }}</div>
              </div>

              <form @submit.prevent="handleLogin">
                <!-- Email -->
                <div class="mb-3">
                  <label class="form-label fw-medium">{{ t('auth.email') }}</label>
                  <div class="input-group input-group-lg">
                    <span class="input-group-text bg-light border-end-0">
                      <i class="fas fa-envelope text-muted"></i>
                    </span>
                    <input
                      type="email"
                      class="form-control border-start-0 ps-0"
                      v-model="form.email"
                      :placeholder="t('auth.email')"
                      required
                      autocomplete="email"
                    />
                  </div>
                </div>

                <!-- Password -->
                <div class="mb-4">
                  <label class="form-label fw-medium">{{ t('auth.password') }}</label>
                  <div class="input-group input-group-lg">
                    <span class="input-group-text bg-light border-end-0">
                      <i class="fas fa-lock text-muted"></i>
                    </span>
                    <input
                      :type="showPassword ? 'text' : 'password'"
                      class="form-control border-start-0 border-end-0 ps-0"
                      v-model="form.password"
                      :placeholder="t('auth.password')"
                      required
                      autocomplete="current-password"
                    />
                    <button
                      type="button"
                      class="input-group-text bg-light border-start-0 toggle-password"
                      @click="showPassword = !showPassword"
                      :title="showPassword ? t('auth.hidePassword') : t('auth.showPassword')"
                    >
                      <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'" style="width: 16px;"></i>
                    </button>
                  </div>
                </div>

                <!-- Submit Button -->
                <button type="submit" class="btn btn-primary btn-lg w-100 fw-semibold mb-3" :disabled="loading">
                  <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
                  <i v-else class="fas fa-sign-in-alt me-2"></i>
                  {{ loading ? t('app.loading') : t('auth.loginBtn') }}
                </button>
              </form>

              <!-- Divider -->
              <div class="divider-text my-4">
                <span>{{ t('auth.orContinueWith') }}</span>
              </div>

              <!-- Google Login -->
              <button @click="handleGoogleLogin" class="btn btn-outline-secondary btn-lg w-100 google-btn">
                <i class="fab fa-google me-2"></i>{{ t('auth.googleLogin') }}
              </button>
            </div>

            <!-- Register Link -->
            <div class="card-footer bg-light text-center py-3 border-0">
              <span class="text-muted">{{ t('auth.noAccount') }}</span>
              <router-link to="/register" class="fw-semibold ms-1 text-decoration-none">
                {{ t('auth.registerBtn') }}
              </router-link>
            </div>
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
    const redirect = (route.query.redirect as string) || '/dashboard'
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
  background: linear-gradient(135deg, #1a1c4e 0%, #2d5196 50%, #3b82c4 100%);
  display: flex;
  align-items: center;
}

.brand-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 72px;
  height: 72px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  font-size: 32px;
  color: #fff;
}

.brand-logo {
  height: 48px;
  width: auto;
  object-fit: contain;
}

.login-card {
  border-radius: 16px;
  overflow: hidden;
  transition: transform 0.2s ease;
}

.input-group-text {
  border-color: #dee2e6;
}

.form-control {
  border-color: #dee2e6;
}

.form-control:focus {
  box-shadow: none;
  border-color: #86b7fe;
}

.form-control:focus + .toggle-password,
.form-control:focus ~ .toggle-password {
  border-color: #86b7fe;
}

.input-group:focus-within .input-group-text {
  border-color: #86b7fe;
}

.toggle-password {
  cursor: pointer;
  transition: color 0.15s;
}

.toggle-password:hover {
  color: #0d6efd;
}

.divider-text {
  position: relative;
  text-align: center;
}

.divider-text::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  border-top: 1px solid #dee2e6;
}

.divider-text span {
  position: relative;
  padding: 0 16px;
  background: #fff;
  color: #6c757d;
  font-size: 0.875rem;
}

.google-btn {
  border-radius: 8px;
  transition: all 0.2s;
}

.google-btn:hover {
  background-color: #f8f9fa;
  border-color: #adb5bd;
}

.card-footer a {
  color: #2d5196;
}

.card-footer a:hover {
  color: #1a1c4e;
}

@media (max-width: 576px) {
  .login-card .card-body {
    padding: 1.5rem !important;
  }
}
</style>
