import { useMutation, useQuery } from '@tanstack/react-query';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useValorantIntegration = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/integrations/valorant`;

  return {
    useData: () => useQuery<{ username: string }>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
    }),
    usePost: () => useMutation<any, unknown, { username: string }, unknown>({
      mutationFn: ({ username }) => {
        return authFetcher(
          getUrl(),
          {
            body: JSON.stringify({ username }),
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
          },
        );
      },
      onSuccess: () => {
        queryClient.invalidateQueries([getUrl()]);
      },
      mutationKey: [getUrl()],
    }),
  };
};