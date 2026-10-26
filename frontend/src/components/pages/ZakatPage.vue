<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { zakatService } from '@/api/services/zakat';
import { financeService } from '@/api/services/finance';
import { feedback } from '@/composables/useFeedback';
import { queryKeys } from '@/api/queryKeys';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { collectZakatSchema, distributeZakatSchema } from '@/schemas/zakat';
import { HeartHandshake, ArrowDown, ArrowUp } from 'lucide-vue-next';
import AsyncJamaahSelect from '../molecules/AsyncJamaahSelect.vue';

const authStore = useAuthStore();
const queryClient = useQueryClient();

const PAGE_SIZE = 50;

// State
const activeTab = ref<'donations' | 'distributions'>('donations');
const showCollectDialog = ref(false);
const showDistributeDialog = ref(false);

// Collect form
const {
  handleSubmit: handleCollectSubmit,
  errors: collectErrors,
  defineField: defineCollectField,
  setFieldValue: setCollectField,
  resetForm: resetCollectForm,
} = useForm({
  validationSchema: toTypedSchema(collectZakatSchema),
  initialValues: { muzakki_id: '', zakat_type: 'MAAL', fund_id: null, amount: 0, unit: 'IDR', description: '' },
});

const [muzakkiId] = defineCollectField('muzakki_id');
const [zakatType] = defineCollectField('zakat_type');
const [collectFundId] = defineCollectField('fund_id');
const [collectAmount] = defineCollectField('amount');
const [collectUnit] = defineCollectField('unit');
const [collectDescription] = defineCollectField('description');

// Distribute form
const {
  handleSubmit: handleDistributeSubmit,
  errors: distributeErrors,
  defineField: defineDistributeField,
  setFieldValue: setDistributeField,
  resetForm: resetDistributeForm,
} = useForm({
  validationSchema: toTypedSchema(distributeZakatSchema),
  initialValues: { mustahik_id: '', fund_id: null, amount: 0, unit: 'IDR', description: '' },
});

const [mustahikId] = defineDistributeField('mustahik_id');
const [distributeFundId] = defineDistributeField('fund_id');
const [distributeAmount] = defineDistributeField('amount');
const [distributeUnit] = defineDistributeField('unit');
const [distributeDescription] = defineDistributeField('description');

// Rice is counted in KG and never posted to a fund; cash is always IDR.
watch(zakatType, (type) => {
  const isRice = type === 'FITRAH_BERAS';
  setCollectField('unit', isRice ? 'KG' : 'IDR');
  if (isRice) setCollectField('fund_id', null);
});

watch(distributeUnit, (unit) => {
  if (unit !== 'IDR') setDistributeField('fund_id', null);
});

// Permissions
const canManage = computed(() => authStore.can(Permission.ZakatManage));
const canReadFunds = computed(() => authStore.can(Permission.LedgerRead));

// Queries
const { data: donations } = useQuery({
  queryKey: queryKeys.zakat.donations(1, PAGE_SIZE),
  queryFn: () => zakatService.getDonations(1, PAGE_SIZE),
});

const { data: distributions } = useQuery({
  queryKey: queryKeys.zakat.distributions(1, PAGE_SIZE),
  queryFn: () => zakatService.getDistributions(1, PAGE_SIZE),
});

const { data: funds } = useQuery({
  queryKey: queryKeys.finance.funds,
  queryFn: financeService.getFunds,
  enabled: canReadFunds,
});

const donationItems = computed(() => donations.value?.items ?? []);
const distributionItems = computed(() => distributions.value?.items ?? []);

const refreshAfterZakat = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.zakat.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.finance.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Mutations
const collectMutation = useMutation({
  mutationFn: zakatService.collectZakat,
  onSuccess: () => {
    feedback.success('Penerimaan zakat tercatat');
    refreshAfterZakat();
    closeCollectDialog();
  },
  onError: (error: unknown) => feedback.error('Penerimaan zakat belum tercatat', error),
});

const distributeMutation = useMutation({
  mutationFn: zakatService.distributeZakat,
  onSuccess: () => {
    feedback.success('Penyaluran zakat tercatat');
    refreshAfterZakat();
    closeDistributeDialog();
  },
  onError: (error: unknown) => feedback.error('Penyaluran zakat belum tercatat', error),
});

const submitCollect = handleCollectSubmit((values) => collectMutation.mutate(values));
const submitDistribute = handleDistributeSubmit((values) => distributeMutation.mutate(values));

const closeCollectDialog = () => {
  showCollectDialog.value = false;
  resetCollectForm();
};

