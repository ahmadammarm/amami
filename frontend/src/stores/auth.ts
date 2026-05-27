import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { User, UserRole } from '../types/auth';
import apiClient from '../api/client';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null);
  const token = ref<string | null>(localStorage.getItem('token'));

  const isAuthenticated = computed(() => !!token.value);
  const userRole = computed(() => user.value?.role_name);

  async function fetchUser(): Promise<void> {
    if (!token.value) return;
    try {
      const response = await apiClient.get('/auth/me');
      // response.data.data because of the common response wrapper on backend
      user.value = response.data.data;
    } catch (error) {
      logout();
      throw error;
    }
  }

  function setToken(newToken: string): void {
    token.value = newToken;
    localStorage.setItem('token', newToken);
  }

  function logout(): void {
    user.value = null;
    token.value = null;
    localStorage.removeItem('token');
  }

  return {
    user,
    token,
    isAuthenticated,
    userRole,
    fetchUser,
    setToken,
    logout,
  };
});
