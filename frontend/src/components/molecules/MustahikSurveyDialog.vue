<script setup lang="ts">
import { computed, watch } from 'vue';
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { feedback } from '../../composables/useFeedback';
import { jamaahService } from '../../api/services/jamaah';
import { queryKeys } from '../../api/queryKeys';
import { mustahikSurveySchema, type MustahikSurveyPayload } from '../../schemas/jamaah';
import { houseStatusLabels, scoreTone } from '../../lib/mustahik';
import type { HouseStatus, Jamaah } from '../../types/jamaah';
import FormDialog from './FormDialog.vue';
import FormField from './FormField.vue';
import StatusBadge from './StatusBadge.vue';

const props = defineProps<{
  jamaah: Jamaah;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const queryClient = useQueryClient();
const houseStatuses = Object.keys(houseStatusLabels) as HouseStatus[];

const emptySurvey: MustahikSurveyPayload = {
  total_income: 0,
  dependents_count: 0,
  house_status: 'OWNED',
  has_disability: false,
  criteria_notes: '',
};

const { data: existingSurvey, isLoading } = useQuery({
  queryKey: computed(() => queryKeys.jamaah.survey(props.jamaah.id)),
  queryFn: () => jamaahService.getSurvey(props.jamaah.id),
});

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: toTypedSchema(mustahikSurveySchema),
  initialValues: emptySurvey,
});

const [totalIncome] = defineField('total_income');
const [dependentsCount] = defineField('dependents_count');
const [houseStatus] = defineField('house_status');
const [hasDisability] = defineField('has_disability');
const [criteriaNotes] = defineField('criteria_notes');

// Prefill the form with the previous survey once it has loaded.
watch(existingSurvey, (survey) => {
  if (!survey) return;
  resetForm({
    values: {
      total_income: survey.total_income,
      dependents_count: survey.dependents_count,
      house_status: survey.house_status,
      has_disability: survey.has_disability,
      criteria_notes: survey.criteria_notes,
    },
  });
}, { immediate: true });

const currentScore = computed(() => existingSurvey.value?.mustahik_score ?? props.jamaah.mustahik_score);

const submitMutation = useMutation({
  mutationFn: (payload: MustahikSurveyPayload) => jamaahService.submitSurvey(props.jamaah.id, payload),
  onSuccess: (survey) => {
    queryClient.invalidateQueries({ queryKey: queryKeys.jamaah.all });
    emit('close');
    feedback.success('Survei tersimpan', `Skor kebutuhan ${props.jamaah.full_name} kini ${survey.mustahik_score.toFixed(1)} dari 100.`);
  },
  onError: (error: unknown) => feedback.error('Survei belum disimpan', error),
});

const onSubmit = handleSubmit((values) => submitMutation.mutate(values));

const onOpenChange = (open: boolean) => {
  if (!open) emit('close');
};
</script>

<template>
  <FormDialog
    :open="true"
    title="Survei mustahik"
    :description="`Kondisi ekonomi ${jamaah.full_name}. Skor kebutuhan dihitung ulang saat survei disimpan.`"
    submit-label="Simpan survei"
    :submitting="submitMutation.isPending.value"
    @update:open="onOpenChange"
    @submit="onSubmit"
  >
    <p v-if="isLoading" class="py-6 text-center text-sm text-jati-muda">Memuat survei sebelumnya…</p>
    <template v-else>
      <div class="flex items-center justify-between rounded-lg border border-batu bg-kapur/60 px-4 py-3 text-sm">
        <span class="text-jati-muda">Skor saat ini</span>
        <StatusBadge :tone="scoreTone(currentScore)" class="tabular">{{ currentScore.toFixed(1) }} / 100</StatusBadge>
      </div>
      <FormField label="Penghasilan keluarga per bulan (Rp)" for="survey-income" :error="errors.total_income" required>
        <input id="survey-income" v-model.number="totalIncome" type="number" min="0" step="1" class="field-control tabular" :aria-invalid="Boolean(errors.total_income)" />
      </FormField>
      <FormField label="Jumlah tanggungan" for="survey-dependents" :error="errors.dependents_count" hint="Anggota keluarga yang ditanggung, tidak termasuk kepala keluarga." required>
        <input id="survey-dependents" v-model.number="dependentsCount" type="number" min="0" max="30" step="1" class="field-control tabular" :aria-invalid="Boolean(errors.dependents_count)" />
      </FormField>
      <FormField label="Tempat tinggal" for="survey-house" :error="errors.house_status" required>
        <select id="survey-house" v-model="houseStatus" class="field-control">
          <option v-for="status in houseStatuses" :key="status" :value="status">{{ houseStatusLabels[status] }}</option>
        </select>
      </FormField>
      <label class="flex items-start gap-3 rounded-lg border border-batu bg-kapur/60 p-3 text-sm">
        <input v-model="hasDisability" type="checkbox" class="mt-0.5 size-4 accent-sajadah" />
        <span class="font-medium text-jati">Ada anggota keluarga penyandang disabilitas</span>
      </label>
      <FormField label="Catatan survei" for="survey-notes" :error="errors.criteria_notes">
        <textarea id="survey-notes" v-model="criteriaNotes" rows="3" class="field-control" />
      </FormField>
    </template>
  </FormDialog>
</template>
