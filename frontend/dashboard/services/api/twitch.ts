import { useQuery } from '@tanstack/react-query';
import { UnwrapPromise } from 'next/dist/lib/coalesced-function';

import { unprotectedApiClient } from '@/services/api/twirp';

export const useTwitchUsers = (
	names: string[],
	ids: string[],
) => useQuery<UnwrapPromise<ReturnType<typeof unprotectedApiClient.twitchSearchUsers>['response']>>({
	queryFn: async () => {
		const call = await unprotectedApiClient.twitchSearchUsers({
			names,
			ids,
		});

		return call.response;
	},
});


