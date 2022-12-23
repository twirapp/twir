import { useMutation, useQuery } from '@tanstack/react-query';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

export const useBotApi = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/bot`;

  return {
    isMod: () => useQuery<boolean>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(`${getUrl()}/checkmod`),
      refetchInterval: 1000,
    }),
    useChangeState: () => useMutation({
      mutationKey: [getUrl],
      mutationFn: (action: 'part' | 'join') => {
        return authFetcher(`${getUrl()}/connection`, {
          method: 'PATCH',
          body: JSON.stringify({ action }),
          headers: {
            'Content-Type': 'application/json',
          },
        });
      },
    }),
  };
};