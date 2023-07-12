import type { RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';

import { protectedApiClient } from '@/api/twirp.js';

type CallFunc<
	Req extends Record<any, any>,
	Res extends Record<any, any>
> = (input: Req, options?: RpcOptions) => UnaryCall<Req, Res>;

const createIntegrationOauth = <
		GetData extends CallFunc<any, any>,
		GetAuthLink extends CallFunc<any, any>,
		UsePostCode extends CallFunc<any, any>,
		UseLogout extends CallFunc<any, any>,
>(opts: {
	integrationName: string,
	getData: GetData,
	getAuthLink: GetAuthLink,
	usePostCode: UsePostCode,
	useLogout: UseLogout,
}) => {
	for (const [key, value] of Object.entries(opts)) {
		if (typeof value === 'function') {
			// eslint-disable-next-line @typescript-eslint/ban-ts-comment
			// @ts-ignore
			opts[key] = value.bind(protectedApiClient);
		}
	}

	const queryClient = useQueryClient();

	const queryKey = `integrations/${opts.integrationName}`;

	return {
		useData: () => useQuery<Awaited<ReturnType<typeof opts.getData>['response']>>({
			queryKey: [queryKey],
			queryFn: async () => {
				const call = await opts.getData({});
				return call.response;
			},
			retry: false,
		}),
		useAuthLink: () => useQuery<Awaited<ReturnType<typeof opts.getAuthLink>['response']>>({
			queryKey: [`${queryKey}/auth-link`],
			queryFn: async () => {
				const call = await opts.getAuthLink({});
				return call.response;
			},
		}),
		usePostCode: () => useMutation({
			mutationKey: [`${queryKey}/post-code`],
			mutationFn: async (req: Parameters<typeof opts.usePostCode>[0]) => {
				await opts.usePostCode(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([queryKey]);
			},
		}),
		useLogout: () => useMutation({
			mutationKey: [`${queryKey}/post-code`],
			mutationFn: async (req: Parameters<typeof opts.useLogout>[0]) => {
				await opts.useLogout(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([queryKey]);
			},
		}),
	};
};

export const useSpotifyIntegration = () => createIntegrationOauth({
	integrationName: 'spotify',
	getData: protectedApiClient.integrationsSpotifyGetData,
	getAuthLink: protectedApiClient.integrationsSpotifyGetAuthLink,
	usePostCode: protectedApiClient.integrationsSpotifyPostCode,
	useLogout: protectedApiClient.integrationsSpotifyLogout,
});
