import { ref } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { Platform } from '~/gql/graphql.js'

const api = vi.hoisted(() => ({
	useChannelPlatformsApi: vi.fn(),
}))

vi.mock('../api.js', () => api)

import { useChannelPlatforms } from './use-channel-platforms.js'

function createQuery() {
	return {
		data: ref({
			channelPlatformBindings: [
				{
					platform: Platform.Kick,
					capabilities: [{ name: 'chat.read' }],
				},
			],
			channelPlatformOptions: [
				{
					platform: Platform.Twitch,
					capabilities: [],
				},
				{
					platform: Platform.Kick,
					capabilities: [{ name: 'chat.read' }],
				},
			],
		}),
		fetching: ref(false),
		error: ref(undefined),
		executeQuery: vi.fn(),
	}
}

describe('useChannelPlatforms', () => {
	beforeEach(() => {
		api.useChannelPlatformsApi.mockReset()
	})

	afterEach(() => {
		vi.restoreAllMocks()
	})

	it('maps cards in server option order and starts OAuth with the returned URL', async () => {
		const query = createQuery()
		const executeConnect = vi.fn().mockResolvedValue({
			data: { channelPlatformConnect: 'https://id.example/authorize' },
		})
		const assign = vi.spyOn(window.location, 'assign').mockImplementation(() => {})

		api.useChannelPlatformsApi.mockReturnValue({
			useQuery: () => query,
			useConnect: () => ({ executeMutation: executeConnect }),
			useDisconnect: () => ({ executeMutation: vi.fn() }),
			useSetEnabled: () => ({ executeMutation: vi.fn() }),
		} as never)

		const platforms = useChannelPlatforms()
		await platforms.connect(Platform.Twitch)

		expect(platforms.cards.value.map(({ platform }) => platform)).toEqual([
			Platform.Twitch,
			Platform.Kick,
		])
		expect(executeConnect).toHaveBeenCalledWith({ platform: Platform.Twitch })
		expect(assign).toHaveBeenCalledWith('https://id.example/authorize')
	})

	it('does not navigate when the connect mutation returns an error', async () => {
		const query = createQuery()
		const executeConnect = vi.fn().mockResolvedValue({ error: new Error('Connection failed') })
		const assign = vi.spyOn(window.location, 'assign').mockImplementation(() => {})

		api.useChannelPlatformsApi.mockReturnValue({
			useQuery: () => query,
			useConnect: () => ({ executeMutation: executeConnect }),
			useDisconnect: () => ({ executeMutation: vi.fn() }),
			useSetEnabled: () => ({ executeMutation: vi.fn() }),
		} as never)

		const platforms = useChannelPlatforms()
		await platforms.connect(Platform.Twitch)

		expect(assign).not.toHaveBeenCalled()
	})
})
