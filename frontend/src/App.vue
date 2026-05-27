<script setup lang="ts">
import { Toaster } from 'vue-sonner';
</script>

<template>
  <Toaster position="top-right" rich-colors />
  <router-view />
</template>
