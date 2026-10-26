<script setup lang="ts">
import { computed, type Component } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
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

interface MenuItem {
  name: string;
  icon: Component;
  path: string;
  /** Omitted for pages every signed-in user may open. */
  permission?: Permission;
}

const authStore = useAuthStore();
const route = useRoute();

const emit = defineEmits(['close']);

const menuItems: MenuItem[] = [
  { name: 'Dashboard', icon: LayoutDashboard, path: '/dashboard' },
  { name: 'Finance', icon: Wallet, path: '/dashboard/finance', permission: Permission.LedgerRead },
  { name: 'Zakat', icon: HandCoins, path: '/dashboard/zakat', permission: Permission.ZakatRead },
  { name: 'Jamaah', icon: Users, path: '/dashboard/jamaah', permission: Permission.JamaahManage },
  { name: 'Qurban', icon: Beef, path: '/dashboard/qurban', permission: Permission.QurbanRead },
  { name: 'Inventory', icon: PackageSearch, path: '/dashboard/inventory', permission: Permission.InventoryRead },
  { name: 'Agenda', icon: CalendarRange, path: '/dashboard/agenda', permission: Permission.AgendaRead },
  { name: 'Settings', icon: Settings, path: '/dashboard/settings', permission: Permission.SettingsManage },
  { name: 'Users', icon: Users, path: '/dashboard/users', permission: Permission.UserManage },
];

const filteredMenu = computed(() => {
  if (!authStore.user) return [];
  return menuItems.filter(item => !item.permission || authStore.can(item.permission));
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
