<template>
  <div ref="cvContent" class="cv-template" :class="[`template-${template}`, { rtl: language === 'ar' }]"
       :dir="language === 'ar' ? 'rtl' : 'ltr'">

    <!-- Modern Template -->
    <div v-if="template === 'modern'" class="modern-layout">
      <div class="modern-sidebar">
        <div class="photo-section" v-if="data.photo">
          <img :src="data.photo" class="profile-photo" />
        </div>
        <div class="sidebar-section" v-if="data.personal_info.email || data.personal_info.phone">
          <h3>{{ language === 'ar' ? 'معلومات الاتصال' : 'Contact' }}</h3>
          <p v-if="data.personal_info.email"><i class="fas fa-envelope"></i> {{ data.personal_info.email }}</p>
          <p v-if="data.personal_info.phone"><i class="fas fa-phone"></i> {{ data.personal_info.phone }}</p>
          <p v-if="data.personal_info.address"><i class="fas fa-map-marker-alt"></i> {{ data.personal_info.address }}</p>
          <p v-if="data.personal_info.linkedin"><i class="fab fa-linkedin"></i> {{ data.personal_info.linkedin }}</p>
          <p v-if="data.personal_info.github"><i class="fab fa-github"></i> {{ data.personal_info.github }}</p>
        </div>
        <div class="sidebar-section" v-if="data.skills.length">
          <h3>{{ language === 'ar' ? 'المهارات' : 'Skills' }}</h3>
          <div v-for="s in data.skills" :key="s.name" class="skill-item">
            <span>{{ s.name }}</span>
            <div class="skill-bar"><div class="skill-fill" :style="{ width: skillPercent(s.level) }"></div></div>
          </div>
        </div>
        <div class="sidebar-section" v-if="data.languages.length">
          <h3>{{ language === 'ar' ? 'اللغات' : 'Languages' }}</h3>
          <p v-for="l in data.languages" :key="l.name">{{ l.name }} - {{ l.level }}</p>
        </div>
      </div>
      <div class="modern-main">
        <div class="name-header">
          <h1>{{ data.personal_info.full_name }}</h1>
          <h2 v-if="data.personal_info.job_title">{{ data.personal_info.job_title }}</h2>
        </div>
        <div class="section" v-if="data.objective">
          <h3>{{ language === 'ar' ? 'الهدف الوظيفي' : 'Career Objective' }}</h3>
          <p>{{ data.objective }}</p>
        </div>
        <div class="section" v-if="data.experiences.length">
          <h3>{{ language === 'ar' ? 'الخبرات العملية' : 'Work Experience' }}</h3>
          <div v-for="exp in data.experiences" :key="exp.title + exp.company" class="entry">
            <div class="entry-header">
              <strong>{{ exp.title }}</strong>
              <span class="date">{{ exp.start_date }} - {{ exp.current ? (language === 'ar' ? 'حالياً' : 'Present') : exp.end_date }}</span>
            </div>
            <div class="entry-sub">{{ exp.company }} {{ exp.location ? `| ${exp.location}` : '' }}</div>
            <p v-if="exp.description">{{ exp.description }}</p>
          </div>
        </div>
        <div class="section" v-if="data.education.length">
          <h3>{{ language === 'ar' ? 'التعليم' : 'Education' }}</h3>
          <div v-for="edu in data.education" :key="edu.degree + edu.institution" class="entry">
            <div class="entry-header">
              <strong>{{ edu.degree }}</strong>
              <span class="date">{{ edu.start_date }} - {{ edu.end_date }}</span>
            </div>
            <div class="entry-sub">{{ edu.institution }} {{ edu.gpa ? `| GPA: ${edu.gpa}` : '' }}</div>
            <p v-if="edu.description">{{ edu.description }}</p>
          </div>
        </div>
        <div class="section" v-if="data.projects.length">
          <h3>{{ language === 'ar' ? 'المشاريع' : 'Projects' }}</h3>
          <div v-for="p in data.projects" :key="p.name" class="entry">
            <strong>{{ p.name }}</strong>
            <p>{{ p.description }}</p>
          </div>
        </div>
        <div class="section" v-if="data.certificates.length">
          <h3>{{ language === 'ar' ? 'الشهادات' : 'Certificates' }}</h3>
          <div v-for="c in data.certificates" :key="c.name" class="entry">
            <strong>{{ c.name }}</strong> - {{ c.issuer }} ({{ c.date }})
          </div>
        </div>
        <div class="section" v-if="data.references.length">
          <h3>{{ language === 'ar' ? 'المراجع' : 'References' }}</h3>
          <div v-for="r in data.references" :key="r.name" class="entry">
            <strong>{{ r.name }}</strong> - {{ r.position }}, {{ r.company }}<br/>
            <small>{{ r.email }} | {{ r.phone }}</small>
          </div>
        </div>
      </div>
    </div>

    <!-- For other templates, use a generic layout with different styling -->
    <div v-else class="generic-layout" :class="`style-${template}`">
      <!-- Header -->
      <div class="cv-header">
        <div class="header-content">
          <img v-if="data.photo" :src="data.photo" class="header-photo" />
          <div class="header-text">
            <h1>{{ data.personal_info.full_name }}</h1>
            <h2 v-if="data.personal_info.job_title">{{ data.personal_info.job_title }}</h2>
            <div class="contact-line">
              <span v-if="data.personal_info.email"><i class="fas fa-envelope"></i> {{ data.personal_info.email }}</span>
              <span v-if="data.personal_info.phone"><i class="fas fa-phone"></i> {{ data.personal_info.phone }}</span>
              <span v-if="data.personal_info.address"><i class="fas fa-map-marker-alt"></i> {{ data.personal_info.address }}</span>
              <span v-if="data.personal_info.linkedin"><i class="fab fa-linkedin"></i> {{ data.personal_info.linkedin }}</span>
              <span v-if="data.personal_info.github"><i class="fab fa-github"></i> {{ data.personal_info.github }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Body -->
      <div class="cv-body">
        <div class="section" v-if="data.objective">
          <h3 class="section-title">{{ language === 'ar' ? 'الهدف الوظيفي' : 'Career Objective' }}</h3>
          <p>{{ data.objective }}</p>
        </div>

        <div class="section" v-if="data.experiences.length">
          <h3 class="section-title">{{ language === 'ar' ? 'الخبرات العملية' : 'Work Experience' }}</h3>
          <div v-for="exp in data.experiences" :key="exp.title" class="entry">
            <div class="entry-header">
              <strong>{{ exp.title }}</strong> | {{ exp.company }}
              <span class="float-end">{{ exp.start_date }} - {{ exp.current ? (language === 'ar' ? 'حالياً' : 'Present') : exp.end_date }}</span>
            </div>
            <p class="entry-location" v-if="exp.location">{{ exp.location }}</p>
            <p v-if="exp.description">{{ exp.description }}</p>
          </div>
        </div>

        <div class="section" v-if="data.education.length">
          <h3 class="section-title">{{ language === 'ar' ? 'التعليم' : 'Education' }}</h3>
          <div v-for="edu in data.education" :key="edu.degree" class="entry">
            <div class="entry-header">
              <strong>{{ edu.degree }}</strong> | {{ edu.institution }}
              <span class="float-end">{{ edu.start_date }} - {{ edu.end_date }}</span>
            </div>
            <p v-if="edu.gpa">GPA: {{ edu.gpa }}</p>
            <p v-if="edu.description">{{ edu.description }}</p>
          </div>
        </div>

        <div class="section columns" v-if="data.skills.length || data.languages.length">
          <div class="col" v-if="data.skills.length">
            <h3 class="section-title">{{ language === 'ar' ? 'المهارات' : 'Skills' }}</h3>
            <div v-for="s in data.skills" :key="s.name" class="skill-tag">
              {{ s.name }} <small>({{ s.level }})</small>
            </div>
          </div>
          <div class="col" v-if="data.languages.length">
            <h3 class="section-title">{{ language === 'ar' ? 'اللغات' : 'Languages' }}</h3>
            <p v-for="l in data.languages" :key="l.name">{{ l.name }} - {{ l.level }}</p>
          </div>
        </div>

        <div class="section" v-if="data.projects.length">
          <h3 class="section-title">{{ language === 'ar' ? 'المشاريع' : 'Projects' }}</h3>
          <div v-for="p in data.projects" :key="p.name" class="entry">
            <strong>{{ p.name }}</strong>
            <p>{{ p.description }}</p>
          </div>
        </div>

        <div class="section" v-if="data.certificates.length">
          <h3 class="section-title">{{ language === 'ar' ? 'الشهادات' : 'Certificates' }}</h3>
          <div v-for="c in data.certificates" :key="c.name" class="entry">
            <strong>{{ c.name }}</strong> - {{ c.issuer }} ({{ c.date }})
          </div>
        </div>

        <div class="section" v-if="data.links.length">
          <h3 class="section-title">{{ language === 'ar' ? 'الروابط' : 'Links' }}</h3>
          <div v-for="l in data.links" :key="l.url" class="entry">
            <strong>{{ l.title }}</strong>: {{ l.url }}
          </div>
        </div>

        <div class="section" v-if="data.references.length">
          <h3 class="section-title">{{ language === 'ar' ? 'المراجع' : 'References' }}</h3>
          <div class="references-grid">
            <div v-for="r in data.references" :key="r.name" class="ref-card">
              <strong>{{ r.name }}</strong><br/>
              {{ r.position }}, {{ r.company }}<br/>
              <small>{{ r.email }} | {{ r.phone }}</small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CVData } from '../stores/cv'

