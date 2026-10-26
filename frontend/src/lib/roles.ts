import type { UserRole } from '../types/auth';

/** Display names for roles; access is always decided by permission codes, never by these. */
export const ROLE_LABELS: Record<UserRole, string> = {
  SUPER_ADMIN: 'Super admin',
  BENDAHARA: 'Bendahara',
  TAKMIR: 'Takmir',
  SEKRETARIS: 'Sekretaris',
  JAMAAH: 'Jamaah',
};

export const roleLabel = (role: UserRole): string => ROLE_LABELS[role] ?? role;
