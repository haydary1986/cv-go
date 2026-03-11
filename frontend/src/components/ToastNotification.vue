<template>
  <div class="toast-container" :class="{ 'toast-container--rtl': isRtl }">
    <TransitionGroup name="toast-fade">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="toast-item d-flex align-items-center shadow-sm mb-2"
        :class="toastClass(toast.type)"
        role="alert"
      >
        <i class="me-2" :class="iconClass(toast.type)"></i>
        <span class="flex-grow-1">{{ toast.message }}</span>
        <button
          type="button"
          class="btn-close btn-close-sm ms-2"
          :class="{ 'btn-close-white': toast.type === 'error' }"
          @click="removeToast(toast.id)"
          aria-label="Close"
        ></button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from '../composables/useToast'

const { locale } = useI18n()
const { toasts, removeToast } = useToast()

const isRtl = computed(() => locale.value === 'ar')

function toastClass(type: string) {
  const map: Record<string, string> = {
    success: 'alert alert-success',
    error: 'alert alert-danger',
    warning: 'alert alert-warning',
    info: 'alert alert-info',
  }
  return map[type] ?? 'alert alert-secondary'
}

function iconClass(type: string) {
  const map: Record<string, string> = {
    success: 'fas fa-check-circle',
    error: 'fas fa-times-circle',
    warning: 'fas fa-exclamation-triangle',
    info: 'fas fa-info-circle',
  }
  return map[type] ?? 'fas fa-bell'
}
</script>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: 9999;
  max-width: 360px;
  width: 100%;
}

.toast-container--rtl {
  right: auto;
  left: 1.5rem;
}

.toast-item {
  padding: 0.75rem 1rem;
  border-radius: 0.375rem;
  font-size: 0.9rem;
  pointer-events: auto;
}

.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.toast-fade-enter-from {
  opacity: 0;
  transform: translateY(0.5rem);
}

.toast-fade-leave-to {
  opacity: 0;
  transform: translateX(1rem);
}

.toast-container--rtl .toast-fade-leave-to {
  transform: translateX(-1rem);
}
</style>
