<script setup lang="ts">
import { computed, watch } from 'vue';
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import { feedback } from '../../composables/useFeedback';
import { jamaahService } from '../../api/services/jamaah';
import { queryKeys } from '../../api/queryKeys';
import { mustahikSurveySchema, type MustahikSurveyPayload } from '../../schemas/jamaah';
import { houseStatusLabels, scoreBadgeClass, scoreLevel } from '../../lib/mustahik';
import type { HouseStatus, Jamaah } from '../../types/jamaah';

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
    feedback.success('Survei tersimpan', `Skor kemiskinan ${props.jamaah.full_name} kini ${survey.mustahik_score.toFixed(2)} dari 100.`);
    queryClient.invalidateQueries({ queryKey: queryKeys.jamaah.all });
    emit('close');
  },
  onError: (error: unknown) => {
    feedback.error('Survei belum disimpan', error);
  },
});

const onSubmit = handleSubmit((values) => submitMutation.mutate(values));
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
    <div class="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
        <div>
          <h3 class="text-lg font-bold text-gray-900">Survei Mustahik</h3>
          <p class="text-xs text-gray-500">{{ jamaah.full_name }}</p>
        </div>
        <button @click="emit('close')" class="text-gray-400 hover:text-gray-600 rounded-full p-1 hover:bg-gray-200 transition-colors">&times;</button>
      </div>

      <div v-if="isLoading" class="p-8 text-center text-sm text-gray-500">Memuat data survei...</div>

      <form v-else @submit="onSubmit" novalidate>
        <div class="p-6 space-y-4">
          <div class="flex items-center justify-between rounded-lg bg-gray-50 p-3 text-sm">
            <span class="text-gray-600">Skor saat ini</span>
            <span :class="['px-2 py-1 rounded-full text-xs font-semibold', scoreBadgeClass[scoreLevel(currentScore)]]">
              {{ currentScore.toFixed(2) }} / 100
            </span>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Penghasilan keluarga per bulan (Rp)</label>
            <input type="number" min="0" step="1" v-model.number="totalIncome" class="block w-full rounded-lg border-gray-300 shadow-sm sm:text-sm p-2.5 border" />
            <p v-if="errors.total_income" class="mt-1 text-xs text-red-600">{{ errors.total_income }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Jumlah tanggungan</label>
            <input type="number" min="0" max="30" step="1" v-model.number="dependentsCount" class="block w-full rounded-lg border-gray-300 shadow-sm sm:text-sm p-2.5 border" />
            <p v-if="errors.dependents_count" class="mt-1 text-xs text-red-600">{{ errors.dependents_count }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status tempat tinggal</label>
            <select v-model="houseStatus" class="block w-full rounded-lg border-gray-300 shadow-sm sm:text-sm p-2.5 border bg-white">
              <option v-for="status in houseStatuses" :key="status" :value="status">{{ houseStatusLabels[status] }}</option>
            </select>
            <p v-if="errors.house_status" class="mt-1 text-xs text-red-600">{{ errors.house_status }}</p>
          </div>
          <div class="flex items-center p-3 bg-purple-50 rounded-lg border border-purple-100">
            <input id="has-disability" type="checkbox" v-model="hasDisability" class="h-4 w-4 text-purple-600 border-gray-300 rounded" />
            <label for="has-disability" class="ml-3 block text-sm font-medium text-purple-900">Ada anggota keluarga penyandang disabilitas</label>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Catatan survei</label>
            <textarea v-model="criteriaNotes" rows="3" class="block w-full rounded-lg border-gray-300 shadow-sm sm:text-sm p-2.5 border"></textarea>
            <p v-if="errors.criteria_notes" class="mt-1 text-xs text-red-600">{{ errors.criteria_notes }}</p>
          </div>
        </div>
        <div class="px-6 py-4 border-t bg-gray-50 flex justify-end space-x-3">
          <button type="button" @click="emit('close')" class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors shadow-sm">
            Batal
          </button>
          <button type="submit" :disabled="submitMutation.isPending.value" class="px-4 py-2 text-sm font-medium text-white bg-purple-600 rounded-lg hover:bg-purple-700 disabled:opacity-50 transition-colors shadow-sm">
            {{ submitMutation.isPending.value ? 'Menyimpan...' : 'Simpan & Hitung Skor' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
