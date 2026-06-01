<script setup lang="ts">
import { useForm } from 'vee-validate';
import { toTypedSchema } from '@vee-validate/zod';
import * as zod from 'zod';
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle, 
  DialogDescription,
  DialogFooter
} from '../atoms/dialog';
import Button from '../atoms/button/Button.vue';
import Input from '../atoms/input/Input.vue';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '../atoms/select';
import { useRolesQuery, useInviteUserMutation } from '../../composables/useUserQuery';
import { Loader2, UserPlus, Shield, Mail, User, Lock } from 'lucide-vue-next';

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void;
  (e: 'success', tempPass: string): void;
}>();

const { data: roles } = useRolesQuery();
const { mutate: inviteUser, isPending: isInviting } = useInviteUserMutation();

const schema = toTypedSchema(
  zod.object({
    username: zod.string().min(3, 'Username must be at least 3 characters'),
    email: zod.string().email('Invalid email address'),
    full_name: zod.string().min(1, 'Full name is required'),
    role_id: zod.number({ required_error: 'Please select a role' }),
    password: zod.string().min(8, 'Initial password must be at least 8 characters'),
  })
);

const { handleSubmit, errors, defineField, resetForm } = useForm({
  validationSchema: schema,
});

const [username, usernameAttrs] = defineField('username');
const [email, emailAttrs] = defineField('email');
const [fullName, fullNameAttrs] = defineField('full_name');
const [roleId, roleIdAttrs] = defineField('role_id');
const [password, passwordAttrs] = defineField('password');

const onSubmit = handleSubmit((values) => {
  inviteUser(values, {
    onSuccess: () => {
      resetForm();
      emit('update:open', false);
      // We pass back the temp pass if we want to show it in a success dialog
      emit('success', values.password);
    }
  });
});
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="max-w-md rounded-2xl p-0 overflow-hidden border-none shadow-2xl">
      <DialogHeader class="px-8 pt-8 pb-6 bg-primary/5">
        <div class="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center mb-4">
          <UserPlus class="w-6 h-6 text-primary" />
        </div>
        <DialogTitle class="text-2xl font-bold text-gray-900">Invite New Staff</DialogTitle>
        <DialogDescription class="text-gray-500 font-medium pt-1">
          Create a new account and assign a role for mosque management.
        </DialogDescription>
      </DialogHeader>

      <form @submit="onSubmit" class="px-8 py-6 space-y-5 bg-white">
        <div class="space-y-1.5">
          <label class="text-sm font-bold text-gray-700 ml-1 flex items-center gap-2">
            <User class="w-3.5 h-3.5" /> Full Name
          </label>
          <Input v-model="fullName" v-bind="fullNameAttrs" placeholder="e.g. Ahmad Ammar" class="rounded-xl h-11 border-gray-100 bg-gray-50/50" />
          <p v-if="errors.full_name" class="text-[10px] font-bold text-red-500 ml-1">{{ errors.full_name }}</p>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-sm font-bold text-gray-700 ml-1 flex items-center gap-2">
              <Shield class="w-3.5 h-3.5" /> Username
            </label>
            <Input v-model="username" v-bind="usernameAttrs" placeholder="ahmad123" class="rounded-xl h-11 border-gray-100 bg-gray-50/50" />
            <p v-if="errors.username" class="text-[10px] font-bold text-red-500 ml-1">{{ errors.username }}</p>
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-bold text-gray-700 ml-1 flex items-center gap-2">
              <Shield class="w-3.5 h-3.5" /> Assigned Role
            </label>
            <Select v-model="roleId" v-bind="roleIdAttrs">
              <SelectTrigger class="rounded-xl h-11 border-gray-100 bg-gray-50/50 font-medium">
                <SelectValue placeholder="Choose Role" />
              </SelectTrigger>
              <SelectContent class="rounded-xl">
                <SelectItem v-for="role in roles" :key="role.id" :value="role.id" class="font-medium">
                  {{ role.name }}
                </SelectItem>
              </SelectContent>
            </Select>
            <p v-if="errors.role_id" class="text-[10px] font-bold text-red-500 ml-1">{{ errors.role_id }}</p>
          </div>
        </div>

        <div class="space-y-1.5">
          <label class="text-sm font-bold text-gray-700 ml-1 flex items-center gap-2">
            <Mail class="w-3.5 h-3.5" /> Email Address
          </label>
          <Input v-model="email" v-bind="emailAttrs" placeholder="ahmad@amami.org" class="rounded-xl h-11 border-gray-100 bg-gray-50/50" />
          <p v-if="errors.email" class="text-[10px] font-bold text-red-500 ml-1">{{ errors.email }}</p>
        </div>

        <div class="space-y-1.5 pb-2">
          <label class="text-sm font-bold text-gray-700 ml-1 flex items-center gap-2">
            <Lock class="w-3.5 h-3.5" /> Initial Password
          </label>
          <Input v-model="password" v-bind="passwordAttrs" type="password" placeholder="••••••••" class="rounded-xl h-11 border-gray-100 bg-gray-50/50" />
          <p v-if="errors.password" class="text-[10px] font-bold text-red-500 ml-1">{{ errors.password }}</p>
          <p class="text-[10px] text-gray-400 ml-1 font-medium italic">User will be forced to change this on first login.</p>
        </div>

        <DialogFooter class="flex flex-row gap-3 pt-4 border-t border-gray-50 sm:justify-end">
          <Button type="button" variant="ghost" @click="emit('update:open', false)" class="flex-1 sm:flex-none rounded-xl h-11 font-bold">Cancel</Button>
          <Button type="submit" :disabled="isInviting" class="flex-1 sm:flex-none rounded-xl h-11 px-8 font-bold shadow-lg shadow-primary/20">
            <Loader2 v-if="isInviting" class="w-4 h-4 mr-2 animate-spin" />
            <UserPlus v-else class="w-4 h-4 mr-2" /> Invite User
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
