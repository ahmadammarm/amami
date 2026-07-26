import apiClient from '../client';
import type { DashboardMetricsResponse } from '../../types/dashboard';

export const dashboardService = {
  async getMetrics(): Promise<DashboardMetricsResponse> {
    const response = await apiClient.get<DashboardMetricsResponse>('/dashboard/metrics');
    return response.data;
  },
};
