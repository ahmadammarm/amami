<script setup lang="ts">
import { ref, computed } from 'vue';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { Search, Plus, Pencil, ClipboardList, Users } from 'lucide-vue-next';
import { jamaahService } from '@/api/services/jamaah';
import { queryKeys } from '@/api/queryKeys';
import { feedback } from '@/composables/useFeedback';
import { useAuthStore } from '@/stores/auth';
import { Permission } from '@/lib/permissions';
import { scoreTone } from '@/lib/mustahik';
import { jamaahSchema, type JamaahPayload } from '@/schemas/jamaah';
import type { Jamaah, JamaahListParams, JamaahSort } from '@/types/jamaah';
import PageHeader from '@/components/molecules/PageHeader.vue';
import DataPanel from '@/components/molecules/DataPanel.vue';
import PaginationBar from '@/components/molecules/PaginationBar.vue';
import EmptyState from '@/components/molecules/EmptyState.vue';
import StatusBadge from '@/components/molecules/StatusBadge.vue';
import FormDialog from '@/components/molecules/FormDialog.vue';
import FormField from '@/components/molecules/FormField.vue';
import MustahikSurveyDialog from '@/components/molecules/MustahikSurveyDialog.vue';
import { Button } from '@/components/atoms/button';

const authStore = useAuthStore();
const queryClient = useQueryClient();

const PAGE_SIZE = 10;
const currentPage = ref(1);
const searchInput = ref('');
const appliedSearch = ref('');
const mustahikOnly = ref(false);
const sort = ref<JamaahSort>('newest');

const listParams = computed<JamaahListParams>(() => ({
  page: currentPage.value,
  limit: PAGE_SIZE,
  query: appliedSearch.value || undefined,
  mustahik: mustahikOnly.value ? true : undefined,
  sort: sort.value,
}));

const { data: jamaahPage, isLoading } = useQuery({
  queryKey: computed(() => queryKeys.jamaah.list(listParams.value)),
  queryFn: () => jamaahService.getJamaahs(listParams.value),
  placeholderData: (previous) => previous,
});

const jamaahs = computed(() => jamaahPage.value?.items ?? []);
const isFiltered = computed(() => Boolean(appliedSearch.value) || mustahikOnly.value);

const canManage = computed(() => authStore.can(Permission.JamaahManage));
const canSurvey = computed(() => authStore.can(Permission.MustahikSurvey));

const applySearch = () => {
  appliedSearch.value = searchInput.value.trim();
  currentPage.value = 1;
};

const resetToFirstPage = () => {
  currentPage.value = 1;
};

// Create / edit
const showForm = ref(false);
const editing = ref<Jamaah | null>(null);
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

const openCreate = () => {
  editing.value = null;
  resetForm({ values: emptyForm });
  showForm.value = true;
};

const openEdit = (j: Jamaah) => {
  editing.value = j;
  resetForm({ values: { full_name: j.full_name, phone: j.phone, address: j.address, is_mustahik: j.is_mustahik } });
  showForm.value = true;
};

const onSaved = (title: string, name: string) => {
  queryClient.invalidateQueries({ queryKey: queryKeys.jamaah.all });
  showForm.value = false;
  feedback.success(title, name);
};

const createMutation = useMutation({
  mutationFn: jamaahService.createJamaah,
  onSuccess: (j) => onSaved('Jamaah terdaftar', `${j.full_name} sudah masuk daftar jamaah.`),
  onError: (error: unknown) => feedback.error('Jamaah belum terdaftar', error),
});

const updateMutation = useMutation({
  mutationFn: (payload: { id: string; data: JamaahPayload }) => jamaahService.updateJamaah(payload.id, payload.data),
  onSuccess: (j) => onSaved('Data jamaah diperbarui', `Perubahan untuk ${j.full_name} tersimpan.`),
  onError: (error: unknown) => feedback.error('Data jamaah belum diperbarui', error),
});

const isSaving = computed(() => createMutation.isPending.value || updateMutation.isPending.value);

