<template>
  <div class="dashboard-page">
    <div class="container py-5">
      <!-- Welcome Header -->
      <div class="welcome-section mb-5">
        <div class="d-flex align-items-center mb-1">
          <div>
            <p class="welcome-date mb-0">{{ new Date().toLocaleDateString(locale === 'ar' ? 'ar-IQ' : 'en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</p>
            <h2 class="welcome-title mb-0">
              {{ locale === 'ar' ? 'مرحباً' : 'Hello' }}<span v-if="userName" class="welcome-name">, {{ userName }}</span>
            </h2>
            <p class="welcome-subtitle mb-0">{{ t('dashboard.subtitle') }}</p>
          </div>
        </div>
      </div>

      <!-- Stats Row -->
      <div class="row g-3 mb-5">
        <div class="col-6 col-sm-3">
          <div class="stat-card">
            <div class="stat-icon-circle stat-icon-blue">
              <i class="fas fa-file-alt"></i>
            </div>
            <div class="stat-number">{{ totalCount }}</div>
            <div class="stat-label">{{ t('dashboard.totalCVs') }}</div>
          </div>
        </div>
        <div class="col-6 col-sm-3">
          <div class="stat-card">
            <div class="stat-icon-circle stat-icon-green">
              <i class="fas fa-check-circle"></i>
            </div>
            <div class="stat-number">{{ approvedCount }}</div>
            <div class="stat-label">{{ t('cv.approved') }}</div>
          </div>
        </div>
        <div class="col-6 col-sm-3">
          <div class="stat-card">
            <div class="stat-icon-circle stat-icon-gold">
              <i class="fas fa-clock"></i>
            </div>
            <div class="stat-number">{{ pendingCount }}</div>
            <div class="stat-label">{{ t('cv.pending') }}</div>
          </div>
        </div>
        <div class="col-6 col-sm-3">
          <div class="stat-card">
            <div class="stat-icon-circle stat-icon-red">
              <i class="fas fa-times-circle"></i>
            </div>
            <div class="stat-number">{{ cvStore.cvs.filter(cv => cv.status === 'rejected').length }}</div>
            <div class="stat-label">{{ t('cv.rejected') }}</div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="quick-actions mb-5">
        <div class="d-flex flex-wrap gap-3">
          <router-link to="/cv/create" class="quick-action-btn">
            <i class="fas fa-plus-circle me-2"></i>{{ t('app.createCV') }}
          </router-link>
          <button class="quick-action-btn quick-action-btn--secondary" @click="exportJSON">
            <i class="fas fa-download me-2"></i>{{ t('admin.exportJSON') }}
          </button>
          <router-link to="/cv/create" class="quick-action-btn quick-action-btn--secondary">
            <i class="fas fa-palette me-2"></i>{{ locale === 'ar' ? 'تصفح القوالب' : 'Browse Templates' }}
          </router-link>
        </div>
      </div>

      <!-- Search & Filter Bar -->
      <div class="toolbar-section mb-4">
        <div class="row align-items-center g-3">
          <div class="col-md-5">
            <div class="search-box">
              <i class="fas fa-search search-icon"></i>
              <input
                type="text"
                class="form-control search-input"
                :placeholder="t('dashboard.searchPlaceholder')"
                v-model="searchQuery"
              />
              <button v-if="searchQuery" class="search-clear" @click="searchQuery = ''">
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
          <div class="col-md-7">
            <div class="d-flex align-items-center justify-content-between flex-wrap gap-2">
              <div class="filter-pills d-flex flex-wrap gap-2">
                <button
                  v-for="filter in statusFilters"
                  :key="filter.value"
                  class="filter-pill"
                  :class="{ active: activeFilter === filter.value }"
                  @click="activeFilter = filter.value"
                >
                  <i :class="filter.icon" class="me-1"></i>
                  {{ filter.label }}
                  <span v-if="filter.count > 0" class="filter-count">{{ filter.count }}</span>
                </button>
              </div>
              <router-link to="/cv/create" class="btn btn-primary btn-create d-none d-md-inline-flex">
                <i class="fas fa-plus me-1"></i>{{ t('app.createCV') }}
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Loading Skeleton -->
      <div v-if="cvStore.loading" class="row g-4">
        <div class="col-sm-6 col-lg-4" v-for="n in 6" :key="n">
          <div class="cv-card-skeleton">
            <div class="skeleton-thumbnail"></div>
            <div class="skeleton-body">
              <div class="skeleton-line skeleton-title"></div>
              <div class="skeleton-line skeleton-badge"></div>
              <div class="skeleton-line skeleton-text"></div>
              <div class="skeleton-line skeleton-text short"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="cvStore.cvs.length === 0" class="empty-state text-center py-5">
        <div class="empty-state-icon mb-4">
          <div class="empty-icon-circle">
            <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
          </div>
        </div>
        <h3 class="empty-title mb-2">{{ t('dashboard.emptyTitle') }}</h3>
        <p class="empty-description mb-4 mx-auto">{{ t('dashboard.emptyDescription') }}</p>
        <router-link to="/cv/create" class="btn-create-first">
          <i class="fas fa-plus-circle me-2"></i>{{ t('dashboard.createFirst') }}
        </router-link>
        <div class="empty-features mt-5">
          <div class="row g-4 justify-content-center">
            <div class="col-md-4">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-palette"></i></div>
                <h6>{{ t('dashboard.featureTemplates') }}</h6>
                <p class="feature-desc mb-0">{{ t('dashboard.featureTemplatesDesc') }}</p>
              </div>
            </div>
            <div class="col-md-4">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-robot"></i></div>
                <h6>{{ t('dashboard.featureAI') }}</h6>
                <p class="feature-desc mb-0">{{ t('dashboard.featureAIDesc') }}</p>
              </div>
            </div>
            <div class="col-md-4">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-share-alt"></i></div>
                <h6>{{ t('dashboard.featureShare') }}</h6>
                <p class="feature-desc mb-0">{{ t('dashboard.featureShareDesc') }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No results for filter -->
      <div v-else-if="filteredCVs.length === 0" class="text-center py-5">
        <div class="empty-filter-icon mb-3">
          <i class="fas fa-filter fa-3x" style="color: #6a6a6a;"></i>
        </div>
        <h5 style="color: #6a6a6a;">{{ t('dashboard.noResults') }}</h5>
        <button class="filter-pill active mt-2" @click="activeFilter = 'all'; searchQuery = ''">
          {{ t('dashboard.clearFilters') }}
        </button>
      </div>

      <!-- CV Cards Grid -->
      <div v-else class="row g-4">
        <div class="col-sm-6 col-lg-4" v-for="cv in filteredCVs" :key="cv.id">
          <div class="cv-card h-100" @click="$router.push(`/cv/${cv.id}`)">
            <!-- Template Preview Thumbnail -->
            <div class="cv-card-thumbnail" :class="`thumb-${cv.template}`">
              <div class="thumb-overlay">
                <button class="btn-preview-overlay" @click.stop="$router.push(`/cv/${cv.id}`)">
                  <i class="fas fa-eye me-1"></i>{{ t('app.preview') }}
                </button>
              </div>
              <div class="thumb-content">
                <div class="thumb-icon">
                  <i class="fas fa-file-alt"></i>
                </div>
                <div class="thumb-template-name">{{ t(`templates.${cv.template}`) }}</div>
              </div>
              <span class="status-dot-label" :class="`status-${cv.status}`">
                <span class="status-dot"></span>
                {{ t(`cv.${cv.status}`) }}
              </span>
            </div>

            <div class="cv-card-body">
              <h5 class="cv-card-title mb-2">{{ cv.title || t('cv.title') }}</h5>
              <div class="cv-card-meta">
                <span class="meta-item" :title="t('cv.language')">
                  <i class="fas fa-globe me-1"></i>{{ cv.language === 'ar' ? t('app.arabic') : t('app.english') }}
                </span>
                <span class="meta-item" :title="t('cv.createdAt')">
                  <i class="fas fa-calendar-alt me-1"></i>{{ formatDate(cv.created_at) }}
                </span>
              </div>
              <div v-if="cv.is_shared" class="shared-indicator mt-2">
                <i class="fas fa-eye me-1"></i>{{ cv.view_count }} {{ t('cv.viewCount') }}
              </div>
            </div>

            <div class="cv-card-footer">
              <div class="cv-card-actions d-flex gap-2">
                <router-link :to="`/cv/${cv.id}`" class="btn-action btn-action-view flex-fill" @click.stop :title="t('app.view')">
                  <i class="fas fa-eye"></i>
                  <span class="d-none d-xl-inline ms-1">{{ t('app.view') }}</span>
                </router-link>
                <router-link :to="`/cv/${cv.id}/edit`" class="btn-action btn-action-edit flex-fill" @click.stop :title="t('app.edit')">
                  <i class="fas fa-edit"></i>
                  <span class="d-none d-xl-inline ms-1">{{ t('app.edit') }}</span>
                </router-link>
                <button @click.stop="handleShare(cv)" class="btn-action btn-action-share" :title="t('app.share')">
                  <i class="fas fa-share-alt"></i>
                </button>
                <button @click.stop="handleDelete(cv.id)" class="btn-action btn-action-delete" :title="t('app.delete')">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Floating Create Button (Mobile) -->
      <router-link to="/cv/create" class="fab-create d-md-none" :title="t('app.createCV')">
        <i class="fas fa-plus"></i>
      </router-link>
    </div>

    <!-- Share Modal -->
    <div class="modal fade" ref="shareModal" tabindex="-1">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content share-modal-content">
          <div class="modal-header border-0 pb-0">
            <h5 class="modal-title" style="font-weight: 700; color: #222222;">
              <i class="fas fa-share-alt me-2" style="color: #1a5276;"></i>{{ t('cv.shareCV') }}
            </h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body" v-if="selectedCV">
            <div class="share-toggle-section mb-4 text-center">
              <button
                class="btn btn-lg px-4"
                :class="selectedCV.is_shared ? 'btn-danger' : 'btn-success'"
                :disabled="sharingLoading"
                @click="toggleShare"
                style="border-radius: 8px; font-weight: 600;"
              >
                <span v-if="sharingLoading" class="spinner-border spinner-border-sm me-2"></span>
                <i v-else :class="selectedCV.is_shared ? 'fas fa-lock' : 'fas fa-share-alt'" class="me-2"></i>
                {{ selectedCV.is_shared ? t('cv.disableShare') : t('cv.enableShare') }}
              </button>
            </div>

            <div v-if="selectedCV.is_shared">
              <div class="mb-3">
                <label class="form-label" style="color: #6a6a6a; font-size: 0.75rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em;">{{ t('cv.shareLink') }}</label>
                <div class="input-group input-group-lg">
                  <input type="text" class="form-control share-link-input" :value="shareUrl" readonly />
                  <button @click="copyShareLink" class="btn" style="background: #1a5276; color: white; border-radius: 0 8px 8px 0; font-weight: 600;">
                    <i :class="copied ? 'fas fa-check' : 'fas fa-copy'" class="me-1"></i>
                    {{ copied ? t('dashboard.copied') : t('cv.copyLink') }}
                  </button>
                </div>
              </div>

              <div class="text-center mb-3" v-if="qrCodeImage">
                <label class="form-label" style="color: #6a6a6a; font-size: 0.75rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; display: block;">{{ t('cv.qrCode') }}</label>
                <div class="qr-code-wrapper">
                  <img :src="qrCodeImage" alt="QR Code" class="img-fluid" />
                </div>
              </div>

              <div class="view-count-section text-center">
                <span class="view-count-badge">
                  <i class="fas fa-eye me-1"></i>{{ selectedCV.view_count }} {{ t('cv.viewCount') }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCVStore, type CV } from '../stores/cv'
import { useAuthStore } from '../stores/auth'
import { cvAPI } from '../services/api'
import QRCode from 'qrcode'

const { t, locale } = useI18n()
const cvStore = useCVStore()
const authStore = useAuthStore()

const shareModal = ref<HTMLElement>()
const selectedCV = ref<CV | null>(null)
const searchQuery = ref('')
const activeFilter = ref('all')
const copied = ref(false)
const qrCodeImage = ref('')
const sharingLoading = ref(false)
let bsModal: any = null

const userName = computed(() => {
  return authStore.user?.full_name_en || authStore.user?.full_name_ar || authStore.user?.email?.split('@')[0] || ''
})

const totalCount = computed(() => cvStore.cvs.length)
const approvedCount = computed(() => cvStore.cvs.filter(cv => cv.status === 'approved').length)
const pendingCount = computed(() => cvStore.cvs.filter(cv => cv.status === 'pending').length)

const statusFilters = computed(() => [
  { value: 'all', label: t('dashboard.all'), icon: 'fas fa-th', count: totalCount.value },
  { value: 'draft', label: t('cv.draft'), icon: 'fas fa-pencil-alt', count: cvStore.cvs.filter(cv => cv.status === 'draft').length },
  { value: 'pending', label: t('cv.pending'), icon: 'fas fa-clock', count: pendingCount.value },
  { value: 'approved', label: t('cv.approved'), icon: 'fas fa-check-circle', count: approvedCount.value },
  { value: 'rejected', label: t('cv.rejected'), icon: 'fas fa-times-circle', count: cvStore.cvs.filter(cv => cv.status === 'rejected').length },
])

const filteredCVs = computed(() => {
  let result = cvStore.cvs
  if (activeFilter.value !== 'all') {
    result = result.filter(cv => cv.status === activeFilter.value)
  }
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    result = result.filter(cv =>
      (cv.title || '').toLowerCase().includes(q) ||
      (cv.data?.personal_info?.full_name || '').toLowerCase().includes(q)
    )
  }
  return result
})

const shareUrl = computed(() => {
  if (!selectedCV.value) return ''
  return `${window.location.origin}/shared/${selectedCV.value.share_token}`
})

// Generate QR code when share URL changes
watch(shareUrl, async (url) => {
  if (url) {
    try {
      qrCodeImage.value = await QRCode.toDataURL(url, { width: 200, margin: 2 })
    } catch {
      qrCodeImage.value = ''
    }
  } else {
    qrCodeImage.value = ''
  }
})

onMounted(() => {
  cvStore.fetchCVs()
})


function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}

async function handleShare(cv: CV) {
  selectedCV.value = cv
  copied.value = false
  if (shareModal.value) {
    const { Modal } = await import('bootstrap')
    bsModal = new Modal(shareModal.value)
    bsModal.show()
  }
}

function copyShareLink() {
  navigator.clipboard.writeText(shareUrl.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}

async function toggleShare() {
  if (!selectedCV.value || sharingLoading.value) return
  sharingLoading.value = true
  try {
    const updated = await cvStore.toggleShare(selectedCV.value.id)
    selectedCV.value = { ...selectedCV.value, is_shared: updated.is_shared, share_token: updated.share_token }
    cvStore.fetchCVs()
  } catch {
    // Error handling
  } finally {
    sharingLoading.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm(t('cv.deleteConfirm'))) return
  await cvStore.deleteCV(id)
}

async function exportJSON() {
  try {
    const res = await cvAPI.exportJSON()
    const blob = new Blob([res.data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'my-cvs.json'
    a.click()
    URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Export failed:', err)
    alert(t('app.error') || 'Export failed')
  }
}
</script>

<style scoped>
/* === Page === */
.dashboard-page {
  min-height: 100vh;
  background: #ffffff;
}

/* === Welcome Section === */
.welcome-section {
  padding: 0;
}

.welcome-date {
  font-size: 0.88rem;
  color: #6a6a6a;
  font-weight: 500;
}

.welcome-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #222222;
}

.welcome-name {
  color: #222222;
}

.welcome-subtitle {
  font-size: 0.95rem;
  color: #6a6a6a;
  margin-top: 4px;
  font-weight: 500;
}

/* === Stats Cards === */
.stat-card {
  background: #ffffff;
  border-radius: 12px;
  padding: 24px 16px;
  text-align: center;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
  transition: box-shadow 0.2s ease;
}

.stat-card:hover {
  box-shadow: rgba(0, 0, 0, 0.08) 0px 4px 12px;
}

.stat-icon-circle {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
  margin: 0 auto 12px;
}

.stat-icon-blue {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.stat-icon-green {
  background: rgba(30, 132, 73, 0.1);
  color: #1e8449;
}

.stat-icon-gold {
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
}

.stat-icon-red {
  background: rgba(192, 57, 43, 0.1);
  color: #c0392b;
}

.stat-number {
  font-size: 1.75rem;
  font-weight: 700;
  color: #222222;
  line-height: 1.2;
}

.stat-label {
  font-size: 0.78rem;
  color: #6a6a6a;
  font-weight: 500;
  margin-top: 4px;
}

/* === Quick Actions === */
.quick-action-btn {
  display: inline-flex;
  align-items: center;
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 600;
  border: 1px solid #222222;
  background: #222222;
  color: #ffffff;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-action-btn:hover {
  background: #000000;
  color: #ffffff;
}

.quick-action-btn--secondary {
  background: #ffffff;
  color: #222222;
  border: 1px solid #dddddd;
}

.quick-action-btn--secondary:hover {
  background: #f7f7f7;
  border-color: #222222;
  color: #222222;
}

/* === Toolbar === */
.toolbar-section {
  padding: 0;
}

.search-box {
  position: relative;
}

.search-icon {
  position: absolute;
  top: 50%;
  left: 14px;
  transform: translateY(-50%);
  color: #6a6a6a;
  font-size: 0.9rem;
  z-index: 2;
}

[dir="rtl"] .search-icon {
  left: auto;
  right: 14px;
}

.search-input {
  padding-left: 40px;
  border-radius: 40px;
  border: 1px solid #dddddd;
  font-size: 0.9rem;
  background: #ffffff;
  height: 44px;
  color: #222222;
  transition: border-color 0.2s, box-shadow 0.2s;
}

[dir="rtl"] .search-input {
  padding-left: 12px;
  padding-right: 40px;
}

.search-input:focus {
  border-color: #222222;
  box-shadow: none;
}

.search-clear {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  border: none;
  background: none;
  color: #6a6a6a;
  cursor: pointer;
  padding: 4px;
  z-index: 2;
}

[dir="rtl"] .search-clear {
  right: auto;
  left: 12px;
}

.search-clear:hover { color: #222222; }

/* Filter Pills */
.filter-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px solid #dddddd;
  background: #ffffff;
  color: #6a6a6a;
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.filter-pill:hover {
  border-color: #222222;
  color: #222222;
}

.filter-pill.active {
  background: #222222;
  color: #ffffff;
  border-color: #222222;
}

.filter-count {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 10px;
  padding: 1px 7px;
  font-size: 0.72rem;
  font-weight: 600;
}

.filter-pill:not(.active) .filter-count {
  background: #f0f0f0;
  color: #6a6a6a;
}

.btn-create {
  font-weight: 600;
  border-radius: 8px;
  padding: 8px 20px;
  background: #1a5276;
  border-color: #1a5276;
  display: inline-flex;
  align-items: center;
  font-size: 0.88rem;
}

.btn-create:hover {
  background: #154360;
  border-color: #154360;
}

/* === Loading Skeleton === */
.cv-card-skeleton {
  border-radius: 12px;
  overflow: hidden;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.skeleton-thumbnail {
  height: 140px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite;
}

.skeleton-body {
  padding: 16px;
}

.skeleton-line {
  height: 14px;
  border-radius: 6px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite;
  margin-bottom: 10px;
}

.skeleton-title { width: 70%; height: 18px; }
.skeleton-badge { width: 35%; height: 24px; border-radius: 12px; }
.skeleton-text { width: 90%; }
.skeleton-text.short { width: 55%; }

@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* === Empty State === */
.empty-state {
  max-width: 700px;
  margin: 0 auto;
}

.empty-icon-circle {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: #f7f7f7;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: #6a6a6a;
}

.empty-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #222222;
}

.empty-description {
  max-width: 400px;
  font-size: 1rem;
  line-height: 1.6;
  color: #6a6a6a;
}

.btn-create-first {
  display: inline-flex;
  align-items: center;
  background: #1a5276;
  color: #ffffff;
  border-radius: 8px;
  padding: 12px 32px;
  font-weight: 600;
  font-size: 1.05rem;
  text-decoration: none;
  transition: all 0.2s ease;
}

.btn-create-first:hover {
  background: #154360;
  color: #ffffff;
  box-shadow: rgba(0, 0, 0, 0.08) 0px 4px 12px;
}

.feature-item {
  text-align: center;
  padding: 20px;
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #f7f7f7;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  font-size: 1.2rem;
  color: #1a5276;
}

.feature-item h6 {
  font-weight: 600;
  color: #222222;
}

.feature-desc {
  color: #6a6a6a;
  font-size: 0.88rem;
}

/* === CV Cards === */
.cv-card {
  background: #ffffff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
  transition: box-shadow 0.25s ease, transform 0.25s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.cv-card:hover {
  transform: translateY(-4px);
  box-shadow: rgba(0, 0, 0, 0.08) 0px 4px 12px;
}

/* Thumbnail */
.cv-card-thumbnail {
  height: 140px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: #f7f7f7;
}

.thumb-modern { background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); }
.thumb-professional { background: linear-gradient(135deg, #1a237e 0%, #3949ab 100%); }
.thumb-minimalist { background: #f0f0f0; }
.thumb-minimalist .thumb-content { color: #222222; }
.thumb-academic { background: linear-gradient(135deg, #3e2723 0%, #795548 100%); }
.thumb-creative { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.thumb-ats { background: linear-gradient(135deg, #37474f 0%, #546e7a 100%); }
.thumb-compact { background: linear-gradient(135deg, #455a64 0%, #78909c 100%); }
.thumb-elegant { background: linear-gradient(135deg, #5d4037 0%, #b8860b 100%); }
.thumb-executive { background: linear-gradient(135deg, #1b2838 0%, #2c3e50 100%); }
.thumb-tech { background: linear-gradient(135deg, #0d1117 0%, #161b22 100%); }
.thumb-designer { background: linear-gradient(45deg, #ff6b6b 0%, #feca57 100%); }
.thumb-medical { background: linear-gradient(135deg, #1b5e20 0%, #2e7d32 100%); }
.thumb-engineer { background: linear-gradient(135deg, #0d47a1 0%, #1565c0 100%); }
.thumb-standard { background: linear-gradient(135deg, #424242 0%, #757575 100%); }
.thumb-teacher { background: linear-gradient(135deg, #e65100 0%, #f57c00 100%); }
.thumb-lawyer { background: linear-gradient(135deg, #212121 0%, #424242 100%); }

.thumb-content {
  text-align: center;
  color: white;
  z-index: 1;
}

.thumb-icon {
  font-size: 2.4rem;
  margin-bottom: 4px;
  opacity: 0.9;
}

.thumb-template-name {
  font-size: 0.8rem;
  opacity: 0.85;
  font-weight: 500;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.thumb-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.25s ease;
  z-index: 2;
}

.cv-card:hover .thumb-overlay {
  opacity: 1;
}

.btn-preview-overlay {
  background: #ffffff;
  color: #222222;
  border: none;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
}

/* Status Dot + Label */
.status-dot-label {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.72rem;
  font-weight: 600;
  z-index: 3;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(4px);
  color: #222222;
}

[dir="rtl"] .status-dot-label {
  right: auto;
  left: 10px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-draft .status-dot { background: #6a6a6a; }
.status-pending .status-dot { background: #c0982b; }
.status-approved .status-dot { background: #1e8449; }
.status-rejected .status-dot { background: #c0392b; }

/* Card Body */
.cv-card-body {
  padding: 16px 16px 8px;
  flex: 1;
}

.cv-card-title {
  font-size: 1.05rem;
  font-weight: 600;
  color: #222222;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cv-card-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.meta-item {
  font-size: 0.82rem;
  color: #6a6a6a;
}

.meta-item i {
  color: #6a6a6a;
}

.shared-indicator {
  display: inline-flex;
  align-items: center;
  font-size: 0.8rem;
  color: #1a5276;
  background: rgba(26, 82, 118, 0.08);
  padding: 3px 10px;
  border-radius: 8px;
}

/* Card Footer */
.cv-card-footer {
  padding: 0 16px 16px;
}

/* Card Actions */
.btn-action {
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 500;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px 10px;
  border: 1px solid #dddddd;
  background: #ffffff;
  color: #222222;
  text-decoration: none;
  cursor: pointer;
}

.btn-action:hover {
  border-color: #222222;
  color: #222222;
}

.btn-action-view:hover {
  background: #1a5276;
  border-color: #1a5276;
  color: #ffffff;
}

.btn-action-edit:hover {
  background: #c0982b;
  border-color: #c0982b;
  color: #ffffff;
}

.btn-action-share:hover {
  background: #1a5276;
  border-color: #1a5276;
  color: #ffffff;
}

.btn-action-delete:hover {
  background: #c0392b;
  border-color: #c0392b;
  color: #ffffff;
}

/* === Floating Action Button === */
.fab-create {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #1a5276;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
  text-decoration: none;
  transition: transform 0.2s, box-shadow 0.2s;
  z-index: 1000;
}

[dir="rtl"] .fab-create {
  right: auto;
  left: 24px;
}

.fab-create:hover {
  transform: scale(1.1);
  box-shadow: rgba(0, 0, 0, 0.08) 0px 4px 12px;
  color: white;
}

/* === Share Modal === */
.share-modal-content {
  border: none;
  border-radius: 12px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.share-link-input {
  border-radius: 8px 0 0 8px;
  background: #f7f7f7;
  font-size: 0.9rem;
  border: 1px solid #dddddd;
  color: #222222;
}

[dir="rtl"] .share-link-input {
  border-radius: 0 8px 8px 0;
}

.qr-code-wrapper {
  display: inline-block;
  padding: 12px;
  background: white;
  border-radius: 12px;
  border: 1px solid #dddddd;
}

.qr-code-wrapper img {
  max-width: 180px;
}

.view-count-badge {
  display: inline-flex;
  align-items: center;
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.88rem;
  font-weight: 500;
}

/* === Responsive === */
@media (max-width: 768px) {
  .welcome-title {
    font-size: 1.3rem;
  }

  .stat-card {
    padding: 16px 8px;
  }

  .stat-number {
    font-size: 1.3rem;
  }

  .stat-label {
    font-size: 0.65rem;
  }

  .filter-pills {
    overflow-x: auto;
    flex-wrap: nowrap !important;
    padding-bottom: 4px;
  }

  .cv-card-thumbnail {
    height: 120px;
  }
}
</style>
