<template>
  <div class="cv-form-page">
    <div class="container py-4" style="max-width: 800px;">

      <!-- Header -->
      <div class="text-center mb-4">
        <h3 class="fw-bold">{{ isEdit ? t('app.edit') : t('app.createCV') }}</h3>
        <p class="text-muted mb-0">{{ t('cv.step') }} {{ currentStep + 1 }} / {{ totalSteps }}</p>
      </div>

      <!-- Progress Bar -->
      <div class="progress mb-4" style="height: 8px; border-radius: 8px;">
        <div
          class="progress-bar"
          role="progressbar"
          :style="{ width: ((currentStep + 1) / totalSteps * 100) + '%' }"
          style="transition: width 0.4s ease; border-radius: 8px;"
        ></div>
      </div>

      <!-- Auto-save indicator -->
      <div v-if="lastSaved" class="text-end mb-2">
        <small class="text-muted">
          <i class="fas fa-check-circle text-success"></i>
          {{ locale === 'ar' ? 'حُفظ تلقائياً' : 'Auto-saved' }} {{ lastSaved }}
        </small>
      </div>

      <!-- Step Card -->
      <div class="card border-0 shadow-sm" style="border-radius: 16px;">
        <div class="card-body p-4 p-md-5">

          <!-- ===== Step 0: Template & Language ===== -->
          <div v-show="currentStep === 0">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-palette"></i></div>
              <h4>{{ t('cv.template') }}</h4>
            </div>

            <div class="row g-3 mb-4">
              <div class="col-md-6">
                <label class="form-label fw-semibold">{{ t('cv.language') }}</label>
                <select class="form-select form-select-lg" v-model="form.language">
                  <option value="ar">{{ t('app.arabic') }}</option>
                  <option value="en">{{ t('app.english') }}</option>
                </select>
              </div>
            </div>

            <label class="form-label fw-semibold">{{ t('cv.chooseTemplate') }}</label>
            <div class="row g-3">
              <div class="col-6 col-sm-4 col-md-3" v-for="tmpl in templates" :key="tmpl">
                <div
                  class="template-card"
                  :class="{ active: form.template === tmpl }"
                  @click="form.template = tmpl"
                >
                  <div class="template-preview">
                    <img :src="getTemplatePreviewDataUrl(tmpl)" :alt="t(`templates.${tmpl}`)" loading="lazy" />
                    <div class="template-check" v-if="form.template === tmpl">
                      <i class="fas fa-check"></i>
                    </div>
                  </div>
                  <div class="template-name">{{ t(`templates.${tmpl}`) }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- ===== Step 1: University Affiliation ===== -->
          <div v-show="currentStep === 1">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-university"></i></div>
              <h4>{{ t('cv.universityAffiliation') }}</h4>
            </div>

            <p class="text-muted mb-4">{{ t('cv.areYouMember') }}</p>

            <div class="row g-3 justify-content-center mb-4">
              <div class="col-6">
                <div
                  class="choice-card"
                  :class="{ active: isUniversityMember }"
                  @click="isUniversityMember = true"
                >
                  <i class="fas fa-university fa-2x mb-2"></i>
                  <div class="fw-semibold">{{ t('cv.yesMember') }}</div>
                  <small class="text-muted">{{ t('cv.memberDescription') }}</small>
                </div>
              </div>
              <div class="col-6">
                <div
                  class="choice-card"
                  :class="{ active: !isUniversityMember }"
                  @click="isUniversityMember = false; selectedFacultyId = null; selectedDepartmentId = null"
                >
                  <i class="fas fa-user-tie fa-2x mb-2"></i>
                  <div class="fw-semibold">{{ t('cv.noExternal') }}</div>
                  <small class="text-muted">{{ t('cv.externalDescription') }}</small>
                </div>
              </div>
            </div>

            <div v-if="isUniversityMember" class="row g-3">
              <div class="col-md-6">
                <label class="form-label fw-semibold">{{ t('cv.selectFaculty') }} *</label>
                <select class="form-select" v-model="selectedFacultyId">
                  <option :value="null" disabled>-- {{ t('cv.selectFaculty') }} --</option>
                  <option v-for="fac in faculties" :key="fac.id" :value="fac.id">{{ fac.name_ar || fac.name }}</option>
                </select>
              </div>
              <div class="col-md-6">
                <label class="form-label fw-semibold">{{ t('cv.selectDepartment') }} *</label>
                <select class="form-select" v-model="selectedDepartmentId" :disabled="!selectedFacultyId">
                  <option :value="null" disabled>-- {{ t('cv.selectDepartment') }} --</option>
                  <option v-for="dept in departments" :key="dept.id" :value="dept.id">{{ dept.name_ar || dept.name }}</option>
                </select>
              </div>
            </div>
          </div>

          <!-- ===== Step 2: Personal Info + Photo ===== -->
          <div v-show="currentStep === 2">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-user"></i></div>
              <h4>{{ t('cv.personalInfo') }}</h4>
            </div>

            <!-- Photo Upload -->
            <div class="text-center mb-4">
              <div v-if="form.data.photo" class="mb-2">
                <img :src="form.data.photo" class="rounded-circle shadow-sm" style="width:100px;height:100px;object-fit:cover" />
              </div>
              <div v-else class="mb-2">
                <div class="photo-placeholder"><i class="fas fa-camera"></i></div>
              </div>
              <input type="file" ref="photoInput" accept="image/*" class="d-none" @change="handlePhotoUpload" />
              <button @click="($refs.photoInput as HTMLInputElement).click()" class="btn btn-sm btn-outline-primary me-1">
                <i class="fas fa-upload me-1"></i>{{ t('cv.uploadPhoto') }}
              </button>
              <button v-if="form.data.photo" @click="form.data.photo = ''" class="btn btn-sm btn-outline-danger">
                <i class="fas fa-times"></i>
              </button>
            </div>

            <div class="row g-3">
              <div class="col-md-6">
                <label class="form-label">{{ t('cv.fullName') }} <span class="text-danger">*</span></label>
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
            </div>
          </div>

          <!-- ===== Step 3: Objective ===== -->
          <div v-show="currentStep === 3">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-bullseye"></i></div>
              <h4>{{ t('cv.objective') }}</h4>
            </div>
            <textarea class="form-control" rows="6" v-model="form.data.objective" :placeholder="t('cv.objective')"></textarea>
            <button @click="improveText('objective')" class="btn btn-sm btn-outline-info mt-2" :disabled="aiLoading">
              <i class="fas fa-magic me-1"></i>{{ t('ai.improveText') }}
            </button>
          </div>

          <!-- ===== Step 4: Experience & Education ===== -->
          <div v-show="currentStep === 4">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-briefcase"></i></div>
              <h4>{{ t('cv.experience') }} & {{ t('cv.education') }}</h4>
            </div>

            <!-- Experience -->
            <h6 class="fw-bold mb-3"><i class="fas fa-briefcase text-primary me-2"></i>{{ t('cv.experience') }}</h6>
            <div v-for="(exp, idx) in form.data.experiences" :key="'exp'+idx" class="item-card mb-3">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span class="badge bg-primary">{{ idx + 1 }}</span>
                <button @click="form.data.experiences.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-trash"></i></button>
              </div>
              <div class="row g-2">
                <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.position')" v-model="exp.title" /></div>
                <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.companyName')" v-model="exp.company" /></div>
                <div class="col-md-4"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.location')" v-model="exp.location" /></div>
                <div class="col-md-3"><input type="month" class="form-control form-control-sm" v-model="exp.start_date" /></div>
                <div class="col-md-3"><input type="month" class="form-control form-control-sm" v-model="exp.end_date" :disabled="exp.current" /></div>
                <div class="col-md-2 d-flex align-items-center">
                  <div class="form-check"><input type="checkbox" class="form-check-input" v-model="exp.current" /><label class="form-check-label small">{{ t('cv.current') }}</label></div>
                </div>
                <div class="col-12"><textarea class="form-control form-control-sm" rows="2" :placeholder="t('cv.description')" v-model="exp.description"></textarea></div>
              </div>
            </div>
            <button @click="addExperience" class="btn btn-sm btn-outline-primary mb-4"><i class="fas fa-plus me-1"></i>{{ t('cv.addExperience') }}</button>

            <hr class="my-4" />

            <!-- Education -->
            <h6 class="fw-bold mb-3"><i class="fas fa-graduation-cap text-success me-2"></i>{{ t('cv.education') }}</h6>
            <div v-for="(edu, idx) in form.data.education" :key="'edu'+idx" class="item-card mb-3">
              <div class="d-flex justify-content-between align-items-center mb-2">
                <span class="badge bg-success">{{ idx + 1 }}</span>
                <button @click="form.data.education.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-trash"></i></button>
              </div>
              <div class="row g-2">
                <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.degree')" v-model="edu.degree" /></div>
                <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.institution')" v-model="edu.institution" /></div>
                <div class="col-md-4"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.location')" v-model="edu.location" /></div>
                <div class="col-md-3"><input type="month" class="form-control form-control-sm" v-model="edu.start_date" /></div>
                <div class="col-md-3"><input type="month" class="form-control form-control-sm" v-model="edu.end_date" /></div>
                <div class="col-md-2"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.gpa')" v-model="edu.gpa" /></div>
                <div class="col-12"><textarea class="form-control form-control-sm" rows="2" :placeholder="t('cv.description')" v-model="edu.description"></textarea></div>
              </div>
            </div>
            <button @click="addEducation" class="btn btn-sm btn-outline-success"><i class="fas fa-plus me-1"></i>{{ t('cv.addEducation') }}</button>
          </div>

          <!-- ===== Step 5: Skills & Languages ===== -->
          <div v-show="currentStep === 5">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-star"></i></div>
              <h4>{{ t('cv.skills') }} & {{ t('cv.languages') }}</h4>
            </div>

            <!-- Skills -->
            <h6 class="fw-bold mb-3"><i class="fas fa-cogs text-primary me-2"></i>{{ t('cv.skills') }}</h6>
            <div v-for="(skill, idx) in form.data.skills" :key="'sk'+idx" class="row g-2 mb-2 align-items-center">
              <div class="col"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.skillName')" v-model="skill.name" /></div>
              <div class="col">
                <select class="form-select form-select-sm" v-model="skill.level">
                  <option value="beginner">{{ t('cv.beginner') }}</option>
                  <option value="intermediate">{{ t('cv.intermediate') }}</option>
                  <option value="advanced">{{ t('cv.advanced') }}</option>
                  <option value="expert">{{ t('cv.expert') }}</option>
                </select>
              </div>
              <div class="col-auto"><button @click="form.data.skills.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-times"></i></button></div>
            </div>
            <button @click="form.data.skills.push({ name: '', level: 'intermediate' })" class="btn btn-sm btn-outline-primary mb-4"><i class="fas fa-plus me-1"></i>{{ t('cv.addSkill') }}</button>

            <hr class="my-4" />

            <!-- Languages -->
            <h6 class="fw-bold mb-3"><i class="fas fa-language text-info me-2"></i>{{ t('cv.languages') }}</h6>
            <div v-for="(lang, idx) in form.data.languages" :key="'lang'+idx" class="row g-2 mb-2 align-items-center">
              <div class="col"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.languageName')" v-model="lang.name" /></div>
              <div class="col">
                <select class="form-select form-select-sm" v-model="lang.level">
                  <option value="basic">{{ t('cv.basic') }}</option>
                  <option value="conversational">{{ t('cv.conversational') }}</option>
                  <option value="fluent">{{ t('cv.fluent') }}</option>
                  <option value="native">{{ t('cv.native') }}</option>
                </select>
              </div>
              <div class="col-auto"><button @click="form.data.languages.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-times"></i></button></div>
            </div>
            <button @click="form.data.languages.push({ name: '', level: 'conversational' })" class="btn btn-sm btn-outline-info"><i class="fas fa-plus me-1"></i>{{ t('cv.addLanguage') }}</button>
          </div>

          <!-- ===== Step 6: Additional (Optional) ===== -->
          <div v-show="currentStep === 6">
            <div class="step-header mb-4">
              <div class="step-icon"><i class="fas fa-plus-circle"></i></div>
              <h4>{{ t('cv.additionalInfo') }}</h4>
              <p class="text-muted small mb-0">{{ t('cv.optionalSection') }}</p>
            </div>

            <!-- Accordion for optional sections -->
            <div class="accordion" id="additionalAccordion">
              <!-- Links -->
              <div class="accordion-item border-0 mb-2">
                <h2 class="accordion-header">
                  <button class="accordion-button collapsed rounded" type="button" data-bs-toggle="collapse" data-bs-target="#linksCollapse">
                    <i class="fas fa-link text-primary me-2"></i>{{ t('cv.links') }}
                    <span class="badge bg-primary ms-2" v-if="form.data.links.length">{{ form.data.links.length }}</span>
                  </button>
                </h2>
                <div id="linksCollapse" class="accordion-collapse collapse" data-bs-parent="#additionalAccordion">
                  <div class="accordion-body">
                    <div v-for="(link, idx) in form.data.links" :key="'lnk'+idx" class="row g-2 mb-2 align-items-center">
                      <div class="col-md-4"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.linkTitle')" v-model="link.title" /></div>
                      <div class="col-md-4"><input type="url" class="form-control form-control-sm" :placeholder="t('cv.linkUrl')" v-model="link.url" /></div>
                      <div class="col-md-3">
                        <select class="form-select form-select-sm" v-model="link.type">
                          <option value="academic">{{ t('cv.academic') }}</option>
                          <option value="social">{{ t('cv.social') }}</option>
                          <option value="portfolio">{{ t('cv.portfolio') }}</option>
                        </select>
                      </div>
                      <div class="col-auto"><button @click="form.data.links.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-times"></i></button></div>
                    </div>
                    <button @click="form.data.links.push({ title: '', url: '', type: 'academic' })" class="btn btn-sm btn-outline-primary"><i class="fas fa-plus me-1"></i>{{ t('cv.addLink') }}</button>
                  </div>
                </div>
              </div>

              <!-- Projects -->
              <div class="accordion-item border-0 mb-2">
                <h2 class="accordion-header">
                  <button class="accordion-button collapsed rounded" type="button" data-bs-toggle="collapse" data-bs-target="#projectsCollapse">
                    <i class="fas fa-project-diagram text-warning me-2"></i>{{ t('cv.projects') }}
                    <span class="badge bg-warning ms-2" v-if="form.data.projects.length">{{ form.data.projects.length }}</span>
                  </button>
                </h2>
                <div id="projectsCollapse" class="accordion-collapse collapse" data-bs-parent="#additionalAccordion">
                  <div class="accordion-body">
                    <div v-for="(proj, idx) in form.data.projects" :key="'prj'+idx" class="item-card mb-3">
                      <div class="d-flex justify-content-between mb-2">
                        <span class="badge bg-warning">{{ idx + 1 }}</span>
                        <button @click="form.data.projects.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-trash"></i></button>
                      </div>
                      <div class="row g-2">
                        <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.projectName')" v-model="proj.name" /></div>
                        <div class="col-md-6"><input type="url" class="form-control form-control-sm" :placeholder="t('cv.url')" v-model="proj.url" /></div>
                        <div class="col-12"><textarea class="form-control form-control-sm" rows="2" :placeholder="t('cv.description')" v-model="proj.description"></textarea></div>
                      </div>
                    </div>
                    <button @click="addProject" class="btn btn-sm btn-outline-warning"><i class="fas fa-plus me-1"></i>{{ t('cv.addProject') }}</button>
                  </div>
                </div>
              </div>

              <!-- Certificates -->
              <div class="accordion-item border-0 mb-2">
                <h2 class="accordion-header">
                  <button class="accordion-button collapsed rounded" type="button" data-bs-toggle="collapse" data-bs-target="#certsCollapse">
                    <i class="fas fa-certificate text-success me-2"></i>{{ t('cv.certificates') }}
                    <span class="badge bg-success ms-2" v-if="form.data.certificates.length">{{ form.data.certificates.length }}</span>
                  </button>
                </h2>
                <div id="certsCollapse" class="accordion-collapse collapse" data-bs-parent="#additionalAccordion">
                  <div class="accordion-body">
                    <div v-for="(cert, idx) in form.data.certificates" :key="'crt'+idx" class="item-card mb-3">
                      <div class="d-flex justify-content-between mb-2">
                        <span class="badge bg-success">{{ idx + 1 }}</span>
                        <button @click="form.data.certificates.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-trash"></i></button>
                      </div>
                      <div class="row g-2">
                        <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.certName')" v-model="cert.name" /></div>
                        <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.certIssuer')" v-model="cert.issuer" /></div>
                        <div class="col-md-6"><input type="month" class="form-control form-control-sm" v-model="cert.date" /></div>
                        <div class="col-md-6"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.credentialId')" v-model="cert.credential_id" /></div>
                      </div>
                    </div>
                    <button @click="addCertificate" class="btn btn-sm btn-outline-success"><i class="fas fa-plus me-1"></i>{{ t('cv.addCertificate') }}</button>
                  </div>
                </div>
              </div>

              <!-- References -->
              <div class="accordion-item border-0 mb-2">
                <h2 class="accordion-header">
                  <button class="accordion-button collapsed rounded" type="button" data-bs-toggle="collapse" data-bs-target="#refsCollapse">
                    <i class="fas fa-user-friends text-info me-2"></i>{{ t('cv.references') }}
                    <span class="badge bg-info ms-2" v-if="form.data.references.length">{{ form.data.references.length }}</span>
                  </button>
                </h2>
                <div id="refsCollapse" class="accordion-collapse collapse" data-bs-parent="#additionalAccordion">
                  <div class="accordion-body">
                    <div v-for="(ref, idx) in form.data.references" :key="'ref'+idx" class="item-card mb-3">
                      <div class="d-flex justify-content-between mb-2">
                        <span class="badge bg-info">{{ idx + 1 }}</span>
                        <button @click="form.data.references.splice(idx, 1)" class="btn btn-sm btn-link text-danger p-0"><i class="fas fa-trash"></i></button>
                      </div>
                      <div class="row g-2">
                        <div class="col-md-4"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.refName')" v-model="ref.name" /></div>
                        <div class="col-md-4"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.refPosition')" v-model="ref.position" /></div>
                        <div class="col-md-4"><input type="text" class="form-control form-control-sm" :placeholder="t('cv.refCompany')" v-model="ref.company" /></div>
                        <div class="col-md-6"><input type="email" class="form-control form-control-sm" :placeholder="t('cv.email')" v-model="ref.email" /></div>
                        <div class="col-md-6"><input type="tel" class="form-control form-control-sm" :placeholder="t('cv.phone')" v-model="ref.phone" /></div>
                      </div>
                    </div>
                    <button @click="addReference" class="btn btn-sm btn-outline-info"><i class="fas fa-plus me-1"></i>{{ t('cv.addReference') }}</button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- ===== Step 7: Disclaimer & Save ===== -->
          <div v-show="currentStep === 7">
            <div class="step-header mb-4">
              <div class="step-icon bg-success"><i class="fas fa-check-circle"></i></div>
              <h4>{{ t('cv.reviewAndSave') }}</h4>
            </div>

            <!-- Summary -->
            <div class="summary-card mb-4">
              <div class="row g-3 text-center">
                <div class="col-4">
                  <div class="fw-bold fs-5">{{ form.data.personal_info.full_name || '—' }}</div>
                  <small class="text-muted">{{ t('cv.fullName') }}</small>
                </div>
                <div class="col-4">
                  <div class="fw-bold fs-5">{{ t(`templates.${form.template}`) }}</div>
                  <small class="text-muted">{{ t('cv.template') }}</small>
                </div>
                <div class="col-4">
                  <div class="fw-bold fs-5">{{ form.language === 'ar' ? t('app.arabic') : t('app.english') }}</div>
                  <small class="text-muted">{{ t('cv.language') }}</small>
                </div>
              </div>
              <hr />
              <div class="row g-2 text-center">
                <div class="col"><span class="badge bg-primary">{{ form.data.experiences.length }}</span> {{ t('cv.experience') }}</div>
                <div class="col"><span class="badge bg-success">{{ form.data.education.length }}</span> {{ t('cv.education') }}</div>
                <div class="col"><span class="badge bg-info">{{ form.data.skills.length }}</span> {{ t('cv.skills') }}</div>
              </div>
            </div>

            <!-- Legal Disclaimer -->
            <div class="disclaimer-card">
              <div class="d-flex align-items-start">
                <i class="fas fa-shield-alt text-warning fa-2x me-3 mt-1 flex-shrink-0"></i>
                <div>
                  <div class="form-check">
                    <input type="checkbox" class="form-check-input" id="disclaimerCheck" v-model="agreedToTerms" />
                    <label class="form-check-label" for="disclaimerCheck">{{ t('cv.disclaimer') }}</label>
                  </div>
                  <div v-if="!agreedToTerms" class="text-danger small mt-2">
                    <i class="fas fa-info-circle me-1"></i>{{ t('cv.mustAgree') }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- ===== Navigation ===== -->
          <div class="d-flex justify-content-between mt-4 pt-3 border-top">
            <button @click="currentStep--" class="btn btn-outline-secondary px-4" :disabled="currentStep === 0">
              <i class="fas fa-arrow-right me-1"></i>{{ t('app.previous') }}
            </button>
            <div>
              <button v-if="currentStep < totalSteps - 1" @click="currentStep++" class="btn btn-primary px-4">
                {{ t('app.next') }} <i class="fas fa-arrow-left ms-1"></i>
              </button>
              <button v-else @click="handleSubmit" class="btn btn-success btn-lg px-5" :disabled="saving || !agreedToTerms">
                <span v-if="saving" class="spinner-border spinner-border-sm me-1"></span>
                <i v-else class="fas fa-check-circle me-1"></i>{{ t('app.save') }}
              </button>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useCVStore, getEmptyCVData } from '../stores/cv'
import { cvAPI, aiAPI, publicAPI } from '../services/api'
import { useToast } from '../composables/useToast'
import { getTemplatePreviewDataUrl } from '../data/templatePreviews'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const cvStore = useCVStore()
const toast = useToast()

const isEdit = ref(false)
const currentStep = ref(0)
const totalSteps = 8
const saving = ref(false)
const aiLoading = ref(false)

// University affiliation
const isUniversityMember = ref(false)
const selectedFacultyId = ref<number | null>(null)
const selectedDepartmentId = ref<number | null>(null)
const faculties = ref<any[]>([])
const departments = ref<any[]>([])

// Legal disclaimer
const agreedToTerms = ref(false)

// Auto-save
const lastSaved = ref('')
let autoSaveTimer: ReturnType<typeof setInterval>

const templates = [
  'academic', 'ats', 'compact', 'creative', 'designer', 'elegant', 'engineer',
  'executive', 'lawyer', 'medical', 'minimalist', 'modern', 'professional',
  'standard', 'teacher', 'tech'
]

const form = reactive({
  title: '',
  language: 'ar',
  template: 'modern',
  data: getEmptyCVData(),
})

async function fetchFaculties() {
  try {
    const res = await publicAPI.getFaculties()
    faculties.value = res.data.faculties || []
  } catch {}
}

async function fetchDepartments(facultyId: number) {
  try {
    const res = await publicAPI.getDepartments(facultyId)
    departments.value = res.data.departments || []
  } catch {}
}

watch(selectedFacultyId, (newVal) => {
  selectedDepartmentId.value = null
  departments.value = []
  if (newVal) fetchDepartments(newVal)
})

onMounted(async () => {
  fetchFaculties()
  const id = route.params.id
  if (id) {
    isEdit.value = true
    const cv = await cvStore.fetchCV(Number(id))
    if (cv) {
      form.title = cv.title
      form.language = cv.language
      form.template = cv.template
      Object.assign(form.data, cv.data)
      if (cv.is_university_member) {
        isUniversityMember.value = true
        selectedFacultyId.value = cv.faculty_id || null
        selectedDepartmentId.value = cv.department_id || null
      }
    }
  } else {
    const draft = localStorage.getItem('cv-draft')
    if (draft) {
      try {
        const parsed = JSON.parse(draft)
        Object.assign(form.data, parsed)
        toast.info(locale.value === 'ar' ? 'تم استعادة المسودة' : 'Draft restored')
      } catch {}
    }
  }

  // Start auto-save timer (every 30 seconds)
  autoSaveTimer = setInterval(async () => {
    if (route.params.id) {
      // Existing CV - save to server silently
      try {
        await cvAPI.update(Number(route.params.id), { data: form.data })
        lastSaved.value = new Date().toLocaleTimeString()
      } catch {}
    } else {
      // New CV - save to localStorage
      localStorage.setItem('cv-draft', JSON.stringify(form.data))
      lastSaved.value = new Date().toLocaleTimeString()
    }
  }, 30000)
})

onUnmounted(() => {
  clearInterval(autoSaveTimer)
})

let saveTimeout: ReturnType<typeof setTimeout> | null = null
watch(form, () => {
  if (!isEdit.value) {
    if (saveTimeout) clearTimeout(saveTimeout)
    saveTimeout = setTimeout(() => localStorage.setItem('cv-draft', JSON.stringify(form.data)), 1000)
  }
}, { deep: true })

function handlePhotoUpload(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) { toast.warning(t('cv.photoTooLarge')); return }
  if (!file.type.startsWith('image/')) { toast.warning(t('cv.invalidPhotoType')); return }
  const reader = new FileReader()
  reader.onload = () => { form.data.photo = reader.result as string }
  reader.readAsDataURL(file)
}

