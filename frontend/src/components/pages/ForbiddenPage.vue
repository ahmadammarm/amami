<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ArrowLeft, Lock } from 'lucide-vue-next';
import { ArchBadge } from '../atoms/arch';
import { Button } from '../atoms/button';

const router = useRouter();
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-kapur px-5">
    <div class="w-full max-w-md text-center">
      <ArchBadge tone="kuningan" :size="64" class="mx-auto">
        <Lock class="size-6" />
      </ArchBadge>
      <h1 class="mt-6 font-display text-3xl font-semibold tracking-wide text-jati">Halaman ini bukan untuk peran Anda</h1>
      <p class="mt-3 leading-relaxed text-jati-muda">
        Peran akun Anda tidak memiliki akses ke halaman ini. Jika Anda memerlukannya, minta admin masjid menyesuaikan peran Anda.
      </p>
      <div class="mt-8 flex flex-col justify-center gap-3 sm:flex-row">
        <Button variant="outline" @click="router.back()"><ArrowLeft /> Kembali</Button>
        <Button @click="router.push({ name: 'Dashboard' })">Ke beranda</Button>
      </div>
    </div>
  </div>
</template>
