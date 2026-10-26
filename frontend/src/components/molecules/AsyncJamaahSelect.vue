<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import { jamaahService } from '../../api/services/jamaah';
import { queryKeys } from '../../api/queryKeys';
import { scoreBadgeClass, scoreLevel } from '../../lib/mustahik';
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
import Button from '../atoms/button/Button.vue';

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
      <Button
        variant="outline"
        role="combobox"
        :aria-expanded="open"
        class="w-full justify-between font-normal h-11 border-gray-300 rounded-lg px-3"
      >
        <span class="truncate">
          {{ selectedJamaah ? selectedJamaah.full_name : (placeholder || 'Select Jamaah...') }}
        </span>
        <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
      </Button>
    </PopoverTrigger>
    <PopoverContent class="w-full p-0 shadow-xl border-gray-100 rounded-xl" align="start">
      <Command>
        <CommandInput
          placeholder="Search jamaah by name..."
          v-model="searchQuery"
          class="h-11 border-none focus:ring-0"
        />
        <CommandEmpty class="py-6 text-center text-sm text-gray-500">
          <span v-if="isLoading">
            <Loader2 class="w-4 h-4 animate-spin mx-auto mb-2" />
            Searching...
          </span>
          <span v-else>No jamaah found.</span>
        </CommandEmpty>
        <CommandList>
          <CommandGroup>
            <CommandItem
              v-for="jamaah in searchResults"
              :key="jamaah.id"
              :value="jamaah.full_name"
              @select="handleSelect(jamaah)"
              class="cursor-pointer py-2.5 rounded-md text-sm"
            >
              <Check
                :class="[
                  'mr-2 h-4 w-4',
                  modelValue === jamaah.id ? 'opacity-100' : 'opacity-0'
                ]"
              />
              {{ jamaah.full_name }}
              <span
                v-if="mustahikOnly"
                :class="['ml-auto px-2 py-0.5 rounded-full text-[10px] font-semibold', scoreBadgeClass[scoreLevel(jamaah.mustahik_score)]]"
              >
                {{ jamaah.mustahik_score.toFixed(1) }}
              </span>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>
</template>
