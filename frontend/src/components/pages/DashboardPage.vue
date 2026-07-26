<script setup lang="ts">
import { computed } from 'vue';
import VueApexCharts from 'vue3-apexcharts';
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
  FileText,
  ArrowUpRight,
  ArrowDownRight,
  CalendarDays,
  ShieldCheck,
} from 'lucide-vue-next';
import { useQuery } from '@tanstack/vue-query';
import { dashboardService } from '@/api/services/dashboard';
import type { MonthlyFinancialPoint, MonthlyJamaahPoint } from '@/types/dashboard';

const authStore = useAuthStore();
const role = computed(() => authStore.userRole);
const userName = computed(() => authStore.user?.full_name || authStore.user?.username || 'Admin');

const { data: metrics, isLoading } = useQuery({
  queryKey: ['dashboard-metrics'],
  queryFn: dashboardService.getMetrics,
  refetchInterval: 60000,
});

const formatCurrency = (amount: number = 0) =>
  new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(amount);

const formatNumber = (n: number = 0) => new Intl.NumberFormat('id-ID').format(n);

const greeting = computed(() => {
  const hour = new Date().getHours();
  if (hour < 12) return 'Selamat Pagi';
  if (hour < 15) return 'Selamat Siang';
  if (hour < 18) return 'Selamat Sore';
  return 'Selamat Malam';
});

const now = new Date();
const dateStr = now.toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });

const isSuperAdmin = computed(() => role.value === 'SUPER_ADMIN');
const isBendahara = computed(() => role.value === 'BENDAHARA');
const isTakmir = computed(() => role.value === 'TAKMIR');
const isSekretaris = computed(() => role.value === 'SEKRETARIS');

// ── Chart: 6-Month Finance Trend (Area) ──────────────────────────────────────
const financeAreaSeries = computed(() => {
  const data = metrics.value?.analytics?.monthly_finance ?? [];
  return [
    { name: 'Pemasukan', data: data.map((d: MonthlyFinancialPoint) => d.income) },
    { name: 'Pengeluaran', data: data.map((d: MonthlyFinancialPoint) => d.expense) },
  ];
});
const financeAreaOptions = computed(() => ({
  chart: {
    type: 'area',
    toolbar: { show: false },
    fontFamily: 'inherit',
    sparkline: { enabled: false },
    animations: { enabled: true, speed: 800 },
  },
  colors: ['#10b981', '#f43f5e'],
  fill: {
    type: 'gradient',
    gradient: { shadeIntensity: 1, opacityFrom: 0.45, opacityTo: 0.05 },
  },
  stroke: { curve: 'smooth', width: 2 },
  dataLabels: { enabled: false },
  xaxis: {
    categories: (metrics.value?.analytics?.monthly_finance ?? []).map((d: MonthlyFinancialPoint) => d.month),
    labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: {
      formatter: (v: number) => `Rp ${formatNumber(v)}`,
      style: { colors: '#9ca3af', fontSize: '11px' },
    },
    stepSize: 10000000,
  },
  grid: { borderColor: '#f3f4f6', strokeDashArray: 4, padding: { left: 0, right: 0 } },
  tooltip: {
    y: { formatter: (v: number) => formatCurrency(v) },
  },
  legend: { position: 'top', horizontalAlign: 'right', fontSize: '13px' },
}));

// ── Chart: Zakat Radial Bar ───────────────────────────────────────────────────
const zakatRadialSeries = computed(() => {
  const donated = metrics.value?.analytics?.zakat_breakdown?.total_donated ?? 0;
  const distributed = metrics.value?.analytics?.zakat_breakdown?.total_distributed ?? 0;
  if (donated === 0) return [0];
  return [Math.min(Math.round((distributed / donated) * 100), 100)];
});
const zakatRadialOptions = computed(() => ({
  chart: { type: 'radialBar', fontFamily: 'inherit' },
  colors: ['#3b82f6'],
  plotOptions: {
    radialBar: {
      hollow: { size: '60%' },
      dataLabels: {
        name: { show: true, color: '#6b7280', fontSize: '13px', offsetY: -8 },
        value: { show: true, color: '#111827', fontSize: '28px', fontWeight: 700, formatter: (v: number) => `${v}%` },
      },
      track: { background: '#e5e7eb' },
    },
  },
  labels: ['Tersalurkan'],
}));

