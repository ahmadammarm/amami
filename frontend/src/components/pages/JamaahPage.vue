<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { jamaahService } from '@/api/services/jamaah';
import { useAuthStore } from '@/stores/auth';
import type { Jamaah, CreateJamaahPayload } from '@/types/jamaah';
import { toast } from 'vue-sonner';
import { Users, Search, Plus, Edit2, ShieldAlert } from 'lucide-vue-next';

const authStore = useAuthStore();
const queryClient = useQueryClient();

// Pagination & Search State
const currentPage = ref(1);
const pageSize = ref(10);
const searchQuery = ref('');

// Fetch Data
const { data: paginatedData, isLoading } = useQuery({
  queryKey: ['jamaahs', currentPage, pageSize, searchQuery],
  queryFn: () => jamaahService.getJamaahs(currentPage.value, pageSize.value, searchQuery.value),
  placeholderData: (prev) => prev,
});

const jamaahs = computed(() => paginatedData.value?.data || []);
const totalRecords = computed(() => paginatedData.value?.total || 0);
const totalPages = computed(() => paginatedData.value?.total_pages || 1);

// Permissions
const canManage = computed(() => {
  return ['SUPER_ADMIN', 'SEKRETARIS', 'TAKMIR'].includes(authStore.userRole);
});

// Form State
const showDialog = ref(false);
const editingId = ref<string | null>(null);
const formState = ref<CreateJamaahPayload>({
  full_name: '',
  phone: '',
  address: '',
  is_mustahik: false,
});

const openCreateDialog = () => {
  editingId.value = null;
  formState.value = { full_name: '', phone: '', address: '', is_mustahik: false };
  showDialog.value = true;
};

const openEditDialog = (j: Jamaah) => {
  editingId.value = j.id;
  formState.value = { 
    full_name: j.full_name, 
    phone: j.phone, 
    address: j.address, 
    is_mustahik: j.is_mustahik 
  };
  showDialog.value = true;
};

// Mutations
const createMutation = useMutation({
  mutationFn: jamaahService.createJamaah,
  onSuccess: () => {
    toast.success('Jamaah registered successfully!');
    queryClient.invalidateQueries({ queryKey: ['jamaahs'] });
    showDialog.value = false;
  },
  onError: () => toast.error('Failed to register jamaah'),
});

const updateMutation = useMutation({
  mutationFn: (payload: { id: string; data: CreateJamaahPayload }) => jamaahService.updateJamaah(payload.id, payload.data),
  onSuccess: () => {
    toast.success('Jamaah updated successfully!');
    queryClient.invalidateQueries({ queryKey: ['jamaahs'] });
    showDialog.value = false;
  },
  onError: () => toast.error('Failed to update jamaah'),
});

const submitForm = () => {
  if (!formState.value.full_name) {
    toast.error('Full name is required');
    return;
  }
  
  if (editingId.value) {
    updateMutation.mutate({ id: editingId.value, data: formState.value });
  } else {
    createMutation.mutate(formState.value);
  }
};

