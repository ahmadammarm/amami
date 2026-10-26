<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ArrowLeft, Lock } from 'lucide-vue-next';
import StatusLayout from '../templates/StatusLayout.vue';
import { Button } from '../atoms/button';

const router = useRouter();
</script>

<template>
  <StatusLayout
    code="403"
    title="Halaman ini bukan untuk peran Anda"
    description="Peran akun Anda tidak memiliki akses ke halaman ini. Jika Anda memerlukannya, minta admin masjid menyesuaikan peran Anda."
  >
    <template #icon><Lock class="size-7" /></template>

    <template #actions>
      <Button variant="outline" @click="router.back()"><ArrowLeft /> Kembali</Button>
      <Button @click="router.push({ name: 'Dashboard' })">Ke beranda dashboard</Button>
    </template>
  </StatusLayout>
</template>
