import type { JamaahListParams } from '../types/jamaah';

/**
 * Central TanStack Query key factory. Each module's `all` key is a prefix of its
 * other keys, so invalidating `all` refreshes every query of that module.
 */
export const queryKeys = {
  dashboard: {
    metrics: ['dashboard', 'metrics'] as const,
  },
  finance: {
    all: ['finance'] as const,
    funds: ['finance', 'funds'] as const,
    transactions: (page: number, limit: number) => ['finance', 'transactions', page, limit] as const,
  },
  zakat: {
    all: ['zakat'] as const,
    donations: (page: number, limit: number) => ['zakat', 'donations', page, limit] as const,
    distributions: (page: number, limit: number) => ['zakat', 'distributions', page, limit] as const,
  },
  jamaah: {
    all: ['jamaah'] as const,
    list: (params: JamaahListParams) => ['jamaah', 'list', params] as const,
    survey: (id: string) => ['jamaah', 'survey', id] as const,
  },
  qurban: {
    all: ['qurban'] as const,
    packages: ['qurban', 'packages'] as const,
    bookings: ['qurban', 'bookings'] as const,
    animals: ['qurban', 'animals'] as const,
  },
  inventory: {
    all: ['inventory'] as const,
    assets: ['inventory', 'assets'] as const,
  },
  agenda: {
    all: ['agenda'] as const,
    list: ['agenda', 'list'] as const,
    documentation: (agendaId: string) => ['agenda', 'documentation', agendaId] as const,
  },
  users: {
    all: ['users'] as const,
    roles: ['users', 'roles'] as const,
  },
  settings: {
    profile: ['settings', 'profile'] as const,
    smtp: ['settings', 'smtp'] as const,
    health: ['settings', 'health'] as const,
    auditLogs: (page: number, limit: number) => ['settings', 'audit-logs', page, limit] as const,
  },
};
