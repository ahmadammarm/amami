<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { jamaahService } from '@/api/services/jamaah';
import { feedback } from '@/composables/useFeedback';
import { queryKeys } from '@/api/queryKeys';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { scoreBadgeClass, scoreLevel } from '@/lib/mustahik';
import { jamaahSchema, type JamaahPayload } from '@/schemas/jamaah';
import type { Jamaah, JamaahListParams, JamaahSort } from '@/types/jamaah';
import MustahikSurveyDialog from '@/components/molecules/MustahikSurveyDialog.vue';
import { Users, Search, Plus, Edit2, ShieldAlert, ClipboardList } from 'lucide-vue-next';
import { RecycleScroller } from 'vue-virtual-scroller';
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css';

const authStore = useAuthStore();
const queryClient = useQueryClient();

// Pagination, Search & Filter State
const currentPage = ref(1);
const pageSize = ref(10);
const searchInput = ref('');
const appliedSearch = ref('');
const mustahikOnly = ref(false);
const sort = ref<JamaahSort>('newest');

const listParams = computed<JamaahListParams>(() => ({
  page: currentPage.value,
  limit: pageSize.value,
  query: appliedSearch.value || undefined,
  mustahik: mustahikOnly.value ? true : undefined,
  sort: sort.value,
}));

// Fetch Data
const { data: paginatedData, isLoading } = useQuery({
  queryKey: computed(() => queryKeys.jamaah.list(listParams.value)),
  queryFn: () => jamaahService.getJamaahs(listParams.value),
  placeholderData: (prev) => prev,
});

const jamaahs = computed(() => paginatedData.value?.items ?? []);
const totalRecords = computed(() => paginatedData.value?.total ?? 0);
const totalPages = computed(() => paginatedData.value?.total_pages || 1);

// Permissions
const canManage = computed(() => authStore.can(Permission.JamaahManage));
const canSurvey = computed(() => authStore.can(Permission.MustahikSurvey));
const showActions = computed(() => canManage.value || canSurvey.value);

// Form State
const showDialog = ref(false);
const editingId = ref<string | null>(null);
const surveyTarget = ref<Jamaah | null>(null);
const emptyForm: JamaahPayload = { full_name: '', phone: '', address: '', is_mustahik: false };

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: toTypedSchema(jamaahSchema),
  initialValues: emptyForm,
});

const [fullName] = defineField('full_name');
const [phone] = defineField('phone');
const [address] = defineField('address');
const [isMustahik] = defineField('is_mustahik');

const openCreateDialog = () => {
  editingId.value = null;
  resetForm({ values: emptyForm });
  showDialog.value = true;
};

const openEditDialog = (j: Jamaah) => {
  editingId.value = j.id;
  resetForm({
    values: { full_name: j.full_name, phone: j.phone, address: j.address, is_mustahik: j.is_mustahik },
  });
  showDialog.value = true;
};

const onSaved = (message: string) => {
  feedback.success(message);
  queryClient.invalidateQueries({ queryKey: queryKeys.jamaah.all });
  showDialog.value = false;
};

// Mutations
const createMutation = useMutation({
  mutationFn: jamaahService.createJamaah,
  onSuccess: () => onSaved('Jamaah terdaftar'),
  onError: (error: unknown) => feedback.error('Jamaah belum terdaftar', error),
});

const updateMutation = useMutation({
  mutationFn: (payload: { id: string; data: JamaahPayload }) => jamaahService.updateJamaah(payload.id, payload.data),
  onSuccess: () => onSaved('Data jamaah diperbarui'),
  onError: (error: unknown) => feedback.error('Data jamaah belum diperbarui', error),
});

const isSaving = computed(() => createMutation.isPending.value || updateMutation.isPending.value);

const submitForm = handleSubmit((values) => {
  if (editingId.value) {
    updateMutation.mutate({ id: editingId.value, data: values });
  } else {
    createMutation.mutate(values);
  }
});

const handleSearch = () => {
  appliedSearch.value = searchInput.value.trim();
  currentPage.value = 1;
};

