<template>
  <div class="guest-form-page">
    <div class="container py-5" style="max-width: 800px;">

      <!-- Header -->
      <div class="text-center mb-4">
        <h3 class="page-title">{{ t('cv.createWithoutAccount') }}</h3>
        <p class="page-step mb-0">{{ t('cv.step') || '' }} {{ currentStep + 1 }} / {{ steps.length }}</p>
      </div>

      <!-- Guest Warning Banner -->
      <div class="guest-warning-card mb-4">
        <div class="guest-warning-icon">
          <i class="fas fa-user-clock"></i>
        </div>
        <div class="guest-warning-content">
          <h6 class="guest-warning-title">{{ t('cv.guestMode') }}</h6>
          <p class="guest-warning-text mb-2">{{ t('cv.guestWarning') }}</p>
          <div class="d-flex gap-2 flex-wrap">
            <router-link to="/register" class="guest-cta-btn guest-cta-btn--primary">
              <i class="fas fa-user-plus me-1"></i>{{ t('auth.registerBtn') }}
            </router-link>
            <router-link to="/login" class="guest-cta-btn guest-cta-btn--outline">
              <i class="fas fa-sign-in-alt me-1"></i>{{ t('auth.loginBtn') }}
            </router-link>
          </div>
        </div>
      </div>

      <!-- Progress Indicator -->
      <div class="progress-track mb-4">
        <div class="progress-track-fill" :style="{ width: ((currentStep + 1) / steps.length * 100) + '%' }"></div>
      </div>

      <!-- Steps Navigation -->
      <div class="steps-nav mb-4">
        <button
          v-for="(step, idx) in steps"
          :key="idx"
          class="step-chip"
          :class="{
            'step-chip--active': currentStep === idx,
            'step-chip--done': currentStep > idx
          }"
          @click="currentStep = idx"
        >
          <span class="step-chip-num">{{ idx + 1 }}</span>
          <span class="step-chip-label d-none d-md-inline">{{ step }}</span>
        </button>
      </div>

      <!-- Step Card -->
      <div class="step-card">
        <div class="step-card-body">

          <!-- Step 0: Guest Info + Basic Info -->
          <div v-show="currentStep === 0">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-palette"></i></div>
              <h4>{{ t('cv.template') }}</h4>
            </div>
            <div class="row g-3">
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.guestName') }} *</label>
                <input type="text" class="form-control form-control-lg field-input" v-model="guestName" required />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.guestEmail') }}</label>
                <input type="email" class="form-control form-control-lg field-input" v-model="guestEmail" />
              </div>
              <div class="col-12"><hr class="section-divider" /></div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.language') }}</label>
                <select class="form-select form-select-lg field-input" v-model="form.language">
                  <option value="ar">{{ t('app.arabic') }}</option>
                  <option value="en">{{ t('app.english') }}</option>
                </select>
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.template') }}</label>
                <select class="form-select form-select-lg field-input" v-model="form.template">
                  <option v-for="tmpl in templates" :key="tmpl" :value="tmpl">{{ t(`templates.${tmpl}`) }}</option>
                </select>
              </div>
              <div class="col-12">
                <label class="form-label field-label">{{ t('cv.chooseTemplate') }}</label>
                <div class="row g-2">
                  <div class="col-4 col-md-2" v-for="tmpl in templates" :key="tmpl">
                    <div
                      class="template-card"
                      :class="{ 'template-card--selected': form.template === tmpl }"
                      @click="form.template = tmpl"
                    >
                      <i class="fas fa-file-alt fa-2x mb-1" :class="form.template === tmpl ? 'text-accent' : 'text-secondary-light'"></i>
                      <small>{{ t(`templates.${tmpl}`) }}</small>
                    </div>
                  </div>
                </div>
              </div>

              <!-- University Affiliation Section -->
              <div class="col-12">
                <hr class="section-divider" />
                <h5 class="section-heading mb-3">
                  <i class="fas fa-university me-2 section-heading-icon"></i>
                  {{ t('cv.universityAffiliation') }}
                </h5>
                <p class="section-desc mb-3">{{ t('cv.areYouMember') }}</p>

                <div class="row g-3">
                  <div class="col-md-6">
                    <div
                      class="affiliation-card"
                      :class="{ 'affiliation-card--selected': isUniversityMember }"
                      @click="isUniversityMember = true"
                    >
                      <i class="fas fa-university fa-2x mb-2" :class="isUniversityMember ? 'text-accent' : 'text-secondary-light'"></i>
                      <h6>{{ t('cv.yesMember') }}</h6>
                      <small class="text-secondary-light">{{ t('cv.memberDescription') }}</small>
                    </div>
                  </div>
                  <div class="col-md-6">
                    <div
                      class="affiliation-card"
                      :class="{ 'affiliation-card--selected': !isUniversityMember }"
                      @click="isUniversityMember = false; selectedFacultyId = null; selectedDepartmentId = null"
                    >
                      <i class="fas fa-user-tie fa-2x mb-2" :class="!isUniversityMember ? 'text-accent' : 'text-secondary-light'"></i>
                      <h6>{{ t('cv.noExternal') }}</h6>
                      <small class="text-secondary-light">{{ t('cv.externalDescription') }}</small>
                    </div>
                  </div>
                </div>

                <!-- Faculty & Department Dropdowns -->
                <div v-if="isUniversityMember" class="row g-3 mt-2">
                  <div class="col-md-6">
                    <label class="form-label field-label">{{ t('cv.selectFaculty') }} *</label>
                    <select class="form-select field-input" v-model="selectedFacultyId">
                      <option :value="null" disabled>{{ t('cv.selectFaculty') }}</option>
                      <option v-for="fac in faculties" :key="fac.id" :value="fac.id">
                        {{ fac.name_ar || fac.name_en || fac.name }}
                      </option>
                    </select>
                  </div>
                  <div class="col-md-6">
                    <label class="form-label field-label">{{ t('cv.selectDepartment') }}</label>
                    <select class="form-select field-input" v-model="selectedDepartmentId" :disabled="!selectedFacultyId">
                      <option :value="null" disabled>{{ t('cv.selectDepartment') }}</option>
                      <option v-for="dept in departments" :key="dept.id" :value="dept.id">
                        {{ dept.name_ar || dept.name_en || dept.name }}
                      </option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Step 1: Personal Info -->
          <div v-show="currentStep === 1">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-user"></i></div>
              <h4>{{ t('cv.personalInfo') }}</h4>
            </div>
            <div class="row g-3">
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.fullName') }} *</label>
                <input type="text" class="form-control field-input" v-model="form.data.personal_info.full_name" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.jobTitle') }}</label>
                <input type="text" class="form-control field-input" v-model="form.data.personal_info.job_title" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.email') }}</label>
                <input type="email" class="form-control field-input" v-model="form.data.personal_info.email" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.phone') }}</label>
                <input type="tel" class="form-control field-input" v-model="form.data.personal_info.phone" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.address') }}</label>
                <input type="text" class="form-control field-input" v-model="form.data.personal_info.address" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.dateOfBirth') }}</label>
                <input type="date" class="form-control field-input" v-model="form.data.personal_info.date_of_birth" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.nationality') }}</label>
                <input type="text" class="form-control field-input" v-model="form.data.personal_info.nationality" />
              </div>
              <div class="col-md-6">
                <label class="form-label field-label">{{ t('cv.linkedin') }}</label>
                <input type="url" class="form-control field-input" v-model="form.data.personal_info.linkedin" />
              </div>
            </div>
          </div>

          <!-- Step 2: Objective -->
          <div v-show="currentStep === 2">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-bullseye"></i></div>
              <h4>{{ t('cv.objective') }}</h4>
            </div>
            <label class="form-label field-label">{{ t('cv.objective') }}</label>
            <textarea class="form-control field-input" rows="5" v-model="form.data.objective"></textarea>
          </div>

          <!-- Step 3: Experience -->
          <div v-show="currentStep === 3">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-briefcase"></i></div>
              <h4>{{ t('cv.experience') }}</h4>
            </div>
            <div v-for="(exp, idx) in form.data.experiences" :key="idx" class="section-entry mb-3">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0 entry-label">{{ t('cv.experience') }} {{ idx + 1 }}</h6>
                <button @click="form.data.experiences.splice(idx, 1)" class="btn-delete">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </div>
              <div class="row g-2">
                <div class="col-md-6">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.position')" v-model="exp.title" />
                </div>
                <div class="col-md-6">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.companyName')" v-model="exp.company" />
                </div>
                <div class="col-md-4">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.location')" v-model="exp.location" />
                </div>
                <div class="col-md-3">
                  <input type="month" class="form-control field-input" v-model="exp.start_date" />
                </div>
                <div class="col-md-3">
                  <input type="month" class="form-control field-input" v-model="exp.end_date" :disabled="exp.current" />
                </div>
                <div class="col-md-2 d-flex align-items-center">
                  <div class="form-check">
                    <input type="checkbox" class="form-check-input" v-model="exp.current" />
                    <label class="form-check-label">{{ t('cv.current') }}</label>
                  </div>
                </div>
                <div class="col-12">
                  <textarea class="form-control field-input" rows="2" :placeholder="t('cv.description')" v-model="exp.description"></textarea>
                </div>
              </div>
            </div>
            <button @click="addExperience" class="btn-add">
              <i class="fas fa-plus me-1"></i>{{ t('cv.addExperience') }}
            </button>
          </div>

          <!-- Step 4: Education -->
          <div v-show="currentStep === 4">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-graduation-cap"></i></div>
              <h4>{{ t('cv.education') }}</h4>
            </div>
            <div v-for="(edu, idx) in form.data.education" :key="idx" class="section-entry mb-3">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0 entry-label">{{ t('cv.education') }} {{ idx + 1 }}</h6>
                <button @click="form.data.education.splice(idx, 1)" class="btn-delete">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </div>
              <div class="row g-2">
                <div class="col-md-6">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.degree')" v-model="edu.degree" />
                </div>
                <div class="col-md-6">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.institution')" v-model="edu.institution" />
                </div>
                <div class="col-md-4">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.location')" v-model="edu.location" />
                </div>
                <div class="col-md-3">
                  <input type="month" class="form-control field-input" v-model="edu.start_date" />
                </div>
                <div class="col-md-3">
                  <input type="month" class="form-control field-input" v-model="edu.end_date" />
                </div>
                <div class="col-md-2">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.gpa')" v-model="edu.gpa" />
                </div>
              </div>
            </div>
            <button @click="addEducation" class="btn-add">
              <i class="fas fa-plus me-1"></i>{{ t('cv.addEducation') }}
            </button>
          </div>

          <!-- Step 5: Skills -->
          <div v-show="currentStep === 5">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-tools"></i></div>
              <h4>{{ t('cv.skills') }}</h4>
            </div>
            <div v-for="(skill, idx) in form.data.skills" :key="idx" class="row g-2 mb-2 align-items-center">
              <div class="col-md-6">
                <input type="text" class="form-control field-input" :placeholder="t('cv.skillName')" v-model="skill.name" />
              </div>
              <div class="col-md-4">
                <select class="form-select field-input" v-model="skill.level">
                  <option value="beginner">{{ t('cv.beginner') }}</option>
                  <option value="intermediate">{{ t('cv.intermediate') }}</option>
                  <option value="advanced">{{ t('cv.advanced') }}</option>
                  <option value="expert">{{ t('cv.expert') }}</option>
                </select>
              </div>
              <div class="col-md-2">
                <button @click="form.data.skills.splice(idx, 1)" class="btn-delete w-100">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </div>
            </div>
            <button @click="form.data.skills.push({ name: '', level: 'intermediate' })" class="btn-add">
              <i class="fas fa-plus me-1"></i>{{ t('cv.addSkill') }}
            </button>
          </div>

          <!-- Step 6: Languages -->
          <div v-show="currentStep === 6">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-language"></i></div>
              <h4>{{ t('cv.languages') }}</h4>
            </div>
            <div v-for="(lang, idx) in form.data.languages" :key="idx" class="row g-2 mb-2 align-items-center">
              <div class="col-md-6">
                <input type="text" class="form-control field-input" :placeholder="t('cv.languageName')" v-model="lang.name" />
              </div>
              <div class="col-md-4">
                <select class="form-select field-input" v-model="lang.level">
                  <option value="basic">{{ t('cv.basic') }}</option>
                  <option value="conversational">{{ t('cv.conversational') }}</option>
                  <option value="fluent">{{ t('cv.fluent') }}</option>
                  <option value="native">{{ t('cv.native') }}</option>
                </select>
              </div>
              <div class="col-md-2">
                <button @click="form.data.languages.splice(idx, 1)" class="btn-delete w-100">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </div>
            </div>
            <button @click="form.data.languages.push({ name: '', level: 'conversational' })" class="btn-add">
              <i class="fas fa-plus me-1"></i>{{ t('cv.addLanguage') }}
            </button>
          </div>

          <!-- Step 7: Certificates -->
          <div v-show="currentStep === 7">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-certificate"></i></div>
              <h4>{{ t('cv.certificates') }}</h4>
            </div>
            <div v-for="(cert, idx) in form.data.certificates" :key="idx" class="section-entry mb-3">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <h6 class="mb-0 entry-label">{{ t('cv.certificates') }} {{ idx + 1 }}</h6>
                <button @click="form.data.certificates.splice(idx, 1)" class="btn-delete">
                  <i class="fas fa-trash-alt"></i>
                </button>
              </div>
              <div class="row g-2">
                <div class="col-md-6">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.certName')" v-model="cert.name" />
                </div>
                <div class="col-md-6">
                  <input type="text" class="form-control field-input" :placeholder="t('cv.certIssuer')" v-model="cert.issuer" />
                </div>
                <div class="col-md-4">
                  <input type="month" class="form-control field-input" v-model="cert.date" />
                </div>
              </div>
            </div>
            <button @click="form.data.certificates.push({ name: '', issuer: '', date: '', expiry_date: '', credential_id: '', url: '' })" class="btn-add">
              <i class="fas fa-plus me-1"></i>{{ t('cv.addCertificate') }}
            </button>

            <!-- Legal Disclaimer -->
            <div class="disclaimer-card mt-4">
              <div class="d-flex align-items-start">
                <div class="disclaimer-icon">
                  <i class="fas fa-shield-alt"></i>
                </div>
                <div>
                  <div class="form-check">
                    <input
                      type="checkbox"
                      class="form-check-input"
                      id="disclaimerCheck"
                      v-model="agreedToTerms"
                    />
                    <label class="form-check-label" for="disclaimerCheck">
                      {{ t('cv.disclaimer') }}
                    </label>
                  </div>
                  <small v-if="!agreedToTerms" class="text-danger mt-1 d-block">
                    {{ t('cv.mustAgree') }}
                  </small>
                </div>
              </div>
            </div>
          </div>

          <!-- Navigation -->
          <div class="d-flex justify-content-between mt-4 pt-3 border-top">
            <button @click="currentStep--" class="nav-btn nav-btn--back" :disabled="currentStep === 0">
              <i class="fas fa-arrow-left me-1"></i>{{ t('app.previous') }}
            </button>
            <div>
              <button v-if="currentStep < steps.length - 1" @click="currentStep++" class="nav-btn nav-btn--next">
                {{ t('app.next') }} <i class="fas fa-arrow-right ms-1"></i>
              </button>
              <button v-else @click="handleSubmit" class="nav-btn nav-btn--submit" :disabled="saving || !agreedToTerms">
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
          <div class="modal-content success-modal">
            <div class="modal-body text-center py-5 px-4">
              <div class="success-icon-wrap mb-3">
                <i class="fas fa-check"></i>
              </div>
              <h4 class="success-title mb-1">{{ t('cv.guestSuccess') }}</h4>
              <div class="input-group my-3 share-link-group">
                <input type="text" class="form-control" :value="shareLink" readonly />
                <button @click="copyLink" class="btn btn-outline-secondary">
                  <i class="fas fa-copy"></i>
                </button>
              </div>
              <div class="d-flex gap-2 justify-content-center flex-wrap">
                <router-link :to="shareLink ? `/shared/${shareToken}` : '/'" class="nav-btn nav-btn--next">
                  <i class="fas fa-eye me-1"></i>{{ t('app.view') }}
                </router-link>
                <router-link to="/register" class="nav-btn nav-btn--submit">
                  <i class="fas fa-user-plus me-1"></i>{{ t('auth.registerBtn') }}
                </router-link>
                <router-link to="/" class="nav-btn nav-btn--back">
                  <i class="fas fa-home me-1"></i>{{ t('app.home') }}
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { getEmptyCVData } from '../stores/cv'
import { cvAPI, publicAPI } from '../services/api'
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
const isUniversityMember = ref(false)
const selectedFacultyId = ref<number | null>(null)
const selectedDepartmentId = ref<number | null>(null)
const faculties = ref<any[]>([])
const departments = ref<any[]>([])
const agreedToTerms = ref(false)

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

