import { AuthUser } from '@tsuwari/shared';
import useSWR from 'swr';
import useSWRMutation from 'swr/mutation';

import { swrAuthFetcher } from '@/services/api';
import { authFetch } from '@/services/api';

export const useProfile = () => {
  return useSWR<AuthUser>('/api/auth/profile', swrAuthFetcher);
};

export const useLogoutMutation = () =>
  useSWRMutation('/api/user', async () => {
    const res = await authFetch('/api/auth/logout', { method: 'POST' });

    if (res.ok) {
      localStorage.removeItem('access_token');
      window.location.replace(window.location.origin);
    }
  });
