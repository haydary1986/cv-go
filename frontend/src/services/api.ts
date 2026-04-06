import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor - add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor - handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authAPI = {
  login: (data: { email: string; password: string }) => api.post('/auth/login', data),
  register: (data: any) => api.post('/auth/register', data),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (data: any) => api.put('/auth/profile', data),
  changePassword: (data: any) => api.put('/auth/change-password', data),
  googleLogin: () => window.location.href = `${api.defaults.baseURL}/auth/google`,
}

// CV API
export const cvAPI = {
  list: (params?: any) => api.get('/cvs', { params }),
  get: (id: number) => api.get(`/cvs/${id}`),
  create: (data: any) => api.post('/cvs', data),
  update: (id: number, data: any) => api.put(`/cvs/${id}`, data),
  delete: (id: number) => api.delete(`/cvs/${id}`),
  toggleShare: (id: number) => api.post(`/cvs/${id}/toggle-share`),
  getShared: (token: string) => api.get(`/shared/${token}`),
  exportJSON: () => api.get('/cvs/export/json', { responseType: 'blob' }),
  exportCSV: () => api.get('/cvs/export/csv', { responseType: 'blob' }),
  exportPDF: (id: number) => api.get(`/cvs/${id}/export/pdf`, { responseType: 'blob' }),
  createGuest: (data: any) => api.post('/guest/cv', data),
  importLinkedIn: (formData: FormData) => api.post('/cvs/import/linkedin', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
}

// AI API
export const aiAPI = {
  improveText: (data: { text: string; language: string }) => api.post('/ai/improve-text', data),
  analyzeCV: (cvId: number) => api.post(`/ai/analyze/${cvId}`),
  coverLetter: (data: { cv_id: number; job_title: string; company: string; language: string }) =>
    api.post('/ai/cover-letter', data),
  suggestJobs: (cvId: number) => api.post(`/ai/suggest-jobs/${cvId}`),
  evaluateResearch: (cvId: number) => api.post(`/ai/evaluate-research/${cvId}`),
  getTips: (cvId: number) => api.post(`/ai/tips/${cvId}`),
}

// Notification API
export const notificationAPI = {
  list: (params?: any) => api.get('/notifications', { params }),
  markAsRead: (id: number) => api.put(`/notifications/${id}/read`),
  markAllAsRead: () => api.put('/notifications/read-all'),
  getUnreadCount: () => api.get('/notifications/unread-count'),
  delete: (id: number) => api.delete(`/notifications/${id}`),
}

// Admin API
export const adminAPI = {
  getDashboardStats: () => api.get('/admin/dashboard'),
  listCVs: (params?: any) => api.get('/admin/cvs', { params }),
  approveCV: (id: number) => api.post(`/admin/cvs/${id}/approve`),
  rejectCV: (id: number, reason: string) => api.post(`/admin/cvs/${id}/reject`, { reason }),
  requestRevision: (id: number, note: string) => api.post(`/admin/cvs/${id}/revision`, { note }),
  listUsers: (params?: any) => api.get('/admin/users', { params }),
  createUser: (data: any) => api.post('/admin/users', data),
  updateUser: (id: number, data: any) => api.put(`/admin/users/${id}`, data),
  deleteUser: (id: number) => api.delete(`/admin/users/${id}`),
  updateCredits: (userId: number, credits: number) =>
    api.put(`/admin/users/${userId}/credits`, { credits }),
  listFaculties: () => api.get('/admin/faculties'),
  createFaculty: (data: any) => api.post('/admin/faculties', data),
  updateFaculty: (id: number, data: any) => api.put(`/admin/faculties/${id}`, data),
  deleteFaculty: (id: number) => api.delete(`/admin/faculties/${id}`),
  listDepartments: (facultyId?: number) =>
    api.get('/admin/departments', { params: { faculty_id: facultyId } }),
  createDepartment: (data: any) => api.post('/admin/departments', data),
  updateDepartment: (id: number, data: any) => api.put(`/admin/departments/${id}`, data),
  deleteDepartment: (id: number) => api.delete(`/admin/departments/${id}`),
  getBranding: () => api.get('/admin/branding'),
  updateBranding: (data: any) => api.put('/admin/branding', data),
  getAISettings: () => api.get('/admin/ai-settings'),
  updateAISettings: (data: any) => api.put('/admin/ai-settings', data),
  getAdSettings: () => api.get('/admin/ad-settings'),
  updateAdSettings: (data: any) => api.put('/admin/ad-settings', data),
  uploadLogo: (formData: FormData) => api.post('/admin/branding/logo', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  sendNotification: (data: any) => api.post('/admin/notifications', data),
  getActivityLogs: (params?: any) => api.get('/admin/activity-logs', { params }),
  getStats: () => api.get('/admin/stats'),
  getAuditTrail: (params?: any) => api.get('/admin/audit-trail', { params }),
  exportUsersCSV: () => api.get('/admin/users/export/csv', { responseType: 'blob' }),
  importUsersCSV: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/admin/users/import/csv', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

// 2FA API
export const twoFAAPI = {
  setup: () => api.post('/auth/2fa/setup'),
  verifySetup: (code: string) => api.post('/auth/2fa/verify-setup', { code }),
  disable: (password: string) => api.post('/auth/2fa/disable', { password }),
  validate: (data: { temp_token: string; code: string }) => api.post('/auth/2fa/validate', data),
}

// Public API
export const publicAPI = {
  getFaculties: () => api.get('/v1/faculties'),
  getDepartments: (facultyId?: number) =>
    api.get('/v1/departments', { params: { faculty_id: facultyId } }),
  getBranding: () => api.get('/branding'),
}

export default api
