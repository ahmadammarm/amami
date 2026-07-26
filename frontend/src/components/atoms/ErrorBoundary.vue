<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue';
import { AlertCircle, RefreshCcw } from 'lucide-vue-next';
import Button from './button/Button.vue';

const hasError = ref(false);
const errorDetails = ref<Error | null>(null);

onErrorCaptured((err: unknown) => {
  hasError.value = true;
  if (err instanceof Error) {
    errorDetails.value = err;
  }
  // Prevent error from propagating further if we want to handle it entirely here
  return false;
});

const reloadPage = () => {
  window.location.reload();
};
</script>

<template>
  <div v-if="hasError" class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-2xl p-8 text-center shadow-xl border border-gray-100">
      <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-6">
        <AlertCircle class="w-8 h-8 text-red-600" />
      </div>
      <h1 class="text-2xl font-bold text-gray-900 mb-2">Something went wrong</h1>
      <p class="text-gray-500 mb-6">
        We encountered an unexpected error while rendering this page. Our team has been notified.
      </p>
      
      <div v-if="errorDetails" class="bg-gray-50 rounded-xl p-4 text-left mb-6 overflow-x-auto text-xs text-gray-600 font-mono border border-gray-100">
        {{ errorDetails.message }}
      </div>

      <Button @click="reloadPage" class="w-full h-11 rounded-xl text-base font-bold shadow-lg shadow-primary/20">
        <RefreshCcw class="w-4 h-4 mr-2" />
        Reload Page
      </Button>
    </div>
  </div>
  <slot v-else />
</template>
