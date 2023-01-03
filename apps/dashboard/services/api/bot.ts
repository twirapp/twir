import { useMutation, useQuery } from '@tanstack/react-query';
import { V1 } from '@tsuwari/types/api';
import { getCookie } from 'cookies-next';
import { authFetcher } from '@/services/api/fetchWrappers';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';
import { queryClient } from './queryClient';

export const useBotApi = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/bot`;

  return {
    botInfo: () =>
      useQuery<V1['CHANNELS']['BOT']['GET']>({
        queryKey: [getUrl()],
        queryFn: () => authFetcher(`${getUrl()}`),
        refetchInterval: 4000,
      }),
    useChangeState: () =>
      useMutation({
        onSuccess: async () => {
          await queryClient.invalidateQueries({ queryKey: [getUrl()] });
        },
        mutationKey: [getUrl()],
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
