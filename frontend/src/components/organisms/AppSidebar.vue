<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { 
  LayoutDashboard, 
  Wallet, 
  HandCoins, 
  Beef, 
  Users, 
  Settings, 
  X,
  PackageSearch,
  CalendarRange
} from 'lucide-vue-next';

const authStore = useAuthStore();
const route = useRoute();

const emit = defineEmits(['close']);

const menuItems = [
  { name: 'Dashboard', icon: LayoutDashboard, path: '/', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] },
  { name: 'Finance', icon: Wallet, path: '/finance', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS'] },
  { name: 'Zakat', icon: HandCoins, path: '/zakat', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] },
  { name: 'Jamaah', icon: Users, path: '/jamaah', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS'] },
  { name: 'Qurban', icon: Beef, path: '/qurban', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] },
  { name: 'Inventory', icon: PackageSearch, path: '/inventory', roles: ['SUPER_ADMIN', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] },
  { name: 'Logistics', icon: CalendarRange, path: '/logistics', roles: ['SUPER_ADMIN', 'BENDAHARA', 'TAKMIR', 'SEKRETARIS', 'JAMAAH'] },
  { name: 'Settings', icon: Settings, path: '/settings', roles: ['SUPER_ADMIN', 'SEKRETARIS'] },
  { name: 'Users', icon: Users, path: '/users', roles: ['SUPER_ADMIN'] },
];

const filteredMenu = computed(() => {
  const role = authStore.userRole;
  if (!role) return [];
  return menuItems.filter(item => item.roles.includes(role));
});
</script>

<template>
  <aside class="w-64 bg-white border-r border-gray-100 flex flex-col h-full shadow-sm relative z-20">
    <div class="p-6 border-b border-gray-50 mb-4 flex items-center justify-between sticky top-0 bg-white z-10">
      <div>
        <h1 class="text-2xl font-bold text-primary flex items-center gap-2">
          amami
        </h1>
        <p class="text-xs text-gray-400 mt-1 uppercase tracking-wider font-semibold">Mosque Intelligence</p>
      </div>
      <button @click="emit('close')" class="lg:hidden p-2 text-gray-400 hover:text-gray-600">
        <X class="w-5 h-5" />
      </button>
    </div>

    <nav class="flex-1 px-4 space-y-1 overflow-y-auto pb-4">
      <router-link
        v-for="item in filteredMenu"
        :key="item.name"
        :to="item.path"
        @click="emit('close')"
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
  </aside>
</template>
