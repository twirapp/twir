import { useMutation, useQuery } from '@tanstack/react-query';
import { AuthUser } from '@tsuwari/shared';
import { deleteCookie } from 'cookies-next';

import { authFetch, authFetcher } from '@/services/api';

export const useProfile = () =>
  useQuery<AuthUser & { apiKey: string }>({
    queryKey: [`/api/auth/profile`],
    queryFn: () => authFetcher(`/api/auth/profile`),
    retry: false,
  });

export const useLogoutMutation = () =>
  useMutation({
    mutationFn: () => {
      return authFetch('/api/auth/logout', { method: 'POST' });
    },
    onSuccess() {
      localStorage.removeItem('access_token');
      deleteCookie('dashboard_id');
      window.location.replace(window.location.origin);
    },
  });

export type Dashboard = {
  id: string,
  name: string,
  displayName: string,
  avatar: string,
  flags: string[]
}

export const useDashboards = () => useQuery<Dashboard[]>({
  queryKey: [`/api/auth/profile/dashboards`],
  queryFn: () => authFetcher(`/api/auth/profile/dashboards`),
  retry: false,
});