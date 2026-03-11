<template>
  <div class="container-fluid py-4">
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary"></div>
    </div>
    <div v-else-if="error" class="text-center py-5">
      <i class="fas fa-exclamation-triangle fa-3x text-warning mb-3"></i>
      <h4>{{ error }}</h4>
    </div>
    <div v-else-if="cv" class="cv-preview-container bg-white shadow mx-auto">
      <CVTemplates :data="cv.data" :template="cv.template" :language="cv.language" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { cvAPI } from '../services/api'
import CVTemplates from '../components/CVTemplates.vue'

const route = useRoute()
const loading = ref(true)
const cv = ref<any>(null)
const error = ref('')

onMounted(async () => {
  try {
    const token = route.params.token as string
    const res = await cvAPI.getShared(token)
    cv.value = res.data.cv
  } catch {
    error.value = 'CV not found or sharing is disabled'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
}
</style>
