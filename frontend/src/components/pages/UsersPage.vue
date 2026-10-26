<script setup lang="ts">
import { computed, ref } from 'vue';
import { UserPlus, Search, MoreHorizontal, Ban, CircleCheck, Trash2, UserCog } from 'lucide-vue-next';
import { useUsersQuery, useUpdateUserStatusMutation, useDeleteUserMutation } from '../../composables/useUserQuery';
import { feedback } from '../../composables/useFeedback';
import { useAuthStore } from '../../stores/auth';
import { roleLabel } from '../../lib/roles';
import type { UserStatus } from '../../types/auth';
import type { InviteUserResponse, UserResponse } from '../../types/user';
import type { BadgeTone } from '../../types/ui';
import PageHeader from '../molecules/PageHeader.vue';
import DataPanel from '../molecules/DataPanel.vue';
import EmptyState from '../molecules/EmptyState.vue';
import StatusBadge from '../molecules/StatusBadge.vue';
import CreateUserDialog from '../molecules/CreateUserDialog.vue';
import { Button } from '../atoms/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '../atoms/dropdown-menu';

const USER_STATUS: Record<UserStatus, { label: string; tone: BadgeTone }> = {
  ACTIVE: { label: 'Aktif', tone: 'sajadah' },
  SUSPENDED: { label: 'Dinonaktifkan', tone: 'bata' },
  PENDING_PASSWORD_CHANGE: { label: 'Menunggu ganti password', tone: 'kuningan' },
};

const authStore = useAuthStore();
const { data: users, isLoading } = useUsersQuery();
const { mutate: updateStatus } = useUpdateUserStatusMutation();
const { mutate: deleteUser } = useDeleteUserMutation();

const showInvite = ref(false);
const searchQuery = ref('');

const filteredUsers = computed(() => {
  const items = users.value?.items ?? [];
  const query = searchQuery.value.trim().toLowerCase();
  if (!query) return items;
  return items.filter(u =>
    [u.full_name, u.username, u.email].some(field => field.toLowerCase().includes(query)),
  );
});

/** The backend rejects deleting yourself or a SUPER_ADMIN, so the action is hidden for them. */
const canDelete = (user: UserResponse) => user.role_name !== 'SUPER_ADMIN' && user.id !== authStore.user?.id;

const onInvited = (result: InviteUserResponse) => {
  feedback.success(
    'Akun pengurus dibuat',
    `Sampaikan password awal "${result.temporary_password}" kepada ${result.user.full_name}. Password ini wajib diganti saat pertama masuk.`,
  );
};

const toggleStatus = async (user: UserResponse) => {
  const suspending = user.status !== 'SUSPENDED';
  const confirmed = await feedback.confirm({
    title: suspending ? 'Nonaktifkan pengguna?' : 'Aktifkan pengguna?',
    message: suspending
      ? `${user.full_name || user.username} tidak bisa masuk sampai diaktifkan kembali.`
      : `${user.full_name || user.username} bisa masuk kembali ke dashboard.`,
    confirmLabel: suspending ? 'Nonaktifkan' : 'Aktifkan',
    danger: suspending,
  });
  if (confirmed) updateStatus({ id: user.id, status: suspending ? 'SUSPENDED' : 'ACTIVE' });
};

const confirmDelete = async (user: UserResponse) => {
  const confirmed = await feedback.confirm({
    title: 'Hapus pengguna?',
    message: `Akun ${user.username} dan profil jamaahnya tidak bisa dipakai lagi. Riwayat datanya tetap tersimpan untuk laporan.`,
    confirmLabel: 'Hapus pengguna',
    danger: true,
  });
  if (confirmed) deleteUser(user.id);
};
</script>

<template>
  <div class="space-y-8">
    <PageHeader title="Pengguna" description="Akun pengurus yang dapat masuk ke dashboard, beserta perannya.">
      <template #actions>
        <Button @click="showInvite = true"><UserPlus /> Undang pengurus</Button>
      </template>
    </PageHeader>

    <DataPanel>
      <template #actions>
        <div class="relative w-full sm:w-72">
          <Search class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-jati-muda" />
          <input
            v-model="searchQuery"
            type="search"
            placeholder="Cari nama, username, atau email"
            aria-label="Cari pengguna"
            class="field-control pl-9"
          />
        </div>
      </template>

      <div v-if="isLoading" class="p-8 text-center text-sm text-jati-muda">Memuat pengguna…</div>
      <EmptyState
        v-else-if="filteredUsers.length === 0"
        :icon="UserCog"
        :title="searchQuery ? 'Tidak ada pengguna yang cocok' : 'Belum ada pengguna'"
        :description="searchQuery ? 'Coba kata kunci lain.' : 'Undang pengurus agar mereka bisa mengelola data masjid.'"
      />
      <div v-else class="overflow-x-auto">
        <table class="tabel min-w-[760px]">
          <thead>
            <tr>
              <th>Pengguna</th>
              <th>Peran</th>
              <th>Status</th>
              <th class="w-16"><span class="sr-only">Aksi</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in filteredUsers" :key="user.id">
              <td>
                <div class="flex items-center gap-3">
                  <span class="flex size-9 shrink-0 items-center justify-center rounded-full bg-sajadah-muda font-semibold text-sajadah">
                    {{ (user.full_name || user.username).charAt(0).toUpperCase() }}
                  </span>
                  <div class="min-w-0">
                    <p class="truncate font-medium">
                      {{ user.full_name || user.username }}
                      <span v-if="user.id === authStore.user?.id" class="text-xs font-normal text-jati-muda">(Anda)</span>
                    </p>
                    <p class="truncate text-xs text-jati-muda">@{{ user.username }} · {{ user.email }}</p>
                  </div>
                </div>
              </td>
              <td>{{ roleLabel(user.role_name) }}</td>
              <td><StatusBadge :tone="USER_STATUS[user.status].tone">{{ USER_STATUS[user.status].label }}</StatusBadge></td>
              <td class="text-right">
                <DropdownMenu v-if="user.id !== authStore.user?.id">
                  <DropdownMenuTrigger as-child>
                    <button type="button" class="rounded-md p-2 text-jati-muda hover:bg-kapur hover:text-jati" :aria-label="`Aksi untuk ${user.username}`">
                      <MoreHorizontal class="size-4" />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" class="w-52 rounded-xl border-batu p-1">
                    <DropdownMenuItem class="cursor-pointer rounded-lg py-2" @click="toggleStatus(user)">
                      <template v-if="user.status === 'SUSPENDED'">
                        <CircleCheck class="mr-2 size-4 text-sajadah" /> Aktifkan
                      </template>
                      <template v-else>
                        <Ban class="mr-2 size-4 text-kuningan" /> Nonaktifkan
                      </template>
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      v-if="canDelete(user)"
                      class="cursor-pointer rounded-lg py-2 text-bata focus:bg-bata-muda focus:text-bata"
                      @click="confirmDelete(user)"
                    >
                      <Trash2 class="mr-2 size-4" /> Hapus
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </DataPanel>

    <CreateUserDialog v-model:open="showInvite" @success="onInvited" />
  </div>
</template>
