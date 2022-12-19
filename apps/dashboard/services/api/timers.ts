import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import useSWR from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useTimers = () => {
  const [selectedDashboard] = useSelectedDashboard();

  return useSWR<ChannelTimer[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/timers` : null,
    swrAuthFetcher,
  );
};
