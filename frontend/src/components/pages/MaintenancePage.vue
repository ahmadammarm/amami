<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch, type Component } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useIntervalFn, useOnline } from '@vueuse/core';
import { Loader2, RotateCw, ServerOff, WifiOff, Wrench } from 'lucide-vue-next';
import { pingServer } from '../../api/client';
import { classifyServiceError, parseServiceIssue, safeReturnPath, type ServiceIssue } from '../../lib/serviceStatus';
import type { ArchTone } from '../atoms/arch';
import StatusLayout from '../templates/StatusLayout.vue';
import { Button } from '../atoms/button';

const RETRY_SECONDS = 15;

const COPY: Record<ServiceIssue, { code: string; title: string; description: string; icon: Component; tone: ArchTone }> = {
  maintenance: {
    code: '503',
    title: 'Sedang dalam pemeliharaan',
    description: 'amami sedang diperbarui oleh pengelola sistem. Data Anda aman, dan halaman akan terbuka kembali setelah pemeliharaan selesai.',
    icon: Wrench,
    tone: 'kuningan',
  },
  unreachable: {
    code: 'SERVER',
    title: 'Server tidak dapat dihubungi',
    description: 'Aplikasi belum bisa terhubung ke server masjid. Biasanya ini hanya sementara.',
    icon: ServerOff,
    tone: 'bata',
  },
  offline: {
    code: 'OFFLINE',
    title: 'Perangkat Anda sedang offline',
    description: 'Periksa koneksi internet atau Wi-Fi Anda. Halaman akan terbuka kembali begitu koneksi tersambung.',
    icon: WifiOff,
    tone: 'outline',
  },
};

const route = useRoute();
const router = useRouter();
const isOnline = useOnline();

const issue = ref<ServiceIssue>(parseServiceIssue(route.query.alasan));
const returnPath = computed(() => safeReturnPath(route.query.kembali));
const copy = computed(() => COPY[issue.value]);

const isChecking = ref(false);
const secondsLeft = ref(RETRY_SECONDS);

const checkServer = async () => {
  if (isChecking.value) return;
  isChecking.value = true;
  try {
    await pingServer();
    await router.replace(returnPath.value);
  } catch (error) {
    // Still down: show the current reason and wait for the next attempt.
    issue.value = classifyServiceError(error) ?? issue.value;
    secondsLeft.value = RETRY_SECONDS;
  } finally {
    isChecking.value = false;
  }
};

const { pause } = useIntervalFn(() => {
  if (isChecking.value) return;
  secondsLeft.value -= 1;
  if (secondsLeft.value <= 0) void checkServer();
}, 1000);

// The moment the device reconnects, check instead of waiting for the countdown.
watch(isOnline, (online) => {
  if (online) void checkServer();
  else issue.value = 'offline';
});

onBeforeUnmount(pause);
</script>

<template>
  <StatusLayout :code="copy.code" :title="copy.title" :description="copy.description" :tone="copy.tone">
    <template #icon><component :is="copy.icon" class="size-7" /></template>

    <p class="text-sm text-jati-muda" aria-live="polite">
      <template v-if="isChecking">Memeriksa koneksi ke server…</template>
      <template v-else>Mencoba lagi otomatis dalam <span class="tabular font-semibold text-jati">{{ secondsLeft }}</span> detik.</template>
    </p>

    <template #actions>
      <Button :disabled="isChecking" @click="checkServer">
        <Loader2 v-if="isChecking" class="animate-spin" />
        <RotateCw v-else />
        Coba lagi sekarang
      </Button>
    </template>
  </StatusLayout>
</template>
