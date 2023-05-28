import { useMutation, useQuery } from '@tanstack/react-query';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { getCookie } from 'cookies-next';
import { useContext } from 'react';

import { authFetcher, queryClient } from '@/services/api';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';
import { ChannelStream } from '@tsuwari/typeorm/entities/ChannelStream';

export const useModerationSettings = () => {
	const dashboard = useContext(SelectedDashboardContext);
	const getUrl = () => `/api/v1/channels/${dashboard.id}/moderation`;

	return {
		useGet: () =>
			useQuery<ChannelModerationSetting[]>({
				queryKey: [getUrl()],
				queryFn: () => authFetcher(getUrl()),
			}),
		useUpdate: () =>
			useMutation({
				mutationKey: [getUrl()],
				mutationFn: (data: ChannelModerationSetting[]) => {
					return authFetcher(getUrl(), {
						body: JSON.stringify({ items: data }),
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
					});
				},
				onSuccess: (result) => {
					queryClient.setQueryData<ChannelModerationSetting[]>([getUrl()], (old) => {
						return result;
					});
				},
			}),
	};
};

export const useModerationManager = () => {
	const dashboard = useContext(SelectedDashboardContext);
	const getUrl = () => `/api/v1/channels/${dashboard.id}/moderation`;

	return {
		useUpdateTitle: () =>
			useMutation({
				mutationKey: [getUrl() + '/title'],
				mutationFn: (title: string) => {
					return authFetcher(getUrl() + '/title', {
						body: JSON.stringify({ title: title }),
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
					});
				},
				onSuccess: (result) => {
					queryClient.setQueryData<string>([getUrl() + '/title'], (old) => {
						return result;
					});
				},
			}),
		useUpdateCategory: () =>
			useMutation({
				mutationKey: [getUrl() + '/category'],
				mutationFn: (categoryId: string) => {
					return authFetcher(getUrl() + '/category', {
						body: JSON.stringify({ categoryId: categoryId }),
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
					});
				},
			}),
	};
};

export const useGetStream = () => {
	const dashboard = useContext(SelectedDashboardContext);
	const getUrl = () => `/api/v1/channels/${dashboard.id}/streams`;

	return {
		useGet: () =>
			useQuery<ChannelStream>({
				queryKey: [getUrl()],
				queryFn: () => authFetcher(getUrl()),
			}),
	};
};
