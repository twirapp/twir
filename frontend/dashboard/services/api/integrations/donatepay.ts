import { useMutation, useQuery } from '@tanstack/react-query';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useDonatePayIntegration = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/integrations/donatepay`;

  return {
    useData: () => useQuery<string>({
      queryKey: [getUrl()],
      queryFn: () => authFetcher(getUrl()),
    }),
    usePost: () => useMutation<any, unknown, { apiKey: string }, unknown>({
      mutationFn: ({ apiKey }) => {
        return authFetcher(
          getUrl(),
          {
            body: JSON.stringify({ apiKey }),
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