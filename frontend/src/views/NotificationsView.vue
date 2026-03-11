<template>
  <div class="container py-4">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h3>{{ t('notifications.title') }}</h3>
      <button @click="markAllRead" class="btn btn-outline-primary btn-sm" v-if="notifStore.notifications.length">
        {{ t('notifications.markAllRead') }}
      </button>
    </div>

    <div v-if="notifStore.notifications.length === 0" class="text-center py-5">
      <i class="fas fa-bell-slash fa-3x text-muted mb-3"></i>
      <h5 class="text-muted">{{ t('notifications.noNotifications') }}</h5>
    </div>

    <div class="list-group">
      <div
        v-for="n in notifStore.notifications"
        :key="n.id"
        class="list-group-item list-group-item-action"
        :class="{ 'bg-light': !n.is_read }"
        @click="markRead(n)"
      >
        <div class="d-flex justify-content-between">
          <div>
            <h6 class="mb-1">
              <i :class="notifIcon(n.type)" class="me-2"></i>
              {{ locale === 'ar' ? n.title_ar : n.title_en }}
            </h6>
            <p class="mb-0 text-muted small">{{ locale === 'ar' ? n.message_ar : n.message_en }}</p>
          </div>
          <small class="text-muted">{{ formatDate(n.created_at) }}</small>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../stores/notification'

const { t, locale } = useI18n()
const notifStore = useNotificationStore()

onMounted(() => {
  notifStore.fetchNotifications()
})

function notifIcon(type: string) {
  const map: Record<string, string> = {
    approval: 'fas fa-check-circle text-success',
    rejection: 'fas fa-times-circle text-danger',
    revision: 'fas fa-edit text-warning',
    announcement: 'fas fa-bullhorn text-info',
  }
  return map[type] || 'fas fa-bell text-primary'
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}

async function markRead(n: any) {
  if (!n.is_read) {
    await notifStore.markAsRead(n.id)
  }
}

async function markAllRead() {
  await notifStore.markAllAsRead()
}
</script>
