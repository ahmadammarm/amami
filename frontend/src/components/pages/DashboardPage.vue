<script setup lang="ts">
import { computed } from 'vue';
import VueApexCharts from 'vue3-apexcharts';
import type { ApexOptions } from 'apexcharts';
import { useQuery } from '@tanstack/vue-query';
import {
  Wallet,
  Users,
  HandCoins,
  CalendarRange,
  PackageSearch,
  TriangleAlert,
} from 'lucide-vue-next';
import { dashboardService } from '@/api/services/dashboard';
import { queryKeys } from '@/api/queryKeys';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { formatNumber, formatRupiah } from '@/lib/format';
import StatCard from '@/components/molecules/StatCard.vue';
import DataPanel from '@/components/molecules/DataPanel.vue';

// Chart colours mirror the design tokens in style.css.
const COLOR = {
  sajadah: '#2d6a57',
  kuningan: '#b08a3e',
  bata: '#a8412f',
  jatiMuda: '#6b6158',
  batu: '#e4e7e1',
  sajadahMuda: '#e6efeb',
} as const;

const authStore = useAuthStore();

const { data: metrics, isLoading } = useQuery({
  queryKey: queryKeys.dashboard.metrics,
  queryFn: dashboardService.getMetrics,
  refetchInterval: 60000,
});

const firstName = computed(() => (authStore.user?.full_name || authStore.user?.username || '').split(' ')[0]);

const greeting = computed(() => {
  const hour = new Date().getHours();
  if (hour < 11) return 'Selamat pagi';
  if (hour < 15) return 'Selamat siang';
  if (hour < 18) return 'Selamat sore';
  return 'Selamat malam';
});

const today = new Date();
const gregorianDate = today.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
const hijriDate = today.toLocaleDateString('id-ID-u-ca-islamic-umalqura', { day: 'numeric', month: 'long', year: 'numeric' });

const show = computed(() => ({
  finance: authStore.can(Permission.LedgerRead),
  zakat: authStore.can(Permission.ZakatRead),
  jamaah: authStore.can(Permission.JamaahManage),
  agenda: authStore.can(Permission.AgendaRead),
  inventory: authStore.can(Permission.InventoryRead),
  qurban: authStore.can(Permission.QurbanRead),
}));

const analytics = computed(() => metrics.value?.analytics);
const brokenAssets = computed(() => metrics.value?.logistics?.broken_assets ?? 0);

const axisLabel = { style: { colors: COLOR.jatiMuda, fontSize: '12px' } };
const baseChart = { toolbar: { show: false }, fontFamily: 'inherit', animations: { enabled: true, speed: 500 } };

