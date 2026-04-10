<template>
  <div class="prof-page">
    <!-- Hero Header -->
    <div class="prof-hero">
      <div class="prof-hero-bg"></div>
      <div class="container position-relative">
        <div class="prof-hero-content">
          <div class="prof-avatar">
            <span class="prof-avatar-text">{{ userInitials }}</span>
          </div>
          <div class="prof-hero-info">
            <h2 class="prof-hero-name">{{ displayName }}</h2>
            <p class="prof-hero-email">{{ authStore.user?.email }}</p>
            <div class="prof-hero-badges">
              <span class="prof-role-badge">
                <i class="fas fa-user-tag me-1"></i>{{ authStore.user?.role }}
              </span>
              <span v-if="authStore.user?.created_at" class="prof-date-badge">
                <i class="fas fa-calendar-alt me-1"></i>{{ t('auth.memberSince') }}: {{ formatDate(authStore.user.created_at) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Bar -->
    <div class="prof-stats-bar">
      <div class="container">
        <div class="prof-stats-grid">
          <div class="prof-stat">
            <div class="prof-stat-icon">
              <i class="fas fa-coins"></i>
            </div>
            <div class="prof-stat-info">
              <span class="prof-stat-value">{{ authStore.user?.ai_credits || 0 }}</span>
              <span class="prof-stat-label">{{ t('auth.aiCredits') }}</span>
            </div>
          </div>
          <div class="prof-stat">
            <div class="prof-stat-icon">
              <i class="fas fa-file-alt"></i>
            </div>
            <div class="prof-stat-info">
              <span class="prof-stat-value">{{ authStore.user?.cvs_count || 0 }}</span>
              <span class="prof-stat-label">{{ locale === 'ar' ? 'السير الذاتية' : 'CVs Created' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="container prof-body">
      <!-- Alerts -->
      <transition name="fade">
        <div v-if="success" class="alert alert-success d-flex align-items-center border-0 prof-alert" role="alert">
          <i class="fas fa-check-circle me-2 flex-shrink-0"></i>
          <div>{{ success }}</div>
          <button type="button" class="btn-close ms-auto" @click="success = ''"></button>
        </div>
      </transition>
      <transition name="fade">
        <div v-if="error" class="alert alert-danger d-flex align-items-center border-0 prof-alert" role="alert">
          <i class="fas fa-exclamation-circle me-2 flex-shrink-0"></i>
          <div>{{ error }}</div>
          <button type="button" class="btn-close ms-auto" @click="error = ''"></button>
        </div>
      </transition>

      <div class="row g-4 justify-content-center">
        <div class="col-lg-9 col-xl-8">
          <!-- Personal Information Card -->
          <div class="prof-card">
            <div class="prof-card-header">
              <div class="prof-card-icon prof-card-icon-blue">
                <i class="fas fa-user-edit"></i>
              </div>
              <h5 class="prof-card-title">{{ t('auth.personalInformation') }}</h5>
            </div>
            <div class="prof-card-body">
              <form @submit.prevent="updateProfile">
                <div class="row g-3">
                  <div class="col-md-6">
                    <label class="prof-label">{{ t('auth.fullNameAr') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-user prof-input-icon"></i>
                      <input type="text" class="form-control prof-input" v-model="form.full_name_ar" />
                    </div>
                  </div>
                  <div class="col-md-6">
                    <label class="prof-label">{{ t('auth.fullNameEn') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-user prof-input-icon"></i>
                      <input type="text" class="form-control prof-input" v-model="form.full_name_en" />
                    </div>
                  </div>
                  <div class="col-md-6">
                    <label class="prof-label">{{ t('auth.email') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-envelope prof-input-icon"></i>
                      <input type="email" class="form-control prof-input prof-input-disabled" :value="authStore.user?.email" disabled />
                    </div>
                    <small class="prof-hint">
                      <i class="fas fa-info-circle me-1"></i>{{ locale === 'ar' ? 'لا يمكن تغيير البريد الإلكتروني' : 'Email cannot be changed' }}
                    </small>
                  </div>
                  <div class="col-md-6">
                    <label class="prof-label">{{ t('auth.phone') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-phone prof-input-icon"></i>
                      <input type="tel" class="form-control prof-input" v-model="form.phone" dir="ltr" />
                    </div>
                  </div>
                </div>
                <div class="mt-4">
                  <button type="submit" class="btn prof-save-btn" :disabled="savingProfile">
                    <span v-if="savingProfile" class="spinner-border spinner-border-sm me-2" role="status"></span>
                    <i v-else class="fas fa-save me-2"></i>
                    {{ savingProfile ? t('auth.saving') : t('app.save') }}
                  </button>
                </div>
              </form>
            </div>
          </div>

          <!-- Security Settings Card -->
          <div class="prof-card">
            <div class="prof-card-header">
              <div class="prof-card-icon prof-card-icon-gold">
                <i class="fas fa-shield-alt"></i>
              </div>
              <h5 class="prof-card-title">{{ t('auth.securitySettings') }}</h5>
            </div>
            <div class="prof-card-body">
              <form @submit.prevent="changePassword">
                <div class="row g-3">
                  <div class="col-md-4">
                    <label class="prof-label">{{ t('auth.currentPassword') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-lock prof-input-icon"></i>
                      <input
                        :type="showCurrentPw ? 'text' : 'password'"
                        class="form-control prof-input prof-input-pw"
                        v-model="pwForm.current_password"
                        required
                      />
                      <button type="button" class="prof-toggle-pw" @click="showCurrentPw = !showCurrentPw">
                        <i class="fas" :class="showCurrentPw ? 'fa-eye-slash' : 'fa-eye'"></i>
                      </button>
                    </div>
                  </div>
                  <div class="col-md-4">
                    <label class="prof-label">{{ t('auth.newPassword') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-key prof-input-icon"></i>
                      <input
                        :type="showNewPw ? 'text' : 'password'"
                        class="form-control prof-input prof-input-pw"
                        v-model="pwForm.new_password"
                        required
                        minlength="6"
                      />
                      <button type="button" class="prof-toggle-pw" @click="showNewPw = !showNewPw">
                        <i class="fas" :class="showNewPw ? 'fa-eye-slash' : 'fa-eye'"></i>
                      </button>
                    </div>
                  </div>
                  <div class="col-md-4">
                    <label class="prof-label">{{ t('auth.confirmPassword') }}</label>
                    <div class="prof-input-wrap">
                      <i class="fas fa-key prof-input-icon"></i>
                      <input
                        :type="showConfirmPw ? 'text' : 'password'"
                        class="form-control prof-input prof-input-pw"
                        :class="pwConfirmClass"
                        v-model="pwForm.confirm_password"
                        required
                      />
                      <button type="button" class="prof-toggle-pw" @click="showConfirmPw = !showConfirmPw">
                        <i class="fas" :class="showConfirmPw ? 'fa-eye-slash' : 'fa-eye'"></i>
                      </button>
                    </div>
                    <small v-if="pwForm.confirm_password && pwForm.new_password !== pwForm.confirm_password" class="text-danger" style="font-size: 0.8rem;">
                      <i class="fas fa-times-circle me-1"></i>{{ t('auth.passwordMismatch') }}
                    </small>
                  </div>
                </div>
                <div class="mt-4">
                  <button type="submit" class="btn prof-pw-btn" :disabled="savingPassword">
                    <span v-if="savingPassword" class="spinner-border spinner-border-sm me-2" role="status"></span>
                    <i v-else class="fas fa-sync-alt me-2"></i>
                    {{ savingPassword ? t('auth.saving') : t('auth.changePassword') }}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { authAPI } from '../services/api'

const { t, locale } = useI18n()
const authStore = useAuthStore()

const success = ref('')
const error = ref('')
const savingProfile = ref(false)
const savingPassword = ref(false)
const showCurrentPw = ref(false)
const showNewPw = ref(false)
const showConfirmPw = ref(false)
const form = reactive({ full_name_ar: '', full_name_en: '', phone: '' })
const pwForm = reactive({ current_password: '', new_password: '', confirm_password: '' })

const userInitials = computed(() => {
  const name = locale.value === 'ar'
    ? authStore.user?.full_name_ar
    : authStore.user?.full_name_en
  if (!name) return '?'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return name.substring(0, 2).toUpperCase()
})

const displayName = computed(() => {
  if (locale.value === 'ar') {
    return authStore.user?.full_name_ar || authStore.user?.full_name_en || ''
  }
  return authStore.user?.full_name_en || authStore.user?.full_name_ar || ''
})

const pwConfirmClass = computed(() => {
  if (!pwForm.confirm_password) return ''
  return pwForm.new_password === pwForm.confirm_password ? 'is-valid' : 'is-invalid'
})

function formatDate(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleDateString(locale.value === 'ar' ? 'ar-IQ' : 'en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    })
  } catch {
    return dateStr
  }
}

onMounted(() => {
  if (authStore.user) {
    form.full_name_ar = authStore.user.full_name_ar || ''
    form.full_name_en = authStore.user.full_name_en || ''
    form.phone = authStore.user.phone || ''
  }
})

async function updateProfile() {
  error.value = ''
  success.value = ''
  savingProfile.value = true
  try {
    await authStore.updateProfile(form)
    success.value = t('auth.profileUpdateSuccess')
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  } finally {
    savingProfile.value = false
  }
}

async function changePassword() {
  error.value = ''
  success.value = ''
  if (pwForm.new_password !== pwForm.confirm_password) {
    error.value = t('auth.passwordMismatch')
    return
  }
  savingPassword.value = true
  try {
    await authAPI.changePassword(pwForm)
    success.value = t('auth.passwordChangeSuccess')
    pwForm.current_password = ''
    pwForm.new_password = ''
    pwForm.confirm_password = ''
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  } finally {
    savingPassword.value = false
  }
}
</script>

<style scoped>
/* ── Page ── */
.prof-page {
  min-height: 100vh;
  background: #f0f2f5;
}

/* ── Hero ── */
.prof-hero {
  position: relative;
  padding: 3rem 0 4rem;
  overflow: hidden;
}

.prof-hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #1a5276 0%, #0d6efd 100%);
}

.prof-hero-bg::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 60px;
  background: #f0f2f5;
  clip-path: ellipse(55% 100% at 50% 100%);
}

.prof-hero-content {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  position: relative;
  z-index: 1;
}

.prof-avatar {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: linear-gradient(135deg, #c0982b, #d4a933);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  border: 4px solid rgba(255, 255, 255, 0.3);
}

.prof-avatar-text {
  color: #fff;
  font-size: 2rem;
  font-weight: 800;
  letter-spacing: 1px;
}

.prof-hero-info {
  flex: 1;
}

.prof-hero-name {
  color: #fff;
  font-size: 1.75rem;
  font-weight: 800;
  margin-bottom: 0.15rem;
}

.prof-hero-email {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  margin-bottom: 0.75rem;
}

.prof-hero-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.prof-role-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 20px;
  background: rgba(192, 152, 43, 0.25);
  color: #f0d78c;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: capitalize;
}

.prof-date-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.12);
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.8rem;
}

/* ── Stats Bar ── */
.prof-stats-bar {
  margin-top: -2rem;
  position: relative;
  z-index: 2;
  padding: 0 1rem;
}

.prof-stats-grid {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
  max-width: 500px;
  margin: 0 auto;
}

.prof-stat {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  background: #fff;
  padding: 1rem 1.25rem;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid #e9ecef;
}

.prof-stat-icon {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  flex-shrink: 0;
}

.prof-stat:first-child .prof-stat-icon {
  background: rgba(192, 152, 43, 0.12);
  color: #c0982b;
}

.prof-stat:last-child .prof-stat-icon {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.prof-stat-info {
  display: flex;
  flex-direction: column;
}

.prof-stat-value {
  font-size: 1.3rem;
  font-weight: 800;
  color: #2c3e50;
  line-height: 1.2;
}

.prof-stat-label {
  font-size: 0.75rem;
  color: #6c757d;
  font-weight: 500;
}

/* ── Body ── */
.prof-body {
  padding-top: 2rem;
  padding-bottom: 3rem;
}

.prof-alert {
  border-radius: 12px;
  margin-bottom: 1.5rem;
}

/* ── Cards ── */
.prof-card {
  background: #fff;
  border-radius: 16px;
  border: 1px solid #e9ecef;
  overflow: hidden;
  margin-bottom: 1.5rem;
  box-shadow: 0 1px 8px rgba(0, 0, 0, 0.04);
}

.prof-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.prof-card-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.95rem;
  flex-shrink: 0;
}

.prof-card-icon-blue {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.prof-card-icon-gold {
  background: rgba(192, 152, 43, 0.12);
  color: #c0982b;
}

.prof-card-title {
  font-weight: 700;
  font-size: 1.05rem;
  color: #2c3e50;
  margin: 0;
}

.prof-card-body {
  padding: 1.5rem;
}

/* ── Inputs ── */
.prof-label {
  font-weight: 600;
  font-size: 0.85rem;
  color: #2c3e50;
  margin-bottom: 0.35rem;
  display: block;
}

.prof-input-wrap {
  position: relative;
}

.prof-input-icon {
  position: absolute;
  top: 50%;
  inset-inline-start: 14px;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 0.9rem;
  pointer-events: none;
  z-index: 2;
}

.prof-input {
  height: 46px;
  padding-inline-start: 42px;
  padding-inline-end: 14px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 0.95rem;
  background: #f9fafb;
  transition: border-color 0.2s, background-color 0.2s, box-shadow 0.2s;
}

.prof-input-pw {
  padding-inline-end: 42px;
}

.prof-input-disabled {
  background: #f3f4f6;
  color: #6c757d;
  cursor: not-allowed;
}

.prof-input:focus {
  border-color: #1a5276;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
  outline: none;
}

.prof-hint {
  font-size: 0.78rem;
  color: #9ca3af;
  margin-top: 4px;
  display: block;
}

.prof-toggle-pw {
  position: absolute;
  top: 50%;
  inset-inline-end: 12px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 4px;
  z-index: 2;
  transition: color 0.15s;
}

.prof-toggle-pw:hover {
  color: #1a5276;
}

/* ── Buttons ── */
.prof-save-btn {
  background: linear-gradient(135deg, #1a5276, #0d6efd);
  color: #fff;
  border: none;
  border-radius: 10px;
  padding: 10px 28px;
  font-weight: 600;
  font-size: 0.95rem;
  transition: transform 0.15s, box-shadow 0.2s;
}

.prof-save-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(26, 82, 118, 0.3);
  color: #fff;
}

.prof-save-btn:disabled {
  opacity: 0.6;
  color: #fff;
}

.prof-pw-btn {
  background: linear-gradient(135deg, #c0982b, #d4a933);
  color: #fff;
  border: none;
  border-radius: 10px;
  padding: 10px 28px;
  font-weight: 600;
  font-size: 0.95rem;
  transition: transform 0.15s, box-shadow 0.2s;
}

.prof-pw-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(192, 152, 43, 0.3);
  color: #fff;
}

.prof-pw-btn:disabled {
  opacity: 0.6;
  color: #fff;
}

/* ── Transitions ── */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ── Responsive ── */
@media (max-width: 768px) {
  .prof-hero {
    padding: 2rem 0 3rem;
  }

  .prof-hero-content {
    flex-direction: column;
    text-align: center;
  }

  .prof-hero-badges {
    justify-content: center;
  }

  .prof-avatar {
    width: 80px;
    height: 80px;
  }

  .prof-avatar-text {
    font-size: 1.5rem;
  }

  .prof-hero-name {
    font-size: 1.35rem;
  }

  .prof-stats-grid {
    flex-direction: column;
    gap: 0.75rem;
  }
}

@media (max-width: 576px) {
  .prof-card-body {
    padding: 1rem;
  }
}
</style>
