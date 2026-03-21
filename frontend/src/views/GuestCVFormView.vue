<template>
  <div class="container py-4">
    <h3 class="mb-3">{{ t('cv.createWithoutAccount') }}</h3>

    <!-- Guest Warning Banner -->
    <div class="alert alert-warning d-flex align-items-start mb-4">
      <i class="fas fa-exclamation-triangle me-2 mt-1"></i>
      <div>
        <strong>{{ t('cv.guestMode') }}</strong><br />
        {{ t('cv.guestWarning') }}
        <div class="mt-2">
          <router-link to="/register" class="btn btn-sm btn-primary me-2">
            <i class="fas fa-user-plus me-1"></i>{{ t('auth.registerBtn') }}
          </router-link>
          <router-link to="/login" class="btn btn-sm btn-outline-primary">
            <i class="fas fa-sign-in-alt me-1"></i>{{ t('auth.loginBtn') }}
          </router-link>
        </div>
      </div>
    </div>

    <!-- Steps Navigation -->
    <ul class="nav nav-pills nav-fill mb-4 flex-nowrap overflow-auto">
      <li class="nav-item" v-for="(step, idx) in steps" :key="idx">
        <a class="nav-link" :class="{ active: currentStep === idx }" href="#" @click.prevent="currentStep = idx">
          {{ step }}
        </a>
      </li>
    </ul>

    <div class="card shadow-sm">
      <div class="card-body p-4">

        <!-- Step 0: Guest Info + Basic Info -->
        <div v-show="currentStep === 0">
          <div class="row g-3">
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.guestName') }} *</label>
              <input type="text" class="form-control" v-model="guestName" required />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.guestEmail') }}</label>
              <input type="email" class="form-control" v-model="guestEmail" />
            </div>
            <hr />
            <div class="col-md-12">
              <label class="form-label">{{ t('cv.title') }} *</label>
              <input type="text" class="form-control" v-model="form.title" required />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.language') }}</label>
              <select class="form-select" v-model="form.language">
                <option value="ar">{{ t('app.arabic') }}</option>
                <option value="en">{{ t('app.english') }}</option>
              </select>
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.template') }}</label>
              <select class="form-select" v-model="form.template">
                <option v-for="tmpl in templates" :key="tmpl" :value="tmpl">{{ t(`templates.${tmpl}`) }}</option>
              </select>
            </div>
            <div class="col-12">
              <label class="form-label">{{ t('cv.chooseTemplate') }}</label>
              <div class="row g-2">
                <div class="col-4 col-md-2" v-for="tmpl in templates" :key="tmpl">
                  <div class="card text-center p-2" :class="{ 'border-primary bg-primary bg-opacity-10': form.template === tmpl }" @click="form.template = tmpl" style="cursor:pointer">
                    <i class="fas fa-file-alt fa-2x mb-1" :class="form.template === tmpl ? 'text-primary' : 'text-muted'"></i>
                    <small>{{ t(`templates.${tmpl}`) }}</small>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 1: Personal Info -->
        <div v-show="currentStep === 1">
          <div class="row g-3">
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.fullName') }} *</label>
              <input type="text" class="form-control" v-model="form.data.personal_info.full_name" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.jobTitle') }}</label>
              <input type="text" class="form-control" v-model="form.data.personal_info.job_title" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.email') }}</label>
              <input type="email" class="form-control" v-model="form.data.personal_info.email" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.phone') }}</label>
              <input type="tel" class="form-control" v-model="form.data.personal_info.phone" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.address') }}</label>
              <input type="text" class="form-control" v-model="form.data.personal_info.address" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.dateOfBirth') }}</label>
              <input type="date" class="form-control" v-model="form.data.personal_info.date_of_birth" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.nationality') }}</label>
              <input type="text" class="form-control" v-model="form.data.personal_info.nationality" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.linkedin') }}</label>
              <input type="url" class="form-control" v-model="form.data.personal_info.linkedin" />
            </div>
          </div>
        </div>

        <!-- Step 2: Objective -->
        <div v-show="currentStep === 2">
          <label class="form-label">{{ t('cv.objective') }}</label>
          <textarea class="form-control" rows="5" v-model="form.data.objective"></textarea>
        </div>

        <!-- Step 3: Experience -->
        <div v-show="currentStep === 3">
          <div v-for="(exp, idx) in form.data.experiences" :key="idx" class="border rounded p-3 mb-3">
            <div class="d-flex justify-content-between mb-2">
              <h6>{{ t('cv.experience') }} {{ idx + 1 }}</h6>
              <button @click="form.data.experiences.splice(idx, 1)" class="btn btn-sm btn-outline-danger">
                <i class="fas fa-trash"></i>
              </button>
            </div>
            <div class="row g-2">
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.position')" v-model="exp.title" />
              </div>
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.companyName')" v-model="exp.company" />
              </div>
              <div class="col-md-4">
                <input type="text" class="form-control" :placeholder="t('cv.location')" v-model="exp.location" />
              </div>
              <div class="col-md-3">
                <input type="month" class="form-control" v-model="exp.start_date" />
              </div>
              <div class="col-md-3">
                <input type="month" class="form-control" v-model="exp.end_date" :disabled="exp.current" />
              </div>
              <div class="col-md-2 d-flex align-items-center">
                <div class="form-check">
                  <input type="checkbox" class="form-check-input" v-model="exp.current" />
                  <label class="form-check-label">{{ t('cv.current') }}</label>
                </div>
              </div>
              <div class="col-12">
                <textarea class="form-control" rows="2" :placeholder="t('cv.description')" v-model="exp.description"></textarea>
              </div>
            </div>
          </div>
          <button @click="addExperience" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addExperience') }}
          </button>
        </div>

        <!-- Step 4: Education -->
        <div v-show="currentStep === 4">
          <div v-for="(edu, idx) in form.data.education" :key="idx" class="border rounded p-3 mb-3">
            <div class="d-flex justify-content-between mb-2">
              <h6>{{ t('cv.education') }} {{ idx + 1 }}</h6>
              <button @click="form.data.education.splice(idx, 1)" class="btn btn-sm btn-outline-danger">
                <i class="fas fa-trash"></i>
              </button>
            </div>
            <div class="row g-2">
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.degree')" v-model="edu.degree" />
              </div>
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.institution')" v-model="edu.institution" />
              </div>
              <div class="col-md-4">
                <input type="text" class="form-control" :placeholder="t('cv.location')" v-model="edu.location" />
              </div>
              <div class="col-md-3">
                <input type="month" class="form-control" v-model="edu.start_date" />
              </div>
              <div class="col-md-3">
                <input type="month" class="form-control" v-model="edu.end_date" />
              </div>
              <div class="col-md-2">
                <input type="text" class="form-control" :placeholder="t('cv.gpa')" v-model="edu.gpa" />
              </div>
            </div>
          </div>
          <button @click="addEducation" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addEducation') }}
          </button>
        </div>

        <!-- Step 5: Skills -->
        <div v-show="currentStep === 5">
          <div v-for="(skill, idx) in form.data.skills" :key="idx" class="row g-2 mb-2">
            <div class="col-md-6">
              <input type="text" class="form-control" :placeholder="t('cv.skillName')" v-model="skill.name" />
            </div>
            <div class="col-md-4">
              <select class="form-select" v-model="skill.level">
                <option value="beginner">{{ t('cv.beginner') }}</option>
                <option value="intermediate">{{ t('cv.intermediate') }}</option>
                <option value="advanced">{{ t('cv.advanced') }}</option>
                <option value="expert">{{ t('cv.expert') }}</option>
              </select>
            </div>
            <div class="col-md-2">
              <button @click="form.data.skills.splice(idx, 1)" class="btn btn-outline-danger w-100">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
          <button @click="form.data.skills.push({ name: '', level: 'intermediate' })" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addSkill') }}
          </button>
        </div>

        <!-- Step 6: Languages -->
        <div v-show="currentStep === 6">
          <div v-for="(lang, idx) in form.data.languages" :key="idx" class="row g-2 mb-2">
            <div class="col-md-6">
              <input type="text" class="form-control" :placeholder="t('cv.languageName')" v-model="lang.name" />
            </div>
            <div class="col-md-4">
              <select class="form-select" v-model="lang.level">
                <option value="basic">{{ t('cv.basic') }}</option>
                <option value="conversational">{{ t('cv.conversational') }}</option>
                <option value="fluent">{{ t('cv.fluent') }}</option>
                <option value="native">{{ t('cv.native') }}</option>
              </select>
            </div>
            <div class="col-md-2">
              <button @click="form.data.languages.splice(idx, 1)" class="btn btn-outline-danger w-100">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
          <button @click="form.data.languages.push({ name: '', level: 'conversational' })" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addLanguage') }}
          </button>
        </div>

        <!-- Step 7: Certificates -->
        <div v-show="currentStep === 7">
          <div v-for="(cert, idx) in form.data.certificates" :key="idx" class="border rounded p-3 mb-3">
            <div class="d-flex justify-content-between mb-2">
              <h6>{{ t('cv.certificates') }} {{ idx + 1 }}</h6>
              <button @click="form.data.certificates.splice(idx, 1)" class="btn btn-sm btn-outline-danger">
                <i class="fas fa-trash"></i>
              </button>
            </div>
            <div class="row g-2">
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.certName')" v-model="cert.name" />
              </div>
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.certIssuer')" v-model="cert.issuer" />
              </div>
              <div class="col-md-4">
                <input type="month" class="form-control" v-model="cert.date" />
              </div>
            </div>
          </div>
          <button @click="form.data.certificates.push({ name: '', issuer: '', date: '', expiry_date: '', credential_id: '', url: '' })" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addCertificate') }}
          </button>
        </div>

        <!-- Navigation -->
        <div class="d-flex justify-content-between mt-4 pt-3 border-top">
          <button @click="currentStep--" class="btn btn-outline-secondary" :disabled="currentStep === 0">
            <i class="fas fa-arrow-left me-1"></i>{{ t('app.previous') }}
          </button>
          <div>
            <button v-if="currentStep < steps.length - 1" @click="currentStep++" class="btn btn-primary">
              {{ t('app.next') }} <i class="fas fa-arrow-right ms-1"></i>
            </button>
            <button v-else @click="handleSubmit" class="btn btn-success" :disabled="saving">
              <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
              <i v-else class="fas fa-save me-1"></i>{{ t('app.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Modal -->
    <div v-if="showSuccess" class="modal fade show d-block" tabindex="-1" style="background:rgba(0,0,0,0.5)">
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-body text-center py-5">
            <i class="fas fa-check-circle fa-4x text-success mb-3"></i>
            <h4>{{ t('cv.guestSuccess') }}</h4>
            <div class="input-group my-3">
              <input type="text" class="form-control" :value="shareLink" readonly />
              <button @click="copyLink" class="btn btn-outline-primary">
                <i class="fas fa-copy"></i>
              </button>
            </div>
            <div class="d-flex gap-2 justify-content-center">
              <router-link :to="shareLink ? `/shared/${shareToken}` : '/'" class="btn btn-primary">
                <i class="fas fa-eye me-1"></i>{{ t('app.view') }}
              </router-link>
              <router-link to="/register" class="btn btn-outline-success">
                <i class="fas fa-user-plus me-1"></i>{{ t('auth.registerBtn') }}
              </router-link>
              <router-link to="/" class="btn btn-outline-secondary">
                <i class="fas fa-home me-1"></i>{{ t('app.home') }}
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { getEmptyCVData } from '../stores/cv'
import { cvAPI } from '../services/api'
import { useToast } from '../composables/useToast'

const { t } = useI18n()
const toast = useToast()

const currentStep = ref(0)
const saving = ref(false)
const showSuccess = ref(false)
const shareLink = ref('')
const shareToken = ref('')
const guestName = ref('')
const guestEmail = ref('')

const templates = [
  'academic', 'ats', 'compact', 'creative', 'designer', 'elegant', 'engineer',
  'executive', 'lawyer', 'medical', 'minimalist', 'modern', 'professional',
  'standard', 'teacher', 'tech'
]

const steps = [
  t('cv.template'), t('cv.personalInfo'), t('cv.objective'),
  t('cv.experience'), t('cv.education'), t('cv.skills'), t('cv.languages'),
  t('cv.certificates')
]

const form = reactive({
  title: '',
  language: 'ar',
  template: 'modern',
  data: getEmptyCVData(),
})

function addExperience() {
  form.data.experiences.push({
    title: '', company: '', location: '', start_date: '', end_date: '', current: false, description: ''
  })
}

function addEducation() {
  form.data.education.push({
    degree: '', institution: '', location: '', start_date: '', end_date: '', gpa: '', description: ''
  })
}

function copyLink() {
  navigator.clipboard.writeText(shareLink.value)
}

async function handleSubmit() {
  if (!form.title || !guestName.value) {
    currentStep.value = 0
    return
  }
  saving.value = true
  try {
    const res = await cvAPI.createGuest({
      ...form,
      guest_name: guestName.value,
      guest_email: guestEmail.value,
    })
    const cv = res.data.cv
    shareToken.value = cv.share_token
    shareLink.value = `${window.location.origin}/shared/${cv.share_token}`
    showSuccess.value = true
  } catch {
    toast.error(t('app.error'))
  } finally {
    saving.value = false
  }
}
</script>
