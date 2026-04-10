<template>
  <div class="admin-audit">
    <!-- Page Header -->
    <div class="page-header">
      <h2 class="page-title">{{ t('admin.auditTrail') || 'Audit Trail' }}</h2>
      <p class="page-subtitle">Track CV changes, approvals, and rejections</p>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-bar-inner">
        <div class="filter-input-wrapper">
          <i class="fas fa-hashtag filter-icon"></i>
          <input type="number" class="form-control filter-input" :placeholder="'CV ID'" v-model="cvIdFilter" @input="debouncedFetch" />
        </div>
        <div class="filter-select-wrapper">
          <select class="form-select filter-select" v-model="actionFilter" @change="fetchEntries">
            <option value="">{{ t('app.all') || 'All Actions' }}</option>
            <option value="create">Create</option>
            <option value="update">Update</option>
            <option value="approve">Approve</option>
            <option value="reject">Reject</option>
            <option value="revision">Revision</option>
            <option value="delete">Delete</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Timeline -->
    <div v-if="entries.length > 0" class="timeline-card">
      <div class="audit-timeline">
        <div v-for="entry in entries" :key="entry.id" class="audit-item">
          <div class="audit-marker" :class="'audit-marker--' + entry.action">
            <i :class="getAuditIcon(entry.action)"></i>
          </div>
          <div class="audit-content">
            <div class="audit-header">
              <div class="audit-header-left">
                <span class="audit-action-badge" :class="'audit-action-badge--' + entry.action">
                  {{ entry.action }}
                </span>
                <span class="audit-actor">
                  {{ entry.user?.full_name_en || entry.user?.email || 'System' }}
                </span>
                <span class="audit-cv-ref">
                  <i class="fas fa-file-alt me-1"></i>CV #{{ entry.cv_id }}
                </span>
              </div>
              <span class="audit-time">
                <i class="fas fa-clock me-1"></i>{{ formatDate(entry.created_at) }}
              </span>
            </div>
            <div v-if="entry.details" class="audit-details">
              <div class="audit-details-content">
                <i class="fas fa-info-circle me-2"></i>{{ entry.details }}
              </div>
            </div>
            <div class="audit-meta">
              <span class="audit-ip"><i class="fas fa-globe me-1"></i>{{ entry.ip }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="entries.length === 0" class="empty-state">
      <div class="empty-state-icon">
        <i class="fas fa-shield-alt"></i>
      </div>
      <h5 class="empty-state-title">{{ t('app.noData') || 'No audit entries found' }}</h5>
      <p class="empty-state-text">Audit entries will appear here when CV actions occur</p>
    </div>

    <!-- Pagination -->
    <nav v-if="totalPages > 1" class="mt-4">
      <ul class="pagination admin-pagination justify-content-center">
        <li class="page-item" :class="{ disabled: page <= 1 }">
          <a class="page-link" href="#" @click.prevent="page--; fetchEntries()">
            <i class="fas fa-chevron-left"></i>
          </a>
        </li>
        <li class="page-item active">
          <span class="page-link">{{ page }} / {{ totalPages }}</span>
        </li>
        <li class="page-item" :class="{ disabled: page >= totalPages }">
          <a class="page-link" href="#" @click.prevent="page++; fetchEntries()">
            <i class="fas fa-chevron-right"></i>
          </a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'

const { t } = useI18n()
const entries = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const cvIdFilter = ref('')
const actionFilter = ref('')

const totalPages = computed(() => Math.ceil(total.value / 50))

let debounceTimer: ReturnType<typeof setTimeout>
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchEntries, 300)
}

async function fetchEntries() {
  try {
    const params: any = { page: page.value, limit: 50 }
    if (cvIdFilter.value) params.cv_id = cvIdFilter.value
    if (actionFilter.value) params.action = actionFilter.value
    const res = await adminAPI.getAuditTrail(params)
    entries.value = res.data.entries || []
    total.value = res.data.total || 0
  } catch {
    entries.value = []
  }
}


function formatDate(d: string) {
  return new Date(d).toLocaleString()
}

