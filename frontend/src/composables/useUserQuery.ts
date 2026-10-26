import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { userService } from '../api/services/user';
import { getErrorMessage } from '../api/client';
import { queryKeys } from '../api/queryKeys';
import type { ManageableUserStatus } from '../types/user';
import { toast } from 'vue-sonner';

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

export const useInviteUserMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: userService.inviteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
      toast.success('User invited successfully');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error, 'Failed to invite user'));
    },
  });
};

export const useUpdateUserStatusMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: ManageableUserStatus }) =>
      userService.updateUserStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
      toast.success('User status updated');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error, 'Failed to update user status'));
    },
  });
};

export const useDeleteUserMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => userService.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.users.all });
      toast.success('User deleted');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error, 'Failed to delete user'));
    },
  });
};
