<template>
  <div>
    <h3 class="mb-4">{{ t('admin.activityLogs') }}</h3>

    <div class="row g-2 mb-3">
      <div class="col-md-4">
        <input type="text" class="form-control" :placeholder="t('app.search')" v-model="search" @input="debouncedFetchLogs" />
      </div>
      <div class="col-md-3">
        <select class="form-select" v-model="actionFilter" @change="fetchLogs">
          <option value="">All Actions</option>
          <option value="login">Login</option>
          <option value="register">Register</option>
          <option value="create_cv">Create CV</option>
          <option value="edit_cv">Edit CV</option>
          <option value="delete_cv">Delete CV</option>
          <option value="view_cv">View CV</option>
          <option value="download_cv">Download CV</option>
        </select>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="logs.length === 0" class="text-center py-5">
      <i class="fas fa-history fa-4x text-muted mb-3 d-block"></i>
      <h5 class="text-muted">No activity logs found</h5>
      <p class="text-muted">There are no logs matching your current filters.</p>
    </div>

    <div v-else class="table-responsive">
      <table class="table table-sm table-hover">
        <thead>
          <tr>
            <th>#</th>
            <th>{{ t('auth.email') }}</th>
            <th>Action</th>
            <th>Details</th>
            <th>IP</th>
            <th>Date</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id">
            <td>{{ log.id }}</td>
            <td>{{ log.user?.email }}</td>
            <td><span class="badge bg-info">{{ log.action }}</span></td>
            <td><small>{{ log.details }}</small></td>
            <td><small>{{ log.ip }}</small></td>
            <td><small>{{ new Date(log.created_at).toLocaleString() }}</small></td>
          </tr>
        </tbody>
      </table>
    </div>

    <nav v-if="totalPages > 1">
      <ul class="pagination pagination-sm">
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

onMounted(fetchLogs)

onUnmounted(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
})
</script>
