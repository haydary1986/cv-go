<template>
  <div class="toast-container" :class="{ 'toast-container--rtl': isRtl }">
    <TransitionGroup name="toast-slide">
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="toast-item d-flex align-items-center mb-2"
        :class="['toast-item--' + toast.type]"
        role="alert"
      >
        <div class="toast-icon-wrap">
          <i :class="iconClass(toast.type)"></i>
        </div>
        <span class="toast-message flex-grow-1">{{ toast.message }}</span>
        <button
          type="button"
          class="toast-close ms-2"
          @click="removeToast(toast.id)"
          aria-label="Close"
        >
          <i class="fas fa-times"></i>
        </button>
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
  top: 1.25rem;
  right: 1.25rem;
  z-index: 9999;
  max-width: 380px;
  width: 100%;
}

.toast-container--rtl {
  right: auto;
  left: 1.25rem;
}

.toast-item {
  padding: 0;
  border-radius: 12px;
  font-size: 0.88rem;
  pointer-events: auto;
  background: #ffffff;
  box-shadow: rgba(0,0,0,0.02) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 6px, rgba(0,0,0,0.1) 0px 4px 8px;
  border-left: 4px solid transparent;
  overflow: hidden;
}

.toast-item--success {
  border-left-color: #2e7d32;
}

.toast-item--error {
  border-left-color: #c62828;
}

.toast-item--warning {
  border-left-color: #c0982b;
}

.toast-item--info {
  border-left-color: #1a5276;
}

.toast-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  min-height: 48px;
  font-size: 1rem;
}

.toast-item--success .toast-icon-wrap {
  color: #2e7d32;
}

.toast-item--error .toast-icon-wrap {
  color: #c62828;
}

.toast-item--warning .toast-icon-wrap {
  color: #c0982b;
}

.toast-item--info .toast-icon-wrap {
  color: #1a5276;
}

.toast-message {
  padding: 0.75rem 0.25rem;
  color: #222222;
  line-height: 1.4;
  font-weight: 500;
}

.toast-close {
  background: none;
  border: none;
  color: #c1c1c1;
  padding: 0.75rem;
  cursor: pointer;
  font-size: 0.75rem;
  transition: color 0.15s;
  line-height: 1;
}

.toast-close:hover {
  color: #222222;
}

/* RTL border flip */
.toast-container--rtl .toast-item {
  border-left: none;
  border-right: 4px solid transparent;
}

.toast-container--rtl .toast-item--success {
  border-right-color: #2e7d32;
}

.toast-container--rtl .toast-item--error {
  border-right-color: #c62828;
}

.toast-container--rtl .toast-item--warning {
  border-right-color: #c0982b;
}

.toast-container--rtl .toast-item--info {
  border-right-color: #1a5276;
}

/* Slide-in from top-right */
.toast-slide-enter-active {
  transition: opacity 0.35s ease, transform 0.35s cubic-bezier(0.21, 1.02, 0.73, 1);
}

.toast-slide-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}

.toast-slide-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-slide-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-container--rtl .toast-slide-enter-from {
  transform: translateX(-100%);
}

.toast-container--rtl .toast-slide-leave-to {
  transform: translateX(-100%);
}
</style>
