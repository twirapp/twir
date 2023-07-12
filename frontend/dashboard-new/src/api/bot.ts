import { useQuery } from '@tanstack/vue-query';
import { BotInfo } from '@twir/grpc/generated/api/api/bots';

import { protectedApiClient } from './twirp.js';

const queryKey = ['botInfo'];

export const useBotInfo = () => useQuery<BotInfo>({
	queryKey,
	queryFn: async () => {
		const call = await protectedApiClient.botInfo({});
		return call.response;
	},
	refetchInterval: 4000,
});
