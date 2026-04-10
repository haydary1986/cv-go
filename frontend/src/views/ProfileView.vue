<template>
  <div class="prof-page">
    <div class="prof-container">
      <!-- Profile Header -->
      <div class="prof-header">
        <div class="prof-avatar">
          <span class="prof-avatar-text">{{ userInitials }}</span>
        </div>
        <div class="prof-header-info">
          <h2 class="prof-name">{{ displayName }}</h2>
          <p class="prof-email">{{ authStore.user?.email }}</p>
          <div class="prof-badges">
            <span class="prof-role-badge">{{ authStore.user?.role }}</span>
            <span v-if="authStore.user?.created_at" class="prof-date-badge">
              <i class="fas fa-calendar-alt me-1"></i>{{ t('auth.memberSince') }}: {{ formatDate(authStore.user.created_at) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Stats Row -->
      <div class="prof-stats-row">
        <div class="prof-stat-card">
          <div class="prof-stat-icon prof-stat-icon-gold">
            <i class="fas fa-coins"></i>
          </div>
          <div class="prof-stat-info">
            <span class="prof-stat-value">{{ authStore.user?.ai_credits || 0 }}</span>
            <span class="prof-stat-label">{{ t('auth.aiCredits') }}</span>
          </div>
        </div>
        <div class="prof-stat-card">
          <div class="prof-stat-icon prof-stat-icon-blue">
            <i class="fas fa-file-alt"></i>
          </div>
          <div class="prof-stat-info">
            <span class="prof-stat-value">{{ authStore.user?.cvs_count || 0 }}</span>
            <span class="prof-stat-label">{{ locale === 'ar' ? 'السير الذاتية' : 'CVs Created' }}</span>
          </div>
        </div>
      </div>

      <!-- Alerts -->
      <transition name="fade">
        <div v-if="success" class="prof-alert prof-alert-success" role="alert">
          <i class="fas fa-check-circle"></i>
          <span>{{ success }}</span>
          <button type="button" class="prof-alert-close" @click="success = ''">
            <i class="fas fa-times"></i>
          </button>
        </div>
      </transition>
      <transition name="fade">
        <div v-if="error" class="prof-alert prof-alert-error" role="alert">
          <i class="fas fa-exclamation-circle"></i>
          <span>{{ error }}</span>
          <button type="button" class="prof-alert-close" @click="error = ''">
            <i class="fas fa-times"></i>
          </button>
        </div>
      </transition>

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
            <div class="prof-form-grid">
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.fullNameAr') }}</label>
                <input type="text" class="prof-input" v-model="form.full_name_ar" />
              </div>
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.fullNameEn') }}</label>
                <input type="text" class="prof-input" v-model="form.full_name_en" />
              </div>
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.email') }}</label>
                <input type="email" class="prof-input prof-input-disabled" :value="authStore.user?.email" disabled />
                <small class="prof-hint">
                  <i class="fas fa-info-circle me-1"></i>{{ locale === 'ar' ? 'لا يمكن تغيير البريد الإلكتروني' : 'Email cannot be changed' }}
                </small>
              </div>
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.phone') }}</label>
                <input type="tel" class="prof-input" v-model="form.phone" dir="ltr" />
              </div>
            </div>
            <div class="prof-card-actions">
              <button type="submit" class="prof-save-btn" :disabled="savingProfile">
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
            <div class="prof-form-grid prof-form-grid-3">
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.currentPassword') }}</label>
                <div class="prof-input-wrap">
                  <input
                    :type="showCurrentPw ? 'text' : 'password'"
                    class="prof-input prof-input-pw"
                    v-model="pwForm.current_password"
                    required
                  />
                  <button type="button" class="prof-toggle-pw" @click="showCurrentPw = !showCurrentPw">
                    <i class="fas" :class="showCurrentPw ? 'fa-eye-slash' : 'fa-eye'"></i>
                  </button>
                </div>
              </div>
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.newPassword') }}</label>
                <div class="prof-input-wrap">
                  <input
                    :type="showNewPw ? 'text' : 'password'"
                    class="prof-input prof-input-pw"
                    v-model="pwForm.new_password"
                    required
                    minlength="6"
                  />
                  <button type="button" class="prof-toggle-pw" @click="showNewPw = !showNewPw">
                    <i class="fas" :class="showNewPw ? 'fa-eye-slash' : 'fa-eye'"></i>
                  </button>
                </div>
              </div>
              <div class="prof-form-field">
                <label class="prof-label">{{ t('auth.confirmPassword') }}</label>
                <div class="prof-input-wrap">
                  <input
                    :type="showConfirmPw ? 'text' : 'password'"
                    class="prof-input prof-input-pw"
                    :class="pwConfirmClass"
                    v-model="pwForm.confirm_password"
                    required
                  />
                  <button type="button" class="prof-toggle-pw" @click="showConfirmPw = !showConfirmPw">
                    <i class="fas" :class="showConfirmPw ? 'fa-eye-slash' : 'fa-eye'"></i>
                  </button>
                </div>
                <small v-if="pwForm.confirm_password && pwForm.new_password !== pwForm.confirm_password" class="prof-pw-mismatch">
                  <i class="fas fa-times-circle me-1"></i>{{ t('auth.passwordMismatch') }}
                </small>
              </div>
            </div>
            <div class="prof-card-actions">
              <button type="submit" class="prof-save-btn" :disabled="savingPassword">
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
.prof-page {
  min-height: 100vh;
  background: #ffffff;
  padding: 40px 24px;
}

.prof-container {
  max-width: 720px;
  margin: 0 auto;
}

/* Profile Header */
.prof-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 32px;
}

.prof-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #1a5276;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.prof-avatar-text {
  color: #ffffff;
  font-size: 1.25rem;
  font-weight: 700;
  letter-spacing: 1px;
}

