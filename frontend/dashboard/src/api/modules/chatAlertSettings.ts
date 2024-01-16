import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query';
import { type ChatAlertsSettings } from '@twir/api/messages/modules_chat_alerts/modules_chat_alerts';
import { MaybeRef, unref } from 'vue';

import { protectedApiClient } from '../twirp';

export { ChatAlertsSettings };

const key = ['chatAlertsSettings'];

export const useChatAlertsSettings = () => useQuery({
	queryKey: key,
	queryFn: async () => {
		const call = await protectedApiClient.modulesChatAlertsGet({});

		return call.response;
	},
});

export const useChatAlertsSettingsUpdate = () => {
	const queryClient = useQueryClient();


	return useMutation({
		mutationKey: key,
		mutationFn: async (opts: MaybeRef<ChatAlertsSettings>) => {
			const data = unref(opts);
			await protectedApiClient.modulesChatAlertsUpdate(data);
		},
		onSuccess: () => {
			queryClient.invalidateQueries(key);
		},
	});
};
