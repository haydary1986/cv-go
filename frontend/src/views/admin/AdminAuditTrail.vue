<template>
  <div>
    <h3 class="mb-4">{{ t('admin.auditTrail') || 'Audit Trail' }}</h3>

    <!-- Filters -->
    <div class="row g-2 mb-3">
      <div class="col-md-3">
        <input type="number" class="form-control" :placeholder="'CV ID'" v-model="cvIdFilter" @input="debouncedFetch" />
      </div>
      <div class="col-md-3">
        <select class="form-select" v-model="actionFilter" @change="fetchEntries">
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

    <!-- Timeline -->
    <div class="list-group">
      <div v-for="entry in entries" :key="entry.id" class="list-group-item">
        <div class="d-flex justify-content-between align-items-start">
          <div>
            <span class="badge me-2" :class="actionBadge(entry.action)">{{ entry.action }}</span>
            <strong>{{ entry.user?.full_name_en || entry.user?.email || 'System' }}</strong>
            <span class="text-muted ms-2">CV #{{ entry.cv_id }}</span>
          </div>
          <small class="text-muted">{{ formatDate(entry.created_at) }}</small>
        </div>
        <div v-if="entry.details" class="mt-1 text-muted small">{{ entry.details }}</div>
        <div class="mt-1 text-muted small">
          <i class="fas fa-globe me-1"></i>{{ entry.ip }}
        </div>
      </div>
    </div>

    <div v-if="entries.length === 0" class="text-center text-muted py-4">
      <i class="fas fa-history fa-3x mb-3"></i>
      <p>{{ t('app.noData') || 'No audit entries found' }}</p>
    </div>

    <!-- Pagination -->
    <nav v-if="totalPages > 1" class="mt-3">
      <ul class="pagination justify-content-center">
        <li class="page-item" :class="{ disabled: page <= 1 }">
          <a class="page-link" href="#" @click.prevent="page--; fetchEntries()">«</a>
        </li>
        <li class="page-item active"><span class="page-link">{{ page }} / {{ totalPages }}</span></li>
        <li class="page-item" :class="{ disabled: page >= totalPages }">
          <a class="page-link" href="#" @click.prevent="page++; fetchEntries()">»</a>
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

function actionBadge(action: string) {
  const map: Record<string, string> = {
    create: 'bg-success', update: 'bg-primary', approve: 'bg-success',
    reject: 'bg-danger', revision: 'bg-warning text-dark', delete: 'bg-danger'
  }
  return map[action] || 'bg-secondary'
}

function formatDate(d: string) {
  return new Date(d).toLocaleString()
}

onMounted(fetchEntries)
</script>
