<script setup lang="ts">
import { useAuthStore } from '../../stores/auth';
import { useMosqueProfileQuery } from '../../composables/useSettingsQuery';
import { Building2, Users, HandHeart, Calendar } from 'lucide-vue-next';

const authStore = useAuthStore();
const { data: profile, isLoading } = useMosqueProfileQuery();

const stats = [
  { name: 'Total Jamaah', value: '1,240', icon: Users, color: 'text-blue-600', bg: 'bg-blue-50' },
  { name: 'Infaq Month', value: 'Rp 12.5M', icon: HandHeart, color: 'text-emerald-600', bg: 'bg-emerald-50' },
  { name: 'Active Events', value: '12', icon: Calendar, color: 'text-amber-600', bg: 'bg-amber-50' },
];
</script>

<template>
  <div class="max-w-6xl mx-auto space-y-8">
    <header class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-bold text-gray-800">Assalamu'alaikum, {{ authStore.user?.username }}!</h2>
        <p class="text-gray-500 font-medium mt-1">Welcome back to the mosque management dashboard.</p>
      </div>
      
      <div v-if="!isLoading && profile" class="flex items-center gap-3 px-4 py-2 bg-white rounded-lg shadow-sm border border-gray-100">
        <Building2 class="w-5 h-5 text-primary" />
        <span class="text-sm font-bold text-gray-700">{{ profile.name }}</span>
      </div>
    </header>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div 
        v-for="stat in stats" 
        :key="stat.name"
        class="bg-white p-6 rounded-lg shadow-md border border-gray-50 flex items-center gap-4 transition-transform hover:scale-[1.02]"
      >
        <div :class="[stat.bg, stat.color, 'p-4 rounded-lg shadow-sm']">
          <component :is="stat.icon" class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs font-bold text-gray-400 uppercase tracking-wider">{{ stat.name }}</p>
          <p class="text-2xl font-bold text-gray-800">{{ stat.value }}</p>
        </div>
      </div>
    </div>

    <section class="bg-white p-8 rounded-lg shadow-md border border-gray-50">
      <h3 class="text-xl font-bold text-gray-800 mb-6 flex items-center gap-2">
        <div class="w-2 h-6 bg-primary rounded-full"></div>
        Quick Overview
      </h3>
      
      <div class="aspect-video bg-gray-50 rounded-lg flex items-center justify-center border-2 border-dashed border-gray-200">
        <p class="text-gray-400 font-medium italic">Financial trends and activity charts will appear here.</p>
      </div>
    </section>
  </div>
</template>
