import { useMutation, useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed } from 'vue'

import { useIntegrationsPageData } from '#layers/dashboard/api/integrations/integrations-page.ts'
import { graphql } from '~/gql'

const DiscordGuildInfoQuery = graphql(`
	query DiscordIntegrationGuildInfo($guildId: String!) {
		discordIntegrationGuildInfo(guildId: $guildId) {
			id
			name
			icon
			channels {
				id
				name
				type
				canSendMessages
			}
			roles {
				id
				name
				color
			}
		}
	}
`)

const DiscordConnectGuildMutation = graphql(`
	mutation DiscordIntegrationConnectGuild($code: String!) {
		discordIntegrationConnectGuild(code: $code)
	}
`)

const DiscordDisconnectGuildMutation = graphql(`
	mutation DiscordIntegrationDisconnectGuild($guildId: String!) {
		discordIntegrationDisconnectGuild(guildId: $guildId)
	}
`)

const DiscordUpdateGuildMutation = graphql(`
	mutation DiscordIntegrationUpdateGuild($guildId: String!, $input: DiscordGuildUpdateInput!) {
		discordIntegrationUpdateGuild(guildId: $guildId, input: $input)
	}
`)

export const useDiscordIntegration = createGlobalState(() => {
	const integrationsPage = useIntegrationsPageData()

	const connectGuildMutation = useMutation(DiscordConnectGuildMutation)
	const disconnectGuildMutation = useMutation(DiscordDisconnectGuildMutation)
	const updateGuildMutation = useMutation(DiscordUpdateGuildMutation)

	const guilds = computed(() => integrationsPage.discordGuilds.value ?? [])
	const authLink = computed(() => integrationsPage.discordAuthLink.value ?? null)
	const isLoading = computed(() => integrationsPage.fetching.value)

	async function refetchData() {
		await integrationsPage.refetch()
	}

	async function connectGuild(code: string) {
		const result = await connectGuildMutation.executeMutation({ code })
		if (!result.error) {
			await refetchData()
		}
		return result
	}

	async function disconnectGuild(guildId: string) {
		const result = await disconnectGuildMutation.executeMutation({ guildId })
		if (!result.error) {
			await refetchData()
		}
		return result
	}

	async function updateGuild(guildId: string, input: {
		liveNotificationEnabled?: boolean
		liveNotificationChannelsIds?: string[]
		liveNotificationShowTitle?: boolean
		liveNotificationShowCategory?: boolean
		liveNotificationShowViewers?: boolean
		liveNotificationMessage?: string
		liveNotificationShowPreview?: boolean
		liveNotificationShowProfileImage?: boolean
		offlineNotificationMessage?: string
		shouldDeleteMessageOnOffline?: boolean
		additionalUsersIdsForLiveCheck?: string[]
	}) {
		const result = await updateGuildMutation.executeMutation({ guildId, input })
		if (!result.error) {
			await refetchData()
		}
		return result
	}

	return {
		guilds,
		authLink,
		isLoading,
		refetchData,
		connectGuild,
		disconnectGuild,
		updateGuild,
		connectGuildMutation,
		disconnectGuildMutation,
		updateGuildMutation,
	}
})

export function useDiscordGuildInfo(guildId: () => string | null) {
	const query = useQuery({
		query: DiscordGuildInfoQuery,
		variables: computed(() => ({
			guildId: guildId() ?? '',
		})),
		pause: computed(() => !guildId()),
	})

	const channels = computed(() => query.data.value?.discordIntegrationGuildInfo?.channels ?? [])
	const roles = computed(() => query.data.value?.discordIntegrationGuildInfo?.roles ?? [])
	const isLoading = computed(() => query.fetching.value)

	return {
		query,
		channels,
		roles,
		isLoading,
	}
}
