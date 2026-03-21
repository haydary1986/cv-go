<template>
  <div class="container py-5 text-center">
    <div class="row justify-content-center">
      <div class="col-md-6">
        <i class="fas fa-exclamation-triangle fa-5x text-warning mb-4"></i>
        <h1 class="display-4 fw-bold text-muted">404</h1>
        <h3 class="mb-3">{{ t('app.pageNotFound') }}</h3>
        <p class="text-muted mb-4">{{ t('app.pageNotFoundDesc') }}</p>
        <router-link to="/" class="btn btn-primary btn-lg">
          <i class="fas fa-home me-2"></i>{{ t('app.goHome') }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>
