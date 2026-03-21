<template>
  <div class="preview-page">
    <div class="container-fluid py-4">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="loading-content text-center py-5">
          <div class="loading-spinner mb-3">
            <div class="spinner-border text-primary" style="width: 3rem; height: 3rem;"></div>
          </div>
          <h5 class="text-muted">{{ t('preview.loadingCV') }}</h5>
          <div class="loading-bar mt-3 mx-auto"></div>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-5">
        <div class="mb-3">
          <i class="fas fa-exclamation-circle fa-3x text-danger"></i>
        </div>
        <h5 class="text-muted mb-3">{{ t('app.error') }}</h5>
        <router-link to="/dashboard" class="btn btn-primary">
          <i class="fas fa-arrow-left me-2"></i>{{ t('app.back') }}
        </router-link>
      </div>

      <div v-else-if="cv">
        <!-- Toolbar -->
        <div class="preview-toolbar mb-4">
          <div class="toolbar-left">
            <router-link to="/dashboard" class="btn btn-toolbar btn-back">
              <i class="fas fa-arrow-left me-2"></i>{{ t('app.back') }}
            </router-link>
            <div class="toolbar-divider"></div>
            <div class="toolbar-title-group">
              <h4 class="toolbar-title mb-0">{{ cv.title }}</h4>
              <div class="toolbar-meta">
                <span class="meta-template">
                  <i class="fas fa-palette me-1"></i>{{ t(`templates.${cv.template}`) }}
                </span>
                <span class="meta-date">
                  <i class="fas fa-calendar-alt me-1"></i>{{ formatDate(cv.created_at) }}
                </span>
              </div>
            </div>
          </div>
          <div class="toolbar-center">
            <span class="status-pill" :class="`status-${cv.status}`">
              <i :class="statusIcon(cv.status)" class="me-1"></i>{{ t(`cv.${cv.status}`) }}
            </span>
          </div>
          <div class="toolbar-right">
            <button @click="handleShare" class="btn btn-toolbar" :class="{ 'btn-shared-active': cv.is_shared }" :title="t('cv.shareCV')">
              <i class="fas fa-share-alt me-1"></i><span class="d-none d-lg-inline">{{ t('app.share') }}</span>
            </button>
            <router-link :to="`/cv/${cv.id}/edit`" class="btn btn-toolbar btn-edit-cv">
              <i class="fas fa-edit me-1"></i><span class="d-none d-lg-inline">{{ t('app.edit') }}</span>
            </router-link>
            <button @click="exportPDF" class="btn btn-toolbar btn-export" :disabled="exporting">
              <span v-if="exporting" class="spinner-border spinner-border-sm me-1"></span>
              <i v-else class="fas fa-file-pdf me-1"></i><span class="d-none d-lg-inline">{{ t('cv.exportPDF') }}</span>
            </button>
            <div class="dropdown">
              <button class="btn btn-toolbar btn-ai-tools dropdown-toggle" data-bs-toggle="dropdown">
                <i class="fas fa-robot me-1"></i><span class="d-none d-lg-inline">{{ t('ai.title') }}</span>
              </button>
              <ul class="dropdown-menu dropdown-menu-end ai-dropdown">
                <li class="dropdown-header">
                  <i class="fas fa-magic me-1"></i>{{ t('ai.title') }}
                </li>
                <li><hr class="dropdown-divider"></li>
                <li>
                  <a class="dropdown-item" href="#" @click.prevent="aiAction('analyze')">
                    <i class="fas fa-chart-bar me-2 text-primary"></i>{{ t('ai.analyzeCV') }}
                  </a>
                </li>
                <li>
                  <a class="dropdown-item" href="#" @click.prevent="showCoverLetterModal = true">
                    <i class="fas fa-envelope me-2 text-success"></i>{{ t('ai.coverLetter') }}
                  </a>
                </li>
                <li>
                  <a class="dropdown-item" href="#" @click.prevent="aiAction('jobs')">
                    <i class="fas fa-briefcase me-2 text-warning"></i>{{ t('ai.suggestJobs') }}
                  </a>
                </li>
                <li>
                  <a class="dropdown-item" href="#" @click.prevent="aiAction('research')">
                    <i class="fas fa-microscope me-2 text-info"></i>{{ t('ai.evaluateResearch') }}
                  </a>
                </li>
                <li>
                  <a class="dropdown-item" href="#" @click.prevent="aiAction('tips')">
                    <i class="fas fa-lightbulb me-2 text-danger"></i>{{ t('ai.getTips') }}
                  </a>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Rejection Alert -->
        <div v-if="cv.status === 'rejected' && cv.reject_note" class="reject-alert mb-4">
          <div class="reject-alert-content">
            <div class="reject-icon">
              <i class="fas fa-exclamation-triangle"></i>
            </div>
            <div class="reject-text">
              <h6 class="reject-title mb-1">{{ t('preview.rejectionReason') }}</h6>
              <p class="reject-message mb-0">{{ cv.reject_note }}</p>
            </div>
            <router-link :to="`/cv/${cv.id}/edit`" class="btn btn-sm btn-outline-danger ms-auto">
              <i class="fas fa-edit me-1"></i>{{ t('preview.fixAndResubmit') }}
            </router-link>
          </div>
        </div>

        <!-- Share Section (visible when shared) -->
        <div v-if="cv.is_shared" class="share-section mb-4">
          <div class="share-section-content">
            <div class="share-info">
              <i class="fas fa-link text-primary me-2"></i>
              <span class="share-label">{{ t('preview.sharedPublicly') }}</span>
              <span class="share-views ms-3">
                <i class="fas fa-eye me-1"></i>{{ cv.view_count }} {{ t('cv.viewCount') }}
              </span>
            </div>
            <div class="share-actions d-flex align-items-center gap-2">
              <div class="input-group input-group-sm share-link-group">
                <input type="text" class="form-control share-url-input" :value="shareUrl" readonly />
                <button @click="copyShareLink" class="btn btn-primary btn-sm">
                  <i :class="copied ? 'fas fa-check' : 'fas fa-copy'" class="me-1"></i>
                  {{ copied ? t('dashboard.copied') : t('cv.copyLink') }}
                </button>
              </div>
              <button @click="handleShare" class="btn btn-outline-secondary btn-sm" :title="t('preview.manageSharing')">
                <i class="fas fa-cog"></i>
              </button>
            </div>
          </div>
        </div>

        <!-- CV Preview -->
        <div class="cv-preview-wrapper">
          <div class="cv-preview-container bg-white shadow mx-auto" ref="cvContainer">
            <CVTemplates :data="cv.data" :template="cv.template" :language="cv.language" />
          </div>
        </div>

        <!-- Share Modal -->
        <div class="modal fade" ref="shareModalRef" tabindex="-1">
          <div class="modal-dialog modal-dialog-centered">
            <div class="modal-content share-modal-content">
              <div class="modal-header border-0 pb-0">
                <h5 class="modal-title fw-bold">
                  <i class="fas fa-share-alt text-primary me-2"></i>{{ t('cv.shareCV') }}
                </h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
              </div>
              <div class="modal-body">
                <div class="share-toggle-area text-center mb-4">
                  <button
                    class="btn btn-lg px-4"
                    :class="cv.is_shared ? 'btn-danger' : 'btn-success'"
                    :disabled="sharingLoading"
                    @click="toggleShare"
                  >
                    <span v-if="sharingLoading" class="spinner-border spinner-border-sm me-2"></span>
                    <i v-else :class="cv.is_shared ? 'fas fa-lock' : 'fas fa-share-alt'" class="me-2"></i>
                    {{ cv.is_shared ? t('cv.disableShare') : t('cv.enableShare') }}
                  </button>
                </div>

                <div v-if="cv.is_shared">
                  <div class="mb-3">
                    <label class="form-label text-muted small text-uppercase fw-bold">{{ t('cv.shareLink') }}</label>
                    <div class="input-group">
                      <input type="text" class="form-control" :value="shareUrl" readonly />
                      <button @click="copyShareLink" class="btn btn-primary">
                        <i :class="copied ? 'fas fa-check' : 'fas fa-copy'" class="me-1"></i>
                        {{ copied ? t('dashboard.copied') : t('cv.copyLink') }}
                      </button>
                    </div>
                  </div>

                  <div class="text-center mb-3" v-if="qrCodeImage">
                    <label class="form-label text-muted small text-uppercase fw-bold d-block">{{ t('cv.qrCode') }}</label>
                    <div class="qr-code-box">
                      <img :src="qrCodeImage" alt="QR Code" class="img-fluid" />
                    </div>
                  </div>

                  <div class="text-center">
                    <span class="views-badge">
                      <i class="fas fa-eye me-1"></i>{{ cv.view_count }} {{ t('cv.viewCount') }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- AI Result Modal -->
        <div class="modal fade" :class="{ show: showAIResult }" :style="{ display: showAIResult ? 'block' : 'none' }" tabindex="-1">
          <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content ai-modal-content">
              <div class="modal-header border-0">
                <h5 class="modal-title fw-bold">
                  <i class="fas fa-robot text-primary me-2"></i>{{ t('ai.title') }}
                </h5>
                <button type="button" class="btn-close" @click="showAIResult = false"></button>
              </div>
              <div class="modal-body">
                <div v-if="aiLoading" class="ai-loading-state text-center py-5">
                  <div class="ai-spinner mb-3">
                    <i class="fas fa-robot fa-3x text-primary ai-pulse"></i>
                  </div>
                  <h6 class="text-muted">{{ t('ai.generating') }}</h6>
                  <div class="ai-loading-bar mt-3 mx-auto"></div>
                </div>
                <div v-else class="ai-result-content" v-html="formatAIResult(aiResult)"></div>
              </div>
              <div v-if="!aiLoading && aiResult" class="modal-footer border-0">
                <button class="btn btn-outline-secondary" @click="copyAIResult">
                  <i class="fas fa-copy me-1"></i>{{ t('cv.copyLink') }}
                </button>
                <button class="btn btn-primary" @click="showAIResult = false">{{ t('app.close') }}</button>
              </div>
            </div>
          </div>
        </div>
        <div v-if="showAIResult" class="modal-backdrop fade show"></div>

        <!-- Cover Letter Modal -->
        <div class="modal fade" :class="{ show: showCoverLetterModal }" :style="{ display: showCoverLetterModal ? 'block' : 'none' }" tabindex="-1">
          <div class="modal-dialog modal-dialog-centered">
            <div class="modal-content cover-letter-modal">
              <div class="modal-header border-0">
                <h5 class="modal-title fw-bold">
                  <i class="fas fa-envelope text-success me-2"></i>{{ t('ai.coverLetter') }}
                </h5>
                <button type="button" class="btn-close" @click="showCoverLetterModal = false"></button>
              </div>
              <div class="modal-body">
                <div class="mb-3">
                  <label class="form-label fw-medium">{{ t('ai.jobTitle') }}</label>
                  <div class="input-group">
                    <span class="input-group-text"><i class="fas fa-briefcase text-muted"></i></span>
                    <input type="text" class="form-control" v-model="coverLetterForm.job_title" :placeholder="t('ai.jobTitle')" />
                  </div>
                </div>
                <div class="mb-3">
                  <label class="form-label fw-medium">{{ t('ai.companyName') }}</label>
                  <div class="input-group">
                    <span class="input-group-text"><i class="fas fa-building text-muted"></i></span>
                    <input type="text" class="form-control" v-model="coverLetterForm.company" :placeholder="t('ai.companyName')" />
                  </div>
                </div>
              </div>
              <div class="modal-footer border-0">
                <button class="btn btn-outline-secondary" @click="showCoverLetterModal = false">{{ t('app.cancel') }}</button>
                <button class="btn btn-success" @click="generateCoverLetter">
                  <i class="fas fa-magic me-1"></i>{{ t('preview.generate') }}
                </button>
              </div>
            </div>
          </div>
        </div>
        <div v-if="showCoverLetterModal" class="modal-backdrop fade show"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useCVStore } from '../stores/cv'
