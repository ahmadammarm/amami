<script setup lang="ts">
import AppSidebar from '../organisms/AppSidebar.vue';
</script>

<template>
  <div class="flex h-screen bg-white font-sans">
    <AppSidebar />
    <main class="flex-1 overflow-y-auto p-8 bg-gray-50/50">
      <router-view />
    </main>
  </div>
</template>
