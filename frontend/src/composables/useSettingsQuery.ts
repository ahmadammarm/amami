import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import { 
  getMosqueProfile, 
  updateMosqueProfile, 
  getSMTPConfig, 
  updateSMTPConfig, 
  testSMTPConnection,
  getSystemHealth,
  getAuditLogs
} from '../api/services/settings';
import { toast } from 'vue-sonner';
import { type MaybeRef, unref } from 'vue';

export const useMosqueProfileQuery = () => {
  return useQuery({
    queryKey: ['mosque-profile'],
    queryFn: getMosqueProfile,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
};

export const useUpdateMosqueProfile = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: updateMosqueProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['mosque-profile'] });
      toast.success('Profile updated successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to update profile');
    },
  });
};

export const useSMTPConfigQuery = () => {
  return useQuery({
    queryKey: ['smtp-config'],
    queryFn: getSMTPConfig,
  });
};

export const useUpdateSMTPMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateSMTPConfig,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['smtp-config'] });
      toast.success('SMTP configuration updated successfully');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'Failed to update SMTP configuration');
    },
  });
};

export const useTestSMTPMutation = () => {
  return useMutation({
    mutationFn: testSMTPConnection,
    onSuccess: () => {
      toast.success('SMTP connection successful');
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.message || 'SMTP connection failed');
    },
  });
};

export const useSystemHealthQuery = () => {
  return useQuery({
    queryKey: ['system-health'],
    queryFn: getSystemHealth,
    refetchInterval: 30000, // Every 30 seconds
  });
};

export const useAuditLogsQuery = (page: MaybeRef<number> = 1, limit: MaybeRef<number> = 10) => {
  return useQuery({
    queryKey: ['audit-logs', page, limit],
    queryFn: () => getAuditLogs(unref(page), unref(limit)),
  });
};
