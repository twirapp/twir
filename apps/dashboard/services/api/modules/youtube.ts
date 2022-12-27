import { useMutation, useQuery } from '@tanstack/react-query';
import { V1 } from '@tsuwari/types/api';
import { getCookie } from 'cookies-next';

import { authFetcher } from '@/services/api';
import { SELECTED_DASHBOARD_KEY } from '@/services/dashboard';

type Youtube = V1['CHANNELS']['MODULES']['YouTube']

export const useYoutubeModule = () => {
  const getUrl = () => `/api/v1/channels/${getCookie(SELECTED_DASHBOARD_KEY)}/modules/youtube-sr`;
  
  return {
    useSettings: () => useQuery<Youtube['GET']>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
      retry: false,
    }),
    useSearch: () => useMutation({
      mutationFn: ({ query, type }: {query: string, type: 'channel' | 'video'}) => {
        return authFetcher(`${getUrl()}/search?type=${type}&query=${query}`);
      },
      mutationKey: [`${getUrl()}/search`],
    }),
    useUpdate: () => useMutation({
      mutationFn: (body: Youtube['POST']) => {
        return authFetcher(`${getUrl()}`, {
          method: 'POST',
          body: JSON.stringify(body),
          headers: {
            'Content-Type': 'application/json',
          },
        });
      },
      mutationKey: [getUrl()],
    }),
  };
};