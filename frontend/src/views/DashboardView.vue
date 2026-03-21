<template>
  <div class="dashboard-page">
    <div class="container py-4">
      <!-- Welcome Header -->
      <div class="welcome-section mb-4">
        <div class="row align-items-center">
          <div class="col-lg-7">
            <div class="d-flex align-items-center mb-2">
              <div class="welcome-avatar me-3">
                <i class="fas fa-user-circle"></i>
              </div>
              <div>
                <h2 class="welcome-title mb-1">
                  {{ t('dashboard.welcome') }}<span v-if="userName" class="text-primary">, {{ userName }}</span>
                </h2>
                <p class="welcome-subtitle text-muted mb-0">{{ t('dashboard.subtitle') }}</p>
              </div>
            </div>
          </div>
          <div class="col-lg-5">
            <div class="row g-2 mt-2 mt-lg-0">
              <div class="col-4">
                <div class="stat-card stat-total">
                  <div class="stat-number">{{ totalCount }}</div>
                  <div class="stat-label">{{ t('dashboard.totalCVs') }}</div>
                </div>
              </div>
              <div class="col-4">
                <div class="stat-card stat-approved">
                  <div class="stat-number">{{ approvedCount }}</div>
                  <div class="stat-label">{{ t('cv.approved') }}</div>
                </div>
              </div>
              <div class="col-4">
                <div class="stat-card stat-pending">
                  <div class="stat-number">{{ pendingCount }}</div>
                  <div class="stat-label">{{ t('cv.pending') }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Toolbar -->
      <div class="toolbar-section mb-4">
        <div class="row align-items-center g-3">
          <div class="col-md-4">
            <div class="search-box">
              <i class="fas fa-search search-icon"></i>
              <input
                type="text"
                class="form-control form-control-lg search-input"
                :placeholder="t('dashboard.searchPlaceholder')"
                v-model="searchQuery"
              />
              <button v-if="searchQuery" class="search-clear" @click="searchQuery = ''">
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
          <div class="col-md-5">
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
          </div>
          <div class="col-md-3 text-md-end">
            <div class="d-flex gap-2 justify-content-md-end">
              <button @click="exportJSON" class="btn btn-outline-secondary">
                <i class="fas fa-download me-1"></i><span class="d-none d-sm-inline">{{ t('admin.exportJSON') }}</span>
              </button>
              <router-link to="/cv/create" class="btn btn-primary btn-create">
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
            <i class="fas fa-file-alt"></i>
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
              <span class="status-badge" :class="`status-${cv.status}`">
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

            <div class="card-footer bg-transparent border-top pt-0 pb-3">
              <div class="cv-card-actions d-flex gap-2">
                <router-link :to="`/cv/${cv.id}`" class="btn btn-sm btn-outline-primary flex-fill" @click.stop>
                  <i class="fas fa-eye me-1"></i>{{ t('app.view') }}
                </router-link>
                <router-link :to="`/cv/${cv.id}/edit`" class="btn btn-sm btn-outline-secondary flex-fill" @click.stop>
                  <i class="fas fa-edit me-1"></i>{{ t('app.edit') }}
                </router-link>
                <button @click.stop="handleShare(cv)" class="btn btn-sm btn-outline-info" :title="t('app.share')">
                  <i class="fas fa-share-alt"></i>
                </button>
                <button @click.stop="handleDelete(cv.id)" class="btn btn-sm btn-outline-danger" :title="t('app.delete')">
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
              <div class="form-check form-switch d-inline-flex align-items-center gap-2">
                <input
                  class="form-check-input share-toggle-input"
                  type="checkbox"
                  :checked="selectedCV.is_shared"
                  @change="toggleShare"
                  id="shareToggle"
                />
                <label class="form-check-label fw-medium" for="shareToggle">
                  {{ selectedCV.is_shared ? t('cv.disableShare') : t('cv.enableShare') }}
                </label>
              </div>
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

              <div class="text-center mb-3" v-if="selectedCV.qr_code_data">
                <label class="form-label text-muted small text-uppercase fw-bold d-block">{{ t('cv.qrCode') }}</label>
                <div class="qr-code-wrapper">
                  <img :src="selectedCV.qr_code_data" alt="QR Code" class="img-fluid" />
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
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useCVStore, type CV } from '../stores/cv'
import { useAuthStore } from '../stores/auth'
import { cvAPI } from '../services/api'

const { t } = useI18n()
const cvStore = useCVStore()
const authStore = useAuthStore()

