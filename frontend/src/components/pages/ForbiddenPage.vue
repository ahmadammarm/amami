<script setup lang="ts">
import { useRouter } from 'vue-router';
import { ShieldAlert, ArrowLeft } from 'lucide-vue-next';
import Button from '../atoms/button/Button.vue';

const router = useRouter();

const goBack = () => {
  router.back();
};

const goHome = () => {
  router.push({ name: 'Dashboard' });
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 p-4">
    <div class="max-w-md w-full text-center space-y-6">
      <div class="bg-red-50 w-24 h-24 rounded-full flex items-center justify-center mx-auto mb-8 relative">
        <div class="absolute inset-0 bg-red-100 rounded-full animate-ping opacity-50"></div>
        <ShieldAlert class="w-12 h-12 text-red-500 relative z-10" />
      </div>
      
      <div class="space-y-2">
        <h1 class="text-4xl font-black text-gray-900 tracking-tight">403 Forbidden</h1>
        <p class="text-gray-500 font-medium text-lg">
          You don't have permission to access this page.
        </p>
      </div>

      <div class="bg-white p-6 rounded-2xl border border-gray-100 shadow-sm text-sm text-gray-600 text-left">
        <p class="font-medium text-gray-900 mb-2">Why am I seeing this?</p>
        <ul class="list-disc list-inside space-y-1 ml-4 text-gray-500">
          <li>Your current role does not grant access to this module.</li>
          <li>The resource might be restricted to administrators only.</li>
          <li>Your session might need to be refreshed.</li>
        </ul>
      </div>

      <div class="flex items-center gap-3 justify-center pt-4">
        <Button variant="outline" @click="goBack" class="rounded-xl h-11 px-6 bg-white border-gray-200 hover:bg-gray-50">
          <ArrowLeft class="w-4 h-4 mr-2" /> Go Back
        </Button>
        <Button @click="goHome" class="rounded-xl h-11 px-8 font-bold shadow-lg shadow-primary/20">
          Return to Dashboard
        </Button>
      </div>
    </div>
  </div>
</template>
