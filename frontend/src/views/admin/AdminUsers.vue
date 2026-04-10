<template>
  <div class="admin-users">
    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h2 class="page-title">{{ t('admin.manageUsers') }}</h2>
        <p class="page-subtitle">Manage system users, roles, and permissions</p>
      </div>
      <div class="header-actions">
        <button @click="showCreateModal = true" class="btn admin-btn admin-btn--primary">
          <i class="fas fa-user-plus me-2"></i>{{ t('admin.addUser') }}
        </button>
        <button @click="exportCSV" class="btn admin-btn admin-btn--outline">
          <i class="fas fa-download me-1"></i>{{ t('admin.exportCSV') }}
        </button>
        <label class="btn admin-btn admin-btn--outline mb-0">
          <i class="fas fa-upload me-1"></i>{{ t('admin.importCSV') }}
          <input type="file" accept=".csv" class="d-none" @change="importCSV" />
        </label>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="filter-bar">
      <div class="filter-bar-inner">
        <div class="search-input-wrapper">
          <i class="fas fa-search search-icon"></i>
          <input type="text" class="form-control search-input" :placeholder="t('app.search')" v-model="search" @input="debouncedFetchUsers" />
        </div>
        <div class="filter-select-wrapper">
          <select class="form-select filter-select" v-model="roleFilter" @change="fetchUsers">
            <option value="">{{ t('app.filter') }}</option>
            <option value="student">{{ t('admin.student') }}</option>
            <option value="teacher">{{ t('admin.teacher') }}</option>
            <option value="admin">{{ t('admin.adminRole') }}</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="users.length === 0" class="empty-state">
      <div class="empty-state-icon">
        <i class="fas fa-users"></i>
      </div>
      <h5 class="empty-state-title">{{ t('admin.noUsersFound') }}</h5>
      <p class="empty-state-text">{{ t('admin.noUsersFoundDesc') }}</p>
    </div>

    <!-- Table -->
    <div v-else class="data-card">
      <div class="table-responsive">
        <table class="table admin-table">
          <thead>
            <tr>
              <th>#</th>
              <th>{{ t('admin.name') }}</th>
              <th>{{ t('auth.email') }}</th>
              <th>{{ t('admin.role') }}</th>
              <th>{{ t('auth.faculty') }}</th>
              <th>{{ t('ai.credits') }}</th>
              <th>{{ t('admin.active') }}</th>
              <th class="text-center">{{ t('app.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id" :class="{ 'row-inactive': u.is_active === false }">
              <td class="text-muted fw-medium">{{ u.id }}</td>
              <td>
                <div class="user-info">
                  <div class="user-avatar" :class="'user-avatar--' + u.role">
                    {{ (u.full_name_en || u.full_name_ar || u.email || '?').charAt(0).toUpperCase() }}
                  </div>
                  <span class="user-name">{{ u.full_name_ar || u.full_name_en || '-' }}</span>
                </div>
              </td>
              <td class="text-muted">{{ u.email }}</td>
              <td>
                <select
                  class="form-select role-select"
                  :class="'role-select--' + u.role"
                  :value="u.role"
                  @change="updateRole(u, ($event.target as HTMLSelectElement).value)"
                >
                  <option value="student">{{ t('admin.student') }}</option>
                  <option value="teacher">{{ t('admin.teacher') }}</option>
                  <option value="admin">{{ t('admin.adminRole') }}</option>
                </select>
              </td>
              <td class="text-muted">{{ u.faculty?.name_en || '-' }}</td>
              <td>
                <span class="credit-badge">
                  <i class="fas fa-coins me-1"></i>{{ u.ai_credits }}
                </span>
              </td>
              <td>
                <div class="status-toggle">
                  <div class="form-check form-switch mb-0">
                    <input
                      class="form-check-input custom-switch"
                      type="checkbox"
                      role="switch"
                      :checked="u.is_active !== false"
                      @change="toggleActive(u)"
                    />
                  </div>
                  <span class="status-label" :class="u.is_active !== false ? 'status-label--active' : 'status-label--inactive'">
                    <span class="status-indicator"></span>
                    {{ u.is_active !== false ? t('admin.active') : t('admin.inactive') }}
                  </span>
                </div>
              </td>
              <td>
                <div class="action-buttons">
                  <button @click="openCreditsModal(u)" class="action-btn action-btn--credits" :title="t('admin.updateCredits')">
                    <i class="fas fa-coins"></i>
                  </button>
                  <button @click="confirmDeleteUser(u)" class="action-btn action-btn--delete" :title="t('admin.deleteUser')">
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Credits Modal -->
    <div class="modal fade" :class="{ show: showCreditsModal }" :style="{ display: showCreditsModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-sm modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--primary">
            <div class="modal-header-icon">
              <i class="fas fa-coins"></i>
            </div>
            <h5 class="modal-title">{{ t('admin.updateCredits') }}</h5>
            <button class="btn-close btn-close-white" @click="showCreditsModal = false"></button>
          </div>
          <div class="modal-body p-4">
            <label class="form-label fw-semibold mb-2">AI Credits</label>
            <input type="number" class="form-control admin-input" v-model.number="newCredits" min="0" />
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="showCreditsModal = false">{{ t('app.cancel') || 'Cancel' }}</button>
            <button class="btn admin-btn admin-btn--primary" @click="updateCredits">{{ t('app.save') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showCreditsModal" class="modal-backdrop fade show"></div>

    <!-- Create User Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-lg modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--primary">
            <div class="modal-header-icon">
              <i class="fas fa-user-plus"></i>
            </div>
            <h5 class="modal-title">{{ t('admin.addUser') }}</h5>
            <button class="btn-close btn-close-white" @click="closeCreateModal"></button>
          </div>
          <div class="modal-body p-4">
            <!-- Error display -->
            <div v-if="createError" class="alert alert-danger alert-dismissible fade show admin-alert">
              <i class="fas fa-exclamation-circle me-2"></i>{{ createError }}
              <button type="button" class="btn-close" @click="createError = ''"></button>
            </div>

            <form @submit.prevent="submitCreateUser">
              <!-- Email & Password row -->
              <div class="row g-3 mb-4">
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('auth.email') }} <span class="text-danger">*</span></label>
                  <input type="email" class="form-control admin-input" v-model="newUser.email" required />
                </div>
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('admin.password') }} <span class="text-danger">*</span></label>
                  <div class="input-group">
                    <input
                      :type="showPassword ? 'text' : 'password'"
                      class="form-control admin-input"
                      v-model="newUser.password"
                      required
                      minlength="6"
                    />
                    <button
                      class="btn btn-outline-secondary"
                      type="button"
                      @click="showPassword = !showPassword"
                      style="border-radius: 0 10px 10px 0;"
                    >
                      <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                    </button>
                  </div>
                </div>
              </div>

              <!-- Name fields row -->
              <div class="row g-3 mb-4">
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('auth.fullNameAr') }}</label>
                  <input type="text" class="form-control admin-input" v-model="newUser.full_name_ar" dir="rtl" />
                </div>
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('auth.fullNameEn') }}</label>
                  <input type="text" class="form-control admin-input" v-model="newUser.full_name_en" />
                </div>
              </div>

              <!-- Phone & Role row -->
              <div class="row g-3 mb-4">
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('auth.phone') }}</label>
                  <input type="tel" class="form-control admin-input" v-model="newUser.phone" />
                </div>
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('admin.role') }}</label>
                  <select class="form-select admin-input" v-model="newUser.role">
                    <option value="student">{{ t('admin.student') }}</option>
                    <option value="teacher">{{ t('admin.teacher') }}</option>
                    <option value="admin">{{ t('admin.adminRole') }}</option>
                  </select>
                </div>
              </div>

              <!-- Faculty & Department row -->
              <div class="row g-3">
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('auth.faculty') }}</label>
                  <select class="form-select admin-input" v-model="newUser.faculty_id" @change="onFacultyChange">
                    <option :value="null">{{ t('auth.selectFaculty') }}</option>
                    <option v-for="f in faculties" :key="f.id" :value="f.id">
                      {{ f.name_en || f.name_ar }}
                    </option>
                  </select>
                </div>
                <div class="col-md-6">
                  <label class="form-label fw-semibold">{{ t('auth.department') }}</label>
                  <select class="form-select admin-input" v-model="newUser.department_id" :disabled="!newUser.faculty_id">
                    <option :value="null">{{ t('auth.selectDepartment') }}</option>
                    <option v-for="d in departments" :key="d.id" :value="d.id">
                      {{ d.name_en || d.name_ar }}
                    </option>
                  </select>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="closeCreateModal" :disabled="creating">{{ t('app.cancel') }}</button>
            <button class="btn admin-btn admin-btn--primary" @click="submitCreateUser" :disabled="creating">
              <span v-if="creating" class="spinner-border spinner-border-sm me-1" role="status"></span>
              {{ creating ? t('app.loading') : t('app.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showCreateModal" class="modal-backdrop fade show"></div>

    <!-- Delete Confirmation Modal -->
    <div class="modal fade" :class="{ show: showDeleteModal }" :style="{ display: showDeleteModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-sm modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--danger">
            <div class="modal-header-icon">
              <i class="fas fa-exclamation-triangle"></i>
            </div>
            <h5 class="modal-title">{{ t('admin.deleteUser') }}</h5>
            <button class="btn-close btn-close-white" @click="showDeleteModal = false"></button>
          </div>
          <div class="modal-body p-4">
            <p class="mb-2">{{ t('admin.deleteUserConfirm') }}</p>
            <div class="delete-user-info">
              <i class="fas fa-user me-2"></i>{{ userToDelete?.email }}
            </div>
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="showDeleteModal = false">{{ t('app.cancel') }}</button>
            <button class="btn admin-btn admin-btn--danger" @click="deleteUser" :disabled="deleting">
              <span v-if="deleting" class="spinner-border spinner-border-sm me-1" role="status"></span>
              <i v-else class="fas fa-trash me-1"></i>{{ t('app.delete') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showDeleteModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI, publicAPI } from '../../services/api'
import { useToast } from '../../composables/useToast'

const { t } = useI18n()
const toast = useToast()

const users = ref<any[]>([])
const search = ref('')
const roleFilter = ref('')

// Credits modal
const showCreditsModal = ref(false)
const selectedUserId = ref(0)
const newCredits = ref(0)

// Create user modal
const showCreateModal = ref(false)
const creating = ref(false)
const createError = ref('')
const showPassword = ref(false)
const faculties = ref<any[]>([])
const departments = ref<any[]>([])
const newUser = ref(getEmptyUser())

// Delete modal
const showDeleteModal = ref(false)
const deleting = ref(false)
const userToDelete = ref<any>(null)

let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

function getEmptyUser() {
  return {
    email: '',
    password: '',
    full_name_ar: '',
    full_name_en: '',
    phone: '',
    role: 'student',
    faculty_id: null as number | null,
    department_id: null as number | null,
  }
}

function debouncedFetchUsers() {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    fetchUsers()
  }, 300)
}

async function fetchUsers() {
  try {
    const res = await adminAPI.listUsers({ search: search.value, role: roleFilter.value })
    users.value = res.data.users || []
  } catch {
    // silent
  }
}

// --- Credits ---
function openCreditsModal(u: any) {
  selectedUserId.value = u.id
  newCredits.value = u.ai_credits
  showCreditsModal.value = true
}

async function updateCredits() {
  try {
    await adminAPI.updateCredits(selectedUserId.value, newCredits.value)
    showCreditsModal.value = false
    toast.success(t('admin.userUpdated'))
    fetchUsers()
  } catch {
    toast.error(t('app.error'))
  }
}

// --- Role update ---
async function updateRole(u: any, newRole: string) {
  const oldRole = u.role
  u.role = newRole
  try {
    await adminAPI.updateUser(u.id, { role: newRole })
    toast.success(t('admin.userUpdated'))
  } catch {
    u.role = oldRole
    toast.error(t('app.error'))
  }
}

// --- Active toggle ---
async function toggleActive(u: any) {
  const newValue = u.is_active === false ? true : false
  const oldValue = u.is_active
  u.is_active = newValue
  try {
    await adminAPI.updateUser(u.id, { is_active: newValue })
    toast.success(t('admin.userUpdated'))
  } catch {
    u.is_active = oldValue
    toast.error(t('app.error'))
  }
}

// --- Create user ---
async function loadFaculties() {
  try {
    const res = await publicAPI.getFaculties()
    faculties.value = res.data.faculties || res.data || []
  } catch {
    // silent
  }
}

async function onFacultyChange() {
  newUser.value.department_id = null
  departments.value = []
  if (newUser.value.faculty_id) {
    try {
      const res = await publicAPI.getDepartments(newUser.value.faculty_id)
      departments.value = res.data.departments || res.data || []
    } catch {
      // silent
    }
  }
}

function closeCreateModal() {
  showCreateModal.value = false
  createError.value = ''
  showPassword.value = false
  newUser.value = getEmptyUser()
  departments.value = []
}

async function submitCreateUser() {
  if (!newUser.value.email || !newUser.value.password) return
  creating.value = true
  createError.value = ''
  try {
    await adminAPI.createUser(newUser.value)
    toast.success(t('admin.userCreated'))
    closeCreateModal()
    fetchUsers()
  } catch (err: any) {
    const msg = err.response?.data?.error || err.response?.data?.message || ''
    if (msg.toLowerCase().includes('exist') || msg.toLowerCase().includes('duplicate')) {
      createError.value = t('admin.emailExists')
    } else {
      createError.value = msg || t('app.error')
    }
  } finally {
    creating.value = false
  }
}

// --- Delete user ---
function confirmDeleteUser(u: any) {
  userToDelete.value = u
  showDeleteModal.value = true
}

async function deleteUser() {
  if (!userToDelete.value) return
  deleting.value = true
  try {
    await adminAPI.deleteUser(userToDelete.value.id)
    toast.success(t('admin.userDeleted'))
    showDeleteModal.value = false
    userToDelete.value = null
    fetchUsers()
  } catch {
    toast.error(t('app.error'))
  } finally {
    deleting.value = false
  }
}

// --- Export / Import ---
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

onMounted(() => {
  fetchUsers()
  loadFaculties()
})

onUnmounted(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
})
</script>

<style scoped>
.admin-users {
  max-width: 1400px;
}

/* ── Page Header ── */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 16px;
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

.header-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* ── Admin Buttons ── */
.admin-btn {
  border-radius: 8px;
  font-weight: 600;
  padding: 8px 18px;
  font-size: 13px;
  border: none;
  transition: all 0.2s;
}

.admin-btn--primary {
  background: #1a5276;
  color: #ffffff;
}
.admin-btn--primary:hover {
  background: #154360;
  color: #ffffff;
}

.admin-btn--outline {
  background: #ffffff;
  color: #1a5276;
  border: 1px solid #1a5276;
}
.admin-btn--outline:hover {
  background: #1a5276;
  color: #ffffff;
}

.admin-btn--danger {
  background: #c62828;
  color: #ffffff;
}
.admin-btn--danger:hover {
  background: #b71c1c;
  color: #ffffff;
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
  min-width: 160px;
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
  max-width: 400px;
  margin: 0 auto;
}

/* ── Data Card & Table ── */
.data-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

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
  padding: 12px 16px;
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

.row-inactive {
  opacity: 0.55;
}

/* ── User Info ── */
.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 14px;
  color: #ffffff;
  flex-shrink: 0;
}

