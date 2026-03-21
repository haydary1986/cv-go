<template>
  <div>
    <h3 class="mb-4">{{ t('admin.branding') }} & {{ t('admin.aiSettings') }}</h3>

    <ul class="nav nav-tabs mb-4">
      <li class="nav-item"><a class="nav-link" :class="{ active: tab === 'branding' }" href="#" @click.prevent="tab = 'branding'">{{ t('admin.branding') }}</a></li>
      <li class="nav-item"><a class="nav-link" :class="{ active: tab === 'ai' }" href="#" @click.prevent="tab = 'ai'">{{ t('admin.aiSettings') }}</a></li>
      <li class="nav-item"><a class="nav-link" :class="{ active: tab === 'ads' }" href="#" @click.prevent="tab = 'ads'">{{ t('admin.adSettings') }}</a></li>
      <li class="nav-item"><a class="nav-link" :class="{ active: tab === 'notif' }" href="#" @click.prevent="tab = 'notif'">{{ t('admin.sendNotification') }}</a></li>
    </ul>

    <!-- Branding -->
    <div v-show="tab === 'branding'">
      <!-- Live Preview -->
      <div class="card mb-4">
        <div class="card-header bg-light">
          <h6 class="mb-0">{{ t('admin.livePreview') }}</h6>
        </div>
        <div class="card-body p-0">
          <div
            class="d-flex align-items-center px-4 py-3"
            :style="{ backgroundColor: branding.primary_color, color: '#fff' }"
          >
            <img
              v-if="branding.logo_url"
              :src="branding.logo_url"
              alt="Logo"
              class="me-3"
              style="height: 36px; width: auto; object-fit: contain; border-radius: 4px;"
            />
            <div
              v-else
              class="me-3 d-flex align-items-center justify-content-center"
              style="height: 36px; width: 36px; background: rgba(255,255,255,0.2); border-radius: 4px; font-size: 18px; font-weight: bold;"
            >
              {{ branding.name ? branding.name.charAt(0).toUpperCase() : 'C' }}
            </div>
            <span class="fw-semibold fs-5">{{ branding.name || t('app.name') }}</span>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h5 class="mb-0">{{ t('admin.brandingSettings') }}</h5>
        </div>
        <div class="card-body">
          <!-- System Name -->
          <div class="mb-4">
            <label class="form-label fw-semibold">{{ t('admin.systemName') }}</label>
            <input type="text" class="form-control" v-model="branding.name" />
          </div>

          <!-- Logo Upload -->
          <div class="mb-4">
            <label class="form-label fw-semibold">{{ t('admin.systemLogo') }}</label>

            <!-- Current Logo Preview -->
            <div v-if="branding.logo_url" class="mb-3 d-flex align-items-center gap-3">
              <div class="border rounded p-2" style="background: #f8f9fa;">
                <img
                  :src="branding.logo_url"
                  alt="Current Logo"
                  style="max-height: 64px; max-width: 200px; object-fit: contain;"
                />
              </div>
              <button
                type="button"
                class="btn btn-outline-danger btn-sm"
                @click="removeLogo"
              >
                <i class="bi bi-trash me-1"></i>{{ t('admin.removeLogo') }}
              </button>
            </div>

            <!-- Upload Zone -->
            <div
              class="upload-zone border rounded-3 p-4 text-center"
              :class="{ 'drag-over': isDragging }"
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop.prevent="handleDrop"
              @click="($refs.logoInput as HTMLInputElement).click()"
              style="cursor: pointer; border-style: dashed !important; transition: all 0.2s;"
            >
              <input
                ref="logoInput"
                type="file"
                class="d-none"
                accept="image/jpeg,image/png,image/gif,image/svg+xml,image/webp"
                @change="uploadLogo"
              />
              <div v-if="logoUploading" class="py-2">
                <div class="spinner-border spinner-border-sm text-primary me-2" role="status"></div>
                <span class="text-muted">{{ t('app.loading') }}</span>
              </div>
              <div v-else>
                <i class="bi bi-cloud-arrow-up fs-1 text-muted d-block mb-2"></i>
                <p class="mb-1 text-muted">{{ t('admin.dragDropLogo') }}</p>
                <small class="text-muted">{{ t('admin.logoHint') }}</small>
              </div>
            </div>
          </div>

          <!-- Color Settings -->
          <div class="mb-3">
            <label class="form-label fw-semibold">{{ t('admin.colorSettings') }}</label>
            <div class="row g-3">
              <div class="col-md-4">
                <label class="form-label text-muted small">{{ t('admin.primaryColor') }}</label>
                <div class="d-flex align-items-center gap-2">
                  <input
                    type="color"
                    class="form-control form-control-color"
                    v-model="branding.primary_color"
                    style="min-width: 48px; height: 40px;"
                  />
                  <div
                    class="rounded flex-grow-1 d-flex align-items-center justify-content-center"
                    :style="{ backgroundColor: branding.primary_color, color: '#fff', height: '40px', fontSize: '13px', fontFamily: 'monospace' }"
                  >
                    {{ branding.primary_color }}
                  </div>
                </div>
              </div>
              <div class="col-md-4">
                <label class="form-label text-muted small">{{ t('admin.secondaryColor') }}</label>
                <div class="d-flex align-items-center gap-2">
                  <input
                    type="color"
                    class="form-control form-control-color"
                    v-model="branding.secondary_color"
                    style="min-width: 48px; height: 40px;"
                  />
                  <div
                    class="rounded flex-grow-1 d-flex align-items-center justify-content-center"
                    :style="{ backgroundColor: branding.secondary_color, color: '#fff', height: '40px', fontSize: '13px', fontFamily: 'monospace' }"
                  >
                    {{ branding.secondary_color }}
                  </div>
                </div>
              </div>
              <div class="col-md-4">
                <label class="form-label text-muted small">{{ t('admin.accentColor') }}</label>
                <div class="d-flex align-items-center gap-2">
                  <input
                    type="color"
                    class="form-control form-control-color"
                    v-model="branding.accent_color"
                    style="min-width: 48px; height: 40px;"
                  />
                  <div
                    class="rounded flex-grow-1 d-flex align-items-center justify-content-center"
                    :style="{ backgroundColor: branding.accent_color, color: '#fff', height: '40px', fontSize: '13px', fontFamily: 'monospace' }"
                  >
                    {{ branding.accent_color }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <button
            @click="saveBranding"
            class="btn btn-primary mt-3"
            :disabled="savingBranding"
          >
            <span v-if="savingBranding">
              <span class="spinner-border spinner-border-sm me-1" role="status"></span>
              {{ t('app.loading') }}
            </span>
            <span v-else>{{ t('app.save') }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- AI Settings -->
    <div v-show="tab === 'ai'" class="card">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.provider') }}</label>
            <select class="form-select" v-model="aiForm.provider">
              <option value="openai">OpenAI</option>
              <option value="gemini">Google Gemini</option>
              <option value="deepseek">DeepSeek</option>
            </select>
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.model') }}</label>
            <select class="form-select" v-model="aiForm.model">
              <optgroup v-if="aiForm.provider === 'openai'" label="OpenAI">
                <option value="gpt-4o">GPT-4o</option>
                <option value="gpt-4o-mini">GPT-4o Mini</option>
                <option value="gpt-4-turbo">GPT-4 Turbo</option>
                <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
              </optgroup>
              <optgroup v-if="aiForm.provider === 'gemini'" label="Gemini">
                <option value="gemini-1.5-pro">Gemini 1.5 Pro</option>
                <option value="gemini-1.5-flash">Gemini 1.5 Flash</option>
              </optgroup>
              <optgroup v-if="aiForm.provider === 'deepseek'" label="DeepSeek">
                <option value="deepseek-chat">DeepSeek Chat</option>
                <option value="deepseek-coder">DeepSeek Coder</option>
              </optgroup>
            </select>
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.maxTokens') }}</label>
            <input type="number" class="form-control" v-model.number="aiForm.max_tokens" />
          </div>
          <div class="col-md-8">
            <label class="form-label">{{ t('admin.apiKey') }}</label>
            <input type="password" class="form-control" v-model="aiForm.api_key" :placeholder="t('admin.apiKey')" />
          </div>
          <div class="col-md-4 d-flex align-items-end">
            <div class="form-check form-switch">
              <input type="checkbox" class="form-check-input" v-model="aiForm.is_active" />
              <label class="form-check-label">{{ t('app.active') }}</label>
            </div>
          </div>
        </div>
        <button @click="saveAI" class="btn btn-primary mt-3">{{ t('app.save') }}</button>
      </div>
    </div>

    <!-- Ad Settings -->
    <div v-show="tab === 'ads'" class="card">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.adType') }}</label>
            <select class="form-select" v-model="adSettings.type">
              <option value="video">Video</option>
              <option value="adsense">Google AdSense</option>
            </select>
          </div>
          <div class="col-md-8" v-if="adSettings.type === 'video'">
            <label class="form-label">{{ t('admin.videoUrl') }}</label>
            <input type="url" class="form-control" v-model="adSettings.video_url" />
          </div>
          <div class="col-md-8" v-else>
            <label class="form-label">{{ t('admin.adsenseCode') }}</label>
            <textarea class="form-control" rows="3" v-model="adSettings.adsense_code"></textarea>
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.duration') }}</label>
            <input type="number" class="form-control" v-model.number="adSettings.duration" />
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.skipAfter') }}</label>
            <input type="number" class="form-control" v-model.number="adSettings.skip_after" />
          </div>
          <div class="col-md-4 d-flex align-items-end">
            <div class="form-check form-switch">
              <input type="checkbox" class="form-check-input" v-model="adSettings.is_active" />
              <label class="form-check-label">{{ t('app.active') }}</label>
            </div>
          </div>
        </div>
        <button @click="saveAds" class="btn btn-primary mt-3">{{ t('app.save') }}</button>
      </div>
    </div>

    <!-- Send Notification -->
    <div v-show="tab === 'notif'" class="card">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-6">
            <label class="form-label">{{ t('notifications.title') }} (AR)</label>
            <input type="text" class="form-control" v-model="notifForm.title_ar" />
          </div>
          <div class="col-md-6">
            <label class="form-label">{{ t('notifications.title') }} (EN)</label>
            <input type="text" class="form-control" v-model="notifForm.title_en" />
          </div>
          <div class="col-md-6">
            <label class="form-label">{{ t('admin.message') }} (AR)</label>
            <textarea class="form-control" rows="3" v-model="notifForm.message_ar"></textarea>
          </div>
          <div class="col-md-6">
            <label class="form-label">{{ t('admin.message') }} (EN)</label>
            <textarea class="form-control" rows="3" v-model="notifForm.message_en"></textarea>
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.target') }}</label>
            <select class="form-select" v-model="notifForm.target" @change="onTargetChange">
              <option value="all">{{ t('app.allUsers') }}</option>
              <option value="faculty">{{ t('admin.byFaculty') }}</option>
              <option value="department">{{ t('admin.byDepartment') }}</option>
            </select>
          </div>
          <div class="col-md-4" v-if="notifForm.target === 'faculty' || notifForm.target === 'department'">
            <label class="form-label">{{ t('cv.selectFaculty') }}</label>
            <select class="form-select" v-model="notifForm.faculty_id" @change="onNotifFacultyChange">
              <option :value="null" disabled>-- {{ t('cv.selectFaculty') }} --</option>
              <option v-for="f in faculties" :key="f.id" :value="f.id">{{ f.name_ar }} - {{ f.name_en }}</option>
            </select>
          </div>
          <div class="col-md-4" v-if="notifForm.target === 'department' && notifForm.faculty_id">
            <label class="form-label">{{ t('cv.selectDepartment') }}</label>
            <select class="form-select" v-model="notifForm.department_id">
              <option :value="null" disabled>-- {{ t('cv.selectDepartment') }} --</option>
              <option v-for="d in notifDepartments" :key="d.id" :value="d.id">{{ d.name_ar }} - {{ d.name_en }}</option>
            </select>
          </div>
        </div>
        <button @click="sendNotif" class="btn btn-primary mt-3" :disabled="sendingNotif">
          <span v-if="sendingNotif" class="spinner-border spinner-border-sm me-1"></span>
          {{ t('admin.sendNotification') }}
        </button>
      </div>
    </div>

    <div v-if="message" class="alert alert-success mt-3">{{ message }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api, { adminAPI, publicAPI } from '../../services/api'
import { useToast } from '../../composables/useToast'

const { t } = useI18n()
const toast = useToast()
const tab = ref('branding')
const message = ref('')

const branding = reactive({ name: '', logo_url: '', primary_color: '#0d6efd', secondary_color: '#6c757d', accent_color: '#198754' })
const aiForm = reactive({ provider: 'openai', model: 'gpt-4o-mini', api_key: '', max_tokens: 2000, is_active: false })
const adSettings = reactive({ type: 'video', video_url: '', adsense_code: '', duration: 5, skip_after: 3, is_active: false })
const notifForm = reactive({ title_ar: '', title_en: '', message_ar: '', message_en: '', target: 'all', faculty_id: null as number | null, department_id: null as number | null })
const faculties = ref<any[]>([])
const notifDepartments = ref<any[]>([])
const sendingNotif = ref(false)

const isDragging = ref(false)
const logoUploading = ref(false)
const savingBranding = ref(false)

onMounted(async () => {
  try {
    const [bRes, aiRes, adRes, fRes] = await Promise.all([
      adminAPI.getBranding(),
      adminAPI.getAISettings(),
      adminAPI.getAdSettings(),
      publicAPI.getFaculties(),
    ])
    if (bRes.data.branding) Object.assign(branding, bRes.data.branding)
    if (aiRes.data.settings) Object.assign(aiForm, aiRes.data.settings)
    if (adRes.data.settings) Object.assign(adSettings, adRes.data.settings)
    if (fRes.data.faculties) faculties.value = fRes.data.faculties
  } catch {}
})

function onTargetChange() {
  notifForm.faculty_id = null
  notifForm.department_id = null
  notifDepartments.value = []
}

async function onNotifFacultyChange() {
  notifForm.department_id = null
  notifDepartments.value = []
  if (notifForm.faculty_id) {
    try {
      const res = await publicAPI.getDepartments(notifForm.faculty_id)
      notifDepartments.value = res.data.departments || []
    } catch {}
  }
}

async function uploadLogo(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return
  const file = input.files[0]
  if (file.size > 2 * 1024 * 1024) {
    toast.error('File too large (max 2MB)')
    return
  }
  const formData = new FormData()
  formData.append('logo', file)
  logoUploading.value = true
  try {
    const res = await api.post('/admin/branding/logo', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    branding.logo_url = res.data.url
    toast.success('Logo uploaded')
  } catch {
    toast.error('Upload failed')
  } finally {
    logoUploading.value = false
    // Reset file input so the same file can be re-selected
    input.value = ''
  }
}

function handleDrop(event: DragEvent) {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (!files?.length) return
  const file = files[0]
  if (!['image/jpeg', 'image/png', 'image/gif', 'image/svg+xml', 'image/webp'].includes(file.type)) {
    toast.error('Invalid file type')
    return
  }
  // Create a synthetic event-like call through the upload function
  if (file.size > 2 * 1024 * 1024) {
    toast.error('File too large (max 2MB)')
    return
  }
  const formData = new FormData()
  formData.append('logo', file)
  logoUploading.value = true
  api.post('/admin/branding/logo', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }).then((res) => {
    branding.logo_url = res.data.url
    toast.success('Logo uploaded')
  }).catch(() => {
    toast.error('Upload failed')
  }).finally(() => {
    logoUploading.value = false
  })
}

function removeLogo() {
  branding.logo_url = ''
}

async function saveBranding() {
  savingBranding.value = true
  try {
    await adminAPI.updateBranding(branding)
    toast.success(t('admin.savedSuccess'))
  } catch {
    toast.error(t('app.error'))
  } finally {
    savingBranding.value = false
  }
}

async function saveAI() { await adminAPI.updateAISettings(aiForm); message.value = t('app.success') }
async function saveAds() { await adminAPI.updateAdSettings(adSettings); message.value = t('app.success') }
async function sendNotif() {
  sendingNotif.value = true
  try {
    const res = await adminAPI.sendNotification(notifForm)
    toast.success(res.data.message || t('app.success'))
    notifForm.title_ar = ''
    notifForm.title_en = ''
    notifForm.message_ar = ''
    notifForm.message_en = ''
  } catch {
    toast.error(t('app.error'))
  } finally {
    sendingNotif.value = false
  }
}
</script>

<style scoped>
.upload-zone {
  border-color: #dee2e6;
  background-color: #fafbfc;
}
.upload-zone:hover,
.upload-zone.drag-over {
  border-color: #0d6efd;
  background-color: #f0f4ff;
}
</style>
