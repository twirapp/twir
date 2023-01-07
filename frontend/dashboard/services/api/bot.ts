import { useMutation, useQuery } from '@tanstack/react-query';
import { V1 } from '@tsuwari/types/api';
import { useContext } from 'react';

import { queryClient } from './queryClient';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useBotApi = () => {
  const dashboard = useContext(SelectedDashboardContext);

  const getUrl = () => `/api/v1/channels/${dashboard.id}/bot`;

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
