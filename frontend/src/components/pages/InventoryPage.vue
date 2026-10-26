<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { inventoryService } from '../../api/services/inventory';
import { getErrorMessage } from '../../api/client';
import { queryKeys } from '../../api/queryKeys';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import { assetSchema, toAssetPayload, type AssetFormValues, type AssetPayload } from '../../schemas/inventory';
import type { Asset, AssetStatus } from '../../types/inventory';
import { PlusIcon, XIcon, TrashIcon, BoxIcon, EditIcon, CheckCircleIcon, WrenchIcon, AlertTriangleIcon, RepeatIcon } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

const authStore = useAuthStore();
const queryClient = useQueryClient();

const canManage = computed(() => authStore.can(Permission.InventoryManage));

const showAddDialog = ref(false);
const editingAsset = ref<Asset | null>(null);

const emptyForm: AssetFormValues = {
  name: '',
  sku: '',
  purchase_date: null,
  purchase_price: 0,
  current_status: 'GOOD',
  location: '',
};

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: toTypedSchema(assetSchema),
  initialValues: emptyForm,
});

const [name] = defineField('name');
const [sku] = defineField('sku');
const [purchaseDate] = defineField('purchase_date');
const [purchasePrice] = defineField('purchase_price');
const [currentStatus] = defineField('current_status');
const [location] = defineField('location');

/** LOANED is set by the loan workflow; it is only shown for an asset that is already on loan. */
const isOnLoan = computed(() => editingAsset.value?.current_status === 'LOANED');

const { data: assets, isLoading } = useQuery({
  queryKey: queryKeys.inventory.assets,
  queryFn: () => inventoryService.getAssets(1, 100)
});

const assetItems = computed(() => assets.value?.items ?? []);
const countByStatus = (...statuses: AssetStatus[]) =>
  assetItems.value.filter(a => statuses.includes(a.current_status)).length;

const totalAssets = computed(() => assets.value?.total ?? 0);
const goodAssets = computed(() => countByStatus('GOOD'));
const loanedAssets = computed(() => countByStatus('LOANED'));
const repairAssets = computed(() => countByStatus('REPAIR', 'NEEDS_REPAIR'));
const brokenAssets = computed(() => countByStatus('BROKEN'));

const refreshAssets = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

