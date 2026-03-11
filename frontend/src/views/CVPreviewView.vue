<template>
  <div class="container-fluid py-4">
    <div v-if="loading" class="text-center py-5">
      <div class="spinner-border text-primary"></div>
    </div>

    <div v-else-if="cv">
      <!-- Toolbar -->
      <div class="d-flex justify-content-between align-items-center mb-3 flex-wrap gap-2">
        <div class="d-flex gap-2 align-items-center">
          <router-link to="/dashboard" class="btn btn-outline-secondary">
            <i class="fas fa-arrow-left me-1"></i>{{ t('app.back') }}
          </router-link>
          <h4 class="mb-0">{{ cv.title }}</h4>
          <span :class="statusBadge(cv.status)" class="badge">{{ t(`cv.${cv.status}`) }}</span>
        </div>
        <div class="d-flex gap-2">
          <button @click="exportPDF" class="btn btn-danger" :disabled="exporting">
            <span v-if="exporting" class="spinner-border spinner-border-sm me-1"></span>
            <i v-else class="fas fa-file-pdf me-1"></i>{{ t('cv.exportPDF') }}
          </button>
          <router-link :to="`/cv/${cv.id}/edit`" class="btn btn-primary">
            <i class="fas fa-edit me-1"></i>{{ t('app.edit') }}
          </router-link>
          <div class="dropdown">
            <button class="btn btn-outline-primary dropdown-toggle" data-bs-toggle="dropdown">
              <i class="fas fa-robot me-1"></i>{{ t('ai.title') }}
            </button>
            <ul class="dropdown-menu">
              <li><a class="dropdown-item" href="#" @click.prevent="aiAction('analyze')">{{ t('ai.analyzeCV') }}</a></li>
              <li><a class="dropdown-item" href="#" @click.prevent="showCoverLetterModal = true">{{ t('ai.coverLetter') }}</a></li>
              <li><a class="dropdown-item" href="#" @click.prevent="aiAction('jobs')">{{ t('ai.suggestJobs') }}</a></li>
              <li><a class="dropdown-item" href="#" @click.prevent="aiAction('research')">{{ t('ai.evaluateResearch') }}</a></li>
              <li><a class="dropdown-item" href="#" @click.prevent="aiAction('tips')">{{ t('ai.getTips') }}</a></li>
            </ul>
          </div>
        </div>
      </div>

      <!-- CV Preview -->
      <div class="cv-preview-container bg-white shadow mx-auto" ref="cvContainer">
        <CVTemplates :data="cv.data" :template="cv.template" :language="cv.language" />
      </div>

      <!-- AI Result Modal -->
      <div class="modal fade" :class="{ show: showAIResult }" :style="{ display: showAIResult ? 'block' : 'none' }" tabindex="-1">
        <div class="modal-dialog modal-lg">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">{{ t('ai.title') }}</h5>
              <button type="button" class="btn-close" @click="showAIResult = false"></button>
            </div>
            <div class="modal-body">
              <div v-if="aiLoading" class="text-center py-4">
                <div class="spinner-border text-primary"></div>
                <p class="mt-2">{{ t('ai.generating') }}</p>
              </div>
              <div v-else class="ai-result" v-html="formatAIResult(aiResult)"></div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="showAIResult" class="modal-backdrop fade show"></div>

      <!-- Cover Letter Modal -->
      <div class="modal fade" :class="{ show: showCoverLetterModal }" :style="{ display: showCoverLetterModal ? 'block' : 'none' }" tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">{{ t('ai.coverLetter') }}</h5>
              <button type="button" class="btn-close" @click="showCoverLetterModal = false"></button>
            </div>
            <div class="modal-body">
              <div class="mb-3">
                <label class="form-label">{{ t('ai.jobTitle') }}</label>
                <input type="text" class="form-control" v-model="coverLetterForm.job_title" />
              </div>
              <div class="mb-3">
                <label class="form-label">{{ t('ai.companyName') }}</label>
                <input type="text" class="form-control" v-model="coverLetterForm.company" />
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="showCoverLetterModal = false">{{ t('app.cancel') }}</button>
              <button class="btn btn-primary" @click="generateCoverLetter">{{ t('app.submit') }}</button>
            </div>
          </div>
        </div>
      </div>
      <div v-if="showCoverLetterModal" class="modal-backdrop fade show"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useCVStore } from '../stores/cv'
import { aiAPI } from '../services/api'
import CVTemplates from '../components/CVTemplates.vue'

const { t } = useI18n()
const route = useRoute()
const cvStore = useCVStore()

const loading = ref(true)
const exporting = ref(false)
const cv = ref<any>(null)
const cvContainer = ref<HTMLElement>()
const showAIResult = ref(false)
const aiLoading = ref(false)
const aiResult = ref('')
const showCoverLetterModal = ref(false)
const coverLetterForm = reactive({ job_title: '', company: '' })

onMounted(async () => {
  const id = Number(route.params.id)
  cv.value = await cvStore.fetchCV(id)
  loading.value = false
})

function statusBadge(status: string) {
  const map: Record<string, string> = {
    draft: 'bg-secondary', pending: 'bg-warning text-dark',
    approved: 'bg-success', rejected: 'bg-danger',
  }
  return map[status] || 'bg-secondary'
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
</script>

<style scoped>
.cv-preview-container {
  max-width: 210mm;
  min-height: 297mm;
}
.ai-result {
  white-space: pre-wrap;
  line-height: 1.8;
}
</style>
