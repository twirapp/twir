import { useMutation, useQuery } from '@tanstack/react-query';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { getCookie } from 'cookies-next';

import { authFetcher, queryClient } from '@/services/api';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

export const useModerationSettings = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/moderation`;

  return {
    useGet: () => useQuery<ChannelModerationSetting[]>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
    }),
    useUpdate: () => useMutation({
      mutationKey: [getUrl()],
      mutationFn: (data: ChannelModerationSetting[]) => {
        return authFetcher(
          getUrl(),
          {
            body: JSON.stringify({ items: data }),
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: (result) => {
        queryClient.setQueryData<ChannelModerationSetting[]>([getUrl()], old => {
          return result;
        });
      },
    }),
  };
};
