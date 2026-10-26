<script setup lang="ts">
import type { Component } from 'vue';
import { ArchBadge } from '../atoms/arch';

defineProps<{
  title: string;
  description?: string;
  icon?: Component;
}>();
</script>

<template>
  <div class="flex flex-col items-center px-6 py-14 text-center">
    <ArchBadge tone="outline" :size="52">
      <component :is="icon" v-if="icon" class="size-5" />
    </ArchBadge>
    <p class="mt-4 font-semibold text-jati">{{ title }}</p>
    <p v-if="description" class="mt-1 max-w-sm text-sm text-jati-muda">{{ description }}</p>
    <div v-if="$slots.default" class="mt-5">
      <slot />
    </div>
  </div>
</template>
