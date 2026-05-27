<script setup lang="ts">
import { computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { 
  LayoutDashboard, 
  Wallet, 
  HandCoins, 
  Beef, 
  Users, 
  Settings, 
  LogOut 
} from 'lucide-vue-next';

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute();

const menuItems = [
  { name: 'Dashboard', icon: LayoutDashboard, path: '/', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] },
  { name: 'Finance', icon: Wallet, path: '/finance', roles: ['SUPER_ADMIN', 'BENDAHARA'] },
  { name: 'Zakat', icon: HandCoins, path: '/zakat', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR'] },
  { name: 'Qurban', icon: Beef, path: '/qurban', roles: ['SUPER_ADMIN', 'TAKMIR'] },
  { name: 'Users', icon: Users, path: '/users', roles: ['SUPER_ADMIN'] },
  { name: 'Settings', icon: Settings, path: '/settings', roles: ['SUPER_ADMIN', 'SEKRETARIS'] },
];

const filteredMenu = computed(() => {
  const role = authStore.userRole;
  if (!role) return [];
  return menuItems.filter(item => item.roles.includes(role));
});

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};
</script>

<template>
  <aside class="w-64 bg-white border-r border-gray-100 flex flex-col h-full shadow-md">
    <div class="p-6 border-b border-gray-50 mb-4">
      <h1 class="text-2xl font-bold text-primary flex items-center gap-2">
        amami
      </h1>
      <p class="text-xs text-gray-400 mt-1 uppercase tracking-wider font-semibold">Mosque Intelligence</p>
    </div>

    <nav class="flex-1 px-4 space-y-1">
      <router-link
        v-for="item in filteredMenu"
        :key="item.name"
        :to="item.path"
        class="flex items-center gap-3 px-4 py-3 text-sm font-medium transition-all duration-200 rounded-lg group"
        :class="[
          route.path === item.path
            ? 'bg-primary/10 text-primary'
            : 'text-gray-500 hover:bg-gray-50 hover:text-primary'
        ]"
      >
        <component :is="item.icon" class="w-5 h-5" />
        {{ item.name }}
      </router-link>
    </nav>

    <div class="p-4 border-t border-gray-50">
      <div v-if="authStore.user" class="mb-4 px-4 py-3 bg-gray-50 rounded-lg flex items-center gap-3">
        <div class="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-xs uppercase">
          {{ authStore.user.username.charAt(0) }}
        </div>
        <div class="overflow-hidden">
          <p class="text-sm font-semibold text-gray-700 truncate">{{ authStore.user.username }}</p>
          <p class="text-xs text-gray-400 font-medium uppercase">{{ authStore.user.role_name }}</p>
        </div>
      </div>
      
      <button
        @click="handleLogout"
        class="w-full flex items-center gap-3 px-4 py-3 text-sm font-medium text-red-500 hover:bg-red-50 rounded-lg transition-colors duration-200"
      >
        <LogOut class="w-5 h-5" />
        Logout
      </button>
    </div>
  </aside>
</template>
