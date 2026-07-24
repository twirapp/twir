import { computed } from 'vue'

import type { Platform } from '~/gql/graphql.js'

import { useChannelPlatformsApi } from '../api.js'

const platformPresentation = {
	TWITCH: { label: 'Twitch', icon: 'simple-icons:twitch', iconClass: 'text-[#9146FF]' },
	KICK: { label: 'Kick', icon: 'simple-icons:kick', iconClass: 'text-[#53FC18]' },
	VK_VIDEO_LIVE: { label: 'VK Video Live', icon: 'simple-icons:vk', iconClass: 'text-[#0077FF]' },
} satisfies Record<Platform, { label: string; icon: string; iconClass: string }>

export function useChannelPlatforms() {
	const api = useChannelPlatformsApi()
	const query = api.useQuery()
	const connectMutation = api.useConnect()
	const disconnectMutation = api.useDisconnect()
	const setEnabledMutation = api.useSetEnabled()

	const cards = computed(() => {
		const options = query.data.value?.channelPlatformOptions ?? []
		const bindings = query.data.value?.channelPlatformBindings ?? []

		return options.map((option) => ({
			platform: option.platform,
			presentation: platformPresentation[option.platform],
			capabilities: option.capabilities,
			binding: bindings.find((binding) => binding.platform === option.platform) ?? null,
		}))
	})

	async function connect(platform: Platform) {
		const result = await connectMutation.executeMutation({ platform })
		if (result.error || !result.data?.channelPlatformConnect) return result.error

		window.location.assign(result.data.channelPlatformConnect)
	}

	async function disconnect(platform: Platform) {
		const result = await disconnectMutation.executeMutation({ platform })
		if (!result.error) {
			await query.executeQuery({ requestPolicy: 'network-only' })
		}

		return result.error
	}

	async function setEnabled(platform: Platform, enabled: boolean) {
		const result = await setEnabledMutation.executeMutation({ platform, enabled })
		if (!result.error) {
			await query.executeQuery({ requestPolicy: 'network-only' })
		}

		return result.error
	}

	return {
		cards,
		fetching: query.fetching,
		error: query.error,
		connect,
		disconnect,
		setEnabled,
	}
}
