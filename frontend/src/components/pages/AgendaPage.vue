<script setup lang="ts">
import { ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { agendaService } from '../../api/services/agenda';
import type { CreateAgendaPayload, Agenda } from '../../types/agenda';
import { PlusIcon, XIcon, CalendarIcon, ClockIcon, MapPinIcon, EditIcon } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

const queryClient = useQueryClient();

const showAddDialog = ref(false);
const editingId = ref<string | null>(null);

const form = ref<CreateAgendaPayload>({
  title: '',
  description: '',
  day: 'Senin',
  time: '19:00',
  location: '',
  status: 'SCHEDULED'
});

const { data: agendas, isLoading } = useQuery({
  queryKey: ['agendas'],
  queryFn: agendaService.getAgendas
});

const createMutation = useMutation({
  mutationFn: agendaService.createAgenda,
  onSuccess: () => {
    toast.success('Agenda berhasil ditambahkan');
    queryClient.invalidateQueries({ queryKey: ['agendas'] });
    closeDialog();
  },
  onError: () => toast.error('Gagal menambahkan agenda')
});

const updateMutation = useMutation({
  mutationFn: ({ id, payload }: { id: string, payload: CreateAgendaPayload }) => agendaService.updateAgenda(id, payload),
  onSuccess: () => {
    toast.success('Agenda berhasil diperbarui');
    queryClient.invalidateQueries({ queryKey: ['agendas'] });
    closeDialog();
  },
  onError: () => toast.error('Gagal memperbarui agenda')
});

const openEditDialog = (agenda: Agenda) => {
  editingId.value = agenda.id;
  form.value = {
    title: agenda.title,
    description: agenda.description,
    day: agenda.day,
    time: agenda.time,
    location: agenda.location,
    status: agenda.status
  };
  showAddDialog.value = true;
};

const closeDialog = () => {
  showAddDialog.value = false;
  editingId.value = null;
  form.value = { title: '', description: '', day: 'Senin', time: '19:00', location: '', status: 'SCHEDULED' };
};

const submitAgenda = () => {
  const payload = {
    ...form.value,
  };
  if (editingId.value) {
    updateMutation.mutate({ id: editingId.value, payload });
  } else {
    createMutation.mutate(payload);
  }
};
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Agenda</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola jadwal kegiatan dan acara masjid</p>
      </div>
      <button @click="showAddDialog = true" class="px-4 py-2 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 focus:ring-4 focus:ring-emerald-100 transition-colors flex items-center gap-2">
        <PlusIcon class="w-5 h-5" />
        Tambah Agenda
      </button>
    </div>

    <!-- Agenda List -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-if="isLoading" v-for="i in 3" :key="i" class="bg-white rounded-xl shadow-sm border border-gray-100 p-6 animate-pulse">
        <div class="h-4 bg-gray-200 rounded w-3/4 mb-4"></div>
        <div class="h-3 bg-gray-200 rounded w-1/2 mb-2"></div>
        <div class="h-3 bg-gray-200 rounded w-5/6"></div>
      </div>
      <div v-else-if="agendas?.length === 0" class="col-span-full bg-white rounded-xl shadow-sm border border-gray-100 p-12 text-center">
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
            <button @click="openEditDialog(agenda)" class="text-blue-500 hover:text-blue-700 ml-1">
              <EditIcon class="w-4 h-4" />
            </button>
          </div>
          <span class="text-xs font-medium px-2 py-1 rounded-full bg-gray-100 text-gray-600">{{ agenda.status }}</span>
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
        <form @submit.prevent="submitAgenda" class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Judul Kegiatan</label>
            <input v-model="form.title" required type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Cth: Pengajian Rutin" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
            <textarea v-model="form.description" class="block w-full rounded-lg border-gray-300 p-2.5 border" rows="3" placeholder="Keterangan acara..."></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Hari</label>
              <select v-model="form.day" required class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="Senin">Senin</option>
                <option value="Selasa">Selasa</option>
                <option value="Rabu">Rabu</option>
                <option value="Kamis">Kamis</option>
                <option value="Jumat">Jumat</option>
                <option value="Sabtu">Sabtu</option>
                <option value="Minggu">Minggu</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Waktu</label>
              <input v-model="form.time" required type="time" class="block w-full rounded-lg border-gray-300 p-2.5 border" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Lokasi</label>
              <input v-model="form.location" type="text" class="block w-full rounded-lg border-gray-300 p-2.5 border" placeholder="Aula Masjid" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <select v-model="form.status" class="block w-full rounded-lg border-gray-300 p-2.5 border">
                <option value="SCHEDULED">Terjadwal</option>
                <option value="ONGOING">Berlangsung</option>
                <option value="COMPLETED">Selesai</option>
              </select>
            </div>
          </div>
          
          <div class="pt-4 flex gap-3">
            <button type="button" @click="closeDialog" class="flex-1 px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors">Batal</button>
            <button type="submit" :disabled="createMutation.isPending.value || updateMutation.isPending.value" class="flex-1 px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-colors disabled:opacity-50">
              {{ (createMutation.isPending.value || updateMutation.isPending.value) ? 'Menyimpan...' : 'Simpan Agenda' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