const resetToFirstPage = () => {
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
          <p class="text-sm text-gray-500">Kelola registrasi warga, status Mustahik, dan skor kemiskinan.</p>
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

    <!-- Search & Filters -->
    <div class="flex flex-col md:flex-row md:items-center gap-3">
      <div class="flex items-center bg-white p-2 rounded-xl shadow-sm border border-gray-100 max-w-md flex-1">
        <Search class="w-5 h-5 text-gray-400 ml-2" />
        <input
          v-model="searchInput"
          @keyup.enter="handleSearch"
          type="text"
          placeholder="Cari nama jamaah..."
          class="flex-1 border-0 focus:ring-0 text-sm p-2 w-full"
        />
        <button @click="handleSearch" class="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors">
          Cari
        </button>
      </div>
      <label class="inline-flex items-center gap-2 bg-white px-3 py-2 rounded-xl shadow-sm border border-gray-100 text-sm text-gray-700">
        <input type="checkbox" v-model="mustahikOnly" @change="resetToFirstPage" class="h-4 w-4 rounded border-gray-300" />
        Hanya Mustahik
      </label>
      <select v-model="sort" @change="resetToFirstPage" class="bg-white px-3 py-2 rounded-xl shadow-sm border border-gray-100 text-sm text-gray-700">
        <option value="newest">Terbaru</option>
        <option value="score_desc">Skor tertinggi</option>
        <option value="name_asc">Nama (A-Z)</option>
      </select>
    </div>

    <!-- Table / Virtual Scroller -->
    <div class="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden flex flex-col h-[600px]">
      <div class="overflow-x-auto flex-1 flex flex-col relative">
        <div class="flex items-center text-xs text-gray-700 uppercase bg-gray-50 border-b min-w-[900px]">
          <div class="px-6 py-4 flex-[2]">Nama Lengkap</div>
          <div class="px-6 py-4 flex-1">No. HP</div>
          <div class="px-6 py-4 flex-[2]">Alamat</div>
          <div class="px-6 py-4 flex-1">Status Zakat</div>
          <div class="px-6 py-4 w-28">Skor</div>
          <div v-if="showActions" class="px-6 py-4 w-28 text-right">Aksi</div>
        </div>

        <div v-if="isLoading" class="p-8 text-center text-gray-500 border-b min-w-[900px]">
          Memuat data jamaah...
        </div>
        <div v-else-if="jamaahs.length === 0" class="p-8 text-center text-gray-500 border-b min-w-[900px]">
          Tidak ada jamaah ditemukan.
        </div>

        <RecycleScroller
          v-else
          class="flex-1 h-full overflow-y-auto min-w-[900px]"
          :items="jamaahs"
          :item-size="65"
          key-field="id"
          v-slot="{ item: j }"
        >
          <div class="flex items-center text-sm border-b hover:bg-gray-50 transition-colors w-full h-[65px]">
            <div class="px-6 flex-[2] font-medium text-gray-900 truncate">{{ j.full_name }}</div>
            <div class="px-6 flex-1 text-gray-600 truncate">{{ j.phone || '-' }}</div>
            <div class="px-6 flex-[2] text-gray-600 truncate max-w-xs">{{ j.address || '-' }}</div>
            <div class="px-6 flex-1 truncate">
                <span v-if="j.is_mustahik" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-purple-100 text-purple-700">
                  <ShieldAlert class="w-3 h-3 mr-1" /> Mustahik
                </span>
                <span v-else class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-700">
                  Umum (Muzakki)
                </span>
            </div>
            <div class="px-6 w-28">
              <span :class="['px-2 py-1 rounded-full text-xs font-semibold', scoreBadgeClass[scoreLevel(j.mustahik_score)]]">
                {{ j.mustahik_score.toFixed(1) }}
              </span>
            </div>
            <div v-if="showActions" class="px-6 w-28 flex justify-end gap-2">
                <button v-if="canSurvey" @click="surveyTarget = j" title="Survei Mustahik" class="text-purple-600 hover:text-purple-900 bg-purple-50 hover:bg-purple-100 p-2 rounded-lg transition-colors">
                  <ClipboardList class="w-4 h-4" />
                </button>
                <button v-if="canManage" @click="openEditDialog(j)" title="Edit" class="text-blue-600 hover:text-blue-900 bg-blue-50 hover:bg-blue-100 p-2 rounded-lg transition-colors">
                  <Edit2 class="w-4 h-4" />
                </button>
            </div>
          </div>
        </RecycleScroller>
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

  <!-- Mustahik Survey Dialog -->
  <MustahikSurveyDialog v-if="surveyTarget" :jamaah="surveyTarget" @close="surveyTarget = null" />

  <!-- Dialog Form -->
  <div v-if="showDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
        <h3 class="text-lg font-bold text-gray-900">{{ editingId ? 'Edit Data Jamaah' : 'Register Jamaah Baru' }}</h3>
        <button @click="showDialog = false" class="text-gray-400 hover:text-gray-600 rounded-full p-1 hover:bg-gray-200 transition-colors">&times;</button>
      </div>
      <form @submit="submitForm" novalidate>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nama Lengkap <span class="text-red-500">*</span></label>
            <input type="text" v-model="fullName" class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2.5 border" placeholder="Cth: Ahmad Abdullah" />
            <p v-if="errors.full_name" class="mt-1 text-xs text-red-600">{{ errors.full_name }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nomor HP / WhatsApp</label>
            <input type="text" v-model="phone" class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2.5 border" placeholder="Cth: 08123456789" />
            <p v-if="errors.phone" class="mt-1 text-xs text-red-600">{{ errors.phone }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Alamat Domisili</label>
            <textarea v-model="address" rows="3" class="block w-full rounded-lg border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm p-2.5 border" placeholder="Jalan, RT/RW, Kelurahan..."></textarea>
            <p v-if="errors.address" class="mt-1 text-xs text-red-600">{{ errors.address }}</p>
          </div>
          <div class="flex items-center mt-4 p-3 bg-purple-50 rounded-lg border border-purple-100">
            <input id="mustahik" type="checkbox" v-model="isMustahik" class="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded" />
            <label for="mustahik" class="ml-3 block text-sm font-medium text-purple-900">
              Tandai sebagai Mustahik (Penerima Zakat)
            </label>
          </div>
        </div>
        <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
          <button type="button" @click="showDialog = false" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors shadow-sm">
            Batal
          </button>
          <button type="submit" :disabled="isSaving" class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors shadow-sm flex items-center">
            <span v-if="isSaving" class="w-4 h-4 mr-2 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
            Simpan Data
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
