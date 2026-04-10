<template>
  <div v-if="visible" class="ai-panel">
    <!-- Header -->
    <div class="ai-panel-header">
      <div class="ai-panel-title">
        <span class="ai-sparkle">&#10024;</span>
        {{ locale === 'ar' ? 'مساعد AI' : 'AI Assistant' }}
      </div>
      <button class="ai-panel-close" @click="visible = false">
        <i class="fas fa-times"></i>
      </button>
    </div>

    <!-- Body -->
    <div class="ai-panel-body">
      <!-- Loading -->
      <div v-if="loading" class="ai-loading text-center py-3">
        <div class="ai-loading-dots">
          <span></span><span></span><span></span>
        </div>
        <small class="cell-muted mt-2 d-block">{{ locale === 'ar' ? 'جاري التحسين...' : 'Improving...' }}</small>
      </div>

      <!-- Suggestion Card -->
      <div v-else-if="suggestion" class="ai-suggestion">
        <div class="ai-suggestion-card">
          <p class="ai-suggestion-text">{{ suggestion }}</p>
        </div>
        <div class="d-flex gap-2 mt-3">
          <button class="ai-btn ai-btn--accept" @click="accept">
            <i class="fas fa-check me-1"></i>{{ locale === 'ar' ? 'قبول' : 'Accept' }}
          </button>
          <button class="ai-btn ai-btn--reject" @click="reject">
            <i class="fas fa-redo me-1"></i>{{ locale === 'ar' ? 'اقتراح آخر' : 'Try again' }}
          </button>
        </div>
      </div>

      <!-- Actions -->
      <div v-else class="ai-actions text-center py-2">
        <button class="ai-btn ai-btn--primary me-2" @click="improve" :disabled="!currentText">
          <i class="fas fa-magic me-1"></i>{{ locale === 'ar' ? 'تحسين النص' : 'Improve' }}
        </button>
        <button class="ai-btn ai-btn--outline" @click="generate" :disabled="!fieldName">
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
.ai-panel {
  position: absolute;
  z-index: 1050;
  width: 370px;
  max-width: 90vw;
  right: 0;
  top: 100%;
  margin-top: 6px;
  border-radius: 12px;
  overflow: hidden;
  background: #ffffff;
  border: none;
  box-shadow: rgba(0,0,0,0.02) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 6px, rgba(0,0,0,0.1) 0px 4px 8px;
}

.ai-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #ebebeb;
  background: #ffffff;
}

.ai-panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #222222;
  display: flex;
  align-items: center;
  gap: 6px;
}

.ai-sparkle {
  font-size: 1.1rem;
  animation: sparkle-pulse 2s ease-in-out infinite;
}

@keyframes sparkle-pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(1.15); }
}

.ai-panel-close {
  background: #f7f7f7;
  border: none;
  color: #222222;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  cursor: pointer;
  transition: background 0.15s;
}

.ai-panel-close:hover {
  background: #ebebeb;
}

.ai-panel-body {
  padding: 16px;
}

.cell-muted {
  color: #6a6a6a;
}

/* Loading dots */
.ai-loading-dots {
  display: inline-flex;
  gap: 4px;
}

.ai-loading-dots span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #c0982b;
  animation: dot-bounce 1.4s infinite ease-in-out;
}

.ai-loading-dots span:nth-child(2) {
  animation-delay: 0.16s;
}

.ai-loading-dots span:nth-child(3) {
  animation-delay: 0.32s;
}

@keyframes dot-bounce {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}

/* Suggestion Card */
.ai-suggestion-card {
  background: #f7f7f7;
  border: 1px solid #ebebeb;
  border-radius: 12px;
  padding: 14px;
}

.ai-suggestion-text {
  font-size: 14px;
  color: #222222;
  white-space: pre-wrap;
  line-height: 1.6;
  margin: 0;
}

/* Buttons */
.ai-btn {
  display: inline-flex;
  align-items: center;
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.ai-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.ai-btn--accept {
  background: #222222;
  color: #fff;
}

.ai-btn--accept:hover:not(:disabled) {
  background: #000000;
}

.ai-btn--reject {
  background: #ffffff;
  border: 1px solid #c1c1c1;
  color: #222222;
}

.ai-btn--reject:hover:not(:disabled) {
  background: #f7f7f7;
}

.ai-btn--primary {
  background: #c0982b;
  color: #fff;
}

.ai-btn--primary:hover:not(:disabled) {
  background: #a8862a;
}

.ai-btn--outline {
  background: #ffffff;
  border: 1px solid #c1c1c1;
  color: #222222;
}

.ai-btn--outline:hover:not(:disabled) {
  background: #f7f7f7;
}
</style>
