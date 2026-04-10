<template>
  <div class="reg-page">
    <div class="reg-split">
      <!-- Left: Form -->
      <div class="reg-form-side">
        <div class="reg-form-container">
          <h2 class="reg-heading">{{ t('auth.createYourAccount') }}</h2>
          <p class="reg-subtitle">{{ t('auth.registerSubtitle') }}</p>

          <!-- Error Alert -->
          <div v-if="error" class="alert alert-danger d-flex align-items-center border-0 reg-alert" role="alert">
            <i class="fas fa-exclamation-circle me-2 flex-shrink-0"></i>
            <div>{{ error }}</div>
          </div>

          <form @submit.prevent="handleRegister">
            <!-- Step 1: Account Info -->
            <div class="reg-section">
              <div class="reg-section-header">
                <span class="reg-step-badge">1</span>
                <span class="reg-step-title">{{ locale === 'ar' ? 'معلومات الحساب' : 'Account Information' }}</span>
              </div>

              <div class="mb-3">
                <label class="reg-label">{{ t('auth.email') }}</label>
                <div class="reg-input-wrap">
                  <i class="fas fa-envelope reg-input-icon"></i>
                  <input
                    type="email"
                    class="form-control reg-input"
                    v-model="form.email"
                    required
                    autocomplete="email"
                  />
                </div>
              </div>

              <div class="row g-3">
                <div class="col-md-6">
                  <label class="reg-label">{{ t('auth.password') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-lock reg-input-icon"></i>
                    <input
                      :type="showPassword ? 'text' : 'password'"
                      class="form-control reg-input reg-input-pw"
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
                  <div class="mt-2" v-if="form.password">
                    <div class="reg-strength-bar">
                      <div class="reg-strength-seg" :class="passwordStrength >= 1 ? strengthColorClass : ''"></div>
                      <div class="reg-strength-seg" :class="passwordStrength >= 2 ? strengthColorClass : ''"></div>
                      <div class="reg-strength-seg" :class="passwordStrength >= 3 ? strengthColorClass : ''"></div>
                    </div>
                    <div class="d-flex justify-content-between align-items-center mt-1">
                      <small :class="passwordStrengthTextClass" class="fw-medium">{{ passwordStrengthLabel }}</small>
                    </div>
                    <div class="reg-pw-reqs mt-2">
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

                <div class="col-md-6">
                  <label class="reg-label">{{ t('auth.confirmPassword') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-lock reg-input-icon"></i>
                    <input
                      :type="showConfirmPassword ? 'text' : 'password'"
                      class="form-control reg-input reg-input-pw"
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
                  <div class="mt-2" v-if="form.confirm_password">
                    <small v-if="passwordsMatch" class="text-success d-flex align-items-center gap-1">
                      <i class="fas fa-check-circle"></i> {{ locale === 'ar' ? 'كلمات المرور متطابقة' : 'Passwords match' }}
                    </small>
                    <small v-else class="text-danger d-flex align-items-center gap-1">
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

              <div class="row g-3">
                <div class="col-md-6">
                  <label class="reg-label">{{ t('auth.fullNameAr') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-user reg-input-icon"></i>
                    <input
                      type="text"
                      class="form-control reg-input"
                      v-model="form.full_name_ar"
                      required
                    />
                  </div>
                </div>

                <div class="col-md-6">
                  <label class="reg-label">{{ t('auth.fullNameEn') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-user reg-input-icon"></i>
                    <input
                      type="text"
                      class="form-control reg-input"
                      v-model="form.full_name_en"
                      required
                    />
                  </div>
                </div>

                <div class="col-md-6">
                  <label class="reg-label">{{ t('auth.phone') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-phone reg-input-icon"></i>
                    <input
                      type="tel"
                      class="form-control reg-input"
                      v-model="form.phone"
                      placeholder="+964XXXXXXXXXX"
                      dir="ltr"
                    />
                  </div>
                  <small class="text-muted mt-1 d-block" style="font-size: 0.75rem;">
                    <i class="fas fa-info-circle me-1"></i>{{ t('auth.phoneHint') }}
                  </small>
                </div>
              </div>
            </div>

            <!-- Step 3: University -->
            <div class="reg-section">
              <div class="reg-section-header">
                <span class="reg-step-badge">3</span>
                <span class="reg-step-title">{{ locale === 'ar' ? 'المعلومات الجامعية' : 'University Details' }}</span>
              </div>

              <div class="row g-3">
                <div class="col-md-6">
                  <label class="reg-label">{{ t('auth.faculty') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-university reg-input-icon"></i>
                    <select
                      class="form-select reg-input reg-select"
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
                      <span class="spinner-border spinner-border-sm text-primary"></span>
                    </span>
                  </div>
                </div>

                <div class="col-md-6" v-if="departments.length || loadingDepartments">
                  <label class="reg-label">{{ t('auth.department') }}</label>
                  <div class="reg-input-wrap">
                    <i class="fas fa-building reg-input-icon"></i>
                    <select
                      class="form-select reg-input reg-select"
                      v-model="form.department_id"
                      :disabled="loadingDepartments"
                    >
                      <option value="">{{ t('auth.selectDepartment') }}</option>
                      <option v-for="d in departments" :key="d.id" :value="d.id">
                        {{ locale === 'ar' ? d.name_ar : d.name_en }}
                      </option>
                    </select>
                    <span v-if="loadingDepartments" class="reg-spinner">
                      <span class="spinner-border spinner-border-sm text-primary"></span>
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Submit -->
            <button type="submit" class="btn reg-submit-btn w-100 fw-semibold" :disabled="loading || !canSubmit">
              <span v-if="loading" class="spinner-border spinner-border-sm me-2" role="status"></span>
              <i v-else class="fas fa-user-plus me-2"></i>
              {{ loading ? t('app.loading') : t('auth.registerBtn') }}
            </button>
          </form>

          <!-- Login Link -->
          <div class="reg-login-link">
            <span>{{ t('auth.hasAccount') }}</span>
            <router-link to="/login" class="fw-semibold">
              {{ t('auth.loginBtn') }}
            </router-link>
          </div>
        </div>
      </div>

      <!-- Right: University Branding -->
      <div class="reg-brand-side">
        <div class="reg-brand-content">
          <div class="reg-brand-logo">
            <img :src="brandingStore.logoUrl || '/logo.svg'" alt="Logo" />
          </div>
          <h1 class="reg-brand-name">{{ brandingStore.systemName }}</h1>
          <p class="reg-brand-tagline">{{ t('auth.registerSubtitle') }}</p>
          <div class="reg-brand-features">
            <div class="reg-brand-feature">
              <i class="fas fa-file-alt"></i>
              <span>{{ locale === 'ar' ? 'سيرة ذاتية احترافية' : 'Professional CV' }}</span>
            </div>
            <div class="reg-brand-feature">
              <i class="fas fa-robot"></i>
              <span>{{ locale === 'ar' ? 'مدعوم بالذكاء الاصطناعي' : 'AI-Powered' }}</span>
            </div>
            <div class="reg-brand-feature">
              <i class="fas fa-globe"></i>
              <span>{{ locale === 'ar' ? 'عربي وإنجليزي' : 'Arabic & English' }}</span>
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
/* ── Layout ── */
.reg-page {
  min-height: 100vh;
}

.reg-split {
  display: flex;
  min-height: 100vh;
}

/* ── Left: Form Side ── */
.reg-form-side {
  flex: 1.3;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 2.5rem 2rem;
  background: #fff;
  overflow-y: auto;
  animation: regFadeIn 0.7s ease-out both;
}

@keyframes regFadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.reg-form-container {
  width: 100%;
  max-width: 580px;
  padding-bottom: 2rem;
}

.reg-heading {
  font-size: 1.85rem;
  font-weight: 800;
  color: #1a5276;
  margin-bottom: 0.25rem;
  letter-spacing: -0.5px;
}

.reg-subtitle {
  color: #6c757d;
  font-size: 0.95rem;
  margin-bottom: 2rem;
}

.reg-alert {
  border-radius: 10px;
  background: #fef2f2;
  color: #b91c1c;
  margin-bottom: 1.5rem;
}

/* ── Sections ── */
.reg-section {
  margin-bottom: 1.75rem;
  padding-bottom: 1.75rem;
  border-bottom: 1px solid #f0f0f0;
}

.reg-section:last-of-type {
  border-bottom: none;
  padding-bottom: 0;
}

.reg-section-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 1rem;
}

.reg-step-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1a5276, #0d6efd);
  color: #fff;
  font-size: 0.8rem;
  font-weight: 700;
  flex-shrink: 0;
}

.reg-step-title {
  font-size: 1rem;
  font-weight: 700;
  color: #2c3e50;
}

/* ── Inputs ── */
.reg-label {
  font-weight: 600;
  font-size: 0.85rem;
  color: #2c3e50;
  margin-bottom: 0.35rem;
  display: block;
}

.reg-input-wrap {
  position: relative;
}

.reg-input-icon {
  position: absolute;
  top: 50%;
  inset-inline-start: 14px;
  transform: translateY(-50%);
  color: #9ca3af;
  font-size: 0.9rem;
  pointer-events: none;
  z-index: 2;
}

.reg-input {
  height: 46px;
  padding-inline-start: 42px;
  padding-inline-end: 14px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 0.95rem;
  background: #f9fafb;
  transition: border-color 0.2s, background-color 0.2s, box-shadow 0.2s;
}

.reg-input-pw {
  padding-inline-end: 42px;
}

.reg-select {
  cursor: pointer;
}

.reg-input:focus {
  border-color: #1a5276;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
  outline: none;
}

.reg-toggle-pw {
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

.reg-toggle-pw:hover {
  color: #1a5276;
}

.reg-spinner {
  position: absolute;
  top: 50%;
  inset-inline-end: 12px;
  transform: translateY(-50%);
  z-index: 2;
}

/* ── Strength Bar ── */
.reg-strength-bar {
  display: flex;
  gap: 4px;
}

.reg-strength-seg {
  flex: 1;
  height: 5px;
  border-radius: 3px;
  background: #e9ecef;
  transition: background-color 0.3s ease;
}

.reg-strength-seg.strength-weak { background-color: #dc3545; }
.reg-strength-seg.strength-medium { background-color: #ffc107; }
.reg-strength-seg.strength-strong { background-color: #198754; }

/* ── Password Requirements ── */
.reg-pw-reqs {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.reg-req {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.75rem;
  transition: color 0.2s;
}

.reg-req.met { color: #198754; }
.reg-req.unmet { color: #adb5bd; }
.reg-req i { font-size: 0.65rem; }

/* ── Submit ── */
.reg-submit-btn {
  height: 50px;
  font-size: 1.05rem;
  border-radius: 12px;
  background: linear-gradient(135deg, #c0982b 0%, #d4a933 100%);
  border: none;
  color: #fff;
  letter-spacing: 0.3px;
  transition: transform 0.15s, box-shadow 0.2s;
  margin-top: 0.5rem;
}

.reg-submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(192, 152, 43, 0.35);
  background: linear-gradient(135deg, #b08a24 0%, #c0982b 100%);
  color: #fff;
}

.reg-submit-btn:disabled {
  opacity: 0.55;
  color: #fff;
}

/* ── Login Link ── */
.reg-login-link {
  text-align: center;
  margin-top: 1.75rem;
  color: #6c757d;
  font-size: 0.95rem;
}

.reg-login-link a {
  color: #1a5276;
  text-decoration: none;
  margin-inline-start: 6px;
  transition: color 0.15s;
}

.reg-login-link a:hover {
  color: #c0982b;
}

/* ── Right: Branding ── */
.reg-brand-side {
  flex: 0.7;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(160deg, #1a5276 0%, #0d6efd 100%);
  position: relative;
  overflow: hidden;
  position: sticky;
  top: 0;
  height: 100vh;
}

.reg-brand-side::before {
  content: '';
  position: absolute;
  top: -25%;
  right: -15%;
  width: 450px;
  height: 450px;
  border-radius: 50%;
  background: rgba(192, 152, 43, 0.08);
}

.reg-brand-content {
  position: relative;
  z-index: 1;
  text-align: center;
  padding: 3rem;
  animation: regBrandFade 0.9s ease-out 0.3s both;
}

@keyframes regBrandFade {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

.reg-brand-logo {
  margin-bottom: 2rem;
}

.reg-brand-logo img {
  height: 110px;
  width: auto;
  object-fit: contain;
  filter: drop-shadow(0 4px 20px rgba(0, 0, 0, 0.15));
}

.reg-brand-name {
  font-size: 1.75rem;
  font-weight: 800;
  color: #fff;
  margin-bottom: 0.5rem;
}

.reg-brand-tagline {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
  max-width: 280px;
  margin: 0 auto 2.5rem;
  line-height: 1.5;
}

.reg-brand-features {
  display: flex;
  flex-direction: column;
  gap: 14px;
  align-items: center;
}

.reg-brand-feature {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.85);
  font-size: 0.95rem;
}

.reg-brand-feature i {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: rgba(192, 152, 43, 0.2);
  color: #c0982b;
  font-size: 0.85rem;
}

/* ── Responsive ── */
@media (max-width: 991.98px) {
  .reg-brand-side {
    display: none;
  }

  .reg-form-side {
    padding: 2rem 1.25rem;
  }
}

@media (max-width: 576px) {
  .reg-heading {
    font-size: 1.4rem;
  }

  .reg-form-container {
    max-width: 100%;
  }
}

/* ── RTL ── */
[dir="rtl"] .reg-split {
  flex-direction: row-reverse;
}
</style>
