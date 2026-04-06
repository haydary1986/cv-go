import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
      meta: { title: 'الرئيسية' },
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { guest: true, title: 'تسجيل الدخول' },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { guest: true, title: 'إنشاء حساب' },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/DashboardView.vue'),
      meta: { requiresAuth: true, title: 'لوحة التحكم' },
    },
    {
      path: '/cv/create',
      name: 'cv-create',
      component: () => import('../views/CVFormView.vue'),
      meta: { requiresAuth: true, title: 'إنشاء سيرة ذاتية' },
    },
    {
      path: '/cv/guest',
      name: 'cv-guest',
      component: () => import('../views/GuestCVFormView.vue'),
      meta: { title: 'إنشاء سيرة ذاتية - ضيف' },
    },
    {
      path: '/cv/:id/edit',
      name: 'cv-edit',
      component: () => import('../views/CVFormView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/cv/:id',
      name: 'cv-view',
      component: () => import('../views/CVPreviewView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/shared/:token',
      name: 'shared-cv',
      component: () => import('../views/SharedCVView.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/notifications',
      name: 'notifications',
      component: () => import('../views/NotificationsView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('../views/admin/AdminDashboard.vue'),
        },
        {
          path: 'cvs',
          name: 'admin-cvs',
          component: () => import('../views/admin/AdminCVs.vue'),
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('../views/admin/AdminUsers.vue'),
        },
        {
          path: 'faculties',
          name: 'admin-faculties',
          component: () => import('../views/admin/AdminFaculties.vue'),
        },
        {
          path: 'settings',
          name: 'admin-settings',
          component: () => import('../views/admin/AdminSettings.vue'),
        },
        {
          path: 'statistics',
          name: 'admin-statistics',
          component: () => import('../views/admin/AdminStatistics.vue'),
        },
        {
          path: 'logs',
          name: 'admin-logs',
          component: () => import('../views/admin/AdminLogs.vue'),
        },
        {
          path: 'audit-trail',
          name: 'admin-audit-trail',
          component: () => import('../views/admin/AdminAuditTrail.vue'),
        },
      ],
    },
    {
      path: '/import/linkedin',
      name: 'linkedin-import',
      component: () => import('../views/LinkedInImport.vue'),
      meta: { requiresAuth: true, title: 'استيراد من LinkedIn' },
    },
    {
      path: '/auth/callback',
      name: 'auth-callback',
      component: () => import('../views/AuthCallback.vue'),
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFoundView.vue'),
    },
  ],
})

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  if (authStore.isAuthenticated && !authStore.user) {
    await authStore.fetchProfile()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: 'dashboard' })
  } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'dashboard' })
  } else {
    next()
  }
})

router.afterEach((to) => {
  const pageTitle = to.meta.title as string
  const baseName = 'جامعة التراث'
  document.title = pageTitle ? `${pageTitle} | ${baseName}` : baseName
})

export default router
