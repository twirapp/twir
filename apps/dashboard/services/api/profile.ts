import { useMutation, useQuery } from '@tanstack/react-query';
import { AuthUser } from '@tsuwari/shared';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api';
import { authFetch } from '@/services/api';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/moderation`;

export const useProfile = () => useQuery<AuthUser>({
  queryKey: [getUrl()],
  queryFn: () => authFetcher(getUrl()),
});

export const useLogoutMutation = () => useMutation({
  mutationFn: () => {
    return authFetch('/api/auth/logout', { method: 'POST' });
  },
  onSuccess() {
    localStorage.removeItem('access_token');
    window.location.replace(window.location.origin);
  },
});