// ── Chart: Qurban Donut (Sapi, Kambing, Domba) ───────────────────────────────
const qurbanDonutSeries = computed(() => [
  metrics.value?.analytics?.qurban_breakdown?.sapi ?? 0,
  metrics.value?.analytics?.qurban_breakdown?.kambing ?? 0,
  metrics.value?.analytics?.qurban_breakdown?.domba ?? 0,
]);
const qurbanDonutOptions = computed(() => ({
  chart: { type: 'donut', fontFamily: 'inherit', animations: { enabled: true } },
  labels: ['Sapi', 'Kambing', 'Domba'],
  colors: ['#f97316', '#f59e0b', '#84cc16'],
  dataLabels: { enabled: true, formatter: (_v: number, opts: { seriesIndex: number; w: { config: { labels: string[] }; globals: { series: number[] } } }) => `${opts.w.config.labels[opts.seriesIndex]}: ${opts.w.globals.series[opts.seriesIndex]}` },
  legend: { position: 'bottom', fontSize: '13px' },
  stroke: { width: 2, colors: ['#ffffff'] },
  plotOptions: { pie: { donut: { size: '62%', labels: { show: true, total: { show: true, label: 'Total Hewan', color: '#6b7280', formatter: (w: { globals: { seriesTotals: number[] } }) => w.globals.seriesTotals.reduce((a: number, b: number) => a + b, 0).toString() } } } } },
  tooltip: { y: { formatter: (v: number) => `${v} ekor` } },
}));

// ── Chart: Asset Status Bar ───────────────────────────────────────────────────
const assetBarSeries = computed(() => [{
  name: 'Aset',
  data: [
    metrics.value?.analytics?.asset_breakdown?.good ?? 0,
    metrics.value?.analytics?.asset_breakdown?.loaned ?? 0,
    metrics.value?.analytics?.asset_breakdown?.needs_repair ?? 0,
    metrics.value?.analytics?.asset_breakdown?.broken ?? 0,
  ],
}]);
const assetBarOptions = computed(() => ({
  chart: { type: 'bar', toolbar: { show: false }, fontFamily: 'inherit', animations: { enabled: true } },
  colors: ['#10b981', '#f59e0b', '#f97316', '#f43f5e'],
  plotOptions: {
    bar: {
      borderRadius: 6,
      columnWidth: '45%',
      distributed: true,
    },
  },
  dataLabels: { enabled: false },
  legend: { show: false },
  xaxis: {
    categories: ['Baik', 'Dipinjam', 'Perlu Perbaikan', 'Rusak'],
    labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    tickAmount: 4,
  },
  grid: { borderColor: '#f3f4f6', strokeDashArray: 4 },
  tooltip: { y: { formatter: (v: number) => `${v} unit` } },
}));

// ── Chart: Jamaah Monthly Growth Bar ─────────────────────────────────────────
const jamaahBarSeries = computed(() => [{
  name: 'Jamaah Baru',
  data: (metrics.value?.analytics?.jamaah_growth ?? []).map((d: MonthlyJamaahPoint) => d.new_members),
}]);
const jamaahBarOptions = computed(() => ({
  chart: { type: 'bar', toolbar: { show: false }, fontFamily: 'inherit' },
  colors: ['#6366f1'],
  plotOptions: {
    bar: { borderRadius: 6, columnWidth: '40%' },
  },
  dataLabels: { enabled: false },
  xaxis: {
    categories: (metrics.value?.analytics?.jamaah_growth ?? []).map((d: MonthlyJamaahPoint) => d.month),
    labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: { style: { colors: '#9ca3af', fontSize: '12px' } },
    tickAmount: 4,
  },
  grid: { borderColor: '#f3f4f6', strokeDashArray: 4 },
  tooltip: { y: { formatter: (v: number) => `${v} orang` } },
}));
</script>

