<script setup lang="ts">
import { computed, ref } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { Plus, Clock, MapPin, Pencil, Trash2, Images, CalendarRange } from 'lucide-vue-next';
import { agendaService } from '../../api/services/agenda';
import { queryKeys } from '../../api/queryKeys';
import { feedback } from '../../composables/useFeedback';
import { useAuthStore } from '../../stores/auth';
import { Permission } from '../../lib/permissions';
import { agendaSchema, type AgendaPayload } from '../../schemas/agenda';
import type { Agenda, AgendaStatus } from '../../types/agenda';
import type { BadgeTone } from '../../types/ui';
import PageHeader from '../molecules/PageHeader.vue';
import EmptyState from '../molecules/EmptyState.vue';
import StatusBadge from '../molecules/StatusBadge.vue';
import FormDialog from '../molecules/FormDialog.vue';
import FormField from '../molecules/FormField.vue';
import AgendaDocumentationDialog from '../molecules/AgendaDocumentationDialog.vue';
import { Button } from '../atoms/button';

const DAYS = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu'];

const AGENDA_STATUS: Record<AgendaStatus, { label: string; tone: BadgeTone }> = {
  SCHEDULED: { label: 'Terjadwal', tone: 'sajadah' },
  ONGOING: { label: 'Berlangsung', tone: 'kuningan' },
  COMPLETED: { label: 'Selesai', tone: 'netral' },
  CANCELLED: { label: 'Dibatalkan', tone: 'bata' },
};

const authStore = useAuthStore();
const queryClient = useQueryClient();

const canManage = computed(() => authStore.can(Permission.AgendaManage));

const { data: agendas, isLoading } = useQuery({
  queryKey: queryKeys.agenda.list,
  queryFn: agendaService.getAgendas,
});

/** Week order (Senin first), then time, so the list reads like a weekly timetable. */
const sortedAgendas = computed(() =>
  [...(agendas.value ?? [])].sort((a, b) => {
    const byDay = DAYS.indexOf(a.day) - DAYS.indexOf(b.day);
    return byDay !== 0 ? byDay : a.time.localeCompare(b.time);
  }),
);

const refreshAgendas = () => {
  queryClient.invalidateQueries({ queryKey: queryKeys.agenda.all });
  queryClient.invalidateQueries({ queryKey: queryKeys.dashboard.metrics });
};

// Create / edit
const showForm = ref(false);
const editing = ref<Agenda | null>(null);
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

const saveMutation = useMutation({
  mutationFn: (payload: AgendaPayload) =>
    editing.value ? agendaService.updateAgenda(editing.value.id, payload) : agendaService.createAgenda(payload),
  onSuccess: (agenda) => {
    const wasEditing = editing.value !== null;
    refreshAgendas();
    showForm.value = false;
    feedback.success(wasEditing ? 'Agenda diperbarui' : 'Agenda ditambahkan', `${agenda.title}, ${agenda.day} pukul ${agenda.time}.`);
  },
  onError: (error: unknown) => feedback.error('Agenda belum disimpan', error),
});

const openForm = (agenda: Agenda | null) => {
  editing.value = agenda;
  resetForm({
    values: agenda
      ? {
        title: agenda.title,
        description: agenda.description,
        day: agenda.day,
        time: agenda.time,
        location: agenda.location,
        status: agenda.status,
      }
      : emptyForm,
  });
  showForm.value = true;
};

const submitForm = handleSubmit((values) => saveMutation.mutate(values));

const deleteMutation = useMutation({
  mutationFn: (id: string) => agendaService.deleteAgenda(id),
  onSuccess: () => {
    refreshAgendas();
    feedback.success('Agenda dihapus');
  },
  onError: (error: unknown) => feedback.error('Agenda belum dihapus', error),
});

