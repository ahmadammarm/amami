import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { Permission } from '../lib/permissions';

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean;
    guest?: boolean;
    /** Permission code the user's role must hold to open the route. */
    permission?: Permission;
  }
}

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
    name: 'Landing',
    component: () => import('../components/pages/LandingPage.vue'),
    meta: { guest: true } // Assuming guests see this, or we just let anyone see it
  },
  {
    path: '/dashboard',
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
        meta: { permission: Permission.LedgerRead },
      },
      {
        path: 'zakat',
        name: 'Zakat',
        component: () => import('../components/pages/ZakatPage.vue'),
        meta: { permission: Permission.ZakatRead },
      },
      {
        path: 'jamaah',
        name: 'Jamaah',
        component: () => import('../components/pages/JamaahPage.vue'),
        meta: { permission: Permission.JamaahManage },
      },
      {
        path: 'qurban',
        name: 'Qurban',
        component: () => import('../components/pages/QurbanPage.vue'),
        meta: { permission: Permission.QurbanRead },
      },
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('../components/pages/InventoryPage.vue'),
        meta: { permission: Permission.InventoryRead },
      },
      {
        path: 'agenda',
        name: 'Agenda',
        component: () => import('../components/pages/AgendaPage.vue'),
        meta: { permission: Permission.AgendaRead },
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../components/pages/SettingsPage.vue'),
        meta: { permission: Permission.SettingsManage },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../components/pages/UsersPage.vue'),
        meta: { permission: Permission.UserManage },
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
  const { permission } = to.meta;
  if (permission && authStore.isAuthenticated && !authStore.can(permission)) {
    return next({ name: 'Forbidden' });
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login' });
  } else if (to.meta.guest && authStore.isAuthenticated && (to.name === 'Login' || to.name === 'Landing')) {
    next({ name: 'Dashboard' });
  } else {
    next();
  }
});

export default router;
