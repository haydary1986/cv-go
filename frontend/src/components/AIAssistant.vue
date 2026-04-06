<template>
  <div v-if="visible" class="ai-assistant-panel card shadow-sm">
    <div class="card-header d-flex justify-content-between align-items-center py-2 px-3">
      <small class="fw-bold"><i class="fas fa-robot me-1"></i> {{ locale === 'ar' ? 'مساعد AI' : 'AI Assistant' }}</small>
      <button class="btn btn-sm btn-link text-muted p-0" @click="visible = false"><i class="fas fa-times"></i></button>
    </div>
    <div class="card-body p-3">
      <div v-if="loading" class="text-center py-2">
        <div class="spinner-border spinner-border-sm text-primary"></div>
        <small class="ms-2 text-muted">{{ locale === 'ar' ? 'جاري التحسين...' : 'Improving...' }}</small>
      </div>
      <div v-else-if="suggestion">
        <p class="small mb-2" style="white-space: pre-wrap;">{{ suggestion }}</p>
        <div class="d-flex gap-2">
          <button class="btn btn-sm btn-success" @click="accept">
            <i class="fas fa-check me-1"></i>{{ locale === 'ar' ? 'قبول' : 'Accept' }}
          </button>
          <button class="btn btn-sm btn-outline-secondary" @click="reject">
            <i class="fas fa-redo me-1"></i>{{ locale === 'ar' ? 'اقتراح آخر' : 'Try again' }}
          </button>
        </div>
      </div>
      <div v-else class="text-center">
        <button class="btn btn-sm btn-primary me-2" @click="improve" :disabled="!currentText">
          <i class="fas fa-magic me-1"></i>{{ locale === 'ar' ? 'تحسين النص' : 'Improve' }}
        </button>
        <button class="btn btn-sm btn-outline-primary" @click="generate" :disabled="!fieldName">
          <i class="fas fa-pen me-1"></i>{{ locale === 'ar' ? 'توليد' : 'Generate' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { aiAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const { locale } = useI18n()
const toast = useToast()

const props = defineProps<{
  currentText: string
  fieldName: string
  language: string
  cvContext?: any
}>()

const emit = defineEmits<{
  (e: 'accept', text: string): void
}>()

const visible = ref(false)
const loading = ref(false)
const suggestion = ref('')

function show() {
  visible.value = true
  suggestion.value = ''
}

async function improve() {
  if (!props.currentText) return
  loading.value = true
  try {
    const res = await aiAPI.improveText({ text: props.currentText, language: props.language })
    suggestion.value = res.data.result
  } catch {
    toast.error(locale.value === 'ar' ? 'فشل تحسين النص' : 'Failed to improve text')
  } finally {
    loading.value = false
  }
}

async function generate() {
  loading.value = true
  try {
    const res = await aiAPI.improveText({
      text: `Generate a professional ${props.fieldName} section for a CV. Language: ${props.language}. Context: ${JSON.stringify(props.cvContext || {})}`,
      language: props.language,
    })
    suggestion.value = res.data.result
  } catch {
    toast.error(locale.value === 'ar' ? 'فشل التوليد' : 'Failed to generate')
  } finally {
    loading.value = false
  }
}

function accept() {
  emit('accept', suggestion.value)
  suggestion.value = ''
  visible.value = false
}

function reject() {
  suggestion.value = ''
  improve()
}

defineExpose({ show })
</script>

<style scoped>
.ai-assistant-panel {
  position: absolute;
  z-index: 1050;
  width: 350px;
  max-width: 90vw;
  right: 0;
  top: 100%;
  margin-top: 4px;
}
</style>
