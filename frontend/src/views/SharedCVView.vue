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
          <button @click="exportPDF" class="btn-pill btn-pill--accent d-inline-flex" :disabled="exporting || !cv">
            <span v-if="exporting" class="spinner-border spinner-border-sm me-1"></span>
            <i v-else class="fas fa-file-pdf me-1"></i>
            {{ t('shared.downloadPDF') }}
          </button>
        </div>
      </div>
    </nav>

    <div class="container-fluid py-4 px-3 px-md-5">
      <!-- Loading State -->
      <div v-if="loading" class="shared-loading text-center">
        <div class="loading-ring mb-4">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <h5 class="fw-normal mb-2">{{ t('shared.loadingCV') }}</h5>
        <p class="text-muted small">{{ t('shared.loadingDesc') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="shared-error text-center">
        <div class="error-icon-wrapper mb-4">
          <div class="error-icon-circle">
            <i class="fas fa-exclamation-triangle"></i>
          </div>
        </div>
        <h4 class="fw-bold mb-2" style="color: var(--uni-primary);">{{ t('shared.cvNotFound') }}</h4>
        <p class="text-muted mb-4 mx-auto" style="max-width: 400px;">{{ t('shared.cvNotFoundDesc') }}</p>
        <div class="d-flex gap-2 justify-content-center">
          <router-link to="/" class="btn-pill btn-pill--outline">
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
      <p class="mb-2 fw-semibold">{{ t('shared.createYourCV') }}</p>
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
:root {
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-secondary: #2c3e50;
}

.shared-cv-page {
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-secondary: #2c3e50;
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f2f5 0%, #e8ecf1 100%);
  display: flex;
  flex-direction: column;
}

/* Top Bar */
.shared-topbar {
  background: var(--uni-primary);
  padding: 0.65rem 0;
  box-shadow: 0 2px 12px rgba(26, 82, 118, 0.18);
  z-index: 1030;
}

.topbar-brand {
  color: #fff;
}

.topbar-logo {
  height: 34px;
  width: auto;
  object-fit: contain;
  filter: brightness(0) invert(1);
}

.topbar-logo-fallback {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--uni-accent);
  font-size: 18px;
}

.topbar-name {
  font-size: 1.05rem;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.01em;
}

/* Shared Badge */
.shared-badge {
  display: inline-flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.15);
  color: var(--uni-accent);
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.3rem 0.75rem;
  border-radius: 999px;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

/* Pill Buttons */
.btn-pill {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 1.1rem;
  border-radius: 999px;
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
  border: 1.5px solid rgba(255, 255, 255, 0.4);
  color: #fff;
}

.btn-pill--outline:hover {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.7);
  color: #fff;
}

.btn-pill--accent {
  background: var(--uni-accent);
  color: #fff;
  box-shadow: 0 2px 8px rgba(192, 152, 43, 0.3);
}

.btn-pill--accent:hover {
  background: #d4a82f;
  box-shadow: 0 4px 14px rgba(192, 152, 43, 0.4);
  color: #fff;
}

.btn-pill--primary {
  background: var(--uni-primary);
  color: #fff;
  box-shadow: 0 2px 8px rgba(26, 82, 118, 0.25);
}

.btn-pill--primary:hover {
  background: #1e6291;
  box-shadow: 0 4px 14px rgba(26, 82, 118, 0.35);
  color: #fff;
}

.btn-pill:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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
  color: var(--uni-accent);
  border-width: 3px;
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
  background: linear-gradient(135deg, #fff8e1, #fff3cd);
  color: var(--uni-accent);
  font-size: 2.2rem;
  box-shadow: 0 6px 20px rgba(192, 152, 43, 0.15);
}

/* CV Title */
.cv-title-label {
  color: var(--uni-secondary);
  font-weight: 500;
}

.cv-title-label i {
  color: var(--uni-accent);
}

/* CV Frame */
.cv-frame {
  max-width: 220mm;
  background: #fff;
  border-radius: 8px;
  padding: 8px;
  box-shadow:
    0 1px 3px rgba(0, 0, 0, 0.06),
    0 8px 30px rgba(26, 82, 118, 0.08);
  border: 1px solid rgba(26, 82, 118, 0.08);
}

.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
  border-radius: 4px;
  overflow: hidden;
}

/* CTA Banner */
.cta-banner {
  background: linear-gradient(135deg, var(--uni-primary), var(--uni-secondary));
  color: #fff;
  padding: 1.5rem 1rem;
  margin-top: 2rem;
}

.cta-banner p {
  font-size: 1rem;
  opacity: 0.9;
}

.cta-banner .btn-pill--primary {
  background: var(--uni-accent);
  box-shadow: 0 2px 10px rgba(192, 152, 43, 0.3);
}

.cta-banner .btn-pill--primary:hover {
  background: #d4a82f;
}

/* Footer */
.shared-footer {
  margin-top: auto;
  background: var(--uni-primary);
  color: rgba(255, 255, 255, 0.7);
  padding: 0.85rem 0;
  border-top: 3px solid var(--uni-accent);
}

.shared-footer small {
  color: rgba(255, 255, 255, 0.7);
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

  .container-fluid {
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
