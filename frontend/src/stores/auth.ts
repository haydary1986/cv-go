import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<any>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isTeacher = computed(() => user.value?.role === 'teacher')
  const userRole = computed(() => user.value?.role || 'student')

  async function login(email: string, password: string) {
    const response = await authAPI.login({ email, password })
    token.value = response.data.token
    user.value = response.data.user
    localStorage.setItem('token', response.data.token)
    return response.data
  }

  async function register(data: any) {
    const response = await authAPI.register(data)
    token.value = response.data.token
    user.value = response.data.user
    localStorage.setItem('token', response.data.token)
    return response.data
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const response = await authAPI.getProfile()
      user.value = response.data.user
    } catch {
      logout()
    }
  }

  async function updateProfile(data: any) {
    const response = await authAPI.updateProfile(data)
    user.value = response.data.user
    return response.data
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  function setGoogleAuth(t: string, u: any) {
    token.value = t
    user.value = u
    localStorage.setItem('token', t)
  }

  return {
    token,
    user,
    isAuthenticated,
    isAdmin,
    isTeacher,
    userRole,
    login,
    register,
    fetchProfile,
    updateProfile,
    logout,
    setGoogleAuth,
  }
})
