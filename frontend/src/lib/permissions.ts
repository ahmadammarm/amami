/**
 * Permission codes granted to roles by the backend (backend/internal/domain/permissions.go).
 * The UI decides visibility from the codes returned by /auth/me, never from role names.
 */
export const Permission = {
  UserManage: 'user:manage',
  LedgerRead: 'ledger:read',
  LedgerWrite: 'ledger:write',
  ZakatRead: 'zakat:read',
  ZakatManage: 'zakat:manage',
  JamaahManage: 'jamaah:manage',
  MustahikSurvey: 'mustahik:survey',
  QurbanRead: 'qurban:read',
  QurbanManage: 'qurban:manage',
  InventoryRead: 'inventory:read',
  InventoryManage: 'inventory:manage',
  AgendaRead: 'agenda:read',
  AgendaManage: 'agenda:manage',
  SettingsManage: 'settings:manage',
} as const;

export type Permission = (typeof Permission)[keyof typeof Permission];
