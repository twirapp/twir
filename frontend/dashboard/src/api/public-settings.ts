import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import type { Settings, UpdateRequest } from '@twir/api/messages/channels_public_settings/channels_public_settings';

import { useProfile } from '@/api/auth';
import { protectedApiClient, unprotectedApiClient } from '@/api/twirp';

export const usePublicSettings = () => {
	const queryClient = useQueryClient();
	const { data: profile } = useProfile();

	return {
		useGet: () => useQuery({
			queryKey: ['usePublicSettings'],
			queryFn: async (): Promise<Settings> => {
				const call = await unprotectedApiClient.getPublicSettings({
					channelId: profile.value!.selectedDashboardId,
				});
				return call.response;
			},
			enabled: () => !!profile.value,
		}),
		useUpdate: () => useMutation({
			mutationKey: ['usePublicSettings'],
			mutationFn: async (data: UpdateRequest) => {
				const call = await protectedApiClient.channelsPublicSettingsUpdate(data);
				return call.response;
			},
			async onSuccess() {
				await queryClient.invalidateQueries(['usePublicSettings']);
			},
		}),
	};
};