.prof-header-info {
  flex: 1;
}

.prof-name {
  font-size: 1.5rem;
  font-weight: 700;
  color: #222222;
  margin-bottom: 2px;
}

.prof-email {
  color: #6a6a6a;
  font-size: 0.9rem;
  margin-bottom: 8px;
}

.prof-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.prof-role-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 12px;
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: capitalize;
}

.prof-date-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 12px;
  background: #f5f5f5;
  color: #6a6a6a;
  font-size: 0.8rem;
}

/* Stats Row */
.prof-stats-row {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
}

.prof-stat-card {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  background: #ffffff;
  padding: 16px 20px;
  border-radius: 12px;
  border: 1px solid #e8e8e8;
}

.prof-stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  flex-shrink: 0;
}

.prof-stat-icon-gold {
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
}

.prof-stat-icon-blue {
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
}

.prof-stat-info {
  display: flex;
  flex-direction: column;
}

.prof-stat-value {
  font-size: 1.25rem;
  font-weight: 700;
  color: #222222;
  line-height: 1.2;
}

.prof-stat-label {
  font-size: 0.75rem;
  color: #6a6a6a;
  font-weight: 500;
}

/* Alerts */
.prof-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 0.9rem;
  margin-bottom: 24px;
}

.prof-alert-success {
  background: #f0fdf4;
  color: #166534;
}

.prof-alert-error {
  background: #fef2f2;
  color: #b91c1c;
}

.prof-alert-close {
  margin-inline-start: auto;
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 4px;
  opacity: 0.6;
}

.prof-alert-close:hover {
  opacity: 1;
}

/* Cards */
.prof-card {
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 24px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.prof-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.prof-card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.prof-card-icon-blue {
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
}

.prof-card-icon-gold {
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
}

.prof-card-title {
  font-weight: 600;
  font-size: 1rem;
  color: #222222;
  margin: 0;
}

.prof-card-body {
  padding: 24px;
}

/* Form Grid */
.prof-form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.prof-form-grid-3 {
  grid-template-columns: 1fr 1fr 1fr;
}

.prof-form-field {
  display: flex;
  flex-direction: column;
}

/* Inputs */
.prof-label {
  font-weight: 500;
  font-size: 0.85rem;
  color: #222222;
  margin-bottom: 6px;
}

.prof-input-wrap {
  position: relative;
}

.prof-input {
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

.prof-input-pw {
  padding-inline-end: 44px;
}

.prof-input-disabled {
  background: #f5f5f5;
  color: #6a6a6a;
  cursor: not-allowed;
}

.prof-input:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 2px rgba(26, 82, 118, 0.15);
}

.prof-input.is-valid {
  border-color: #198754;
}

.prof-input.is-invalid {
  border-color: #dc3545;
}

.prof-hint {
  font-size: 0.78rem;
  color: #6a6a6a;
  margin-top: 4px;
}

.prof-pw-mismatch {
  color: #dc3545;
  font-size: 0.8rem;
  margin-top: 4px;
}

.prof-toggle-pw {
  position: absolute;
  top: 50%;
  inset-inline-end: 12px;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #b0b0b0;
  cursor: pointer;
  padding: 4px;
  z-index: 2;
  transition: color 0.15s;
}

.prof-toggle-pw:hover {
  color: #1a5276;
}

/* Actions */
.prof-card-actions {
  margin-top: 20px;
}

.prof-save-btn {
  height: 44px;
  padding: 0 28px;
  font-size: 0.95rem;
  font-weight: 600;
  border-radius: 8px;
  background: #1a5276;
  border: none;
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.2s, box-shadow 0.2s;
  display: inline-flex;
  align-items: center;
}

.prof-save-btn:hover:not(:disabled) {
  background: #164666;
  box-shadow:
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.12) 0px 4px 12px;
}

.prof-save-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .prof-header {
    flex-direction: column;
    text-align: center;
  }

  .prof-badges {
    justify-content: center;
  }

  .prof-stats-row {
    flex-direction: column;
    gap: 12px;
  }

  .prof-form-grid {
    grid-template-columns: 1fr;
  }

  .prof-form-grid-3 {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 576px) {
  .prof-card-body {
    padding: 16px;
  }
}
</style>
