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
      class="max-h-[92dvh] gap-0 overflow-hidden rounded-2xl border-batu p-0 shadow-lg"
      :class="size === 'lg' ? 'sm:max-w-2xl' : 'sm:max-w-md'"
    >
      <form novalidate class="flex max-h-[92dvh] flex-col" @submit.prevent="emit('submit', $event)">
        <div class="shrink-0 border-b border-batu px-5 pt-5 pb-4 pr-12 sm:px-6 sm:pt-6">
          <DialogTitle class="font-display text-xl font-semibold tracking-wide text-jati">{{ title }}</DialogTitle>
          <DialogDescription v-if="description" class="mt-1 text-sm text-jati-muda">{{ description }}</DialogDescription>
        </div>

        <div class="scrollbar-tipis flex-1 space-y-4 overflow-y-auto overscroll-contain px-5 py-5 sm:px-6">
          <slot />
        </div>

        <div class="flex shrink-0 flex-col-reverse gap-2 border-t border-batu bg-kapur/60 px-5 py-4 sm:flex-row sm:justify-end sm:px-6">
          <Button type="button" variant="outline" class="w-full sm:w-auto" :disabled="submitting" @click="emit('update:open', false)">Batal</Button>
          <Button type="submit" class="w-full sm:w-auto" :variant="tone === 'danger' ? 'destructive' : 'default'" :disabled="submitting">
            {{ submitting ? 'Menyimpan…' : submitLabel }}
          </Button>
        </div>
      </form>
    </DialogContent>
  </Dialog>
</template>
