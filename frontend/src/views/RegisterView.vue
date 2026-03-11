<template>
  <div class="container py-5">
    <div class="row justify-content-center">
      <div class="col-md-6">
        <div class="card shadow-sm">
          <div class="card-body p-4">
            <h3 class="text-center mb-4">{{ t('auth.registerTitle') }}</h3>

            <div v-if="error" class="alert alert-danger">{{ error }}</div>

            <form @submit.prevent="handleRegister">
              <div class="row g-3">
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.fullNameAr') }}</label>
                  <input type="text" class="form-control" v-model="form.full_name_ar" required />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.fullNameEn') }}</label>
                  <input type="text" class="form-control" v-model="form.full_name_en" required />
                </div>
                <div class="col-12">
                  <label class="form-label">{{ t('auth.email') }}</label>
                  <input type="email" class="form-control" v-model="form.email" required />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.password') }}</label>
                  <input type="password" class="form-control" v-model="form.password" required minlength="6" />
                  <!-- Password strength indicator -->
                  <div class="mt-1" v-if="form.password">
                    <div class="d-flex align-items-center gap-2">
                      <div class="progress flex-grow-1" style="height: 6px;">
                        <div
                          class="progress-bar"
                          :class="passwordStrengthBarClass"
                          :style="{ width: passwordStrengthPercent + '%' }"
                        ></div>
                      </div>
                      <small :class="passwordStrengthTextClass">{{ passwordStrengthLabel }}</small>
                    </div>
                    <small class="text-muted d-block mt-1">
                      {{ locale === 'ar' ? 'الحد الأدنى 6 أحرف' : 'Minimum 6 characters' }}
                      <span v-if="form.password.length >= 6" class="text-success ms-1">&#10003;</span>
                      <span v-else class="text-danger ms-1">&#10007;</span>
                    </small>
                  </div>
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.confirmPassword') }}</label>
                  <div class="input-group">
                    <input type="password" class="form-control" :class="confirmPasswordInputClass" v-model="form.confirm_password" required />
                    <span v-if="form.confirm_password" class="input-group-text" :class="passwordsMatch ? 'text-success bg-success-subtle' : 'text-danger bg-danger-subtle'">
                      <span v-if="passwordsMatch">&#10003;</span>
                      <span v-else>&#10007;</span>
                    </span>
                  </div>
                  <small v-if="form.confirm_password && !passwordsMatch" class="text-danger">
                    {{ locale === 'ar' ? 'كلمات المرور غير متطابقة' : 'Passwords do not match' }}
                  </small>
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.phone') }}</label>
                  <input type="tel" class="form-control" v-model="form.phone" placeholder="+964XXXXXXXXXX" />
                  <small class="text-muted">{{ locale === 'ar' ? 'مثال: +964XXXXXXXXXX' : 'e.g. +964XXXXXXXXXX' }}</small>
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.faculty') }}</label>
                  <div class="position-relative">
                    <select class="form-select" v-model="form.faculty_id" @change="loadDepartments" :disabled="loadingFaculties">
                      <option value="">--</option>
                      <option v-for="f in faculties" :key="f.id" :value="f.id">
                        {{ locale === 'ar' ? f.name_ar : f.name_en }}
                      </option>
                    </select>
                    <div v-if="loadingFaculties" class="position-absolute top-50 end-0 translate-middle-y me-4">
                      <span class="spinner-border spinner-border-sm text-primary"></span>
                    </div>
                  </div>
                </div>
                <div class="col-md-6" v-if="departments.length || loadingDepartments">
                  <label class="form-label">{{ t('auth.department') }}</label>
                  <div class="position-relative">
                    <select class="form-select" v-model="form.department_id" :disabled="loadingDepartments">
                      <option value="">--</option>
                      <option v-for="d in departments" :key="d.id" :value="d.id">
                        {{ locale === 'ar' ? d.name_ar : d.name_en }}
                      </option>
                    </select>
                    <div v-if="loadingDepartments" class="position-absolute top-50 end-0 translate-middle-y me-4">
                      <span class="spinner-border spinner-border-sm text-primary"></span>
                    </div>
                  </div>
                </div>
              </div>

              <button type="submit" class="btn btn-primary w-100 mt-4" :disabled="loading">
                <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
                {{ t('auth.registerBtn') }}
              </button>
            </form>

            <p class="text-center mt-3 mb-0">
              {{ t('auth.hasAccount') }}
              <router-link to="/login">{{ t('auth.loginBtn') }}</router-link>
            </p>
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

const passwordStrengthPercent = computed(() => {
  return [0, 33, 66, 100][passwordStrength.value]
})

const passwordStrengthBarClass = computed(() => {
  return ['', 'bg-danger', 'bg-warning', 'bg-success'][passwordStrength.value]
})

const passwordStrengthTextClass = computed(() => {
  return ['', 'text-danger', 'text-warning', 'text-success'][passwordStrength.value]
})

const passwordStrengthLabel = computed(() => {
  const labels = locale.value === 'ar'
    ? ['', 'ضعيفة', 'متوسطة', 'قوية']
    : ['', 'Weak', 'Medium', 'Strong']
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