<template>
  <div class="space-y-8 pb-12 animate-in fade-in duration-500">

    <!-- ── Hero Header ─────────────────────────────────────────────────────── -->
    <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-emerald-700 via-teal-700 to-teal-900 p-8 text-white shadow-xl">
      <!-- Decorative circles -->
      <div class="absolute -right-16 -top-16 h-64 w-64 rounded-full bg-white/5"></div>
      <div class="absolute -bottom-10 -right-4 h-40 w-40 rounded-full bg-white/5"></div>
      <div class="absolute top-8 right-40 h-20 w-20 rounded-full bg-white/5"></div>

      <div class="relative z-10 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm font-medium text-emerald-200">{{ dateStr }}</p>
          <h1 class="mt-1 text-3xl font-bold tracking-tight">{{ greeting }}, {{ userName }}!</h1>
          <p class="mt-2 max-w-lg text-emerald-100">
            <span v-if="isSuperAdmin">Sistem kontrol utama Amami — seluruh operasional terpantau secara real-time.</span>
            <span v-else-if="isBendahara">Pantau keuangan dan distribusi zakat dengan akurat dan transparan.</span>
            <span v-else-if="isTakmir">Kelola logistik, inventaris, dan kegiatan masjid di satu tempat.</span>
            <span v-else-if="isSekretaris">Pusat administrasi jamaah, agenda, dan dokumen masjid.</span>
          </p>
        </div>
        <div class="flex flex-col items-start gap-2 md:items-end">
          <div class="inline-flex items-center rounded-full bg-white/10 px-3 py-1 text-sm font-medium text-white backdrop-blur-sm border border-white/20">
            <ShieldCheck class="mr-2 h-4 w-4" />
            {{ role }}
          </div>
        </div>
      </div>
    </div>

    <!-- Skeleton Loader -->
    <div v-if="isLoading" class="grid gap-6 md:grid-cols-4">
      <div v-for="i in 4" :key="i" class="h-32 animate-pulse rounded-2xl bg-gray-100"></div>
    </div>

    <template v-else>

      <!-- ── KPI Summary Strip ──────────────────────────────────────────────── -->
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <!-- Saldo Kas -->
        <div v-if="isSuperAdmin || isBendahara" class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Total Saldo Kas</p>
            <div class="rounded-xl bg-emerald-50 p-2 text-emerald-600"><Wallet class="h-5 w-5" /></div>
          </div>
          <p class="mt-4 text-2xl font-bold text-gray-900">{{ formatCurrency(metrics?.finance?.total_balance) }}</p>
          <p class="mt-1 flex items-center gap-1 text-xs font-medium text-emerald-600">
            <ArrowUpRight class="h-3 w-3" /> {{ formatCurrency(metrics?.finance?.monthly_inflow) }} bulan ini
          </p>
        </div>

        <!-- Total Jamaah -->
        <div v-if="isSuperAdmin || isSekretaris || isTakmir" class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Total Jamaah</p>
            <div class="rounded-xl bg-blue-50 p-2 text-blue-600"><Users class="h-5 w-5" /></div>
          </div>
          <p class="mt-4 text-2xl font-bold text-gray-900">{{ formatNumber(metrics?.jamaah?.total) }}</p>
          <p class="mt-1 flex items-center gap-1 text-xs font-medium text-blue-600">
            <ArrowUpRight class="h-3 w-3" /> +{{ metrics?.jamaah?.newly_registered ?? 0 }} bulan ini
          </p>
        </div>

        <!-- Aset Dipinjam -->
        <div v-if="isSuperAdmin || isTakmir" class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Aset Dipinjam</p>
            <div class="rounded-xl bg-amber-50 p-2 text-amber-600"><BookOpen class="h-5 w-5" /></div>
          </div>
          <p class="mt-4 text-2xl font-bold text-gray-900">{{ metrics?.logistics?.active_loans ?? 0 }}</p>
          <p class="mt-1 text-xs text-gray-400">peminjaman aktif saat ini</p>
        </div>

        <!-- Agenda Mendatang -->
        <div v-if="isSuperAdmin || isTakmir || isSekretaris" class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Agenda Mendatang</p>
            <div class="rounded-xl bg-violet-50 p-2 text-violet-600"><CalendarDays class="h-5 w-5" /></div>
          </div>
          <p class="mt-4 text-2xl font-bold text-gray-900">{{ metrics?.logistics?.upcoming_agendas ?? 0 }}</p>
          <p class="mt-1 text-xs text-gray-400">kegiatan akan datang</p>
        </div>

        <!-- Zakat Terkumpul -->
        <div v-if="isSuperAdmin || isBendahara" class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Zakat Terkumpul</p>
            <div class="rounded-xl bg-teal-50 p-2 text-teal-600"><HeartHandshake class="h-5 w-5" /></div>
          </div>
          <p class="mt-4 text-2xl font-bold text-gray-900">{{ formatCurrency(metrics?.zakat?.collected) }}</p>
          <p class="mt-1 text-xs text-gray-400">{{ metrics?.zakat?.mustahik ?? 0 }} mustahik terdaftar</p>
        </div>

        <!-- Aset Bermasalah -->
        <div v-if="isSuperAdmin || isTakmir" class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-gray-500">Aset Bermasalah</p>
            <div class="rounded-xl bg-red-50 p-2 text-red-500"><AlertTriangle class="h-5 w-5" /></div>
          </div>
          <p class="mt-4 text-2xl font-bold text-red-600">{{ metrics?.logistics?.broken_assets ?? 0 }}</p>
          <p class="mt-1 flex items-center gap-1 text-xs font-medium text-red-400">
            <ArrowDownRight class="h-3 w-3" /> perlu perbaikan / rusak
          </p>
        </div>
      </div>

      <!-- ── Finance Section: SUPER_ADMIN & BENDAHARA ──────────────────────── -->
      <div v-if="isSuperAdmin || isBendahara" class="space-y-4">
        <div class="flex items-center gap-2">
          <TrendingUp class="h-5 w-5 text-emerald-600" />
          <h2 class="text-lg font-bold text-gray-900">Tren Keuangan 6 Bulan Terakhir</h2>
        </div>
        <div class="grid gap-6 lg:grid-cols-3">
          <!-- Area Chart (full width on its row) -->
          <div class="lg:col-span-2 rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
            <h3 class="text-sm font-semibold text-gray-500 mb-4">Pemasukan vs Pengeluaran</h3>
            <VueApexCharts type="area" height="260" :options="financeAreaOptions" :series="financeAreaSeries" />
          </div>
          <!-- Zakat Radial Bar -->
          <div class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100 flex flex-col">
            <h3 class="text-sm font-semibold text-gray-500 mb-2">Distribusi Zakat</h3>
            <div class="flex-1 flex flex-col items-center justify-center">
              <VueApexCharts type="radialBar" height="220" :options="zakatRadialOptions" :series="zakatRadialSeries" />
              <div class="mt-2 space-y-1 w-full text-sm text-gray-600">
                <div class="flex justify-between px-4">
                  <span>Terkumpul</span>
                  <span class="font-semibold text-gray-900">{{ formatCurrency(metrics?.analytics?.zakat_breakdown?.total_donated) }}</span>
                </div>
                <div class="flex justify-between px-4">
                  <span>Tersalurkan</span>
                  <span class="font-semibold text-teal-600">{{ formatCurrency(metrics?.analytics?.zakat_breakdown?.total_distributed) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ── Logistics Section: SUPER_ADMIN & TAKMIR ───────────────────────── -->
      <div v-if="isSuperAdmin || isTakmir" class="space-y-4">
        <div class="flex items-center gap-2">
          <Package class="h-5 w-5 text-orange-500" />
          <h2 class="text-lg font-bold text-gray-900">Agenda, Inventaris & Qurban</h2>
        </div>
        <div class="grid gap-6 md:grid-cols-2">
          <!-- Qurban Donut -->
          <div class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
            <h3 class="text-sm font-semibold text-gray-500 mb-4">Komposisi Hewan Qurban</h3>
            <VueApexCharts type="donut" height="280" :options="qurbanDonutOptions" :series="qurbanDonutSeries" />
          </div>
          <!-- Asset Status Bar -->
          <div class="rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
            <h3 class="text-sm font-semibold text-gray-500 mb-4">Status Aset Inventaris</h3>
            <VueApexCharts type="bar" height="280" :options="assetBarOptions" :series="assetBarSeries" />
          </div>
        </div>
      </div>

      <!-- ── Jamaah Section: SUPER_ADMIN & SEKRETARIS ──────────────────────── -->
      <div v-if="isSuperAdmin || isSekretaris" class="space-y-4">
        <div class="flex items-center gap-2">
          <Users class="h-5 w-5 text-indigo-500" />
          <h2 class="text-lg font-bold text-gray-900">Pertumbuhan Jamaah</h2>
        </div>
        <div class="grid gap-6 lg:grid-cols-3">
          <!-- Jamaah Growth Bar -->
          <div class="lg:col-span-2 rounded-2xl bg-white p-6 shadow-sm border border-gray-100">
            <h3 class="text-sm font-semibold text-gray-500 mb-4">Jamaah Baru per Bulan (6 Bulan Terakhir)</h3>
            <VueApexCharts type="bar" height="240" :options="jamaahBarOptions" :series="jamaahBarSeries" />
          </div>
          <!-- Stats sidebar -->
          <div class="rounded-2xl bg-gradient-to-br from-indigo-600 to-violet-700 p-6 text-white shadow-sm flex flex-col justify-between">
            <div>
              <p class="text-sm font-medium text-indigo-200">Total Jamaah</p>
              <p class="mt-2 text-5xl font-bold">{{ formatNumber(metrics?.jamaah?.total) }}</p>
              <p class="mt-2 text-sm text-indigo-200">anggota komunitas terdaftar</p>
            </div>
            <div class="mt-6 space-y-3">
              <div class="flex items-center justify-between rounded-xl bg-white/10 px-4 py-2.5">
                <span class="text-sm text-indigo-100">Mustahik</span>
                <span class="font-bold">{{ metrics?.zakat?.mustahik ?? 0 }} KK</span>
              </div>
              <div class="flex items-center justify-between rounded-xl bg-white/10 px-4 py-2.5">
                <span class="text-sm text-indigo-100">Baru bulan ini</span>
                <span class="font-bold">+{{ metrics?.jamaah?.newly_registered ?? 0 }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

    </template>
  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(12px); }
  to   { opacity: 1; transform: translateY(0); }
}
.animate-in { animation: fadeIn 0.5s ease-out forwards; }
</style>
