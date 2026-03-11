<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h3>{{ t('admin.manageUsers') }}</h3>
      <div class="d-flex gap-2">
        <button @click="exportCSV" class="btn btn-outline-success btn-sm"><i class="fas fa-download me-1"></i>{{ t('admin.exportCSV') }}</button>
        <label class="btn btn-outline-primary btn-sm mb-0">
          <i class="fas fa-upload me-1"></i>{{ t('admin.importCSV') }}
          <input type="file" accept=".csv" class="d-none" @change="importCSV" />
        </label>
      </div>
    </div>

    <div class="row g-2 mb-3">
      <div class="col-md-4">
        <input type="text" class="form-control" :placeholder="t('app.search')" v-model="search" @input="debouncedFetchUsers" />
      </div>
      <div class="col-md-3">
        <select class="form-select" v-model="roleFilter" @change="fetchUsers">
          <option value="">{{ t('app.filter') }}</option>
          <option value="student">Student</option>
          <option value="teacher">Teacher</option>
          <option value="admin">Admin</option>
        </select>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="users.length === 0" class="text-center py-5">
      <i class="fas fa-users fa-4x text-muted mb-3 d-block"></i>
      <h5 class="text-muted">No users found</h5>
      <p class="text-muted">Try adjusting your search or filter criteria.</p>
    </div>

    <div v-else class="table-responsive">
      <table class="table table-hover">
        <thead>
          <tr>
            <th>#</th>
            <th>{{ t('auth.fullNameAr') }}</th>
            <th>{{ t('auth.email') }}</th>
            <th>Role</th>
            <th>{{ t('auth.faculty') }}</th>
            <th>{{ t('ai.credits') }}</th>
            <th>{{ t('app.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>{{ u.id }}</td>
            <td>{{ u.full_name_ar || u.full_name_en }}</td>
            <td>{{ u.email }}</td>
            <td><span class="badge" :class="u.role === 'admin' ? 'bg-danger' : u.role === 'teacher' ? 'bg-info' : 'bg-secondary'">{{ u.role }}</span></td>
            <td>{{ u.faculty?.name_en || '-' }}</td>
            <td>{{ u.ai_credits }}</td>
            <td>
              <button @click="openCreditsModal(u)" class="btn btn-sm btn-outline-primary">
                <i class="fas fa-coins"></i> {{ t('admin.updateCredits') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Credits Modal -->
    <div class="modal fade" :class="{ show: showCreditsModal }" :style="{ display: showCreditsModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-sm">
        <div class="modal-content">
          <div class="modal-header">
            <h5>{{ t('admin.updateCredits') }}</h5>
            <button class="btn-close" @click="showCreditsModal = false"></button>
          </div>
          <div class="modal-body">
            <input type="number" class="form-control" v-model.number="newCredits" min="0" />
          </div>
          <div class="modal-footer">
            <button class="btn btn-primary" @click="updateCredits">{{ t('app.save') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showCreditsModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'

const { t } = useI18n()
const users = ref<any[]>([])
const search = ref('')
const roleFilter = ref('')
const showCreditsModal = ref(false)
const selectedUserId = ref(0)
const newCredits = ref(0)

let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

function debouncedFetchUsers() {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    fetchUsers()
  }, 300)
}

async function fetchUsers() {
  const res = await adminAPI.listUsers({ search: search.value, role: roleFilter.value })
  users.value = res.data.users || []
}

function openCreditsModal(u: any) {
  selectedUserId.value = u.id
  newCredits.value = u.ai_credits
  showCreditsModal.value = true
}

async function updateCredits() {
  await adminAPI.updateCredits(selectedUserId.value, newCredits.value)
  showCreditsModal.value = false
  fetchUsers()
}

async function exportCSV() {
  const res = await adminAPI.exportUsersCSV()
  const blob = new Blob([res.data], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url; a.download = 'users.csv'; a.click()
  URL.revokeObjectURL(url)
}

async function importCSV(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  await adminAPI.importUsersCSV(file)
  fetchUsers()
}

onMounted(fetchUsers)

onUnmounted(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
})
</script>
