import { defineStore } from 'pinia'
import { ref } from 'vue'
import { notificationAPI } from '../services/api'

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<any[]>([])
  const unreadCount = ref(0)

  async function fetchNotifications(page = 1) {
    const response = await notificationAPI.list({ page, limit: 20 })
    notifications.value = response.data.notifications || []
    return response.data
  }

  async function fetchUnreadCount() {
    const response = await notificationAPI.getUnreadCount()
    unreadCount.value = response.data.count || 0
  }

  async function markAsRead(id: number) {
    await notificationAPI.markAsRead(id)
    const n = notifications.value.find((n) => n.id === id)
    if (n) n.is_read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }

  async function markAllAsRead() {
    await notificationAPI.markAllAsRead()
    notifications.value.forEach((n) => (n.is_read = true))
    unreadCount.value = 0
  }

  return { notifications, unreadCount, fetchNotifications, fetchUnreadCount, markAsRead, markAllAsRead }
})
