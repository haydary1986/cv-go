<template>
  <div class="auth-callback">
    <div class="auth-callback-spinner"></div>
    <p class="auth-callback-text">{{ t('app.loading') }}</p>
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

<style scoped>
.auth-callback {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  gap: 16px;
}

.auth-callback-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e8e8e8;
  border-top-color: #1a5276;
  border-radius: 50%;
  animation: authSpin 0.8s linear infinite;
}

@keyframes authSpin {
  to {
    transform: rotate(360deg);
  }
}

.auth-callback-text {
  color: #6a6a6a;
  font-size: 0.95rem;
  margin: 0;
}
</style>
