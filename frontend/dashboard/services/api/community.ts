import { useMutation, useQuery } from '@tanstack/react-query';
import { useContext } from 'react';

import { authFetcher } from '@/services/api/fetchWrappers';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

type User = {
  'id': string,
  'name': string,
  'displayName': string,
  'watched': number,
  'messages': number,
  'emotes': number,
  'avatarUrl': string,
}

export type SortyByField = 'messages' | 'watched' | 'emotes'

export const useCommunity = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/community`;

  return {
    useUsers: (limit= 50, page = 1, sortBy: SortyByField, order: 'asc' | 'desc' = 'desc') => useQuery<User[]>({
      queryKey: [`${getUrl()}/users`],
      queryFn: () => {
        const q = new URLSearchParams({
          limit: limit.toString(),
          page: page.toString(),
          sortBy,
          order,
        });

        return authFetcher(`${getUrl()}/users?${q}`);
      },
      retry: false,
      refetchInterval: 1000,
    }),
    useResetStats: () => useMutation({
      mutationFn: (field: SortyByField) => {
        return authFetcher(
          `${getUrl()}/users/stats`,
          {
            method: 'DELETE',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({
              field,
            }),
          },
        );
      },
      mutationKey: [`${getUrl()}/users`],
    }),
  };
};