import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { publicAPI } from '../services/api'

export const useBrandingStore = defineStore('branding', () => {
  const branding = ref({
    name: 'جامعة التراث',
    logo_url: '',
    primary_color: '#1a5276',
    secondary_color: '#2c3e50',
    accent_color: '#c0982b',
  })
  const loaded = ref(false)

  async function fetchBranding() {
    try {
      const res = await publicAPI.getBranding()
      if (res.data.branding) {
        branding.value = { ...branding.value, ...res.data.branding }
      }
      loaded.value = true
      applyColors()
    } catch {
      loaded.value = true
    }
  }

  function applyColors() {
    const root = document.documentElement
    root.style.setProperty('--bs-primary', branding.value.primary_color)
    root.style.setProperty('--bs-secondary', branding.value.secondary_color)
    root.style.setProperty('--app-primary', branding.value.primary_color)
    root.style.setProperty('--app-secondary', branding.value.secondary_color)
    root.style.setProperty('--app-accent', branding.value.accent_color)
  }

  const systemName = computed(() => branding.value.name || 'جامعة التراث')
  const logoUrl = computed(() => branding.value.logo_url || '')

  return { branding, loaded, fetchBranding, applyColors, systemName, logoUrl }
})
