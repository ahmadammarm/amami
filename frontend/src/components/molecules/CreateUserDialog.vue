<script setup lang="ts">
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { useRolesQuery, useInviteUserMutation } from '../../composables/useUserQuery';
import { roleLabel } from '../../lib/roles';
import type { InviteUserResponse } from '../../types/user';
import FormDialog from './FormDialog.vue';
import FormField from './FormField.vue';

defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void;
  (e: 'success', result: InviteUserResponse): void;
}>();

const { data: roles } = useRolesQuery();
const { mutate: inviteUser, isPending: isInviting } = useInviteUserMutation();

const schema = toTypedSchema(
  zod.object({
    full_name: zod.string().trim().min(1, 'Nama lengkap wajib diisi').max(100),
    username: zod.string().trim().min(3, 'Username minimal 3 karakter').max(50),
    email: zod.string().trim().email('Format email tidak valid'),
    role_id: zod.number({ required_error: 'Pilih peran', invalid_type_error: 'Pilih peran' }).int().positive('Pilih peran'),
    password: zod.string().min(8, 'Password awal minimal 8 karakter'),
  })
);

const emptyForm = { full_name: '', username: '', email: '', role_id: 0, password: '' };

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: schema,
  initialValues: emptyForm,
});

const [fullName] = defineField('full_name');
const [username] = defineField('username');
const [email] = defineField('email');
const [roleId] = defineField('role_id');
const [password] = defineField('password');

const onSubmit = handleSubmit((values) => {
  inviteUser(values, {
    onSuccess: (result) => {
      resetForm({ values: emptyForm });
      emit('update:open', false);
      emit('success', result);
    },
  });
});
</script>

<template>
  <FormDialog
    :open="open"
    title="Undang pengurus"
    description="Akun baru langsung aktif. Pengguna wajib mengganti password awal saat pertama masuk."
    submit-label="Buat akun"
    :submitting="isInviting"
    @update:open="emit('update:open', $event)"
    @submit="onSubmit"
  >
    <FormField label="Nama lengkap" for="invite-name" :error="errors.full_name" required>
      <input id="invite-name" v-model="fullName" type="text" class="field-control" :aria-invalid="Boolean(errors.full_name)" />
    </FormField>
    <div class="grid gap-4 sm:grid-cols-2">
      <FormField label="Username" for="invite-username" :error="errors.username" required>
        <input id="invite-username" v-model="username" type="text" autocomplete="off" class="field-control" :aria-invalid="Boolean(errors.username)" />
      </FormField>
      <FormField label="Peran" for="invite-role" :error="errors.role_id" required>
        <select id="invite-role" v-model.number="roleId" class="field-control" :aria-invalid="Boolean(errors.role_id)">
          <option :value="0" disabled>Pilih peran</option>
          <option v-for="role in roles" :key="role.id" :value="role.id">{{ roleLabel(role.name) }}</option>
        </select>
      </FormField>
    </div>
    <FormField label="Email" for="invite-email" :error="errors.email" required>
      <input id="invite-email" v-model="email" type="email" autocomplete="off" class="field-control" :aria-invalid="Boolean(errors.email)" />
    </FormField>
    <FormField label="Password awal" for="invite-password" :error="errors.password" hint="Sampaikan kepada pengguna secara langsung." required>
      <input id="invite-password" v-model="password" type="password" autocomplete="new-password" class="field-control" :aria-invalid="Boolean(errors.password)" />
    </FormField>
  </FormDialog>
</template>
