<template>
  <div class="linkedin-page">
    <div class="container py-5" style="max-width: 720px;">

      <!-- Page Header -->
      <div class="text-center mb-4">
        <div class="linkedin-icon-wrap mb-3">
          <i class="fab fa-linkedin-in"></i>
        </div>
        <h3 class="page-title">
          {{ locale === 'ar' ? 'استيراد من LinkedIn' : 'Import from LinkedIn' }}
        </h3>
      </div>

      <!-- Step-by-Step Instructions Card -->
      <div class="instructions-card mb-4">
        <div class="instructions-header">
          <i class="fas fa-list-ol me-2"></i>
          {{ locale === 'ar' ? 'كيفية تصدير بياناتك من LinkedIn:' : 'How to export your LinkedIn data:' }}
        </div>
        <div class="instructions-body">
          <div class="instruction-step" v-for="(step, idx) in [
            locale === 'ar' ? 'افتح LinkedIn وانتقل إلى الإعدادات' : 'Open LinkedIn and go to Settings',
            locale === 'ar' ? 'انقر على &quot;الخصوصية&quot; ثم &quot;الحصول على نسخة من بياناتك&quot;' : 'Click &quot;Data Privacy&quot; then &quot;Get a copy of your data&quot;',
            locale === 'ar' ? 'اختر &quot;التحميل الكامل&quot; وانتظر البريد' : 'Select &quot;Download larger data archive&quot; and wait for email',
            locale === 'ar' ? 'حمّل الملف وارفعه هنا' : 'Download the file and upload it here'
          ]" :key="idx">
            <div class="step-number">{{ idx + 1 }}</div>
            <span v-html="step"></span>
          </div>
        </div>
      </div>

      <!-- Upload Drop Zone -->
      <div
        class="upload-zone"
        :class="{ 'upload-zone--active': isDragging }"
        @dragover.prevent="isDragging = true"
        @dragleave="isDragging = false"
        @drop.prevent="handleDrop"
      >
        <div class="upload-zone-inner">
          <div class="upload-icon">
            <i class="fas fa-cloud-upload-alt"></i>
          </div>
          <p class="upload-title mb-1">
            {{ locale === 'ar' ? 'اسحب ملف LinkedIn ZIP هنا' : 'Drop your LinkedIn ZIP file here' }}
          </p>
          <p class="upload-subtitle mb-3">
            {{ locale === 'ar' ? 'أو' : 'or' }}
          </p>
          <label class="upload-browse-btn">
            <i class="fas fa-folder-open me-1"></i>
            {{ locale === 'ar' ? 'اختر ملف' : 'Browse Files' }}
            <input type="file" accept=".zip" class="d-none" @change="handleFileSelect" />
          </label>
          <p class="upload-hint mt-2">ZIP {{ locale === 'ar' ? 'فقط — الحد الأقصى 10 ميجابايت' : 'only — Max 10MB' }}</p>
        </div>
      </div>

      <!-- Loading / Progress -->
      <div v-if="loading" class="loading-card text-center my-4">
        <div class="progress-track mb-3">
          <div class="progress-track-fill" style="width: 70%;"></div>
        </div>
        <p class="loading-text mb-0">
          <i class="fas fa-cog fa-spin me-1"></i>
          {{ locale === 'ar' ? 'جاري تحليل البيانات...' : 'Analyzing data...' }}
        </p>
      </div>

      <!-- Data Preview -->
      <div v-if="preview" class="preview-section mt-4">
        <h5 class="preview-title mb-3">
          <i class="fas fa-check-circle me-2" style="color: #1e8449;"></i>
          {{ locale === 'ar' ? 'معاينة البيانات المستوردة:' : 'Imported Data Preview:' }}
        </h5>
        <div class="row g-3 mb-4">
          <div class="col-6 col-md-4">
            <div class="preview-stat-card">
              <i class="fas fa-user"></i>
              <div class="preview-stat-label">{{ locale === 'ar' ? 'الاسم' : 'Name' }}</div>
              <div class="preview-stat-value">{{ preview.personal_info?.full_name || '—' }}</div>
            </div>
          </div>
          <div class="col-6 col-md-4">
            <div class="preview-stat-card">
              <i class="fas fa-envelope"></i>
              <div class="preview-stat-label">{{ locale === 'ar' ? 'البريد' : 'Email' }}</div>
              <div class="preview-stat-value">{{ preview.personal_info?.email || '—' }}</div>
            </div>
          </div>
          <div class="col-4 col-md-4">
            <div class="preview-stat-card">
              <i class="fas fa-briefcase"></i>
              <div class="preview-stat-label">{{ locale === 'ar' ? 'الخبرات' : 'Experiences' }}</div>
              <div class="preview-stat-value">{{ preview.experiences?.length || 0 }}</div>
            </div>
          </div>
          <div class="col-4 col-md-6">
            <div class="preview-stat-card">
              <i class="fas fa-graduation-cap"></i>
              <div class="preview-stat-label">{{ locale === 'ar' ? 'التعليم' : 'Education' }}</div>
              <div class="preview-stat-value">{{ preview.education?.length || 0 }}</div>
            </div>
          </div>
          <div class="col-4 col-md-6">
            <div class="preview-stat-card">
              <i class="fas fa-tools"></i>
              <div class="preview-stat-label">{{ locale === 'ar' ? 'المهارات' : 'Skills' }}</div>
              <div class="preview-stat-value">{{ preview.skills?.length || 0 }}</div>
            </div>
          </div>
        </div>
        <button class="create-cv-btn" @click="createCV">
          <i class="fas fa-plus me-2"></i>
          {{ locale === 'ar' ? 'إنشاء سيرة ذاتية من هذه البيانات' : 'Create CV from this data' }}
        </button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { cvAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const { locale } = useI18n()
const router = useRouter()
const toast = useToast()

const isDragging = ref(false)
const loading = ref(false)
const preview = ref<any>(null)

function handleDrop(e: DragEvent) {
  isDragging.value = false
  const file = e.dataTransfer?.files[0]
  if (file) processFile(file)
}

function handleFileSelect(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) processFile(file)
}

