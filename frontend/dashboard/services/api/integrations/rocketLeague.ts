import { useMutation, useQuery } from '@tanstack/react-query';
import { useContext } from 'react';

import { queryClient } from '..';
import { authFetcher } from '../fetchWrappers';

import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useRocketLeagueIntegration = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/integrations/rocketleague`;
  
  return {
		useData: () =>
			useQuery<{ username: string; code: string }>({
				queryKey: [getUrl()],
				queryFn: () => authFetcher(getUrl()),
			}),
		usePost: () =>
			useMutation<any, unknown, { username: string; code: string }, unknown>({
        mutationFn: ({ username, code }) => {
          return authFetcher(
            getUrl(),
            {
              body: JSON.stringify({ username, code }),
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