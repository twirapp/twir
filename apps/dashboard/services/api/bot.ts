import { useQuery } from '@tanstack/react-query';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

export const useBotApi = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/bot/checkmod`;

  return {
    isMod: useQuery<boolean>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
      refetchInterval: 1000,
    }),
  };
};