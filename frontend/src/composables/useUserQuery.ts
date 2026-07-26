import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { userService } from '../api/services/user';
import { toast } from 'vue-sonner';

export const useUsersQuery = () => {
  return useQuery({
    queryKey: ['users'],
    queryFn: () => userService.getAllUsers(1, 100),
  });
};

export const useRolesQuery = () => {
  return useQuery({
    queryKey: ['roles'],
    queryFn: userService.getRoles,
  });
};

export const useInviteUserMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: userService.inviteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      toast.success('User invited successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to invite user');
    },
  });
};

export const useUpdateUserStatusMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) => 
      userService.updateUserStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      toast.success('User status updated');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to update user status');
    },
  });
};

export const useDeleteUserMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => userService.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      toast.success('User permanently deleted');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to delete user');
    },
  });
};