const submitForm = handleSubmit((values) => {
  if (editing.value) {
    updateMutation.mutate({ id: editing.value.id, data: values });
  } else {
    createMutation.mutate(values);
  }
});
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Jamaah" description="Daftar warga dan jamaah, termasuk status mustahik dan skor kebutuhannya.">
      <template v-if="canManage" #actions>
        <Button @click="openCreate"><Plus /> Daftarkan jamaah</Button>
      </template>
    </PageHeader>

    <DataPanel>
      <template #actions>
        <form class="flex w-full flex-col gap-2 sm:flex-row sm:items-center" role="search" @submit.prevent="applySearch">
          <div class="relative flex-1 sm:w-72 sm:flex-none">
            <Search class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-jati-muda" />
            <input
              v-model="searchInput"
              type="search"
              placeholder="Cari nama jamaah"
              aria-label="Cari nama jamaah"
              class="field-control pl-9"
            />
          </div>
          <select v-model="sort" class="field-control sm:w-44" aria-label="Urutkan" @change="resetToFirstPage">
            <option value="newest">Terbaru</option>
            <option value="score_desc">Skor tertinggi</option>
            <option value="name_asc">Nama A–Z</option>
          </select>
          <label class="inline-flex items-center gap-2 whitespace-nowrap text-sm text-jati">
            <input v-model="mustahikOnly" type="checkbox" class="size-4 accent-sajadah" @change="resetToFirstPage" />
            Hanya mustahik
          </label>
        </form>
      </template>

      <div v-if="isLoading" class="p-8 text-center text-sm text-jati-muda">Memuat jamaah…</div>
      <EmptyState
        v-else-if="jamaahs.length === 0"
        :icon="Users"
        :title="isFiltered ? 'Tidak ada jamaah yang cocok' : 'Belum ada jamaah terdaftar'"
        :description="isFiltered ? 'Ubah kata kunci atau matikan filter mustahik.' : 'Daftarkan jamaah agar bisa dicatat sebagai muzakki, mustahik, atau shohibul qurban.'"
      />
      <div v-else class="overflow-x-auto">
        <table class="tabel min-w-[820px]">
          <thead>
            <tr>
              <th>Nama</th>
              <th>No. HP</th>
              <th>Alamat</th>
              <th>Status</th>
              <th class="text-right">Skor</th>
              <th class="w-24"><span class="sr-only">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="j in jamaahs" :key="j.id">
              <td class="font-medium">{{ j.full_name }}</td>
              <td class="tabular text-jati-muda">{{ j.phone || '–' }}</td>
              <td class="max-w-xs truncate text-jati-muda">{{ j.address || '–' }}</td>
              <td>
                <StatusBadge :tone="j.is_mustahik ? 'sajadah' : 'netral'">{{ j.is_mustahik ? 'Mustahik' : 'Umum' }}</StatusBadge>
              </td>
              <td class="text-right">
                <StatusBadge :tone="scoreTone(j.mustahik_score)" class="tabular">{{ j.mustahik_score.toFixed(1) }}</StatusBadge>
              </td>
              <td>
                <div class="flex justify-end gap-1">
                  <button
                    v-if="canSurvey"
                    type="button"
                    class="rounded-md p-2 text-jati-muda hover:bg-kapur hover:text-sajadah"
                    :aria-label="`Survei mustahik ${j.full_name}`"
                    title="Survei mustahik"
                    @click="surveyTarget = j"
                  >
                    <ClipboardList class="size-4" />
                  </button>
                  <button
                    v-if="canManage"
                    type="button"
                    class="rounded-md p-2 text-jati-muda hover:bg-kapur hover:text-sajadah"
                    :aria-label="`Ubah data ${j.full_name}`"
                    title="Ubah data"
                    @click="openEdit(j)"
                  >
                    <Pencil class="size-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <template v-if="jamaahPage && jamaahPage.total > 0" #footer>
        <PaginationBar v-model:page="currentPage" :page-size="PAGE_SIZE" :total="jamaahPage.total" :total-pages="jamaahPage.total_pages" noun="jamaah" />
      </template>
    </DataPanel>

    <MustahikSurveyDialog v-if="surveyTarget" :jamaah="surveyTarget" @close="surveyTarget = null" />

    <FormDialog
      v-model:open="showForm"
      :title="editing ? 'Ubah data jamaah' : 'Daftarkan jamaah'"
      :submit-label="editing ? 'Simpan perubahan' : 'Daftarkan'"
      :submitting="isSaving"
      @submit="submitForm"
    >
      <FormField label="Nama lengkap" for="jamaah-name" :error="errors.full_name" required>
        <input id="jamaah-name" v-model="fullName" type="text" class="field-control" :aria-invalid="Boolean(errors.full_name)" />
      </FormField>
      <FormField label="No. HP / WhatsApp" for="jamaah-phone" :error="errors.phone">
        <input id="jamaah-phone" v-model="phone" type="tel" class="field-control" placeholder="08xxxxxxxxxx" />
      </FormField>
      <FormField label="Alamat" for="jamaah-address" :error="errors.address">
        <textarea id="jamaah-address" v-model="address" rows="3" class="field-control" placeholder="Jalan, RT/RW, kelurahan" />
      </FormField>
      <label class="flex items-start gap-3 rounded-lg border border-batu bg-kapur/60 p-3 text-sm">
        <input v-model="isMustahik" type="checkbox" class="mt-0.5 size-4 accent-sajadah" />
        <span>
          <span class="font-medium text-jati">Mustahik (penerima zakat)</span>
          <span class="block text-jati-muda">Hanya mustahik yang dapat dipilih saat menyalurkan zakat.</span>
        </span>
      </label>
    </FormDialog>
  </div>
</template>
