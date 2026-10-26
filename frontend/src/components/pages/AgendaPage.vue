<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { agendaService } from '../../api/services/agenda';
import { feedback } from '../../composables/useFeedback';
import { queryKeys } from '../../api/queryKeys';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import { agendaSchema, type AgendaPayload } from '../../schemas/agenda';
import type { Agenda } from '../../types/agenda';
import AgendaDocumentationDialog from '../molecules/AgendaDocumentationDialog.vue';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '../atoms/alert-dialog';
import { PlusIcon, XIcon, CalendarIcon, ClockIcon, MapPinIcon, EditIcon, Trash2Icon, ImagesIcon } from 'lucide-vue-next';

const authStore = useAuthStore();
const queryClient = useQueryClient();

const canManage = computed(() => authStore.can(Permission.AgendaManage));

const showAddDialog = ref(false);
const editingId = ref<string | null>(null);
const deleteTarget = ref<Agenda | null>(null);
const documentationTarget = ref<Agenda | null>(null);

const emptyForm: AgendaPayload = {
  title: '',
  description: '',
  day: 'Senin',
  time: '19:00',
  location: '',
  status: 'SCHEDULED',
};

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: toTypedSchema(agendaSchema),
  initialValues: emptyForm,
});

const [title] = defineField('title');
const [description] = defineField('description');
const [day] = defineField('day');
const [time] = defineField('time');
const [location] = defineField('location');
const [status] = defineField('status');

const { data: agendas, isLoading } = useQuery({
  queryKey: queryKeys.agenda.list,
  queryFn: agendaService.getAgendas
});

const refreshAgendas = () => queryClient.invalidateQueries({ queryKey: queryKeys.agenda.all });

const createMutation = useMutation({
  mutationFn: agendaService.createAgenda,
  onSuccess: () => {
    feedback.success('Agenda ditambahkan');
    refreshAgendas();
    closeDialog();
  },
  onError: (error: unknown) => feedback.error('Agenda belum ditambahkan', error)
});

const updateMutation = useMutation({
  mutationFn: ({ id, payload }: { id: string, payload: AgendaPayload }) => agendaService.updateAgenda(id, payload),
  onSuccess: () => {
    feedback.success('Agenda diperbarui');
    refreshAgendas();
    closeDialog();
  },
  onError: (error: unknown) => feedback.error('Agenda belum diperbarui', error)
});

const deleteMutation = useMutation({
  mutationFn: (id: string) => agendaService.deleteAgenda(id),
  onSuccess: () => {
    feedback.success('Agenda dihapus');
    refreshAgendas();
    deleteTarget.value = null;
  },
  onError: (error: unknown) => feedback.error('Agenda belum dihapus', error)
});

const isSaving = computed(() => createMutation.isPending.value || updateMutation.isPending.value);

const openCreateDialog = () => {
  editingId.value = null;
  resetForm({ values: emptyForm });
  showAddDialog.value = true;
};

const openEditDialog = (agenda: Agenda) => {
  editingId.value = agenda.id;
  resetForm({
    values: {
      title: agenda.title,
      description: agenda.description,
      day: agenda.day,
      time: agenda.time,
      location: agenda.location,
      status: agenda.status,
    },
  });
  showAddDialog.value = true;
};

const closeDialog = () => {
  showAddDialog.value = false;
  editingId.value = null;
};

const submitAgenda = handleSubmit((values) => {
  if (editingId.value) {
    updateMutation.mutate({ id: editingId.value, payload: values });
  } else {
    createMutation.mutate(values);
  }
});

const confirmDelete = () => {
  if (deleteTarget.value) deleteMutation.mutate(deleteTarget.value.id);
};

