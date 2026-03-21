<template>
  <div class="shared-cv-page">
    <!-- Top Navbar -->
    <nav class="navbar navbar-expand navbar-light bg-white shadow-sm sticky-top">
      <div class="container-fluid px-3 px-md-4">
        <router-link to="/" class="navbar-brand d-flex align-items-center gap-2 fw-bold">
          <div class="nav-brand-icon">
            <i class="fas fa-file-alt"></i>
          </div>
          <span class="brand-text">CV Builder</span>
        </router-link>
        <div class="d-flex align-items-center gap-2">
          <button @click="printCV" class="btn btn-outline-secondary btn-sm d-none d-sm-inline-flex align-items-center" :disabled="!cv">
            <i class="fas fa-print me-1"></i>
            {{ t('shared.printCV') }}
          </button>
          <button @click="exportPDF" class="btn btn-outline-danger btn-sm d-inline-flex align-items-center" :disabled="exporting || !cv">
            <span v-if="exporting" class="spinner-border spinner-border-sm me-1"></span>
            <i v-else class="fas fa-file-pdf me-1"></i>
            {{ t('shared.downloadPDF') }}
          </button>
          <router-link to="/dashboard" class="btn btn-primary btn-sm d-inline-flex align-items-center">
            <i class="fas fa-plus me-1"></i>
            <span class="d-none d-sm-inline">{{ t('shared.createYourCV') }}</span>
            <span class="d-sm-none">{{ t('app.createCV') }}</span>
          </router-link>
        </div>
      </div>
    </nav>

    <div class="container-fluid py-4 px-3 px-md-4">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state text-center py-5">
        <div class="loading-spinner mb-4">
          <div class="spinner-border text-primary" style="width: 3rem; height: 3rem;" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <h5 class="text-muted fw-normal mb-2">{{ t('shared.loadingCV') }}</h5>
        <p class="text-muted small">{{ t('shared.loadingDesc') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state text-center py-5">
        <div class="error-icon-wrapper mb-4">
          <div class="error-icon-circle">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
        </div>
        <h4 class="fw-semibold mb-2">{{ t('shared.cvNotFound') }}</h4>
        <p class="text-muted mb-4 mx-auto" style="max-width: 400px;">{{ t('shared.cvNotFoundDesc') }}</p>
        <div class="d-flex gap-2 justify-content-center">
          <router-link to="/" class="btn btn-outline-secondary">
            <i class="fas fa-home me-2"></i>{{ t('shared.goHome') }}
          </router-link>
          <router-link to="/dashboard" class="btn btn-primary">
            <i class="fas fa-plus me-2"></i>{{ t('shared.createYourCV') }}
          </router-link>
        </div>
      </div>

      <!-- CV Content -->
      <div v-else-if="cv">
        <div class="text-center mb-3">
          <h5 class="text-muted fw-normal">
            <i class="fas fa-file-alt me-2 text-primary"></i>{{ cv.title || t('shared.untitledCV') }}
          </h5>
        </div>
        <div class="cv-preview-container bg-white shadow mx-auto" ref="cvContainer">
          <CVTemplates :data="cv.data" :template="cv.template" :language="cv.language" />
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="shared-footer text-center py-3" v-if="!loading">
      <small class="text-muted">
        <i class="fas fa-file-alt me-1"></i>{{ t('shared.poweredBy') }}
      </small>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { cvAPI } from '../services/api'
import CVTemplates from '../components/CVTemplates.vue'

const { t } = useI18n()
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

function printCV() {
  window.print()
}

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
  background: #f0f2f5;
  display: flex;
  flex-direction: column;
}

.nav-brand-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #2d5196, #3b82c4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 16px;
}

.brand-text {
  font-size: 1.1rem;
  color: #1a1c4e;
}

.loading-state {
  max-width: 400px;
  margin: 0 auto;
  padding-top: 10vh !important;
}

.error-state {
  padding-top: 10vh !important;
}

.error-icon-circle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #fff3cd;
  color: #856d00;
  font-size: 2rem;
  box-shadow: 0 4px 14px rgba(133, 109, 0, 0.15);
}

.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
  border-radius: 4px;
  overflow: hidden;
}

.shared-footer {
  margin-top: auto;
  background: #fff;
  border-top: 1px solid #e9ecef;
}

/* Print styles */
@media print {
  .navbar,
  .shared-footer {
    display: none !important;
  }

  .shared-cv-page {
    background: #fff;
  }

  .cv-preview-container {
    box-shadow: none !important;
    margin: 0 !important;
    max-width: 100% !important;
  }

  .container-fluid {
    padding: 0 !important;
  }
}

@media (max-width: 576px) {
  .brand-text {
    display: none;
  }
}
</style>
