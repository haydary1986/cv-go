<template>
  <div class="container py-4">
    <h3 class="mb-4">{{ isEdit ? t('app.edit') : t('app.createCV') }}</h3>

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

        <!-- Step 0: Basic Info -->
        <div v-show="currentStep === 0">
          <div class="row g-3">
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
            <!-- Template Preview Cards -->
            <div class="col-12">
              <label class="form-label">{{ t('cv.chooseTemplate') }}</label>
              <div class="row g-2">
                <div class="col-4 col-md-2" v-for="tmpl in templates" :key="tmpl">
                  <div class="card text-center p-2 cursor-pointer" :class="{ 'border-primary bg-primary bg-opacity-10': form.template === tmpl }" @click="form.template = tmpl" style="cursor:pointer">
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
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.website') }}</label>
              <input type="url" class="form-control" v-model="form.data.personal_info.website" />
            </div>
            <div class="col-md-6">
              <label class="form-label">{{ t('cv.github') }}</label>
              <input type="url" class="form-control" v-model="form.data.personal_info.github" />
            </div>
          </div>
        </div>

        <!-- Step 2: Photo -->
        <div v-show="currentStep === 2">
          <div class="text-center">
            <div v-if="form.data.photo" class="mb-3">
              <img :src="form.data.photo" class="rounded-circle" style="width:150px;height:150px;object-fit:cover" />
            </div>
            <div v-else class="mb-3">
              <i class="fas fa-user-circle fa-5x text-muted"></i>
            </div>
            <input type="file" ref="photoInput" accept="image/*" class="d-none" @change="handlePhotoUpload" />
            <button @click="($refs.photoInput as HTMLInputElement).click()" class="btn btn-outline-primary me-2">
              <i class="fas fa-upload me-1"></i>{{ t('cv.uploadPhoto') }}
            </button>
            <button v-if="form.data.photo" @click="form.data.photo = ''" class="btn btn-outline-danger">
              <i class="fas fa-times me-1"></i>{{ t('cv.removePhoto') }}
            </button>
          </div>
        </div>

        <!-- Step 3: Objective -->
        <div v-show="currentStep === 3">
          <label class="form-label">{{ t('cv.objective') }}</label>
          <textarea class="form-control" rows="5" v-model="form.data.objective"></textarea>
          <button @click="improveText('objective')" class="btn btn-sm btn-outline-info mt-2" :disabled="aiLoading">
            <i class="fas fa-magic me-1"></i>{{ t('ai.improveText') }}
          </button>
        </div>

        <!-- Step 4: Experience -->
        <div v-show="currentStep === 4">
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

        <!-- Step 5: Education -->
        <div v-show="currentStep === 5">
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
              <div class="col-12">
                <textarea class="form-control" rows="2" :placeholder="t('cv.description')" v-model="edu.description"></textarea>
              </div>
            </div>
          </div>
          <button @click="addEducation" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addEducation') }}
          </button>
        </div>

        <!-- Step 6: Skills -->
        <div v-show="currentStep === 6">
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

        <!-- Step 7: Languages -->
        <div v-show="currentStep === 7">
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

        <!-- Step 8: Links -->
        <div v-show="currentStep === 8">
          <div v-for="(link, idx) in form.data.links" :key="idx" class="row g-2 mb-2">
            <div class="col-md-4">
              <input type="text" class="form-control" :placeholder="t('cv.linkTitle')" v-model="link.title" />
            </div>
            <div class="col-md-4">
              <input type="url" class="form-control" :placeholder="t('cv.linkUrl')" v-model="link.url" />
            </div>
            <div class="col-md-3">
              <select class="form-select" v-model="link.type">
                <option value="academic">{{ t('cv.academic') }}</option>
                <option value="social">{{ t('cv.social') }}</option>
                <option value="portfolio">{{ t('cv.portfolio') }}</option>
              </select>
            </div>
            <div class="col-md-1">
              <button @click="form.data.links.splice(idx, 1)" class="btn btn-outline-danger">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
          <button @click="form.data.links.push({ title: '', url: '', type: 'academic' })" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addLink') }}
          </button>
        </div>

        <!-- Step 9: Projects -->
        <div v-show="currentStep === 9">
          <div v-for="(proj, idx) in form.data.projects" :key="idx" class="border rounded p-3 mb-3">
            <div class="d-flex justify-content-between mb-2">
              <h6>{{ t('cv.projects') }} {{ idx + 1 }}</h6>
              <button @click="form.data.projects.splice(idx, 1)" class="btn btn-sm btn-outline-danger">
                <i class="fas fa-trash"></i>
              </button>
            </div>
            <div class="row g-2">
              <div class="col-md-6">
                <input type="text" class="form-control" :placeholder="t('cv.projectName')" v-model="proj.name" />
              </div>
              <div class="col-md-6">
                <input type="url" class="form-control" placeholder="URL" v-model="proj.url" />
              </div>
              <div class="col-md-3">
                <input type="month" class="form-control" v-model="proj.start_date" />
              </div>
              <div class="col-md-3">
                <input type="month" class="form-control" v-model="proj.end_date" />
              </div>
              <div class="col-12">
                <textarea class="form-control" rows="2" :placeholder="t('cv.description')" v-model="proj.description"></textarea>
              </div>
            </div>
          </div>
          <button @click="addProject" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addProject') }}
          </button>
        </div>

        <!-- Step 10: Certificates -->
        <div v-show="currentStep === 10">
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
                <input type="month" class="form-control" :placeholder="t('cv.certDate')" v-model="cert.date" />
              </div>
              <div class="col-md-4">
                <input type="text" class="form-control" placeholder="Credential ID" v-model="cert.credential_id" />
              </div>
              <div class="col-md-4">
                <input type="url" class="form-control" placeholder="URL" v-model="cert.url" />
              </div>
            </div>
          </div>
          <button @click="addCertificate" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addCertificate') }}
          </button>
        </div>

        <!-- Step 11: References -->
        <div v-show="currentStep === 11">
          <div v-for="(ref, idx) in form.data.references" :key="idx" class="border rounded p-3 mb-3">
            <div class="d-flex justify-content-between mb-2">
              <h6>{{ t('cv.references') }} {{ idx + 1 }}</h6>
              <button @click="form.data.references.splice(idx, 1)" class="btn btn-sm btn-outline-danger">
                <i class="fas fa-trash"></i>
              </button>
            </div>
            <div class="row g-2">
              <div class="col-md-4">
                <input type="text" class="form-control" :placeholder="t('cv.refName')" v-model="ref.name" />
              </div>
              <div class="col-md-4">
                <input type="text" class="form-control" :placeholder="t('cv.refPosition')" v-model="ref.position" />
              </div>
              <div class="col-md-4">
                <input type="text" class="form-control" :placeholder="t('cv.refCompany')" v-model="ref.company" />
              </div>
              <div class="col-md-6">
                <input type="email" class="form-control" :placeholder="t('cv.email')" v-model="ref.email" />
              </div>
              <div class="col-md-6">
                <input type="tel" class="form-control" :placeholder="t('cv.phone')" v-model="ref.phone" />
              </div>
            </div>
          </div>
          <button @click="addReference" class="btn btn-outline-primary">
            <i class="fas fa-plus me-1"></i>{{ t('cv.addReference') }}
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useCVStore, getEmptyCVData } from '../stores/cv'
import { aiAPI } from '../services/api'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const cvStore = useCVStore()

