import { useAuthStore } from '../stores/auth'
import { useNotificationStore } from '../stores/notification'
import { useToast } from '../composables/useToast'

let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let reconnectAttempts = 0
const MAX_RECONNECT_DELAY = 30000

export function connectWebSocket() {
  const authStore = useAuthStore()
  const token = localStorage.getItem('token')

  if (!token || !authStore.isAuthenticated) return

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  const url = `${protocol}//${host}/api/ws`

  try {
    ws = new WebSocket(url)

    ws.onopen = () => {
      reconnectAttempts = 0
      // Send token as first message for authentication (not in URL)
      ws?.send(JSON.stringify({ type: 'auth', token }))
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        if (data.type === 'notification') {
          const notifStore = useNotificationStore()
          notifStore.fetchUnreadCount()
          const toast = useToast()
          const title = data.payload?.title_ar || data.payload?.title_en || 'New notification'
          toast.info(title)
        }
      } catch {
        // ignore parse errors
      }
    }

    ws.onclose = () => {
      ws = null
      scheduleReconnect()
    }

    ws.onerror = () => {
      ws?.close()
    }
  } catch {
    scheduleReconnect()
  }
}

function scheduleReconnect() {
  if (reconnectTimer) return
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), MAX_RECONNECT_DELAY)
  reconnectAttempts++
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null
    const authStore = useAuthStore()
    if (authStore.isAuthenticated) {
      connectWebSocket()
    }
  }, delay)
}

export function disconnectWebSocket() {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  reconnectAttempts = 0
  if (ws) {
    ws.close()
    ws = null
  }
}
