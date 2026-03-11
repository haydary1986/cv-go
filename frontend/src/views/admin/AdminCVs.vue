<template>
  <div>
    <h3 class="mb-4">{{ t('admin.manageCVs') }}</h3>

    <!-- Filters -->
    <div class="row g-2 mb-3">
      <div class="col-md-4">
        <input type="text" class="form-control" :placeholder="t('app.search')" v-model="search" @input="fetchCVs" />
      </div>
      <div class="col-md-3">
        <select class="form-select" v-model="statusFilter" @change="fetchCVs">
          <option value="">{{ t('admin.byStatus') }}</option>
          <option value="draft">{{ t('cv.draft') }}</option>
          <option value="pending">{{ t('cv.pending') }}</option>
          <option value="approved">{{ t('cv.approved') }}</option>
          <option value="rejected">{{ t('cv.rejected') }}</option>
        </select>
      </div>
    </div>

    <!-- Table -->
    <div class="table-responsive">
      <table class="table table-hover">
        <thead>
          <tr>
            <th>#</th>
            <th>{{ t('cv.title') }}</th>
            <th>{{ t('auth.email') }}</th>
            <th>{{ t('cv.language') }}</th>
            <th>{{ t('cv.template') }}</th>
            <th>{{ t('cv.status') }}</th>
            <th>{{ t('cv.createdAt') }}</th>
            <th>{{ t('app.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cv in cvs" :key="cv.id">
            <td>{{ cv.id }}</td>
            <td>
              {{ cv.title }}
              <span v-if="cv.is_guest" class="badge bg-info ms-1">{{ t('cv.guestMode') }}</span>
            </td>
            <td>
              <template v-if="cv.is_guest">
                <span class="text-info">{{ cv.guest_name }}</span>
                <br v-if="cv.guest_email" /><small class="text-muted">{{ cv.guest_email }}</small>
                <br /><small class="text-muted">IP: {{ cv.guest_ip }}</small>
              </template>
              <template v-else>{{ cv.user?.email }}</template>
            </td>
            <td>{{ cv.language }}</td>
            <td>{{ cv.template }}</td>
            <td><span :class="statusBadge(cv.status)" class="badge">{{ cv.status }}</span></td>
            <td>{{ new Date(cv.created_at).toLocaleDateString() }}</td>
            <td>
              <div class="btn-group btn-group-sm">
                <router-link :to="cv.is_guest ? `/shared/${cv.share_token}` : `/cv/${cv.id}`" class="btn btn-outline-primary"><i class="fas fa-eye"></i></router-link>
                <button v-if="cv.status !== 'approved'" @click="approve(cv.id)" class="btn btn-outline-success"><i class="fas fa-check"></i></button>
                <button @click="openReject(cv.id)" class="btn btn-outline-danger"><i class="fas fa-times"></i></button>
                <button @click="openRevision(cv.id)" class="btn btn-outline-warning"><i class="fas fa-edit"></i></button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Reject Modal -->
    <div class="modal fade" :class="{ show: showRejectModal }" :style="{ display: showRejectModal ? 'block' : 'none' }">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5>{{ t('admin.reject') }}</h5>
            <button class="btn-close" @click="showRejectModal = false"></button>
          </div>
          <div class="modal-body">
            <textarea class="form-control" rows="3" :placeholder="t('admin.rejectReason')" v-model="rejectReason"></textarea>
          </div>
          <div class="modal-footer">
            <button class="btn btn-danger" @click="reject">{{ t('admin.reject') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showRejectModal" class="modal-backdrop fade show"></div>

    <!-- Revision Modal -->
    <div class="modal fade" :class="{ show: showRevisionModal }" :style="{ display: showRevisionModal ? 'block' : 'none' }">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5>{{ t('admin.requestRevision') }}</h5>
            <button class="btn-close" @click="showRevisionModal = false"></button>
          </div>
          <div class="modal-body">
            <textarea class="form-control" rows="3" :placeholder="t('admin.revisionNote')" v-model="revisionNote"></textarea>
          </div>
          <div class="modal-footer">
            <button class="btn btn-warning" @click="requestRevision">{{ t('admin.requestRevision') }}</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showRevisionModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'

const { t } = useI18n()
const cvs = ref<any[]>([])
const search = ref('')
const statusFilter = ref('')
const showRejectModal = ref(false)
const showRevisionModal = ref(false)
const selectedCVId = ref(0)
const rejectReason = ref('')
const revisionNote = ref('')

function statusBadge(s: string) {
  return { draft: 'bg-secondary', pending: 'bg-warning text-dark', approved: 'bg-success', rejected: 'bg-danger' }[s] || 'bg-secondary'
}

async function fetchCVs() {
  const res = await adminAPI.listCVs({ search: search.value, status: statusFilter.value })
  cvs.value = res.data.cvs || []
}

async function approve(id: number) {
  await adminAPI.approveCV(id)
  fetchCVs()
}

function openReject(id: number) { selectedCVId.value = id; rejectReason.value = ''; showRejectModal.value = true }
function openRevision(id: number) { selectedCVId.value = id; revisionNote.value = ''; showRevisionModal.value = true }

async function reject() {
  await adminAPI.rejectCV(selectedCVId.value, rejectReason.value)
  showRejectModal.value = false
  fetchCVs()
}

async function requestRevision() {
  await adminAPI.requestRevision(selectedCVId.value, revisionNote.value)
  showRevisionModal.value = false
  fetchCVs()
}

onMounted(fetchCVs)
</script>
