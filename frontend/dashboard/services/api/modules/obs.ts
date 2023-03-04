import { useMutation, useQuery } from '@tanstack/react-query';
import { V1 } from '@tsuwari/types/api';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export type OBS = V1['CHANNELS']['MODULES']['OBS']

type OBSData = {
  sources: string[],
  scenes: string[],
  audioSources: string[],
}

export const useObsModule = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/modules/obs-websocket`;
  
  return {
    useSettings: () => useQuery<OBS['GET']>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
      retry: false,
    }),
    useUpdate: () => useMutation({
      mutationFn: (body: OBS['POST']) => {
        return authFetcher(`${getUrl()}`, {
          method: 'POST',
          body: JSON.stringify(body),
          headers: {
            'Content-Type': 'application/json',
          },
        });
      },
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: [getUrl()] });
        queryClient.refetchQueries({ queryKey: [`${getUrl()}/data`] });
      },
      mutationKey: [getUrl()],
    }),
    useData: () => useQuery<OBSData>({
      queryKey: [`${getUrl()}/data`],
      queryFn: () => authFetcher(`${getUrl()}/data`),
      retry: false,
    }),
  };
};