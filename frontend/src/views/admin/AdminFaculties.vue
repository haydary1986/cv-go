<template>
  <div>
    <h3 class="mb-4">{{ t('admin.manageFaculties') }}</h3>

    <div class="row g-4">
      <!-- Faculties -->
      <div class="col-md-6">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">{{ t('admin.manageFaculties') }}</h5>
            <button @click="showFacultyModal = true; editFaculty = null; facultyForm = { name_ar: '', name_en: '' }" class="btn btn-sm btn-primary"><i class="fas fa-plus"></i></button>
          </div>
          <div class="card-body">
            <div v-for="f in faculties" :key="f.id" class="d-flex justify-content-between align-items-center border-bottom py-2">
              <div>
                <strong>{{ f.name_ar }}</strong> / {{ f.name_en }}
                <br/><small class="text-muted">{{ f.departments?.length || 0 }} {{ t('admin.manageDepartments') }}</small>
              </div>
              <div class="btn-group btn-group-sm">
                <button @click="editFac(f)" class="btn btn-outline-primary"><i class="fas fa-edit"></i></button>
                <button @click="deleteFac(f.id)" class="btn btn-outline-danger"><i class="fas fa-trash"></i></button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Departments -->
      <div class="col-md-6">
        <div class="card">
          <div class="card-header d-flex justify-content-between align-items-center">
            <h5 class="mb-0">{{ t('admin.manageDepartments') }}</h5>
            <button @click="showDeptModal = true; editDept = null; deptForm = { faculty_id: 0, name_ar: '', name_en: '' }" class="btn btn-sm btn-primary"><i class="fas fa-plus"></i></button>
          </div>
          <div class="card-body">
            <div v-for="d in departments" :key="d.id" class="d-flex justify-content-between align-items-center border-bottom py-2">
              <div>
                <strong>{{ d.name_ar }}</strong> / {{ d.name_en }}
                <br/><small class="text-muted">{{ d.faculty?.name_en }}</small>
              </div>
              <div class="btn-group btn-group-sm">
                <button @click="editDepartment(d)" class="btn btn-outline-primary"><i class="fas fa-edit"></i></button>
                <button @click="deleteDept(d.id)" class="btn btn-outline-danger"><i class="fas fa-trash"></i></button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Faculty Modal -->
    <div class="modal fade" :class="{ show: showFacultyModal }" :style="{ display: showFacultyModal ? 'block' : 'none' }">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5>{{ editFaculty ? t('app.edit') : t('admin.addFaculty') }}</h5>
            <button class="btn-close" @click="showFacultyModal = false"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">{{ t('admin.facultyNameAr') }}</label>
              <input type="text" class="form-control" v-model="facultyForm.name_ar" />
            </div>
            <div class="mb-3">
              <label class="form-label">{{ t('admin.facultyNameEn') }}</label>
              <input type="text" class="form-control" v-model="facultyForm.name_en" />
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-primary" @click="saveFaculty">{{ t('app.save') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showFacultyModal" class="modal-backdrop fade show"></div>

    <!-- Department Modal -->
    <div class="modal fade" :class="{ show: showDeptModal }" :style="{ display: showDeptModal ? 'block' : 'none' }">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5>{{ editDept ? t('app.edit') : t('admin.addDepartment') }}</h5>
            <button class="btn-close" @click="showDeptModal = false"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label class="form-label">{{ t('auth.faculty') }}</label>
              <select class="form-select" v-model="deptForm.faculty_id">
                <option v-for="f in faculties" :key="f.id" :value="f.id">{{ f.name_en }}</option>
              </select>
            </div>
            <div class="mb-3">
              <label class="form-label">{{ t('admin.departmentNameAr') }}</label>
              <input type="text" class="form-control" v-model="deptForm.name_ar" />
            </div>
            <div class="mb-3">
              <label class="form-label">{{ t('admin.departmentNameEn') }}</label>
              <input type="text" class="form-control" v-model="deptForm.name_en" />
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-primary" @click="saveDept">{{ t('app.save') }}</button>
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
