<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { Plus, Pencil, Trash2, PackageSearch, TriangleAlert } from 'lucide-vue-next';
import { inventoryService } from '../../api/services/inventory';
import { queryKeys } from '../../api/queryKeys';
import { feedback } from '../../composables/useFeedback';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import { formatDate, formatNumber, formatRupiah } from '../../lib/format';
import { assetSchema, toAssetPayload, type AssetFormValues, type AssetPayload } from '../../schemas/inventory';
import type { Asset, AssetStatus } from '../../types/inventory';
import type { BadgeTone } from '../../types/ui';
import PageHeader from '../molecules/PageHeader.vue';
import DataPanel from '../molecules/DataPanel.vue';
import EmptyState from '../molecules/EmptyState.vue';
import StatCard from '../molecules/StatCard.vue';
import StatusBadge from '../molecules/StatusBadge.vue';
import FormDialog from '../molecules/FormDialog.vue';
import FormField from '../molecules/FormField.vue';
import { Button } from '../atoms/button';

const ASSET_STATUS: Record<AssetStatus, { label: string; tone: BadgeTone }> = {
  GOOD: { label: 'Baik', tone: 'sajadah' },
  LOANED: { label: 'Dipinjam', tone: 'kuningan' },
  NEEDS_REPAIR: { label: 'Perlu perbaikan', tone: 'kuningan' },
  REPAIR: { label: 'Sedang diperbaiki', tone: 'netral' },
  BROKEN: { label: 'Rusak', tone: 'bata' },
};

/** LOANED is set by the loan workflow, so it is never offered as a manual choice. */
const EDITABLE_STATUSES: AssetStatus[] = ['GOOD', 'NEEDS_REPAIR', 'REPAIR', 'BROKEN'];

const authStore = useAuthStore();
const queryClient = useQueryClient();

const canManage = computed(() => authStore.can(Permission.InventoryManage));

const { data: assets, isLoading } = useQuery({
  queryKey: queryKeys.inventory.assets,
  queryFn: () => inventoryService.getAssets(1, 100),
});

const assetItems = computed(() => assets.value?.items ?? []);
const countByStatus = (...statuses: AssetStatus[]) =>
  assetItems.value.filter(a => statuses.includes(a.current_status)).length;
const problemCount = computed(() => countByStatus('NEEDS_REPAIR', 'REPAIR', 'BROKEN'));

const refreshAssets = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.inventory.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Create / edit
const showForm = ref(false);
const editing = ref<Asset | null>(null);
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

const isOnLoan = computed(() => editing.value?.current_status === 'LOANED');

const saveMutation = useMutation({
  mutationFn: (payload: AssetPayload) =>
    editing.value ? inventoryService.updateAsset(editing.value.id, payload) : inventoryService.createAsset(payload),
  onSuccess: (asset) => {
    const wasEditing = editing.value !== null;
    refreshAssets();
    showForm.value = false;
    feedback.success(wasEditing ? 'Aset diperbarui' : 'Aset ditambahkan', `${asset.name} (${asset.sku}).`);
  },
  onError: (error: unknown) => feedback.error('Aset belum disimpan', error),
});

const openForm = (asset: Asset | null) => {
  editing.value = asset;
  resetForm({
    values: asset
      ? {
        name: asset.name,
        sku: asset.sku,
        purchase_date: asset.purchase_date ? asset.purchase_date.slice(0, 10) : null,
        purchase_price: asset.purchase_price,
        current_status: asset.current_status,
        location: asset.location,
      }
      : emptyForm,
  });
  showForm.value = true;
};

const submitForm = handleSubmit((values) => saveMutation.mutate(toAssetPayload(values)));

const deleteMutation = useMutation({
  mutationFn: inventoryService.deleteAsset,
  onSuccess: () => {
    refreshAssets();
    feedback.success('Aset dihapus');
  },
  onError: (error: unknown) => feedback.error('Aset belum dihapus', error),
});