async function fetchFaculties() {
  try {
    const res = await publicAPI.getFaculties()
    faculties.value = res.data.faculties || res.data || []
  } catch {
    faculties.value = []
  }
}

async function fetchDepartments(facultyId: number) {
  try {
    const res = await publicAPI.getDepartments(facultyId)
    departments.value = res.data.departments || res.data || []
  } catch {
    departments.value = []
  }
}

watch(selectedFacultyId, (newVal) => {
  selectedDepartmentId.value = null
  departments.value = []
  if (newVal) {
    fetchDepartments(newVal)
  }
})

onMounted(() => {
  fetchFaculties()
})

function copyLink() {
  navigator.clipboard.writeText(shareLink.value)
}

async function handleSubmit() {
  // Use personal info name as fallback for guest name
  if (!guestName.value) {
    guestName.value = form.data.personal_info.full_name
  }
  if (!guestName.value) {
    toast.warning(t('cv.guestNameRequired') || 'Please enter your name')
    currentStep.value = 0
    return
  }
  // Set title from guest name or personal info name
  form.title = form.data.personal_info.full_name || guestName.value
  saving.value = true
  try {
    const payload = {
      title: form.title,
      language: form.language,
      template: form.template,
      data: form.data,
      guest_name: guestName.value,
      guest_email: guestEmail.value,
      is_university_member: isUniversityMember.value,
      faculty_id: isUniversityMember.value ? selectedFacultyId.value : null,
      department_id: isUniversityMember.value ? selectedDepartmentId.value : null,
    }
    const res = await cvAPI.createGuest(payload)
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

<style scoped>
.guest-form-page {
  min-height: 100vh;
  background: #ffffff;
}

.page-title {
  font-weight: 700;
  color: #222222;
}

.page-step {
  color: #6a6a6a;
  font-weight: 500;
}

/* Guest Warning Card */
.guest-warning-card {
  display: flex;
  gap: 1rem;
  padding: 1.25rem;
  border-radius: 12px;
  background: #ffffff;
  border-left: 4px solid #c0982b;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px;
}

.guest-warning-icon {
  width: 44px;
  height: 44px;
  min-width: 44px;
  border-radius: 50%;
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
}

.guest-warning-title {
  font-weight: 600;
  color: #222222;
  margin-bottom: 0.25rem;
}

.guest-warning-text {
  color: #6a6a6a;
  font-size: 0.88rem;
  line-height: 1.5;
}

.guest-cta-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 1rem;
  border-radius: 8px;
  font-size: 0.82rem;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;
  border: none;
}

.guest-cta-btn--primary {
  background: #1a5276;
  color: #ffffff;
}

.guest-cta-btn--primary:hover {
  background: #154360;
  color: #ffffff;
}

.guest-cta-btn--outline {
  background: transparent;
  border: 1px solid #dddddd;
  color: #222222;
}

.guest-cta-btn--outline:hover {
  border-color: #222222;
  color: #222222;
}

/* Progress Track */
.progress-track {
  height: 4px;
  background: #f0f0f0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-track-fill {
  height: 100%;
  background: #1a5276;
  border-radius: 4px;
  transition: width 0.4s ease;
}

/* Steps Navigation */
.steps-nav {
  display: flex;
  gap: 0.35rem;
  overflow-x: auto;
  padding-bottom: 0.25rem;
  -webkit-overflow-scrolling: touch;
}

.step-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.35rem 0.65rem;
  border-radius: 20px;
  font-size: 0.78rem;
  font-weight: 500;
  border: 1px solid #dddddd;
  background: #ffffff;
  color: #6a6a6a;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  flex-shrink: 0;
}

