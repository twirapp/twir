import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';
import type { BotJoinPartRequest_Action } from '@twir/api/messages/bots/bots';

import { protectedApiClient } from './twirp.js';

const queryKey = ['botInfo'];


export const useBotInfo = () => useQuery({
	queryKey,
	queryFn: async () => {
		const call = await protectedApiClient.botInfo({});
		return call.response;
	},
	refetchInterval: 4000,
});

type Action = 'join' | 'part'

export const useBotJoinPart = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (action: Action) => {
			const call = await protectedApiClient.botJoinPart({
				action: action === 'join' ? BotJoinPartRequest_Action.JOIN : BotJoinPartRequest_Action.LEAVE,
			});
			return call.response;
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey });
		},
	});
};
