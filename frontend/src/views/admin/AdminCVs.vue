<template>
  <div class="admin-cvs">
    <!-- Page Header -->
    <div class="page-header">
      <h2 class="page-title">{{ t('admin.manageCVs') }}</h2>
      <p class="page-subtitle">Review, approve, and manage submitted CVs</p>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-bar-inner">
        <div class="search-input-wrapper">
          <i class="fas fa-search search-icon"></i>
          <input
            type="text"
            class="form-control search-input"
            :placeholder="t('app.search')"
            v-model="search"
            @input="debouncedFetchCVs"
          />
        </div>
        <div class="filter-select-wrapper">
          <select class="form-select filter-select" v-model="statusFilter" @change="fetchCVs">
            <option value="">{{ t('admin.byStatus') }}</option>
            <option value="draft">{{ t('cv.draft') }}</option>
            <option value="pending">{{ t('cv.pending') }}</option>
            <option value="approved">{{ t('cv.approved') }}</option>
            <option value="rejected">{{ t('cv.rejected') }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div class="data-card">
      <div class="table-responsive">
        <table class="table admin-table">
          <thead>
            <tr>
              <th>#</th>
              <th>{{ t('cv.title') }}</th>
              <th>{{ t('auth.email') }}</th>
              <th>{{ t('cv.language') }}</th>
              <th>{{ t('cv.template') }}</th>
              <th>{{ t('cv.status') }}</th>
              <th>{{ t('cv.createdAt') }}</th>
              <th class="text-center">{{ t('app.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="cv in cvs" :key="cv.id">
              <td class="text-muted fw-medium">{{ cv.id }}</td>
              <td>
                <span class="fw-semibold text-dark">{{ cv.title }}</span>
                <span v-if="cv.is_guest" class="badge badge-guest ms-2">
                  <i class="fas fa-user-secret me-1"></i>{{ t('cv.guestMode') }}
                </span>
              </td>
              <td>
                <template v-if="cv.is_guest">
                  <span class="fw-medium" style="color: #1a5276;">{{ cv.guest_name }}</span>
                  <br v-if="cv.guest_email" /><small class="text-muted">{{ cv.guest_email }}</small>
                  <br /><small class="text-muted"><i class="fas fa-globe me-1"></i>{{ cv.guest_ip }}</small>
                </template>
                <template v-else>
                  <span class="text-dark">{{ cv.user?.email }}</span>
                </template>
              </td>
              <td>
                <span class="badge badge-lang">{{ cv.language }}</span>
              </td>
              <td class="text-muted">{{ cv.template }}</td>
              <td>
                <span :class="'status-pill status-pill--' + cv.status">
                  <span class="status-dot"></span>
                  {{ cv.status }}
                </span>
              </td>
              <td class="text-muted">{{ new Date(cv.created_at).toLocaleDateString() }}</td>
              <td>
                <div class="action-buttons">
                  <router-link
                    :to="cv.is_guest ? `/shared/${cv.share_token}` : `/cv/${cv.id}`"
                    class="action-btn action-btn--view"
                    :title="t('app.view')"
                  >
                    <i class="fas fa-eye"></i>
                  </router-link>
                  <button
                    v-if="cv.status !== 'approved'"
                    @click="approve(cv.id)"
                    class="action-btn action-btn--approve"
                    :title="t('admin.approve') || 'Approve'"
                  >
                    <i class="fas fa-check"></i>
                  </button>
                  <button
                    @click="openReject(cv.id)"
                    class="action-btn action-btn--reject"
                    :title="t('admin.reject')"
                  >
                    <i class="fas fa-times"></i>
                  </button>
                  <button
                    @click="openRevision(cv.id)"
                    class="action-btn action-btn--revision"
                    :title="t('admin.requestRevision')"
                  >
                    <i class="fas fa-edit"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Reject Modal -->
    <div class="modal fade" :class="{ show: showRejectModal }" :style="{ display: showRejectModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--danger">
            <div class="modal-header-icon">
              <i class="fas fa-times-circle"></i>
            </div>
            <h5 class="modal-title">{{ t('admin.reject') }}</h5>
            <button class="btn-close btn-close-white" @click="showRejectModal = false"></button>
          </div>
          <div class="modal-body p-4">
            <label class="form-label fw-semibold mb-2">{{ t('admin.rejectReason') }}</label>
            <textarea
              class="form-control admin-textarea"
              rows="4"
              :placeholder="t('admin.rejectReason')"
              v-model="rejectReason"
            ></textarea>
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="showRejectModal = false">{{ t('app.cancel') || 'Cancel' }}</button>
            <button class="btn admin-btn admin-btn--danger" @click="reject">
              <i class="fas fa-times me-1"></i>{{ t('admin.reject') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showRejectModal" class="modal-backdrop fade show"></div>

    <!-- Revision Modal -->
    <div class="modal fade" :class="{ show: showRevisionModal }" :style="{ display: showRevisionModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--warning">
            <div class="modal-header-icon">
              <i class="fas fa-edit"></i>
            </div>
            <h5 class="modal-title">{{ t('admin.requestRevision') }}</h5>
            <button class="btn-close btn-close-white" @click="showRevisionModal = false"></button>
          </div>
          <div class="modal-body p-4">
            <label class="form-label fw-semibold mb-2">{{ t('admin.revisionNote') }}</label>
            <textarea
              class="form-control admin-textarea"
              rows="4"
              :placeholder="t('admin.revisionNote')"
              v-model="revisionNote"
            ></textarea>
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="showRevisionModal = false">{{ t('app.cancel') || 'Cancel' }}</button>
            <button class="btn admin-btn admin-btn--warning" @click="requestRevision">
              <i class="fas fa-paper-plane me-1"></i>{{ t('admin.requestRevision') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showRevisionModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'
import { useToast } from '../../composables/useToast'

const { t } = useI18n()
const toast = useToast()
const cvs = ref<any[]>([])
const search = ref('')
const statusFilter = ref('')
const showRejectModal = ref(false)
const showRevisionModal = ref(false)
const selectedCVId = ref(0)
const rejectReason = ref('')
const revisionNote = ref('')

let debounceTimer: ReturnType<typeof setTimeout>
function debouncedFetchCVs() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchCVs, 300)
}

async function fetchCVs() {
  try {
    const res = await adminAPI.listCVs({ search: search.value, status: statusFilter.value })
    cvs.value = res.data.cvs || []
  } catch {
    toast.error(t('app.error'))
  }
}

async function approve(id: number) {
  try {
    await adminAPI.approveCV(id)
    await fetchCVs()
  } catch {
    toast.error(t('app.error'))
  }
}

function openReject(id: number) { selectedCVId.value = id; rejectReason.value = ''; showRejectModal.value = true }
function openRevision(id: number) { selectedCVId.value = id; revisionNote.value = ''; showRevisionModal.value = true }

async function reject() {
  try {
    await adminAPI.rejectCV(selectedCVId.value, rejectReason.value)
    showRejectModal.value = false
    await fetchCVs()
  } catch {
    toast.error(t('app.error'))
  }
}

async function requestRevision() {
  try {
    await adminAPI.requestRevision(selectedCVId.value, revisionNote.value)
    showRevisionModal.value = false
    await fetchCVs()
  } catch {
    toast.error(t('app.error'))
  }
}

onMounted(fetchCVs)
</script>

<style scoped>
.admin-cvs {
  max-width: 1400px;
}

/* ── Page Header ── */
.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1a5276;
  margin: 0;
}

.page-subtitle {
  color: #6c757d;
  font-size: 14px;
  margin: 4px 0 0;
}

/* ── Filter Bar ── */
.filter-bar {
  margin-bottom: 20px;
}

.filter-bar-inner {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input-wrapper {
  position: relative;
  flex: 1;
  min-width: 240px;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  inset-inline-start: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: #adb5bd;
  font-size: 14px;
  pointer-events: none;
}

.search-input {
  padding-inline-start: 40px;
  border-radius: 10px;
  border: 1px solid #dee2e6;
  height: 42px;
  font-size: 14px;
  background: #ffffff;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.search-input:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

.filter-select-wrapper {
  min-width: 180px;
}

.filter-select {
  border-radius: 10px;
  border: 1px solid #dee2e6;
  height: 42px;
  font-size: 14px;
  background: #ffffff;
}

.filter-select:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

/* ── Data Card ── */
.data-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

/* ── Table ── */
.admin-table {
  margin: 0;
  font-size: 14px;
}

.admin-table thead th {
  background: #f8f9fb;
  border-bottom: 2px solid #e9ecef;
  color: #1a5276;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 14px 16px;
  white-space: nowrap;
}

.admin-table tbody td {
  padding: 14px 16px;
  vertical-align: middle;
  border-bottom: 1px solid #f0f2f5;
}

.admin-table tbody tr:hover {
  background: #f8fafc;
}

.admin-table tbody tr:nth-child(even) {
  background: #fafbfc;
}

.admin-table tbody tr:nth-child(even):hover {
  background: #f0f4f8;
}

/* ── Badges ── */
.badge-guest {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 6px;
}

.badge-lang {
  background: #f0f2f5;
  color: #2c3e50;
  font-size: 12px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 6px;
}

/* ── Status Pills ── */
.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  text-transform: capitalize;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.status-pill--draft {
  background: #f0f2f5;
  color: #6c757d;
}
.status-pill--draft .status-dot { background: #6c757d; }

.status-pill--pending {
  background: #fff8e1;
  color: #b8860b;
}
.status-pill--pending .status-dot { background: #daa520; }

.status-pill--approved {
  background: #e8f5e9;
  color: #2e7d32;
}
.status-pill--approved .status-dot { background: #43a047; }

.status-pill--rejected {
  background: #fce4ec;
  color: #c62828;
}
.status-pill--rejected .status-dot { background: #e53935; }

/* ── Action Buttons ── */
.action-buttons {
  display: flex;
  gap: 6px;
  justify-content: center;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;
}

.action-btn--view {
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
  border-color: rgba(26, 82, 118, 0.15);
}
.action-btn--view:hover {
  background: #1a5276;
  color: #fff;
}

.action-btn--approve {
  background: rgba(46, 125, 50, 0.08);
  color: #2e7d32;
  border-color: rgba(46, 125, 50, 0.15);
}
.action-btn--approve:hover {
  background: #2e7d32;
  color: #fff;
}

.action-btn--reject {
  background: rgba(198, 40, 40, 0.08);
  color: #c62828;
  border-color: rgba(198, 40, 40, 0.15);
}
.action-btn--reject:hover {
  background: #c62828;
  color: #fff;
}

.action-btn--revision {
  background: rgba(192, 120, 43, 0.08);
  color: #c0782b;
  border-color: rgba(192, 120, 43, 0.15);
}
.action-btn--revision:hover {
  background: #c0782b;
  color: #fff;
}

/* ── Modal ── */
.admin-modal {
  border: none;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.admin-modal-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px 24px;
  border: none;
  color: #ffffff;
}

.admin-modal-header--danger {
  background: linear-gradient(135deg, #c62828, #e53935);
}

.admin-modal-header--warning {
  background: linear-gradient(135deg, #c0782b, #e09138);
}

.modal-header-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}

.admin-textarea {
  border-radius: 10px;
  border: 1px solid #dee2e6;
  resize: vertical;
  font-size: 14px;
}

.admin-textarea:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

.admin-modal-footer {
  border-top: 1px solid #f0f0f0;
  padding: 16px 24px;
  gap: 8px;
}

.admin-btn {
  border-radius: 8px;
  font-weight: 600;
  padding: 8px 20px;
  font-size: 14px;
  border: none;
  color: #ffffff;
}

.admin-btn--danger { background: #c62828; }
.admin-btn--danger:hover { background: #b71c1c; color: #ffffff; }

.admin-btn--warning { background: #c0782b; }
.admin-btn--warning:hover { background: #a8692a; color: #ffffff; }
</style>
