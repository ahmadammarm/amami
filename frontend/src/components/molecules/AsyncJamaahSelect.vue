<script setup lang="ts">
import { ref, watch } from 'vue';
import { useQuery } from '@tanstack/vue-query';
import apiClient from '../../api/client';
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
  modelValue: string;
  placeholder?: string;
  mustahikOnly?: boolean;
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

// Fetch Jamaah asynchronously based on search query
const { data: searchResults, isLoading } = useQuery({
  queryKey: ['jamaah-search', debouncedSearch],
  queryFn: async () => {
    // Only search if there's a query or just load the first 20
    const query = debouncedSearch.value ? `&search=${encodeURIComponent(debouncedSearch.value)}` : '';
    const mustahikParam = props.mustahikOnly ? '&is_mustahik=true' : '';
    const response = await apiClient.get(`/jamaah?page=1&limit=20${query}${mustahikParam}`);
    // Handle both old and new response shapes for safety
    const data = response.data.data;
    return data.items || data || [];
  },
});

const selectedJamaah = ref<any>(null);

// When modelValue changes externally or we load initial data, we might need to fetch the single user if not in list
// For simplicity in this demo, we assume the user selects from the list.
const handleSelect = (jamaah: any) => {
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
              <span class="text-xs text-gray-400 ml-2">({{ jamaah.nik }})</span>
            </CommandItem>
          </CommandGroup>
        </CommandList>
      </Command>
    </PopoverContent>
  </Popover>
</template>
