<script setup lang="ts">
import type { BadgeTone } from '../../types/ui';

withDefaults(defineProps<{ tone?: BadgeTone }>(), { tone: 'netral' });

const toneClass: Record<BadgeTone, string> = {
  sajadah: 'bg-sajadah-muda text-sajadah-tua',
  kuningan: 'bg-kuningan-muda text-kuningan-tua',
  bata: 'bg-bata-muda text-bata',
  netral: 'bg-batu/70 text-jati-muda',
};
</script>

<template>
  <span class="inline-flex items-center gap-1 whitespace-nowrap rounded-full px-2.5 py-0.5 text-xs font-semibold" :class="toneClass[tone]">
    <slot />
  </span>
</template>
