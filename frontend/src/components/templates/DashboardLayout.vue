<script setup lang="ts">
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { Menu } from 'lucide-vue-next';
import AppSidebar from '../organisms/AppSidebar.vue';
import BrandMark from '../molecules/BrandMark.vue';

const route = useRoute();
const isSidebarOpen = ref(false);

const closeSidebar = () => {
  isSidebarOpen.value = false;
};

// The mobile drawer never stays open across navigation.
watch(() => route.fullPath, closeSidebar);
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-kapur">
    <!-- Mobile backdrop -->
    <div
      v-if="isSidebarOpen"
      class="fixed inset-0 z-40 bg-jati/40 lg:hidden"
      aria-hidden="true"
      @click="closeSidebar"
    />

    <!-- Sidebar: static on desktop, drawer on mobile -->
    <div
      class="fixed inset-y-0 left-0 z-50 transition-transform duration-300 ease-out lg:static lg:translate-x-0"
      :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full'"
    >
      <AppSidebar @close="closeSidebar" />
    </div>

    <div class="flex min-w-0 flex-1 flex-col">
      <!-- Mobile top bar -->
      <header class="flex h-14 items-center gap-3 border-b border-batu bg-white px-4 lg:hidden">
        <button
          type="button"
          class="rounded-md p-2 text-jati hover:bg-kapur"
          aria-label="Buka menu"
          @click="isSidebarOpen = true"
        >
          <Menu class="size-5" />
        </button>
        <BrandMark />
      </header>

      <main class="flex-1 overflow-y-auto">
        <div class="mx-auto w-full max-w-7xl px-4 py-6 sm:px-6 lg:px-10 lg:py-10">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>