import { aiAPI } from '../services/api'
import CVTemplates from '../components/CVTemplates.vue'
import QRCode from 'qrcode'

const { t } = useI18n()
const route = useRoute()
const cvStore = useCVStore()

const loading = ref(true)
const error = ref(false)
const exporting = ref(false)
const cv = ref<any>(null)
const cvContainer = ref<HTMLElement>()
const shareModalRef = ref<HTMLElement>()
const showAIResult = ref(false)
const aiLoading = ref(false)
const aiResult = ref('')
const showCoverLetterModal = ref(false)
const coverLetterForm = reactive({ job_title: '', company: '' })
const copied = ref(false)
const qrCodeImage = ref('')
const sharingLoading = ref(false)
let bsModal: any = null

const shareUrl = computed(() => {
  if (!cv.value) return ''
  return `${window.location.origin}/shared/${cv.value.share_token}`
})

// Generate QR code whenever share URL changes
watch(shareUrl, async (url) => {
  if (url) {
    try {
      qrCodeImage.value = await QRCode.toDataURL(url, { width: 200, margin: 2 })
    } catch {
      qrCodeImage.value = ''
    }
  }
}, { immediate: true })

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    cv.value = await cvStore.fetchCV(id)
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
})

function statusIcon(status: string) {
  const map: Record<string, string> = {
    draft: 'fas fa-pencil-alt',
    pending: 'fas fa-clock',
    approved: 'fas fa-check-circle',
    rejected: 'fas fa-times-circle',
  }
  return map[status] || 'fas fa-circle'
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}

