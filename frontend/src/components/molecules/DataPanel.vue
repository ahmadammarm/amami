<script setup lang="ts">
defineProps<{
  title?: string;
  description?: string;
}>();
</script>

<template>
  <section class="overflow-hidden rounded-xl border border-batu bg-white shadow-sm">
    <header
      v-if="title || $slots.actions"
      class="flex flex-col gap-3 border-b border-batu px-5 py-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <h2 v-if="title" class="font-semibold text-jati">{{ title }}</h2>
        <p v-if="description" class="mt-0.5 text-sm text-jati-muda">{{ description }}</p>
      </div>
      <div v-if="$slots.actions" class="flex flex-wrap items-center gap-2">
        <slot name="actions" />
      </div>
    </header>
    <slot />
    <footer v-if="$slots.footer" class="border-t border-batu bg-kapur/60 px-5 py-3">
      <slot name="footer" />
    </footer>
  </section>
</template>
