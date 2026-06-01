<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/auth';
import apiClient from '../../api/client';
import axios from 'axios';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { Lock, Mail, Loader2 } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

const router = useRouter();
const authStore = useAuthStore();
const isLoading = ref(false);

const schema = toTypedSchema(
  zod.object({
    email: zod.string().min(1, 'Email is required').email('Invalid email format'),
    password: zod.string().min(1, 'Password is required'),
  })
);

const { handleSubmit, errors, defineField } = useForm({
  validationSchema: schema,
});

const [email, emailAttrs] = defineField('email');
const [password, passwordAttrs] = defineField('password');

const onSubmit = handleSubmit(async (values) => {
  isLoading.value = true;
  try {
    const response = await apiClient.post('/auth/login', values);
    const { token, requires_password_change } = response.data.data;
    
    authStore.setToken(token, requires_password_change);
    
    if (requires_password_change) {
      router.push({ name: 'ForceChangePassword' });
    } else {
      await authStore.fetchUser();
      // Redirect immediately to dashboard with success query param
      router.push({ path: '/', query: { loginSuccess: 'true' } });
    }
  } catch (error) {
    let message = 'Invalid email or password';
    if (axios.isAxiosError(error)) {
      message = error.response?.data?.message || message;
    }
    toast.error(message);
  } finally {
    isLoading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 font-sans p-4">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-8">
      <div class="text-center mb-10">
        <h1 class="text-3xl font-bold text-primary mb-2">amami</h1>
        <p class="text-gray-400 font-medium">Please sign in to your account</p>
      </div>

      <form @submit="onSubmit" class="space-y-6">
        <div>
          <label for="email" class="block text-sm font-semibold text-gray-700 mb-2">Email Address</label>
          <div class="relative">
            <Mail class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="email"
              v-bind="emailAttrs"
              type="email"
              id="email"
              placeholder="admin@amami.org"
              class="w-full pl-10 pr-4 py-3 bg-gray-50 border border-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all text-sm"
              :class="{ 'border-red-500 focus:ring-red-500/10 focus:border-red-500': errors.email }"
            />
          </div>
          <p v-if="errors.email" class="mt-1.5 text-xs font-medium text-red-500">{{ errors.email }}</p>
        </div>

        <div>
          <label for="password" class="block text-sm font-semibold text-gray-700 mb-2">Password</label>
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

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full py-3 bg-primary text-white font-bold rounded-lg hover:bg-primary/90 focus:outline-none focus:ring-4 focus:ring-primary/20 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 shadow-sm"
        >
          <Loader2 v-if="isLoading" class="w-5 h-5 animate-spin" />
          {{ isLoading ? 'Signing in...' : 'Sign In' }}
        </button>
      </form>

      <div class="mt-8 text-center text-xs text-gray-400 font-medium">
        &copy; 2026 amami • Mosque Intelligence ERP
      </div>
    </div>
  </div>
</template>