const handleSearch = () => {
  currentPage.value = 1;
};
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 bg-white p-6 rounded-2xl shadow-sm border border-gray-100">
      <div class="flex items-center gap-4">
        <div class="h-12 w-12 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center">
          <Users class="w-6 h-6" />
        </div>
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-gray-900">Database Jamaah</h1>
          <p class="text-sm text-gray-500">Kelola registrasi warga dan status Mustahik.</p>
        </div>
      </div>
      <button
        v-if="canManage"
        @click="openCreateDialog"
        class="inline-flex items-center justify-center rounded-lg bg-blue-600 text-white hover:bg-blue-700 h-10 px-4 text-sm font-medium transition-colors shadow-sm"
      >
        <Plus class="w-4 h-4 mr-2" /> Register Jamaah
      </button>
    </div>

    <!-- Search Bar -->
    <div class="flex items-center bg-white p-2 rounded-xl shadow-sm border border-gray-100 max-w-md">
      <Search class="w-5 h-5 text-gray-400 ml-2" />
      <input 
        v-model="searchQuery"
        @keyup.enter="handleSearch"
        type="text" 
        placeholder="Cari nama jamaah..." 
        class="flex-1 border-0 focus:ring-0 text-sm p-2 w-full"
      />
      <button @click="handleSearch" class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors">
        Cari
      </button>
    </div>

    <!-- Table -->
    <div class="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-sm text-left">
          <thead class="text-xs text-gray-700 uppercase bg-gray-50 border-b">
            <tr>
              <th class="px-6 py-4">Nama Lengkap</th>
              <th class="px-6 py-4">No. HP</th>
              <th class="px-6 py-4">Alamat</th>
              <th class="px-6 py-4">Status Zakat</th>
              <th v-if="canManage" class="px-6 py-4 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading" class="border-b">
              <td colspan="5" class="px-6 py-8 text-center text-gray-500">Memuat data jamaah...</td>
            </tr>
            <tr v-else-if="jamaahs.length === 0" class="border-b">
              <td colspan="5" class="px-6 py-8 text-center text-gray-500">Tidak ada jamaah ditemukan.</td>
            </tr>
            <tr v-else v-for="j in jamaahs" :key="j.id" class="border-b hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 font-medium text-gray-900">{{ j.full_name }}</td>
              <td class="px-6 py-4 text-gray-600">{{ j.phone || '-' }}</td>
              <td class="px-6 py-4 text-gray-600 truncate max-w-xs">{{ j.address || '-' }}</td>
              <td class="px-6 py-4">
                <span v-if="j.is_mustahik" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-700">
                  <ShieldAlert class="w-3 h-3 mr-1" /> Mustahik (Penerima)
                </span>
                <span v-else class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-700">
                  Umum (Muzakki)
                </span>
              </td>
              <td v-if="canManage" class="px-6 py-4 text-right">
                <button @click="openEditDialog(j)" class="text-blue-600 hover:text-blue-900 bg-blue-50 hover:bg-blue-100 p-2 rounded-lg transition-colors">
                  <Edit2 class="w-4 h-4" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- Pagination Controls -->
      <div class="px-6 py-4 border-t flex items-center justify-between text-sm text-gray-500 bg-gray-50">
        <div>
          Menampilkan <span class="font-medium">{{ totalRecords === 0 ? 0 : (currentPage - 1) * pageSize + 1 }}</span> s/d 
          <span class="font-medium">{{ Math.min(currentPage * pageSize, totalRecords) }}</span> dari 
          <span class="font-medium">{{ totalRecords }}</span> jamaah
        </div>
        <div class="flex space-x-2">
          <button @click="currentPage--" :disabled="currentPage === 1" class="px-3 py-1.5 border border-gray-300 rounded-md hover:bg-white disabled:opacity-50 transition-colors bg-gray-50">Sebelumnya</button>
          <button @click="currentPage++" :disabled="currentPage >= totalPages" class="px-3 py-1.5 border border-gray-300 rounded-md hover:bg-white disabled:opacity-50 transition-colors bg-gray-50">Selanjutnya</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Dialog Form -->
  <div v-if="showDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
        <h3 class="text-lg font-bold text-gray-900">{{ editingId ? 'Edit Data Jamaah' : 'Register Jamaah Baru' }}</h3>
        <button @click="showDialog = false" class="text-gray-400 hover:text-gray-600 rounded-full p-1 hover:bg-gray-200 transition-colors">&times;</button>
      </div>
      <form @submit.prevent="submitForm">
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap <span class="text-red-500">*</span></label>
            <input type="text" v-model="formState.full_name" required class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2.5 border" placeholder="Cth: Ahmad Abdullah" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nomor HP / WhatsApp</label>
            <input type="text" v-model="formState.phone" class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2.5 border" placeholder="Cth: 08123456789" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Alamat Domisili</label>
            <textarea v-model="formState.address" rows="3" class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2.5 border" placeholder="Jalan, RT/RW, Kelurahan..."></textarea>
          </div>
          <div class="flex items-center mt-4 p-3 bg-purple-50 rounded-lg border border-purple-100">
            <input id="mustahik" type="checkbox" v-model="formState.is_mustahik" class="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded" />
            <label for="mustahik" class="ml-3 block text-sm font-medium text-purple-900">
              Tandai sebagai Mustahik (Penerima Zakat)
            </label>
          </div>
        </div>
        <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
          <button type="button" @click="showDialog = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors shadow-sm">
            Batal
          </button>
          <button type="submit" :disabled="createMutation.isPending.value || updateMutation.isPending.value" class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors shadow-sm flex items-center">
            <span v-if="createMutation.isPending.value || updateMutation.isPending.value" class="w-4 h-4 mr-2 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            Simpan Data
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