function getAuditIcon(action: string): string {
  const map: Record<string, string> = {
    create: 'fas fa-plus',
    update: 'fas fa-edit',
    approve: 'fas fa-check',
    reject: 'fas fa-times',
    revision: 'fas fa-redo',
    delete: 'fas fa-trash',
  }
  return map[action] || 'fas fa-circle'
}

onMounted(fetchEntries)
</script>

<style scoped>
.admin-audit {
  max-width: 1000px;
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

.filter-input-wrapper {
  position: relative;
  min-width: 160px;
  max-width: 200px;
}

.filter-icon {
  position: absolute;
  inset-inline-start: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: #adb5bd;
  font-size: 13px;
  pointer-events: none;
}

.filter-input {
  padding-inline-start: 36px;
  border-radius: 10px;
  border: 1px solid #dee2e6;
  height: 42px;
  font-size: 14px;
  background: #ffffff;
}

.filter-input:focus {
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
}

.filter-select:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

/* ── Empty State ── */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
}

.empty-state-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: #f0f4f8;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 32px;
  color: #1a5276;
}

.empty-state-title {
  color: #2c3e50;
  font-weight: 600;
}

.empty-state-text {
  color: #6c757d;
}

/* ── Timeline Card ── */
.timeline-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  padding: 24px;
}

/* ── Audit Timeline ── */
.audit-timeline {
  position: relative;
}

.audit-timeline::before {
  content: '';
  position: absolute;
  inset-inline-start: 19px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #e9ecef;
}

.audit-item {
  display: flex;
  gap: 16px;
  padding: 16px 0;
  position: relative;
}

.audit-marker {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
  z-index: 1;
  border: 2px solid #ffffff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.audit-marker--create { background: rgba(46, 125, 50, 0.1); color: #2e7d32; }
.audit-marker--update { background: rgba(26, 82, 118, 0.1); color: #1a5276; }
.audit-marker--approve { background: rgba(46, 125, 50, 0.1); color: #2e7d32; }
.audit-marker--reject { background: rgba(198, 40, 40, 0.1); color: #c62828; }
.audit-marker--revision { background: rgba(192, 152, 43, 0.1); color: #c0982b; }
.audit-marker--delete { background: rgba(198, 40, 40, 0.1); color: #c62828; }

.audit-content {
  flex: 1;
  min-width: 0;
}

.audit-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 6px;
}

.audit-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.audit-action-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.audit-action-badge--create { background: rgba(46, 125, 50, 0.08); color: #2e7d32; }
.audit-action-badge--update { background: rgba(26, 82, 118, 0.08); color: #1a5276; }
.audit-action-badge--approve { background: rgba(46, 125, 50, 0.08); color: #2e7d32; }
.audit-action-badge--reject { background: rgba(198, 40, 40, 0.08); color: #c62828; }
.audit-action-badge--revision { background: rgba(192, 152, 43, 0.08); color: #c0982b; }
.audit-action-badge--delete { background: rgba(198, 40, 40, 0.08); color: #c62828; }

.audit-actor {
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
}

.audit-cv-ref {
  font-size: 13px;
  color: #1a5276;
  font-weight: 500;
  padding: 2px 8px;
  background: rgba(26, 82, 118, 0.05);
  border-radius: 4px;
}

.audit-time {
  font-size: 12px;
  color: #adb5bd;
}

.audit-details {
  margin: 6px 0;
}

.audit-details-content {
  padding: 8px 12px;
  background: #f8f9fb;
  border-radius: 8px;
  border-inline-start: 3px solid #1a5276;
  font-size: 13px;
  color: #495057;
}

.audit-meta {
  display: flex;
  gap: 16px;
  margin-top: 6px;
}

.audit-ip {
  font-size: 12px;
  color: #adb5bd;
}

/* ── Pagination ── */
.admin-pagination .page-link {
  border-radius: 8px;
  margin: 0 2px;
  border: 1px solid #dee2e6;
  color: #1a5276;
  font-size: 14px;
  font-weight: 500;
  padding: 8px 14px;
}

.admin-pagination .page-item.active .page-link {
  background: #1a5276;
  border-color: #1a5276;
  color: #ffffff;
}

.admin-pagination .page-item.disabled .page-link {
  color: #adb5bd;
}
</style>
