import { useMutation, useQuery } from '@tanstack/react-query';
import { AuthUser } from '@tsuwari/shared';
import { deleteCookie } from 'cookies-next';

import { authFetcher, FetcherError } from '@/services/api';
import { authFetch } from '@/services/api';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

export const useProfile = () => useQuery<AuthUser>({
  queryKey: [`/api/auth/profile`],
  queryFn: () => authFetcher(`/api/auth/profile`),
  retry: false,
});

export const useLogoutMutation = () => useMutation({
  mutationFn: () => {
    return authFetch('/api/auth/logout', { method: 'POST' });
  },
  onSuccess() {
    localStorage.removeItem('access_token');
    localStorage.removeItem(SELECTED_DASHBOARD_KEY);
    deleteCookie(SELECTED_DASHBOARD_KEY);
    window.location.replace(window.location.origin);
  },
});
