<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Manajemen Qurban</h1>
        <p class="text-sm text-gray-500 mt-1">Kelola paket qurban, pendaftaran shohibul, dan hewan qurban</p>
      </div>
      <div class="flex gap-2">
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
            <p class="text-2xl font-bold text-gray-900">{{ packages?.total || 0 }}</p>
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
            <p class="text-2xl font-bold text-gray-900">{{ bookings?.length || 0 }}</p>
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
            <p class="text-2xl font-bold text-gray-900">{{ animals?.length || 0 }}</p>
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
                <th class="px-6 py-4 font-medium text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="isLoadingPackages">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Memuat data...</td>
              </tr>
              <tr v-else-if="!packages?.items?.length">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Belum ada paket qurban.</td>
              </tr>
              <tr v-else v-for="pkg in packages?.items" :key="pkg.id" class="hover:bg-gray-50 transition-colors">
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
                <td class="px-6 py-4 text-right">
                  <button @click="openEditPackageDialog(pkg)" class="text-blue-600 hover:text-blue-900 p-2 hover:bg-blue-50 rounded-lg transition-colors mr-2">
                    <EditIcon class="w-4 h-4" />
                  </button>
                  <button @click="deletePackageMutation.mutate(pkg.id)" class="text-red-600 hover:text-red-900 p-2 hover:bg-red-50 rounded-lg transition-colors">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Bookings Tab -->
        <div v-if="activeTab === 'bookings'" class="p-6">
          <div class="mb-4 flex justify-end">
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
                <tr v-else-if="!bookings?.length">
                  <td colspan="5" class="px-6 py-8 text-center text-gray-500">Belum ada pendaftaran qurban.</td>
                </tr>
                <tr v-else v-for="b in bookings" :key="b.id" class="hover:bg-gray-50 transition-colors">
                  <td class="px-6 py-4 text-gray-600">{{ new Date(b.booking_date).toLocaleDateString('id-ID') }}</td>
                  <td class="px-6 py-4 font-medium text-gray-900">ID: {{ b.shohibul_id.substring(0, 8) }}...</td>
                  <td class="px-6 py-4 text-gray-600">ID Paket: {{ b.package_id }}</td>
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
                <th class="px-6 py-4 font-medium text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-if="isLoadingAnimals">
                <td colspan="6" class="px-6 py-8 text-center text-gray-500">Memuat data...</td>
              </tr>
              <tr v-else-if="!animals?.length">
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
                <td class="px-6 py-4 text-right">
                  <button @click="openEditAnimalDialog(animal)" class="text-blue-600 hover:text-blue-900 p-2 hover:bg-blue-50 rounded-lg transition-colors mr-2">
                    <EditIcon class="w-4 h-4" />
                  </button>
                  <button @click="deleteAnimalMutation.mutate(animal.id)" class="text-red-600 hover:text-red-900 p-2 hover:bg-red-50 rounded-lg transition-colors">
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
        <form @submit.prevent="submitPackage" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Paket</label>
            <input v-model="packageForm.name" type="text" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all" placeholder="Contoh: Paket A Sapi">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Tipe</label>
            <select v-model="packageForm.type" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <option value="SAPI">Sapi</option>
              <option value="KAMBING">Kambing</option>
              <option value="DOMBA">Domba</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tahun Hijriah</label>
              <input v-model.number="packageForm.year_hijri" type="number" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Total Stok</label>
              <input v-model.number="packageForm.stock_total" type="number" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Harga (Rp)</label>
            <input v-model.number="packageForm.price" type="number" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
          </div>
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closePackageDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="createPackageMutation.isPending.value || updatePackageMutation.isPending.value" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ (createPackageMutation.isPending.value || updatePackageMutation.isPending.value) ? 'Menyimpan...' : 'Simpan' }}
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
        <form @submit.prevent="submitAnimal" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">No Tagging</label>
            <input v-model="animalForm.tag_number" type="text" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Tipe</label>
              <select v-model="animalForm.type" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
                <option value="SAPI">Sapi</option>
                <option value="KAMBING">Kambing</option>
                <option value="DOMBA">Domba</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Berat (kg)</label>
              <input v-model.number="animalForm.weight" type="number" step="0.1" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Vendor / Penjual</label>
            <input v-model="animalForm.vendor_info" type="text" class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select v-model="animalForm.status" required class="w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:bg-white focus:ring-2 focus:ring-emerald-500 outline-none transition-all">
              <option value="AVAILABLE">Tersedia</option>
              <option value="ASSIGNED">Dialokasikan</option>
              <option value="SLAUGHTERED">Disembelih</option>
            </select>
          </div>
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeAnimalDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="createAnimalMutation.isPending.value || updateAnimalMutation.isPending.value" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ (createAnimalMutation.isPending.value || updateAnimalMutation.isPending.value) ? 'Menyimpan...' : 'Simpan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
    
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { qurbanService } from '../../api/services/qurban';
import type { CreatePackagePayload, CreateAnimalPayload, QurbanPackage, QurbanAnimal } from '../../types/qurban';
import { PlusIcon, XIcon, TrashIcon, EditIcon, PackageIcon, UsersIcon, DogIcon } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

const queryClient = useQueryClient();

const tabs = [
  { id: 'packages', name: 'Paket Qurban' },
  { id: 'bookings', name: 'Pendaftaran' },
  { id: 'animals', name: 'Hewan Qurban' }
];
const activeTab = ref('packages');

