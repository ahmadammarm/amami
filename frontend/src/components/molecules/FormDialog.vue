<script setup lang="ts">
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '../atoms/dialog';
import { Button } from '../atoms/button';

withDefaults(defineProps<{
  open: boolean;
  title: string;
  description?: string;
  submitLabel: string;
  submitting?: boolean;
  /** `danger` styles the submit button for destructive actions. */
  tone?: 'default' | 'danger';
  size?: 'md' | 'lg';
}>(), {
  submitting: false,
  tone: 'default',
  size: 'md',
});

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void;
  (e: 'submit', event: Event): void;
}>();
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent
      class="max-h-[90vh] gap-0 overflow-hidden rounded-2xl border-batu p-0 shadow-lg"
      :class="size === 'lg' ? 'sm:max-w-2xl' : 'sm:max-w-md'"
    >
      <form novalidate class="flex max-h-[90vh] flex-col" @submit.prevent="emit('submit', $event)">
        <div class="border-b border-batu px-6 pt-6 pb-4 pr-12">
          <DialogTitle class="font-display text-xl font-semibold tracking-wide text-jati">{{ title }}</DialogTitle>
          <DialogDescription v-if="description" class="mt-1 text-sm text-jati-muda">{{ description }}</DialogDescription>
        </div>

        <div class="flex-1 space-y-4 overflow-y-auto px-6 py-5">
          <slot />
        </div>

        <div class="flex justify-end gap-2 border-t border-batu bg-kapur/60 px-6 py-4">
          <Button type="button" variant="outline" :disabled="submitting" @click="emit('update:open', false)">Batal</Button>
          <Button type="submit" :variant="tone === 'danger' ? 'destructive' : 'default'" :disabled="submitting">
            {{ submitting ? 'Menyimpan…' : submitLabel }}
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
