<template>
  <div class="reg-page">
    <div class="reg-card">
      <!-- Logo -->
      <div class="reg-logo">
        <img :src="brandingStore.logoUrl || '/logo.svg'" alt="Logo" />
      </div>

      <h2 class="reg-heading">{{ t('auth.createYourAccount') }}</h2>
      <p class="reg-subtitle">{{ t('auth.registerSubtitle') }}</p>

      <!-- Error Alert -->
      <div v-if="error" class="reg-alert" role="alert">
        <i class="fas fa-exclamation-circle"></i>
        <span>{{ error }}</span>
      </div>

      <form @submit.prevent="handleRegister">
        <!-- Step 1: Account Info -->
        <div class="reg-section">
          <div class="reg-section-header">
            <span class="reg-step-badge">1</span>
            <span class="reg-step-title">{{ locale === 'ar' ? 'معلومات الحساب' : 'Account Information' }}</span>
          </div>

          <div class="reg-field">
            <label class="reg-label">{{ t('auth.email') }}</label>
            <input
              type="email"
              class="reg-input"
              v-model="form.email"
              required
              autocomplete="email"
            />
          </div>

          <div class="reg-row">
            <div class="reg-col">
              <label class="reg-label">{{ t('auth.password') }}</label>
              <div class="reg-input-wrap">
                <input
                  :type="showPassword ? 'text' : 'password'"
                  class="reg-input reg-input-pw"
                  v-model="form.password"
                  required
                  minlength="6"
                  autocomplete="new-password"
                />
                <button
                  type="button"
                  class="reg-toggle-pw"
                  @click="showPassword = !showPassword"
                >
                  <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                </button>
              </div>
              <!-- Password Strength -->
              <div class="reg-strength-area" v-if="form.password">
                <div class="reg-strength-bar">
                  <div class="reg-strength-seg" :class="passwordStrength >= 1 ? strengthColorClass : ''"></div>
                  <div class="reg-strength-seg" :class="passwordStrength >= 2 ? strengthColorClass : ''"></div>
                  <div class="reg-strength-seg" :class="passwordStrength >= 3 ? strengthColorClass : ''"></div>
                </div>
                <small :class="passwordStrengthTextClass" class="reg-strength-label">{{ passwordStrengthLabel }}</small>
                <div class="reg-pw-reqs">
                  <div class="reg-req" :class="form.password.length >= 6 ? 'met' : 'unmet'">
                    <i class="fas" :class="form.password.length >= 6 ? 'fa-check-circle' : 'fa-circle'"></i>
                    <span>{{ t('auth.minChars') }}</span>
                  </div>
                  <div class="reg-req" :class="/[A-Z]/.test(form.password) ? 'met' : 'unmet'">
                    <i class="fas" :class="/[A-Z]/.test(form.password) ? 'fa-check-circle' : 'fa-circle'"></i>
                    <span>{{ t('auth.uppercase') }}</span>
                  </div>
                  <div class="reg-req" :class="/[0-9]/.test(form.password) ? 'met' : 'unmet'">
                    <i class="fas" :class="/[0-9]/.test(form.password) ? 'fa-check-circle' : 'fa-circle'"></i>
                    <span>{{ t('auth.number') }}</span>
                  </div>
                  <div class="reg-req" :class="/[^A-Za-z0-9]/.test(form.password) ? 'met' : 'unmet'">
                    <i class="fas" :class="/[^A-Za-z0-9]/.test(form.password) ? 'fa-check-circle' : 'fa-circle'"></i>
                    <span>{{ t('auth.specialChar') }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="reg-col">
              <label class="reg-label">{{ t('auth.confirmPassword') }}</label>
              <div class="reg-input-wrap">
                <input
                  :type="showConfirmPassword ? 'text' : 'password'"
                  class="reg-input reg-input-pw"
                  :class="confirmPasswordInputClass"
                  v-model="form.confirm_password"
                  required
                  autocomplete="new-password"
                />
                <button
                  type="button"
                  class="reg-toggle-pw"
                  @click="showConfirmPassword = !showConfirmPassword"
                >
                  <i class="fas" :class="showConfirmPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                </button>
              </div>
              <div v-if="form.confirm_password" class="reg-match-status">
                <small v-if="passwordsMatch" class="reg-match-ok">
                  <i class="fas fa-check-circle"></i> {{ locale === 'ar' ? 'كلمات المرور متطابقة' : 'Passwords match' }}
                </small>
                <small v-else class="reg-match-err">
                  <i class="fas fa-times-circle"></i> {{ t('auth.passwordMismatch') }}
                </small>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 2: Personal Info -->
        <div class="reg-section">
          <div class="reg-section-header">
            <span class="reg-step-badge">2</span>
            <span class="reg-step-title">{{ locale === 'ar' ? 'المعلومات الشخصية' : 'Personal Information' }}</span>
          </div>

          <div class="reg-row">
            <div class="reg-col">
              <label class="reg-label">{{ t('auth.fullNameAr') }}</label>
              <input
                type="text"
                class="reg-input"
                v-model="form.full_name_ar"
                required
              />
            </div>

            <div class="reg-col">
              <label class="reg-label">{{ t('auth.fullNameEn') }}</label>
              <input
                type="text"
                class="reg-input"
                v-model="form.full_name_en"
                required
              />
            </div>
          </div>

          <div class="reg-field" style="max-width: 50%;">
            <label class="reg-label">{{ t('auth.phone') }}</label>
            <input
              type="tel"
              class="reg-input"
              v-model="form.phone"
              placeholder="+964XXXXXXXXXX"
              dir="ltr"
            />
            <small class="reg-hint">
              <i class="fas fa-info-circle me-1"></i>{{ t('auth.phoneHint') }}
            </small>
          </div>
        </div>

        <!-- Step 3: University -->
        <div class="reg-section reg-section-last">
          <div class="reg-section-header">
            <span class="reg-step-badge">3</span>
            <span class="reg-step-title">{{ locale === 'ar' ? 'المعلومات الجامعية' : 'University Details' }}</span>
          </div>

          <div class="reg-row">
            <div class="reg-col">
              <label class="reg-label">{{ t('auth.faculty') }}</label>
              <div class="reg-input-wrap">
                <select
                  class="reg-input reg-select"
                  v-model="form.faculty_id"
                  @change="loadDepartments"
                  :disabled="loadingFaculties"
                >
                  <option value="">{{ t('auth.selectFaculty') }}</option>
                  <option v-for="f in faculties" :key="f.id" :value="f.id">
                    {{ locale === 'ar' ? f.name_ar : f.name_en }}
                  </option>
                </select>
                <span v-if="loadingFaculties" class="reg-spinner">
                  <span class="spinner-border spinner-border-sm"></span>
                </span>
              </div>
            </div>

            <div class="reg-col" v-if="departments.length || loadingDepartments">
              <label class="reg-label">{{ t('auth.department') }}</label>
              <div class="reg-input-wrap">
                <select
                  class="reg-input reg-select"
                  v-model="form.department_id"
                  :disabled="loadingDepartments"
                >
                  <option value="">{{ t('auth.selectDepartment') }}</option>
                  <option v-for="d in departments" :key="d.id" :value="d.id">
                    {{ locale === 'ar' ? d.name_ar : d.name_en }}
                  </option>
                </select>
                <span v-if="loadingDepartments" class="reg-spinner">
                  <span class="spinner-border spinner-border-sm"></span>
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Submit -->
        <button type="submit" class="reg-submit-btn" :disabled="loading || !canSubmit">
          <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
          <i v-else class="fas fa-user-plus me-2"></i>
          {{ loading ? t('app.loading') : t('auth.registerBtn') }}
        </button>
      </form>

      <!-- Login Link -->
      <div class="reg-login-link">
        <span>{{ t('auth.hasAccount') }}</span>
        <router-link to="/login">
          {{ t('auth.loginBtn') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { useBrandingStore } from '../stores/branding'
import { publicAPI } from '../services/api'

const { t, locale } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const brandingStore = useBrandingStore()

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
.reg-page {
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  background: #ffffff;
  padding: 40px 24px;
}

.reg-card {
  width: 100%;
  max-width: 560px;
  background: #ffffff;
  border-radius: 20px;
  padding: 48px 40px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
  animation: regFadeIn 0.5s ease-out both;
}

@keyframes regFadeIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Logo */
.reg-logo {
  text-align: center;
  margin-bottom: 32px;
}

.reg-logo img {
  height: 48px;
  width: auto;
  object-fit: contain;
}

/* Heading */
.reg-heading {
  font-size: 1.75rem;
  font-weight: 700;
  color: #222222;
  margin-bottom: 4px;
  text-align: center;
}

.reg-subtitle {
  color: #6a6a6a;
  font-size: 0.95rem;
  margin-bottom: 32px;
  text-align: center;
}

/* Alert */
.reg-alert {
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

.reg-alert i {
  flex-shrink: 0;
}

/* Sections */
.reg-section {
  margin-bottom: 24px;
  padding-bottom: 24px;
  border-bottom: 1px solid #f0f0f0;
}

.reg-section-last {
  border-bottom: none;
  padding-bottom: 0;
}

.reg-section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.reg-step-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #1a5276;
  color: #ffffff;
  font-size: 0.8rem;
  font-weight: 700;
  flex-shrink: 0;
}

.reg-step-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: #222222;
}

/* Layout helpers */
.reg-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.reg-col {
  flex: 1;
  min-width: 0;
}

.reg-field {
  margin-bottom: 16px;
}

/* Inputs */
.reg-label {
  display: block;
  font-weight: 500;
  font-size: 0.85rem;
  color: #222222;
  margin-bottom: 6px;
}

.reg-input-wrap {
  position: relative;
}

.reg-input {
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

.reg-input-pw {
  padding-inline-end: 44px;
}

.reg-select {
  cursor: pointer;
  appearance: auto;
}

.reg-input:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 2px rgba(26, 82, 118, 0.15);
}

.reg-input.is-valid {
  border-color: #198754;
}

.reg-input.is-invalid {
  border-color: #dc3545;
}

.reg-toggle-pw {
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

.reg-toggle-pw:hover {
  color: #1a5276;
}

.reg-spinner {
  position: absolute;
  top: 50%;
  inset-inline-end: 12px;
  transform: translateY(-50%);
  z-index: 2;
  color: #1a5276;
}

.reg-hint {
  font-size: 0.78rem;
  color: #6a6a6a;
  margin-top: 4px;
  display: block;
}

/* Strength Bar */
.reg-strength-area {
  margin-top: 8px;
}

.reg-strength-bar {
  display: flex;
  gap: 4px;
}

.reg-strength-seg {
  flex: 1;
  height: 4px;
  border-radius: 2px;
  background: #e8e8e8;
  transition: background-color 0.3s ease;
}

.reg-strength-seg.strength-weak { background-color: #dc3545; }
.reg-strength-seg.strength-medium { background-color: #ffc107; }
.reg-strength-seg.strength-strong { background-color: #198754; }

.reg-strength-label {
  display: block;
  margin-top: 4px;
  font-weight: 500;
}

/* Password Requirements */
.reg-pw-reqs {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-top: 8px;
}

.reg-req {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  transition: color 0.2s;
}

.reg-req.met { color: #198754; }
.reg-req.unmet { color: #b0b0b0; }
.reg-req i { font-size: 0.65rem; }

/* Match status */
.reg-match-status {
  margin-top: 6px;
}

.reg-match-ok {
  color: #198754;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.8rem;
}

.reg-match-err {
  color: #dc3545;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.8rem;
}

/* Submit */
.reg-submit-btn {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 8px;
  background: #c0982b;
  border: none;
  color: #ffffff;
  cursor: pointer;
  transition: background-color 0.2s, box-shadow 0.2s;
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.reg-submit-btn:hover:not(:disabled) {
  background: #a98825;
  box-shadow:
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.12) 0px 4px 12px;
}

.reg-submit-btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

/* Login Link */
.reg-login-link {
  text-align: center;
  margin-top: 28px;
  color: #6a6a6a;
  font-size: 0.9rem;
}

.reg-login-link a {
  color: #1a5276;
  text-decoration: none;
  font-weight: 600;
  margin-inline-start: 4px;
  transition: color 0.15s;
}

.reg-login-link a:hover {
  color: #c0982b;
}

/* Responsive */
@media (max-width: 576px) {
  .reg-card {
    padding: 32px 24px;
    border-radius: 16px;
  }

  .reg-heading {
    font-size: 1.5rem;
  }

  .reg-row {
    flex-direction: column;
    gap: 0;
  }

  .reg-col {
    margin-bottom: 16px;
  }

  .reg-field[style] {
    max-width: 100% !important;
  }
}
</style>
