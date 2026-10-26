<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import { jamaahService } from '../../api/services/jamaah';
import { queryKeys } from '../../api/queryKeys';
import { scoreTone } from '../../lib/mustahik';
import type { Jamaah, JamaahListParams } from '../../types/jamaah';
import { Check, ChevronsUpDown, Loader2 } from 'lucide-vue-next';
import { useDebounceFn } from '@vueuse/core';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '../atoms/popover';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '../atoms/command';
import StatusBadge from './StatusBadge.vue';

const props = defineProps<{
  /** Selected jamaah id; empty or undefined when nothing is selected. */
  modelValue?: string;
  placeholder?: string;
  /** Only list registered mustahik, highest poverty score first (zakat targeting). */
  mustahikOnly?: boolean;
  muzakkiOnly?: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const open = ref(false);
const searchQuery = ref('');
const debouncedSearch = ref('');

const updateSearch = useDebounceFn((val: string) => {
  debouncedSearch.value = val;
}, 300);

watch(searchQuery, (val) => {
  updateSearch(val);
});

const params = computed<JamaahListParams>(() => ({
  page: 1,
  limit: 20,
  query: debouncedSearch.value || undefined,
  mustahik: props.mustahikOnly ? true : props.muzakkiOnly ? false : undefined,
  sort: props.mustahikOnly ? 'score_desc' : 'name_asc',
}));

const { data: searchResults, isLoading } = useQuery({
  queryKey: computed(() => queryKeys.jamaah.list(params.value)),
  queryFn: async () => (await jamaahService.getJamaahs(params.value)).items,
});

const selectedJamaah = ref<Jamaah | null>(null);

// Forget the selection when the parent clears the value (e.g. after a form reset).
watch(() => props.modelValue, (value) => {
  if (!value) selectedJamaah.value = null;
});

const handleSelect = (jamaah: Jamaah) => {
  selectedJamaah.value = jamaah;
  emit('update:modelValue', jamaah.id);
  open.value = false;
};
</script>

<template>
  <Popover v-model:open="open">
    <PopoverTrigger as-child>
      <button
        type="button"
        role="combobox"
        :aria-expanded="open"
        class="field-control flex items-center justify-between gap-2 text-left"
      >
        <span class="truncate" :class="selectedJamaah ? 'text-jati' : 'text-jati-muda/70'">
          {{ selectedJamaah ? selectedJamaah.full_name : (placeholder || 'Pilih jamaah…') }}
        </span>
        <ChevronsUpDown class="size-4 shrink-0 text-jati-muda" />
      </button>
    </PopoverTrigger>
    <PopoverContent class="w-(--reka-popover-trigger-width) min-w-72 rounded-xl border-batu p-0 shadow-lg" align="start">
      <Command>
        <CommandInput v-model="searchQuery" placeholder="Ketik nama jamaah…" />
        <CommandEmpty class="py-6 text-center text-sm text-jati-muda">
          <span v-if="isLoading" class="inline-flex items-center gap-2">
            <Loader2 class="size-4 animate-spin" /> Mencari…
          </span>
          <span v-else>Tidak ada jamaah yang cocok.</span>
        </CommandEmpty>
        <CommandList>
          <CommandGroup>
            <CommandItem
              v-for="jamaah in searchResults"
              :key="jamaah.id"
              :value="jamaah.full_name"
              class="cursor-pointer rounded-md py-2.5 text-sm"
              @select="handleSelect(jamaah)"
            >
              <Check class="mr-2 size-4" :class="modelValue === jamaah.id ? 'opacity-100 text-sajadah' : 'opacity-0'" />
              <span class="flex-1 truncate">{{ jamaah.full_name }}</span>
              <StatusBadge v-if="mustahikOnly" :tone="scoreTone(jamaah.mustahik_score)" class="tabular">
                {{ jamaah.mustahik_score.toFixed(1) }}
              </StatusBadge>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>
</template>
