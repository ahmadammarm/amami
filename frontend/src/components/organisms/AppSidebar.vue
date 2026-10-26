<script setup lang="ts">
import { computed, type Component } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  LayoutDashboard,
  Wallet,
  HandCoins,
  Beef,
  Users,
  UserCog,
  Settings,
  X,
  PackageSearch,
  CalendarRange,
  LogOut,
  ChevronsUpDown,
} from 'lucide-vue-next';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import { roleLabel } from '../../lib/roles';
import { feedback } from '../../composables/useFeedback';
import BrandMark from '../molecules/BrandMark.vue';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../atoms/dropdown-menu';

interface MenuItem {
  name: string;
  icon: Component;
  path: string;
  /** Omitted for pages every signed-in user may open. */
  permission?: Permission;
}

interface MenuGroup {
  label: string;
  items: MenuItem[];
}

const authStore = useAuthStore();
const route = useRoute();
const router = useRouter();

const emit = defineEmits<{ (e: 'close'): void }>();

const menuGroups: MenuGroup[] = [
  {
    label: 'Operasional',
    items: [
      { name: 'Beranda', icon: LayoutDashboard, path: '/dashboard' },
      { name: 'Keuangan', icon: Wallet, path: '/dashboard/finance', permission: Permission.LedgerRead },
      { name: 'Zakat', icon: HandCoins, path: '/dashboard/zakat', permission: Permission.ZakatRead },
      { name: 'Jamaah', icon: Users, path: '/dashboard/jamaah', permission: Permission.JamaahManage },
      { name: 'Qurban', icon: Beef, path: '/dashboard/qurban', permission: Permission.QurbanRead },
      { name: 'Inventaris', icon: PackageSearch, path: '/dashboard/inventory', permission: Permission.InventoryRead },
      { name: 'Agenda', icon: CalendarRange, path: '/dashboard/agenda', permission: Permission.AgendaRead },
    ],
  },
  {
    label: 'Administrasi',
    items: [
      { name: 'Pengaturan', icon: Settings, path: '/dashboard/settings', permission: Permission.SettingsManage },
      { name: 'Pengguna', icon: UserCog, path: '/dashboard/users', permission: Permission.UserManage },
    ],
  },
];

const visibleGroups = computed(() => {
  if (!authStore.user) return [];
  return menuGroups
    .map(group => ({
      ...group,
      items: group.items.filter(item => !item.permission || authStore.can(item.permission)),
    }))
    .filter(group => group.items.length > 0);
});

const displayName = computed(() => authStore.user?.full_name || authStore.user?.username || '');
const initial = computed(() => displayName.value.charAt(0).toUpperCase() || 'A');
const userRoleLabel = computed(() => (authStore.user ? roleLabel(authStore.user.role_name) : ''));

const isActive = (path: string) => route.path === path;

const signOut = async () => {
  const confirmed = await feedback.confirm({
    title: 'Keluar dari amami?',
    message: 'Anda perlu masuk lagi untuk membuka dashboard.',
    confirmLabel: 'Keluar',
  });
  if (!confirmed) return;
  authStore.logout();
  router.push({ name: 'Login' });
};
</script>

<template>
  <aside class="flex h-full w-64 flex-col bg-sajadah-tua text-white">
    <div class="flex items-center justify-between px-5 pt-6 pb-5">
      <RouterLink :to="{ name: 'Dashboard' }" class="rounded-md" @click="emit('close')">
        <BrandMark variant="light" show-tagline />
      </RouterLink>
      <button
        type="button"
        class="rounded-md p-1.5 text-white/70 hover:bg-white/10 hover:text-white lg:hidden"
        aria-label="Tutup menu"
        @click="emit('close')"
      >
        <X class="size-5" />
      </button>
    </div>

    <nav class="flex-1 space-y-6 overflow-y-auto px-3 pb-6" aria-label="Menu utama">
      <div v-for="group in visibleGroups" :key="group.label">
        <p class="px-3 pb-2 text-[10px] font-semibold uppercase tracking-[0.2em] text-white/45">{{ group.label }}</p>
        <ul class="space-y-0.5">
          <li v-for="item in group.items" :key="item.path">
            <RouterLink
              :to="item.path"
              class="relative flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors"
              :class="isActive(item.path) ? 'bg-white/10 text-white' : 'text-white/70 hover:bg-white/5 hover:text-white'"
              :aria-current="isActive(item.path) ? 'page' : undefined"
              @click="emit('close')"
            >
              <span
                v-if="isActive(item.path)"
                class="absolute inset-y-2 left-0 w-[3px] rounded-full bg-kuningan"
                aria-hidden="true"
              />
              <component :is="item.icon" class="size-[18px] shrink-0" :class="isActive(item.path) ? 'text-kuningan' : ''" />
              {{ item.name }}
            </RouterLink>
          </li>
        </ul>
      </div>
    </nav>

    <div v-if="authStore.user" class="border-t border-white/10 p-3">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <button
            type="button"
            class="flex w-full items-center gap-3 rounded-lg px-2 py-2 text-left transition-colors hover:bg-white/10"
          >
            <span class="flex size-9 shrink-0 items-center justify-center rounded-full bg-kuningan/20 font-semibold text-kuningan">
              {{ initial }}
            </span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-sm font-semibold">{{ displayName }}</span>
              <span class="block truncate text-[11px] uppercase tracking-wider text-white/55">{{ userRoleLabel }}</span>
            </span>
            <ChevronsUpDown class="size-4 text-white/50" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent side="top" align="start" class="w-56 rounded-xl border-batu p-1">
          <div class="px-2 py-2">
            <p class="truncate text-sm font-semibold text-jati">{{ displayName }}</p>
            <p class="truncate text-xs text-jati-muda">{{ authStore.user.email }}</p>
          </div>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            v-if="authStore.can(Permission.SettingsManage)"
            class="cursor-pointer rounded-lg py-2"
            @click="router.push({ name: 'Settings' }); emit('close')"
          >
            <Settings class="mr-2 size-4" /> Pengaturan
          </DropdownMenuItem>
          <DropdownMenuItem class="cursor-pointer rounded-lg py-2 text-bata focus:bg-bata-muda focus:text-bata" @click="signOut">
            <LogOut class="mr-2 size-4" /> Keluar
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  </aside>
</template>