const shareModal = ref<HTMLElement>()
const selectedCV = ref<CV | null>(null)
const searchQuery = ref('')
const activeFilter = ref('all')
const copied = ref(false)
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
  if (!selectedCV.value) return
  const updated = await cvStore.toggleShare(selectedCV.value.id)
  selectedCV.value = updated
  cvStore.fetchCVs()
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
/* === Welcome Section === */
.welcome-section {
  background: linear-gradient(135deg, #f8f9ff 0%, #eef1ff 100%);
  border-radius: 16px;
  padding: 24px 28px;
  border: 1px solid #e4e8f8;
}

.welcome-avatar {
  font-size: 2.8rem;
  color: #6366f1;
  line-height: 1;
}

.welcome-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.welcome-subtitle {
  font-size: 0.95rem;
}

/* Stats Cards */
.stat-card {
  text-align: center;
  background: white;
  border-radius: 12px;
  padding: 14px 8px;
  border: 1px solid #e8ecf3;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.stat-number {
  font-size: 1.6rem;
  font-weight: 800;
  line-height: 1.2;
}

.stat-label {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-total .stat-number { color: #6366f1; }
.stat-approved .stat-number { color: #22c55e; }
.stat-pending .stat-number { color: #f59e0b; }

/* === Toolbar === */
.toolbar-section {
  background: white;
  border-radius: 12px;
  padding: 16px 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.06);
  border: 1px solid #e8ecf3;
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
  border-radius: 10px;
  border: 1px solid #e2e8f0;
  font-size: 0.95rem;
  background: #f8fafc;
  transition: border-color 0.2s, box-shadow 0.2s, background 0.2s;
}

[dir="rtl"] .search-input {
  padding-left: 12px;
  padding-right: 40px;
}

.search-input:focus {
  background: white;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
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
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.filter-pill:hover {
  border-color: #6366f1;
  color: #6366f1;
  background: #f0f0ff;
}

.filter-pill.active {
  background: #6366f1;
  color: white;
  border-color: #6366f1;
}

.filter-count {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 10px;
  padding: 1px 7px;
  font-size: 0.75rem;
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
  background: #6366f1;
  border-color: #6366f1;
}

.btn-create:hover {
  background: #4f46e5;
  border-color: #4f46e5;
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
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.5s infinite;
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
.skeleton-btn {
  height: 32px;
  border-radius: 6px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
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
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #eef1ff 0%, #e0e4ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  font-size: 3rem;
  color: #6366f1;
  animation: empty-float 3s ease-in-out infinite;
}

@keyframes empty-float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}

.empty-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1e293b;
}

.empty-description {
  max-width: 400px;
  font-size: 1rem;
  line-height: 1.6;
}

.btn-create-empty {
  background: #6366f1;
  border-color: #6366f1;
  border-radius: 12px;
  padding: 12px 32px;
  font-weight: 600;
  font-size: 1.05rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-create-empty:hover {
  background: #4f46e5;
  border-color: #4f46e5;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.3);
}

.feature-item {
  text-align: center;
  padding: 20px;
}

.feature-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #f0f0ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 12px;
  font-size: 1.2rem;
  color: #6366f1;
}

/* === CV Cards === */
.cv-card {
  border: none;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06), 0 1px 2px rgba(0, 0, 0, 0.04);
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  cursor: pointer;
}

.cv-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.1), 0 4px 10px rgba(0, 0, 0, 0.06);
}

/* Thumbnail Backgrounds by Template */
.cv-card-thumbnail {
  height: 140px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

/* Status Badge */
.status-badge {
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

[dir="rtl"] .status-badge {
  right: auto;
  left: 10px;
}

.status-draft {
  background: rgba(100, 116, 139, 0.9);
  color: white;
}

.status-pending {
  background: rgba(245, 158, 11, 0.95);
  color: #1e293b;
}

.status-approved {
  background: rgba(34, 197, 94, 0.9);
  color: white;
}

.status-rejected {
  background: rgba(239, 68, 68, 0.9);
  color: white;
}

/* Card Body */
.cv-card-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: #1e293b;
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
  color: #0ea5e9;
  background: #f0f9ff;
  padding: 3px 10px;
  border-radius: 8px;
}

/* Card Actions */
.cv-card-actions .btn {
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 500;
  transition: all 0.2s;
}

.cv-card-actions .btn:hover {
  transform: translateY(-1px);
}

/* === Floating Action Button === */
.fab-create {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #6366f1;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.3rem;
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
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
  box-shadow: 0 8px 28px rgba(99, 102, 241, 0.5);
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
  border: 2px solid #e8ecf3;
}

.qr-code-wrapper img {
  max-width: 180px;
}

.view-count-badge {
  display: inline-flex;
  align-items: center;
  background: #f0f9ff;
  color: #0284c7;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.88rem;
  font-weight: 500;
}

/* === Responsive === */
@media (max-width: 768px) {
  .welcome-section {
    padding: 16px 18px;
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
