<script setup lang="ts">
import { computed, type Component } from 'vue';
import { CheckIcon, TriangleAlertIcon, CircleHelpIcon } from 'lucide-vue-next';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '../atoms/dialog';
import { Button } from '../atoms/button';
import { ArchBadge, type ArchTone } from '../atoms/arch';
import { useFeedbackQueue } from '../../composables/useFeedback';

const { current, settle } = useFeedbackQueue();

const appearance = computed<{ tone: ArchTone; icon: Component }>(() => {
  const request = current.value;
  if (!request || request.kind === 'success') return { tone: 'sajadah', icon: CheckIcon };
  if (request.kind === 'error' || request.danger) return { tone: 'bata', icon: TriangleAlertIcon };
  return { tone: 'kuningan', icon: CircleHelpIcon };
});

const onOpenChange = (open: boolean) => {
  if (!open) settle(false);
};
</script>

<template>
  <Dialog :open="current !== null" @update:open="onOpenChange">
    <DialogContent
      v-if="current"
      :show-close-button="false"
      class="max-w-sm gap-0 rounded-2xl border-batu p-0 text-center shadow-lg"
    >
      <div class="flex flex-col items-center px-7 pt-8 pb-6">
        <ArchBadge :tone="appearance.tone" :size="44">
          <component :is="appearance.icon" class="size-5" :stroke-width="2.25" />
        </ArchBadge>
        <DialogTitle class="mt-5 text-lg font-semibold text-jati">{{ current.title }}</DialogTitle>
        <DialogDescription v-if="current.message" class="mt-2 text-sm leading-relaxed text-jati-muda">
          {{ current.message }}
        </DialogDescription>
      </div>

      <div class="flex gap-2 border-t border-batu bg-kapur/60 px-5 py-4" :class="current.kind === 'confirm' ? 'justify-stretch' : 'justify-center'">
        <Button
          v-if="current.kind === 'confirm'"
          variant="outline"
          class="flex-1"
          @click="settle(false)"
        >
          {{ current.cancelLabel }}
        </Button>
        <Button
          :variant="current.danger || current.kind === 'error' ? 'destructive' : 'default'"
          :class="current.kind === 'confirm' ? 'flex-1' : 'min-w-28'"
          autofocus
          @click="settle(true)"
        >
          {{ current.confirmLabel }}
        </Button>
      </div>
    </DialogContent>
  </Dialog>
</template>
