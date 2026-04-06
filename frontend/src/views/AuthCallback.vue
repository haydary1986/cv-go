<template>
  <div class="container py-5 text-center">
    <div class="spinner-border text-primary mb-3"></div>
    <p>{{ t('app.loading') }}</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  // Extract token from URL fragment (hash) for security
  const hash = window.location.hash.substring(1) // remove #
  const params = new URLSearchParams(hash)
  const token = params.get('token')
  if (token) {
    // Basic JWT format validation (3 base64 parts separated by dots)
    if (!/^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$/.test(token)) {
      router.push('/login?error=invalid_token')
      return
    }
    authStore.setGoogleAuth(token, null)
    try {
      await authStore.fetchProfile()
      router.push('/dashboard')
    } catch {
      router.push('/login?error=auth_failed')
    }
  } else {
    router.push('/login')
  }
})
</script>