// Queries
const { data: packages, isLoading: isLoadingPackages } = useQuery({
  queryKey: ['qurban_packages'],
  queryFn: () => qurbanService.getPackages(1, 100)
});

const { data: bookings, isLoading: isLoadingBookings } = useQuery({
  queryKey: ['qurban_bookings'],
  queryFn: qurbanService.getBookings
});

const { data: animals, isLoading: isLoadingAnimals } = useQuery({
  queryKey: ['qurban_animals'],
  queryFn: qurbanService.getAnimals
});

// Dialog States
const showPackageDialog = ref(false);
const editingPackageId = ref<number | null>(null);
const packageForm = ref<CreatePackagePayload>({
  name: '',
  type: 'SAPI',
  price: 0,
  year_hijri: 1447,
  stock_total: 7
});

const showAnimalDialog = ref(false);
const editingAnimalId = ref<string | null>(null);
const animalForm = ref<CreateAnimalPayload>({
  tag_number: '',
  type: 'SAPI',
  weight: 0,
  status: 'AVAILABLE',
  vendor_info: ''
});

// Package Mutations
const createPackageMutation = useMutation({
  mutationFn: qurbanService.createPackage,
  onSuccess: () => {
    toast.success('Paket berhasil dibuat');
    queryClient.invalidateQueries({ queryKey: ['qurban_packages'] });
    closePackageDialog();
  },
  onError: () => toast.error('Gagal membuat paket')
});

const updatePackageMutation = useMutation({
  mutationFn: ({ id, payload }: { id: number, payload: CreatePackagePayload }) => qurbanService.updatePackage(id, payload),
  onSuccess: () => {
    toast.success('Paket diperbarui');
    queryClient.invalidateQueries({ queryKey: ['qurban_packages'] });
    closePackageDialog();
  },
  onError: () => toast.error('Gagal memperbarui paket')
});

const deletePackageMutation = useMutation({
  mutationFn: qurbanService.deletePackage,
  onSuccess: () => {
    toast.success('Paket dihapus');
    queryClient.invalidateQueries({ queryKey: ['qurban_packages'] });
  },
  onError: () => toast.error('Gagal menghapus paket')
});

// Animal Mutations
const createAnimalMutation = useMutation({
  mutationFn: qurbanService.createAnimal,
  onSuccess: () => {
    toast.success('Hewan berhasil didata');
    queryClient.invalidateQueries({ queryKey: ['qurban_animals'] });
    closeAnimalDialog();
  },
  onError: () => toast.error('Gagal mendata hewan')
});

const updateAnimalMutation = useMutation({
  mutationFn: ({ id, payload }: { id: string, payload: CreateAnimalPayload }) => qurbanService.updateAnimal(id, payload),
  onSuccess: () => {
    toast.success('Data hewan diperbarui');
    queryClient.invalidateQueries({ queryKey: ['qurban_animals'] });
    closeAnimalDialog();
  },
  onError: () => toast.error('Gagal memperbarui data hewan')
});

const deleteAnimalMutation = useMutation({
  mutationFn: qurbanService.deleteAnimal,
  onSuccess: () => {
    toast.success('Data hewan dihapus');
    queryClient.invalidateQueries({ queryKey: ['qurban_animals'] });
  },
  onError: () => toast.error('Gagal menghapus data hewan')
});

// Handlers
const openPackageDialog = () => {
  editingPackageId.value = null;
  packageForm.value = { name: '', type: 'SAPI', price: 0, year_hijri: 1447, stock_total: 1 };
  showPackageDialog.value = true;
};

const openEditPackageDialog = (pkg: QurbanPackage) => {
  editingPackageId.value = pkg.id;
  packageForm.value = { 
    name: pkg.name, 
    type: pkg.type, 
    price: pkg.price, 
    year_hijri: pkg.year_hijri, 
    stock_total: pkg.stock_total 
  };
  showPackageDialog.value = true;
};

const closePackageDialog = () => { showPackageDialog.value = false; };

const submitPackage = () => {
  if (editingPackageId.value) {
    updatePackageMutation.mutate({ id: editingPackageId.value, payload: packageForm.value });
  } else {
    createPackageMutation.mutate(packageForm.value);
  }
};

const openAnimalDialog = () => {
  editingAnimalId.value = null;
  animalForm.value = { tag_number: '', type: 'SAPI', weight: 0, status: 'AVAILABLE', vendor_info: '' };
  showAnimalDialog.value = true;
};

const openEditAnimalDialog = (animal: QurbanAnimal) => {
  editingAnimalId.value = animal.id;
  animalForm.value = { 
    tag_number: animal.tag_number, 
    type: animal.type, 
    weight: animal.weight, 
    status: animal.status, 
    vendor_info: animal.vendor_info 
  };
  showAnimalDialog.value = true;
};

const closeAnimalDialog = () => { showAnimalDialog.value = false; };

const submitAnimal = () => {
  if (editingAnimalId.value) {
    updateAnimalMutation.mutate({ id: editingAnimalId.value, payload: animalForm.value });
  } else {
    createAnimalMutation.mutate(animalForm.value);
  }
};

const openBookingDialog = () => {
  toast.info('Pendaftaran Qurban akan mengarahkan ke form pemilihan Jamaah dan Pembayaran.');
};

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(amount);
};
</script>
