<template>
  <div class="notifications-page">
    <div class="container py-4">
      <div class="row justify-content-center">
        <div class="col-lg-8 col-xl-7">
          <!-- Header -->
          <div class="notif-header mb-4">
            <div class="d-flex flex-column flex-sm-row justify-content-between align-items-start align-items-sm-center gap-3">
              <div class="d-flex align-items-center gap-3">
                <div class="notif-header-icon">
                  <i class="fas fa-bell"></i>
                </div>
                <div>
                  <h3 class="fw-bold mb-0 notif-title">
                    {{ t('notifications.title') }}
                    <span v-if="unreadCount > 0" class="unread-badge">{{ unreadCount }}</span>
                  </h3>
                  <p class="text-muted-light mb-0 small mt-1" v-if="notifStore.notifications.length">
                    {{ unreadCount > 0
                      ? (locale === 'ar' ? `${unreadCount} إشعار غير مقروء` : `${unreadCount} unread notification${unreadCount > 1 ? 's' : ''}`)
                      : (locale === 'ar' ? 'جميع الإشعارات مقروءة' : 'All notifications read')
                    }}
                  </p>
                </div>
              </div>
              <div class="d-flex gap-2" v-if="notifStore.notifications.length">
                <button
                  @click="deleteAllRead"
                  class="btn btn-notif-action btn-notif-delete d-inline-flex align-items-center"
                  v-if="hasReadNotifications"
                >
                  <i class="fas fa-trash-alt me-1"></i>
                  {{ t('notifications.deleteAllRead') }}
                </button>
                <button
                  @click="markAllRead"
                  class="btn btn-notif-action btn-notif-markread d-inline-flex align-items-center"
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
              <h5 class="fw-semibold empty-title mb-2">{{ t('notifications.noNotifications') }}</h5>
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
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-secondary: #2c3e50;
  --uni-primary-light: #e8f0f7;
  --uni-success: #1e8449;
  --uni-warning: #d4a017;
  --uni-danger: #c0392b;
  --uni-info: #2471a3;
  min-height: 100vh;
  background: #f5f7fa;
}

/* === Header === */
.notif-header {
  background: linear-gradient(135deg, var(--uni-primary) 0%, var(--uni-secondary) 100%);
  border-radius: 16px;
  padding: 24px 28px;
  color: white;
}

.notif-header-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  color: var(--uni-accent);
  flex-shrink: 0;
}

.notif-title {
  color: white;
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
  background: var(--uni-accent);
  color: var(--uni-primary);
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0 7px;
}

.text-muted-light {
  color: rgba(255, 255, 255, 0.65);
}

.btn-notif-action {
  font-size: 0.82rem;
  font-weight: 500;
  border-radius: 8px;
  padding: 6px 14px;
  transition: all 0.2s;
}

.btn-notif-delete {
  background: rgba(192, 57, 43, 0.15);
  color: #f1a9a0;
  border: 1px solid rgba(241, 169, 160, 0.3);
}

.btn-notif-delete:hover {
  background: var(--uni-danger);
  color: white;
  border-color: var(--uni-danger);
}

.btn-notif-markread {
  background: rgba(192, 152, 43, 0.15);
  color: var(--uni-accent);
  border: 1px solid rgba(192, 152, 43, 0.3);
}

.btn-notif-markread:hover {
  background: var(--uni-accent);
  color: var(--uni-primary);
  border-color: var(--uni-accent);
}

/* === Empty State === */
.empty-state {
  background: white;
  border-radius: 16px;
  border: 1px solid #e4e8f0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.empty-bell-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--uni-primary-light) 0%, #d5e8f7 100%);
  color: var(--uni-primary);
}

.empty-title {
  color: var(--uni-secondary);
}

/* === Notification Cards === */
.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.notification-card {
  display: flex;
  background: white;
  border-radius: 14px;
  overflow: hidden;
  transition: all 0.2s ease;
  cursor: pointer;
  border: 1px solid #e4e8f0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.03);
}

.notification-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(26, 82, 118, 0.1);
  border-color: #d0d8e4;
}

/* Color Strip */
.notification-color-strip {
  width: 4px;
  flex-shrink: 0;
}

.strip-approval {
  background: var(--uni-success);
}

.strip-rejection {
  background: var(--uni-danger);
}

.strip-revision {
  background: var(--uni-warning);
}

.strip-announcement {
  background: var(--uni-info);
}

[dir="rtl"] .notification-color-strip {
  order: 1;
}

/* Unread state */
.notification-unread {
  background: #f8fbff;
  border-color: rgba(26, 82, 118, 0.15);
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
  width: 42px;
  height: 42px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.95rem;
  flex-shrink: 0;
}

.icon-success {
  background: linear-gradient(135deg, #d4efdf, #a9dfbf);
  color: var(--uni-success);
}

.icon-danger {
  background: linear-gradient(135deg, #fadbd8, #f1948a);
  color: var(--uni-danger);
}

.icon-warning {
  background: linear-gradient(135deg, #fdebd0, #f9e79f);
  color: #9a7d0a;
}

.icon-info {
  background: linear-gradient(135deg, #d4e6f1, #a9cce3);
  color: var(--uni-info);
}

.icon-primary {
  background: linear-gradient(135deg, var(--uni-primary-light), #bdd7ea);
  color: var(--uni-primary);
}

.notification-title {
  color: var(--uni-secondary);
  font-size: 0.92rem;
  line-height: 1.4;
}

.link-indicator {
  color: #aeb6bf;
  font-size: 0.7rem;
}

.notification-message {
  color: #7f8c8d;
  font-size: 0.84rem;
  line-height: 1.5;
}

.notification-time {
  color: #aeb6bf;
  font-size: 0.78rem;
}

.unread-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--uni-accent);
  display: inline-block;
  box-shadow: 0 0 0 3px rgba(192, 152, 43, 0.2);
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

/* === Responsive === */
@media (max-width: 576px) {
  .notif-header {
    padding: 18px 20px;
  }

  .notif-header-icon {
    width: 40px;
    height: 40px;
    font-size: 1rem;
  }

  .notification-body {
    padding: 12px 14px;
  }

  .notification-icon-wrapper {
    width: 36px;
    height: 36px;
    font-size: 0.85rem;
    border-radius: 10px;
  }

  .notification-title {
    font-size: 0.86rem;
  }
}
</style>
