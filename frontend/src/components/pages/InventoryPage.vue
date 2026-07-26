<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { inventoryService } from '../../api/services/inventory';
import type { CreateAssetPayload, Asset } from '../../types/inventory';
import { PlusIcon, XIcon, TrashIcon, BoxIcon, EditIcon, CheckCircleIcon, WrenchIcon, AlertTriangleIcon, ActivityIcon, RepeatIcon } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

const queryClient = useQueryClient();

const showAddDialog = ref(false);
const editingId = ref<string | null>(null);

const form = ref<CreateAssetPayload>({
  name: '',
  sku: '',
  purchase_date: null,
  purchase_price: 0,
  current_status: 'GOOD',
  location: ''
});

const { data: assets, isLoading } = useQuery({
  queryKey: ['assets'],
  queryFn: () => inventoryService.getAssets(1, 100)
});

const totalAssets = computed(() => assets.value?.total || 0);
const goodAssets = computed(() => assets.value?.items?.filter(a => a.current_status === 'GOOD').length || 0);
const loanedAssets = computed(() => assets.value?.items?.filter(a => a.current_status === 'LOANED').length || 0);
const repairAssets = computed(() => assets.value?.items?.filter(a => a.current_status === 'REPAIR' || a.current_status === 'NEEDS_REPAIR').length || 0);
const brokenAssets = computed(() => assets.value?.items?.filter(a => a.current_status === 'BROKEN').length || 0);

const createMutation = useMutation({
  mutationFn: inventoryService.createAsset,
  onSuccess: () => {
    toast.success('Aset berhasil ditambahkan');
    queryClient.invalidateQueries({ queryKey: ['assets'] });
    closeDialog();
  },
  onError: () => toast.error('Gagal menambahkan aset')
});

const updateMutation = useMutation({
  mutationFn: ({ id, payload }: { id: string, payload: CreateAssetPayload }) => inventoryService.updateAsset(id, payload),
  onSuccess: () => {
    toast.success('Aset berhasil diperbarui');
    queryClient.invalidateQueries({ queryKey: ['assets'] });
    closeDialog();
  },
  onError: () => toast.error('Gagal memperbarui aset')
});

const deleteMutation = useMutation({
  mutationFn: inventoryService.deleteAsset,
  onSuccess: () => {
    toast.success('Aset berhasil dihapus');
    queryClient.invalidateQueries({ queryKey: ['assets'] });
  },
  onError: () => toast.error('Gagal menghapus aset')
});

const openEditDialog = (asset: Asset) => {
  editingId.value = asset.id;
  form.value = {
    name: asset.name,
    sku: asset.sku,
    purchase_date: asset.purchase_date,
    purchase_price: asset.purchase_price,
    current_status: asset.current_status,
    location: asset.location
  };
  showAddDialog.value = true;
};

const closeDialog = () => {
  showAddDialog.value = false;
  editingId.value = null;
  form.value = { name: '', sku: '', purchase_date: null, purchase_price: 0, current_status: 'GOOD', location: '' };
};

const submitAsset = () => {
  if (editingId.value) {
    updateMutation.mutate({ id: editingId.value, payload: form.value });
  } else {
    createMutation.mutate(form.value);
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
      <button @click="showAddDialog = true" class="px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 focus:ring-4 focus:ring-emerald-100 transition-colors flex items-center gap-2">
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
              <th class="px-6 py-4 font-medium text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-if="isLoading" class="animate-pulse">
              <td colspan="5" class="px-6 py-4 text-center text-gray-500">Memuat data...</td>
            </tr>
            <tr v-else-if="!assets?.items?.length">
              <td colspan="5" class="px-6 py-12 text-center text-gray-500">
                <BoxIcon class="w-12 h-12 mx-auto text-gray-300 mb-3" />
                Belum ada data aset
              </td>
            </tr>
            <tr v-else v-for="asset in assets?.items" :key="asset.id" class="hover:bg-gray-50 transition-colors">
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
              <td class="px-6 py-4 text-right">
                <button @click="openEditDialog(asset)" class="text-blue-600 hover:text-blue-900 p-2 hover:bg-blue-50 rounded-lg transition-colors mr-2">
                  <EditIcon class="w-4 h-4" />
                </button>
                <button @click="deleteMutation.mutate(asset.id)" class="text-red-600 hover:text-red-900 p-2 hover:bg-red-50 rounded-lg transition-colors">
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
          <h3 class="font-bold text-gray-900">{{ editingId ? 'Edit Aset' : 'Tambah Aset Baru' }}</h3>
          <button @click="closeDialog" class="text-gray-400 hover:text-gray-600">
            <XIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit.prevent="submitAsset" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Aset</label>
            <input v-model="form.name" required type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Cth: Sound System Yamaha" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">SKU / Kode Aset</label>
            <input v-model="form.sku" required type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Cth: AST-001" />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi</label>
              <input v-model="form.location" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Gudang Utama" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select v-model="form.current_status" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="GOOD">Baik (Tersedia)</option>
                <option value="LOANED">Dipinjam</option>
                <option value="NEEDS_REPAIR">Perlu Perbaikan</option>
                <option value="REPAIR">Sedang Diperbaiki</option>
                <option value="BROKEN">Rusak Total</option>
              </select>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Harga Beli</label>
            <input v-model.number="form.purchase_price" required type="number" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
          </div>
          
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="createMutation.isPending.value || updateMutation.isPending.value" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ (createMutation.isPending.value || updateMutation.isPending.value) ? 'Menyimpan...' : 'Simpan Aset' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
