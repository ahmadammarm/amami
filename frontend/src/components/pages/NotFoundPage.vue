<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ArrowLeft, Compass } from 'lucide-vue-next';
import { useAuthStore } from '../../stores/auth';
import StatusLayout from '../templates/StatusLayout.vue';
import { Button } from '../atoms/button';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const home = computed(() =>
  authStore.isAuthenticated
    ? { label: 'Ke beranda dashboard', to: { name: 'Dashboard' } }
    : { label: 'Ke halaman depan', to: { name: 'Landing' } },
);

/** history.state.back is set by Vue Router when there is an in-app page to return to. */
const canGoBack = computed(() => typeof window !== 'undefined' && window.history.state?.back != null);
</script>

<template>
  <StatusLayout
    code="404"
    title="Halaman tidak ditemukan"
    description="Alamat yang Anda buka tidak ada. Mungkin salah ketik, atau halamannya sudah dipindahkan."
  >
    <template #icon><Compass class="size-7" /></template>

    <p class="mx-auto max-w-full truncate rounded-lg border border-batu bg-white px-3 py-2 font-mono text-xs text-jati-muda">
      {{ route.fullPath }}
    </p>

    <template #actions>
      <Button v-if="canGoBack" variant="outline" @click="router.back()"><ArrowLeft /> Kembali</Button>
      <Button @click="router.push(home.to)">{{ home.label }}</Button>
    </template>
  </StatusLayout>
</template>
