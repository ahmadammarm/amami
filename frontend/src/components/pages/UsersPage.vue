<script setup lang="ts">
import { ref } from 'vue';
import { useUsersQuery, useUpdateUserStatusMutation, useDeleteUserMutation } from '../../composables/useUserQuery';
import { 
  Users as UsersIcon, 
  UserPlus, 
  Search, 
  MoreVertical, 
  Shield, 
  Mail, 
  Circle,
  RefreshCcw,
  Ban,
  CheckCircle2,
  AlertCircle,
  Trash2
} from 'lucide-vue-next';
import Button from '../atoms/button/Button.vue';
import Input from '../atoms/input/Input.vue';
import CreateUserDialog from '../molecules/CreateUserDialog.vue';
import SuccessDialog from '../atoms/SuccessDialog.vue';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '../atoms/alert-dialog';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../atoms/dropdown-menu';

const { data: users, isLoading, isRefetching, refetch } = useUsersQuery();
const { mutate: updateStatus } = useUpdateUserStatusMutation();
const { mutate: deleteUser, isPending: isDeleting } = useDeleteUserMutation();

const isInviteDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const userToDelete = ref<any>(null);
const searchQuery = ref('');
const successDialog = ref({
  open: false,
  title: '',
  description: '',
});

const handleInviteSuccess = (tempPass: string) => {
  successDialog.value = {
    open: true,
    title: 'User Invited Successfully',
    description: `The account has been created. Please share this initial password with the user: "${tempPass}". They will be forced to change it on their first login.`,
  };
};

const toggleUserStatus = (user: any) => {
  const newStatus = user.status === 'ACTIVE' ? 'SUSPENDED' : 'ACTIVE';
  updateStatus({ id: user.id, status: newStatus });
};

const confirmDelete = (user: any) => {
  userToDelete.value = user;
  isDeleteDialogOpen.value = true;
};

const executeDelete = () => {
  if (userToDelete.value) {
    deleteUser(userToDelete.value.id, {
      onSuccess: () => {
        isDeleteDialogOpen.value = false;
        userToDelete.value = null;
      }
    });
  }
};

const getStatusColor = (status: string) => {
  switch (status) {
    case 'ACTIVE': return 'text-emerald-600 bg-emerald-50';
    case 'SUSPENDED': return 'text-red-600 bg-red-50';
    case 'PENDING_PASSWORD_CHANGE': return 'text-amber-600 bg-amber-50';
    default: return 'text-gray-600 bg-gray-50';
  }
};
</script>

