<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { zakatService } from '@/api/services/zakat';
import { jamaahService } from '@/api/services/jamaah';
import { financeService } from '@/api/services/finance';
import { useAuthStore } from '@/stores/auth';
import type { CollectZakatPayload, DistributeZakatPayload } from '@/types/zakat';
import { toast } from 'vue-sonner';
import { HeartHandshake, ArrowDown, ArrowUp } from 'lucide-vue-next';

const authStore = useAuthStore();
const queryClient = useQueryClient();

// State
const activeTab = ref<'donations' | 'distributions'>('donations');
const showCollectDialog = ref(false);
const showDistributeDialog = ref(false);

const collectForm = ref<CollectZakatPayload>({
  muzakki_id: '',
  fund_id: 0,
  amount: 0,
  zakat_type: 'FITRAH_UANG',
  unit: 'IDR',
  description: ''
});

const distributeForm = ref<DistributeZakatPayload>({
  mustahik_id: '',
  fund_id: 0,
  amount: 0,
  unit: 'IDR',
  description: ''
});

// Watcher to dynamically set unit and fund requirement
import { watch } from 'vue';
watch(() => collectForm.value.zakat_type, (newType) => {
  if (newType === 'FITRAH_BERAS') {
    collectForm.value.unit = 'KG';
    collectForm.value.fund_id = undefined; // Nullify fund_id for non-cash
  } else {
    collectForm.value.unit = 'IDR';
    if (!collectForm.value.fund_id) collectForm.value.fund_id = 0; // Require fund
  }
});

watch(() => distributeForm.value.unit, (newUnit) => {
  if (newUnit === 'KG') {
    distributeForm.value.fund_id = undefined;
  } else {
    if (!distributeForm.value.fund_id) distributeForm.value.fund_id = 0;
  }
});

// Permissions
const canManage = computed(() => {
  return ['SUPER_ADMIN', 'BENDAHARA'].includes(authStore.userRole);
});

// Queries
const { data: donations } = useQuery({
  queryKey: ['zakat-donations'],
  queryFn: () => zakatService.getDonations(1, 50),
});

const { data: distributions } = useQuery({
  queryKey: ['zakat-distributions'],
  queryFn: () => zakatService.getDistributions(1, 50),
});

const { data: jamaahs } = useQuery({
  queryKey: ['jamaahs-all'],
  queryFn: () => jamaahService.getJamaahs(1, 100, ''),
});

const { data: funds } = useQuery({
  queryKey: ['funds'],
  queryFn: financeService.getFunds,
});

const muzakkis = computed(() => jamaahs.value?.data || []);
const mustahiks = computed(() => (jamaahs.value?.data || []).filter(j => j.is_mustahik));

// Mutations
const collectMutation = useMutation({
  mutationFn: zakatService.collectZakat,
  onSuccess: () => {
    toast.success('Zakat collected successfully!');
    queryClient.invalidateQueries({ queryKey: ['zakat-donations'] });
    queryClient.invalidateQueries({ queryKey: ['funds'] }); // Update ledger
    queryClient.invalidateQueries({ queryKey: ['transactions'] }); // Update finance history
    queryClient.invalidateQueries({ queryKey: ['dashboard-metrics'] }); // Update dashboard
    showCollectDialog.value = false;
  },
  onError: () => toast.error('Failed to collect zakat'),
});

