<template>
  <div class="dashboard-page">
    <div class="container py-4">
      <!-- Welcome Header -->
      <div class="welcome-section mb-4">
        <div class="row align-items-center">
          <div class="col-lg-6">
            <div class="d-flex align-items-center mb-2">
              <div class="welcome-avatar me-3">
                <i class="fas fa-user-graduate"></i>
              </div>
              <div>
                <p class="welcome-date mb-1">{{ new Date().toLocaleDateString(locale === 'ar' ? 'ar-IQ' : 'en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }) }}</p>
                <h2 class="welcome-title mb-0">
                  {{ locale === 'ar' ? 'مرحباً' : 'Welcome' }}<span v-if="userName" class="welcome-name">, {{ userName }}</span>
                </h2>
                <p class="welcome-subtitle mb-0">{{ t('dashboard.subtitle') }}</p>
              </div>
            </div>
          </div>
          <div class="col-lg-6">
            <div class="row g-3 mt-2 mt-lg-0">
              <div class="col-6 col-sm-3">
                <div class="stat-card stat-total">
                  <div class="stat-icon-wrap"><i class="fas fa-file-alt"></i></div>
                  <div class="stat-number">{{ totalCount }}</div>
                  <div class="stat-label">{{ t('dashboard.totalCVs') }}</div>
                </div>
              </div>
              <div class="col-6 col-sm-3">
                <div class="stat-card stat-approved">
                  <div class="stat-icon-wrap"><i class="fas fa-check-circle"></i></div>
                  <div class="stat-number">{{ approvedCount }}</div>
                  <div class="stat-label">{{ t('cv.approved') }}</div>
                </div>
              </div>
              <div class="col-6 col-sm-3">
                <div class="stat-card stat-pending">
                  <div class="stat-icon-wrap"><i class="fas fa-clock"></i></div>
                  <div class="stat-number">{{ pendingCount }}</div>
                  <div class="stat-label">{{ t('cv.pending') }}</div>
                </div>
              </div>
              <div class="col-6 col-sm-3">
                <div class="stat-card stat-rejected">
                  <div class="stat-icon-wrap"><i class="fas fa-times-circle"></i></div>
                  <div class="stat-number">{{ cvStore.cvs.filter(cv => cv.status === 'rejected').length }}</div>
                  <div class="stat-label">{{ t('cv.rejected') }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="quick-actions mb-4">
        <div class="row g-3">
          <div class="col-sm-4">
            <router-link to="/cv/create" class="quick-action-card">
              <div class="qa-icon"><i class="fas fa-plus-circle"></i></div>
              <div class="qa-text">
                <h6 class="mb-0">{{ t('app.createCV') }}</h6>
                <small>{{ locale === 'ar' ? 'ابدأ سيرتك الذاتية الجديدة' : 'Start a new CV from scratch' }}</small>
              </div>
              <i class="fas fa-arrow-right qa-arrow"></i>
            </router-link>
          </div>
          <div class="col-sm-4">
            <div class="quick-action-card" @click="exportJSON" role="button">
              <div class="qa-icon qa-icon-export"><i class="fas fa-download"></i></div>
              <div class="qa-text">
                <h6 class="mb-0">{{ t('admin.exportJSON') }}</h6>
                <small>{{ locale === 'ar' ? 'تصدير السير الذاتية' : 'Export your CVs as JSON' }}</small>
              </div>
              <i class="fas fa-arrow-right qa-arrow"></i>
            </div>
          </div>
          <div class="col-sm-4">
            <router-link to="/cv/create" class="quick-action-card">
              <div class="qa-icon qa-icon-template"><i class="fas fa-palette"></i></div>
              <div class="qa-text">
                <h6 class="mb-0">{{ locale === 'ar' ? 'تصفح القوالب' : 'Browse Templates' }}</h6>
                <small>{{ locale === 'ar' ? 'اختر من القوالب الجاهزة' : 'Choose from ready templates' }}</small>
              </div>
              <i class="fas fa-arrow-right qa-arrow"></i>
            </router-link>
          </div>
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
          <div class="card cv-card-skeleton">
            <div class="skeleton-thumbnail"></div>
            <div class="card-body">
              <div class="skeleton-line skeleton-title"></div>
              <div class="skeleton-line skeleton-badge"></div>
              <div class="skeleton-line skeleton-text"></div>
              <div class="skeleton-line skeleton-text short"></div>
            </div>
            <div class="card-footer bg-transparent border-0">
              <div class="d-flex gap-2">
                <div class="skeleton-btn flex-fill"></div>
                <div class="skeleton-btn flex-fill"></div>
              </div>
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
        <p class="empty-description text-muted mb-4 mx-auto">{{ t('dashboard.emptyDescription') }}</p>
        <router-link to="/cv/create" class="btn btn-primary btn-lg btn-create-empty">
          <i class="fas fa-plus-circle me-2"></i>{{ t('dashboard.createFirst') }}
        </router-link>
        <div class="empty-features mt-5">
          <div class="row g-4 justify-content-center">
            <div class="col-md-4">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-palette"></i></div>
                <h6>{{ t('dashboard.featureTemplates') }}</h6>
                <p class="text-muted small mb-0">{{ t('dashboard.featureTemplatesDesc') }}</p>
              </div>
            </div>
            <div class="col-md-4">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-robot"></i></div>
                <h6>{{ t('dashboard.featureAI') }}</h6>
                <p class="text-muted small mb-0">{{ t('dashboard.featureAIDesc') }}</p>
              </div>
            </div>
            <div class="col-md-4">
              <div class="feature-item">
                <div class="feature-icon"><i class="fas fa-share-alt"></i></div>
                <h6>{{ t('dashboard.featureShare') }}</h6>
                <p class="text-muted small mb-0">{{ t('dashboard.featureShareDesc') }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- No results for filter -->
      <div v-else-if="filteredCVs.length === 0" class="text-center py-5">
        <div class="empty-filter-icon mb-3">
          <i class="fas fa-filter fa-3x text-muted"></i>
        </div>
        <h5 class="text-muted">{{ t('dashboard.noResults') }}</h5>
        <button class="btn btn-outline-primary mt-2" @click="activeFilter = 'all'; searchQuery = ''">
          {{ t('dashboard.clearFilters') }}
        </button>
      </div>

      <!-- CV Cards Grid -->
      <div v-else class="row g-4">
        <div class="col-sm-6 col-lg-4" v-for="cv in filteredCVs" :key="cv.id">
          <div class="card cv-card h-100" @click="$router.push(`/cv/${cv.id}`)">
            <!-- Template Preview Thumbnail -->
            <div class="cv-card-thumbnail" :class="`thumb-${cv.template}`">
              <div class="thumb-overlay">
                <button class="btn btn-light btn-sm rounded-pill" @click.stop="$router.push(`/cv/${cv.id}`)">
                  <i class="fas fa-eye me-1"></i>{{ t('app.preview') }}
                </button>
              </div>
              <div class="thumb-content">
                <div class="thumb-icon">
                  <i class="fas fa-file-alt"></i>
                </div>
                <div class="thumb-template-name">{{ t(`templates.${cv.template}`) }}</div>
              </div>
              <span class="status-pill" :class="`status-${cv.status}`">
                <i :class="statusIcon(cv.status)" class="me-1"></i>{{ t(`cv.${cv.status}`) }}
              </span>
            </div>

            <div class="card-body pb-2">
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

            <div class="card-footer bg-transparent border-top-0 pt-0 pb-3">
              <div class="cv-card-actions d-flex gap-2">
                <router-link :to="`/cv/${cv.id}`" class="btn btn-sm btn-action btn-action-view flex-fill" @click.stop :title="t('app.view')">
                  <i class="fas fa-eye"></i>
                  <span class="d-none d-xl-inline ms-1">{{ t('app.view') }}</span>
                </router-link>
                <router-link :to="`/cv/${cv.id}/edit`" class="btn btn-sm btn-action btn-action-edit flex-fill" @click.stop :title="t('app.edit')">
                  <i class="fas fa-edit"></i>
                  <span class="d-none d-xl-inline ms-1">{{ t('app.edit') }}</span>
                </router-link>
                <button @click.stop="handleShare(cv)" class="btn btn-sm btn-action btn-action-share" :title="t('app.share')">
                  <i class="fas fa-share-alt"></i>
                </button>
                <button @click.stop="handleDelete(cv.id)" class="btn btn-sm btn-action btn-action-delete" :title="t('app.delete')">
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
            <h5 class="modal-title fw-bold">
              <i class="fas fa-share-alt text-primary me-2"></i>{{ t('cv.shareCV') }}
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
              >
                <span v-if="sharingLoading" class="spinner-border spinner-border-sm me-2"></span>
                <i v-else :class="selectedCV.is_shared ? 'fas fa-lock' : 'fas fa-share-alt'" class="me-2"></i>
                {{ selectedCV.is_shared ? t('cv.disableShare') : t('cv.enableShare') }}
              </button>
            </div>

            <div v-if="selectedCV.is_shared">
              <div class="mb-3">
                <label class="form-label text-muted small text-uppercase fw-bold">{{ t('cv.shareLink') }}</label>
                <div class="input-group input-group-lg">
                  <input type="text" class="form-control share-link-input" :value="shareUrl" readonly />
                  <button @click="copyShareLink" class="btn btn-primary">
                    <i :class="copied ? 'fas fa-check' : 'fas fa-copy'" class="me-1"></i>
                    {{ copied ? t('dashboard.copied') : t('cv.copyLink') }}
                  </button>
                </div>
              </div>

              <div class="text-center mb-3" v-if="qrCodeImage">
                <label class="form-label text-muted small text-uppercase fw-bold d-block">{{ t('cv.qrCode') }}</label>
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

function statusIcon(status: string) {
  const map: Record<string, string> = {
    draft: 'fas fa-pencil-alt',
    pending: 'fas fa-clock',
    approved: 'fas fa-check-circle',
    rejected: 'fas fa-times-circle',
  }
  return map[status] || 'fas fa-circle'
}

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
/* === Color Variables === */
.dashboard-page {
  --uni-primary: #1a5276;
  --uni-accent: #c0982b;
  --uni-secondary: #2c3e50;
  --uni-primary-light: #e8f0f7;
  --uni-accent-light: #fdf5e6;
  --uni-success: #1e8449;
  --uni-warning: #d4a017;
  --uni-danger: #c0392b;
  --uni-bg: #f5f7fa;
  min-height: 100vh;
  background: var(--uni-bg);
}

/* === Welcome Section === */
.welcome-section {
  background: linear-gradient(135deg, var(--uni-primary) 0%, var(--uni-secondary) 100%);
  border-radius: 16px;
  padding: 28px 32px;
  color: white;
  position: relative;
  overflow: hidden;
}

.welcome-section::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -10%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(192, 152, 43, 0.15) 0%, transparent 70%);
  border-radius: 50%;
}

[dir="rtl"] .welcome-section::before {
  right: auto;
  left: -10%;
}

.welcome-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  border: 2px solid var(--uni-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.6rem;
  color: var(--uni-accent);
  flex-shrink: 0;
}

.welcome-date {
  font-size: 0.82rem;
  opacity: 0.75;
  letter-spacing: 0.3px;
}

.welcome-title {
  font-size: 1.6rem;
  font-weight: 700;
  color: white;
}

.welcome-name {
  color: var(--uni-accent);
}

.welcome-subtitle {
  font-size: 0.9rem;
  opacity: 0.8;
  margin-top: 2px;
}

/* Stats Cards */
.stat-card {
  text-align: center;
  background: rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(8px);
  border-radius: 12px;
  padding: 16px 8px;
  border: 1px solid rgba(255, 255, 255, 0.15);
  transition: transform 0.2s, background 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.18);
}

.stat-icon-wrap {
  font-size: 1.1rem;
  margin-bottom: 6px;
  opacity: 0.8;
}

.stat-total .stat-icon-wrap { color: #93c5fd; }
.stat-approved .stat-icon-wrap { color: #86efac; }
.stat-pending .stat-icon-wrap { color: #fde68a; }
.stat-rejected .stat-icon-wrap { color: #fca5a5; }

.stat-number {
  font-size: 1.7rem;
  font-weight: 800;
  line-height: 1.2;
  color: white;
}

.stat-label {
  font-size: 0.7rem;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-top: 2px;
}

/* === Quick Actions === */
.quick-action-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  background: white;
  border-radius: 14px;
  border: 1px solid #e4e8f0;
  text-decoration: none;
  color: var(--uni-secondary);
  transition: all 0.25s ease;
  cursor: pointer;
  height: 100%;
}

.quick-action-card:hover {
  border-color: var(--uni-accent);
  box-shadow: 0 6px 20px rgba(26, 82, 118, 0.1);
  transform: translateY(-2px);
  color: var(--uni-secondary);
}

.qa-icon {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--uni-primary) 0%, #2471a3 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.15rem;
  flex-shrink: 0;
}

.qa-icon-export {
  background: linear-gradient(135deg, var(--uni-accent) 0%, #d4a017 100%);
}

.qa-icon-template {
  background: linear-gradient(135deg, var(--uni-secondary) 0%, #34495e 100%);
}

.qa-text {
  flex: 1;
  min-width: 0;
}

.qa-text h6 {
  font-weight: 600;
  font-size: 0.92rem;
  color: var(--uni-secondary);
}

.qa-text small {
  color: #7f8c8d;
  font-size: 0.78rem;
}

.qa-arrow {
  color: #bdc3c7;
  font-size: 0.85rem;
  transition: transform 0.2s;
}

[dir="rtl"] .qa-arrow {
  transform: scaleX(-1);
}

.quick-action-card:hover .qa-arrow {
  color: var(--uni-accent);
  transform: translateX(3px);
}

[dir="rtl"] .quick-action-card:hover .qa-arrow {
  transform: scaleX(-1) translateX(3px);
}

/* === Toolbar === */
.toolbar-section {
  background: white;
  border-radius: 14px;
  padding: 16px 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid #e4e8f0;
}

.search-box {
  position: relative;
}

.search-icon {
  position: absolute;
  top: 50%;
  left: 14px;
  transform: translateY(-50%);
  color: #94a3b8;
  font-size: 0.9rem;
  z-index: 2;
}

[dir="rtl"] .search-icon {
  left: auto;
  right: 14px;
}

.search-input {
  padding-left: 40px;
  border-radius: 24px;
  border: 1px solid #e2e8f0;
  font-size: 0.9rem;
  background: #f8fafc;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
  height: 42px;
}

[dir="rtl"] .search-input {
  padding-left: 12px;
  padding-right: 40px;
}

.search-input:focus {
  background: white;
  border-color: var(--uni-primary);
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

.search-clear {
  position: absolute;
  top: 50%;
  right: 12px;
  transform: translateY(-50%);
  border: none;
  background: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  z-index: 2;
}

[dir="rtl"] .search-clear {
  right: auto;
  left: 12px;
}

.search-clear:hover { color: #64748b; }

/* Filter Pills */
.filter-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #64748b;
  font-size: 0.82rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.filter-pill:hover {
  border-color: var(--uni-primary);
  color: var(--uni-primary);
  background: var(--uni-primary-light);
}

.filter-pill.active {
  background: var(--uni-primary);
  color: white;
  border-color: var(--uni-primary);
}

.filter-count {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 10px;
  padding: 1px 7px;
  font-size: 0.72rem;
  font-weight: 600;
}

.filter-pill:not(.active) .filter-count {
  background: #e2e8f0;
  color: #475569;
}

.btn-create {
  font-weight: 600;
  border-radius: 10px;
  padding: 8px 20px;
  background: var(--uni-primary);
  border-color: var(--uni-primary);
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
  border: none;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
}

.skeleton-thumbnail {
  height: 140px;
  background: linear-gradient(90deg, #eef2f7 25%, #e4e9f0 50%, #eef2f7 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite;
}

.skeleton-line {
  height: 14px;
  border-radius: 6px;
  background: linear-gradient(90deg, #eef2f7 25%, #e4e9f0 50%, #eef2f7 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite;
  margin-bottom: 10px;
}

.skeleton-title { width: 70%; height: 18px; }
.skeleton-badge { width: 35%; height: 24px; border-radius: 12px; }
.skeleton-text { width: 90%; }
.skeleton-text.short { width: 55%; }
.skeleton-btn {
  height: 32px;
  border-radius: 6px;
  background: linear-gradient(90deg, #eef2f7 25%, #e4e9f0 50%, #eef2f7 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite;
}

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
  width: 140px;
  height: 140px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--uni-primary-light) 0%, #d5e8f7 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: var(--uni-primary);
  animation: empty-float 3s ease-in-out infinite;
}

@keyframes empty-float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.empty-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--uni-secondary);
}

.empty-description {
  max-width: 400px;
  font-size: 1rem;
  line-height: 1.6;
}

.btn-create-empty {
  background: var(--uni-primary);
  border-color: var(--uni-primary);
  border-radius: 12px;
  padding: 12px 32px;
  font-weight: 600;
  font-size: 1.05rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-create-empty:hover {
  background: #154360;
  border-color: #154360;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(26, 82, 118, 0.3);
}

.feature-item {
  text-align: center;
  padding: 20px;
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: var(--uni-primary-light);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  font-size: 1.2rem;
  color: var(--uni-primary);
}

/* === CV Cards === */
.cv-card {
  border: none;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(26, 82, 118, 0.06), 0 1px 2px rgba(0, 0, 0, 0.04);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  cursor: pointer;
}

.cv-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 28px rgba(26, 82, 118, 0.12), 0 4px 10px rgba(0, 0, 0, 0.06);
}

/* Thumbnail Backgrounds by Template */
.cv-card-thumbnail {
  height: 140px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: linear-gradient(135deg, var(--uni-primary) 0%, var(--uni-secondary) 100%);
}

.thumb-modern { background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); }
.thumb-professional { background: linear-gradient(135deg, #1a237e 0%, #3949ab 100%); }
.thumb-minimalist { background: linear-gradient(135deg, #f5f5f5 0%, #e0e0e0 100%); }
.thumb-minimalist .thumb-content { color: #333; }
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
  background: rgba(26, 82, 118, 0.6);
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

/* Status Pills */
.status-pill {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  z-index: 3;
  backdrop-filter: blur(4px);
}

[dir="rtl"] .status-pill {
  right: auto;
  left: 10px;
}

.status-draft {
  background: rgba(100, 116, 139, 0.9);
  color: white;
}

.status-pending {
  background: rgba(212, 160, 23, 0.95);
  color: #1e293b;
}

.status-approved {
  background: rgba(30, 132, 73, 0.9);
  color: white;
}

.status-rejected {
  background: rgba(192, 57, 43, 0.9);
  color: white;
}

/* Card Body */
.cv-card-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--uni-secondary);
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
  color: #64748b;
}

.meta-item i {
  color: #94a3b8;
}

.shared-indicator {
  display: inline-flex;
  align-items: center;
  font-size: 0.8rem;
  color: var(--uni-primary);
  background: var(--uni-primary-light);
  padding: 3px 10px;
  border-radius: 8px;
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
}

.btn-action:hover {
  transform: translateY(-1px);
}

.btn-action-view {
  background: transparent;
  border: 1px solid var(--uni-primary);
  color: var(--uni-primary);
}

.btn-action-view:hover {
  background: var(--uni-primary);
  color: white;
}

.btn-action-edit {
  background: transparent;
  border: 1px solid var(--uni-accent);
  color: #9a7b22;
}

.btn-action-edit:hover {
  background: var(--uni-accent);
  color: white;
}

.btn-action-share {
  background: transparent;
  border: 1px solid #7fb3d8;
  color: #2980b9;
}

.btn-action-share:hover {
  background: #2980b9;
  color: white;
}

.btn-action-delete {
  background: transparent;
  border: 1px solid #e8a9a3;
  color: var(--uni-danger);
}

.btn-action-delete:hover {
  background: var(--uni-danger);
  color: white;
}

/* === Floating Action Button === */
.fab-create {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--uni-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  box-shadow: 0 6px 20px rgba(26, 82, 118, 0.4);
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
  box-shadow: 0 8px 28px rgba(26, 82, 118, 0.5);
  color: white;
}

/* === Share Modal === */
.share-modal-content {
  border: none;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.share-toggle-input {
  width: 3em;
  height: 1.5em;
  cursor: pointer;
}

.share-link-input {
  border-radius: 10px 0 0 10px;
  background: #f8fafc;
  font-size: 0.9rem;
}

[dir="rtl"] .share-link-input {
  border-radius: 0 10px 10px 0;
}

.qr-code-wrapper {
  display: inline-block;
  padding: 12px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e4e8f0;
}

.qr-code-wrapper img {
  max-width: 180px;
}

.view-count-badge {
  display: inline-flex;
  align-items: center;
  background: var(--uni-primary-light);
  color: var(--uni-primary);
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.88rem;
  font-weight: 500;
}

/* === Responsive === */
@media (max-width: 768px) {
  .welcome-section {
    padding: 18px 20px;
  }

  .welcome-title {
    font-size: 1.2rem;
  }

  .stat-card {
    padding: 10px 4px;
  }

  .stat-number {
    font-size: 1.3rem;
  }

  .stat-label {
    font-size: 0.6rem;
  }

  .filter-pills {
    overflow-x: auto;
    flex-wrap: nowrap !important;
    padding-bottom: 4px;
  }

  .cv-card-thumbnail {
    height: 120px;
  }

  .quick-action-card {
    padding: 14px 16px;
  }

  .qa-icon {
    width: 38px;
    height: 38px;
    font-size: 1rem;
  }
}
</style>
