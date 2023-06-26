import { useMutation, useQuery } from '@tanstack/react-query';
import { V1 } from '@twir/types/api';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export type TTS = V1['CHANNELS']['MODULES']['TTS']


type UserSettings = {
  rate: number,
  pitch: number,
  voice: string,
  userLogin: string,
  userName: string,
  userAvatar: string,
  userId: string,
}

export const useTtsModule = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/modules/tts`;

  return {
    useInfo: () => useQuery({
      queryKey: [`${getUrl()}/info`],
      queryFn: () => authFetcher(`${getUrl()}/info`),
      retry: false,
    }),
    useSettings: () => useQuery<TTS['GET']>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
      retry: false,
    }),
    useUpdate: () => useMutation({
      mutationFn: (body: TTS['POST']) => {
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
      },
      mutationKey: [getUrl()],
    }),
    useUsersSettings: () => useQuery<UserSettings[]>({
      queryKey: [`${getUrl()}/users`],
      queryFn: () => authFetcher(`${getUrl()}/users`),
      retry: false,
    }),
    useUsersDelete: () => useMutation({
      mutationFn: (usersIds: string[]) => {
        return authFetcher(`${getUrl()}/users/`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ usersIds }),
        });
      },
      onSuccess: () => {
        queryClient.refetchQueries({ queryKey: [`${getUrl()}/users`] });
      },
      mutationKey: [`${getUrl()}/users`],
    }),
  };
};
