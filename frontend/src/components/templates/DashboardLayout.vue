<script setup lang="ts">
import { ref } from 'vue';
import AppSidebar from '../organisms/AppSidebar.vue';
import { Menu, X } from 'lucide-vue-next';

const isSidebarOpen = ref(false);

const toggleSidebar = () => {
  isSidebarOpen.value = !isSidebarOpen.value;
};

const closeSidebar = () => {
  isSidebarOpen.value = false;
};
</script>

<template>
  <div class="flex h-screen bg-white font-sans overflow-hidden">
    <!-- Mobile Backdrop -->
    <div 
      v-if="isSidebarOpen" 
      class="fixed inset-0 bg-black/40 z-40 lg:hidden transition-opacity duration-300"
      @click="closeSidebar"
    ></div>

    <!-- Sidebar -->
    <div 
      class="fixed inset-y-0 left-0 z-50 transform lg:relative lg:translate-x-0 transition-transform duration-300 ease-in-out"
      :class="[isSidebarOpen ? 'translate-x-0' : '-translate-x-full']"
    >
      <AppSidebar @close="closeSidebar" />
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 overflow-hidden">
      <!-- Top Mobile Header -->
      <header class="lg:hidden flex items-center justify-between p-4 border-b border-gray-100 bg-white">
        <h1 class="text-xl font-bold text-primary">amami</h1>
        <button 
          @click="toggleSidebar" 
          class="p-2 text-gray-500 hover:bg-gray-50 rounded-lg"
          aria-label="Toggle Menu"
        >
          <Menu v-if="!isSidebarOpen" class="w-6 h-6" />
          <X v-else class="w-6 h-6" />
        </button>
      </header>

      <main class="flex-1 overflow-y-auto p-4 md:p-8 bg-gray-50/50">
        <router-view />
      </main>
    </div>
  </div>
</template>
