<template>
  <div class="admin-logs">
    <!-- Page Header -->
    <div class="page-header">
      <h2 class="page-title">{{ t('admin.activityLogs') }}</h2>
      <p class="page-subtitle">Track user activity and system events</p>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-bar-inner">
        <div class="search-input-wrapper">
          <i class="fas fa-search search-icon"></i>
          <input type="text" class="form-control search-input" :placeholder="t('app.search')" v-model="search" @input="debouncedFetchLogs" />
        </div>
        <div class="filter-select-wrapper">
          <select class="form-select filter-select" v-model="actionFilter" @change="fetchLogs">
            <option value="">{{ t('admin.allActions') }}</option>
            <option value="login">{{ t('auth.loginBtn') }}</option>
            <option value="register">{{ t('auth.registerBtn') }}</option>
            <option value="create_cv">{{ t('app.createCV') }}</option>
            <option value="edit_cv">{{ t('app.edit') }}</option>
            <option value="delete_cv">{{ t('app.delete') }}</option>
            <option value="view_cv">{{ t('app.view') }}</option>
            <option value="download_cv">{{ t('app.download') }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="logs.length === 0" class="empty-state">
      <div class="empty-state-icon">
        <i class="fas fa-history"></i>
      </div>
      <h5 class="empty-state-title">{{ t('admin.noLogsFound') }}</h5>
      <p class="empty-state-text">{{ t('admin.noLogsFoundDesc') }}</p>
    </div>

    <!-- Timeline -->
    <div v-else class="timeline-card">
      <div class="timeline">
        <div v-for="log in logs" :key="log.id" class="timeline-item">
          <div class="timeline-marker" :class="'timeline-marker--' + getActionType(log.action)">
            <i :class="getActionIcon(log.action)"></i>
          </div>
          <div class="timeline-content">
            <div class="timeline-header">
              <div class="timeline-user">
                <span class="timeline-user-email">{{ log.user?.email || 'System' }}</span>
                <span class="timeline-action-badge" :class="'timeline-action-badge--' + getActionType(log.action)">
                  {{ log.action }}
                </span>
              </div>
              <span class="timeline-time">
                <i class="fas fa-clock me-1"></i>{{ new Date(log.created_at).toLocaleString() }}
              </span>
            </div>
            <div v-if="log.details" class="timeline-details">{{ log.details }}</div>
            <div class="timeline-meta">
              <span class="timeline-ip"><i class="fas fa-globe me-1"></i>{{ log.ip }}</span>
              <span class="timeline-id">#{{ log.id }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <nav v-if="totalPages > 1" class="mt-4">
      <ul class="pagination admin-pagination justify-content-center">
        <li class="page-item" v-for="p in totalPages" :key="p" :class="{ active: p === page }">
          <a class="page-link" href="#" @click.prevent="page = p; fetchLogs()">{{ p }}</a>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'

const { t } = useI18n()
const logs = ref<any[]>([])
const search = ref('')
const actionFilter = ref('')
const page = ref(1)
const totalPages = ref(1)

let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

function debouncedFetchLogs() {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    page.value = 1
    fetchLogs()
  }, 300)
}

async function fetchLogs() {
  const res = await adminAPI.getActivityLogs({ search: search.value, action: actionFilter.value, page: page.value })
  logs.value = res.data.logs || []
  const total = res.data.total || 0
  totalPages.value = Math.ceil(total / 50) || 1
}

function getActionType(action: string): string {
  if (action.includes('login') || action.includes('register')) return 'auth'
  if (action.includes('create')) return 'create'
  if (action.includes('edit') || action.includes('update')) return 'edit'
  if (action.includes('delete')) return 'delete'
  if (action.includes('view') || action.includes('download')) return 'view'
  return 'default'
}

function getActionIcon(action: string): string {
  const map: Record<string, string> = {
    login: 'fas fa-sign-in-alt',
    register: 'fas fa-user-plus',
    create_cv: 'fas fa-plus',
    edit_cv: 'fas fa-edit',
    delete_cv: 'fas fa-trash',
    view_cv: 'fas fa-eye',
    download_cv: 'fas fa-download',
  }
  return map[action] || 'fas fa-circle'
}

onMounted(fetchLogs)

onUnmounted(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
})
</script>

<style scoped>
.admin-logs {
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

/* ── Timeline ── */
.timeline {
  position: relative;
}

.timeline::before {
  content: '';
  position: absolute;
  inset-inline-start: 19px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #e9ecef;
}

.timeline-item {
  display: flex;
  gap: 16px;
  padding: 16px 0;
  position: relative;
}

.timeline-item + .timeline-item {
  border-top: none;
}

.timeline-marker {
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

.timeline-marker--auth { background: rgba(26, 82, 118, 0.1); color: #1a5276; }
.timeline-marker--create { background: rgba(46, 125, 50, 0.1); color: #2e7d32; }
.timeline-marker--edit { background: rgba(192, 152, 43, 0.1); color: #c0982b; }
.timeline-marker--delete { background: rgba(198, 40, 40, 0.1); color: #c62828; }
.timeline-marker--view { background: rgba(102, 16, 242, 0.1); color: #6610f2; }
.timeline-marker--default { background: rgba(108, 117, 125, 0.1); color: #6c757d; }

.timeline-content {
  flex: 1;
  min-width: 0;
}

.timeline-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 4px;
}

.timeline-user {
  display: flex;
  align-items: center;
  gap: 8px;
}

.timeline-user-email {
  font-weight: 600;
  color: #2c3e50;
  font-size: 14px;
}

.timeline-action-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.timeline-action-badge--auth { background: rgba(26, 82, 118, 0.08); color: #1a5276; }
.timeline-action-badge--create { background: rgba(46, 125, 50, 0.08); color: #2e7d32; }
.timeline-action-badge--edit { background: rgba(192, 152, 43, 0.08); color: #c0982b; }
.timeline-action-badge--delete { background: rgba(198, 40, 40, 0.08); color: #c62828; }
.timeline-action-badge--view { background: rgba(102, 16, 242, 0.08); color: #6610f2; }
.timeline-action-badge--default { background: rgba(108, 117, 125, 0.08); color: #6c757d; }

.timeline-time {
  font-size: 12px;
  color: #adb5bd;
}

.timeline-details {
  font-size: 13px;
  color: #6c757d;
  margin: 4px 0;
  padding: 6px 10px;
  background: #f8f9fb;
  border-radius: 6px;
}

.timeline-meta {
  display: flex;
  gap: 16px;
  margin-top: 6px;
}

.timeline-ip,
.timeline-id {
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
}

.admin-pagination .page-item.active .page-link {
  background: #1a5276;
  border-color: #1a5276;
  color: #ffffff;
}
</style>
