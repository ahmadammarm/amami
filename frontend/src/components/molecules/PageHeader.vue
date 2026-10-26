<script setup lang="ts">
defineProps<{
  title: string;
  description?: string;
}>();
</script>

<template>
  <header class="flex flex-col gap-4 border-b border-batu pb-6 sm:flex-row sm:items-end sm:justify-between">
    <div class="min-w-0">
      <h1 class="font-display text-3xl font-semibold tracking-wide text-jati">{{ title }}</h1>
      <p v-if="description" class="mt-1.5 max-w-2xl text-sm text-jati-muda">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
      <slot name="actions" />
    </div>
  </header>
</template>