<template>
  <div class="max-w-6xl mx-auto space-y-8 pb-12">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div class="space-y-1">
        <div class="flex items-center gap-2 text-primary font-semibold text-sm tracking-wider uppercase">
          <Shield class="w-4 h-4" />
          Administrative
        </div>
        <h1 class="text-3xl font-bold text-gray-900">User Management</h1>
        <p class="text-gray-500 font-medium">Manage your mosque staff, roles, and system access.</p>
      </div>
      
      <div class="flex items-center gap-3">
        <Button 
          variant="outline" 
          @click="refetch" 
          :disabled="isLoading || isRefetching"
          class="bg-white border-gray-200 h-11 rounded-xl"
        >
          <RefreshCcw :class="['w-4 h-4', (isLoading || isRefetching) && 'animate-spin']" />
        </Button>
        <Button 
          @click="isInviteDialogOpen = true"
          class="h-11 px-6 rounded-xl font-bold shadow-lg shadow-primary/20"
        >
          <UserPlus class="w-4 h-4 mr-2" />
          Invite New User
        </Button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden min-h-[400px]">
      <!-- Filters/Search -->
      <div class="p-6 border-b border-gray-50 bg-gray-50/30 flex items-center gap-4">
        <div class="relative flex-1 max-w-sm">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
          <Input 
            v-model="searchQuery" 
            placeholder="Search by name, email, or username..." 
            class="pl-10 h-10 rounded-lg border-gray-200 bg-white"
          />
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[800px]">
          <thead>
            <tr class="bg-gray-50/50 text-[10px] uppercase tracking-widest text-gray-400 font-bold border-b border-gray-50">
              <th class="px-8 py-4">User Details</th>
              <th class="px-8 py-4">Role</th>
              <th class="px-8 py-4">Status</th>
              <th class="px-8 py-4 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-50">
            <template v-if="isLoading">
              <tr v-for="i in 5" :key="i">
                <td colspan="4" class="px-8 py-6">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-full bg-gray-100 animate-pulse"></div>
                    <div class="space-y-2 flex-1">
                      <div class="h-4 bg-gray-100 animate-pulse rounded w-1/4"></div>
                      <div class="h-3 bg-gray-100 animate-pulse rounded w-1/3"></div>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
            <template v-else>
              <tr v-for="user in users?.items" :key="user.id" class="hover:bg-gray-50/30 transition-colors group">
                <td class="px-8 py-5">
                  <div class="flex items-center gap-4">
                    <div class="w-11 h-11 rounded-full bg-primary/10 flex items-center justify-center text-primary font-bold text-lg">
                      {{ user.full_name?.[0] || user.username[0].toUpperCase() }}
                    </div>
                    <div class="flex flex-col">
                      <span class="text-sm font-bold text-gray-900 leading-none mb-1">{{ user.full_name || 'No Name' }}</span>
                      <div class="flex items-center gap-2 text-xs text-gray-400 font-medium">
                        <span class="text-gray-600">@{{ user.username }}</span>
                        <Circle class="w-1 h-1 fill-gray-300 stroke-none" />
                        <span class="flex items-center gap-1"><Mail class="w-3 h-3" /> {{ user.email }}</span>
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-5">
                  <div class="flex items-center gap-2 text-xs font-bold text-gray-700 bg-gray-100 w-fit px-3 py-1 rounded-lg">
                    <Shield class="w-3.5 h-3.5 text-gray-400" />
                    {{ user.role_name }}
                  </div>
                </td>
                <td class="px-8 py-5">
                  <div :class="[getStatusColor(user.status), 'px-3 py-1 rounded-full text-[10px] font-black tracking-tight inline-flex items-center gap-1.5 uppercase']">
                    <CheckCircle2 v-if="user.status === 'ACTIVE'" class="w-3 h-3" />
                    <AlertCircle v-else-if="user.status === 'PENDING_PASSWORD_CHANGE'" class="w-3 h-3" />
                    <Ban v-else class="w-3 h-3" />
                    {{ user.status.replace(/_/g, ' ') }}
                  </div>
                </td>
                <td class="px-8 py-5 text-right">
                  <DropdownMenu>
                    <DropdownMenuTrigger as-child>
                      <Button variant="ghost" size="icon" class="h-9 w-9 rounded-lg hover:bg-white hover:shadow-sm border border-transparent hover:border-gray-100">
                        <MoreVertical class="w-4 h-4 text-gray-400" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" class="w-48 rounded-xl p-1 shadow-xl border-gray-100">
                      <DropdownMenuItem @click="toggleUserStatus(user)" class="rounded-lg font-bold text-xs py-2.5 cursor-pointer">
                        <template v-if="user.status === 'ACTIVE'">
                          <Ban class="w-4 h-4 mr-2 text-amber-500" /> Suspend User
                        </template>
                        <template v-else>
                          <CheckCircle2 class="w-4 h-4 mr-2 text-emerald-500" /> Activate User
                        </template>
                      </DropdownMenuItem>
                      <DropdownMenuItem v-if="user.role_name !== 'SUPER_ADMIN'" @click="confirmDelete(user)" class="rounded-lg font-bold text-xs py-2.5 cursor-pointer text-red-600 focus:bg-red-50 focus:text-red-700">
                        <Trash2 class="w-4 h-4 mr-2" /> Hapus User
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </td>
              </tr>
              <tr v-if="!users?.items?.length">
                <td colspan="4" class="py-24 text-center">
                  <div class="max-w-xs mx-auto space-y-4">
                    <div class="p-4 bg-gray-50 rounded-full w-16 h-16 flex items-center justify-center mx-auto">
                      <UsersIcon class="w-8 h-8 text-gray-300" />
                    </div>
                    <div>
                      <h3 class="font-bold text-gray-900">No users found</h3>
                      <p class="text-sm text-gray-400 font-medium mt-1">Start by inviting your first mosque staff member.</p>
                    </div>
                    <Button variant="outline" size="sm" @click="isInviteDialogOpen = true" class="rounded-lg h-10 px-6">
                      Invite Someone
                    </Button>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modals -->
    <CreateUserDialog 
      v-model:open="isInviteDialogOpen" 
      @success="handleInviteSuccess"
    />

    <SuccessDialog
      v-model:open="successDialog.open"
      :title="successDialog.title"
      :description="successDialog.description"
    />

    <AlertDialog v-model:open="isDeleteDialogOpen">
      <AlertDialogContent class="rounded-2xl">
        <AlertDialogHeader>
          <AlertDialogTitle class="text-xl">Hapus Pengguna Secara Permanen?</AlertDialogTitle>
          <AlertDialogDescription class="text-base leading-relaxed">
            Anda yakin ingin menghapus <strong>{{ userToDelete?.username }}</strong>? Seluruh data yang berkaitan dengan pengguna ini (seperti identitas jamaah) juga akan ikut terhapus. Aksi ini tidak dapat dibatalkan.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter class="gap-3">
          <AlertDialogCancel class="rounded-xl h-11 font-bold" :disabled="isDeleting">Batal</AlertDialogCancel>
          <AlertDialogAction @click="executeDelete" :disabled="isDeleting" class="bg-red-500 hover:bg-red-600 focus:ring-red-500 rounded-xl h-11 font-bold shadow-lg shadow-red-200">
            {{ isDeleting ? 'Menghapus...' : 'Ya, Hapus' }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
