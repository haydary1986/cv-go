<template>
  <div class="register-page">
    <div class="container">
      <div class="row justify-content-center align-items-center min-vh-100 py-4">
        <div class="col-sm-11 col-md-9 col-lg-7 col-xl-6">
          <!-- Branding -->
          <div class="text-center mb-4">
            <div class="brand-icon mb-3">
              <i class="fas fa-file-alt"></i>
            </div>
            <h2 class="fw-bold text-white mb-1">CV Builder</h2>
            <p class="text-white-50">{{ t('auth.registerSubtitle') }}</p>
          </div>

          <!-- Register Card -->
          <div class="card register-card border-0 shadow-lg">
            <div class="card-body p-4 p-md-5">
              <h4 class="text-center fw-semibold mb-1">{{ t('auth.createYourAccount') }}</h4>
              <p class="text-center text-muted mb-4">{{ t('auth.registerSubtitle') }}</p>

              <!-- Error Alert -->
              <div v-if="error" class="alert alert-danger d-flex align-items-center border-0 shadow-sm" role="alert">
                <i class="fas fa-exclamation-circle me-2 flex-shrink-0"></i>
                <div>{{ error }}</div>
              </div>

              <form @submit.prevent="handleRegister">
                <div class="row g-3">
                  <!-- Full Name Arabic -->
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.fullNameAr') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-user text-muted"></i>
                      </span>
                      <input
                        type="text"
                        class="form-control border-start-0 ps-0"
                        v-model="form.full_name_ar"
                        required
                      />
                    </div>
                  </div>

                  <!-- Full Name English -->
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.fullNameEn') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-user text-muted"></i>
                      </span>
                      <input
                        type="text"
                        class="form-control border-start-0 ps-0"
                        v-model="form.full_name_en"
                        required
                      />
                    </div>
                  </div>

                  <!-- Email -->
                  <div class="col-12">
                    <label class="form-label fw-medium">{{ t('auth.email') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-envelope text-muted"></i>
                      </span>
                      <input
                        type="email"
                        class="form-control border-start-0 ps-0"
                        v-model="form.email"
                        required
                        autocomplete="email"
                      />
                    </div>
                  </div>

                  <!-- Password -->
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.password') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-lock text-muted"></i>
                      </span>
                      <input
                        :type="showPassword ? 'text' : 'password'"
                        class="form-control border-start-0 border-end-0 ps-0"
                        v-model="form.password"
                        required
                        minlength="6"
                        autocomplete="new-password"
                      />
                      <button
                        type="button"
                        class="input-group-text bg-light border-start-0 toggle-password"
                        @click="showPassword = !showPassword"
                      >
                        <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'" style="width: 16px;"></i>
                      </button>
                    </div>
                    <!-- Password Strength -->
                    <div class="mt-2" v-if="form.password">
                      <div class="password-strength-bar">
                        <div class="strength-segment" :class="passwordStrength >= 1 ? strengthColorClass : ''"></div>
                        <div class="strength-segment" :class="passwordStrength >= 2 ? strengthColorClass : ''"></div>
                        <div class="strength-segment" :class="passwordStrength >= 3 ? strengthColorClass : ''"></div>
                      </div>
                      <div class="d-flex justify-content-between align-items-center mt-1">
                        <small :class="passwordStrengthTextClass" class="fw-medium">{{ passwordStrengthLabel }}</small>
                      </div>
                      <!-- Password Requirements Checklist -->
                      <div class="password-requirements mt-2">
                        <div class="req-item" :class="form.password.length >= 6 ? 'met' : 'unmet'">
                          <i class="fas" :class="form.password.length >= 6 ? 'fa-check-circle' : 'fa-circle'"></i>
                          <span>{{ t('auth.minChars') }}</span>
                        </div>
                        <div class="req-item" :class="/[A-Z]/.test(form.password) ? 'met' : 'unmet'">
                          <i class="fas" :class="/[A-Z]/.test(form.password) ? 'fa-check-circle' : 'fa-circle'"></i>
                          <span>{{ t('auth.uppercase') }}</span>
                        </div>
                        <div class="req-item" :class="/[0-9]/.test(form.password) ? 'met' : 'unmet'">
                          <i class="fas" :class="/[0-9]/.test(form.password) ? 'fa-check-circle' : 'fa-circle'"></i>
                          <span>{{ t('auth.number') }}</span>
                        </div>
                        <div class="req-item" :class="/[^A-Za-z0-9]/.test(form.password) ? 'met' : 'unmet'">
                          <i class="fas" :class="/[^A-Za-z0-9]/.test(form.password) ? 'fa-check-circle' : 'fa-circle'"></i>
                          <span>{{ t('auth.specialChar') }}</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Confirm Password -->
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.confirmPassword') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-lock text-muted"></i>
                      </span>
                      <input
                        :type="showConfirmPassword ? 'text' : 'password'"
                        class="form-control border-start-0 border-end-0 ps-0"
                        :class="confirmPasswordInputClass"
                        v-model="form.confirm_password"
                        required
                        autocomplete="new-password"
                      />
                      <button
                        type="button"
                        class="input-group-text bg-light border-start-0 toggle-password"
                        @click="showConfirmPassword = !showConfirmPassword"
                      >
                        <i class="fas" :class="showConfirmPassword ? 'fa-eye-slash' : 'fa-eye'" style="width: 16px;"></i>
                      </button>
                    </div>
                    <div class="mt-2" v-if="form.confirm_password">
                      <small v-if="passwordsMatch" class="text-success d-flex align-items-center gap-1">
                        <i class="fas fa-check-circle"></i> {{ locale === 'ar' ? 'كلمات المرور متطابقة' : 'Passwords match' }}
                      </small>
                      <small v-else class="text-danger d-flex align-items-center gap-1">
                        <i class="fas fa-times-circle"></i> {{ t('auth.passwordMismatch') }}
                      </small>
                    </div>
                  </div>

                  <!-- Phone -->
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.phone') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-phone text-muted"></i>
                      </span>
                      <input
                        type="tel"
                        class="form-control border-start-0 ps-0"
                        v-model="form.phone"
                        placeholder="+964XXXXXXXXXX"
                        dir="ltr"
                      />
                    </div>
                    <small class="text-muted mt-1 d-block">
                      <i class="fas fa-info-circle me-1"></i>{{ t('auth.phoneHint') }}
                    </small>
                  </div>

                  <!-- Faculty -->
                  <div class="col-md-6">
                    <label class="form-label fw-medium">{{ t('auth.faculty') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-university text-muted"></i>
                      </span>
                      <select
                        class="form-select border-start-0 ps-0"
                        v-model="form.faculty_id"
                        @change="loadDepartments"
                        :disabled="loadingFaculties"
                      >
                        <option value="">{{ t('auth.selectFaculty') }}</option>
                        <option v-for="f in faculties" :key="f.id" :value="f.id">
                          {{ locale === 'ar' ? f.name_ar : f.name_en }}
                        </option>
                      </select>
                      <span v-if="loadingFaculties" class="input-group-text bg-light border-start-0">
                        <span class="spinner-border spinner-border-sm text-primary"></span>
                      </span>
                    </div>
                  </div>

                  <!-- Department -->
                  <div class="col-md-6" v-if="departments.length || loadingDepartments">
                    <label class="form-label fw-medium">{{ t('auth.department') }}</label>
                    <div class="input-group">
                      <span class="input-group-text bg-light border-end-0">
                        <i class="fas fa-building text-muted"></i>
                      </span>
                      <select
                        class="form-select border-start-0 ps-0"
                        v-model="form.department_id"
                        :disabled="loadingDepartments"
                      >
                        <option value="">{{ t('auth.selectDepartment') }}</option>
                        <option v-for="d in departments" :key="d.id" :value="d.id">
                          {{ locale === 'ar' ? d.name_ar : d.name_en }}
                        </option>
                      </select>
                      <span v-if="loadingDepartments" class="input-group-text bg-light border-start-0">
                        <span class="spinner-border spinner-border-sm text-primary"></span>
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Submit Button -->
                <button type="submit" class="btn btn-primary btn-lg w-100 fw-semibold mt-4" :disabled="loading || !canSubmit">
                  <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
                  <i v-else class="fas fa-user-plus me-2"></i>
                  {{ loading ? t('app.loading') : t('auth.registerBtn') }}
                </button>
              </form>
            </div>

            <!-- Login Link -->
            <div class="card-footer bg-light text-center py-3 border-0">
              <span class="text-muted">{{ t('auth.hasAccount') }}</span>
              <router-link to="/login" class="fw-semibold ms-1 text-decoration-none">
                {{ t('auth.loginBtn') }}
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { publicAPI } from '../services/api'

const { t, locale } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const loadingFaculties = ref(false)
const loadingDepartments = ref(false)
const error = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const faculties = ref<any[]>([])
const departments = ref<any[]>([])

const form = reactive({
  full_name_ar: '',
  full_name_en: '',
  email: '',
  password: '',
  confirm_password: '',
  phone: '',
  faculty_id: '',
  department_id: '',
})

// Password strength computation
const passwordStrength = computed(() => {
  const p = form.password
  if (!p) return 0
  let score = 0
  if (p.length >= 6) score++
  if (p.length >= 10) score++
  if (/[A-Z]/.test(p)) score++
  if (/[0-9]/.test(p)) score++
  if (/[^A-Za-z0-9]/.test(p)) score++
  if (score <= 1) return 1 // weak
  if (score <= 3) return 2 // medium
  return 3 // strong
})

const strengthColorClass = computed(() => {
  return ['', 'strength-weak', 'strength-medium', 'strength-strong'][passwordStrength.value]
})

const passwordStrengthTextClass = computed(() => {
  return ['', 'text-danger', 'text-warning', 'text-success'][passwordStrength.value]
})

const passwordStrengthLabel = computed(() => {
  const labels = ['', t('auth.weak'), t('auth.medium'), t('auth.strong')]
  return labels[passwordStrength.value]
})

// Password match
const passwordsMatch = computed(() => {
  return form.password === form.confirm_password && form.confirm_password.length > 0
})

const confirmPasswordInputClass = computed(() => {
  if (!form.confirm_password) return ''
  return passwordsMatch.value ? 'is-valid' : 'is-invalid'
})

// Can submit check
const canSubmit = computed(() => {
  return form.full_name_ar && form.full_name_en && form.email &&
    form.password.length >= 6 && passwordsMatch.value
})

onMounted(async () => {
  loadingFaculties.value = true
  try {
    const res = await publicAPI.getFaculties()
    faculties.value = res.data.faculties || []
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  } finally {
    loadingFaculties.value = false
  }
})

async function loadDepartments() {
  if (!form.faculty_id) {
    departments.value = []
    return
  }
  loadingDepartments.value = true
  try {
    const res = await publicAPI.getDepartments(Number(form.faculty_id))
    departments.value = res.data.departments || []
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  } finally {
    loadingDepartments.value = false
  }
}

async function handleRegister() {
  if (form.password !== form.confirm_password) {
    error.value = t('auth.passwordMismatch')
    return
  }
  loading.value = true
  error.value = ''
  try {
    await authStore.register({
      ...form,
      faculty_id: form.faculty_id ? Number(form.faculty_id) : undefined,
      department_id: form.department_id ? Number(form.department_id) : undefined,
    })
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
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

.register-card {
  border-radius: 16px;
  overflow: hidden;
}

.input-group-text {
  border-color: #dee2e6;
}

.form-control,
.form-select {
  border-color: #dee2e6;
}

.form-control:focus,
.form-select:focus {
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

/* Password strength bar */
.password-strength-bar {
  display: flex;
  gap: 4px;
}

.strength-segment {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: #e9ecef;
  transition: background-color 0.3s ease;
}

.strength-segment.strength-weak {
  background-color: #dc3545;
}

.strength-segment.strength-medium {
  background-color: #ffc107;
}

.strength-segment.strength-strong {
  background-color: #198754;
}

/* Password requirements */
.password-requirements {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.req-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  transition: color 0.2s;
}

.req-item.met {
  color: #198754;
}

.req-item.unmet {
  color: #adb5bd;
}

.req-item i {
  font-size: 0.65rem;
}

.card-footer a {
  color: #2d5196;
}

.card-footer a:hover {
  color: #1a1c4e;
}

@media (max-width: 576px) {
  .register-card .card-body {
    padding: 1.5rem !important;
  }
}
</style>
