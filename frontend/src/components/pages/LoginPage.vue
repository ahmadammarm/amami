<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { Loader2 } from 'lucide-vue-next';
import { useAuthStore } from '../../stores/auth';
import { request } from '../../api/client';
import { feedback } from '../../composables/useFeedback';
import type { LoginResponse } from '../../types/auth';
import AuthLayout from '../templates/AuthLayout.vue';
import FormField from '../molecules/FormField.vue';
import { Button } from '../atoms/button';

const router = useRouter();
const authStore = useAuthStore();
const isLoading = ref(false);

const schema = toTypedSchema(
  zod.object({
    email: zod.string().trim().min(1, 'Email wajib diisi').email('Format email tidak valid'),
    password: zod.string().min(1, 'Password wajib diisi'),
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
    const { token, requires_password_change } = await request<LoginResponse>({
      method: 'POST',
      url: '/auth/login',
      data: values,
    });

    authStore.setToken(token, requires_password_change);

    if (requires_password_change) {
      router.push({ name: 'ForceChangePassword' });
    } else {
      await authStore.fetchUser();
      router.push({ name: 'Dashboard' });
    }
  } catch (error) {
    feedback.error('Gagal masuk', error);
  } finally {
    isLoading.value = false;
  }
});
</script>

<template>
  <AuthLayout
    title="Masuk ke amami"
    description="Gunakan email dan password akun pengurus Anda."
  >
    <form novalidate class="space-y-5" @submit="onSubmit">
      <FormField label="Email" for="email" :error="errors.email">
        <input
          id="email"
          v-model="email"
          v-bind="emailAttrs"
          type="email"
          autocomplete="username"
          placeholder="nama@masjid.org"
          class="field-control"
          :aria-invalid="Boolean(errors.email)"
        />
      </FormField>

      <FormField label="Password" for="password" :error="errors.password">
        <input
          id="password"
          v-model="password"
          v-bind="passwordAttrs"
          type="password"
          autocomplete="current-password"
          class="field-control"
          :aria-invalid="Boolean(errors.password)"
        />
      </FormField>

      <Button type="submit" size="lg" class="w-full" :disabled="isLoading">
        <Loader2 v-if="isLoading" class="animate-spin" />
        {{ isLoading ? 'Memeriksa…' : 'Masuk' }}
      </Button>
    </form>

    <p class="mt-8 text-xs leading-relaxed text-jati-muda">
      Belum punya akun? Akun pengurus dibuat oleh admin masjid melalui menu Pengguna.
    </p>
  </AuthLayout>
</template>
