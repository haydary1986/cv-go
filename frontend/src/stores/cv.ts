import { defineStore } from 'pinia'
import { ref } from 'vue'
import { cvAPI } from '../services/api'

export interface CVData {
  personal_info: {
    full_name: string
    email: string
    phone: string
    address: string
    date_of_birth: string
    nationality: string
    linkedin: string
    website: string
    github: string
    job_title: string
  }
  objective: string
  experiences: Array<{
    title: string
    company: string
    location: string
    start_date: string
    end_date: string
    current: boolean
    description: string
  }>
  education: Array<{
    degree: string
    institution: string
    location: string
    start_date: string
    end_date: string
    gpa: string
    description: string
  }>
  skills: Array<{ name: string; level: string }>
  languages: Array<{ name: string; level: string }>
  links: Array<{ title: string; url: string; type: string }>
  projects: Array<{
    name: string
    description: string
    url: string
    start_date: string
    end_date: string
  }>
  certificates: Array<{
    name: string
    issuer: string
    date: string
    expiry_date: string
    credential_id: string
    url: string
  }>
  references: Array<{
    name: string
    position: string
    company: string
    email: string
    phone: string
  }>
  photo: string
}

export interface CV {
  id: number
  user_id: number
  language: string
  template: string
  status: string
  reject_note: string
  title: string
  data: CVData
  share_token: string
  is_shared: boolean
  view_count: number
  qr_code_data: string
  created_at: string
  updated_at: string
}

export function getEmptyCVData(): CVData {
  return {
    personal_info: {
      full_name: '',
      email: '',
      phone: '',
      address: '',
      date_of_birth: '',
      nationality: '',
      linkedin: '',
      website: '',
      github: '',
      job_title: '',
    },
    objective: '',
    experiences: [],
    education: [],
    skills: [],
    languages: [],
    links: [],
    projects: [],
    certificates: [],
    references: [],
    photo: '',
  }
}

export const useCVStore = defineStore('cv', () => {
  const cvs = ref<CV[]>([])
  const currentCV = ref<CV | null>(null)
  const loading = ref(false)
  const totalPages = ref(1)

  async function fetchCVs(page = 1) {
    loading.value = true
    try {
      const response = await cvAPI.list({ page, limit: 12 })
      cvs.value = response.data.cvs || []
      totalPages.value = response.data.total_pages || 1
    } finally {
      loading.value = false
    }
  }

  async function fetchCV(id: number) {
    loading.value = true
    try {
      const response = await cvAPI.get(id)
      currentCV.value = response.data.cv
      return response.data.cv
    } finally {
      loading.value = false
    }
  }

  async function createCV(data: { title: string; language: string; template: string; data: CVData }) {
    const response = await cvAPI.create(data)
    return response.data.cv
  }

  async function updateCV(id: number, data: any) {
    const response = await cvAPI.update(id, data)
    currentCV.value = response.data.cv
    return response.data.cv
  }

  async function deleteCV(id: number) {
    await cvAPI.delete(id)
    cvs.value = cvs.value.filter((cv) => cv.id !== id)
  }

  async function toggleShare(id: number) {
    const response = await cvAPI.toggleShare(id)
    if (currentCV.value?.id === id) {
      currentCV.value = response.data.cv
    }
    return response.data.cv
  }

  return {
    cvs,
    currentCV,
    loading,
    totalPages,
    fetchCVs,
    fetchCV,
    createCV,
    updateCV,
    deleteCV,
    toggleShare,
  }
})
