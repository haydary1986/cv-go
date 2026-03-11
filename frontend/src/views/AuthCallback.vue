<template>
  <div class="container py-5 text-center">
    <div class="spinner-border text-primary mb-3"></div>
    <p>{{ t('app.loading') }}</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

onMounted(() => {
  const token = route.query.token as string
  if (token) {
    authStore.setGoogleAuth(token, null)
    authStore.fetchProfile().then(() => {
      router.push('/dashboard')
    })
  } else {
    router.push('/login')
  }
})
</script>
