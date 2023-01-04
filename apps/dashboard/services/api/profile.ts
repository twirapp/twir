import { useMutation, useQuery } from '@tanstack/react-query';
import { AuthUser } from '@tsuwari/shared';
import { deleteCookie } from 'cookies-next';
import { useContext } from 'react';

import { authFetch, authFetcher } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useProfile = () => useQuery<AuthUser & { apiKey: string }>({
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
      deleteCookie('selectedDashboard');
      window.location.replace(window.location.origin);
    },
  });
