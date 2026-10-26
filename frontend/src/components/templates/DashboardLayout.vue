<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { onKeyStroke, useMediaQuery } from '@vueuse/core';
import { Menu } from 'lucide-vue-next';
import AppSidebar from '../organisms/AppSidebar.vue';
import BrandMark from '../molecules/BrandMark.vue';

const route = useRoute();
const isSidebarOpen = ref(false);
const isDesktop = useMediaQuery('(min-width: 1024px)');

/** On mobile the closed drawer is off-screen; inert keeps its links out of the tab order. */
const isSidebarHidden = computed(() => !isDesktop.value && !isSidebarOpen.value);

const closeSidebar = () => {
  isSidebarOpen.value = false;
};

// The mobile drawer never stays open across navigation or a resize to desktop.
watch(() => route.fullPath, closeSidebar);
watch(isDesktop, closeSidebar);

onKeyStroke('Escape', () => {
  if (isSidebarOpen.value) closeSidebar();
});
</script>

<template>
  <!-- Fixed to the viewport: the document itself never scrolls, only <main> does,
       so mobile browsers cannot drag the page past its content. -->
  <div class="fixed inset-0 flex overflow-hidden bg-kapur">
    <!-- Mobile backdrop -->
    <div
      v-if="isSidebarOpen"
      class="fixed inset-0 z-40 bg-jati/40 lg:hidden"
      aria-hidden="true"
      @click="closeSidebar"
    />

    <!-- Sidebar: static on desktop, drawer on mobile -->
    <div
      class="fixed inset-y-0 left-0 z-50 max-w-[85vw] transition-transform duration-300 ease-out lg:static lg:max-w-none lg:translate-x-0"
      :class="isSidebarOpen ? 'translate-x-0 shadow-lg' : '-translate-x-full'"
      :inert="isSidebarHidden"
    >
      <AppSidebar @close="closeSidebar" />
    </div>

    <div class="flex min-w-0 flex-1 flex-col">
      <!-- Mobile top bar -->
      <header class="flex h-14 shrink-0 items-center gap-2 border-b border-batu bg-white px-2 sm:px-4 lg:hidden">
        <button
          type="button"
          class="tombol-ikon text-jati"
          aria-label="Buka menu"
          :aria-expanded="isSidebarOpen"
          @click="isSidebarOpen = true"
        >
          <Menu class="size-5" />
        </button>
        <BrandMark />
      </header>

      <main class="scrollbar-tipis min-h-0 flex-1 overflow-x-hidden overflow-y-auto overscroll-contain">
        <div class="mx-auto w-full max-w-7xl px-4 py-6 sm:px-6 lg:px-10 lg:py-10 2xl:max-w-[88rem]">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>
