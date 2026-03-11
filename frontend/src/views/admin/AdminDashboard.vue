<template>
  <div>
    <h3 class="mb-4">{{ t('admin.dashboard') }}</h3>

    <!-- Stats Cards -->
    <div class="row g-3 mb-4">
      <div class="col-md-3">
        <div class="card bg-primary text-white">
          <div class="card-body text-center">
            <i class="fas fa-users fa-2x mb-2"></i>
            <h3>{{ stats.total_users }}</h3>
            <small>{{ t('admin.totalUsers') }}</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card bg-success text-white">
          <div class="card-body text-center">
            <i class="fas fa-file-alt fa-2x mb-2"></i>
            <h3>{{ stats.total_cvs }}</h3>
            <small>{{ t('admin.totalCVs') }}</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card bg-info text-white">
          <div class="card-body text-center">
            <i class="fas fa-bolt fa-2x mb-2"></i>
            <h3>{{ stats.active_today }}</h3>
            <small>{{ t('admin.activeToday') }}</small>
          </div>
        </div>
      </div>
      <div class="col-md-3">
        <div class="card bg-warning text-dark">
          <div class="card-body text-center">
            <i class="fas fa-clock fa-2x mb-2"></i>
            <h3>{{ stats.pending_cvs }}</h3>
            <small>{{ t('admin.pendingReview') }}</small>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="row g-3">
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">{{ t('admin.byStatus') }}</div>
          <div class="card-body">
            <canvas ref="statusChart"></canvas>
          </div>
        </div>
      </div>
      <div class="col-md-6">
        <div class="card">
          <div class="card-header">{{ t('admin.monthlyTrends') }}</div>
          <div class="card-body">
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
const stats = ref<any>({ total_users: 0, total_cvs: 0, active_today: 0, pending_cvs: 0 })
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
