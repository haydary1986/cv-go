<template>
  <div class="container py-4">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h2>{{ t('app.myCVs') }}</h2>
      <div class="d-flex gap-2">
        <button @click="exportJSON" class="btn btn-outline-secondary btn-sm">
          <i class="fas fa-download me-1"></i>{{ t('admin.exportJSON') }}
        </button>
        <router-link to="/cv/create" class="btn btn-primary">
          <i class="fas fa-plus me-1"></i>{{ t('app.createCV') }}
        </router-link>
      </div>
    </div>

    <div v-if="cvStore.loading" class="text-center py-5">
      <div class="spinner-border text-primary"></div>
      <p class="mt-2">{{ t('app.loading') }}</p>
    </div>

    <div v-else-if="cvStore.cvs.length === 0" class="text-center py-5">
      <i class="fas fa-file-alt fa-4x text-muted mb-3"></i>
      <h4 class="text-muted">{{ t('app.noData') }}</h4>
      <router-link to="/cv/create" class="btn btn-primary mt-3">
        {{ t('app.createCV') }}
      </router-link>
    </div>

    <div v-else class="row g-4">
      <div class="col-md-6 col-lg-4" v-for="cv in cvStore.cvs" :key="cv.id">
        <div class="card h-100 shadow-sm">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-start mb-2">
              <h5 class="card-title mb-0">{{ cv.title || t('cv.title') }}</h5>
              <span :class="statusBadge(cv.status)" class="badge">
                {{ t(`cv.${cv.status}`) }}
              </span>
            </div>
            <p class="text-muted small mb-1">
              <i class="fas fa-language me-1"></i>{{ cv.language === 'ar' ? t('app.arabic') : t('app.english') }}
            </p>
            <p class="text-muted small mb-1">
              <i class="fas fa-palette me-1"></i>{{ t(`templates.${cv.template}`) }}
            </p>
            <p class="text-muted small mb-2">
              <i class="fas fa-calendar me-1"></i>{{ formatDate(cv.created_at) }}
            </p>
            <div v-if="cv.is_shared" class="small text-info mb-2">
              <i class="fas fa-eye me-1"></i>{{ cv.view_count }} {{ t('cv.viewCount') }}
            </div>
          </div>
          <div class="card-footer bg-white border-top-0 d-flex gap-2">
            <router-link :to="`/cv/${cv.id}`" class="btn btn-sm btn-outline-primary flex-fill">
              <i class="fas fa-eye me-1"></i>{{ t('app.view') }}
            </router-link>
            <router-link :to="`/cv/${cv.id}/edit`" class="btn btn-sm btn-outline-secondary flex-fill">
              <i class="fas fa-edit me-1"></i>{{ t('app.edit') }}
            </router-link>
            <button @click="handleShare(cv)" class="btn btn-sm btn-outline-info">
              <i class="fas fa-share-alt"></i>
            </button>
            <button @click="handleDelete(cv.id)" class="btn btn-sm btn-outline-danger">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Share Modal -->
    <div class="modal fade" ref="shareModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ t('cv.shareCV') }}</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal"></button>
          </div>
          <div class="modal-body" v-if="selectedCV">
            <div class="mb-3">
              <label class="form-label">{{ t('cv.shareLink') }}</label>
              <div class="input-group">
                <input type="text" class="form-control" :value="shareUrl" readonly />
                <button @click="copyShareLink" class="btn btn-outline-primary">
                  <i class="fas fa-copy"></i>
                </button>
              </div>
            </div>
            <div class="text-center mb-3" v-if="selectedCV.qr_code_data">
              <img :src="selectedCV.qr_code_data" alt="QR Code" class="img-fluid" style="max-width: 200px" />
            </div>
            <div class="d-flex justify-content-between align-items-center">
              <span>{{ t('cv.viewCount') }}: {{ selectedCV.view_count }}</span>
              <button @click="toggleShare" class="btn btn-sm" :class="selectedCV.is_shared ? 'btn-danger' : 'btn-success'">
                {{ selectedCV.is_shared ? t('cv.disableShare') : t('cv.enableShare') }}
              </button>
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
import { cvAPI } from '../services/api'

const { t } = useI18n()
const cvStore = useCVStore()

const shareModal = ref<HTMLElement>()
const selectedCV = ref<CV | null>(null)
let bsModal: any = null

const shareUrl = computed(() => {
  if (!selectedCV.value) return ''
  return `${window.location.origin}/shared/${selectedCV.value.share_token}`
})

onMounted(() => {
  cvStore.fetchCVs()
})

function statusBadge(status: string) {
  const map: Record<string, string> = {
    draft: 'bg-secondary',
    pending: 'bg-warning text-dark',
    approved: 'bg-success',
    rejected: 'bg-danger',
  }
  return map[status] || 'bg-secondary'
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}

async function handleShare(cv: CV) {
  selectedCV.value = cv
  if (shareModal.value) {
    const { Modal } = await import('bootstrap')
    bsModal = new Modal(shareModal.value)
    bsModal.show()
  }
}

function copyShareLink() {
  navigator.clipboard.writeText(shareUrl.value)
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
