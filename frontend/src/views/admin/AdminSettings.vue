<template>
  <div class="admin-settings">
    <!-- Page Header -->
    <div class="page-header">
      <h2 class="page-title">{{ t('admin.branding') }} & {{ t('admin.aiSettings') }}</h2>
      <p class="page-subtitle">Configure system appearance, AI, ads, and notifications</p>
    </div>

    <!-- Tab Navigation -->
    <div class="settings-tabs">
      <button
        class="settings-tab"
        :class="{ active: tab === 'branding' }"
        @click="tab = 'branding'"
      >
        <i class="fas fa-palette me-2"></i>{{ t('admin.branding') }}
      </button>
      <button
        class="settings-tab"
        :class="{ active: tab === 'ai' }"
        @click="tab = 'ai'"
      >
        <i class="fas fa-robot me-2"></i>{{ t('admin.aiSettings') }}
      </button>
      <button
        class="settings-tab"
        :class="{ active: tab === 'ads' }"
        @click="tab = 'ads'"
      >
        <i class="fas fa-ad me-2"></i>{{ t('admin.adSettings') }}
      </button>
      <button
        class="settings-tab"
        :class="{ active: tab === 'notif' }"
        @click="tab = 'notif'"
      >
        <i class="fas fa-bell me-2"></i>{{ t('admin.sendNotification') }}
      </button>
    </div>

    <!-- Branding -->
    <div v-show="tab === 'branding'">
      <!-- Live Preview -->
      <div class="settings-card mb-4">
        <div class="settings-card-header settings-card-header--preview">
          <i class="fas fa-eye me-2"></i>{{ t('admin.livePreview') }}
        </div>
        <div class="settings-card-body p-0">
          <div
            class="preview-bar"
            :style="{ backgroundColor: branding.primary_color }"
          >
            <img
              v-if="branding.logo_url"
              :src="branding.logo_url"
              alt="Logo"
              class="preview-logo"
            />
            <div
              v-else
              class="preview-logo-placeholder"
            >
              {{ branding.name ? branding.name.charAt(0).toUpperCase() : 'C' }}
            </div>
            <span class="preview-name">{{ branding.name || t('app.name') }}</span>
          </div>
        </div>
      </div>

      <!-- Branding Settings -->
      <div class="settings-card">
        <div class="settings-card-header">
          <i class="fas fa-paint-brush me-2"></i>{{ t('admin.brandingSettings') }}
        </div>
        <div class="settings-card-body">
          <!-- System Name -->
          <div class="settings-field">
            <label class="settings-label">{{ t('admin.systemName') }}</label>
            <input type="text" class="form-control admin-input" v-model="branding.name" />
          </div>

          <!-- Logo Upload -->
          <div class="settings-field">
            <label class="settings-label">{{ t('admin.systemLogo') }}</label>

            <div v-if="branding.logo_url" class="logo-preview-row">
              <div class="logo-preview-box">
                <img
                  :src="branding.logo_url"
                  alt="Current Logo"
                  class="logo-preview-img"
                />
              </div>
              <button type="button" class="btn btn-sm admin-btn-outline-danger" @click="removeLogo">
                <i class="fas fa-trash me-1"></i>{{ t('admin.removeLogo') }}
              </button>
            </div>

            <div
              class="upload-zone"
              :class="{ 'drag-over': isDragging }"
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop.prevent="handleDrop"
              @click="logoInput?.click()"
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
              <div v-else class="upload-zone-content">
                <div class="upload-zone-icon">
                  <i class="fas fa-cloud-upload-alt"></i>
                </div>
                <p class="mb-1 text-muted">{{ t('admin.dragDropLogo') }}</p>
                <small class="text-muted">{{ t('admin.logoHint') }}</small>
              </div>
            </div>
          </div>

          <!-- Color Settings -->
          <div class="settings-field">
            <label class="settings-label">{{ t('admin.colorSettings') }}</label>
            <div class="row g-3">
              <div class="col-md-4">
                <div class="color-field">
                  <label class="color-field-label">{{ t('admin.primaryColor') }}</label>
                  <div class="color-picker-row">
                    <input
                      type="color"
                      class="color-input"
                      v-model="branding.primary_color"
                    />
                    <div
                      class="color-preview"
                      :style="{ backgroundColor: branding.primary_color }"
                    >
                      {{ branding.primary_color }}
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-md-4">
                <div class="color-field">
                  <label class="color-field-label">{{ t('admin.secondaryColor') }}</label>
                  <div class="color-picker-row">
                    <input
                      type="color"
                      class="color-input"
                      v-model="branding.secondary_color"
                    />
                    <div
                      class="color-preview"
                      :style="{ backgroundColor: branding.secondary_color }"
                    >
                      {{ branding.secondary_color }}
                    </div>
                  </div>
                </div>
              </div>
              <div class="col-md-4">
                <div class="color-field">
                  <label class="color-field-label">{{ t('admin.accentColor') }}</label>
                  <div class="color-picker-row">
                    <input
                      type="color"
                      class="color-input"
                      v-model="branding.accent_color"
                    />
                    <div
                      class="color-preview"
                      :style="{ backgroundColor: branding.accent_color }"
                    >
                      {{ branding.accent_color }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <button
            @click="saveBranding"
            class="btn admin-btn admin-btn--primary mt-2"
            :disabled="savingBranding"
          >
            <span v-if="savingBranding">
              <span class="spinner-border spinner-border-sm me-1" role="status"></span>
              {{ t('app.loading') }}
            </span>
            <span v-else><i class="fas fa-save me-1"></i>{{ t('app.save') }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- AI Settings -->
    <div v-show="tab === 'ai'">
      <div class="settings-card">
        <div class="settings-card-header">
          <i class="fas fa-robot me-2"></i>{{ t('admin.aiSettings') }}
        </div>
        <div class="settings-card-body">
          <div class="row g-4">
            <div class="col-md-4">
              <label class="settings-label">{{ t('admin.provider') }}</label>
              <select class="form-select admin-input" v-model="aiForm.provider">
                <option value="openai">OpenAI</option>
                <option value="gemini">Google Gemini</option>
                <option value="deepseek">DeepSeek</option>
              </select>
            </div>
            <div class="col-md-4">
              <label class="settings-label">{{ t('admin.model') }}</label>
              <select class="form-select admin-input" v-model="aiForm.model">
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
              <label class="settings-label">{{ t('admin.maxTokens') }}</label>
              <input type="number" class="form-control admin-input" v-model.number="aiForm.max_tokens" />
            </div>
            <div class="col-md-8">
              <label class="settings-label">{{ t('admin.apiKey') }}</label>
              <input type="password" class="form-control admin-input" v-model="aiForm.api_key" :placeholder="t('admin.apiKey')" />
            </div>
            <div class="col-md-4 d-flex align-items-end">
              <div class="form-check form-switch">
                <input type="checkbox" class="form-check-input custom-switch" v-model="aiForm.is_active" />
                <label class="form-check-label fw-semibold">{{ t('app.active') }}</label>
              </div>
            </div>
          </div>
          <button @click="saveAI" class="btn admin-btn admin-btn--primary mt-4">
            <i class="fas fa-save me-1"></i>{{ t('app.save') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Ad Settings -->
    <div v-show="tab === 'ads'">
      <div class="settings-card">
        <div class="settings-card-header">
          <i class="fas fa-ad me-2"></i>{{ t('admin.adSettings') }}
        </div>
        <div class="settings-card-body">
          <div class="row g-4">
            <div class="col-md-4">
              <label class="settings-label">{{ t('admin.adType') }}</label>
              <select class="form-select admin-input" v-model="adSettings.type">
                <option value="video">Video</option>
                <option value="adsense">Google AdSense</option>
              </select>
            </div>
            <div class="col-md-8" v-if="adSettings.type === 'video'">
              <label class="settings-label">{{ t('admin.videoUrl') }}</label>
              <input type="url" class="form-control admin-input" v-model="adSettings.video_url" />
            </div>
            <div class="col-md-8" v-else>
              <label class="settings-label">{{ t('admin.adsenseCode') }}</label>
              <textarea class="form-control admin-textarea" rows="3" v-model="adSettings.adsense_code"></textarea>
            </div>
            <div class="col-md-4">
              <label class="settings-label">{{ t('admin.duration') }}</label>
              <input type="number" class="form-control admin-input" v-model.number="adSettings.duration" />
            </div>
            <div class="col-md-4">
              <label class="settings-label">{{ t('admin.skipAfter') }}</label>
              <input type="number" class="form-control admin-input" v-model.number="adSettings.skip_after" />
            </div>
            <div class="col-md-4 d-flex align-items-end">
              <div class="form-check form-switch">
                <input type="checkbox" class="form-check-input custom-switch" v-model="adSettings.is_active" />
                <label class="form-check-label fw-semibold">{{ t('app.active') }}</label>
              </div>
            </div>
          </div>
          <button @click="saveAds" class="btn admin-btn admin-btn--primary mt-4">
            <i class="fas fa-save me-1"></i>{{ t('app.save') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Send Notification -->
    <div v-show="tab === 'notif'">
      <div class="settings-card">
        <div class="settings-card-header">
          <i class="fas fa-bell me-2"></i>{{ t('admin.sendNotification') }}
        </div>
        <div class="settings-card-body">
          <div class="row g-4">
            <div class="col-md-6">
              <label class="settings-label">{{ t('notifications.title') }} (AR)</label>
              <input type="text" class="form-control admin-input" v-model="notifForm.title_ar" dir="rtl" />
            </div>
            <div class="col-md-6">
              <label class="settings-label">{{ t('notifications.title') }} (EN)</label>
              <input type="text" class="form-control admin-input" v-model="notifForm.title_en" />
            </div>
            <div class="col-md-6">
              <label class="settings-label">{{ t('admin.message') }} (AR)</label>
              <textarea class="form-control admin-textarea" rows="3" v-model="notifForm.message_ar" dir="rtl"></textarea>
            </div>
            <div class="col-md-6">
              <label class="settings-label">{{ t('admin.message') }} (EN)</label>
              <textarea class="form-control admin-textarea" rows="3" v-model="notifForm.message_en"></textarea>
            </div>
            <div class="col-md-4">
              <label class="settings-label">{{ t('admin.target') }}</label>
              <select class="form-select admin-input" v-model="notifForm.target" @change="onTargetChange">
                <option value="all">{{ t('app.allUsers') }}</option>
                <option value="faculty">{{ t('admin.byFaculty') }}</option>
                <option value="department">{{ t('admin.byDepartment') }}</option>
              </select>
            </div>
            <div class="col-md-4" v-if="notifForm.target === 'faculty' || notifForm.target === 'department'">
              <label class="settings-label">{{ t('cv.selectFaculty') }}</label>
              <select class="form-select admin-input" v-model="notifForm.faculty_id" @change="onNotifFacultyChange">
                <option :value="null" disabled>-- {{ t('cv.selectFaculty') }} --</option>
                <option v-for="f in faculties" :key="f.id" :value="f.id">{{ f.name_ar }} - {{ f.name_en }}</option>
              </select>
            </div>
            <div class="col-md-4" v-if="notifForm.target === 'department' && notifForm.faculty_id">
              <label class="settings-label">{{ t('cv.selectDepartment') }}</label>
              <select class="form-select admin-input" v-model="notifForm.department_id">
                <option :value="null" disabled>-- {{ t('cv.selectDepartment') }} --</option>
                <option v-for="d in notifDepartments" :key="d.id" :value="d.id">{{ d.name_ar }} - {{ d.name_en }}</option>
              </select>
            </div>
          </div>
          <button @click="sendNotif" class="btn admin-btn admin-btn--primary mt-4" :disabled="sendingNotif">
            <span v-if="sendingNotif" class="spinner-border spinner-border-sm me-1"></span>
            <i v-else class="fas fa-paper-plane me-1"></i>
            {{ t('admin.sendNotification') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="message" class="alert admin-alert-success mt-3">
      <i class="fas fa-check-circle me-2"></i>{{ message }}
    </div>
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

const logoInput = ref<HTMLInputElement>()
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
  } catch {
    toast.error(t('app.error'))
  }
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
    } catch {
      toast.error(t('app.error'))
    }
  }
}

async function doUpload(file: File) {
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
  }
}

function uploadLogo(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  doUpload(file)
  input.value = ''
}

function handleDrop(event: DragEvent) {
  isDragging.value = false
  const file = event.dataTransfer?.files?.[0]
  if (!file) return
  if (!['image/jpeg', 'image/png', 'image/gif', 'image/svg+xml', 'image/webp'].includes(file.type)) {
    toast.error('Invalid file type')
    return
  }
  doUpload(file)
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
.admin-settings {
  max-width: 1000px;
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

/* ── Tab Navigation ── */
.settings-tabs {
  display: flex;
  gap: 6px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.settings-tab {
  padding: 10px 20px;
  border-radius: 10px;
  border: 1px solid #dee2e6;
  background: #ffffff;
  color: #6c757d;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.settings-tab:hover {
  color: #1a5276;
  border-color: #1a5276;
}

.settings-tab.active {
  background: #1a5276;
  color: #ffffff;
  border-color: #1a5276;
}

/* ── Settings Card ── */
.settings-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  margin-bottom: 20px;
}

.settings-card-header {
  padding: 18px 24px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 15px;
  font-weight: 600;
  color: #1a5276;
}

.settings-card-header--preview {
  background: #f8f9fb;
  color: #6c757d;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.settings-card-body {
  padding: 24px;
}

/* ── Settings Field ── */
.settings-field {
  margin-bottom: 28px;
}

.settings-label {
  display: block;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 8px;
  font-size: 14px;
}

/* ── Preview Bar ── */
.preview-bar {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  gap: 14px;
  color: #ffffff;
}

.preview-logo {
  height: 36px;
  width: auto;
  object-fit: contain;
  border-radius: 6px;
}

.preview-logo-placeholder {
  height: 36px;
  width: 36px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
}

.preview-name {
  font-weight: 600;
  font-size: 18px;
}

/* ── Logo Preview ── */
.logo-preview-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.logo-preview-box {
  padding: 10px;
  border: 1px solid #e9ecef;
  border-radius: 10px;
  background: #f8f9fa;
}

.logo-preview-img {
  max-height: 64px;
  max-width: 200px;
  object-fit: contain;
}

/* ── Upload Zone ── */
.upload-zone {
  border: 2px dashed #dee2e6;
  border-radius: 12px;
  padding: 32px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: #fafbfc;
}

.upload-zone:hover,
.upload-zone.drag-over {
  border-color: #1a5276;
  background: rgba(26, 82, 118, 0.03);
}

.upload-zone-icon {
  font-size: 36px;
  color: #adb5bd;
  margin-bottom: 8px;
}

/* ── Color Picker ── */
.color-field {
  padding: 14px;
  background: #f8f9fb;
  border-radius: 10px;
}

.color-field-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: #6c757d;
  margin-bottom: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.color-picker-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.color-input {
  width: 48px;
  height: 40px;
  border: 2px solid #dee2e6;
  border-radius: 8px;
  cursor: pointer;
  padding: 2px;
}

.color-preview {
  flex: 1;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  font-size: 13px;
  font-family: monospace;
  font-weight: 500;
}

/* ── Form Elements ── */
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

.admin-textarea {
  border-radius: 10px;
  border: 1px solid #dee2e6;
  font-size: 14px;
  resize: vertical;
}

.admin-textarea:focus {
  border-color: #1a5276;
  box-shadow: 0 0 0 3px rgba(26, 82, 118, 0.1);
}

.custom-switch:checked {
  background-color: #1a5276;
  border-color: #1a5276;
}

/* ── Buttons ── */
.admin-btn {
  border-radius: 8px;
  font-weight: 600;
  padding: 10px 24px;
  font-size: 14px;
  border: none;
}

.admin-btn--primary {
  background: #1a5276;
  color: #ffffff;
}
.admin-btn--primary:hover {
  background: #154360;
  color: #ffffff;
}

.admin-btn-outline-danger {
  border: 1px solid #c62828;
  color: #c62828;
  background: transparent;
  border-radius: 8px;
  font-weight: 500;
}
.admin-btn-outline-danger:hover {
  background: #c62828;
  color: #ffffff;
}

/* ── Alert ── */
.admin-alert-success {
  background: #e8f5e9;
  color: #2e7d32;
  border: 1px solid #c8e6c9;
  border-radius: 10px;
  padding: 14px 20px;
  font-weight: 500;
}
</style>
