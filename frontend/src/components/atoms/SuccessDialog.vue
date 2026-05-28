<script setup lang="ts">
import { 
  Dialog, 
  DialogContent, 
  DialogTitle, 
  DialogDescription,
  DialogFooter 
} from '../atoms/dialog';
import { Button } from '../atoms/button';
import { CheckCircle2 } from 'lucide-vue-next';

interface Props {
  open: boolean;
  title: string;
  description?: string;
  buttonText?: string;
}

const props = withDefaults(defineProps<Props>(), {
  description: '',
  buttonText: 'Continue',
});

const emit = defineEmits(['update:open', 'confirm']);

const handleClose = () => {
  emit('update:open', false);
  emit('confirm');
};
</script>

<template>
  <Dialog :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogContent class="sm:max-w-md rounded-2xl p-0 overflow-hidden border-none shadow-2xl">
      <div class="bg-emerald-500 h-24 flex items-center justify-center relative">
        <div class="absolute -bottom-10 bg-white p-2 rounded-full shadow-lg">
          <CheckCircle2 class="w-16 h-16 text-emerald-500" />
        </div>
      </div>
      
      <div class="pt-14 pb-8 px-8 text-center space-y-2">
        <DialogTitle class="text-2xl font-bold text-gray-900">{{ title }}</DialogTitle>
        <DialogDescription v-if="description" class="text-gray-500 font-medium leading-relaxed">
          {{ description }}
        </DialogDescription>
      </div>

      <DialogFooter class="p-6 pt-0 sm:justify-center">
        <Button 
          @click="handleClose" 
          class="w-full sm:w-48 h-12 rounded-xl bg-emerald-600 hover:bg-emerald-700 font-bold text-base shadow-lg shadow-emerald-200 transition-all active:scale-95"
        >
          {{ buttonText }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
