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
    path: '/forbidden',
    name: 'Forbidden',
    component: () => import('../components/pages/ForbiddenPage.vue'),
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
        path: 'finance',
        name: 'Finance',
        component: () => import('../components/pages/FinancePage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS'] }
      },
      {
        path: 'zakat',
        name: 'Zakat',
        component: () => import('../components/pages/ZakatPage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] }
      },
      {
        path: 'jamaah',
        name: 'Jamaah',
        component: () => import('../components/pages/JamaahPage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS'] }
      },
      {
        path: 'qurban',
        name: 'Qurban',
        component: () => import('../components/pages/QurbanPage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] }
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('../components/pages/InventoryPage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] }
      },
      {
        path: 'agenda',
        name: 'Agenda',
        component: () => import('../components/pages/AgendaPage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../components/pages/SettingsPage.vue'),
        meta: { requiresRole: ['SUPER_ADMIN', 'SEKRETARIS'] }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../components/pages/UsersPage.vue'),
        meta: { requiresRole: 'SUPER_ADMIN' }
      },
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

  // 3. RBAC Check
  if (to.meta.requiresRole && authStore.isAuthenticated) {
    const requiredRoles = Array.isArray(to.meta.requiresRole) 
      ? to.meta.requiresRole 
      : [to.meta.requiresRole];
      
    if (!requiredRoles.includes(authStore.userRole)) {
      return next({ name: 'Forbidden' });
    }
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