.step-chip-num {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #f0f0f0;
  color: #6a6a6a;
  font-size: 0.7rem;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.step-chip--active {
  border-color: #222222;
  background: #222222;
  color: #ffffff;
}

.step-chip--active .step-chip-num {
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

.step-chip--done {
  border-color: #1e8449;
  color: #1e8449;
}

.step-chip--done .step-chip-num {
  background: #1e8449;
  color: #ffffff;
}

/* Step Card */
.step-card {
  background: #ffffff;
  border-radius: 12px;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.step-card-body {
  padding: 32px;
}

/* Step Header */
.step-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.step-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1rem;
}

.step-header h4 {
  margin: 0;
  color: #222222;
  font-weight: 700;
}

/* Field styles */
.field-label {
  font-weight: 600;
  color: #222222;
  font-size: 0.9rem;
}

.field-input {
  border: 1px solid #dddddd;
  border-radius: 8px;
  color: #222222;
  transition: border-color 0.2s;
}

.field-input:focus {
  border-color: #222222;
  box-shadow: none;
}

.section-divider {
  border-color: #ebebeb;
  margin: 1rem 0;
}

.section-heading {
  color: #222222;
  font-weight: 700;
}

.section-heading-icon {
  color: #c0982b;
}

.section-desc {
  color: #6a6a6a;
  font-size: 0.9rem;
}

.text-secondary-light {
  color: #6a6a6a;
}

.text-accent {
  color: #c0982b !important;
}

/* Template Card */
.template-card {
  text-align: center;
  padding: 0.75rem 0.5rem;
  border-radius: 8px;
  border: 1px solid #dddddd;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #ffffff;
}

.template-card:hover {
  border-color: #222222;
}

.template-card--selected {
  border-color: #c0982b;
  background: rgba(192, 152, 43, 0.04);
}

/* Affiliation Card */
.affiliation-card {
  text-align: center;
  padding: 1.5rem 1rem;
  border-radius: 12px;
  border: 1px solid #dddddd;
  cursor: pointer;
  transition: all 0.2s ease;
  background: #ffffff;
  height: 100%;
}

.affiliation-card:hover {
  border-color: #222222;
}

.affiliation-card--selected {
  border-color: #1a5276;
  background: rgba(26, 82, 118, 0.03);
}

.affiliation-card h6 {
  color: #222222;
  font-weight: 600;
}

/* Section Entry */
.section-entry {
  border: 1px solid #ebebeb;
  border-radius: 12px;
  padding: 1rem;
  background: #fafafa;
  transition: border-color 0.2s;
}

.section-entry:hover {
  border-color: #dddddd;
}

.entry-label {
  font-weight: 600;
  color: #222222;
}

/* Buttons */
.btn-add {
  display: inline-flex;
  align-items: center;
  padding: 0.5rem 1.2rem;
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  border: 1px dashed #dddddd;
  background: transparent;
  color: #1a5276;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add:hover {
  border-color: #1a5276;
  background: rgba(26, 82, 118, 0.03);
}

.btn-delete {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.35rem 0.65rem;
  border-radius: 8px;
  font-size: 0.78rem;
  border: 1px solid #dddddd;
  background: #ffffff;
  color: #c0392b;
  cursor: pointer;
  transition: all 0.15s;
}

.btn-delete:hover {
  background: #c0392b;
  color: #ffffff;
  border-color: #c0392b;
}

/* Disclaimer */
.disclaimer-card {
  border-radius: 12px;
  padding: 1rem 1.25rem;
  background: #fafafa;
  border-left: 4px solid #c0982b;
}

.disclaimer-icon {
  width: 36px;
  height: 36px;
  min-width: 36px;
  border-radius: 50%;
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  margin-right: 0.75rem;
}

/* Navigation Buttons */
.nav-btn {
  display: inline-flex;
  align-items: center;
  padding: 0.55rem 1.5rem;
  border-radius: 8px;
  font-size: 0.88rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s ease;
}

.nav-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.nav-btn--back {
  background: transparent;
  border: 1px solid #dddddd;
  color: #6a6a6a;
}

.nav-btn--back:hover:not(:disabled) {
  border-color: #222222;
  color: #222222;
}

.nav-btn--next {
  background: #222222;
  color: #ffffff;
}

.nav-btn--next:hover:not(:disabled) {
  background: #000000;
  color: #ffffff;
}

.nav-btn--submit {
  background: #1a5276;
  color: #ffffff;
}

.nav-btn--submit:hover:not(:disabled) {
  background: #154360;
  color: #ffffff;
}

/* Success Modal */
.success-modal {
  border: none;
  border-radius: 12px;
  overflow: hidden;
  box-shadow:
    rgba(0, 0, 0, 0.02) 0px 0px 0px 1px,
    rgba(0, 0, 0, 0.04) 0px 2px 6px,
    rgba(0, 0, 0, 0.1) 0px 4px 8px;
}

.success-icon-wrap {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: rgba(30, 132, 73, 0.1);
  color: #1e8449;
  font-size: 1.8rem;
}

.success-title {
  font-weight: 700;
  color: #222222;
}

.share-link-group .form-control {
  border-radius: 8px 0 0 8px;
  border: 1px solid #dddddd;
}

.share-link-group .btn {
  border-radius: 0 8px 8px 0;
}

@media (max-width: 768px) {
  .step-card-body {
    padding: 20px;
  }
}
</style>
