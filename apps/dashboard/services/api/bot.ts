import useSWR, { useSWRConfig } from 'swr';

import { swrAuthFetcher } from '@/services/api/fetchWrappers';
import { useSelectedDashboard } from '@/services/dashboard';

export const useBotApi = () => {
  const [selectedDashboard] = useSelectedDashboard();
  const { mutate } = useSWRConfig();

  return {
    isMod() {
      return useSWR<boolean>(
        selectedDashboard ? `/api/v1/channels/${selectedDashboard.channelId}/bot/checkmod` : null,
        swrAuthFetcher,
      );
    },
  };
};