async function handleShare() {
  copied.value = false
  if (shareModalRef.value) {
    const { Modal } = await import('bootstrap')
    bsModal = new Modal(shareModalRef.value)
    bsModal.show()
  }
}

function copyShareLink() {
  navigator.clipboard.writeText(shareUrl.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

async function toggleShare() {
  if (!cv.value || sharingLoading.value) return
  sharingLoading.value = true
  try {
    const updated = await cvStore.toggleShare(cv.value.id)
    cv.value = { ...cv.value, is_shared: updated.is_shared, share_token: updated.share_token }
  } catch {
    // Revert checkbox visual state on error
  } finally {
    sharingLoading.value = false
  }
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

async function aiAction(type: string) {
  showAIResult.value = true
  aiLoading.value = true
  aiResult.value = ''
  try {
    let res
    const id = cv.value.id
    switch (type) {
      case 'analyze': res = await aiAPI.analyzeCV(id); break
      case 'jobs': res = await aiAPI.suggestJobs(id); break
      case 'research': res = await aiAPI.evaluateResearch(id); break
      case 'tips': res = await aiAPI.getTips(id); break
    }
    aiResult.value = res?.data?.result || ''
  } catch (err: any) {
    aiResult.value = err.response?.data?.error || t('app.error')
  } finally {
    aiLoading.value = false
  }
}

async function generateCoverLetter() {
  showCoverLetterModal.value = false
  showAIResult.value = true
  aiLoading.value = true
  aiResult.value = ''
  try {
    const res = await aiAPI.coverLetter({
      cv_id: cv.value.id,
      job_title: coverLetterForm.job_title,
      company: coverLetterForm.company,
      language: cv.value.language,
    })
    aiResult.value = res.data.result || ''
  } catch (err: any) {
    aiResult.value = err.response?.data?.error || t('app.error')
  } finally {
    aiLoading.value = false
  }
}

function formatAIResult(text: string) {
  return text.replace(/\n/g, '<br>')
}

function copyAIResult() {
  navigator.clipboard.writeText(aiResult.value)
}
</script>

<style scoped>
/* === Loading State === */
.loading-state {
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-bar {
  width: 200px;
  height: 4px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.loading-bar::after {
  content: '';
  position: absolute;
  top: 0;
  left: -50%;
  width: 50%;
  height: 100%;
  background: #6366f1;
  border-radius: 4px;
  animation: loading-slide 1.2s ease-in-out infinite;
}

@keyframes loading-slide {
  0% { left: -50%; }
  100% { left: 100%; }
}

/* === Toolbar === */
.preview-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  background: white;
  border-radius: 14px;
  padding: 12px 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid #e8ecf3;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  flex: 1;
}

.toolbar-center {
  flex-shrink: 0;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.toolbar-divider {
  width: 1px;
  height: 32px;
  background: #e2e8f0;
  flex-shrink: 0;
}

.toolbar-title {
  font-size: 1.15rem;
  font-weight: 700;
  color: #1e293b;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.toolbar-meta {
  display: flex;
  gap: 12px;
  font-size: 0.78rem;
  color: #94a3b8;
  margin-top: 2px;
}

.btn-toolbar {
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #475569;
  font-weight: 500;
  font-size: 0.88rem;
  border-radius: 10px;
  padding: 8px 14px;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-toolbar:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
  color: #1e293b;
}

.btn-back:hover {
  border-color: #6366f1;
  color: #6366f1;
}

.btn-edit-cv {
  border-color: #6366f1;
  color: #6366f1;
  background: #f0f0ff;
}

.btn-edit-cv:hover {
  background: #6366f1;
  color: white;
}

.btn-export {
  border-color: #dc2626;
  color: #dc2626;
  background: #fef2f2;
}

.btn-export:hover {
  background: #dc2626;
  color: white;
}

.btn-ai-tools {
  border-color: #6366f1;
  color: #6366f1;
  background: #f0f0ff;
}

.btn-ai-tools:hover {
  background: #6366f1;
  color: white;
}

.btn-shared-active {
  border-color: #0ea5e9;
  color: #0ea5e9;
  background: #f0f9ff;
}

/* Status Pill */
.status-pill {
  display: inline-flex;
  align-items: center;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.82rem;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.status-draft {
  background: #f1f5f9;
  color: #64748b;
}

.status-pending {
  background: #fef3c7;
  color: #92400e;
}

.status-approved {
  background: #dcfce7;
  color: #166534;
}

.status-rejected {
  background: #fee2e2;
  color: #991b1b;
}

/* === Rejection Alert === */
.reject-alert {
  background: linear-gradient(135deg, #fef2f2, #fee2e2);
  border: 1px solid #fca5a5;
  border-radius: 14px;
  overflow: hidden;
}

.reject-alert-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
}

.reject-icon {
  flex-shrink: 0;
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: #fee2e2;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #dc2626;
  font-size: 1.2rem;
}

.reject-title {
  color: #991b1b;
  font-weight: 700;
  font-size: 0.9rem;
}

.reject-message {
  color: #7f1d1d;
  font-size: 0.88rem;
  line-height: 1.5;
}

/* === Share Section === */
.share-section {
  background: linear-gradient(135deg, #f0f9ff, #e0f2fe);
  border: 1px solid #bae6fd;
  border-radius: 14px;
  padding: 14px 20px;
}

.share-section-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.share-info {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  color: #0369a1;
}

.share-label {
  font-weight: 600;
}

.share-views {
  color: #64748b;
  font-size: 0.82rem;
}

.share-link-group {
  max-width: 380px;
}

.share-url-input {
  background: white;
  font-size: 0.82rem;
}

/* === CV Preview === */
.cv-preview-wrapper {
  padding: 20px 0;
}

.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
  border-radius: 4px;
  transition: box-shadow 0.3s;
}

.cv-preview-container:hover {
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.12) !important;
}

/* === Share Modal === */
.share-modal-content {
  border: none;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.share-toggle {
  width: 3em;
  height: 1.5em;
  cursor: pointer;
}

.qr-code-box {
  display: inline-block;
  padding: 12px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e8ecf3;
}

.qr-code-box img {
  max-width: 180px;
}

.views-badge {
  display: inline-flex;
  align-items: center;
  background: #f0f9ff;
  color: #0284c7;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.88rem;
  font-weight: 500;
}

/* === AI Modal === */
.ai-modal-content {
  border: none;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.ai-dropdown {
  border: none;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.12);
  padding: 8px;
}

.ai-dropdown .dropdown-item {
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 0.9rem;
  transition: background 0.15s;
}

.ai-dropdown .dropdown-item:hover {
  background: #f0f0ff;
}

.ai-dropdown .dropdown-header {
  color: #6366f1;
  font-weight: 700;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.ai-pulse {
  animation: ai-pulse 2s ease-in-out infinite;
}

@keyframes ai-pulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.05); }
}

.ai-loading-bar {
  width: 180px;
  height: 4px;
  background: #e2e8f0;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.ai-loading-bar::after {
  content: '';
  position: absolute;
  top: 0;
  left: -50%;
  width: 50%;
  height: 100%;
  background: #6366f1;
  border-radius: 4px;
  animation: loading-slide 1.2s ease-in-out infinite;
}

.ai-result-content {
  white-space: pre-wrap;
  line-height: 1.8;
  font-size: 0.95rem;
  color: #334155;
}

/* === Cover Letter Modal === */
.cover-letter-modal {
  border: none;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.cover-letter-modal .input-group-text {
  background: #f8fafc;
  border-color: #e2e8f0;
}

/* === Responsive === */
@media (max-width: 992px) {
  .preview-toolbar {
    padding: 10px 14px;
  }

  .toolbar-left {
    flex-basis: 100%;
    order: 1;
  }

  .toolbar-center {
    order: 2;
  }

  .toolbar-right {
    order: 3;
    flex-wrap: wrap;
    justify-content: flex-end;
  }

  .toolbar-divider {
    display: none;
  }
}

@media (max-width: 768px) {
  .reject-alert-content {
    flex-direction: column;
    text-align: center;
  }

  .share-section-content {
    flex-direction: column;
    text-align: center;
  }

  .share-link-group {
    max-width: 100%;
    width: 100%;
  }

  .btn-toolbar {
    padding: 6px 10px;
    font-size: 0.82rem;
  }
}
</style>
