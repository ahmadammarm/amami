<script setup lang="ts">
import { computed } from 'vue';
import heroBanner from '../../assets/hero-banner.jpg';
import { useMosqueProfileQuery } from '../../composables/useSettingsQuery';
import { ArchFrame } from '../atoms/arch';
import BrandMark from '../molecules/BrandMark.vue';

defineProps<{
  title: string;
  description: string;
}>();

const { data: profile } = useMosqueProfileQuery();
const mosqueName = computed(() => profile.value?.name || 'Masjid kami');
</script>

<template>
  <div class="grid min-h-dvh bg-kapur lg:grid-cols-[1fr_1.1fr]">
    <!-- Mosque panel (desktop) -->
    <aside class="pola-girih relative hidden flex-col justify-between overflow-hidden bg-sajadah-tua p-12 text-white lg:flex">
      <RouterLink :to="{ name: 'Landing' }" class="w-fit rounded-md">
        <BrandMark variant="light" show-tagline />
      </RouterLink>
      <div class="mx-auto w-full max-w-xs">
        <ArchFrame :src="heroBanner" :alt="`Halaman ${mosqueName}`" class="w-full" />
      </div>
      <div>
        <p class="text-xs font-semibold uppercase tracking-[0.2em] text-kuningan">Panel pengurus</p>
        <p class="mt-2 font-display text-3xl font-semibold tracking-wide">{{ mosqueName }}</p>
      </div>
    </aside>

    <!-- Form -->
    <main class="flex items-center justify-center px-5 py-12 sm:px-10">
      <div class="w-full max-w-sm">
        <RouterLink :to="{ name: 'Landing' }" class="mb-10 inline-flex rounded-md lg:hidden">
          <BrandMark />
        </RouterLink>
        <h1 class="font-display text-3xl font-semibold tracking-wide text-jati">{{ title }}</h1>
        <p class="mt-2 text-sm leading-relaxed text-jati-muda">{{ description }}</p>
        <div class="mt-8">
          <slot />
        </div>
      </div>
    </main>
  </div>
</template>