async function processFile(file: File) {
  if (!file.name.endsWith('.zip')) {
    toast.error(locale.value === 'ar' ? 'يرجى رفع ملف ZIP' : 'Please upload a ZIP file')
    return
  }
  if (file.size > 10 * 1024 * 1024) {
    toast.error(locale.value === 'ar' ? 'الملف كبير جداً (الحد 10MB)' : 'File too large (max 10MB)')
    return
  }

  loading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const res = await cvAPI.importLinkedIn(formData)
    preview.value = res.data.data
    toast.success(locale.value === 'ar' ? 'تم تحليل البيانات بنجاح' : 'Data analyzed successfully')
  } catch {
    toast.error(locale.value === 'ar' ? 'فشل تحليل الملف' : 'Failed to parse file')
  } finally {
    loading.value = false
  }
}

async function createCV() {
  if (!preview.value) return
  // Store preview data and navigate to CV creation form
  localStorage.setItem('cv-import-data', JSON.stringify(preview.value))
  router.push('/cv/create')
}
</script>

<style scoped>
.linkedin-page {
  min-height: 80vh;
  background: #ffffff;
}

.linkedin-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: rgba(0, 119, 181, 0.08);
  color: #0077b5;
  font-size: 1.5rem;
}

.page-title {
  font-weight: 700;
  color: #222222;
}

/* Instructions Card */
.instructions-card {
  border-radius: 12px;
  overflow: hidden;
  background: #ffffff;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.instructions-header {
  background: #222222;
  color: #ffffff;
  padding: 0.75rem 1.25rem;
  font-weight: 600;
  font-size: 0.9rem;
}

.instructions-body {
  padding: 1rem 1.25rem;
}

.instruction-step {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 0;
  font-size: 0.9rem;
  color: #222222;
}

.instruction-step + .instruction-step {
  border-top: 1px solid #ebebeb;
}

.step-number {
  width: 28px;
  height: 28px;
  min-width: 28px;
  border-radius: 50%;
  background: #1a5276;
  color: #ffffff;
  font-weight: 700;
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Upload Zone */
.upload-zone {
  border: 2px dashed #c1c1c1;
  border-radius: 12px;
  padding: 2.5rem 1.5rem;
  text-align: center;
  background: #fafafa;
  transition: all 0.25s ease;
  cursor: pointer;
}

.upload-zone--active {
  border-color: #1a5276;
  background: rgba(26, 82, 118, 0.03);
}

.upload-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #f0f0f0;
  color: #6a6a6a;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.upload-zone--active .upload-icon {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.upload-title {
  font-weight: 600;
  color: #222222;
  font-size: 1rem;
}

.upload-subtitle {
  color: #6a6a6a;
  font-size: 0.85rem;
}

.upload-browse-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.5rem 1.5rem;
  border-radius: 8px;
  background: #1a5276;
  color: #ffffff;
  font-weight: 600;
  font-size: 0.88rem;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.upload-browse-btn:hover {
  background: #154360;
}

.upload-hint {
  font-size: 0.75rem;
  color: #6a6a6a;
}

/* Loading */
.loading-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px;
}

.progress-track {
  height: 6px;
  background: #f0f0f0;
  border-radius: 6px;
  overflow: hidden;
}

.progress-track-fill {
  height: 100%;
  background: #1a5276;
  border-radius: 6px;
  animation: progress-pulse 1.5s ease-in-out infinite;
}

@keyframes progress-pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.loading-text {
  color: #6a6a6a;
  font-size: 0.9rem;
}

/* Preview */
.preview-title {
  color: #222222;
  font-weight: 700;
}

.preview-stat-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 1rem;
  text-align: center;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px;
  transition: box-shadow 0.2s;
}

.preview-stat-card:hover {
  box-shadow: rgba(0, 0, 0, 0.08) 0px 4px 12px;
}

.preview-stat-card i {
  font-size: 1.25rem;
  color: #1a5276;
  margin-bottom: 0.5rem;
  display: block;
}

.preview-stat-label {
  font-size: 0.72rem;
  color: #6a6a6a;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-weight: 600;
  margin-bottom: 0.2rem;
}

.preview-stat-value {
  font-weight: 700;
  color: #222222;
  font-size: 0.9rem;
  word-break: break-all;
}

.create-cv-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.65rem 2rem;
  border-radius: 8px;
  background: #1a5276;
  color: #ffffff;
  font-weight: 600;
  font-size: 0.95rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.create-cv-btn:hover {
  background: #154360;
}
</style>
