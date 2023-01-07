import { useQuery } from '@tanstack/react-query';
import { type HelixCustomRewardData } from '@twurple/api/lib/interfaces/helix/channelPoints.external';
import { getCookie } from 'cookies-next';
import { useContext } from 'react';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';


export const useRewards = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/rewards`;

  return () => useQuery<HelixCustomRewardData[]>({
    queryKey: [getUrl()],
    queryFn: () => authFetcher(getUrl()),
  });
};