<template>
  <div class="admin-statistics">
    <!-- Page Header -->
    <div class="page-header">
      <h2 class="page-title">{{ t('admin.statistics') }}</h2>
      <p class="page-subtitle">Data analytics and insights for the CV system</p>
    </div>

    <div class="row g-4">
      <!-- By Faculty -->
      <div class="col-lg-6">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-university me-2 chart-icon chart-icon--primary"></i>
              {{ t('admin.byFaculty') }}
            </div>
          </div>
          <div class="chart-card-body">
            <canvas ref="facultyChart"></canvas>
          </div>
        </div>
      </div>

      <!-- By Language -->
      <div class="col-lg-6">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-language me-2 chart-icon chart-icon--accent"></i>
              {{ t('admin.byLanguage') }}
            </div>
          </div>
          <div class="chart-card-body">
            <canvas ref="languageChart"></canvas>
          </div>
        </div>
      </div>

      <!-- By Status -->
      <div class="col-lg-6">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-chart-pie me-2 chart-icon chart-icon--green"></i>
              {{ t('admin.byStatus') }}
            </div>
          </div>
          <div class="chart-card-body">
            <canvas ref="statusChart"></canvas>
          </div>
        </div>
      </div>

      <!-- Monthly Trends -->
      <div class="col-lg-6">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-chart-line me-2 chart-icon chart-icon--purple"></i>
              {{ t('admin.monthlyTrends') }}
            </div>
          </div>
          <div class="chart-card-body">
            <canvas ref="trendChart"></canvas>
          </div>
        </div>
      </div>

      <!-- Top Users -->
      <div class="col-12">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-trophy me-2 chart-icon chart-icon--gold"></i>
              {{ t('admin.topUsers') }}
            </div>
          </div>
          <div class="chart-card-body p-0">
            <div class="table-responsive">
              <table class="table admin-table mb-0">
                <thead>
                  <tr>
                    <th style="width: 50px;">#</th>
                    <th>{{ t('auth.email') }}</th>
                    <th>{{ t('admin.name') }}</th>
                    <th class="text-end">{{ t('admin.cvCount') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(u, index) in topUsers" :key="u.email">
                    <td>
                      <span class="rank-badge" :class="'rank-badge--' + (index + 1)">{{ index + 1 }}</span>
                    </td>
                    <td class="text-muted">{{ u.email }}</td>
                    <td class="fw-semibold text-dark">{{ u.name }}</td>
                    <td class="text-end">
                      <span class="count-badge">{{ u.count }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
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

<style scoped>
.admin-statistics {
  max-width: 1400px;
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

/* ── Chart Cards ── */
.chart-card {
  background: #ffffff;
  border-radius: 12px;
  border: 1px solid #e9ecef;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  overflow: hidden;
  transition: box-shadow 0.2s;
}

.chart-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.chart-card-header {
  padding: 18px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.chart-card-title {
  font-size: 15px;
  font-weight: 600;
  color: #2c3e50;
  display: flex;
  align-items: center;
}

.chart-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
}

.chart-icon--primary { background: rgba(26, 82, 118, 0.1); color: #1a5276; }
.chart-icon--accent { background: rgba(192, 152, 43, 0.1); color: #c0982b; }
.chart-icon--green { background: rgba(25, 135, 84, 0.1); color: #198754; }
.chart-icon--purple { background: rgba(102, 16, 242, 0.1); color: #6610f2; }
.chart-icon--gold { background: rgba(192, 152, 43, 0.1); color: #c0982b; }

.chart-card-body {
  padding: 24px;
}

/* ── Table ── */
.admin-table {
  font-size: 14px;
}

.admin-table thead th {
  background: #f8f9fb;
  border-bottom: 2px solid #e9ecef;
  color: #1a5276;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 14px 20px;
}

.admin-table tbody td {
  padding: 12px 20px;
  vertical-align: middle;
  border-bottom: 1px solid #f0f2f5;
}

.admin-table tbody tr:hover {
  background: #f8fafc;
}

/* ── Rank Badge ── */
.rank-badge {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  background: #f0f2f5;
  color: #6c757d;
}

.rank-badge--1 { background: linear-gradient(135deg, #c0982b, #d4a93a); color: #fff; }
.rank-badge--2 { background: linear-gradient(135deg, #95a5a6, #bdc3c7); color: #fff; }
.rank-badge--3 { background: linear-gradient(135deg, #cd7f32, #d4943a); color: #fff; }

/* ── Count Badge ── */
.count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 700;
  background: rgba(26, 82, 118, 0.08);
  color: #1a5276;
}
</style>
