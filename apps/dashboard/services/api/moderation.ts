import { useQuery } from '@tanstack/react-query';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api';
import { SELECTED_DASHBOARD_KEY, useSelectedDashboard } from '@/services/dashboard';

export const useModerationSettings = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/moderation`;

  return {
    getAll: useQuery<ChannelModerationSetting[]>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
    }),
  };
};
