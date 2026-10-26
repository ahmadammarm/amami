<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { Loader2 } from 'lucide-vue-next';
import { useAuthStore } from '../../stores/auth';
import { feedback } from '../../composables/useFeedback';
import AuthLayout from '../templates/AuthLayout.vue';
import FormField from '../molecules/FormField.vue';
import { Button } from '../atoms/button';

const router = useRouter();
const authStore = useAuthStore();
const isLoading = ref(false);

const schema = toTypedSchema(
  zod.object({
    password: zod.string().min(8, 'Password minimal 8 karakter'),
    confirmPassword: zod.string().min(1, 'Ulangi password baru'),
  }).refine((data) => data.password === data.confirmPassword, {
    message: 'Kedua password tidak sama',
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
    router.push({ name: 'Dashboard' });
    feedback.success('Password tersimpan', 'Mulai sekarang, masuk dengan password baru Anda.');
  } catch (error) {
    feedback.error('Password belum tersimpan', error);
  } finally {
    isLoading.value = false;
  }
});

const signOut = () => {
  authStore.logout();
  router.push({ name: 'Login' });
};
</script>

<template>
  <AuthLayout
    title="Buat password baru"
    description="Anda masuk dengan password sementara dari admin. Buat password sendiri sebelum membuka dashboard."
  >
    <form novalidate class="space-y-5" @submit="onSubmit">
      <FormField label="Password baru" for="password" :error="errors.password" hint="Minimal 8 karakter.">
        <input
          id="password"
          v-model="password"
          v-bind="passwordAttrs"
          type="password"
          autocomplete="new-password"
          class="field-control"
          :aria-invalid="Boolean(errors.password)"
        />
      </FormField>

      <FormField label="Ulangi password baru" for="confirmPassword" :error="errors.confirmPassword">
        <input
          id="confirmPassword"
          v-model="confirmPassword"
          v-bind="confirmPasswordAttrs"
          type="password"
          autocomplete="new-password"
          class="field-control"
          :aria-invalid="Boolean(errors.confirmPassword)"
        />
      </FormField>

      <Button type="submit" size="lg" class="w-full" :disabled="isLoading">
        <Loader2 v-if="isLoading" class="animate-spin" />
        {{ isLoading ? 'Menyimpan…' : 'Simpan password' }}
      </Button>

      <button type="button" class="w-full text-sm font-medium text-jati-muda hover:text-jati" @click="signOut">
        Keluar tanpa mengganti
      </button>
    </form>
  </AuthLayout>
</template>
