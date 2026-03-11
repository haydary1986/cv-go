<template>
  <div class="container py-4">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h3>{{ t('notifications.title') }}</h3>
      <div class="d-flex gap-2">
        <button @click="deleteAllRead" class="btn btn-outline-danger btn-sm" v-if="hasReadNotifications">
          <i class="fas fa-trash me-1"></i>
          {{ locale === 'ar' ? 'حذف المقروءة' : 'Delete all read' }}
        </button>
        <button @click="markAllRead" class="btn btn-outline-primary btn-sm" v-if="notifStore.notifications.length">
          {{ t('notifications.markAllRead') }}
        </button>
      </div>
    </div>

    <div v-if="notifStore.notifications.length === 0" class="text-center py-5">
      <i class="fas fa-bell-slash fa-3x text-muted mb-3"></i>
      <h5 class="text-muted">{{ t('notifications.noNotifications') }}</h5>
    </div>

    <div class="list-group">
      <component
        :is="n.cv_id ? 'router-link' : 'div'"
        v-for="n in notifStore.notifications"
        :key="n.id"
        :to="n.cv_id ? `/cv/${n.cv_id}/edit` : undefined"
        class="list-group-item list-group-item-action text-decoration-none"
        :class="{ 'bg-light': !n.is_read }"
        @click="markRead(n)"
      >
        <div class="d-flex justify-content-between">
          <div>
            <h6 class="mb-1">
              <i :class="notifIcon(n.type)" class="me-2"></i>
              {{ locale === 'ar' ? n.title_ar : n.title_en }}
              <i v-if="n.cv_id" class="fas fa-external-link-alt ms-1 small text-muted"></i>
            </h6>
            <p class="mb-0 text-muted small">{{ locale === 'ar' ? n.message_ar : n.message_en }}</p>
          </div>
          <small class="text-muted text-nowrap">{{ timeAgo(n.created_at) }}</small>
        </div>
      </component>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../stores/notification'

const { t, locale } = useI18n()
const notifStore = useNotificationStore()

const hasReadNotifications = computed(() => {
  return notifStore.notifications.some((n: any) => n.is_read)
})

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

function timeAgo(date: string): string {
  const now = new Date()
  const past = new Date(date)
  const diffMs = now.getTime() - past.getTime()
  const diffSeconds = Math.floor(diffMs / 1000)
  const diffMinutes = Math.floor(diffSeconds / 60)
  const diffHours = Math.floor(diffMinutes / 60)
  const diffDays = Math.floor(diffHours / 24)

  const isAr = locale.value === 'ar'

  if (diffSeconds < 60) {
    return isAr ? 'الآن' : 'Just now'
  } else if (diffMinutes < 60) {
    return isAr ? `منذ ${diffMinutes} دقيقة` : `${diffMinutes}m ago`
  } else if (diffHours < 24) {
    return isAr ? `منذ ${diffHours} ساعة` : `${diffHours}h ago`
  } else if (diffDays < 7) {
    return isAr ? `منذ ${diffDays} يوم` : `${diffDays}d ago`
  } else {
    return past.toLocaleDateString()
  }
}

async function markRead(n: any) {
  if (!n.is_read) {
    await notifStore.markAsRead(n.id)
  }
}

async function markAllRead() {
  await notifStore.markAllAsRead()
}

async function deleteAllRead() {
  await notifStore.deleteAllRead()
}
</script>