const closeDistributeDialog = () => {
  showDistributeDialog.value = false;
  resetDistributeForm();
};

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
          <tr v-if="donationItems.length === 0" class="border-b">
            <td colspan="5" class="px-6 py-8 text-center text-gray-500">Belum ada data penerimaan zakat.</td>
          </tr>
          <tr v-else v-for="d in donationItems" :key="d.id" class="border-b hover:bg-gray-50">
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
          <tr v-if="distributionItems.length === 0" class="border-b">
            <td colspan="4" class="px-6 py-8 text-center text-gray-500">Belum ada data penyaluran zakat.</td>
          </tr>
          <tr v-else v-for="d in distributionItems" :key="d.id" class="border-b hover:bg-gray-50">
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
          <button @click="closeCollectDialog" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <form @submit="submitCollect" novalidate>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Muzakki (Donatur)</label>
              <AsyncJamaahSelect v-model="muzakkiId" :muzakki-only="true" placeholder="Cari Muzakki (Donatur)..." />
              <p v-if="collectErrors.muzakki_id" class="mt-1 text-xs text-red-600">{{ collectErrors.muzakki_id }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Jenis Zakat</label>
              <select v-model="zakatType" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="FITRAH_BERAS">Zakat Fitrah (Beras)</option>
                <option value="MAAL">Zakat Maal</option>
                <option value="INFAQ">Infaq / Sedekah</option>
              </select>
            </div>
            <div v-if="collectUnit === 'IDR'">
              <label class="block text-sm font-medium text-gray-700 mb-1">Rekening Kas Tujuan</label>
              <select v-model="collectFundId" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option :value="null" disabled>Pilih Kas</option>
                <option v-for="f in funds" :key="f.id" :value="f.id">{{ f.name }} ({{ formatCurrency(f.current_balance) }})</option>
              </select>
              <p v-if="collectErrors.fund_id" class="mt-1 text-xs text-red-600">{{ collectErrors.fund_id }}</p>
            </div>
            <div class="flex gap-4">
              <div class="flex-1">
                <label class="block text-sm font-medium text-gray-700 mb-1">Jumlah</label>
                <input type="number" v-model.number="collectAmount" min="0" :step="collectUnit === 'IDR' ? 1 : 0.01" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
                <p v-if="collectErrors.amount" class="mt-1 text-xs text-red-600">{{ collectErrors.amount }}</p>
              </div>
              <div class="w-1/3">
                <label class="block text-sm font-medium text-gray-700 mb-1">Satuan</label>
                <select v-model="collectUnit" disabled class="block w-full rounded-lg border-gray-300 bg-gray-50 p-2.5 border text-gray-500">
                  <option value="IDR">Rupiah (Rp)</option>
                  <option value="KG">Kilogram (KG)</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Keterangan</label>
              <textarea v-model="collectDescription" rows="2" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Catatan tambahan..."></textarea>
              <p v-if="collectErrors.description" class="mt-1 text-xs text-red-600">{{ collectErrors.description }}</p>
            </div>
          </div>
          <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
            <button type="button" @click="closeCollectDialog" class="px-4 py-2 text-sm text-gray-700 border rounded-lg hover:bg-gray-50">Batal</button>
            <button type="submit" :disabled="collectMutation.isPending.value" class="px-4 py-2 text-sm text-white bg-green-600 rounded-lg hover:bg-green-700 disabled:opacity-50">Simpan Penerimaan</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modals (Distribute Zakat) -->
    <div v-if="showDistributeDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
          <h3 class="text-lg font-bold text-gray-900">Salurkan Zakat</h3>
          <button @click="closeDistributeDialog" class="text-gray-400 hover:text-gray-600">&times;</button>
        </div>
        <form @submit="submitDistribute" novalidate>
          <div class="p-6 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Mustahik (Penerima)</label>
              <AsyncJamaahSelect v-model="mustahikId" :mustahik-only="true" placeholder="Cari Mustahik..." />
              <p v-if="distributeErrors.mustahik_id" class="mt-1 text-xs text-red-600">{{ distributeErrors.mustahik_id }}</p>
              <p v-else class="text-xs text-gray-500 mt-1">Hanya Mustahik terdaftar, diurutkan dari skor kemiskinan tertinggi.</p>
            </div>
            <div class="flex items-center gap-4 border-b border-gray-100 pb-4">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="distributeUnit" value="IDR" class="text-blue-600 focus:ring-blue-500" />
                <span class="text-sm font-medium text-gray-700">Uang Tunai</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" v-model="distributeUnit" value="KG" class="text-blue-600 focus:ring-blue-500" />
                <span class="text-sm font-medium text-gray-700">Beras (KG)</span>
              </label>
            </div>
            <div v-if="distributeUnit === 'IDR'">
              <label class="block text-sm font-medium text-gray-700 mb-1">Rekening Sumber Kas</label>
              <select v-model="distributeFundId" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option :value="null" disabled>Pilih Kas</option>
                <option v-for="f in funds" :key="f.id" :value="f.id">{{ f.name }} (Tersedia: {{ formatCurrency(f.current_balance) }})</option>
              </select>
              <p v-if="distributeErrors.fund_id" class="mt-1 text-xs text-red-600">{{ distributeErrors.fund_id }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Jumlah</label>
              <div class="flex">
                <input type="number" v-model.number="distributeAmount" min="0" :step="distributeUnit === 'IDR' ? 1 : 0.01" class="block w-full rounded-l-lg border-gray-300 p-2.5 border focus:ring-blue-500 focus:border-blue-500" />
                <span class="inline-flex items-center px-4 rounded-r-lg border border-l-0 border-gray-300 bg-gray-50 text-gray-500 text-sm font-medium">
                  {{ distributeUnit === 'IDR' ? 'Rp' : 'KG' }}
                </span>
              </div>
              <p v-if="distributeErrors.amount" class="mt-1 text-xs text-red-600">{{ distributeErrors.amount }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Keterangan</label>
              <textarea v-model="distributeDescription" rows="2" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Catatan tambahan..."></textarea>
              <p v-if="distributeErrors.description" class="mt-1 text-xs text-red-600">{{ distributeErrors.description }}</p>
            </div>
          </div>
          <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
            <button type="button" @click="closeDistributeDialog" class="px-4 py-2 text-sm text-gray-700 border rounded-lg hover:bg-gray-50">Batal</button>
            <button type="submit" :disabled="distributeMutation.isPending.value" class="px-4 py-2 text-sm text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50">Salurkan Bantuan</button>
          </div>
        </form>
      </div>
    </div>

  </div>
</template>