defineProps<{
  data: CVData
  template: string
  language: string
}>()

function skillPercent(level: string): string {
  const map: Record<string, string> = {
    beginner: '25%', intermediate: '50%', advanced: '75%', expert: '95%'
  }
  return map[level] || '50%'
}
</script>

<style scoped>
.cv-template {
  font-family: 'Inter', 'Cairo', sans-serif;
  font-size: 11pt;
  line-height: 1.5;
  color: #333;
  background: white;
  min-height: 297mm;
  width: 210mm;
  margin: 0 auto;
}
.cv-template.rtl { font-family: 'Cairo', 'Inter', sans-serif; }

/* Modern Template */
.modern-layout { display: flex; min-height: 297mm; }
.modern-sidebar { width: 35%; background: #2c3e50; color: white; padding: 30px 20px; }
.modern-main { width: 65%; padding: 30px; }
.profile-photo { width: 120px; height: 120px; border-radius: 50%; object-fit: cover; display: block; margin: 0 auto 20px; border: 3px solid #3498db; }
.sidebar-section { margin-bottom: 20px; }
.sidebar-section h3 { font-size: 13pt; border-bottom: 2px solid #3498db; padding-bottom: 5px; margin-bottom: 10px; }
.sidebar-section p { font-size: 9pt; margin: 4px 0; word-break: break-all; }
.sidebar-section i { width: 16px; margin-right: 6px; }
.skill-item { margin-bottom: 8px; }
.skill-item span { font-size: 9pt; }
.skill-bar { height: 6px; background: #34495e; border-radius: 3px; margin-top: 2px; }
.skill-fill { height: 100%; background: #3498db; border-radius: 3px; }
.name-header h1 { font-size: 24pt; color: #2c3e50; margin: 0; }
.name-header h2 { font-size: 14pt; color: #3498db; margin: 0 0 20px; font-weight: 400; }
.section { margin-bottom: 20px; }
.section h3 { font-size: 14pt; color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 5px; margin-bottom: 10px; }
.entry { margin-bottom: 12px; }
.entry-header { display: flex; justify-content: space-between; }
.entry-sub { color: #666; font-size: 10pt; }
.date { color: #888; font-size: 9pt; }

/* Generic Layout Styles */
.generic-layout { padding: 20px 30px; }
.cv-header { margin-bottom: 20px; padding-bottom: 15px; border-bottom: 2px solid #333; }
.header-content { display: flex; align-items: center; gap: 20px; }
.header-photo { width: 80px; height: 80px; border-radius: 50%; object-fit: cover; }
.header-text h1 { font-size: 22pt; margin: 0; }
.header-text h2 { font-size: 13pt; color: #666; margin: 0; font-weight: 400; }
.contact-line { display: flex; flex-wrap: wrap; gap: 12px; margin-top: 8px; font-size: 9pt; color: #555; }
.contact-line i { margin-right: 3px; }
.section-title { font-size: 13pt; border-bottom: 1px solid #ddd; padding-bottom: 4px; margin-bottom: 10px; text-transform: uppercase; letter-spacing: 1px; }
.columns { display: flex; gap: 30px; }
.columns .col { flex: 1; }
.skill-tag { display: inline-block; background: #e9ecef; padding: 3px 10px; border-radius: 12px; margin: 2px 4px; font-size: 9pt; }
.references-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.ref-card { border: 1px solid #ddd; padding: 10px; border-radius: 4px; font-size: 9pt; }
.entry-location { color: #888; font-size: 9pt; margin: 2px 0; }

/* Professional */
.style-professional .cv-header { background: #1a237e; color: white; padding: 25px 30px; margin: -20px -30px 20px; }
.style-professional .header-text h1 { color: white; }
.style-professional .header-text h2 { color: #b3c5ff; }
.style-professional .contact-line { color: #ccc; }
.style-professional .section-title { color: #1a237e; border-color: #1a237e; }

/* Minimalist */
.style-minimalist .cv-header { border-bottom: 1px solid #eee; }
.style-minimalist .section-title { text-transform: none; letter-spacing: 0; border: none; font-weight: 600; }
.style-minimalist .header-photo { display: none; }

/* Academic */
.style-academic { font-family: 'Georgia', 'Cairo', serif; }
.style-academic .section-title { font-variant: small-caps; letter-spacing: 2px; }
.style-academic .cv-header { text-align: center; }
.style-academic .header-content { flex-direction: column; }

/* Creative */
.style-creative .cv-header { background: linear-gradient(135deg, #667eea, #764ba2); color: white; padding: 30px; margin: -20px -30px 20px; border-radius: 0 0 20px 20px; }
.style-creative .header-text h1 { color: white; }
.style-creative .header-text h2 { color: #e0d4ff; }
.style-creative .contact-line { color: #ddd; }
.style-creative .section-title { color: #667eea; border-color: #667eea; }
.style-creative .skill-tag { background: #667eea; color: white; }

/* ATS */
.style-ats { font-family: 'Arial', sans-serif; }
.style-ats .header-photo { display: none; }
.style-ats .cv-header { border-bottom: 2px solid #000; }
.style-ats .section-title { text-transform: uppercase; border-bottom: 1px solid #000; color: #000; }
.style-ats .skill-tag { background: none; border: none; padding: 0; border-radius: 0; }

/* Compact */
.style-compact { font-size: 9pt; line-height: 1.3; padding: 15px 20px; }
.style-compact .section { margin-bottom: 10px; }
.style-compact .entry { margin-bottom: 6px; }

/* Elegant */
.style-elegant .cv-header { border-bottom: 3px double #b8860b; }
.style-elegant .section-title { color: #b8860b; border-color: #b8860b; }
.style-elegant .header-photo { border: 2px solid #b8860b; }

/* Executive */
.style-executive .cv-header { background: #1b2838; color: white; padding: 30px; margin: -20px -30px 20px; }
.style-executive .header-text h1 { color: white; font-size: 26pt; }
.style-executive .header-text h2 { color: #8899aa; }
.style-executive .contact-line { color: #aab; }
.style-executive .section-title { color: #1b2838; }

/* Tech */
.style-tech { font-family: 'Courier New', 'Cairo', monospace; }
.style-tech .cv-header { background: #0d1117; color: #58a6ff; padding: 25px; margin: -20px -30px 20px; }
.style-tech .header-text h1 { color: #f0f6fc; }
.style-tech .header-text h2 { color: #58a6ff; }
.style-tech .contact-line { color: #8b949e; }
.style-tech .section-title { color: #58a6ff; border-color: #30363d; }
.style-tech .skill-tag { background: #21262d; color: #58a6ff; border: 1px solid #30363d; }

/* Designer */
.style-designer .cv-header { background: linear-gradient(45deg, #ff6b6b, #feca57); padding: 30px; margin: -20px -30px 20px; }
.style-designer .header-text h1 { color: white; font-size: 28pt; }
.style-designer .section-title { color: #ff6b6b; }

/* Medical */
.style-medical .cv-header { border-bottom: 3px solid #2e7d32; }
.style-medical .section-title { color: #2e7d32; border-color: #2e7d32; }
.style-medical .header-photo { border: 2px solid #2e7d32; }

/* Engineer */
.style-engineer .cv-header { border-bottom: 3px solid #1565c0; border-top: 3px solid #1565c0; padding: 20px 0; }
.style-engineer .section-title { color: #1565c0; border-color: #1565c0; }

/* Standard */
.style-standard .cv-header { text-align: center; }
.style-standard .header-content { flex-direction: column; }
.style-standard .contact-line { justify-content: center; }

/* Teacher */
.style-teacher .cv-header { background: #f57c00; color: white; padding: 25px; margin: -20px -30px 20px; }
.style-teacher .header-text h1 { color: white; }
.style-teacher .header-text h2 { color: #ffe0b2; }
.style-teacher .section-title { color: #f57c00; border-color: #f57c00; }

/* Lawyer */
.style-lawyer { font-family: 'Georgia', 'Cairo', serif; }
.style-lawyer .cv-header { border-bottom: 3px double #333; }
.style-lawyer .section-title { font-variant: small-caps; color: #333; border-color: #333; }

/* RTL adjustments */
.rtl .sidebar-section i { margin-right: 0; margin-left: 6px; }
.rtl .contact-line i { margin-right: 0; margin-left: 3px; }
.rtl .entry-header { flex-direction: row-reverse; }
.rtl .float-end { float: left !important; }
.rtl .modern-layout { flex-direction: row-reverse; }

@media print {
  .cv-template { box-shadow: none; margin: 0; width: 100%; }
  .modern-sidebar { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
}
</style>
