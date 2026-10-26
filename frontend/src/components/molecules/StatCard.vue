<script setup lang="ts">
import type { Component } from 'vue';

withDefaults(defineProps<{
  label: string;
  value: string | number;
  hint?: string;
  icon?: Component;
  /** `alert` highlights a value that needs attention (e.g. broken assets). */
  tone?: 'default' | 'alert';
}>(), {
  tone: 'default',
});
</script>

<template>
  <div class="rounded-xl border border-batu bg-white p-5 shadow-sm">
    <div class="flex items-center justify-between gap-3">
      <p class="text-xs font-semibold uppercase tracking-[0.12em] text-jati-muda">{{ label }}</p>
      <component :is="icon" v-if="icon" class="size-4" :class="tone === 'alert' ? 'text-bata' : 'text-sajadah'" />
    </div>
    <p class="tabular mt-3 text-2xl font-semibold text-jati" :class="{ 'text-bata': tone === 'alert' }">{{ value }}</p>
    <p v-if="hint" class="mt-1 text-xs text-jati-muda">{{ hint }}</p>
  </div>
</template>
