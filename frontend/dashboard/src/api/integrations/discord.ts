import { useMutation, useQueries, useQuery, useQueryClient } from '@tanstack/vue-query';
import type {
	PostCodeRequest,
	UpdateMessage,
} from '@twir/api/messages/integrations_discord/integrations_discord';

import { protectedApiClient } from '@/api/twirp';

export const useDiscordIntegration = () => {
	const queryClient = useQueryClient();
	const dataQueryKey = ['integrationsDiscordGetData'];

	return {
		getConnectLink: () => useQuery({
			queryKey: ['integrationsDiscordGetConnectLink'],
			queryFn: async () => {
				const call = await protectedApiClient.integrationsDiscordGetAuthLink({});
				return call.response;
			},
		}),
		useData: () => useQuery({
			queryKey: dataQueryKey,
			queryFn: async () => {
				const call = await protectedApiClient.integrationsDiscordGetData({});
				return call.response;
			},
		}),
		updateData: () => useMutation({
			mutationKey: ['integrationsDiscordUpdater'],
			mutationFn: async (data: UpdateMessage) => {
				const call = await protectedApiClient.integrationsDiscordUpdate(data);
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(dataQueryKey);
			},
		}),
		usePostCode: () => useMutation({
			mutationKey: ['integrationsDiscordPostCode'],
			mutationFn: async (req: PostCodeRequest) => {
				const call = await protectedApiClient.integrationDiscordConnectGuild(req);
				return call.response;
			},
		}),
		disconnectGuild: () => useMutation({
			mutationKey: ['integrationsDiscordDisconnectGuild'],
			mutationFn: async (guildId: string) => {
				const call = await protectedApiClient.integrationsDiscordDisconnectGuild({ guildId });
				return call.response;
			},
			onSuccess: async () => {
				await queryClient.invalidateQueries(dataQueryKey);
			},
		}),
		getGuildChannels: (guildId: string) => useQuery({
			queryKey: ['integrationsDiscordGetGuildChannels', guildId],
			queryFn: async ({ queryKey }) => {
				return await getGuildChannelsFn(queryKey.at(1)!);
			},
		}),
		getGuildInfo: (guildId: string) => useQuery({
			queryKey: ['integrationsDiscordGetGuildInfo', guildId],
			queryFn: async ({ queryKey }) => {
				const call = await protectedApiClient.integrationsDiscordGetGuildInfo({
					guildId: queryKey[1],
				});
				return call.response;
			},
		}),
		getGuildsInfo: (guildsIds: string[]) => useQueries({
			queries: guildsIds.map((guildId) => ({
				queryKey: ['integrationsDiscordGetGuildInfo', guildId],
				queryFn: async ({ queryKey }: { queryKey: any[] }) => {
					const call = await protectedApiClient.integrationsDiscordGetGuildInfo({
						guildId: queryKey[1],
					});
					return call.response;
				},
			})),
		}),
	};
};

export const getGuildChannelsFn = async (guildId: string) => {
	const call = await protectedApiClient.integrationsDiscordGetGuildChannels({
		guildId: guildId,
	});
	return call.response;
};
