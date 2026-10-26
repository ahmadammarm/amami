<script setup lang="ts">
import { Toaster } from './components/atoms/sonner';
import ErrorBoundary from './components/atoms/ErrorBoundary.vue';
</script>

<template>
  <Toaster position="top-right" rich-colors close-button />
  <ErrorBoundary>
    <router-view />
  </ErrorBoundary>
</template>
