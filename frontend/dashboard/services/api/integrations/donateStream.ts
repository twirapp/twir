import { useMutation, useQuery } from '@tanstack/react-query';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';

export const useDonateStreamIntegration = () => {
	const dashboard = useContext(SelectedDashboardContext);
	const getUrl = () => `/api/v1/channels/${dashboard.id}/integrations/donate-stream`;

	return {
		useData: () =>
			useQuery<string>({
				queryKey: [getUrl()],
				queryFn: () => authFetcher(getUrl()),
			}),
		usePost: () =>
			useMutation<any, unknown, { id: string; secret: string }, unknown>({
				mutationFn: ({ id, secret }) => {
					return authFetcher(`${getUrl()}/${id}`, {
						body: JSON.stringify({ secret }),
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
					});
				},
				onSuccess: () => {
					queryClient.invalidateQueries([getUrl()]);
				},
				mutationKey: [getUrl()],
			}),
	};
};
