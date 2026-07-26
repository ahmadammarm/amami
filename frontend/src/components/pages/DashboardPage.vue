<script setup lang="ts">
import { computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { 
  Wallet, 
  Users, 
  Activity, 
  Package, 
  BookOpen, 
  HeartHandshake, 
  TrendingUp, 
  AlertTriangle,
  Server,
  FileText
} from 'lucide-vue-next';
import { useQuery } from '@tanstack/vue-query';
import { dashboardService } from '@/api/services/dashboard';

const authStore = useAuthStore();
const role = computed(() => authStore.userRole);
const userName = computed(() => authStore.user?.Username || 'User');

// Fetch Metrics from Backend
const { data: metrics, isLoading } = useQuery({
  queryKey: ['dashboard-metrics'],
  queryFn: dashboardService.getMetrics,
  refetchInterval: 60000, // refresh every minute
});

// Format utilities
const formatCurrency = (amount: number = 0) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amount);
};

const greeting = computed(() => {
  const hour = new Date().getHours();
  if (hour < 12) return 'Selamat Pagi';
  if (hour < 15) return 'Selamat Siang';
  if (hour < 18) return 'Selamat Sore';
  return 'Selamat Malam';
});

// Role checks
const isSuperAdmin = computed(() => role.value === 'SUPER_ADMIN');
const isBendahara = computed(() => role.value === 'BENDAHARA');
const isTakmir = computed(() => role.value === 'TAKMIR');
const isSekretaris = computed(() => role.value === 'SEKRETARIS');
</script>

