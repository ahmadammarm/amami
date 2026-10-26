import { request } from '../client';
import type { DashboardMetricsResponse } from '../../types/dashboard';

export const dashboardService = {
  getMetrics(): Promise<DashboardMetricsResponse> {
    return request<DashboardMetricsResponse>({ method: 'GET', url: '/dashboard/metrics' });
  },
};
