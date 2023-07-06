import { useMutation, useQuery } from '@tanstack/react-query';
import { BotJoinPartRequest_Action } from '@twir/grpc/generated/api/api/bots';
import { UnwrapPromise } from 'next/dist/lib/coalesced-function';

import { queryClient } from '@/services/api/queryClient';
import { protectedApiClient } from '@/services/api/twirp';

const queryKey = ['botInfo'];

export const useBotInfo = () => useQuery<UnwrapPromise<ReturnType<typeof protectedApiClient.botInfo>['response']>>({
	queryKey,
	queryFn: async () => {
		const call = await protectedApiClient.botInfo({});
		return call.response;
	},
	refetchInterval: 4000,
});

export const useBotJoinPart = () => useMutation({
	onSuccess: async () => {
		await queryClient.invalidateQueries({ queryKey });
	},
	mutationFn: async (action: 'part' | 'join') => {
		return protectedApiClient.botJoinPart({
			action: action === 'join'
				? BotJoinPartRequest_Action.JOIN
				: BotJoinPartRequest_Action.LEAVE,
		});
	},
});