// Finance: income vs expense over six months
const financeSeries = computed(() => {
  const points = analytics.value?.monthly_finance ?? [];
  return [
    { name: 'Pemasukan', data: points.map(p => p.income) },
    { name: 'Pengeluaran', data: points.map(p => p.expense) },
  ];
});
const financeOptions = computed<ApexOptions>(() => ({
  chart: { ...baseChart, type: 'area' },
  colors: [COLOR.sajadah, COLOR.bata],
  fill: { type: 'gradient', gradient: { opacityFrom: 0.28, opacityTo: 0.02 } },
  stroke: { curve: 'smooth', width: 2 },
  dataLabels: { enabled: false },
  xaxis: {
    categories: (analytics.value?.monthly_finance ?? []).map(p => p.month),
    labels: axisLabel,
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: { labels: { ...axisLabel, formatter: (v: number) => formatRupiah(v) } },
  grid: { borderColor: COLOR.batu, strokeDashArray: 4 },
  tooltip: { y: { formatter: (v: number) => formatRupiah(v) } },
  legend: { position: 'top', horizontalAlign: 'left', fontSize: '13px', labels: { colors: COLOR.jatiMuda } },
}));

// Zakat: share of collected zakat already distributed
const zakatDonated = computed(() => analytics.value?.zakat_breakdown?.total_donated ?? 0);
const zakatDistributed = computed(() => analytics.value?.zakat_breakdown?.total_distributed ?? 0);
const zakatSeries = computed(() => {
  if (zakatDonated.value === 0) return [0];
  return [Math.min(Math.round((zakatDistributed.value / zakatDonated.value) * 100), 100)];
});
const zakatOptions = computed<ApexOptions>(() => ({
  chart: { ...baseChart, type: 'radialBar' },
  colors: [COLOR.kuningan],
  plotOptions: {
    radialBar: {
      hollow: { size: '62%' },
      track: { background: COLOR.sajadahMuda },
      dataLabels: {
        name: { show: true, color: COLOR.jatiMuda, fontSize: '12px', offsetY: -6 },
        value: { show: true, color: '#2b211b', fontSize: '26px', fontWeight: 600, formatter: (v: number) => `${v}%` },
      },
    },
  },
  labels: ['Tersalurkan'],
}));

// Qurban: animals per type
const QURBAN_LABELS = ['Sapi', 'Kambing', 'Domba'];
const qurbanSeries = computed(() => [
  analytics.value?.qurban_breakdown?.sapi ?? 0,
  analytics.value?.qurban_breakdown?.kambing ?? 0,
  analytics.value?.qurban_breakdown?.domba ?? 0,
]);
const qurbanOptions = computed<ApexOptions>(() => ({
  chart: { ...baseChart, type: 'donut' },
  labels: QURBAN_LABELS,
  colors: [COLOR.sajadah, COLOR.kuningan, COLOR.jatiMuda],
  dataLabels: { enabled: false },
  legend: { position: 'bottom', fontSize: '13px', labels: { colors: COLOR.jatiMuda } },
  stroke: { width: 2, colors: ['#ffffff'] },
  plotOptions: {
    pie: {
      donut: {
        size: '64%',
        labels: {
          show: true,
          total: {
            show: true,
            label: 'Total hewan',
            color: COLOR.jatiMuda,
            formatter: () => formatNumber(qurbanSeries.value.reduce((a, b) => a + b, 0)),
          },
        },
      },
    },
  },
  tooltip: { y: { formatter: (v: number) => `${v} ekor` } },
}));

// Inventory: assets per condition
const assetSeries = computed(() => [{
  name: 'Aset',
  data: [
    analytics.value?.asset_breakdown?.good ?? 0,
    analytics.value?.asset_breakdown?.loaned ?? 0,
    analytics.value?.asset_breakdown?.needs_repair ?? 0,
    analytics.value?.asset_breakdown?.broken ?? 0,
  ],
}]);
const assetOptions = computed<ApexOptions>(() => ({
  chart: { ...baseChart, type: 'bar' },
  colors: [COLOR.sajadah, COLOR.kuningan, COLOR.jatiMuda, COLOR.bata],
  plotOptions: { bar: { borderRadius: 6, columnWidth: '48%', distributed: true } },
  dataLabels: { enabled: false },
  legend: { show: false },
  xaxis: {
    categories: ['Baik', 'Dipinjam', 'Perlu perbaikan', 'Rusak'],
    labels: axisLabel,
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: { labels: axisLabel, tickAmount: 4 },
  grid: { borderColor: COLOR.batu, strokeDashArray: 4 },
  tooltip: { y: { formatter: (v: number) => `${v} unit` } },
}));

// Jamaah: new registrations per month
const jamaahSeries = computed(() => [{
  name: 'Jamaah baru',
  data: (analytics.value?.jamaah_growth ?? []).map(p => p.new_members),
}]);
const jamaahOptions = computed<ApexOptions>(() => ({
  chart: { ...baseChart, type: 'bar' },
  colors: [COLOR.sajadah],
  plotOptions: { bar: { borderRadius: 6, columnWidth: '42%' } },
  dataLabels: { enabled: false },
  xaxis: {
    categories: (analytics.value?.jamaah_growth ?? []).map(p => p.month),
    labels: axisLabel,
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: { labels: axisLabel, tickAmount: 4 },
  grid: { borderColor: COLOR.batu, strokeDashArray: 4 },
  tooltip: { y: { formatter: (v: number) => `${v} orang` } },
}));
</script>

<template>
  <div class="space-y-8">
    <header class="border-b border-batu pb-6">
      <p class="text-sm text-jati-muda">
        {{ gregorianDate }} <span class="text-kuningan">·</span> {{ hijriDate }}
      </p>
      <h1 class="mt-2 font-display text-3xl font-semibold tracking-wide text-jati sm:text-4xl">
        {{ greeting }}, {{ firstName }}
      </h1>
      <p class="mt-1.5 text-sm text-jati-muda">Ringkasan kegiatan dan catatan masjid hari ini.</p>
    </header>

    <div v-if="isLoading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="i in 6" :key="i" class="h-[118px] animate-pulse rounded-xl bg-batu/60" />
    </div>

    <template v-else-if="metrics">
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard
          v-if="show.finance"
          label="Saldo kas"
          :value="formatRupiah(metrics.finance.total_balance)"
          :hint="`${formatRupiah(metrics.finance.monthly_inflow)} masuk bulan ini`"
          :icon="Wallet"
        />
        <StatCard
          v-if="show.zakat"
          label="Zakat terkumpul"
          :value="formatRupiah(metrics.zakat.collected)"
          :hint="`${formatNumber(metrics.zakat.mustahik)} mustahik terdaftar`"
          :icon="HandCoins"
        />
        <StatCard
          v-if="show.jamaah"
          label="Jamaah"
          :value="formatNumber(metrics.jamaah.total)"
          :hint="`${formatNumber(metrics.jamaah.newly_registered)} baru bulan ini`"
          :icon="Users"
        />
        <StatCard
          v-if="show.agenda"
          label="Agenda terjadwal"
          :value="formatNumber(metrics.logistics.upcoming_agendas)"
          hint="kegiatan belum berlangsung"
          :icon="CalendarRange"
        />
        <StatCard
          v-if="show.inventory"
          label="Aset dipinjam"
          :value="formatNumber(metrics.logistics.active_loans)"
          hint="belum dikembalikan"
          :icon="PackageSearch"
        />
        <StatCard
          v-if="show.inventory"
          label="Aset bermasalah"
          :value="formatNumber(brokenAssets)"
          hint="perlu perbaikan atau rusak"
          :icon="TriangleAlert"
          :tone="brokenAssets > 0 ? 'alert' : 'default'"
        />
      </div>

      <div v-if="show.finance || show.zakat" class="grid gap-6 lg:grid-cols-3">
        <DataPanel
          v-if="show.finance"
          title="Arus kas enam bulan terakhir"
          description="Pemasukan dan pengeluaran semua kas."
          class="lg:col-span-2"
        >
          <div class="px-3 pt-2 pb-3">
            <VueApexCharts type="area" height="280" :options="financeOptions" :series="financeSeries" />
          </div>
        </DataPanel>
        <DataPanel v-if="show.zakat" title="Penyaluran zakat" description="Porsi zakat terkumpul yang sudah disalurkan.">
          <div class="px-3 pt-2">
            <VueApexCharts type="radialBar" height="230" :options="zakatOptions" :series="zakatSeries" />
          </div>
          <dl class="shaf mx-5 mb-5 text-sm">
            <div class="flex justify-between py-2">
              <dt class="text-jati-muda">Terkumpul</dt>
              <dd class="tabular font-semibold">{{ formatRupiah(zakatDonated) }}</dd>
            </div>
            <div class="flex justify-between py-2">
              <dt class="text-jati-muda">Tersalurkan</dt>
              <dd class="tabular font-semibold text-sajadah">{{ formatRupiah(zakatDistributed) }}</dd>
            </div>
          </dl>
        </DataPanel>
      </div>

      <div v-if="show.qurban || show.inventory" class="grid gap-6 md:grid-cols-2">
        <DataPanel v-if="show.qurban" title="Hewan qurban" description="Jumlah hewan menurut jenisnya.">
          <div class="px-3 py-3">
            <VueApexCharts type="donut" height="280" :options="qurbanOptions" :series="qurbanSeries" />
          </div>
        </DataPanel>
        <DataPanel v-if="show.inventory" title="Kondisi aset" description="Sebaran aset menurut kondisinya.">
          <div class="px-3 py-3">
            <VueApexCharts type="bar" height="280" :options="assetOptions" :series="assetSeries" />
          </div>
        </DataPanel>
      </div>

      <DataPanel v-if="show.jamaah" title="Pertumbuhan jamaah" description="Pendaftaran jamaah baru per bulan.">
        <div class="px-3 py-3">
          <VueApexCharts type="bar" height="240" :options="jamaahOptions" :series="jamaahSeries" />
        </div>
      </DataPanel>
    </template>
  </div>
</template>
