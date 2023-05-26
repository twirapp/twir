import { useContext } from 'react';
import { SelectedDashboardContext } from '../selectedDashboardProvider';
import { useMutation, useQuery } from '@tanstack/react-query';
import { ChannelCategoryAlias } from '@tsuwari/typeorm/entities/ChannelCategoryAlias';
import { authFetcher } from './fetchWrappers';
import { queryClient } from './queryClient';

export interface ChannelCategoryAliasDto {
	category: string;
	alias: string;
}
export const useCategoriesAliases = () => {
	const dashboard = useContext(SelectedDashboardContext);
	const getUrl = () => `/api/v1/channels/${dashboard.id}/categories-aliases`;

	return {
		useGet: () =>
			useQuery<ChannelCategoryAlias[]>({
				queryKey: [getUrl()],
				queryFn: () => authFetcher(getUrl()),
			}),
		useUpdate: () =>
			useMutation({
				mutationKey: [getUrl()],
				mutationFn: (data: ChannelCategoryAlias) => {
					return authFetcher(getUrl(), {
						body: JSON.stringify({ ...data }),
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
					});
				},
				onSuccess: (result) => {
					queryClient.setQueryData<ChannelCategoryAlias>([getUrl()], (old) => {
						return result;
					});
				},
			}),
	};
};
