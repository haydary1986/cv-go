<template>
  <div class="profile-page">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-lg-9 col-xl-8">
          <!-- Profile Header -->
          <div class="profile-header card border-0 shadow-sm mb-4">
            <div class="card-body p-4">
              <div class="d-flex flex-column flex-sm-row align-items-center gap-3">
                <div class="avatar-circle">
                  <span class="avatar-initials">{{ userInitials }}</span>
                </div>
                <div class="text-center text-sm-start flex-grow-1">
                  <h4 class="mb-1 fw-bold">{{ displayName }}</h4>
                  <p class="text-muted mb-2">{{ authStore.user?.email }}</p>
                  <div class="d-flex flex-wrap gap-2 justify-content-center justify-content-sm-start">
                    <span class="badge role-badge">
                      <i class="fas fa-user-tag me-1"></i>{{ authStore.user?.role }}
                    </span>
                    <span class="badge bg-info-subtle text-info-emphasis">
                      <i class="fas fa-coins me-1"></i>{{ t('auth.aiCredits') }}: {{ authStore.user?.ai_credits || 0 }}
                    </span>
                    <span v-if="authStore.user?.created_at" class="badge bg-light text-muted">
                      <i class="fas fa-calendar-alt me-1"></i>{{ t('auth.memberSince') }}: {{ formatDate(authStore.user.created_at) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Alerts -->
          <transition name="fade">
            <div v-if="success" class="alert alert-success d-flex align-items-center border-0 shadow-sm mb-4" role="alert">
              <i class="fas fa-check-circle me-2 flex-shrink-0"></i>
              <div>{{ success }}</div>
              <button type="button" class="btn-close ms-auto" @click="success = ''"></button>
            </div>
          </transition>
          <transition name="fade">
            <div v-if="error" class="alert alert-danger d-flex align-items-center border-0 shadow-sm mb-4" role="alert">
              <i class="fas fa-exclamation-circle me-2 flex-shrink-0"></i>
              <div>{{ error }}</div>
              <button type="button" class="btn-close ms-auto" @click="error = ''"></button>
            </div>
          </transition>

          <!-- Personal Information Card -->
          <div class="card border-0 shadow-sm mb-4">
            <div class="card-header bg-white border-bottom py-3">
              <h5 class="mb-0 fw-semibold">
                <i class="fas fa-user-edit text-primary me-2"></i>{{ t('auth.personalInformation') }}
              </h5>
            </div>
            <div class="card-body p-4">
              <form @submit.prevent="updateProfile">
                <div class="row g-3">
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.fullNameAr') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-user text-muted"></i>
                      </span>
                      <input type="text" class="form-control border-start-0 ps-0" v-model="form.full_name_ar" />
                    </div>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.fullNameEn') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-user text-muted"></i>
                      </span>
                      <input type="text" class="form-control border-start-0 ps-0" v-model="form.full_name_en" />
                    </div>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.email') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-envelope text-muted"></i>
                      </span>
                      <input type="email" class="form-control border-start-0 ps-0 bg-light" :value="authStore.user?.email" disabled />
                    </div>
                    <small class="text-muted"><i class="fas fa-info-circle me-1"></i>{{ locale === 'ar' ? 'لا يمكن تغيير البريد الإلكتروني' : 'Email cannot be changed' }}</small>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.phone') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-phone text-muted"></i>
                      </span>
                      <input type="tel" class="form-control border-start-0 ps-0" v-model="form.phone" dir="ltr" />
                    </div>
                  </div>
                </div>
                <div class="mt-4">
                  <button type="submit" class="btn btn-primary px-4" :disabled="savingProfile">
                    <span v-if="savingProfile" class="spinner-border spinner-border-sm me-2" role="status"></span>
                    <i v-else class="fas fa-save me-2"></i>
                    {{ savingProfile ? t('auth.saving') : t('app.save') }}
                  </button>
                </div>
              </form>
            </div>
          </div>

          <!-- Security Settings Card -->
          <div class="card border-0 shadow-sm mb-4">
            <div class="card-header bg-white border-bottom py-3">
              <h5 class="mb-0 fw-semibold">
                <i class="fas fa-shield-alt text-warning me-2"></i>{{ t('auth.securitySettings') }}
              </h5>
            </div>
            <div class="card-body p-4">
              <form @submit.prevent="changePassword">
                <div class="row g-3">
                  <div class="col-md-4">
                    <label class="form-label fw-medium">{{ t('auth.currentPassword') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-lock text-muted"></i>
                      </span>
                      <input
                        :type="showCurrentPw ? 'text' : 'password'"
                        class="form-control border-start-0 border-end-0 ps-0"
                        v-model="pwForm.current_password"
                        required
                      />
                      <button type="button" class="input-group-text bg-light border-start-0 toggle-password" @click="showCurrentPw = !showCurrentPw">
                        <i class="fas" :class="showCurrentPw ? 'fa-eye-slash' : 'fa-eye'" style="width: 16px;"></i>
                      </button>
                    </div>
                  </div>
                  <div class="col-md-4">
                    <label class="form-label fw-medium">{{ t('auth.newPassword') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-key text-muted"></i>
                      </span>
                      <input
                        :type="showNewPw ? 'text' : 'password'"
                        class="form-control border-start-0 border-end-0 ps-0"
                        v-model="pwForm.new_password"
                        required
                        minlength="6"
                      />
                      <button type="button" class="input-group-text bg-light border-start-0 toggle-password" @click="showNewPw = !showNewPw">
                        <i class="fas" :class="showNewPw ? 'fa-eye-slash' : 'fa-eye'" style="width: 16px;"></i>
                      </button>
                    </div>
                  </div>
                  <div class="col-md-4">
                    <label class="form-label fw-medium">{{ t('auth.confirmPassword') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-key text-muted"></i>
                      </span>
                      <input
                        :type="showConfirmPw ? 'text' : 'password'"
                        class="form-control border-start-0 border-end-0 ps-0"
                        :class="pwConfirmClass"
                        v-model="pwForm.confirm_password"
                        required
                      />
                      <button type="button" class="input-group-text bg-light border-start-0 toggle-password" @click="showConfirmPw = !showConfirmPw">
                        <i class="fas" :class="showConfirmPw ? 'fa-eye-slash' : 'fa-eye'" style="width: 16px;"></i>
                      </button>
                    </div>
                    <small v-if="pwForm.confirm_password && pwForm.new_password !== pwForm.confirm_password" class="text-danger">
                      <i class="fas fa-times-circle me-1"></i>{{ t('auth.passwordMismatch') }}
                    </small>
                  </div>
                </div>
                <div class="mt-4">
                  <button type="submit" class="btn btn-warning px-4" :disabled="savingPassword">
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
.profile-page {
  min-height: 100vh;
  background: #f4f6f9;
}

.profile-header {
  border-radius: 16px;
  overflow: hidden;
}

.profile-header .card-body {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.avatar-circle {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2d5196, #3b82c4);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 14px rgba(45, 81, 150, 0.3);
}

.avatar-initials {
  color: #fff;
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: 1px;
}

.role-badge {
  background: linear-gradient(135deg, #2d5196, #3b82c4);
  color: #fff;
  text-transform: capitalize;
}

.card {
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  border-radius: 12px 12px 0 0 !important;
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

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 576px) {
  .avatar-circle {
    width: 64px;
    height: 64px;
  }

  .avatar-initials {
    font-size: 1.25rem;
  }
}
</style>
