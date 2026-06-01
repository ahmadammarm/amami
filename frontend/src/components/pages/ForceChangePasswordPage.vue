<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { Lock, Loader2, ShieldAlert } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import axios from 'axios';

const router = useRouter();
const authStore = useAuthStore();
const isLoading = ref(false);

const schema = toTypedSchema(
  zod.object({
    password: zod.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: zod.string().min(1, 'Please confirm your password'),
  }).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ['confirmPassword'],
  })
);

const { handleSubmit, errors, defineField } = useForm({
  validationSchema: schema,
});

const [password, passwordAttrs] = defineField('password');
const [confirmPassword, confirmPasswordAttrs] = defineField('confirmPassword');

const onSubmit = handleSubmit(async (values) => {
  isLoading.value = true;
  try {
    await authStore.changePassword(values.password);
    toast.success('Password changed successfully. Welcome to amami!');
    router.push({ name: 'Dashboard' });
  } catch (error) {
    let message = 'Failed to change password';
    if (axios.isAxiosError(error)) {
      message = error.response?.data?.message || message;
    }
    toast.error(message);
  } finally {
    isLoading.value = false;
  }
});

const handleLogout = () => {
  authStore.logout();
  router.push({ name: 'Login' });
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 font-sans p-4">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-8 border-t-4 border-primary">
      <div class="text-center mb-10">
        <div class="bg-primary/10 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
          <ShieldAlert class="w-8 h-8 text-primary" />
        </div>
        <h1 class="text-2xl font-bold text-gray-900 mb-2">Secure Your Account</h1>
        <p class="text-gray-500 text-sm font-medium px-4">
          This is your first login with a temporary password. Please set a new secure password to continue.
        </p>
      </div>

      <form @submit="onSubmit" class="space-y-6">
        <div>
          <label for="password" class="block text-sm font-semibold text-gray-700 mb-2">New Password</label>
          <div class="relative">
            <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="password"
              v-bind="passwordAttrs"
              type="password"
              id="password"
              placeholder="••••••••"
              class="w-full pl-10 pr-4 py-3 bg-gray-50 border border-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm"
              :class="{ 'border-red-500 focus:ring-red-500/10 focus:border-red-500': errors.password }"
            />
          </div>
          <p v-if="errors.password" class="mt-1.5 text-xs font-medium text-red-500">{{ errors.password }}</p>
        </div>

        <div>
          <label for="confirmPassword" class="block text-sm font-semibold text-gray-700 mb-2">Confirm New Password</label>
          <div class="relative">
            <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="confirmPassword"
              v-bind="confirmPasswordAttrs"
              type="password"
              id="confirmPassword"
              placeholder="••••••••"
              class="w-full pl-10 pr-4 py-3 bg-gray-50 border border-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm"
              :class="{ 'border-red-500 focus:ring-red-500/10 focus:border-red-500': errors.confirmPassword }"
            />
          </div>
          <p v-if="errors.confirmPassword" class="mt-1.5 text-xs font-medium text-red-500">{{ errors.confirmPassword }}</p>
        </div>

        <div class="flex flex-col gap-3 pt-2">
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full py-3 bg-primary text-white font-bold rounded-lg hover:bg-primary/90 focus:outline-none focus:ring-4 focus:ring-primary/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 shadow-sm"
          >
            <Loader2 v-if="isLoading" class="w-5 h-5 animate-spin" />
            {{ isLoading ? 'Updating Password...' : 'Update Password' }}
          </button>
          
          <button
            @click="handleLogout"
            type="button"
            class="w-full py-2.5 text-gray-400 hover:text-gray-600 font-semibold text-sm transition-colors"
          >
            Cancel and Logout
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
