<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue';
import { RotateCw, TriangleAlert } from 'lucide-vue-next';
import StatusLayout from '../templates/StatusLayout.vue';
import { Button } from '../atoms/button';

const errorMessage = ref<string | null>(null);

onErrorCaptured((err: unknown) => {
  errorMessage.value = err instanceof Error ? err.message : String(err);
  console.error(err);
  // Stop propagation: this boundary replaces the broken view.
  return false;
});

const reloadPage = () => {
  window.location.reload();
};
</script>

<template>
  <StatusLayout
    v-if="errorMessage !== null"
    title="Terjadi kesalahan"
    description="Halaman ini gagal ditampilkan. Muat ulang halaman; jika masih terjadi, sampaikan pesan di bawah kepada admin masjid."
    tone="bata"
  >
    <template #icon><TriangleAlert class="size-7" /></template>

    <p class="scrollbar-tipis max-h-32 overflow-auto rounded-lg border border-batu bg-white px-3 py-2 text-left font-mono text-xs break-words text-jati-muda">
      {{ errorMessage }}
    </p>

    <template #actions>
      <Button @click="reloadPage"><RotateCw /> Muat ulang halaman</Button>
    </template>
  </StatusLayout>
  <slot v-else />
</template>