const confirmDelete = async (agenda: Agenda) => {
  const confirmed = await feedback.confirm({
    title: 'Hapus agenda?',
    message: `${agenda.title} akan dihapus dari jadwal. Dokumentasinya tetap tersimpan untuk arsip.`,
    confirmLabel: 'Hapus agenda',
    danger: true,
  });
  if (confirmed) deleteMutation.mutate(agenda.id);
};
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Agenda" description="Jadwal kajian dan kegiatan rutin masjid, beserta dokumentasinya.">
      <template v-if="canManage" #actions>
        <Button @click="openForm(null)"><Plus /> Tambah agenda</Button>
      </template>
    </PageHeader>

    <div v-if="isLoading" class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <div v-for="i in 3" :key="i" class="h-44 animate-pulse rounded-xl bg-batu/60" />
    </div>

    <div v-else-if="sortedAgendas.length === 0" class="rounded-xl border border-batu bg-white">
      <EmptyState
        :icon="CalendarRange"
        title="Belum ada agenda"
        :description="canManage ? 'Tambahkan kajian atau kegiatan rutin agar jamaah dan pengurus tahu jadwalnya.' : 'Jadwal kegiatan yang dibuat pengurus akan tampil di sini.'"
      >
        <Button v-if="canManage" @click="openForm(null)"><Plus /> Tambah agenda</Button>
      </EmptyState>
    </div>

    <ul v-else class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      <li
        v-for="agenda in sortedAgendas"
        :key="agenda.id"
        class="flex flex-col rounded-xl border border-batu bg-white p-5 shadow-sm"
        :class="{ 'opacity-70': agenda.status === 'CANCELLED' || agenda.status === 'COMPLETED' }"
      >
        <div class="flex items-start justify-between gap-3">
          <p class="text-xs font-semibold uppercase tracking-[0.14em] text-kuningan">{{ agenda.day }}</p>
          <StatusBadge :tone="AGENDA_STATUS[agenda.status].tone">{{ AGENDA_STATUS[agenda.status].label }}</StatusBadge>
        </div>
        <h2 class="mt-2 text-lg font-semibold text-jati">{{ agenda.title }}</h2>
        <p v-if="agenda.description" class="mt-1 line-clamp-2 text-sm text-jati-muda">{{ agenda.description }}</p>
        <dl class="mt-4 space-y-1.5 text-sm text-jati-muda">
          <div class="flex items-center gap-2">
            <dt><Clock class="size-4 text-sajadah" /><span class="sr-only">Waktu</span></dt>
            <dd class="tabular">{{ agenda.time }}</dd>
          </div>
          <div class="flex items-center gap-2">
            <dt><MapPin class="size-4 text-sajadah" /><span class="sr-only">Lokasi</span></dt>
            <dd>{{ agenda.location || 'Lokasi belum ditentukan' }}</dd>
          </div>
        </dl>
        <div class="mt-5 flex items-center justify-between border-t border-batu pt-4">
          <button type="button" class="inline-flex items-center gap-1.5 text-sm font-medium text-sajadah hover:text-sajadah-tua" @click="documentationTarget = agenda">
            <Images class="size-4" /> Dokumentasi
          </button>
          <div v-if="canManage" class="flex gap-1">
            <button type="button" class="rounded-md p-2 text-jati-muda hover:bg-kapur hover:text-sajadah" :aria-label="`Ubah ${agenda.title}`" title="Ubah" @click="openForm(agenda)">
              <Pencil class="size-4" />
            </button>
            <button type="button" class="rounded-md p-2 text-jati-muda hover:bg-bata-muda hover:text-bata" :aria-label="`Hapus ${agenda.title}`" title="Hapus" @click="confirmDelete(agenda)">
              <Trash2 class="size-4" />
            </button>
          </div>
        </div>
      </li>
    </ul>

    <AgendaDocumentationDialog
      v-if="documentationTarget"
      :agenda="documentationTarget"
      :can-manage="canManage"
      @close="documentationTarget = null"
    />

    <FormDialog
      v-model:open="showForm"
      :title="editing ? 'Ubah agenda' : 'Tambah agenda'"
      :submit-label="editing ? 'Simpan perubahan' : 'Tambah agenda'"
      :submitting="saveMutation.isPending.value"
      @submit="submitForm"
    >
      <FormField label="Nama kegiatan" for="agenda-title" :error="errors.title" required>
        <input id="agenda-title" v-model="title" type="text" class="field-control" placeholder="Contoh: Kajian tafsir ba'da Maghrib" :aria-invalid="Boolean(errors.title)" />
      </FormField>
      <FormField label="Deskripsi" for="agenda-description" :error="errors.description">
        <textarea id="agenda-description" v-model="description" rows="3" class="field-control" placeholder="Pemateri, tema, atau catatan untuk jamaah" />
      </FormField>
      <div class="grid grid-cols-2 gap-4">
        <FormField label="Hari" for="agenda-day" :error="errors.day" required>
          <select id="agenda-day" v-model="day" class="field-control">
            <option v-for="d in DAYS" :key="d" :value="d">{{ d }}</option>
          </select>
        </FormField>
        <FormField label="Waktu" for="agenda-time" :error="errors.time" required>
          <input id="agenda-time" v-model="time" type="time" class="field-control tabular" />
        </FormField>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <FormField label="Lokasi" for="agenda-location" :error="errors.location">
          <input id="agenda-location" v-model="location" type="text" class="field-control" placeholder="Ruang utama" />
        </FormField>
        <FormField label="Status" for="agenda-status" required>
          <select id="agenda-status" v-model="status" class="field-control">
            <option v-for="(s, value) in AGENDA_STATUS" :key="value" :value="value">{{ s.label }}</option>
          </select>
        </FormField>
      </div>
    </FormDialog>
  </div>
</template>
