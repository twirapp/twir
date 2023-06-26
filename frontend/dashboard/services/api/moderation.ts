import { useMutation, useQuery } from '@tanstack/react-query';
import { ChannelModerationSetting } from '@twir/typeorm/entities/ChannelModerationSetting';
import { getCookie } from 'cookies-next';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useModerationSettings = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/moderation`;

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
