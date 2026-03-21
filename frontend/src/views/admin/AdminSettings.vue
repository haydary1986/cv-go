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
    <div v-show="tab === 'branding'" class="card">
      <div class="card-body">
        <div class="row g-3">
          <div class="col-md-6">
            <label class="form-label">{{ t('admin.brandName') }}</label>
            <input type="text" class="form-control" v-model="branding.name" />
          </div>
          <div class="col-md-6">
            <label class="form-label">{{ t('admin.logoUrl') }}</label>
            <input type="url" class="form-control" v-model="branding.logo_url" />
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.primaryColor') }}</label>
            <input type="color" class="form-control form-control-color w-100" v-model="branding.primary_color" />
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.secondaryColor') }}</label>
            <input type="color" class="form-control form-control-color w-100" v-model="branding.secondary_color" />
          </div>
          <div class="col-md-4">
            <label class="form-label">{{ t('admin.accentColor') }}</label>
            <input type="color" class="form-control form-control-color w-100" v-model="branding.accent_color" />
          </div>
        </div>
        <button @click="saveBranding" class="btn btn-primary mt-3">{{ t('app.save') }}</button>
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
            <select class="form-select" v-model="notifForm.target">
              <option value="all">{{ t('app.allUsers') }}</option>
              <option value="faculty">{{ t('admin.byFaculty') }}</option>
              <option value="department">{{ t('admin.byDepartment') }}</option>
            </select>
          </div>
        </div>
        <button @click="sendNotif" class="btn btn-primary mt-3">{{ t('admin.sendNotification') }}</button>
      </div>
    </div>

    <div v-if="message" class="alert alert-success mt-3">{{ message }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'

const { t } = useI18n()
const tab = ref('branding')
const message = ref('')

const branding = reactive({ name: '', logo_url: '', primary_color: '#0d6efd', secondary_color: '#6c757d', accent_color: '#198754' })
const aiForm = reactive({ provider: 'openai', model: 'gpt-4o-mini', api_key: '', max_tokens: 2000, is_active: false })
const adSettings = reactive({ type: 'video', video_url: '', adsense_code: '', duration: 5, skip_after: 3, is_active: false })
const notifForm = reactive({ title_ar: '', title_en: '', message_ar: '', message_en: '', target: 'all' })

onMounted(async () => {
  try {
    const [bRes, aiRes, adRes] = await Promise.all([
      adminAPI.getBranding(),
      adminAPI.getAISettings(),
      adminAPI.getAdSettings(),
    ])
    if (bRes.data.branding) Object.assign(branding, bRes.data.branding)
    if (aiRes.data.settings) Object.assign(aiForm, aiRes.data.settings)
    if (adRes.data.settings) Object.assign(adSettings, adRes.data.settings)
  } catch {}
})

async function saveBranding() { await adminAPI.updateBranding(branding); message.value = t('app.success') }
async function saveAI() { await adminAPI.updateAISettings(aiForm); message.value = t('app.success') }
async function saveAds() { await adminAPI.updateAdSettings(adSettings); message.value = t('app.success') }
async function sendNotif() { await adminAPI.sendNotification(notifForm); message.value = t('app.success') }
</script>
