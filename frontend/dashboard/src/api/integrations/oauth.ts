import type { RpcOptions, UnaryCall } from '@protobuf-ts/runtime-rpc';
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query';

import { createIntegrationOauth } from './oauth';

import { protectedApiClient } from '@/api/twirp.js';

type CallFunc<
	Req extends Record<any, any>,
	Res extends Record<any, any>
> = (input: Req, options?: RpcOptions) => UnaryCall<Req, Res>;

export const createIntegrationOauth = <
	GetData extends CallFunc<any, any>,
	GetAuthLink extends CallFunc<any, any>,
	UsePostCode extends CallFunc<any, any>,
	UseLogout extends CallFunc<any, any>,
	UpdateData extends CallFunc<any, any>,
>(opts: {
	integrationName: string,
	getData: GetData,
	getAuthLink: GetAuthLink,
	usePostCode: UsePostCode,
	useLogout: UseLogout,
	updateData?: UpdateData,
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
			retry: false,
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
			mutationKey: [`${queryKey}/logout`],
			mutationFn: async (req: Parameters<typeof opts.useLogout>[0]) => {
				await opts.useLogout(req);
			},
			onSuccess: () => {
				queryClient.invalidateQueries([queryKey]);
			},
		}),
		update: opts.updateData ? () => useMutation({
			mutationKey: [`${queryKey}/update`],
			mutationFn: async (req: Parameters<typeof opts.updateData>[0]) => {
				const call = await opts.updateData!(req);
				return call.response;
			},
		}) : undefined,
	};
};

export const useSpotifyIntegration = () => createIntegrationOauth({
	integrationName: 'spotify',
	getData: protectedApiClient.integrationsSpotifyGetData,
	getAuthLink: protectedApiClient.integrationsSpotifyGetAuthLink,
	usePostCode: protectedApiClient.integrationsSpotifyPostCode,
	useLogout: protectedApiClient.integrationsSpotifyLogout,
});

export const useLastfmIntegration = () => createIntegrationOauth({
	integrationName: 'lastfm',
	getData: protectedApiClient.integrationsLastFMGetData,
	getAuthLink: protectedApiClient.integrationsLastFMGetAuthLink,
	usePostCode: protectedApiClient.integrationsLastFMPostCode,
	useLogout: protectedApiClient.integrationsLastFMLogout,
});


export const useVKIntegration = () => createIntegrationOauth({
	integrationName: 'vk',
	getData: protectedApiClient.integrationsVKGetData,
	getAuthLink: protectedApiClient.integrationsVKGetAuthLink,
	usePostCode: protectedApiClient.integrationsVKPostCode,
	useLogout: protectedApiClient.integrationsVKLogout,
});


export const useStreamlabsIntegration = () => createIntegrationOauth({
	integrationName: 'streamlabs',
	getData: protectedApiClient.integrationsStreamlabsGetData,
	getAuthLink: protectedApiClient.integrationsStreamlabsGetAuthLink,
	usePostCode: protectedApiClient.integrationsStreamlabsPostCode,
	useLogout: protectedApiClient.integrationsStreamlabsLogout,
});

export const useFaceitIntegration = () => createIntegrationOauth({
	integrationName: 'faceit',
	getData: protectedApiClient.integrationsFaceitGetData,
	getAuthLink: protectedApiClient.integrationsFaceitGetAuthLink,
	usePostCode: protectedApiClient.integrationsFaceitPostCode,
	useLogout: protectedApiClient.integrationsFaceitLogout,
	updateData: protectedApiClient.integrationsFaceitUpdate,
});

export const useDonationAlertsIntegration = () => createIntegrationOauth({
	integrationName: 'donationalerts',
	getData: protectedApiClient.integrationsDonationAlertsGetData,
	getAuthLink: protectedApiClient.integrationsDonationAlertsGetAuthLink,
	usePostCode: protectedApiClient.integrationsDonationAlertsPostCode,
	useLogout: protectedApiClient.integrationsDonationAlertsLogout,
});

export const useValorantIntegration = () => createIntegrationOauth({
	integrationName: 'valorant',
	getData: protectedApiClient.integrationsValorantGetData,
	getAuthLink: protectedApiClient.integrationsValorantGetAuthLink,
	usePostCode: protectedApiClient.integrationsValorantPostCode,
	useLogout: protectedApiClient.integrationsValorantLogout,
});

export const useNightbotIntegration = () => createIntegrationOauth({
	integrationName: 'nightbot',
	getData: protectedApiClient.integrationsNightbotGetData,
	getAuthLink: protectedApiClient.integrationsNightbotGetAuthLink,
	usePostCode: protectedApiClient.integrationsNightbotPostCode,
	useLogout: protectedApiClient.integrationsNightbotLogout,
});

