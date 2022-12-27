import { useQuery } from '@tanstack/react-query';
import { type HelixCustomRewardData } from '@twurple/api/lib/interfaces/helix/channelPoints.external';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';


export const useRewards = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/rewards`;

  return () => useQuery<HelixCustomRewardData[]>({
    queryKey: [getUrl()],
    queryFn: () => authFetcher(getUrl()),
  });
};