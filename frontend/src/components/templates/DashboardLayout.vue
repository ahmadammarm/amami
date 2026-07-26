<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import AppSidebar from '../organisms/AppSidebar.vue';
import { Menu, X, LogOut, Settings as SettingsIcon, ChevronDown } from 'lucide-vue-next';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '../atoms/dropdown-menu';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '../atoms/alert-dialog';

const authStore = useAuthStore();
const router = useRouter();
const isSidebarOpen = ref(false);
const isLogoutDialogOpen = ref(false);

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

const closeSidebar = () => {
  isSidebarOpen.value = false;
};

const handleLogout = () => {
  authStore.logout();
  router.push('/login');
};

const userInitials = computed(() => {
  if (authStore.user?.full_name) {
    return authStore.user.full_name.charAt(0).toUpperCase();
  }
  return authStore.user?.username?.charAt(0).toUpperCase() || 'U';
});

const canViewSettings = computed(() => {
  return ['SUPER_ADMIN', 'SEKRETARIS'].includes(authStore.userRole || '');
});

</script>

<template>
  <div class="flex h-screen bg-gray-50/50 font-sans overflow-hidden">
    <!-- Mobile Backdrop -->
    <div 
      v-if="isSidebarOpen" 
      class="fixed inset-0 bg-black/40 z-40 lg:hidden transition-opacity duration-300"
      @click="closeSidebar"
    ></div>

    <!-- Sticky Sidebar -->
    <div 
      class="fixed inset-y-0 left-0 z-50 transform lg:sticky lg:translate-x-0 transition-transform duration-300 ease-in-out h-full"
      :class="[isSidebarOpen ? 'translate-x-0' : '-translate-x-full']"
    >
      <AppSidebar @close="closeSidebar" />
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden h-full">
      <!-- Top Header -->
      <header class="sticky top-0 z-30 flex items-center justify-between lg:justify-end px-4 py-3 bg-white border-b border-gray-100 shadow-sm min-h-[64px]">
        
        <!-- Mobile Logo & Menu Toggle (Left) -->
        <div class="lg:hidden flex items-center gap-3">
          <button 
            @click="toggleSidebar" 
            class="p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-700 rounded-lg transition-colors"
            aria-label="Toggle Menu"
          >
            <Menu v-if="!isSidebarOpen" class="w-5 h-5" />
            <X v-else class="w-5 h-5" />
          </button>
          
          <h1 class="text-xl font-bold text-primary">amami</h1>
        </div>

        <!-- User Profile Dropdown (Right) -->
        <div class="flex items-center gap-4">
          <DropdownMenu v-if="authStore.user">
            <DropdownMenuTrigger as-child>
              <button class="flex items-center gap-2 hover:bg-gray-50 p-1.5 pr-2.5 rounded-full transition-colors outline-none focus:ring-2 focus:ring-primary/20">
                <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold text-xs uppercase border border-primary/20">
                  {{ userInitials }}
                </div>
                <div class="hidden sm:block text-left">
                  <p class="text-xs font-bold text-gray-700 leading-none truncate max-w-[120px]">{{ authStore.user.username }}</p>
                </div>
                <ChevronDown class="w-3 h-3 text-gray-400 hidden sm:block" />
              </button>
            </DropdownMenuTrigger>
            
            <DropdownMenuContent align="end" class="w-56 rounded-xl p-2 shadow-xl border-gray-100">
              <DropdownMenuLabel class="font-normal px-2 py-2">
                <div class="flex flex-col space-y-1">
                  <p class="text-sm font-bold text-gray-900 leading-none">{{ authStore.user.full_name || authStore.user.username }}</p>
                  <p class="text-xs font-medium text-gray-500 leading-none">{{ authStore.user.email }}</p>
                  <div class="mt-2 inline-flex items-center w-fit px-2 py-0.5 rounded-md bg-gray-100 text-[10px] font-bold text-gray-600 uppercase tracking-wider">
                    {{ authStore.user.role_name.replace('_', ' ') }}
                  </div>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator class="bg-gray-100 my-1" />
              
              <DropdownMenuItem 
                v-if="canViewSettings"
                @click="router.push('/settings')" 
                class="rounded-lg font-bold text-xs py-2.5 cursor-pointer text-gray-600 hover:text-gray-900 hover:bg-gray-50"
              >
                <SettingsIcon class="w-4 h-4 mr-2" /> System Settings
              </DropdownMenuItem>
              
              <DropdownMenuItem 
                @click="isLogoutDialogOpen = true" 
                class="rounded-lg font-bold text-xs py-2.5 cursor-pointer text-red-600 focus:bg-red-50 focus:text-red-700 mt-1"
              >
                <LogOut class="w-4 h-4 mr-2" /> Sign out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>

      <!-- Scrollable Main Content -->
      <main class="flex-1 overflow-y-auto p-4 md:p-8 bg-gray-50/50">
        <router-view />
      </main>
    </div>

    <!-- Logout Confirmation Dialog -->
    <AlertDialog v-model:open="isLogoutDialogOpen">
      <AlertDialogContent class="rounded-2xl max-w-sm">
        <AlertDialogHeader>
          <AlertDialogTitle class="text-xl">Sign out of amami?</AlertDialogTitle>
          <AlertDialogDescription class="text-gray-500 font-medium">
            You will need to sign in again to access the mosque management dashboard.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter class="sm:justify-end gap-2 pt-2 border-t border-gray-50 mt-4">
          <AlertDialogCancel class="rounded-xl h-10 mt-0 font-bold border-gray-200">Cancel</AlertDialogCancel>
          <AlertDialogAction 
            @click="handleLogout"
            class="rounded-xl h-10 bg-red-500 hover:bg-red-600 focus:ring-red-500 font-bold border-none"
          >
            Sign out
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
