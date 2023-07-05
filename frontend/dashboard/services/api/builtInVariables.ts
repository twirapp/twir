import { useQuery } from '@tanstack/react-query';
import { UnwrapPromise } from 'next/dist/lib/coalesced-function';

import { protectedApiClient } from '@/services/api/twirp';

export const useBuiltInVariables = () => useQuery<UnwrapPromise<ReturnType<typeof protectedApiClient.builtInVariablesGetAll>['response']>>({
	queryKey: ['variablesList'],
	queryFn: async () => {
		const call = await protectedApiClient.builtInVariablesGetAll({});

		return call.response;
	},
});