const createMutation = useMutation({
  mutationFn: inventoryService.createAsset,
  onSuccess: () => {
    toast.success('Aset berhasil ditambahkan');
    refreshAssets();
    closeDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal menambahkan aset'))
});

const updateMutation = useMutation({
  mutationFn: ({ id, payload }: { id: string, payload: AssetPayload }) => inventoryService.updateAsset(id, payload),
  onSuccess: () => {
    toast.success('Aset berhasil diperbarui');
    refreshAssets();
    closeDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal memperbarui aset'))
});

const deleteMutation = useMutation({
  mutationFn: inventoryService.deleteAsset,
  onSuccess: () => {
    toast.success('Aset berhasil dihapus');
    refreshAssets();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal menghapus aset'))
});

const isSaving = computed(() => createMutation.isPending.value || updateMutation.isPending.value);

const openCreateDialog = () => {
  editingAsset.value = null;
  resetForm({ values: emptyForm });
  showAddDialog.value = true;
};

const openEditDialog = (asset: Asset) => {
  editingAsset.value = asset;
  resetForm({
    values: {
      name: asset.name,
      sku: asset.sku,
      purchase_date: asset.purchase_date ? asset.purchase_date.slice(0, 10) : null,
      purchase_price: asset.purchase_price,
      current_status: asset.current_status,
      location: asset.location,
    },
  });
  showAddDialog.value = true;
};

const closeDialog = () => {
  showAddDialog.value = false;
  editingAsset.value = null;
};

const submitAsset = handleSubmit((values) => {
  const payload = toAssetPayload(values);
  if (editingAsset.value) {
    updateMutation.mutate({ id: editingAsset.value.id, payload });
  } else {
    createMutation.mutate(payload);
  }
});

const confirmDelete = (asset: Asset) => {
  if (window.confirm(`Hapus aset "${asset.name}" (${asset.sku})?`)) {
    deleteMutation.mutate(asset.id);
  }
};

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(amount);
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Inventaris & Aset</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola daftar aset masjid dan peminjaman</p>
      </div>
      <button v-if="canManage" @click="openCreateDialog" class="px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 focus:ring-4 focus:ring-emerald-100 transition-colors flex items-center gap-2">
        <PlusIcon class="w-5 h-5" />
        Tambah Aset
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-blue-50 text-blue-600 rounded-xl">
            <BoxIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Total Aset</p>
            <p class="text-2xl font-bold text-gray-900">{{ totalAssets }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-emerald-50 text-emerald-600 rounded-xl">
            <CheckCircleIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Kondisi Baik</p>
            <p class="text-2xl font-bold text-gray-900">{{ goodAssets }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-purple-50 text-purple-600 rounded-xl">
            <RepeatIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Dipinjam</p>
            <p class="text-2xl font-bold text-gray-900">{{ loanedAssets }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-amber-50 text-amber-600 rounded-xl">
            <WrenchIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Perlu Perbaikan</p>
            <p class="text-2xl font-bold text-gray-900">{{ repairAssets }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-red-50 text-red-600 rounded-xl">
            <AlertTriangleIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Rusak</p>
            <p class="text-2xl font-bold text-gray-900">{{ brokenAssets }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Asset List -->
    <div class="bg-white shadow-sm rounded-xl border border-gray-100 overflow-hidden">
      <div class="p-6 border-b border-gray-100">
        <h3 class="font-semibold text-gray-900">Daftar Aset</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <thead class="text-xs text-gray-500 bg-gray-50 uppercase">
            <tr>
              <th class="px-6 py-4 font-medium">Nama Aset / SKU</th>
              <th class="px-6 py-4 font-medium">Lokasi</th>
              <th class="px-6 py-4 font-medium">Status</th>
              <th class="px-6 py-4 font-medium">Harga Beli</th>
              <th v-if="canManage" class="px-6 py-4 font-medium text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-if="isLoading" class="animate-pulse">
              <td colspan="5" class="px-6 py-4 text-center text-gray-500">Memuat data...</td>
            </tr>
            <tr v-else-if="assetItems.length === 0">
              <td colspan="5" class="px-6 py-12 text-center text-gray-500">
                <BoxIcon class="w-12 h-12 mx-auto text-gray-300 mb-3" />
                Belum ada data aset
              </td>
            </tr>
            <tr v-else v-for="asset in assetItems" :key="asset.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4">
                <div class="font-medium text-gray-900">{{ asset.name }}</div>
                <div class="text-xs text-gray-500">{{ asset.sku }}</div>
              </td>
              <td class="px-6 py-4 text-gray-600">{{ asset.location }}</td>
              <td class="px-6 py-4">
                <span :class="[
                  'px-2.5 py-1 text-xs font-medium rounded-full',
                  asset.current_status === 'GOOD' ? 'bg-emerald-50 text-emerald-700' :
                  asset.current_status === 'LOANED' ? 'bg-purple-50 text-purple-700' :
                  (asset.current_status === 'REPAIR' || asset.current_status === 'NEEDS_REPAIR') ? 'bg-amber-50 text-amber-700' :
                  'bg-red-50 text-red-700'
                ]">
                  {{ asset.current_status === 'NEEDS_REPAIR' ? 'NEEDS REPAIR' : asset.current_status }}
                </span>
              </td>
              <td class="px-6 py-4 text-gray-600">{{ formatCurrency(asset.purchase_price) }}</td>
              <td v-if="canManage" class="px-6 py-4 text-right">
                <button @click="openEditDialog(asset)" class="text-blue-600 hover:text-blue-900 p-2 hover:bg-blue-50 rounded-lg transition-colors mr-2" title="Edit">
                  <EditIcon class="w-4 h-4" />
                </button>
                <button
                  @click="confirmDelete(asset)"
                  :disabled="asset.current_status === 'LOANED' || deleteMutation.isPending.value"
                  :title="asset.current_status === 'LOANED' ? 'Aset sedang dipinjam' : 'Hapus'"
                  class="text-red-600 hover:text-red-900 p-2 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add Asset Dialog -->
    <div v-if="showAddDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-bold text-gray-900">{{ editingAsset ? 'Edit Aset' : 'Tambah Aset Baru' }}</h3>
          <button @click="closeDialog" class="text-gray-400 hover:text-gray-600">
            <XIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit="submitAsset" class="p-6 space-y-4" novalidate>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Aset</label>
            <input v-model="name" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Cth: Sound System Yamaha" />
            <p v-if="errors.name" class="mt-1 text-xs text-red-600">{{ errors.name }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">SKU / Kode Aset</label>
            <input v-model="sku" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Cth: AST-001" />
            <p v-if="errors.sku" class="mt-1 text-xs text-red-600">{{ errors.sku }}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi</label>
              <input v-model="location" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Gudang Utama" />
              <p v-if="errors.location" class="mt-1 text-xs text-red-600">{{ errors.location }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select v-model="currentStatus" :disabled="isOnLoan" class="block w-full rounded-lg border-gray-300 p-2.5 border disabled:bg-gray-50 disabled:text-gray-500">
                <option value="GOOD">Baik (Tersedia)</option>
                <option v-if="isOnLoan" value="LOANED">Dipinjam</option>
                <option value="NEEDS_REPAIR">Perlu Perbaikan</option>
                <option value="REPAIR">Sedang Diperbaiki</option>
                <option value="BROKEN">Rusak Total</option>
              </select>
              <p v-if="isOnLoan" class="mt-1 text-xs text-gray-500">Status berubah otomatis saat aset dikembalikan.</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Harga Beli</label>
              <input v-model.number="purchasePrice" type="number" min="0" step="1" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
              <p v-if="errors.purchase_price" class="mt-1 text-xs text-red-600">{{ errors.purchase_price }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tanggal Beli</label>
              <input v-model="purchaseDate" type="date" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
            </div>
          </div>

          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="isSaving" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ isSaving ? 'Menyimpan...' : 'Simpan Aset' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
