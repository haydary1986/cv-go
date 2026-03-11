<template>
  <div class="d-flex">
    <!-- Mobile overlay -->
    <div
      v-if="showSidebar"
      class="sidebar-overlay d-lg-none"
      @click="showSidebar = false"
    ></div>

    <!-- Sidebar -->
    <div class="admin-sidebar bg-dark text-white p-3" :class="{ 'show': showSidebar }">
      <h5 class="mb-4"><i class="fas fa-cog me-2"></i>{{ t('admin.dashboard') }}</h5>
      <nav class="nav flex-column">
        <router-link to="/admin" class="nav-link text-white" exact-active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-tachometer-alt me-2"></i>{{ t('app.dashboard') }}
        </router-link>
        <router-link to="/admin/cvs" class="nav-link text-white" active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-file-alt me-2"></i>{{ t('admin.manageCVs') }}
        </router-link>
        <router-link to="/admin/users" class="nav-link text-white" active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-users me-2"></i>{{ t('admin.manageUsers') }}
        </router-link>
        <router-link to="/admin/faculties" class="nav-link text-white" active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-university me-2"></i>{{ t('admin.manageFaculties') }}
        </router-link>
        <router-link to="/admin/statistics" class="nav-link text-white" active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-chart-bar me-2"></i>{{ t('admin.statistics') }}
        </router-link>
        <router-link to="/admin/settings" class="nav-link text-white" active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-sliders-h me-2"></i>{{ t('admin.branding') }}
        </router-link>
        <router-link to="/admin/logs" class="nav-link text-white" active-class="active bg-primary rounded" @click="closeSidebarOnMobile">
          <i class="fas fa-history me-2"></i>{{ t('admin.activityLogs') }}
        </router-link>
      </nav>
    </div>

    <!-- Main Content -->
    <div class="admin-content flex-grow-1">
      <button class="btn btn-dark d-lg-none m-2 sidebar-toggle" @click="showSidebar = !showSidebar">
        <i class="fas fa-bars"></i>
      </button>
      <div class="p-4">
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
.admin-sidebar { width: 250px; min-height: calc(100vh - 60px); position: sticky; top: 60px; }
.admin-content { min-height: calc(100vh - 60px); }
.nav-link { padding: 8px 12px; margin: 2px 0; transition: all 0.2s; }
.nav-link:hover { background: rgba(255,255,255,0.1); border-radius: 6px; }

.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.sidebar-toggle {
  position: relative;
  z-index: 1;
}

@media (max-width: 991px) {
  .admin-sidebar {
    position: fixed;
    left: -250px;
    top: 60px;
    bottom: 0;
    z-index: 1000;
    transition: left 0.3s ease;
    overflow-y: auto;
  }
  .admin-sidebar.show { left: 0; }
}
</style>
