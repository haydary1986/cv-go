<template>
  <div class="shared-cv-page">
    <!-- Top Toolbar -->
    <nav class="navbar navbar-light bg-white shadow-sm sticky-top">
      <div class="container-fluid">
        <router-link to="/" class="navbar-brand fw-bold">
          <i class="fas fa-file-alt text-primary me-2"></i>CV-Go
        </router-link>
        <div class="d-flex gap-2">
          <button @click="exportPDF" class="btn btn-outline-danger btn-sm" :disabled="exporting || !cv">
            <span v-if="exporting" class="spinner-border spinner-border-sm me-1"></span>
            <i v-else class="fas fa-file-pdf me-1"></i>Download PDF
          </button>
          <router-link to="/dashboard" class="btn btn-primary btn-sm">
            <i class="fas fa-plus me-1"></i>Create your CV
          </router-link>
        </div>
      </div>
    </nav>

    <div class="container-fluid py-4">
      <!-- Loading -->
      <div v-if="loading" class="text-center py-5">
        <div class="spinner-border text-primary mb-3"></div>
        <p class="text-muted">Loading CV...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-5">
        <div class="mb-4">
          <i class="fas fa-exclamation-triangle fa-4x text-warning"></i>
        </div>
        <h4 class="mb-2">{{ error }}</h4>
        <p class="text-muted mb-4">The CV you're looking for may have been removed or sharing may be disabled.</p>
        <router-link to="/" class="btn btn-primary">
          <i class="fas fa-home me-2"></i>Go Home
        </router-link>
      </div>

      <!-- CV Content -->
      <div v-else-if="cv">
        <h5 class="text-center text-muted mb-3">{{ cv.title || 'Untitled CV' }}</h5>
        <div class="cv-preview-container bg-white shadow mx-auto" ref="cvContainer">
          <CVTemplates :data="cv.data" :template="cv.template" :language="cv.language" />
        </div>
      </div>
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
const exporting = ref(false)
const cvContainer = ref<HTMLElement>()

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

async function exportPDF() {
  exporting.value = true
  try {
    const html2canvas = (await import('html2canvas')).default
    const { jsPDF } = await import('jspdf')
    const element = cvContainer.value
    if (!element) return

    const canvas = await html2canvas(element, { scale: 2, useCORS: true })
    const imgData = canvas.toDataURL('image/png')
    const pdf = new jsPDF('p', 'mm', 'a4')
    const pdfWidth = pdf.internal.pageSize.getWidth()
    const pdfHeight = pdf.internal.pageSize.getHeight()
    const imgWidth = pdfWidth
    const imgHeight = (canvas.height * pdfWidth) / canvas.width

    let heightLeft = imgHeight
    let position = 0

    pdf.addImage(imgData, 'PNG', 0, position, imgWidth, imgHeight)
    heightLeft -= pdfHeight

    while (heightLeft > 0) {
      position = -(imgHeight - heightLeft)
      pdf.addPage()
      pdf.addImage(imgData, 'PNG', 0, position, imgWidth, imgHeight)
      heightLeft -= pdfHeight
    }

    const userName = cv.value?.data?.personal_info?.full_name || 'CV'
    pdf.save(`${userName}-CV.pdf`)
  } finally {
    exporting.value = false
  }
}
</script>

<style scoped>
.shared-cv-page {
  min-height: 100vh;
  background: #f5f5f5;
}
.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
}
</style>
