<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { qurbanService } from '../../api/services/qurban';
import { getErrorMessage } from '../../api/client';
import { queryKeys } from '../../api/queryKeys';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import {
  animalSchema,
  bookingSchema,
  packageSchema,
  type AnimalPayload,
  type BookingPayload,
  type PackagePayload,
} from '../../schemas/qurban';
import type { QurbanAnimal, QurbanPackage } from '../../types/qurban';
import AsyncJamaahSelect from '../molecules/AsyncJamaahSelect.vue';
import { PlusIcon, XIcon, TrashIcon, EditIcon, PackageIcon, UsersIcon, DogIcon } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

type TabId = 'packages' | 'bookings' | 'animals';

const authStore = useAuthStore();
const queryClient = useQueryClient();

const canManage = computed(() => authStore.can(Permission.QurbanManage));

const tabs: { id: TabId; name: string }[] = [
  { id: 'packages', name: 'Paket Qurban' },
  { id: 'bookings', name: 'Pendaftaran' },
  { id: 'animals', name: 'Hewan Qurban' }
];
const activeTab = ref<TabId>('packages');

// Queries
const { data: packages, isLoading: isLoadingPackages } = useQuery({
  queryKey: queryKeys.qurban.packages,
  queryFn: () => qurbanService.getPackages(1, 100)
});

const { data: bookings, isLoading: isLoadingBookings } = useQuery({
  queryKey: queryKeys.qurban.bookings,
  queryFn: qurbanService.getBookings
});

const { data: animals, isLoading: isLoadingAnimals } = useQuery({
  queryKey: queryKeys.qurban.animals,
  queryFn: qurbanService.getAnimals
});

const packageItems = computed(() => packages.value?.items ?? []);
const availablePackages = computed(() => packageItems.value.filter(p => p.stock_remaining > 0));

const refreshQurban = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.qurban.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Package form
const showPackageDialog = ref(false);
const editingPackageId = ref<number | null>(null);
const emptyPackage: PackagePayload = { name: '', type: 'SAPI', price: 0, year_hijri: 1447, stock_total: 1 };

const {
  handleSubmit: handlePackageSubmit,
  errors: packageErrors,
  defineField: definePackageField,
  resetForm: resetPackageForm,
} = useForm({ validationSchema: toTypedSchema(packageSchema), initialValues: emptyPackage });

const [packageName] = definePackageField('name');
const [packageType] = definePackageField('type');
const [packagePrice] = definePackageField('price');
const [packageYear] = definePackageField('year_hijri');
const [packageStock] = definePackageField('stock_total');

// Animal form
const showAnimalDialog = ref(false);
const editingAnimalId = ref<string | null>(null);
const emptyAnimal: AnimalPayload = { tag_number: '', type: 'SAPI', weight: 0, status: 'AVAILABLE', vendor_info: '' };

const {
  handleSubmit: handleAnimalSubmit,
  errors: animalErrors,
  defineField: defineAnimalField,
  resetForm: resetAnimalForm,
} = useForm({ validationSchema: toTypedSchema(animalSchema), initialValues: emptyAnimal });

const [animalTag] = defineAnimalField('tag_number');
const [animalType] = defineAnimalField('type');
const [animalWeight] = defineAnimalField('weight');
const [animalStatus] = defineAnimalField('status');
const [animalVendor] = defineAnimalField('vendor_info');

// Booking form
const showBookingDialog = ref(false);
const emptyBooking: BookingPayload = { shohibul_id: '', package_id: 0, payment_status: 'PENDING' };

const {
  handleSubmit: handleBookingSubmit,
  errors: bookingErrors,
  defineField: defineBookingField,
  resetForm: resetBookingForm,
} = useForm({ validationSchema: toTypedSchema(bookingSchema), initialValues: emptyBooking });

const [shohibulId] = defineBookingField('shohibul_id');
const [bookingPackageId] = defineBookingField('package_id');
const [paymentStatus] = defineBookingField('payment_status');

const selectedPackagePrice = computed(
  () => packageItems.value.find(p => p.id === bookingPackageId.value)?.price ?? 0,
);

