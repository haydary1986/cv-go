<template>
  <div class="admin-wrapper">
    <!-- Mobile overlay -->
    <div
      v-if="showSidebar"
      class="sidebar-overlay d-lg-none"
      @click="showSidebar = false"
    ></div>

    <!-- Sidebar -->
    <aside class="admin-sidebar" :class="{ 'show': showSidebar }">
      <!-- Logo / Brand -->
      <div class="sidebar-brand">
        <div class="brand-icon">
          <i class="fas fa-university"></i>
        </div>
        <div class="brand-text">
          <span class="brand-title">{{ t('admin.dashboard') }}</span>
          <span class="brand-subtitle">Admin Panel</span>
        </div>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <div class="nav-section-label">{{ t('app.menu') || 'MENU' }}</div>

        <router-link to="/admin" class="sidebar-link" exact-active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-tachometer-alt"></i></span>
          <span class="sidebar-link-text">{{ t('app.dashboard') }}</span>
        </router-link>

        <router-link to="/admin/cvs" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-file-alt"></i></span>
          <span class="sidebar-link-text">{{ t('admin.manageCVs') }}</span>
        </router-link>

        <router-link to="/admin/users" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-users"></i></span>
          <span class="sidebar-link-text">{{ t('admin.manageUsers') }}</span>
        </router-link>

        <router-link to="/admin/faculties" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-university"></i></span>
          <span class="sidebar-link-text">{{ t('admin.manageFaculties') }}</span>
        </router-link>

        <div class="nav-section-label mt-3">{{ t('admin.analytics') || 'ANALYTICS' }}</div>

        <router-link to="/admin/statistics" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-chart-bar"></i></span>
          <span class="sidebar-link-text">{{ t('admin.statistics') }}</span>
        </router-link>

        <router-link to="/admin/logs" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-history"></i></span>
          <span class="sidebar-link-text">{{ t('admin.activityLogs') }}</span>
        </router-link>

        <router-link to="/admin/audit" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-shield-alt"></i></span>
          <span class="sidebar-link-text">{{ t('admin.auditTrail') || 'Audit Trail' }}</span>
        </router-link>

        <div class="nav-section-label mt-3">{{ t('admin.system') || 'SYSTEM' }}</div>

        <router-link to="/admin/settings" class="sidebar-link" active-class="active" @click="closeSidebarOnMobile">
          <span class="sidebar-link-icon"><i class="fas fa-sliders-h"></i></span>
          <span class="sidebar-link-text">{{ t('admin.branding') }}</span>
        </router-link>
      </nav>

      <!-- Sidebar Footer -->
      <div class="sidebar-footer">
        <div class="sidebar-footer-text">
          <i class="fas fa-lock me-1"></i> Admin Access
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="admin-content">
      <div class="admin-topbar d-lg-none">
        <button class="sidebar-toggle-btn" @click="showSidebar = !showSidebar">
          <i class="fas fa-bars"></i>
        </button>
        <span class="topbar-title">{{ t('admin.dashboard') }}</span>
      </div>
      <div class="admin-content-inner">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const showSidebar = ref(false)

function closeSidebarOnMobile() {
  if (window.innerWidth < 992) {
    showSidebar.value = false
  }
}
</script>

<style scoped>
.admin-wrapper {
  display: flex;
  min-height: calc(100vh - 60px);
}

/* ── Sidebar ── */
.admin-sidebar {
  width: 260px;
  min-width: 260px;
  flex-shrink: 0;
  min-height: calc(100vh - 60px);
  position: sticky;
  top: 60px;
  background: linear-gradient(180deg, #0a1628 0%, #0f1f3d 50%, #132744 100%);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  z-index: 1000;
}

/* Brand */
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 20px 20px;
  border-bottom: 1px solid rgba(192, 152, 43, 0.15);
  margin-bottom: 8px;
}

.brand-icon {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: linear-gradient(135deg, #c0982b, #d4a93a);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: #0a1628;
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-title {
  color: #ffffff;
  font-weight: 700;
  font-size: 15px;
  line-height: 1.2;
}

.brand-subtitle {
  color: rgba(255, 255, 255, 0.4);
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  padding: 8px 12px;
}

.nav-section-label {
  color: rgba(192, 152, 43, 0.6);
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  padding: 12px 12px 6px;
}

.sidebar-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  margin: 2px 0;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.65);
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
  font-size: 14px;
  font-weight: 500;
}

.sidebar-link:hover {
  color: #ffffff;
  background: rgba(255, 255, 255, 0.06);
}

.sidebar-link.active {
  color: #ffffff;
  background: rgba(192, 152, 43, 0.12);
}

.sidebar-link.active::before {
  content: '';
  position: absolute;
  inset-inline-start: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: #c0982b;
  border-radius: 0 3px 3px 0;
}

[dir="rtl"] .sidebar-link.active::before {
  border-radius: 3px 0 0 3px;
}

.sidebar-link-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.04);
  transition: all 0.2s ease;
}

.sidebar-link.active .sidebar-link-icon {
  background: rgba(192, 152, 43, 0.2);
  color: #c0982b;
}

.sidebar-link-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Footer */
.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.sidebar-footer-text {
  color: rgba(255, 255, 255, 0.3);
  font-size: 11px;
  text-align: center;
}

/* ── Main Content ── */
.admin-content {
  flex-grow: 1;
  min-height: calc(100vh - 60px);
  overflow-x: hidden;
  min-width: 0;
  background: #f8fafc;
}

.admin-content-inner {
  padding: 28px 32px;
}

/* ── Mobile Top Bar ── */
.admin-topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #ffffff;
  border-bottom: 1px solid #e9ecef;
  position: sticky;
  top: 60px;
  z-index: 10;
}

.sidebar-toggle-btn {
  width: 36px;
  height: 36px;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  background: #ffffff;
  color: #1a5276;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.sidebar-toggle-btn:hover {
  background: #1a5276;
  color: #ffffff;
  border-color: #1a5276;
}

.topbar-title {
  font-weight: 600;
  color: #1a5276;
  font-size: 16px;
}

/* ── Overlay ── */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(10, 22, 40, 0.6);
  backdrop-filter: blur(2px);
  z-index: 999;
}

/* ── Responsive ── */
@media (max-width: 991px) {
  .admin-sidebar {
    position: fixed;
    inset-inline-start: -260px;
    top: 60px;
    bottom: 0;
    transition: inset-inline-start 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: none;
  }

  .admin-sidebar.show {
    inset-inline-start: 0;
    box-shadow: 4px 0 24px rgba(0, 0, 0, 0.3);
  }

  .admin-content-inner {
    padding: 20px 16px;
  }
}
</style>
