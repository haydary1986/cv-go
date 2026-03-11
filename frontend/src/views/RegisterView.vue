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
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.confirmPassword') }}</label>
                  <input type="password" class="form-control" v-model="form.confirm_password" required />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.phone') }}</label>
                  <input type="tel" class="form-control" v-model="form.phone" />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.faculty') }}</label>
                  <select class="form-select" v-model="form.faculty_id" @change="loadDepartments">
                    <option value="">--</option>
                    <option v-for="f in faculties" :key="f.id" :value="f.id">
                      {{ locale === 'ar' ? f.name_ar : f.name_en }}
                    </option>
                  </select>
                </div>
                <div class="col-md-6" v-if="departments.length">
                  <label class="form-label">{{ t('auth.department') }}</label>
                  <select class="form-select" v-model="form.department_id">
                    <option value="">--</option>
                    <option v-for="d in departments" :key="d.id" :value="d.id">
                      {{ locale === 'ar' ? d.name_ar : d.name_en }}
                    </option>
                  </select>
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { publicAPI } from '../services/api'

const { t, locale } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
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

onMounted(async () => {
  try {
    const res = await publicAPI.getFaculties()
    faculties.value = res.data.faculties || []
  } catch {}
})

async function loadDepartments() {
  if (!form.faculty_id) {
    departments.value = []
    return
  }
  try {
    const res = await publicAPI.getDepartments(Number(form.faculty_id))
    departments.value = res.data.departments || []
  } catch {}
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