const confirmDelete = async (asset: Asset) => {
  const confirmed = await feedback.confirm({
    title: 'Hapus aset?',
    message: `${asset.name} (${asset.sku}) akan dihapus dari daftar aset. Riwayat peminjamannya tetap tersimpan.`,
    confirmLabel: 'Hapus aset',
    danger: true,
  });
  if (confirmed) deleteMutation.mutate(asset.id);
};
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Inventaris" description="Aset milik masjid beserta lokasi dan kondisinya.">
      <template v-if="canManage" #actions>
        <Button @click="openForm(null)"><Plus /> Tambah aset</Button>
      </template>
    </PageHeader>

    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <StatCard label="Total aset" :value="formatNumber(assets?.total ?? 0)" />
      <StatCard label="Kondisi baik" :value="formatNumber(countByStatus('GOOD'))" />
      <StatCard label="Dipinjam" :value="formatNumber(countByStatus('LOANED'))" />
      <StatCard
        label="Bermasalah"
        :value="formatNumber(problemCount)"
        hint="perlu perbaikan, diperbaiki, atau rusak"
        :icon="TriangleAlert"
        :tone="problemCount > 0 ? 'alert' : 'default'"
      />
    </div>

    <DataPanel title="Daftar aset">
      <div v-if="isLoading" class="p-8 text-center text-sm text-jati-muda">Memuat aset…</div>
      <EmptyState
        v-else-if="assetItems.length === 0"
        :icon="PackageSearch"
        title="Belum ada aset tercatat"
        :description="canManage ? 'Catat aset masjid agar kondisi dan peminjamannya bisa dipantau.' : 'Aset yang dicatat pengurus akan tampil di sini.'"
      />
      <div v-else class="overflow-x-auto">
        <table class="tabel min-w-[760px]">
          <thead>
            <tr>
              <th>Aset</th>
              <th>Lokasi</th>
              <th>Kondisi</th>
              <th class="text-right">Harga beli</th>
              <th v-if="canManage" class="w-24"><span class="sr-only">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="asset in assetItems" :key="asset.id">
              <td>
                <p class="font-medium">{{ asset.name }}</p>
                <p class="font-mono text-xs text-jati-muda">{{ asset.sku }}</p>
              </td>
              <td class="text-jati-muda">{{ asset.location || '–' }}</td>
              <td><StatusBadge :tone="ASSET_STATUS[asset.current_status].tone">{{ ASSET_STATUS[asset.current_status].label }}</StatusBadge></td>
              <td class="tabular text-right">
                {{ formatRupiah(asset.purchase_price) }}
                <p v-if="asset.purchase_date" class="text-xs text-jati-muda">{{ formatDate(asset.purchase_date) }}</p>
              </td>
              <td v-if="canManage">
                <div class="flex justify-end gap-1">
                  <button type="button" class="rounded-md p-2 text-jati-muda hover:bg-kapur hover:text-sajadah" :aria-label="`Ubah ${asset.name}`" title="Ubah" @click="openForm(asset)">
                    <Pencil class="size-4" />
                  </button>
                  <button
                    type="button"
                    class="rounded-md p-2 text-jati-muda hover:bg-bata-muda hover:text-bata disabled:cursor-not-allowed disabled:opacity-30 disabled:hover:bg-transparent disabled:hover:text-jati-muda"
                    :disabled="asset.current_status === 'LOANED'"
                    :aria-label="`Hapus ${asset.name}`"
                    :title="asset.current_status === 'LOANED' ? 'Aset sedang dipinjam' : 'Hapus'"
                    @click="confirmDelete(asset)"
                  >
                    <Trash2 class="size-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </DataPanel>

    <FormDialog
      v-model:open="showForm"
      :title="editing ? 'Ubah aset' : 'Tambah aset'"
      :submit-label="editing ? 'Simpan perubahan' : 'Tambah aset'"
      :submitting="saveMutation.isPending.value"
      @submit="submitForm"
    >
      <FormField label="Nama aset" for="asset-name" :error="errors.name" required>
        <input id="asset-name" v-model="name" type="text" class="field-control" placeholder="Contoh: Sound system Yamaha" :aria-invalid="Boolean(errors.name)" />
      </FormField>
      <div class="grid grid-cols-2 gap-4">
        <FormField label="Kode aset (SKU)" for="asset-sku" :error="errors.sku" required>
          <input id="asset-sku" v-model="sku" type="text" class="field-control font-mono" placeholder="AST-001" :aria-invalid="Boolean(errors.sku)" />
        </FormField>
        <FormField label="Lokasi" for="asset-location" :error="errors.location">
          <input id="asset-location" v-model="location" type="text" class="field-control" placeholder="Gudang" />
        </FormField>
      </div>
      <FormField
        label="Kondisi"
        for="asset-status"
        :hint="isOnLoan ? 'Sedang dipinjam. Kondisi berubah otomatis saat aset dikembalikan.' : undefined"
        required
      >
        <select id="asset-status" v-model="currentStatus" class="field-control" :disabled="isOnLoan">
          <option v-if="isOnLoan" value="LOANED">{{ ASSET_STATUS.LOANED.label }}</option>
          <option v-for="status in EDITABLE_STATUSES" :key="status" :value="status">{{ ASSET_STATUS[status].label }}</option>
        </select>
      </FormField>
      <div class="grid grid-cols-2 gap-4">
        <FormField label="Harga beli (Rp)" for="asset-price" :error="errors.purchase_price">
          <input id="asset-price" v-model.number="purchasePrice" type="number" min="0" step="1" class="field-control tabular" />
        </FormField>
        <FormField label="Tanggal beli" for="asset-date">
          <input id="asset-date" v-model="purchaseDate" type="date" class="field-control" />
        </FormField>
      </div>
    </FormDialog>
  </div>
</template>
