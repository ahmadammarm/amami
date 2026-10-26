<script setup lang="ts">
import { computed } from 'vue';
import { ChevronLeft, ChevronRight } from 'lucide-vue-next';
import { formatNumber } from '../../lib/format';

const props = defineProps<{
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
  /** What is being counted, e.g. "transaksi". */
  noun: string;
}>();

const emit = defineEmits<{ (e: 'update:page', value: number): void }>();

const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.pageSize + 1));
const to = computed(() => Math.min(props.page * props.pageSize, props.total));
</script>

<template>
  <div class="flex flex-col gap-3 text-sm text-jati-muda sm:flex-row sm:items-center sm:justify-between">
    <p class="tabular">
      {{ formatNumber(from) }}–{{ formatNumber(to) }} dari {{ formatNumber(total) }} {{ noun }}
    </p>
    <div class="flex gap-2">
      <button
        type="button"
        class="inline-flex items-center gap-1 rounded-lg border border-batu-tua bg-white px-3 py-1.5 font-medium text-jati transition-colors hover:border-sajadah disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="page <= 1"
        @click="emit('update:page', page - 1)"
      >
        <ChevronLeft class="size-4" /> Sebelumnya
      </button>
      <button
        type="button"
        class="inline-flex items-center gap-1 rounded-lg border border-batu-tua bg-white px-3 py-1.5 font-medium text-jati transition-colors hover:border-sajadah disabled:cursor-not-allowed disabled:opacity-40"
        :disabled="page >= totalPages"
        @click="emit('update:page', page + 1)"
      >
        Berikutnya <ChevronRight class="size-4" />
      </button>
    </div>
  </div>
</template>
