<script setup lang="ts">
import { Toaster } from 'vue-sonner';
import ErrorBoundary from './components/atoms/ErrorBoundary.vue';
</script>

<template>
  <Toaster position="top-right" rich-colors />
  <ErrorBoundary>
    <router-view />
  </ErrorBoundary>
</template>
