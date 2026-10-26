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
import { queryKeys } from '../api/queryKeys';
import { feedback } from './useFeedback';
import { computed, type MaybeRef, unref } from 'vue';

export const useMosqueProfileQuery = () => {
  return useQuery({
    queryKey: queryKeys.settings.profile,
    queryFn: getMosqueProfile,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
};

/** The page shows its own success dialog; only failures are reported here. */
export const useUpdateMosqueProfile = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateMosqueProfile,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.settings.profile });
    },
    onError: (error: unknown) => {
      feedback.error('Profil masjid belum disimpan', error);
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
      feedback.success('Pengaturan email tersimpan');
    },
    onError: (error: unknown) => {
      feedback.error('Pengaturan email belum disimpan', error);
    },
  });
};

export const useTestSMTPMutation = () => {
  return useMutation({
    mutationFn: testSMTPConnection,
    onSuccess: () => {
      feedback.success('Koneksi email berhasil', 'Server SMTP dapat dihubungi dengan pengaturan ini.');
    },
    onError: (error: unknown) => {
      feedback.error('Koneksi email gagal', error);
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
