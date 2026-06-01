<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { useMosqueProfileQuery, useAuditLogsQuery } from '../../composables/useSettingsQuery';
import { Building2, Users, HandHeart, Calendar, History, User, Clock, Edit3, Loader2 } from 'lucide-vue-next';
import SuccessDialog from '../atoms/SuccessDialog.vue';

const authStore = useAuthStore();
const route = useRoute();
const router = useRouter();
const { data: profile, isLoading } = useMosqueProfileQuery();

// Fetch last 5 activity logs
const { data: auditData, isLoading: isLogsLoading } = useAuditLogsQuery(1, 5);
const recentActivities = computed(() => auditData.value?.logs || []);

const showLoginSuccess = ref(false);

const stats = [
  { name: 'Total Jamaah', value: '1,240', icon: Users, color: 'text-blue-600', bg: 'bg-blue-50' },
  { name: 'Infaq Month', value: 'Rp 12.5M', icon: HandHeart, color: 'text-emerald-600', bg: 'bg-emerald-50' },
  { name: 'Active Events', value: '12', icon: Calendar, color: 'text-amber-600', bg: 'bg-amber-50' },
];

onMounted(() => {
  if (route.query.loginSuccess === 'true') {
    showLoginSuccess.value = true;
    const query = { ...route.query };
    delete query.loginSuccess;
    router.replace({ query });
  }
});
</script>

<template>
  <div class="max-w-6xl mx-auto space-y-8">
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-3xl font-bold text-gray-900 leading-tight">Assalamu'alaikum, {{ authStore.user?.username }}!</h2>
        <p class="text-gray-500 font-medium mt-1">Welcome back to the mosque management dashboard.</p>
      </div>
      
      <div v-if="!isLoading && profile" class="flex items-center gap-3 px-4 py-2.5 bg-white rounded-xl shadow-sm border border-gray-100">
        <div class="p-1.5 bg-primary/10 rounded-lg">
          <Building2 class="w-5 h-5 text-primary" />
        </div>
        <span class="text-sm font-bold text-gray-700 tracking-tight">{{ profile.name }}</span>
      </div>
    </header>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div 
        v-for="stat in stats" 
        :key="stat.name"
        class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100 flex items-center gap-5 transition-all hover:shadow-md hover:translate-y-[-2px]"
      >
        <div :class="[stat.bg, stat.color, 'p-4 rounded-xl shadow-sm']">
          <component :is="stat.icon" class="w-7 h-7" />
        </div>
        <div>
          <p class="text-xs font-bold text-gray-400 uppercase tracking-widest">{{ stat.name }}</p>
          <p class="text-2xl font-bold text-gray-900 mt-0.5">{{ stat.value }}</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Left Column: Overview -->
      <section class="lg:col-span-2 bg-white p-8 rounded-2xl shadow-sm border border-gray-100">
        <h3 class="text-xl font-bold text-gray-900 mb-6 flex items-center gap-3">
          <div class="w-1.5 h-6 bg-primary rounded-full"></div>
          Financial Performance
        </h3>
        
        <div class="aspect-[16/9] bg-gray-50/50 rounded-xl flex items-center justify-center border-2 border-dashed border-gray-200">
          <div class="text-center space-y-2">
            <p class="text-gray-400 font-medium italic">Activity charts and trends visualization will appear here.</p>
            <p class="text-xs text-gray-300 uppercase tracking-widest font-bold">Module Under Construction</p>
          </div>
        </div>
      </section>

      <!-- Right Column: Recent Activities -->
      <section class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden flex flex-col">
        <div class="px-6 py-5 border-b border-gray-50 bg-gray-50/30 flex items-center justify-between">
          <h3 class="text-lg font-bold text-gray-900 flex items-center gap-2">
            <History class="w-5 h-5 text-primary" />
            Recent Activities
          </h3>
          <router-link v-if="authStore.userRole === 'SUPER_ADMIN'" to="/settings" class="text-xs font-bold text-primary hover:underline uppercase tracking-tighter">View All</router-link>
        </div>
        
        <div class="flex-1 overflow-y-auto max-h-[500px]">
          <div v-if="isLogsLoading" class="p-8 text-center">
            <Loader2 class="w-8 h-8 animate-spin text-primary/40 mx-auto" />
          </div>
          <div v-else-if="!recentActivities?.length" class="p-12 text-center text-gray-400 italic text-sm">
            No recent activities found.
          </div>
          <div v-else class="divide-y divide-gray-50">
            <div 
              v-for="log in recentActivities" 
              :key="log.id"
              class="px-6 py-4 hover:bg-gray-50/50 transition-colors"
            >
              <div class="flex items-start gap-3">
                <div class="mt-1 p-1.5 bg-gray-100 rounded-lg text-gray-500">
                  <User v-if="log.action === 'LOGIN'" class="w-3.5 h-3.5" />
                  <Edit3 v-else class="w-3.5 h-3.5" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-sm font-bold text-gray-800 truncate">
                    {{ log.user }} <span class="font-medium text-gray-500">{{ log.action.toLowerCase() }}d</span>
                  </p>
                  <p class="text-xs text-gray-500 font-medium truncate mt-0.5">
                    {{ log.entity }}
                  </p>
                  <div class="flex items-center gap-1.5 mt-1.5 text-[10px] text-gray-400 font-bold uppercase tracking-tight">
                    <Clock class="w-3 h-3" />
                    {{ log.timestamp }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <!-- Welcome Success Dialog -->
    <SuccessDialog
      v-model:open="showLoginSuccess"
      title="Assalamu'alaikum!"
      :description="`Welcome back, ${authStore.user?.username}. You have successfully signed in to the amami ERP system.`"
      buttonText="Start Working"
    />
  </div>
</template>
