<script setup lang="ts">
import ErrorBoundary from './components/organisms/ErrorBoundary.vue';
import FeedbackDialog from './components/organisms/FeedbackDialog.vue';
</script>

<template>
  <ErrorBoundary>
    <router-view />
  </ErrorBoundary>
  <FeedbackDialog />
</template>
