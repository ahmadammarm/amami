import { useQuery } from '@tanstack/vue-query';
import { getMosqueProfile } from '../api/services/settings';

export const useMosqueProfileQuery = () => {
  return useQuery({
    queryKey: ['mosque-profile'],
    queryFn: getMosqueProfile,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
};