.user-avatar--admin { background: linear-gradient(135deg, #c62828, #e53935); }
.user-avatar--teacher { background: linear-gradient(135deg, #1a5276, #2980b9); }
.user-avatar--student { background: linear-gradient(135deg, #0f7b5f, #27ae60); }

.user-name {
  font-weight: 600;
  color: #2c3e50;
}

/* ── Role Select ── */
.role-select {
  width: auto;
  min-width: 110px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 8px;
  padding: 4px 28px 4px 10px;
  height: 32px;
}

.role-select--admin {
  color: #c62828;
  border-color: rgba(198, 40, 40, 0.3);
  background-color: rgba(198, 40, 40, 0.04);
}

.role-select--teacher {
  color: #1a5276;
  border-color: rgba(26, 82, 118, 0.3);
  background-color: rgba(26, 82, 118, 0.04);
}

.role-select--student {
  color: #0f7b5f;
  border-color: rgba(15, 123, 95, 0.3);
  background-color: rgba(15, 123, 95, 0.04);
}

/* ── Credit Badge ── */
.credit-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
}

/* ── Status Toggle ── */
.status-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
}

.custom-switch {
  cursor: pointer;
}

.custom-switch:checked {
  background-color: #1a5276;
  border-color: #1a5276;
}

.status-label {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  font-weight: 500;
}

.status-indicator {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.status-label--active { color: #2e7d32; }
.status-label--active .status-indicator { background: #43a047; }

.status-label--inactive { color: #9e9e9e; }
.status-label--inactive .status-indicator { background: #bdbdbd; }

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
}

.action-btn--credits {
  background: rgba(192, 152, 43, 0.08);
  color: #c0982b;
  border-color: rgba(192, 152, 43, 0.15);
}
.action-btn--credits:hover {
  background: #c0982b;
  color: #fff;
}

.action-btn--delete {
  background: rgba(198, 40, 40, 0.08);
  color: #c62828;
  border-color: rgba(198, 40, 40, 0.15);
}
.action-btn--delete:hover {
  background: #c62828;
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

.admin-modal-header--primary {
  background: linear-gradient(135deg, #1a5276, #2980b9);
}

.admin-modal-header--danger {
  background: linear-gradient(135deg, #c62828, #e53935);
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

.admin-input {
  border-radius: 10px;
  border: 1px solid #dee2e6;
  height: 42px;
  font-size: 14px;
}

.admin-input:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

.admin-alert {
  border-radius: 10px;
  border: none;
  font-size: 14px;
}

.admin-modal-footer {
  border-top: 1px solid #f0f0f0;
  padding: 16px 24px;
  gap: 8px;
}

.delete-user-info {
  padding: 10px 14px;
  background: #fce4ec;
  border-radius: 8px;
  color: #c62828;
  font-size: 14px;
  font-weight: 500;
}
</style>
