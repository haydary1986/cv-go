<template>
  <div class="admin-faculties">
    <!-- Page Header -->
    <div class="page-header">
      <h2 class="page-title">{{ t('admin.manageFaculties') }}</h2>
      <p class="page-subtitle">Organize faculties and departments</p>
    </div>

    <div class="row g-4">
      <!-- Faculties Panel -->
      <div class="col-lg-6">
        <div class="panel-card">
          <div class="panel-card-header">
            <div class="panel-card-title">
              <i class="fas fa-university me-2"></i>{{ t('admin.manageFaculties') }}
            </div>
            <button
              @click="showFacultyModal = true; editFaculty = null; Object.assign(facultyForm, { name_ar: '', name_en: '' })"
              class="panel-add-btn"
            >
              <i class="fas fa-plus"></i>
            </button>
          </div>
          <div class="panel-card-body">
            <div v-for="f in faculties" :key="f.id" class="faculty-item">
              <div class="faculty-item-info">
                <div class="faculty-icon">
                  <i class="fas fa-building"></i>
                </div>
                <div>
                  <div class="faculty-name">{{ f.name_ar }} <span class="text-muted">/</span> {{ f.name_en }}</div>
                  <div class="faculty-meta">
                    <span class="dept-count-badge">
                      <i class="fas fa-sitemap me-1"></i>{{ f.departments?.length || 0 }} {{ t('admin.manageDepartments') }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="faculty-item-actions">
                <button @click="editFac(f)" class="action-btn action-btn--edit">
                  <i class="fas fa-edit"></i>
                </button>
                <button @click="deleteFac(f.id)" class="action-btn action-btn--delete">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
            <div v-if="faculties.length === 0" class="empty-panel">
              <i class="fas fa-university fa-2x mb-2"></i>
              <p class="mb-0">No faculties yet</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Departments Panel -->
      <div class="col-lg-6">
        <div class="panel-card">
          <div class="panel-card-header panel-card-header--secondary">
            <div class="panel-card-title">
              <i class="fas fa-sitemap me-2"></i>{{ t('admin.manageDepartments') }}
            </div>
            <button
              @click="showDeptModal = true; editDept = null; Object.assign(deptForm, { faculty_id: 0, name_ar: '', name_en: '' })"
              class="panel-add-btn"
            >
              <i class="fas fa-plus"></i>
            </button>
          </div>
          <div class="panel-card-body">
            <div v-for="d in departments" :key="d.id" class="faculty-item">
              <div class="faculty-item-info">
                <div class="dept-icon">
                  <i class="fas fa-folder"></i>
                </div>
                <div>
                  <div class="faculty-name">{{ d.name_ar }} <span class="text-muted">/</span> {{ d.name_en }}</div>
                  <div class="faculty-meta">
                    <span class="parent-badge">
                      <i class="fas fa-university me-1"></i>{{ d.faculty?.name_en }}
                    </span>
                  </div>
                </div>
              </div>
              <div class="faculty-item-actions">
                <button @click="editDepartment(d)" class="action-btn action-btn--edit">
                  <i class="fas fa-edit"></i>
                </button>
                <button @click="deleteDept(d.id)" class="action-btn action-btn--delete">
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
            <div v-if="departments.length === 0" class="empty-panel">
              <i class="fas fa-sitemap fa-2x mb-2"></i>
              <p class="mb-0">No departments yet</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Faculty Modal -->
    <div class="modal fade" :class="{ show: showFacultyModal }" :style="{ display: showFacultyModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--primary">
            <div class="modal-header-icon">
              <i class="fas fa-university"></i>
            </div>
            <h5 class="modal-title">{{ editFaculty ? t('app.edit') : t('admin.addFaculty') }}</h5>
            <button class="btn-close btn-close-white" @click="showFacultyModal = false"></button>
          </div>
          <div class="modal-body p-4">
            <div class="mb-4">
              <label class="form-label fw-semibold">{{ t('admin.facultyNameAr') }}</label>
              <input type="text" class="form-control admin-input" v-model="facultyForm.name_ar" dir="rtl" />
            </div>
            <div>
              <label class="form-label fw-semibold">{{ t('admin.facultyNameEn') }}</label>
              <input type="text" class="form-control admin-input" v-model="facultyForm.name_en" />
            </div>
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="showFacultyModal = false">{{ t('app.cancel') || 'Cancel' }}</button>
            <button class="btn admin-btn admin-btn--primary" @click="saveFaculty">
              <i class="fas fa-save me-1"></i>{{ t('app.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showFacultyModal" class="modal-backdrop fade show"></div>

    <!-- Department Modal -->
    <div class="modal fade" :class="{ show: showDeptModal }" :style="{ display: showDeptModal ? 'block' : 'none' }">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content admin-modal">
          <div class="modal-header admin-modal-header admin-modal-header--secondary">
            <div class="modal-header-icon">
              <i class="fas fa-sitemap"></i>
            </div>
            <h5 class="modal-title">{{ editDept ? t('app.edit') : t('admin.addDepartment') }}</h5>
            <button class="btn-close btn-close-white" @click="showDeptModal = false"></button>
          </div>
          <div class="modal-body p-4">
            <div class="mb-4">
              <label class="form-label fw-semibold">{{ t('auth.faculty') }}</label>
              <select class="form-select admin-input" v-model="deptForm.faculty_id">
                <option v-for="f in faculties" :key="f.id" :value="f.id">{{ f.name_en }}</option>
              </select>
            </div>
            <div class="mb-4">
              <label class="form-label fw-semibold">{{ t('admin.departmentNameAr') }}</label>
              <input type="text" class="form-control admin-input" v-model="deptForm.name_ar" dir="rtl" />
            </div>
            <div>
              <label class="form-label fw-semibold">{{ t('admin.departmentNameEn') }}</label>
              <input type="text" class="form-control admin-input" v-model="deptForm.name_en" />
            </div>
          </div>
          <div class="modal-footer admin-modal-footer">
            <button class="btn btn-light" @click="showDeptModal = false">{{ t('app.cancel') || 'Cancel' }}</button>
            <button class="btn admin-btn admin-btn--primary" @click="saveDept">
              <i class="fas fa-save me-1"></i>{{ t('app.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showDeptModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'

const { t } = useI18n()
const faculties = ref<any[]>([])
const departments = ref<any[]>([])
const showFacultyModal = ref(false)
const showDeptModal = ref(false)
const editFaculty = ref<any>(null)
const editDept = ref<any>(null)
const facultyForm = reactive({ name_ar: '', name_en: '' })
const deptForm = reactive({ faculty_id: 0, name_ar: '', name_en: '' })

async function loadData() {
  const [fRes, dRes] = await Promise.all([adminAPI.listFaculties(), adminAPI.listDepartments()])
  faculties.value = fRes.data.faculties || []
  departments.value = dRes.data.departments || []
}

function editFac(f: any) {
  editFaculty.value = f
  facultyForm.name_ar = f.name_ar
  facultyForm.name_en = f.name_en
  showFacultyModal.value = true
}

function editDepartment(d: any) {
  editDept.value = d
  deptForm.faculty_id = d.faculty_id
  deptForm.name_ar = d.name_ar
  deptForm.name_en = d.name_en
  showDeptModal.value = true
}

async function saveFaculty() {
  if (editFaculty.value) {
    await adminAPI.updateFaculty(editFaculty.value.id, facultyForm)
  } else {
    await adminAPI.createFaculty(facultyForm)
  }
  showFacultyModal.value = false
  loadData()
}

async function deleteFac(id: number) {
  if (confirm(t('app.deleteConfirm'))) { await adminAPI.deleteFaculty(id); loadData() }
}

async function saveDept() {
  if (editDept.value) {
    await adminAPI.updateDepartment(editDept.value.id, deptForm)
  } else {
    await adminAPI.createDepartment(deptForm)
  }
  showDeptModal.value = false
  loadData()
}

async function deleteDept(id: number) {
  if (confirm(t('app.deleteConfirm'))) { await adminAPI.deleteDepartment(id); loadData() }
}

onMounted(loadData)
</script>

<style scoped>
.admin-faculties {
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

/* ── Panel Card ── */
.panel-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  overflow: hidden;
}

.panel-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  background: linear-gradient(135deg, #1a5276, #2980b9);
  color: #ffffff;
}

.panel-card-header--secondary {
  background: linear-gradient(135deg, #2c3e50, #4a6274);
}

.panel-card-title {
  font-size: 15px;
  font-weight: 600;
}

.panel-add-btn {
  width: 34px;
  height: 34px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.15);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.panel-add-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.panel-card-body {
  padding: 8px;
}

/* ── Faculty Item ── */
.faculty-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-radius: 8px;
  transition: background 0.15s;
  margin: 2px 0;
}

.faculty-item:hover {
  background: #f8fafc;
}

.faculty-item + .faculty-item {
  border-top: 1px solid #f0f2f5;
}

.faculty-item-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.faculty-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(26, 82, 118, 0.1), rgba(26, 82, 118, 0.05));
  color: #1a5276;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.dept-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(44, 62, 80, 0.1), rgba(44, 62, 80, 0.05));
  color: #2c3e50;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  flex-shrink: 0;
}

.faculty-name {
  font-weight: 600;
  font-size: 14px;
  color: #2c3e50;
}

.faculty-meta {
  margin-top: 2px;
}

.dept-count-badge {
  font-size: 12px;
  color: #1a5276;
  background: rgba(26, 82, 118, 0.06);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.parent-badge {
  font-size: 12px;
  color: #6c757d;
  font-weight: 500;
}

.faculty-item-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

/* ── Action Buttons ── */
.action-btn {
  width: 30px;
  height: 30px;
  border-radius: 6px;
  border: 1px solid transparent;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn--edit {
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
  border-color: rgba(26, 82, 118, 0.15);
}
.action-btn--edit:hover { background: #1a5276; color: #fff; }

.action-btn--delete {
  background: rgba(198, 40, 40, 0.08);
  color: #c62828;
  border-color: rgba(198, 40, 40, 0.15);
}
.action-btn--delete:hover { background: #c62828; color: #fff; }

/* ── Empty Panel ── */
.empty-panel {
  text-align: center;
  padding: 40px 20px;
  color: #adb5bd;
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

.admin-modal-header--secondary {
  background: linear-gradient(135deg, #2c3e50, #4a6274);
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

.admin-btn--primary {
  background: #1a5276;
}
.admin-btn--primary:hover {
  background: #154360;
  color: #ffffff;
}
</style>
