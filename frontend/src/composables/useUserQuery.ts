import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { userService } from '../api/services/user';
import { queryKeys } from '../api/queryKeys';
import type { ManageableUserStatus } from '../types/user';
import { feedback } from './useFeedback';

export const useUsersQuery = () => {
  return useQuery({
    queryKey: queryKeys.users.all,
    queryFn: () => userService.getAllUsers(1, 100),
  });
};

export const useRolesQuery = () => {
  return useQuery({
    queryKey: queryKeys.users.roles,
    queryFn: userService.getRoles,
    staleTime: Infinity,
  });
};

/** The success dialog is shown by the caller, which knows the temporary password. */
export const useInviteUserMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: userService.inviteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
    },
    onError: (error: unknown) => {
      feedback.error('Pengguna belum diundang', error);
    },
  });
};

export const useUpdateUserStatusMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: ManageableUserStatus }) =>
      userService.updateUserStatus(id, status),
    onSuccess: (_data, { status }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
      feedback.success(status === 'ACTIVE' ? 'Pengguna diaktifkan' : 'Pengguna dinonaktifkan');
    },
    onError: (error: unknown) => {
      feedback.error('Status pengguna belum diubah', error);
    },
  });
};

export const useDeleteUserMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => userService.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
      feedback.success('Pengguna dihapus');
    },
    onError: (error: unknown) => {
      feedback.error('Pengguna belum dihapus', error);
    },
  });
};
