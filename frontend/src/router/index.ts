import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../components/pages/LoginPage.vue'),
    meta: { guest: true },
  },
  {
    path: '/force-change-password',
    name: 'ForceChangePassword',
    component: () => import('../components/pages/ForceChangePasswordPage.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/',
    component: () => import('../components/templates/DashboardLayout.vue'),
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../components/pages/DashboardPage.vue'),
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../components/pages/SettingsPage.vue'),
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../components/pages/UsersPage.vue'),
        meta: { requiresRole: 'SUPER_ADMIN' }
      },
      // Other protected routes will go here
    ],
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore();
  
  if (authStore.token && !authStore.user) {
    try {
      await authStore.fetchUser();
    } catch (e) {
      return next({ name: 'Login' });
    }
  }

  // 1. If authenticated but needs password change, force them to that page
  if (authStore.isAuthenticated && authStore.requiresPasswordChange && to.name !== 'ForceChangePassword') {
    return next({ name: 'ForceChangePassword' });
  }

  // 2. If trying to access ForceChangePassword but doesn't need it, redirect to home
  if (authStore.isAuthenticated && !authStore.requiresPasswordChange && to.name === 'ForceChangePassword') {
    return next({ name: 'Dashboard' });
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login' });
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: 'Dashboard' });
  } else {
    next();
  }
});

export default router;
