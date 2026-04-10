<template>
  <div class="admin-dashboard">
    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h2 class="page-title">{{ t('admin.dashboard') }}</h2>
        <p class="page-subtitle">Overview of your university CV system</p>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="row g-4 mb-4">
      <div class="col-sm-6 col-xl-3">
        <div class="stat-card">
          <div class="stat-card-body">
            <div class="stat-card-info">
              <span class="stat-card-label">{{ t('admin.totalUsers') }}</span>
              <h3 class="stat-card-value">{{ stats.total_users }}</h3>
            </div>
            <div class="stat-card-icon stat-card-icon--blue">
              <i class="fas fa-users"></i>
            </div>
          </div>
          <div class="stat-card-footer">
            <i class="fas fa-arrow-up me-1"></i> Registered accounts
          </div>
        </div>
      </div>

      <div class="col-sm-6 col-xl-3">
        <div class="stat-card">
          <div class="stat-card-body">
            <div class="stat-card-info">
              <span class="stat-card-label">{{ t('admin.totalCVs') }}</span>
              <h3 class="stat-card-value">{{ stats.total_cvs }}</h3>
            </div>
            <div class="stat-card-icon stat-card-icon--green">
              <i class="fas fa-file-alt"></i>
            </div>
          </div>
          <div class="stat-card-footer">
            <i class="fas fa-file me-1"></i> Total CVs created
          </div>
        </div>
      </div>

      <div class="col-sm-6 col-xl-3">
        <div class="stat-card">
          <div class="stat-card-body">
            <div class="stat-card-info">
              <span class="stat-card-label">{{ t('admin.activeToday') }}</span>
              <h3 class="stat-card-value">{{ stats.active_today }}</h3>
            </div>
            <div class="stat-card-icon stat-card-icon--purple">
              <i class="fas fa-bolt"></i>
            </div>
          </div>
          <div class="stat-card-footer">
            <i class="fas fa-clock me-1"></i> Active today
          </div>
        </div>
      </div>

      <div class="col-sm-6 col-xl-3">
        <div class="stat-card">
          <div class="stat-card-body">
            <div class="stat-card-info">
              <span class="stat-card-label">{{ t('admin.pendingReview') }}</span>
              <h3 class="stat-card-value">{{ stats.pending_cvs }}</h3>
            </div>
            <div class="stat-card-icon stat-card-icon--gold">
              <i class="fas fa-clock"></i>
            </div>
          </div>
          <div class="stat-card-footer">
            <i class="fas fa-hourglass-half me-1"></i> Awaiting review
          </div>
        </div>
      </div>

      <div class="col-sm-6 col-xl-3" v-if="stats.guest_cvs > 0">
        <div class="stat-card">
          <div class="stat-card-body">
            <div class="stat-card-info">
              <span class="stat-card-label">{{ t('cv.guestCVs') }}</span>
              <h3 class="stat-card-value">{{ stats.guest_cvs }}</h3>
            </div>
            <div class="stat-card-icon stat-card-icon--gray">
              <i class="fas fa-user-secret"></i>
            </div>
          </div>
          <div class="stat-card-footer">
            <i class="fas fa-user-clock me-1"></i> Guest submissions
          </div>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="row g-4">
      <div class="col-lg-6">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-chart-pie me-2"></i>{{ t('admin.byStatus') }}
            </div>
          </div>
          <div class="chart-card-body">
            <canvas ref="statusChart"></canvas>
          </div>
        </div>
      </div>
      <div class="col-lg-6">
        <div class="chart-card">
          <div class="chart-card-header">
            <div class="chart-card-title">
              <i class="fas fa-chart-line me-2"></i>{{ t('admin.monthlyTrends') }}
            </div>
          </div>
          <div class="chart-card-body">
            <canvas ref="trendChart"></canvas>
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
const stats = ref<any>({ total_users: 0, total_cvs: 0, active_today: 0, pending_cvs: 0, guest_cvs: 0 })
const statusChart = ref<HTMLCanvasElement>()
const trendChart = ref<HTMLCanvasElement>()

onMounted(async () => {
  try {
    const [dashRes, statsRes] = await Promise.all([
      adminAPI.getDashboardStats(),
      adminAPI.getStats()
    ])
    stats.value = dashRes.data

    // Status chart
    if (statusChart.value) {
      const statusCounts = dashRes.data.status_counts || []
      new Chart(statusChart.value, {
        type: 'doughnut',
        data: {
          labels: statusCounts.map((s: any) => s.Status),
          datasets: [{ data: statusCounts.map((s: any) => s.Count), backgroundColor: ['#6c757d', '#ffc107', '#198754', '#dc3545'] }]
        }
      })
    }

    // Trend chart
    if (trendChart.value) {
      const trends = statsRes.data.monthly_trends || []
      new Chart(trendChart.value, {
        type: 'line',
        data: {
          labels: trends.map((t: any) => t.month),
          datasets: [{ label: 'CVs', data: trends.map((t: any) => t.count), borderColor: '#0d6efd', fill: false }]
        }
      })
    }
  } catch {}
})
</script>

<style scoped>
.admin-dashboard {
  max-width: 1400px;
}

/* ── Page Header ── */
.page-header {
  margin-bottom: 28px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #222222;
  margin: 0;
}

.page-subtitle {
  color: #6a6a6a;
  font-size: 14px;
  margin: 4px 0 0;
}

/* ── Stat Cards ── */
.stat-card {
  background: #ffffff;
  border-radius: 12px;
  border: none;
  box-shadow: rgba(0,0,0,0.02) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 6px, rgba(0,0,0,0.1) 0px 4px 8px;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-card-body {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px;
}

.stat-card-info {
  flex: 1;
}

.stat-card-label {
  font-size: 13px;
  font-weight: 600;
  color: #6a6a6a;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-card-value {
  font-size: 32px;
  font-weight: 700;
  color: #222222;
  margin: 6px 0 0;
  line-height: 1;
}

.stat-card-icon {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  flex-shrink: 0;
}

.stat-card-icon--blue {
  background: rgba(26, 82, 118, 0.1);
  color: #1a5276;
}

.stat-card-icon--green {
  background: rgba(46, 125, 50, 0.1);
  color: #2e7d32;
}

.stat-card-icon--purple {
  background: rgba(102, 16, 242, 0.1);
  color: #6610f2;
}

.stat-card-icon--gold {
  background: rgba(192, 152, 43, 0.1);
  color: #c0982b;
}

.stat-card-icon--gray {
  background: rgba(108, 117, 125, 0.1);
  color: #6c757d;
}

.stat-card-footer {
  padding: 10px 24px;
  font-size: 12px;
  color: #6a6a6a;
  border-top: 1px solid #f0f0f0;
  font-weight: 500;
}

/* ── Chart Cards ── */
.chart-card {
  background: #ffffff;
  border-radius: 12px;
  border: none;
  box-shadow: rgba(0,0,0,0.02) 0px 0px 0px 1px, rgba(0,0,0,0.04) 0px 2px 6px, rgba(0,0,0,0.1) 0px 4px 8px;
  overflow: hidden;
}

.chart-card-header {
  padding: 18px 24px;
  border-bottom: 1px solid #f0f0f0;
}

.chart-card-title {
  font-size: 15px;
  font-weight: 600;
  color: #222222;
  margin: 0;
}

.chart-card-body {
  padding: 24px;
}
</style>
