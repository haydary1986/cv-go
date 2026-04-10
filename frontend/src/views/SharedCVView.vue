<template>
  <div class="shared-cv-page">
    <!-- Top Bar -->
    <nav class="shared-topbar sticky-top">
      <div class="container-fluid px-3 px-md-4 d-flex align-items-center justify-content-between">
        <router-link to="/" class="topbar-brand d-flex align-items-center gap-2 text-decoration-none">
          <template v-if="brandingStore.logoUrl">
            <img :src="brandingStore.logoUrl" alt="Logo" class="topbar-logo" />
          </template>
          <div v-else class="topbar-logo-fallback">
            <i class="fas fa-university"></i>
          </div>
          <span class="topbar-name">{{ brandingStore.systemName }}</span>
        </router-link>

        <div class="d-flex align-items-center gap-2">
          <span class="shared-badge d-none d-sm-inline-flex">
            <i class="fas fa-share-alt me-1"></i>
            {{ t('shared.sharedCV') || 'Shared CV' }}
          </span>
          <button @click="printCV" class="btn-pill btn-pill--outline d-none d-sm-inline-flex" :disabled="!cv">
            <i class="fas fa-print me-1"></i>
            {{ t('shared.printCV') }}
          </button>
          <button @click="exportPDF" class="btn-pill btn-pill--primary d-inline-flex" :disabled="exporting || !cv">
            <span v-if="exporting" class="spinner-border spinner-border-sm me-1"></span>
            <i v-else class="fas fa-file-pdf me-1"></i>
            {{ t('shared.downloadPDF') }}
          </button>
        </div>
      </div>
    </nav>

    <div class="shared-body py-5 px-3 px-md-5">
      <!-- Loading State -->
      <div v-if="loading" class="shared-loading text-center">
        <div class="loading-ring mb-4">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <h5 class="loading-title mb-2">{{ t('shared.loadingCV') }}</h5>
        <p class="loading-desc">{{ t('shared.loadingDesc') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="shared-error text-center">
        <div class="error-icon-wrapper mb-4">
          <div class="error-icon-circle">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
        </div>
        <h4 class="error-title mb-2">{{ t('shared.cvNotFound') }}</h4>
        <p class="error-desc mb-4 mx-auto">{{ t('shared.cvNotFoundDesc') }}</p>
        <div class="d-flex gap-2 justify-content-center">
          <router-link to="/" class="btn-pill btn-pill--outline-dark">
            <i class="fas fa-home me-2"></i>{{ t('shared.goHome') }}
          </router-link>
          <router-link to="/dashboard" class="btn-pill btn-pill--primary">
            <i class="fas fa-plus me-2"></i>{{ t('shared.createYourCV') }}
          </router-link>
        </div>
      </div>

      <!-- CV Content -->
      <div v-else-if="cv">
        <div class="text-center mb-4">
          <h5 class="cv-title-label">
            <i class="fas fa-file-alt me-2"></i>{{ cv.title || t('shared.untitledCV') }}
          </h5>
        </div>
        <div class="cv-frame mx-auto">
          <div class="cv-preview-container" ref="cvContainer">
            <CVTemplates :data="cv.data" :template="cv.template" :language="cv.language" />
          </div>
        </div>
      </div>
    </div>

    <!-- CTA Banner -->
    <div v-if="cv && !loading" class="cta-banner text-center">
      <p class="cta-text mb-2">{{ t('shared.createYourCV') }}</p>
      <router-link to="/dashboard" class="btn-pill btn-pill--primary">
        <i class="fas fa-plus me-1"></i>
        {{ t('shared.createYourCV') }}
      </router-link>
    </div>

    <!-- Footer -->
    <footer class="shared-footer text-center" v-if="!loading">
      <small>
        <i class="fas fa-university me-1"></i>{{ t('shared.poweredBy') }}
      </small>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { cvAPI } from '../services/api'
import { useBrandingStore } from '../stores/branding'
import CVTemplates from '../components/CVTemplates.vue'

const { t } = useI18n()
const route = useRoute()
const brandingStore = useBrandingStore()
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
  background: #ffffff;
  display: flex;
  flex-direction: column;
}

/* Top Bar */
.shared-topbar {
  background: #ffffff;
  padding: 0.65rem 0;
  border-bottom: 1px solid #ebebeb;
  z-index: 1030;
}

.topbar-brand {
  color: #222222;
}

.topbar-logo {
  height: 34px;
  width: auto;
  object-fit: contain;
}

.topbar-logo-fallback {
  width: 38px;
  height: 38px;
  border-radius: 8px;
  background: rgba(26, 82, 118, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1a5276;
  font-size: 18px;
}

.topbar-name {
  font-size: 1.05rem;
  font-weight: 700;
  color: #222222;
}

/* Shared Badge */
.shared-badge {
  display: inline-flex;
  align-items: center;
  background: #f7f7f7;
  color: #6a6a6a;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.3rem 0.75rem;
  border-radius: 20px;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

/* Pill Buttons */
.btn-pill {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 1.1rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.btn-pill--outline {
  background: transparent;
  border: 1px solid #dddddd;
  color: #222222;
}

.btn-pill--outline:hover {
  background: #f7f7f7;
  border-color: #222222;
  color: #222222;
}

.btn-pill--outline-dark {
  background: transparent;
  border: 1px solid #222222;
  color: #222222;
}

.btn-pill--outline-dark:hover {
  background: #222222;
  color: #ffffff;
}

.btn-pill--primary {
  background: #1a5276;
  color: #ffffff;
}

.btn-pill--primary:hover {
  background: #154360;
  color: #ffffff;
}

.btn-pill:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Shared Body */
.shared-body {
  flex: 1;
}

/* Loading State */
.shared-loading {
  max-width: 400px;
  margin: 0 auto;
  padding-top: 12vh;
}

.loading-ring .spinner-border {
  width: 3rem;
  height: 3rem;
  color: #1a5276;
  border-width: 3px;
}

.loading-title {
  color: #222222;
  font-weight: 500;
}

.loading-desc {
  color: #6a6a6a;
  font-size: 0.9rem;
}

/* Error State */
.shared-error {
  padding-top: 12vh;
}

.error-icon-circle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 90px;
  height: 90px;
  border-radius: 50%;
  background: #f7f7f7;
  color: #c0982b;
  font-size: 2.2rem;
}

.error-title {
  font-weight: 700;
  color: #222222;
}

.error-desc {
  color: #6a6a6a;
  max-width: 400px;
}

/* CV Title */
.cv-title-label {
  color: #222222;
  font-weight: 600;
}

.cv-title-label i {
  color: #c0982b;
}

/* CV Frame */
.cv-frame {
  max-width: 800px;
  background: #ffffff;
  border-radius: 12px;
  padding: 8px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
  border-radius: 8px;
  overflow: hidden;
}

/* CTA Banner */
.cta-banner {
  background: #f7f7f7;
  padding: 2rem 1rem;
  margin-top: 2rem;
  border-top: 1px solid #ebebeb;
}

.cta-text {
  font-size: 1rem;
  color: #222222;
  font-weight: 600;
}

/* Footer */
.shared-footer {
  margin-top: auto;
  padding: 1rem 0;
  border-top: 1px solid #ebebeb;
}

.shared-footer small {
  color: #6a6a6a;
}

/* Print styles */
@media print {
  .shared-topbar,
  .shared-footer,
  .cta-banner,
  .shared-badge,
  .btn-pill {
    display: none !important;
  }

  .shared-cv-page {
    background: #fff;
  }

  .cv-frame {
    box-shadow: none !important;
    border: none !important;
    padding: 0 !important;
    max-width: 100% !important;
  }

  .cv-preview-container {
    box-shadow: none !important;
    margin: 0 !important;
    max-width: 100% !important;
  }

  .container-fluid,
  .shared-body {
    padding: 0 !important;
  }
}

@media (max-width: 576px) {
  .topbar-name {
    display: none;
  }

  .btn-pill {
    padding: 0.35rem 0.8rem;
    font-size: 0.8rem;
  }
}
</style>
