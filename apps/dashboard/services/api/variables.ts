import { ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import useSWR from 'swr';

import { swrAuthFetcher } from '@/services/api';
import { useSelectedDashboard } from '@/services/dashboard';

export const useVariables = () => {
  const [selectedDashboard] = useSelectedDashboard();

  return useSWR<ChannelCustomvar[]>(
    selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/variables` : null,
    swrAuthFetcher,
  );
};
