<script setup lang="ts">
import { Users } from 'lucide-vue-next';
</script>

<template>
  <div class="max-w-6xl mx-auto space-y-8">
    <header>
      <h2 class="text-3xl font-bold text-gray-900 leading-tight">Jamaah</h2>
      <p class="text-gray-500 font-medium mt-1">Community database and scoring.</p>
    </header>

    <div class="bg-white p-12 rounded-2xl shadow-sm border border-gray-100 flex flex-col items-center justify-center text-center space-y-4">
      <div class="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
        <Users class="w-8 h-8 text-primary" />
      </div>
      <h3 class="text-xl font-bold text-gray-900">Module Under Construction</h3>
      <p class="text-gray-500 font-medium max-w-md">
        The Jamaah registry module is currently being developed.
      </p>
    </div>
  </div>
</template>