const isEdit = ref(false)
const currentStep = ref(0)
const saving = ref(false)
const aiLoading = ref(false)

const templates = [
  'academic', 'ats', 'compact', 'creative', 'designer', 'elegant', 'engineer',
  'executive', 'lawyer', 'medical', 'minimalist', 'modern', 'professional',
  'standard', 'teacher', 'tech'
]

const steps = [
  t('cv.template'), t('cv.personalInfo'), t('cv.photo'), t('cv.objective'),
  t('cv.experience'), t('cv.education'), t('cv.skills'), t('cv.languages'),
  t('cv.links'), t('cv.projects'), t('cv.certificates'), t('cv.references')
]

const form = reactive({
  title: '',
  language: 'ar',
  template: 'modern',
  data: getEmptyCVData(),
})

onMounted(async () => {
  const id = route.params.id
  if (id) {
    isEdit.value = true
    const cv = await cvStore.fetchCV(Number(id))
    if (cv) {
      form.title = cv.title
      form.language = cv.language
      form.template = cv.template
      Object.assign(form.data, cv.data)
    }
  } else {
    // Restore from localStorage
    const saved = localStorage.getItem('cv_draft')
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        Object.assign(form, parsed)
      } catch {}
    }
  }
})

// Auto-save draft with debounce (1 second delay)
let saveTimeout: ReturnType<typeof setTimeout> | null = null
watch(form, () => {
  if (!isEdit.value) {
    if (saveTimeout) clearTimeout(saveTimeout)
    saveTimeout = setTimeout(() => {
      localStorage.setItem('cv_draft', JSON.stringify(form))
    }, 1000)
  }
}, { deep: true })

function handlePhotoUpload(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  // Validate file size (max 2MB)
  if (file.size > 2 * 1024 * 1024) {
    alert(t('cv.photoTooLarge') || 'Photo must be less than 2MB')
    return
  }
  // Validate file type
  if (!file.type.startsWith('image/')) {
    alert(t('cv.invalidPhotoType') || 'Please upload an image file')
    return
  }
  const reader = new FileReader()
  reader.onload = () => {
    form.data.photo = reader.result as string
  }
  reader.readAsDataURL(file)
}

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

function addProject() {
  form.data.projects.push({ name: '', description: '', url: '', start_date: '', end_date: '' })
}

function addCertificate() {
  form.data.certificates.push({ name: '', issuer: '', date: '', expiry_date: '', credential_id: '', url: '' })
}

function addReference() {
  form.data.references.push({ name: '', position: '', company: '', email: '', phone: '' })
}

async function improveText(field: string) {
  aiLoading.value = true
  try {
    const text = field === 'objective' ? form.data.objective : ''
    const res = await aiAPI.improveText({ text, language: form.language })
    if (field === 'objective') form.data.objective = res.data.result
  } catch {}
  aiLoading.value = false
}

async function handleSubmit() {
  if (!form.title) {
    currentStep.value = 0
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await cvStore.updateCV(Number(route.params.id), form)
    } else {
      await cvStore.createCV(form)
      localStorage.removeItem('cv_draft')
    }
    router.push('/dashboard')
  } catch (err) {
    console.error(err)
  } finally {
    saving.value = false
  }
}
</script>
