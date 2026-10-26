<script setup lang="ts">
import { ARCH_PATH } from '../atoms/arch';

withDefaults(defineProps<{
  /** `light` is for dark backgrounds such as the sidebar. */
  variant?: 'dark' | 'light';
  showTagline?: boolean;
}>(), {
  variant: 'dark',
  showTagline: false,
});
</script>

<template>
  <span class="inline-flex items-center gap-2.5">
    <svg viewBox="0 0 40 50" class="h-9 w-[29px] shrink-0" aria-hidden="true">
      <path :d="ARCH_PATH" :class="variant === 'light' ? 'fill-white/10 stroke-kuningan' : 'fill-sajadah stroke-sajadah'" stroke-width="1.5" />
      <text
        x="20"
        y="41"
        text-anchor="middle"
        class="font-display fill-white"
        font-size="22"
        font-weight="600"
      >a</text>
    </svg>
    <span class="flex flex-col leading-none">
      <span class="font-display text-2xl font-semibold tracking-wide" :class="variant === 'light' ? 'text-white' : 'text-jati'">amami</span>
      <span
        v-if="showTagline"
        class="mt-1 text-[10px] font-semibold uppercase tracking-[0.18em]"
        :class="variant === 'light' ? 'text-white/60' : 'text-jati-muda'"
      >Manajemen Masjid</span>
    </span>
  </span>
</template>
