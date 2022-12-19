import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import useSWR from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useKeywords = () => {
  const [selectedDashboard] = useSelectedDashboard();

  return useSWR<ChannelKeyword[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/keywords` : null,
    swrAuthFetcher,
  );
};