function addExperience() {
  form.data.experiences.push({ title: '', company: '', location: '', start_date: '', end_date: '', current: false, description: '' })
}
function addEducation() {
  form.data.education.push({ degree: '', institution: '', location: '', start_date: '', end_date: '', gpa: '', description: '' })
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
  // Set title from person's name automatically
  form.title = form.data.personal_info.full_name || t('cv.untitled')
  if (!agreedToTerms.value) { toast.warning(t('cv.mustAgree')); return }
  saving.value = true
  try {
    const payload = {
      ...form,
      is_university_member: isUniversityMember.value,
      faculty_id: isUniversityMember.value ? selectedFacultyId.value : null,
      department_id: isUniversityMember.value ? selectedDepartmentId.value : null,
    }
    let cvId: number
    if (isEdit.value) {
      const cv = await cvStore.updateCV(Number(route.params.id), payload)
      cvId = cv.id
    } else {
      const cv = await cvStore.createCV(payload)
      cvId = cv.id
      localStorage.removeItem('cv-draft')
    }
    router.push(`/cv/${cvId}`)
  } catch {
    toast.error(t('app.error'))
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.cv-form-page {
  background: #f5f7fb;
  min-height: 100vh;
}

.step-header {
  text-align: center;
}

.step-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
  margin: 0 auto 12px;
}

.step-icon.bg-success {
  background: linear-gradient(135deg, #10b981, #059669) !important;
}

.template-card {
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  overflow: hidden;
  background: #ffffff;
}

.template-card:hover {
  border-color: #6366f1;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.12);
}

.template-card.active {
  border-color: #6366f1;
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.2);
}

.template-preview {
  position: relative;
  aspect-ratio: 160 / 220;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  overflow: hidden;
}

.template-preview img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.template-check {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: #6366f1;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.4);
}

.rtl .template-check {
  right: auto;
  left: 8px;
}

.template-name {
  padding: 10px 8px;
  text-align: center;
  font-size: 0.85rem;
  font-weight: 500;
  color: #334155;
}

.template-card.active .template-name {
  color: #6366f1;
  background: #f0f0ff;
}

.choice-card {
  border: 2px solid #e2e8f0;
  border-radius: 14px;
  padding: 24px 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  color: #64748b;
}

.choice-card:hover {
  border-color: #6366f1;
}

.choice-card.active {
  border-color: #6366f1;
  background: #f0f0ff;
  color: #6366f1;
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.15);
}

.photo-placeholder {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: #f1f5f9;
  border: 2px dashed #cbd5e1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  color: #94a3b8;
}

.item-card {
  background: #f8fafc;
  border: 1px solid #e8ecf3;
  border-radius: 10px;
  padding: 12px 16px;
}

.summary-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 24px;
}

.disclaimer-card {
  background: #fffbeb;
  border: 1px solid #fbbf24;
  border-radius: 14px;
  padding: 20px;
}

.accordion-button {
  background: #f8fafc;
  border: 1px solid #e8ecf3;
  font-weight: 600;
  font-size: 0.95rem;
}

.accordion-button:not(.collapsed) {
  background: #f0f0ff;
  color: #6366f1;
  box-shadow: none;
}

.accordion-body {
  padding-top: 16px;
}

.progress-bar {
  background: linear-gradient(90deg, #6366f1, #8b5cf6);
}
</style>
