<script setup lang="ts">
import { useId } from 'vue';
import { cn } from '@/lib/utils';
import { ARCH_CLIP_PATH } from './arch';

const props = defineProps<{
  src: string;
  alt: string;
  class?: string;
  imgClass?: string;
}>();

// Each instance needs its own clipPath id; the shape uses objectBoundingBox
// units, so it scales with whatever size the frame is given.
const clipId = `arch-${useId()}`;
</script>

<template>
  <div :class="cn('relative aspect-[4/5]', props.class)">
    <svg width="0" height="0" class="absolute" aria-hidden="true" focusable="false">
      <clipPath :id="clipId" clipPathUnits="objectBoundingBox">
        <path :d="ARCH_CLIP_PATH" />
      </clipPath>
    </svg>
    <img
      :src="src"
      :alt="alt"
      :class="cn('h-full w-full object-cover', imgClass)"
      :style="{ clipPath: `url(#${clipId})` }"
    />
  </div>
</template>
