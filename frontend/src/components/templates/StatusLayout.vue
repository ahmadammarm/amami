<script setup lang="ts">
import type { ArchTone } from '../atoms/arch';
import { ArchBadge } from '../atoms/arch';
import BrandMark from '../molecules/BrandMark.vue';

withDefaults(defineProps<{
  /** Short status label shown above the title, e.g. "404". */
  code?: string;
  title: string;
  description: string;
  tone?: ArchTone;
}>(), {
  tone: 'kuningan',
});
</script>

<template>
  <div class="flex min-h-dvh flex-col bg-kapur">
    <header class="px-5 py-5 sm:px-8">
      <RouterLink :to="{ name: 'Landing' }" class="inline-flex rounded-md">
        <BrandMark />
      </RouterLink>
    </header>

    <main class="flex flex-1 items-center justify-center px-5 pb-16">
      <div class="w-full max-w-md text-center">
        <ArchBadge :tone="tone" :size="72" class="mx-auto">
          <slot name="icon" />
        </ArchBadge>
        <p v-if="code" class="mt-6 font-mono text-sm tracking-[0.3em] text-kuningan-tua">{{ code }}</p>
        <h1 class="mt-2 font-display text-3xl font-semibold tracking-wide text-balance text-jati sm:text-4xl">{{ title }}</h1>
        <p class="mt-3 leading-relaxed text-pretty text-jati-muda">{{ description }}</p>
        <div v-if="$slots.default" class="mt-6">
          <slot />
        </div>
        <div v-if="$slots.actions" class="mt-8 flex flex-col-reverse justify-center gap-3 sm:flex-row">
          <slot name="actions" />
        </div>
      </div>
    </main>
  </div>
</template>
