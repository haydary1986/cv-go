<template>
  <div class="container py-4">
    <div class="row justify-content-center">
      <div class="col-md-8">
        <div class="card shadow-sm">
          <div class="card-header">
            <h4 class="mb-0">{{ t('app.profile') }}</h4>
          </div>
          <div class="card-body">
            <div v-if="success" class="alert alert-success">{{ success }}</div>
            <div v-if="error" class="alert alert-danger">{{ error }}</div>

            <form @submit.prevent="updateProfile">
              <div class="row g-3">
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.fullNameAr') }}</label>
                  <input type="text" class="form-control" v-model="form.full_name_ar" />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.fullNameEn') }}</label>
                  <input type="text" class="form-control" v-model="form.full_name_en" />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.email') }}</label>
                  <input type="email" class="form-control" :value="authStore.user?.email" disabled />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.phone') }}</label>
                  <input type="tel" class="form-control" v-model="form.phone" />
                </div>
              </div>
              <button type="submit" class="btn btn-primary mt-3">{{ t('app.save') }}</button>
            </form>

            <hr class="my-4" />

            <h5>{{ t('auth.changePassword') }}</h5>
            <form @submit.prevent="changePassword">
              <div class="row g-3">
                <div class="col-md-4">
                  <label class="form-label">{{ t('auth.currentPassword') }}</label>
                  <input type="password" class="form-control" v-model="pwForm.current_password" required />
                </div>
                <div class="col-md-4">
                  <label class="form-label">{{ t('auth.newPassword') }}</label>
                  <input type="password" class="form-control" v-model="pwForm.new_password" required minlength="6" />
                </div>
                <div class="col-md-4">
                  <label class="form-label">{{ t('auth.confirmPassword') }}</label>
                  <input type="password" class="form-control" v-model="pwForm.confirm_password" required />
                </div>
              </div>
              <button type="submit" class="btn btn-warning mt-3">{{ t('auth.changePassword') }}</button>
            </form>

            <hr class="my-4" />
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <strong>{{ t('ai.credits') }}:</strong>
                <span class="badge bg-primary ms-2">{{ authStore.user?.ai_credits || 0 }}</span>
              </div>
              <div>
                <strong>{{ t('cv.status') }}:</strong>
                <span class="badge bg-success ms-2">{{ authStore.user?.role }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { authAPI } from '../services/api'

const { t } = useI18n()
const authStore = useAuthStore()

const success = ref('')
const error = ref('')
const form = reactive({ full_name_ar: '', full_name_en: '', phone: '' })
const pwForm = reactive({ current_password: '', new_password: '', confirm_password: '' })

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
  try {
    await authStore.updateProfile(form)
    success.value = t('app.success')
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  }
}

async function changePassword() {
  error.value = ''
  success.value = ''
  if (pwForm.new_password !== pwForm.confirm_password) {
    error.value = t('auth.passwordMismatch')
    return
  }
  try {
    await authAPI.changePassword(pwForm)
    success.value = t('app.success')
    pwForm.current_password = ''
    pwForm.new_password = ''
    pwForm.confirm_password = ''
  } catch (err: any) {
    error.value = err.response?.data?.error || t('app.error')
  }
}
</script>
