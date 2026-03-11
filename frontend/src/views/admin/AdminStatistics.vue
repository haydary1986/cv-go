<template>
  <div>
    <h3 class="mb-4">{{ t('admin.statistics') }}</h3>

    <div class="row g-3">
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">{{ t('admin.byFaculty') }}</div>
          <div class="card-body"><canvas ref="facultyChart"></canvas></div>
        </div>
      </div>
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">{{ t('admin.byLanguage') }}</div>
          <div class="card-body"><canvas ref="languageChart"></canvas></div>
        </div>
      </div>
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">{{ t('admin.byStatus') }}</div>
          <div class="card-body"><canvas ref="statusChart"></canvas></div>
        </div>
      </div>
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">{{ t('admin.monthlyTrends') }}</div>
          <div class="card-body"><canvas ref="trendChart"></canvas></div>
        </div>
      </div>
      <div class="col-12">
        <div class="card">
          <div class="card-header">{{ t('admin.topUsers') }}</div>
          <div class="card-body">
            <table class="table table-sm">
              <thead><tr><th>{{ t('auth.email') }}</th><th>Name</th><th>CVs</th></tr></thead>
              <tbody>
                <tr v-for="u in topUsers" :key="u.email">
                  <td>{{ u.email }}</td><td>{{ u.name }}</td><td><span class="badge bg-primary">{{ u.count }}</span></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '../../services/api'
import { Chart, registerables } from 'chart.js'
Chart.register(...registerables)

const { t } = useI18n()
const facultyChart = ref<HTMLCanvasElement>()
const languageChart = ref<HTMLCanvasElement>()
const statusChart = ref<HTMLCanvasElement>()
const trendChart = ref<HTMLCanvasElement>()
const topUsers = ref<any[]>([])

const colors = ['#0d6efd', '#198754', '#dc3545', '#ffc107', '#6610f2', '#0dcaf0', '#fd7e14', '#20c997']

onMounted(async () => {
  try {
    const res = await adminAPI.getStats()
    const data = res.data
    topUsers.value = data.top_users || []

    if (facultyChart.value && data.by_faculty?.length) {
      new Chart(facultyChart.value, {
        type: 'bar',
        data: {
          labels: data.by_faculty.map((f: any) => f.name || 'N/A'),
          datasets: [{ label: 'CVs', data: data.by_faculty.map((f: any) => f.count), backgroundColor: colors }]
        }
      })
    }

    if (languageChart.value && data.by_language?.length) {
      new Chart(languageChart.value, {
        type: 'pie',
        data: {
          labels: data.by_language.map((l: any) => l.language),
          datasets: [{ data: data.by_language.map((l: any) => l.count), backgroundColor: colors }]
        }
      })
    }

    if (statusChart.value && data.by_status?.length) {
      new Chart(statusChart.value, {
        type: 'doughnut',
        data: {
          labels: data.by_status.map((s: any) => s.status),
          datasets: [{ data: data.by_status.map((s: any) => s.count), backgroundColor: ['#6c757d', '#ffc107', '#198754', '#dc3545'] }]
        }
      })
    }

    if (trendChart.value && data.monthly_trends?.length) {
      new Chart(trendChart.value, {
        type: 'line',
        data: {
          labels: data.monthly_trends.map((t: any) => t.month),
          datasets: [{ label: 'CVs', data: data.monthly_trends.map((t: any) => t.count), borderColor: '#0d6efd', tension: 0.3, fill: true, backgroundColor: 'rgba(13,110,253,0.1)' }]
        }
      })
    }
  } catch {}
})
</script>
