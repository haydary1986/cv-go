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
        <small class="text-muted mt-2 d-block">{{ locale === 'ar' ? 'جاري التحسين...' : 'Improving...' }}</small>
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
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  position: absolute;
  z-index: 1050;
  width: 370px;
  max-width: 90vw;
  right: 0;
  top: 100%;
  margin-top: 6px;
  border-radius: 16px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(26, 82, 118, 0.1);
  box-shadow: 0 8px 32px rgba(26, 82, 118, 0.15), 0 2px 8px rgba(0, 0, 0, 0.06);
}

.ai-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.7rem 1rem;
  background: linear-gradient(135deg, var(--uni-primary), #1e6291);
  color: #fff;
}

.ai-panel-title {
  font-size: 0.85rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 0.35rem;
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
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: #fff;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.7rem;
  cursor: pointer;
  transition: background 0.15s;
}

.ai-panel-close:hover {
  background: rgba(255, 255, 255, 0.35);
}

.ai-panel-body {
  padding: 1rem;
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
  background: var(--uni-accent);
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
  background: rgba(26, 82, 118, 0.04);
  border: 1px solid rgba(26, 82, 118, 0.1);
  border-radius: 10px;
  padding: 0.85rem;
}

.ai-suggestion-text {
  font-size: 0.85rem;
  color: #2c3e50;
  white-space: pre-wrap;
  line-height: 1.6;
  margin: 0;
}

/* Buttons */
.ai-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 0.9rem;
  border-radius: 8px;
  font-size: 0.8rem;
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
  background: #28a745;
  color: #fff;
}

.ai-btn--accept:hover:not(:disabled) {
  background: #218838;
}

.ai-btn--reject {
  background: transparent;
  border: 1.5px solid #dc3545;
  color: #dc3545;
}

.ai-btn--reject:hover:not(:disabled) {
  background: rgba(220, 53, 69, 0.08);
}

.ai-btn--primary {
  background: var(--uni-accent);
  color: #fff;
}

.ai-btn--primary:hover:not(:disabled) {
  background: #d4a82f;
}

.ai-btn--outline {
  background: transparent;
  border: 1.5px solid var(--uni-primary);
  color: var(--uni-primary);
}

.ai-btn--outline:hover:not(:disabled) {
  background: rgba(26, 82, 118, 0.06);
}
</style>
