export type UserRole = 'SUPER_ADMIN' | 'BENDAHARA' | 'TAKMIR' | 'SEKRETARIS' | 'JAMAAH';

export interface User {
  id: string;
  username: string;
  email: string;
  role_name: UserRole;
  status: 'ACTIVE' | 'SUSPENDED';
}

export interface LoginResponse {
  token: string;
  user: User;
}