const distributeMutation = useMutation({
  mutationFn: zakatService.distributeZakat,
  onSuccess: () => {
    toast.success('Zakat distributed successfully!');
    queryClient.invalidateQueries({ queryKey: ['zakat-distributions'] });
    queryClient.invalidateQueries({ queryKey: ['funds'] }); // Update ledger
    queryClient.invalidateQueries({ queryKey: ['transactions'] }); // Update finance history
    queryClient.invalidateQueries({ queryKey: ['dashboard-metrics'] }); // Update dashboard
    showDistributeDialog.value = false;
  },
  onError: () => toast.error('Failed to distribute zakat (Check fund balance)'),
});

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amount);
};
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
      <div class="flex items-center gap-4">
        <div class="h-12 w-12 bg-green-50 text-green-600 rounded-xl flex items-center justify-center">
          <HeartHandshake class="w-6 h-6" />
        </div>
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-gray-900">Manajemen Zakat</h1>
          <p class="text-sm text-gray-500">Penerimaan & Penyaluran Zakat Fitrah & Maal.</p>
        </div>
      </div>
      <div class="flex gap-2" v-if="canManage">
        <button @click="showCollectDialog = true" class="inline-flex items-center justify-center rounded-lg bg-green-600 text-white hover:bg-green-700 h-10 px-4 text-sm font-medium transition-colors shadow-sm">
          <ArrowDown class="w-4 h-4 mr-2" /> Terima Zakat
        </button>
        <button @click="showDistributeDialog = true" class="inline-flex items-center justify-center rounded-lg bg-blue-600 text-white hover:bg-blue-700 h-10 px-4 text-sm font-medium transition-colors shadow-sm">
          <ArrowUp class="w-4 h-4 mr-2" /> Salurkan Zakat
        </button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex space-x-1 rounded-xl bg-gray-100 p-1 w-full max-w-md">
      <button 
        @click="activeTab = 'donations'"
        :class="['flex-1 rounded-lg py-2 text-sm font-medium transition-all', activeTab === 'donations' ? 'bg-white text-gray-900 shadow' : 'text-gray-500 hover:text-gray-700']"
      >
        Penerimaan Zakat
      </button>
      <button 
        @click="activeTab = 'distributions'"
        :class="['flex-1 rounded-lg py-2 text-sm font-medium transition-all', activeTab === 'distributions' ? 'bg-white text-gray-900 shadow' : 'text-gray-500 hover:text-gray-700']"
      >
        Penyaluran Zakat
      </button>
    </div>

    <!-- Data Tables -->
    <div class="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
      <!-- Donations Table -->
      <table v-if="activeTab === 'donations'" class="w-full text-sm text-left">
        <thead class="text-xs text-gray-700 uppercase bg-gray-50 border-b">
          <tr>
            <th class="px-6 py-4">Waktu</th>
            <th class="px-6 py-4">Nama Muzakki</th>
            <th class="px-6 py-4">Jenis Zakat</th>
            <th class="px-6 py-4">Keterangan</th>
            <th class="px-6 py-4 text-right">Nominal / Jumlah</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!donations?.data?.length" class="border-b">
            <td colspan="5" class="px-6 py-8 text-center text-gray-500">Belum ada data penerimaan zakat.</td>
          </tr>
          <tr v-else v-for="d in donations?.data" :key="d.id" class="border-b hover:bg-gray-50">
            <td class="px-6 py-4 text-gray-500">{{ new Date(d.created_at).toLocaleString('id-ID') }}</td>
            <td class="px-6 py-4 font-medium text-gray-900">{{ d.muzakki_name }}</td>
            <td class="px-6 py-4">
              <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-green-100 text-green-700">
                {{ d.zakat_type.replace('_', ' ') }}
              </span>
            </td>
            <td class="px-6 py-4 text-gray-600 truncate max-w-xs">{{ d.description || '-' }}</td>
            <td class="px-6 py-4 text-right font-medium text-green-600">
              <span v-if="d.unit === 'IDR'">+ {{ formatCurrency(d.amount) }}</span>
              <span v-else>+ {{ d.amount }} {{ d.unit }}</span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Distributions Table -->
      <table v-if="activeTab === 'distributions'" class="w-full text-sm text-left">
        <thead class="text-xs text-gray-700 uppercase bg-gray-50 border-b">
          <tr>
            <th class="px-6 py-4">Waktu Penyaluran</th>
            <th class="px-6 py-4">Nama Mustahik</th>
            <th class="px-6 py-4">Keterangan</th>
            <th class="px-6 py-4 text-right">Nominal / Jumlah Bantuan</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="!distributions?.data?.length" class="border-b">
            <td colspan="4" class="px-6 py-8 text-center text-gray-500">Belum ada data penyaluran zakat.</td>
          </tr>
          <tr v-else v-for="d in distributions?.data" :key="d.id" class="border-b hover:bg-gray-50">
            <td class="px-6 py-4 text-gray-500">{{ new Date(d.distributed_at).toLocaleString('id-ID') }}</td>
            <td class="px-6 py-4 font-medium text-gray-900">{{ d.mustahik_name }}</td>
            <td class="px-6 py-4 text-gray-600 truncate max-w-xs">{{ d.description || '-' }}</td>
            <td class="px-6 py-4 text-right font-medium text-blue-600">
              <span v-if="d.unit === 'IDR'">- {{ formatCurrency(d.amount) }}</span>
              <span v-else>- {{ d.amount }} {{ d.unit }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modals (Collect Zakat) -->
    <div v-if="showCollectDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
          <h3 class="text-lg font-bold text-gray-900">Terima Zakat</h3>
          <button @click="showCollectDialog = false" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <form @submit.prevent="collectMutation.mutate(collectForm)">
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Muzakki (Donatur)</label>
              <select v-model="collectForm.muzakki_id" required class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="" disabled>Pilih Muzakki</option>
                <option v-for="m in muzakkis" :key="m.id" :value="m.id">{{ m.full_name }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Jenis Zakat</label>
              <select v-model="collectForm.zakat_type" required class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="FITRAH_UANG">Zakat Fitrah (Uang)</option>
                <option value="FITRAH_BERAS">Zakat Fitrah (Beras)</option>
                <option value="MAAL">Zakat Maal</option>
                <option value="PROFESI">Zakat Profesi</option>
                <option value="INFAQ">Infaq / Sedekah</option>
              </select>
            </div>
            <div v-if="collectForm.zakat_type !== 'FITRAH_BERAS'">
              <label class="block text-sm font-medium text-gray-700 mb-1">Rekening Kas Tujuan</label>
              <select v-model="collectForm.fund_id" required class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option :value="0" disabled>Pilih Kas</option>
                <option v-for="f in funds" :key="f.id" :value="f.id">{{ f.name }} ({{ formatCurrency(f.current_balance) }})</option>
              </select>
            </div>
            <div class="flex gap-4">
              <div class="flex-1">
                <label class="block text-sm font-medium text-gray-700 mb-1">Jumlah</label>
                <input type="number" v-model="collectForm.amount" required min="1" step="0.01" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
              </div>
              <div class="w-1/3">
                <label class="block text-sm font-medium text-gray-700 mb-1">Satuan</label>
                <select v-model="collectForm.unit" disabled class="block w-full rounded-lg border-gray-300 bg-gray-50 p-2.5 border text-gray-500">
                  <option value="IDR">Rupiah (Rp)</option>
                  <option value="KG">Kilogram (KG)</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Keterangan</label>
              <textarea v-model="collectForm.description" rows="2" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Catatan tambahan..."></textarea>
            </div>
          </div>
          <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
            <button type="button" @click="showCollectDialog = false" class="px-4 py-2 text-sm text-gray-700 border rounded-lg hover:bg-gray-50">Batal</button>
            <button type="submit" class="px-4 py-2 text-sm text-white bg-green-600 rounded-lg hover:bg-green-700">Simpan Penerimaan</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modals (Distribute Zakat) -->
    <div v-if="showDistributeDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
          <h3 class="text-lg font-bold text-gray-900">Salurkan Zakat</h3>
          <button @click="showDistributeDialog = false" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <form @submit.prevent="distributeMutation.mutate(distributeForm)">
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Mustahik (Penerima)</label>
              <select v-model="distributeForm.mustahik_id" required class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="" disabled>Pilih Mustahik</option>
                <option v-for="m in mustahiks" :key="m.id" :value="m.id">{{ m.full_name }}</option>
              </select>
              <p class="text-xs text-gray-500 mt-1">Hanya menampilkan jamaah yang ditandai sebagai Mustahik.</p>
            </div>
            <div class="flex items-center gap-4 border-b border-gray-100 pb-4">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="distributeForm.unit" value="IDR" class="text-blue-600 focus:ring-blue-500" />
                <span class="text-sm font-medium text-gray-700">Uang Tunai</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="distributeForm.unit" value="KG" class="text-blue-600 focus:ring-blue-500" />
                <span class="text-sm font-medium text-gray-700">Beras (KG)</span>
              </label>
            </div>
            <div v-if="distributeForm.unit === 'IDR'">
              <label class="block text-sm font-medium text-gray-700 mb-1">Rekening Sumber Kas</label>
              <select v-model="distributeForm.fund_id" required class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option :value="0" disabled>Pilih Kas</option>
                <option v-for="f in funds" :key="f.id" :value="f.id">{{ f.name }} (Tersedia: {{ formatCurrency(f.current_balance) }})</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Jumlah</label>
              <div class="flex">
                <input type="number" v-model="distributeForm.amount" required min="1" step="0.01" class="block w-full rounded-l-lg border-gray-300 p-2.5 border focus:ring-blue-500 focus:border-blue-500" />
                <span class="inline-flex items-center px-4 rounded-r-lg border border-l-0 border-gray-300 bg-gray-50 text-gray-500 text-sm font-medium">
                  {{ distributeForm.unit === 'IDR' ? 'Rp' : 'KG' }}
                </span>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Keterangan</label>
              <textarea v-model="distributeForm.description" rows="2" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Catatan tambahan..."></textarea>
            </div>
          </div>
          <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
            <button type="button" @click="showDistributeDialog = false" class="px-4 py-2 text-sm text-gray-700 border rounded-lg hover:bg-gray-50">Batal</button>
            <button type="submit" class="px-4 py-2 text-sm text-white bg-blue-600 rounded-lg hover:bg-blue-700">Salurkan Bantuan</button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>
