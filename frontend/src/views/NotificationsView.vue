<template>
  <div class="notifications-page">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-lg-8 col-xl-7">
          <!-- Header -->
          <div class="d-flex flex-column flex-sm-row justify-content-between align-items-start align-items-sm-center mb-4 gap-2">
            <div>
              <h3 class="fw-bold mb-1">
                <i class="fas fa-bell text-primary me-2"></i>{{ t('notifications.title') }}
              </h3>
              <p class="text-muted mb-0 small" v-if="notifStore.notifications.length">
                {{ unreadCount > 0
                  ? (locale === 'ar' ? `${unreadCount} إشعار غير مقروء` : `${unreadCount} unread notification${unreadCount > 1 ? 's' : ''}`)
                  : (locale === 'ar' ? 'جميع الإشعارات مقروءة' : 'All notifications read')
                }}
              </p>
            </div>
            <div class="d-flex gap-2" v-if="notifStore.notifications.length">
              <button
                @click="deleteAllRead"
                class="btn btn-outline-danger btn-sm d-inline-flex align-items-center"
                v-if="hasReadNotifications"
              >
                <i class="fas fa-trash-alt me-1"></i>
                {{ t('notifications.deleteAllRead') }}
              </button>
              <button
                @click="markAllRead"
                class="btn btn-outline-primary btn-sm d-inline-flex align-items-center"
                v-if="unreadCount > 0"
              >
                <i class="fas fa-check-double me-1"></i>
                {{ t('notifications.markAllRead') }}
              </button>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="notifStore.notifications.length === 0" class="empty-state card border-0 shadow-sm">
            <div class="card-body text-center py-5">
              <div class="empty-icon-circle mb-4">
                <i class="fas fa-bell-slash"></i>
              </div>
              <h5 class="fw-semibold text-muted mb-2">{{ t('notifications.noNotifications') }}</h5>
              <p class="text-muted small mb-0">{{ t('notifications.noNotificationsDesc') }}</p>
            </div>
          </div>

          <!-- Notifications List -->
          <div class="notifications-list" v-else>
            <component
              :is="n.cv_id ? 'router-link' : 'div'"
              v-for="n in notifStore.notifications"
              :key="n.id"
              :to="n.cv_id ? `/cv/${n.cv_id}/edit` : undefined"
              class="notification-card card border-0 shadow-sm mb-2 text-decoration-none"
              :class="{ 'notification-unread': !n.is_read }"
              @click="markRead(n)"
            >
              <div class="card-body p-3">
                <div class="d-flex gap-3 align-items-start">
                  <!-- Icon -->
                  <div class="notification-icon-wrapper flex-shrink-0" :class="notifIconBg(n.type)">
                    <i :class="notifIconClass(n.type)"></i>
                  </div>

                  <!-- Content -->
                  <div class="flex-grow-1 min-width-0">
                    <div class="d-flex justify-content-between align-items-start gap-2">
                      <div class="min-width-0">
                        <h6 class="mb-1" :class="{ 'fw-bold': !n.is_read, 'fw-medium': n.is_read }">
                          {{ locale === 'ar' ? n.title_ar : n.title_en }}
                          <i v-if="n.cv_id" class="fas fa-external-link-alt ms-1 small text-muted"></i>
                        </h6>
                        <p class="mb-0 text-muted small text-truncate-2">{{ locale === 'ar' ? n.message_ar : n.message_en }}</p>
                      </div>
                      <div class="d-flex flex-column align-items-end flex-shrink-0 gap-1">
                        <small class="text-muted text-nowrap">
                          <i class="fas fa-clock me-1"></i>{{ timeAgo(n.created_at) }}
                        </small>
                        <span v-if="!n.is_read" class="unread-dot"></span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </component>
          </div>
        </div>
      </div>
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

const unreadCount = computed(() => {
  return notifStore.notifications.filter((n: any) => !n.is_read).length
})

onMounted(() => {
  notifStore.fetchNotifications()
})

function notifIconClass(type: string) {
  const map: Record<string, string> = {
    approval: 'fas fa-check',
    rejection: 'fas fa-times',
    revision: 'fas fa-pen',
    announcement: 'fas fa-bullhorn',
  }
  return map[type] || 'fas fa-bell'
}

function notifIconBg(type: string) {
  const map: Record<string, string> = {
    approval: 'icon-success',
    rejection: 'icon-danger',
    revision: 'icon-warning',
    announcement: 'icon-info',
  }
  return map[type] || 'icon-primary'
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
    return past.toLocaleDateString(isAr ? 'ar-IQ' : 'en-US', {
      month: 'short',
      day: 'numeric',
    })
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

<style scoped>
.notifications-page {
  min-height: 100vh;
  background: #f4f6f9;
}

.empty-state {
  border-radius: 16px;
}

.empty-icon-circle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #f0f2f5;
  color: #adb5bd;
  font-size: 2rem;
}

.notification-card {
  border-radius: 12px;
  transition: all 0.2s ease;
  cursor: pointer;
  display: block;
}

.notification-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.notification-unread {
  border-left: 3px solid #0d6efd !important;
  background: #f0f6ff;
}

[dir="rtl"] .notification-unread {
  border-left: none !important;
  border-right: 3px solid #0d6efd !important;
}

.notification-icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.95rem;
}

.icon-success {
  background: #d1e7dd;
  color: #198754;
}

.icon-danger {
  background: #f8d7da;
  color: #dc3545;
}

.icon-warning {
  background: #fff3cd;
  color: #856d00;
}

.icon-info {
  background: #cff4fc;
  color: #0dcaf0;
}

.icon-primary {
  background: #cfe2ff;
  color: #0d6efd;
}

.unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #0d6efd;
  display: inline-block;
}

.text-truncate-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.min-width-0 {
  min-width: 0;
}

@media (max-width: 576px) {
  .notification-card .card-body {
    padding: 0.75rem !important;
  }

  .notification-icon-wrapper {
    width: 36px;
    height: 36px;
    font-size: 0.85rem;
  }
}
</style>
