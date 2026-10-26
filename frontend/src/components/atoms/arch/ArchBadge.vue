<script setup lang="ts">
import { computed } from 'vue';
import { cn } from '@/lib/utils';
import { ARCH_PATH, type ArchTone } from './arch';

const props = withDefaults(defineProps<{
  tone?: ArchTone;
  /** Width in pixels; the height follows the arch proportions (4:5). */
  size?: number;
  class?: string;
}>(), {
  tone: 'sajadah',
  size: 48,
});

const toneClass: Record<ArchTone, { shape: string; content: string }> = {
  sajadah: { shape: 'fill-sajadah-muda stroke-sajadah/30', content: 'text-sajadah' },
  bata: { shape: 'fill-bata-muda stroke-bata/30', content: 'text-bata' },
  kuningan: { shape: 'fill-kuningan-muda stroke-kuningan/40', content: 'text-kuningan' },
  outline: { shape: 'fill-transparent stroke-batu-tua', content: 'text-jati-muda' },
};

const style = computed(() => ({ width: `${props.size}px`, height: `${props.size * 1.25}px` }));
</script>

<template>
  <span :class="cn('relative inline-flex shrink-0 items-end justify-center', props.class)" :style="style" aria-hidden="true">
    <svg viewBox="0 0 40 50" class="absolute inset-0 h-full w-full" preserveAspectRatio="xMidYMid meet">
      <path :d="ARCH_PATH" :class="toneClass[tone].shape" stroke-width="1" vector-effect="non-scaling-stroke" />
    </svg>
    <span :class="cn('relative mb-[22%] flex items-center justify-center', toneClass[tone].content)">
      <slot />
    </span>
  </span>
</template>
