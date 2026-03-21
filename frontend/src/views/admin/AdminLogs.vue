<template>
  <div>
    <h3 class="mb-4">{{ t('admin.activityLogs') }}</h3>

    <div class="row g-2 mb-3">
      <div class="col-md-4">
        <input type="text" class="form-control" :placeholder="t('app.search')" v-model="search" @input="debouncedFetchLogs" />
      </div>
      <div class="col-md-3">
        <select class="form-select" v-model="actionFilter" @change="fetchLogs">
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

    <!-- Empty State -->
    <div v-if="logs.length === 0" class="text-center py-5">
      <i class="fas fa-history fa-4x text-muted mb-3 d-block"></i>
      <h5 class="text-muted">{{ t('admin.noLogsFound') }}</h5>
      <p class="text-muted">{{ t('admin.noLogsFoundDesc') }}</p>
    </div>

    <div v-else class="table-responsive">
      <table class="table table-sm table-hover">
        <thead>
          <tr>
            <th>#</th>
            <th>{{ t('auth.email') }}</th>
            <th>{{ t('admin.action') }}</th>
            <th>{{ t('admin.details') }}</th>
            <th>{{ t('admin.ip') }}</th>
            <th>{{ t('admin.date') }}</th>
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