<template>
  <div class="space-y-8 animate-in fade-in duration-500 pb-10">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-center md:justify-between bg-gradient-to-r from-emerald-700 to-teal-800 rounded-2xl p-8 text-white shadow-lg overflow-hidden relative">
      <!-- Decorative background element -->
      <div class="absolute -right-20 -top-20 opacity-10">
        <svg width="300" height="300" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v20"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
      </div>
      
      <div class="z-10 relative">
        <h1 class="text-3xl font-bold tracking-tight">{{ greeting }}, {{ userName }}!</h1>
        <p class="mt-2 text-emerald-100 text-lg max-w-2xl">
          <span v-if="isSuperAdmin">Sistem kontrol utama Amami. Seluruh operasional masjid terpantau dengan baik.</span>
          <span v-if="isBendahara">Ringkasan keuangan dan Zakat hari ini. Transparansi dan akuntabilitas adalah prioritas kita.</span>
          <span v-if="isTakmir">Ringkasan operasional masjid. Mari kita pastikan semua logistik dan inventaris berjalan lancar.</span>
          <span v-if="isSekretaris">Pusat data jamaah dan jadwal masjid. Pastikan semua administrasi terorganisir.</span>
        </p>
        <div class="mt-4 inline-flex items-center rounded-full bg-emerald-900/50 px-3 py-1 text-sm font-medium text-emerald-100 backdrop-blur-sm border border-emerald-500/30">
          <Activity class="w-4 h-4 mr-2" /> 
          Peran Anda saat ini: <span class="font-bold ml-1">{{ role }}</span>
        </div>
      </div>
    </div>

    <!-- BENDAHARA / SUPER ADMIN: Finance & Zakat -->
    <div v-if="isBendahara || isSuperAdmin" class="space-y-4">
      <div class="flex items-center text-gray-800">
        <Wallet class="w-6 h-6 mr-2 text-emerald-600" />
        <h2 class="text-xl font-bold">Keuangan & Zakat</h2>
      </div>
      <div class="grid gap-6 md:grid-cols-3">
        <!-- Kas Card -->
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between text-gray-500">
            <h3 class="text-sm font-medium">Total Saldo Kas</h3>
            <TrendingUp class="w-4 h-4 text-emerald-500" />
          </div>
          <div class="mt-4 text-3xl font-bold text-gray-900">{{ formatCurrency(metrics?.finance.total_balance) }}</div>
          <p class="mt-2 text-sm text-emerald-600 font-medium">+ {{ formatCurrency(metrics?.finance.monthly_inflow) }} bulan ini</p>
        </div>

        <!-- Zakat Card -->
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between text-gray-500">
            <h3 class="text-sm font-medium">Zakat Terkumpul</h3>
            <HeartHandshake class="w-4 h-4 text-blue-500" />
          </div>
          <div class="mt-4 text-3xl font-bold text-gray-900">{{ formatCurrency(metrics?.zakat.collected) }}</div>
          <div class="mt-4 h-2 w-full bg-gray-100 rounded-full overflow-hidden">
            <div class="h-full bg-blue-500 rounded-full" style="width: 15%"></div>
          </div>
          <p class="mt-2 text-xs text-gray-400">Distribusi belum dimulai</p>
        </div>

        <!-- Mustahik Card -->
        <div class="rounded-2xl border border-gray-100 bg-gradient-to-br from-gray-900 to-gray-800 p-6 text-white shadow-sm hover:shadow-lg transition-shadow">
          <div class="flex items-center justify-between text-gray-300">
            <h3 class="text-sm font-medium">Verifikasi Mustahik</h3>
            <Users class="w-4 h-4 text-gray-300" />
          </div>
          <div class="mt-4 text-4xl font-bold">{{ metrics?.zakat.mustahik || 0 }} <span class="text-lg font-normal text-gray-400">KK</span></div>
          <p class="mt-2 text-sm text-gray-300">Terdaftar dan terverifikasi untuk menerima zakat.</p>
        </div>
      </div>
    </div>

    <!-- TAKMIR / SUPER ADMIN: Logistics & Qurban -->
    <div v-if="isTakmir || isSuperAdmin" class="space-y-4">
      <div class="flex items-center text-gray-800">
        <Package class="w-6 h-6 mr-2 text-orange-500" />
        <h2 class="text-xl font-bold">Logistik & Qurban</h2>
      </div>
      <div class="grid gap-6 md:grid-cols-3">
        <!-- Qurban Card -->
        <div class="rounded-2xl border border-orange-100 bg-orange-50 p-6 shadow-sm">
          <h3 class="text-sm font-medium text-orange-800">Hewan Qurban</h3>
          <div class="mt-4 flex items-baseline gap-4">
            <div class="text-3xl font-bold text-orange-900">{{ metrics?.qurban.cows || 0 }} <span class="text-sm font-normal text-orange-700">Sapi</span></div>
            <div class="text-3xl font-bold text-orange-900">{{ metrics?.qurban.goats || 0 }} <span class="text-sm font-normal text-orange-700">Kambing</span></div>
          </div>
          <p class="mt-2 text-sm text-orange-700">Tahun ini</p>
        </div>

        <!-- Inventory Alerts -->
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between text-gray-500">
              <h3 class="text-sm font-medium">Status Inventaris</h3>
              <AlertTriangle class="w-4 h-4 text-red-500" />
            </div>
            <div class="mt-4 flex items-baseline gap-2">
              <span class="text-3xl font-bold text-red-600">{{ metrics?.logistics.broken_assets || 0 }}</span>
              <span class="text-sm text-gray-500">aset butuh perbaikan</span>
            </div>
          </div>
          <button class="mt-4 w-full rounded-lg bg-gray-50 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100">Lihat Detail Inventaris</button>
        </div>
        
        <!-- Peminjaman -->
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
          <div class="flex items-center justify-between text-gray-500">
            <h3 class="text-sm font-medium">Barang Dipinjam</h3>
            <BookOpen class="w-4 h-4" />
          </div>
          <div class="mt-4 text-4xl font-bold text-gray-900">{{ metrics?.logistics.active_loans || 0 }}</div>
          <p class="mt-2 text-sm text-gray-500">Peminjaman aktif oleh jamaah saat ini.</p>
        </div>
      </div>
    </div>

    <!-- SEKRETARIS / SUPER ADMIN: Administration & Jamaah -->
    <div v-if="isSekretaris || isSuperAdmin" class="space-y-4">
      <div class="flex items-center text-gray-800">
        <FileText class="w-6 h-6 mr-2 text-blue-600" />
        <h2 class="text-xl font-bold">Administrasi & Jamaah</h2>
      </div>
      <div class="grid gap-6 md:grid-cols-2">
        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm flex items-center justify-between">
          <div>
            <h3 class="text-sm font-medium text-gray-500">Total Jamaah Terdaftar</h3>
            <div class="mt-2 text-4xl font-bold text-gray-900">{{ metrics?.jamaah.total || 0 }}</div>
            <p class="mt-1 text-sm text-green-600">+{{ metrics?.jamaah.newly_registered || 0 }} bulan ini</p>
          </div>
          <div class="h-16 w-16 bg-blue-50 text-blue-600 rounded-full flex items-center justify-center">
            <Users class="w-8 h-8" />
          </div>
        </div>

        <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm flex items-center justify-between">
          <div>
            <h3 class="text-sm font-medium text-gray-500">Agenda Mendatang</h3>
            <div class="mt-2 text-4xl font-bold text-gray-900">{{ metrics?.logistics.upcoming_agendas || 0 }}</div>
            <p class="mt-1 text-sm text-gray-500">Kegiatan yang akan datang</p>
          </div>
          <div class="h-16 w-16 bg-purple-50 text-purple-600 rounded-full flex items-center justify-center">
            <BookOpen class="w-8 h-8" />
          </div>
        </div>
      </div>
    </div>

    <!-- SUPER ADMIN ONLY: System Health -->
    <div v-if="isSuperAdmin" class="space-y-4">
      <div class="flex items-center text-gray-800">
        <Server class="w-6 h-6 mr-2 text-gray-600" />
        <h2 class="text-xl font-bold">Status Sistem (Admin)</h2>
      </div>
      <div class="rounded-2xl border border-gray-200 bg-gray-50 p-6 shadow-inner flex flex-wrap items-center gap-8 text-sm">
        <div class="flex items-center gap-2">
          <div :class="['w-3 h-3 rounded-full', metrics?.system.db_status === 'Healthy' ? 'bg-green-500 animate-pulse' : 'bg-red-500']"></div>
          <span class="text-gray-600">Database: <span class="font-semibold text-gray-900">{{ metrics?.system.db_status || 'Unknown' }}</span></span>
        </div>
        <div class="flex items-center gap-2">
          <Activity class="w-4 h-4 text-gray-400" />
          <span class="text-gray-600">Uptime: <span class="font-semibold text-gray-900">{{ metrics?.system.uptime || '0%' }}</span></span>
        </div>
        <div class="flex items-center gap-2">
          <Users class="w-4 h-4 text-gray-400" />
          <span class="text-gray-600">Active Sessions: <span class="font-semibold text-gray-900">{{ metrics?.system.active_sessions || 0 }}</span></span>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
/* Minor animation utility since tailwind-animate is not installed by default */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-in {
  animation: fadeIn 0.5s ease-out forwards;
}
</style>
