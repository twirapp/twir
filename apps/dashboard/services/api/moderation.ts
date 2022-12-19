import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import useSWR from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useModerationSettings = () => {
  const [selectedDashboard] = useSelectedDashboard();

  return useSWR<ChannelModerationSetting[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/moderation` : null,
    swrAuthFetcher,
  );
};
