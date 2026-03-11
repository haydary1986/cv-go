<template>
  <div class="container py-5">
    <div class="row justify-content-center">
      <div class="col-md-5">
        <div class="card shadow-sm">
          <div class="card-body p-4">
            <h3 class="text-center mb-4">{{ t('auth.loginTitle') }}</h3>

            <div v-if="error" class="alert alert-danger">{{ error }}</div>

            <form @submit.prevent="handleLogin">
              <div class="mb-3">
                <label class="form-label">{{ t('auth.email') }}</label>
                <input type="email" class="form-control" v-model="form.email" required />
              </div>
              <div class="mb-3">
                <label class="form-label">{{ t('auth.password') }}</label>
                <input type="password" class="form-control" v-model="form.password" required />
              </div>
              <button type="submit" class="btn btn-primary w-100 mb-3" :disabled="loading">
                <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
                {{ t('auth.loginBtn') }}
              </button>
            </form>

            <div class="text-center mb-3">
              <hr />
              <button @click="handleGoogleLogin" class="btn btn-outline-danger w-100">
                <i class="fab fa-google me-2"></i>{{ t('auth.googleLogin') }}
              </button>
            </div>

            <p class="text-center mb-0">
              {{ t('auth.noAccount') }}
              <router-link to="/register">{{ t('auth.registerBtn') }}</router-link>
            </p>
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
import { authAPI } from '../services/api'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loading = ref(false)
const error = ref('')
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