const onDeleteDialogChange = (open: boolean) => {
  if (!open) deleteTarget.value = null;
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Agenda</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola jadwal kegiatan dan acara masjid</p>
      </div>
      <button v-if="canManage" @click="openCreateDialog" class="px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 focus:ring-4 focus:ring-emerald-100 transition-colors flex items-center gap-2">
        <PlusIcon class="w-5 h-5" />
        Tambah Agenda
      </button>
    </div>

    <!-- Agenda List -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <template v-if="isLoading">
        <div v-for="i in 3" :key="i" class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 animate-pulse">
          <div class="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2 mb-2"></div>
          <div class="h-3 bg-gray-200 rounded w-5/6"></div>
        </div>
      </template>
      <div v-else-if="!agendas || agendas.length === 0" class="col-span-full bg-white rounded-xl shadow-sm border border-gray-100 p-12 text-center">
        <CalendarIcon class="w-12 h-12 mx-auto text-gray-300 mb-3" />
        <h3 class="font-semibold text-gray-900">Belum ada agenda</h3>
        <p class="text-sm text-gray-500">Tambahkan agenda kegiatan baru</p>
      </div>
      <div v-else v-for="agenda in agendas" :key="agenda.id" class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 hover:shadow-md transition-shadow relative overflow-hidden">
        <div class="absolute top-0 right-0 w-2 h-full" :class="[
          agenda.status === 'SCHEDULED' ? 'bg-blue-500' :
          agenda.status === 'ONGOING' ? 'bg-emerald-500' :
          agenda.status === 'COMPLETED' ? 'bg-gray-400' : 'bg-red-500'
        ]"></div>
        <div class="flex justify-between items-start mb-4">
          <div class="flex items-center gap-2">
            <h3 class="font-bold text-gray-900">{{ agenda.title }}</h3>
            <button v-if="canManage" @click="openEditDialog(agenda)" class="text-blue-500 hover:text-blue-700 ml-1" title="Edit">
              <EditIcon class="w-4 h-4" />
            </button>
            <button v-if="canManage" @click="deleteTarget = agenda" class="text-red-500 hover:text-red-700" title="Hapus">
              <Trash2Icon class="w-4 h-4" />
            </button>
          </div>
          <span class="text-xs font-medium px-2 py-1 rounded-full bg-gray-100 text-gray-600 mr-2">{{ agenda.status }}</span>
        </div>
        <p class="text-sm text-gray-600 mb-4 line-clamp-2">{{ agenda.description }}</p>
        <div class="space-y-2 text-sm text-gray-500">
          <div class="flex items-center gap-1.5 mt-1 text-gray-500">
            <CalendarIcon class="w-4 h-4" />
            <span>{{ agenda.day }}</span>
            <ClockIcon class="w-4 h-4 ml-2" />
            <span>{{ agenda.time }}</span>
          </div>
          <div class="flex items-center gap-2">
            <MapPinIcon class="w-4 h-4 text-emerald-600" />
            <span>{{ agenda.location || 'Lokasi tidak ditentukan' }}</span>
          </div>
        </div>
        <button @click="documentationTarget = agenda" class="mt-4 inline-flex items-center gap-1.5 text-sm font-medium text-emerald-700 hover:text-emerald-900">
          <ImagesIcon class="w-4 h-4" />
          Dokumentasi
        </button>
      </div>
    </div>

    <!-- Add Agenda Dialog -->
    <div v-if="showAddDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-gray-900/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="font-bold text-gray-900">{{ editingId ? 'Edit Agenda' : 'Tambah Agenda Baru' }}</h3>
          <button @click="closeDialog" class="text-gray-400 hover:text-gray-600">
            <XIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit="submitAgenda" class="p-6 space-y-4" novalidate>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Judul Kegiatan</label>
            <input v-model="title" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Cth: Pengajian Rutin" />
            <p v-if="errors.title" class="mt-1 text-xs text-red-600">{{ errors.title }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
            <textarea v-model="description" class="block w-full rounded-lg border-gray-300 p-2.5 border" rows="3" placeholder="Keterangan acara..."></textarea>
            <p v-if="errors.description" class="mt-1 text-xs text-red-600">{{ errors.description }}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Hari</label>
              <select v-model="day" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="Senin">Senin</option>
                <option value="Selasa">Selasa</option>
                <option value="Rabu">Rabu</option>
                <option value="Kamis">Kamis</option>
                <option value="Jumat">Jumat</option>
                <option value="Sabtu">Sabtu</option>
                <option value="Minggu">Minggu</option>
              </select>
              <p v-if="errors.day" class="mt-1 text-xs text-red-600">{{ errors.day }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Waktu</label>
              <input v-model="time" type="time" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
              <p v-if="errors.time" class="mt-1 text-xs text-red-600">{{ errors.time }}</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi</label>
              <input v-model="location" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Aula Masjid" />
              <p v-if="errors.location" class="mt-1 text-xs text-red-600">{{ errors.location }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select v-model="status" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="SCHEDULED">Terjadwal</option>
                <option value="ONGOING">Berlangsung</option>
                <option value="COMPLETED">Selesai</option>
                <option value="CANCELLED">Dibatalkan</option>
              </select>
            </div>
          </div>

          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="isSaving" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ isSaving ? 'Menyimpan...' : 'Simpan Agenda' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Documentation Dialog -->
    <AgendaDocumentationDialog
      v-if="documentationTarget"
      :agenda="documentationTarget"
      :can-manage="canManage"
      @close="documentationTarget = null"
    />

    <!-- Delete Confirmation -->
    <AlertDialog :open="deleteTarget !== null" @update:open="onDeleteDialogChange">
      <AlertDialogContent class="rounded-2xl max-w-sm">
        <AlertDialogHeader>
          <AlertDialogTitle>Hapus agenda?</AlertDialogTitle>
          <AlertDialogDescription>
            "{{ deleteTarget?.title }}" akan dihapus dari jadwal. Dokumentasinya tetap tersimpan untuk arsip.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Batal</AlertDialogCancel>
          <AlertDialogAction @click="confirmDelete" :disabled="deleteMutation.isPending.value" class="bg-red-500 hover:bg-red-600">
            Hapus
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