// Package Mutations
const createPackageMutation = useMutation({
  mutationFn: qurbanService.createPackage,
  onSuccess: () => {
    toast.success('Paket berhasil dibuat');
    refreshQurban();
    closePackageDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal membuat paket'))
});

const updatePackageMutation = useMutation({
  mutationFn: ({ id, payload }: { id: number, payload: PackagePayload }) => qurbanService.updatePackage(id, payload),
  onSuccess: () => {
    toast.success('Paket diperbarui');
    refreshQurban();
    closePackageDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal memperbarui paket'))
});

const deletePackageMutation = useMutation({
  mutationFn: qurbanService.deletePackage,
  onSuccess: () => {
    toast.success('Paket dihapus');
    refreshQurban();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal menghapus paket'))
});

// Animal Mutations
const createAnimalMutation = useMutation({
  mutationFn: qurbanService.createAnimal,
  onSuccess: () => {
    toast.success('Hewan berhasil didata');
    refreshQurban();
    closeAnimalDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal mendata hewan'))
});

const updateAnimalMutation = useMutation({
  mutationFn: ({ id, payload }: { id: string, payload: AnimalPayload }) => qurbanService.updateAnimal(id, payload),
  onSuccess: () => {
    toast.success('Data hewan diperbarui');
    refreshQurban();
    closeAnimalDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal memperbarui data hewan'))
});

const deleteAnimalMutation = useMutation({
  mutationFn: qurbanService.deleteAnimal,
  onSuccess: () => {
    toast.success('Data hewan dihapus');
    refreshQurban();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal menghapus data hewan'))
});

// Booking Mutation
const createBookingMutation = useMutation({
  mutationFn: qurbanService.createBooking,
  onSuccess: (booking) => {
    toast.success(`Pendaftaran ${booking.shohibul_name} berhasil`);
    refreshQurban();
    closeBookingDialog();
  },
  onError: (error: unknown) => toast.error(getErrorMessage(error, 'Gagal mendaftarkan shohibul'))
});

const isSavingPackage = computed(() => createPackageMutation.isPending.value || updatePackageMutation.isPending.value);
const isSavingAnimal = computed(() => createAnimalMutation.isPending.value || updateAnimalMutation.isPending.value);

// Handlers
const openPackageDialog = () => {
  editingPackageId.value = null;
  resetPackageForm({ values: emptyPackage });
  showPackageDialog.value = true;
};

const openEditPackageDialog = (pkg: QurbanPackage) => {
  editingPackageId.value = pkg.id;
  resetPackageForm({
    values: {
      name: pkg.name,
      type: pkg.type,
      price: pkg.price,
      year_hijri: pkg.year_hijri,
      stock_total: pkg.stock_total,
    },
  });
  showPackageDialog.value = true;
};

const closePackageDialog = () => { showPackageDialog.value = false; };

const submitPackage = handlePackageSubmit((values) => {
  if (editingPackageId.value) {
    updatePackageMutation.mutate({ id: editingPackageId.value, payload: values });
  } else {
    createPackageMutation.mutate(values);
  }
});

const confirmDeletePackage = (pkg: QurbanPackage) => {
  if (window.confirm(`Hapus paket "${pkg.name}"?`)) deletePackageMutation.mutate(pkg.id);
};

const openAnimalDialog = () => {
  editingAnimalId.value = null;
  resetAnimalForm({ values: emptyAnimal });
  showAnimalDialog.value = true;
};

const openEditAnimalDialog = (animal: QurbanAnimal) => {
  editingAnimalId.value = animal.id;
  resetAnimalForm({
    values: {
      tag_number: animal.tag_number,
      type: animal.type,
      weight: animal.weight,
      status: animal.status,
      vendor_info: animal.vendor_info,
    },
  });
  showAnimalDialog.value = true;
};

const closeAnimalDialog = () => { showAnimalDialog.value = false; };

const submitAnimal = handleAnimalSubmit((values) => {
  if (editingAnimalId.value) {
    updateAnimalMutation.mutate({ id: editingAnimalId.value, payload: values });
  } else {
    createAnimalMutation.mutate(values);
  }
});

const confirmDeleteAnimal = (animal: QurbanAnimal) => {
  if (window.confirm(`Hapus data hewan ${animal.tag_number}?`)) deleteAnimalMutation.mutate(animal.id);
};

const openBookingDialog = () => {
  resetBookingForm({ values: emptyBooking });
  showBookingDialog.value = true;
};

const closeBookingDialog = () => { showBookingDialog.value = false; };

const submitBooking = handleBookingSubmit((values) => createBookingMutation.mutate(values));

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amount);
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Manajemen Qurban</h1>
        <p class="text-sm text-gray-500 mt-1">Kelola paket qurban, pendaftaran shohibul, dan hewan qurban</p>
      </div>
      <div v-if="canManage" class="flex gap-2">
        <button @click="openPackageDialog" class="inline-flex items-center justify-center rounded-xl text-sm font-medium bg-emerald-600 text-white hover:bg-emerald-700 h-10 px-4 transition-colors">
          <PlusIcon class="w-4 h-4 mr-2" />
          Tambah Paket
        </button>
        <button @click="openAnimalDialog" class="inline-flex items-center justify-center rounded-xl text-sm font-medium bg-amber-600 text-white hover:bg-amber-700 h-10 px-4 transition-colors">
          <PlusIcon class="w-4 h-4 mr-2" />
          Tambah Hewan
        </button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-emerald-50 text-emerald-600 rounded-xl">
            <PackageIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Total Paket</p>
            <p class="text-2xl font-bold text-gray-900">{{ packages?.total ?? 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-blue-50 text-blue-600 rounded-xl">
            <UsersIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Shohibul Qurban</p>
            <p class="text-2xl font-bold text-gray-900">{{ bookings?.length ?? 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
        <div class="flex items-center gap-4">
          <div class="p-3 bg-amber-50 text-amber-600 rounded-xl">
            <DogIcon class="w-6 h-6" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-500">Hewan Qurban</p>
            <p class="text-2xl font-bold text-gray-900">{{ animals?.length ?? 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div class="border-b border-gray-100">
        <nav class="flex gap-4 px-6" aria-label="Tabs">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              activeTab === tab.id
                ? 'border-emerald-500 text-emerald-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-colors'
            ]"
          >
            {{ tab.name }}
          </button>
        </nav>
      </div>

      <div class="p-0">
        <!-- Packages Tab -->
        <div v-if="activeTab === 'packages'" class="overflow-x-auto">
          <table class="w-full text-sm text-left">
            <thead class="text-xs text-gray-500 uppercase bg-gray-50">
              <tr>
                <th class="px-6 py-4 font-medium">Nama Paket</th>
                <th class="px-6 py-4 font-medium">Tipe</th>
                <th class="px-6 py-4 font-medium">Harga</th>
                <th class="px-6 py-4 font-medium">Tahun (H)</th>
                <th class="px-6 py-4 font-medium">Stok</th>
                <th v-if="canManage" class="px-6 py-4 font-medium text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="isLoadingPackages">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Memuat data...</td>
              </tr>
              <tr v-else-if="packageItems.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Belum ada paket qurban.</td>
              </tr>
              <tr v-else v-for="pkg in packageItems" :key="pkg.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 font-medium text-gray-900">{{ pkg.name }}</td>
                <td class="px-6 py-4 text-gray-600">{{ pkg.type }}</td>
                <td class="px-6 py-4 text-gray-600">{{ formatCurrency(pkg.price) }}</td>
                <td class="px-6 py-4 text-gray-600">{{ pkg.year_hijri }}</td>
                <td class="px-6 py-4">
                  <span :class="[
                    'px-2 py-1 rounded-full text-xs font-medium',
                    pkg.stock_remaining > 0 ? 'bg-emerald-100 text-emerald-700' : 'bg-red-100 text-red-700'
                  ]">
                    Sisa: {{ pkg.stock_remaining }} / {{ pkg.stock_total }}
                  </span>
                </td>
                <td v-if="canManage" class="px-6 py-4 text-right">
                  <button @click="openEditPackageDialog(pkg)" class="text-blue-600 hover:text-blue-900 p-2 hover:bg-blue-50 rounded-lg transition-colors mr-2" title="Edit">
                    <EditIcon class="w-4 h-4" />
                  </button>
                  <button @click="confirmDeletePackage(pkg)" class="text-red-600 hover:text-red-900 p-2 hover:bg-red-50 rounded-lg transition-colors" title="Hapus">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Bookings Tab -->
        <div v-if="activeTab === 'bookings'" class="p-6">
          <div v-if="canManage" class="mb-4 flex justify-end">
            <button @click="openBookingDialog" class="inline-flex items-center justify-center rounded-xl text-sm font-medium bg-blue-600 text-white hover:bg-blue-700 h-10 px-4 transition-colors">
              <PlusIcon class="w-4 h-4 mr-2" />
              Pendaftaran Baru
            </button>
          </div>
          <div class="overflow-x-auto border border-gray-100 rounded-xl">
            <table class="w-full text-sm text-left">
              <thead class="text-xs text-gray-500 uppercase bg-gray-50">
                <tr>
                  <th class="px-6 py-4 font-medium">Tanggal</th>
                  <th class="px-6 py-4 font-medium">Shohibul</th>
                  <th class="px-6 py-4 font-medium">Paket</th>
                  <th class="px-6 py-4 font-medium">Total</th>
                  <th class="px-6 py-4 font-medium">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-if="isLoadingBookings">
                  <td colspan="5" class="px-6 py-8 text-center text-gray-500">Memuat data...</td>
                </tr>
                <tr v-else-if="!bookings || bookings.length === 0">
                  <td colspan="5" class="px-6 py-8 text-center text-gray-500">Belum ada pendaftaran qurban.</td>
                </tr>
                <tr v-else v-for="b in bookings" :key="b.id" class="hover:bg-gray-50 transition-colors">
                  <td class="px-6 py-4 text-gray-600">{{ new Date(b.booking_date).toLocaleDateString('id-ID') }}</td>
                  <td class="px-6 py-4 font-medium text-gray-900">{{ b.shohibul_name }}</td>
                  <td class="px-6 py-4 text-gray-600">{{ b.package_name }}</td>
                  <td class="px-6 py-4 text-gray-600">{{ formatCurrency(b.total_amount) }}</td>
                  <td class="px-6 py-4">
                    <span :class="[
                      'px-2 py-1 rounded-full text-xs font-medium',
                      b.payment_status === 'PAID' ? 'bg-emerald-100 text-emerald-700' : 'bg-amber-100 text-amber-700'
                    ]">
                      {{ b.payment_status }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Animals Tab -->
        <div v-if="activeTab === 'animals'" class="overflow-x-auto">
          <table class="w-full text-sm text-left">
            <thead class="text-xs text-gray-500 uppercase bg-gray-50">
              <tr>
                <th class="px-6 py-4 font-medium">No. Tag</th>
                <th class="px-6 py-4 font-medium">Tipe</th>
                <th class="px-6 py-4 font-medium">Berat (kg)</th>
                <th class="px-6 py-4 font-medium">Vendor</th>
                <th class="px-6 py-4 font-medium">Status</th>
                <th v-if="canManage" class="px-6 py-4 font-medium text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="isLoadingAnimals">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Memuat data...</td>
              </tr>
              <tr v-else-if="!animals || animals.length === 0">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Belum ada data hewan qurban.</td>
              </tr>
              <tr v-else v-for="animal in animals" :key="animal.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-6 py-4 font-medium text-gray-900">{{ animal.tag_number }}</td>
                <td class="px-6 py-4 text-gray-600">{{ animal.type }}</td>
                <td class="px-6 py-4 text-gray-600">{{ animal.weight }}</td>
                <td class="px-6 py-4 text-gray-600">{{ animal.vendor_info }}</td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-700">
                    {{ animal.status }}
                  </span>
                </td>
                <td v-if="canManage" class="px-6 py-4 text-right">
                  <button @click="openEditAnimalDialog(animal)" class="text-blue-600 hover:text-blue-900 p-2 hover:bg-blue-50 rounded-lg transition-colors mr-2" title="Edit">
                    <EditIcon class="w-4 h-4" />
                  </button>
                  <button @click="confirmDeleteAnimal(animal)" class="text-red-600 hover:text-red-900 p-2 hover:bg-red-50 rounded-lg transition-colors" title="Hapus">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Package Dialog -->
    <div v-if="showPackageDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-bold text-gray-900">{{ editingPackageId ? 'Edit Paket' : 'Tambah Paket Qurban' }}</h3>
          <button @click="closePackageDialog" class="text-gray-400 hover:text-gray-600">
            <XIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit="submitPackage" class="p-6 space-y-4" novalidate>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Paket</label>
            <input v-model="packageName" type="text" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all" placeholder="Contoh: Paket A Sapi">
            <p v-if="packageErrors.name" class="mt-1 text-xs text-red-600">{{ packageErrors.name }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tipe</label>
            <select v-model="packageType" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <option value="SAPI">Sapi</option>
              <option value="SAPI_KOLEKTIF">Sapi Kolektif (1/7)</option>
              <option value="KAMBING">Kambing</option>
              <option value="DOMBA">Domba</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tahun Hijriah</label>
              <input v-model.number="packageYear" type="number" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <p v-if="packageErrors.year_hijri" class="mt-1 text-xs text-red-600">{{ packageErrors.year_hijri }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Total Stok</label>
              <input v-model.number="packageStock" type="number" min="1" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <p v-if="packageErrors.stock_total" class="mt-1 text-xs text-red-600">{{ packageErrors.stock_total }}</p>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Harga (Rp)</label>
            <input v-model.number="packagePrice" type="number" min="0" step="1" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
            <p v-if="packageErrors.price" class="mt-1 text-xs text-red-600">{{ packageErrors.price }}</p>
          </div>
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closePackageDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="isSavingPackage" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ isSavingPackage ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Animal Dialog -->
    <div v-if="showAnimalDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-bold text-gray-900">{{ editingAnimalId ? 'Edit Hewan' : 'Tambah Hewan' }}</h3>
          <button @click="closeAnimalDialog" class="text-gray-400 hover:text-gray-600">
            <XIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit="submitAnimal" class="p-6 space-y-4" novalidate>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">No Tagging</label>
            <input v-model="animalTag" type="text" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
            <p v-if="animalErrors.tag_number" class="mt-1 text-xs text-red-600">{{ animalErrors.tag_number }}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tipe</label>
              <select v-model="animalType" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
                <option value="SAPI">Sapi</option>
                <option value="KAMBING">Kambing</option>
                <option value="DOMBA">Domba</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Berat (kg)</label>
              <input v-model.number="animalWeight" type="number" step="0.1" min="0" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <p v-if="animalErrors.weight" class="mt-1 text-xs text-red-600">{{ animalErrors.weight }}</p>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Vendor / Penjual</label>
            <input v-model="animalVendor" type="text" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select v-model="animalStatus" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <option value="AVAILABLE">Tersedia</option>
              <option value="ASSIGNED">Dialokasikan</option>
              <option value="SLAUGHTERED">Disembelih</option>
            </select>
          </div>
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeAnimalDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="isSavingAnimal" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ isSavingAnimal ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Booking Dialog -->
    <div v-if="showBookingDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-bold text-gray-900">Pendaftaran Shohibul Qurban</h3>
          <button @click="closeBookingDialog" class="text-gray-400 hover:text-gray-600">
            <XIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit="submitBooking" class="p-6 space-y-4" novalidate>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Shohibul</label>
            <AsyncJamaahSelect v-model="shohibulId" placeholder="Cari jamaah..." />
            <p v-if="bookingErrors.shohibul_id" class="mt-1 text-xs text-red-600">{{ bookingErrors.shohibul_id }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Paket</label>
            <select v-model.number="bookingPackageId" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <option :value="0" disabled>Pilih paket</option>
              <option v-for="pkg in availablePackages" :key="pkg.id" :value="pkg.id">
                {{ pkg.name }} ({{ formatCurrency(pkg.price) }}, sisa {{ pkg.stock_remaining }})
              </option>
            </select>
            <p v-if="bookingErrors.package_id" class="mt-1 text-xs text-red-600">{{ bookingErrors.package_id }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status Pembayaran</label>
            <select v-model="paymentStatus" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <option value="PENDING">Belum Lunas</option>
              <option value="PAID">Lunas</option>
            </select>
          </div>
          <div v-if="selectedPackagePrice > 0" class="flex items-center justify-between rounded-lg bg-gray-50 p-3 text-sm">
            <span class="text-gray-600">Total</span>
            <span class="font-semibold text-gray-900">{{ formatCurrency(selectedPackagePrice) }}</span>
          </div>
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeBookingDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="createBookingMutation.isPending.value" class="flex-1 px-4 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50">
              {{ createBookingMutation.isPending.value ? 'Menyimpan...' : 'Daftarkan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
