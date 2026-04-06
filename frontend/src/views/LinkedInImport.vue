<template>
  <div class="container py-4">
    <div class="row justify-content-center">
      <div class="col-md-8">
        <div class="card">
          <div class="card-header">
            <h4 class="mb-0">
              <i class="fab fa-linkedin me-2"></i>
              {{ locale === 'ar' ? 'استيراد من LinkedIn' : 'Import from LinkedIn' }}
            </h4>
          </div>
          <div class="card-body">
            <!-- Instructions -->
            <div class="alert alert-info">
              <h6>{{ locale === 'ar' ? 'كيفية تصدير بياناتك من LinkedIn:' : 'How to export your LinkedIn data:' }}</h6>
              <ol class="mb-0">
                <li>{{ locale === 'ar' ? 'افتح LinkedIn وانتقل إلى الإعدادات' : 'Open LinkedIn and go to Settings' }}</li>
                <li>{{ locale === 'ar' ? 'انقر على "الخصوصية" ثم "الحصول على نسخة من بياناتك"' : 'Click "Data Privacy" then "Get a copy of your data"' }}</li>
                <li>{{ locale === 'ar' ? 'اختر "التحميل الكامل" وانتظر البريد' : 'Select "Download larger data archive" and wait for email' }}</li>
                <li>{{ locale === 'ar' ? 'حمّل الملف وارفعه هنا' : 'Download the file and upload it here' }}</li>
              </ol>
            </div>

            <!-- Upload Zone -->
            <div
              class="border border-2 border-dashed rounded p-5 text-center mb-3"
              :class="{ 'border-primary bg-light': isDragging }"
              @dragover.prevent="isDragging = true"
              @dragleave="isDragging = false"
              @drop.prevent="handleDrop"
            >
              <i class="fas fa-cloud-upload-alt fa-3x mb-3 text-muted"></i>
              <p class="mb-2">{{ locale === 'ar' ? 'اسحب ملف LinkedIn ZIP هنا' : 'Drag LinkedIn ZIP file here' }}</p>
              <p class="text-muted small mb-3">{{ locale === 'ar' ? 'أو' : 'or' }}</p>
              <label class="btn btn-outline-primary">
                {{ locale === 'ar' ? 'اختر ملف' : 'Choose File' }}
                <input type="file" accept=".zip" class="d-none" @change="handleFileSelect" />
              </label>
            </div>

            <div v-if="loading" class="text-center py-3">
              <div class="spinner-border text-primary"></div>
              <p class="mt-2">{{ locale === 'ar' ? 'جاري تحليل البيانات...' : 'Analyzing data...' }}</p>
            </div>

            <!-- Preview -->
            <div v-if="preview" class="mt-4">
              <h5>{{ locale === 'ar' ? 'معاينة البيانات المستوردة:' : 'Imported Data Preview:' }}</h5>
              <div class="card mb-2">
                <div class="card-body">
                  <p><strong>{{ locale === 'ar' ? 'الاسم:' : 'Name:' }}</strong> {{ preview.personal_info?.full_name }}</p>
                  <p><strong>{{ locale === 'ar' ? 'البريد:' : 'Email:' }}</strong> {{ preview.personal_info?.email }}</p>
                  <p><strong>{{ locale === 'ar' ? 'الخبرات:' : 'Experiences:' }}</strong> {{ preview.experiences?.length || 0 }}</p>
                  <p><strong>{{ locale === 'ar' ? 'التعليم:' : 'Education:' }}</strong> {{ preview.education?.length || 0 }}</p>
                  <p><strong>{{ locale === 'ar' ? 'المهارات:' : 'Skills:' }}</strong> {{ preview.skills?.length || 0 }}</p>
                </div>
              </div>
              <button class="btn btn-primary" @click="createCV">
                <i class="fas fa-plus me-1"></i>
                {{ locale === 'ar' ? 'إنشاء سيرة ذاتية من هذه البيانات' : 'Create CV from this data' }}
              </button>
            </div>
          </div>
        </div>
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
