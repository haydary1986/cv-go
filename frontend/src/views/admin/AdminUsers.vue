<template>
  <div>
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h3>{{ t('admin.manageUsers') }}</h3>
      <div class="d-flex gap-2">
        <button @click="showCreateModal = true" class="btn btn-primary btn-sm">
          <i class="fas fa-user-plus me-1"></i>{{ t('admin.addUser') }}
        </button>
        <button @click="exportCSV" class="btn btn-outline-success btn-sm">
          <i class="fas fa-download me-1"></i>{{ t('admin.exportCSV') }}
        </button>
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
          <option value="student">{{ t('admin.student') }}</option>
          <option value="teacher">{{ t('admin.teacher') }}</option>
          <option value="admin">{{ t('admin.adminRole') }}</option>
        </select>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="users.length === 0" class="text-center py-5">
      <i class="fas fa-users fa-4x text-muted mb-3 d-block"></i>
      <h5 class="text-muted">{{ t('admin.noUsersFound') }}</h5>
      <p class="text-muted">{{ t('admin.noUsersFoundDesc') }}</p>
    </div>

    <div v-else class="table-responsive">
      <table class="table table-hover align-middle">
        <thead>
          <tr>
            <th>#</th>
            <th>{{ t('admin.name') }}</th>
            <th>{{ t('auth.email') }}</th>
            <th>{{ t('admin.role') }}</th>
            <th>{{ t('auth.faculty') }}</th>
            <th>{{ t('ai.credits') }}</th>
            <th>{{ t('admin.active') }}</th>
            <th>{{ t('app.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id" :class="{ 'opacity-50': u.is_active === false }">
            <td>{{ u.id }}</td>
            <td>{{ u.full_name_ar || u.full_name_en || '-' }}</td>
            <td>{{ u.email }}</td>
            <td>
              <select
                class="form-select form-select-sm d-inline-block"
                style="width: auto; min-width: 100px;"
                :value="u.role"
                @change="updateRole(u, ($event.target as HTMLSelectElement).value)"
                :class="{
                  'text-danger border-danger': u.role === 'admin',
                  'text-primary border-primary': u.role === 'teacher',
                  'text-secondary border-secondary': u.role === 'student',
                }"
              >
                <option value="student">{{ t('admin.student') }}</option>
                <option value="teacher">{{ t('admin.teacher') }}</option>
                <option value="admin">{{ t('admin.adminRole') }}</option>
              </select>
            </td>
            <td>{{ u.faculty?.name_en || '-' }}</td>
            <td>{{ u.ai_credits }}</td>
            <td>
              <div class="form-check form-switch mb-0">
                <input
                  class="form-check-input"
                  type="checkbox"
                  role="switch"
                  :checked="u.is_active !== false"
                  @change="toggleActive(u)"
                  style="cursor: pointer;"
                />
                <label class="form-check-label small" :class="u.is_active !== false ? 'text-success' : 'text-muted'">
                  {{ u.is_active !== false ? t('admin.active') : t('admin.inactive') }}
                </label>
              </div>
            </td>
            <td>
              <div class="btn-group btn-group-sm">
                <button @click="openCreditsModal(u)" class="btn btn-outline-primary" :title="t('admin.updateCredits')">
                  <i class="fas fa-coins"></i>
                </button>
                <button @click="confirmDeleteUser(u)" class="btn btn-outline-danger" :title="t('admin.deleteUser')">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
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

    <!-- Create User Modal -->
    <div class="modal fade" :class="{ show: showCreateModal }" :style="{ display: showCreateModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ t('admin.addUser') }}</h5>
            <button class="btn-close" @click="closeCreateModal"></button>
          </div>
          <div class="modal-body">
            <!-- Error display -->
            <div v-if="createError" class="alert alert-danger alert-dismissible fade show">
              {{ createError }}
              <button type="button" class="btn-close" @click="createError = ''"></button>
            </div>

            <form @submit.prevent="submitCreateUser">
              <!-- Email & Password row -->
              <div class="row g-3 mb-3">
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.email') }} <span class="text-danger">*</span></label>
                  <input type="email" class="form-control" v-model="newUser.email" required />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('admin.password') }} <span class="text-danger">*</span></label>
                  <div class="input-group">
                    <input
                      :type="showPassword ? 'text' : 'password'"
                      class="form-control"
                      v-model="newUser.password"
                      required
                      minlength="6"
                    />
                    <button
                      class="btn btn-outline-secondary"
                      type="button"
                      @click="showPassword = !showPassword"
                    >
                      <i class="fas" :class="showPassword ? 'fa-eye-slash' : 'fa-eye'"></i>
                    </button>
                  </div>
                </div>
              </div>

              <!-- Name fields row -->
              <div class="row g-3 mb-3">
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.fullNameAr') }}</label>
                  <input type="text" class="form-control" v-model="newUser.full_name_ar" dir="rtl" />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.fullNameEn') }}</label>
                  <input type="text" class="form-control" v-model="newUser.full_name_en" />
                </div>
              </div>

              <!-- Phone & Role row -->
              <div class="row g-3 mb-3">
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.phone') }}</label>
                  <input type="tel" class="form-control" v-model="newUser.phone" />
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('admin.role') }}</label>
                  <select class="form-select" v-model="newUser.role">
                    <option value="student">{{ t('admin.student') }}</option>
                    <option value="teacher">{{ t('admin.teacher') }}</option>
                    <option value="admin">{{ t('admin.adminRole') }}</option>
                  </select>
                </div>
              </div>

              <!-- Faculty & Department row -->
              <div class="row g-3 mb-3">
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.faculty') }}</label>
                  <select class="form-select" v-model="newUser.faculty_id" @change="onFacultyChange">
                    <option :value="null">{{ t('auth.selectFaculty') }}</option>
                    <option v-for="f in faculties" :key="f.id" :value="f.id">
                      {{ f.name_en || f.name_ar }}
                    </option>
                  </select>
                </div>
                <div class="col-md-6">
                  <label class="form-label">{{ t('auth.department') }}</label>
                  <select class="form-select" v-model="newUser.department_id" :disabled="!newUser.faculty_id">
                    <option :value="null">{{ t('auth.selectDepartment') }}</option>
                    <option v-for="d in departments" :key="d.id" :value="d.id">
                      {{ d.name_en || d.name_ar }}
                    </option>
                  </select>
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="closeCreateModal" :disabled="creating">{{ t('app.cancel') }}</button>
            <button class="btn btn-primary" @click="submitCreateUser" :disabled="creating">
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
      <div class="modal-dialog modal-sm">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">{{ t('admin.deleteUser') }}</h5>
            <button class="btn-close" @click="showDeleteModal = false"></button>
          </div>
          <div class="modal-body">
            <p>{{ t('admin.deleteUserConfirm') }}</p>
            <p class="text-muted small mb-0">{{ userToDelete?.email }}</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="showDeleteModal = false">{{ t('app.cancel') }}</button>
            <button class="btn btn-danger" @click="deleteUser" :disabled="deleting">
              <span v-if="deleting" class="spinner-border spinner-border-sm me-1" role="status"></span>
              {{ t('app.delete') }}
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
