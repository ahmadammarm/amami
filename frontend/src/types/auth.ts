export type UserRole = 'SUPER_ADMIN' | 'BENDAHARA' | 'TAKMIR' | 'SEKRETARIS' | 'JAMAAH';

export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role_name: UserRole;
  status: 'ACTIVE' | 'SUSPENDED' | 'PENDING_PASSWORD_CHANGE';
}

export interface LoginResponse {
  token: string;
  requires_password_change: boolean;
}
