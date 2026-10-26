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
import { getErrorMessage } from '../api/client';
import { queryKeys } from '../api/queryKeys';
import { toast } from 'vue-sonner';
import { computed, type MaybeRef, unref } from 'vue';

export const useMosqueProfileQuery = () => {
  return useQuery({
    queryKey: queryKeys.settings.profile,
    queryFn: getMosqueProfile,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
};

export const useUpdateMosqueProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateMosqueProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.settings.profile });
      toast.success('Profile updated successfully');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error, 'Failed to update profile'));
    },
  });
};

export const useSMTPConfigQuery = () => {
  return useQuery({
    queryKey: queryKeys.settings.smtp,
    queryFn: getSMTPConfig,
  });
};

export const useUpdateSMTPMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateSMTPConfig,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.settings.smtp });
      toast.success('SMTP configuration updated successfully');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error, 'Failed to update SMTP configuration'));
    },
  });
};

export const useTestSMTPMutation = () => {
  return useMutation({
    mutationFn: testSMTPConnection,
    onSuccess: () => {
      toast.success('SMTP connection successful');
    },
    onError: (error: unknown) => {
      toast.error(getErrorMessage(error, 'SMTP connection failed'));
    },
  });
};

export const useSystemHealthQuery = () => {
  return useQuery({
    queryKey: queryKeys.settings.health,
    queryFn: getSystemHealth,
    refetchInterval: 30000, // Every 30 seconds
  });
};

export const useAuditLogsQuery = (page: MaybeRef<number> = 1, limit: MaybeRef<number> = 10) => {
  return useQuery({
    queryKey: computed(() => queryKeys.settings.auditLogs(unref(page), unref(limit))),
    queryFn: () => getAuditLogs(unref(page), unref(limit)),
  });
};
