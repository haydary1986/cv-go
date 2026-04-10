<template>
  <div class="notifications-page">
    <div class="container py-5">
      <div class="row justify-content-center">
        <div class="col-lg-8 col-xl-7">
          <!-- Header -->
          <div class="notif-header mb-4">
            <div class="d-flex flex-column flex-sm-row justify-content-between align-items-start align-items-sm-center gap-3">
              <div>
                <h3 class="notif-title mb-0">
                  {{ t('notifications.title') }}
                  <span v-if="unreadCount > 0" class="unread-badge">{{ unreadCount }}</span>
                </h3>
                <p class="notif-subtitle mb-0 mt-1" v-if="notifStore.notifications.length">
                  {{ unreadCount > 0
                    ? (locale === 'ar' ? `${unreadCount} إشعار غير مقروء` : `${unreadCount} unread notification${unreadCount > 1 ? 's' : ''}`)
                    : (locale === 'ar' ? 'جميع الإشعارات مقروءة' : 'All notifications read')
                  }}
                </p>
              </div>
              <div class="d-flex gap-2" v-if="notifStore.notifications.length">
                <button
                  @click="deleteAllRead"
                  class="btn-notif-action btn-notif-delete d-inline-flex align-items-center"
                  v-if="hasReadNotifications"
                >
                  <i class="fas fa-trash-alt me-1"></i>
                  {{ t('notifications.deleteAllRead') }}
                </button>
                <button
                  @click="markAllRead"
                  class="btn-notif-action btn-notif-markread d-inline-flex align-items-center"
                  v-if="unreadCount > 0"
                >
                  <i class="fas fa-check-double me-1"></i>
                  {{ t('notifications.markAllRead') }}
                </button>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="notifStore.notifications.length === 0" class="empty-state">
            <div class="empty-content text-center py-5">
              <div class="empty-bell-wrap mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
                  <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
                  <line x1="1" y1="1" x2="23" y2="23"></line>
                </svg>
              </div>
              <h5 class="empty-title mb-2">{{ t('notifications.noNotifications') }}</h5>
              <p class="empty-desc mb-0">{{ t('notifications.noNotificationsDesc') }}</p>
            </div>
          </div>

          <!-- Notifications List -->
          <div class="notifications-list" v-else>
            <component
              :is="n.cv_id ? 'router-link' : 'div'"
              v-for="n in notifStore.notifications"
              :key="n.id"
              :to="n.cv_id ? `/cv/${n.cv_id}/edit` : undefined"
              class="notification-card text-decoration-none"
              :class="[
                { 'notification-unread': !n.is_read },
                `notification-type-${n.type}`
              ]"
              @click="markRead(n)"
            >
              <div class="notification-color-strip" :class="`strip-${n.type}`"></div>
              <div class="notification-body">
                <div class="d-flex gap-3 align-items-start">
                  <!-- Icon -->
                  <div class="notification-icon-wrapper" :class="notifIconBg(n.type)">
                    <i :class="notifIconClass(n.type)"></i>
                  </div>

                  <!-- Content -->
                  <div class="flex-grow-1 min-width-0">
                    <div class="d-flex justify-content-between align-items-start gap-2">
                      <div class="min-width-0">
                        <h6 class="mb-1 notification-title" :class="{ 'fw-bold': !n.is_read, 'fw-medium': n.is_read }">
                          {{ locale === 'ar' ? n.title_ar : n.title_en }}
                          <i v-if="n.cv_id" class="fas fa-external-link-alt ms-1 small link-indicator"></i>
                        </h6>
                        <p class="mb-0 notification-message text-truncate-2">{{ locale === 'ar' ? n.message_ar : n.message_en }}</p>
                      </div>
                      <div class="d-flex flex-column align-items-end flex-shrink-0 gap-1">
                        <small class="notification-time text-nowrap">
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
/* === Page === */
.notifications-page {
  min-height: 100vh;
  background: #ffffff;
}

/* === Header === */
.notif-header {
  padding: 0;
}

.notif-title {
  font-weight: 700;
  color: #222222;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.unread-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  border-radius: 12px;
  background: #c0982b;
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0 7px;
}

.notif-subtitle {
  color: #6a6a6a;
  font-size: 0.88rem;
  font-weight: 500;
}

.btn-notif-action {
  font-size: 0.82rem;
  font-weight: 600;
  border-radius: 8px;
  padding: 6px 14px;
  transition: all 0.2s;
  cursor: pointer;
  border: none;
}

.btn-notif-delete {
  background: rgba(192, 57, 43, 0.08);
  color: #c0392b;
}

.btn-notif-delete:hover {
  background: #c0392b;
  color: #ffffff;
}

.btn-notif-markread {
  background: transparent;
  color: #1a5276;
  border: none;
  padding: 6px 14px;
}

.btn-notif-markread:hover {
  background: rgba(26, 82, 118, 0.06);
}

/* === Empty State === */
.empty-state {
  background: #ffffff;
  border-radius: 12px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.empty-bell-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: #f7f7f7;
  color: #6a6a6a;
}

.empty-title {
  color: #222222;
  font-weight: 600;
}

.empty-desc {
  color: #6a6a6a;
  font-size: 0.9rem;
}

/* === Notification Cards === */
.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.notification-card {
  display: flex;
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s ease;
  cursor: pointer;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px;
}

.notification-card:hover {
  box-shadow: rgba(0, 0, 0, 0.08) 0px 4px 12px;
}

/* Color Strip */
.notification-color-strip {
  width: 4px;
  flex-shrink: 0;
}

.strip-approval { background: #1e8449; }
.strip-rejection { background: #c0392b; }
.strip-revision { background: #c0982b; }
.strip-announcement { background: #1a5276; }

[dir="rtl"] .notification-color-strip {
  order: 1;
}

/* Unread state */
.notification-unread {
  background: #fafafa;
}

.notification-unread .notification-color-strip {
  width: 5px;
}

/* Body */
.notification-body {
  flex: 1;
  padding: 16px 18px;
  min-width: 0;
}

.notification-icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.icon-success {
  background: rgba(30, 132, 73, 0.1);
  color: #1e8449;
}

.icon-danger {
  background: rgba(192, 57, 43, 0.1);
  color: #c0392b;
}

.icon-warning {
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
}

.icon-info {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.icon-primary {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.notification-title {
  color: #222222;
  font-size: 0.92rem;
  line-height: 1.4;
}

.link-indicator {
  color: #6a6a6a;
  font-size: 0.7rem;
}

.notification-message {
  color: #6a6a6a;
  font-size: 0.84rem;
  line-height: 1.5;
}

.notification-time {
  color: #6a6a6a;
  font-size: 0.78rem;
}

.unread-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: #c0982b;
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

.fw-medium {
  font-weight: 500;
}

/* === Responsive === */
@media (max-width: 576px) {
  .notification-body {
    padding: 12px 14px;
  }

  .notification-icon-wrapper {
    width: 36px;
    height: 36px;
    font-size: 0.85rem;
  }

  .notification-title {
    font-size: 0.86rem;
  }
}
</style>
