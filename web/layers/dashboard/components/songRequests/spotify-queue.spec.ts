import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import { SpotifySongRequestStatus } from '~/gql/graphql.js'

import SpotifyQueue from './spotify-queue.vue'

const api = vi.hoisted(() => ({
	useSongRequestsApi: vi.fn(),
}))

vi.mock('~~/layers/dashboard/api/song-requests', () => api)

const i18n = createI18n<[typeof en], 'en'>({
	legacy: false,
	locale: 'en',
	messages: { en },
})

const mountOptions = {
	global: {
		plugins: [i18n],
		stubs: {
			NuxtIcon: { props: ['name'], template: '<i :data-icon="name" />' },
			NuxtLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
		},
	},
}

function createSettingsQuery(capabilities: Record<string, unknown>) {
	return {
		data: ref({
			songRequests: {
				mode: 'SPOTIFY',
				spotifyCapabilities: capabilities,
			},
		}),
		fetching: ref(false),
		error: ref(undefined),
	}
}

function createQueueQuery(requests: Array<Record<string, unknown>>, currentDevice: unknown = null) {
	return {
		data: ref({
			spotifySongRequestsQueue: {
				currentDevice,
				requests,
			},
		}),
		fetching: ref(false),
		error: ref(undefined),
	}
}

const connectedCapabilities = {
	connected: true,
	hasPlaybackScope: true,
	canUseSpotify: true,
	activeDevice: { id: 'device-1', name: 'Desktop', type: 'Computer', isActive: true },
	selectedDevice: { id: 'device-1', name: 'Desktop', type: 'Computer', isActive: true },
}

const queuedRequest = {
	id: 'request-1',
	title: 'Song',
	artist: 'Artist',
	album: 'Album',
	durationMs: 180000,
	requesterName: 'viewer',
	requesterDisplayName: 'Viewer',
	source: 'chat',
	queuePosition: 1,
	status: SpotifySongRequestStatus.Queued,
	createdAt: new Date().toISOString(),
}

function mockApi(capabilities: Record<string, unknown>, requests: Array<Record<string, unknown>>, currentDevice: unknown = null) {
	const mutations = {
		skip: vi.fn().mockResolvedValue({ error: undefined }),
		cancel: vi.fn().mockResolvedValue({ error: undefined }),
		refreshDevice: vi.fn().mockResolvedValue({ error: undefined }),
	}

	api.useSongRequestsApi.mockReturnValue({
		useSongRequestQuery: () => createSettingsQuery(capabilities),
		useSpotifyQueueQuery: () => createQueueQuery(requests, currentDevice),
		useSpotifySkipMutation: () => ({ executeMutation: mutations.skip }),
		useSpotifyCancelMutation: () => ({ executeMutation: mutations.cancel }),
		useSpotifyRefreshDeviceMutation: () => ({ executeMutation: mutations.refreshDevice }),
	} as never)

	return mutations
}

describe('SpotifyQueue', () => {
	beforeEach(() => {
		api.useSongRequestsApi.mockReset()
		document.body.replaceChildren()
	})

	it('prompts to connect Spotify when not connected', () => {
		mockApi({ connected: false, hasPlaybackScope: false, canUseSpotify: false }, [])

		const wrapper = mount(SpotifyQueue, mountOptions)

		expect(wrapper.text()).toContain('Spotify is not connected')
		expect(wrapper.text()).toContain('Connect Spotify')
	})

	it('prompts to reconnect when playback scope is missing', () => {
		mockApi({ connected: true, hasPlaybackScope: false, canUseSpotify: false }, [])

		const wrapper = mount(SpotifyQueue, mountOptions)

		expect(wrapper.text()).toContain('playback permissions are missing')
		expect(wrapper.text()).toContain('Reconnect Spotify')
	})

	it('shows the current device and refresh button when usable', () => {
		mockApi(connectedCapabilities, [], connectedCapabilities.activeDevice)

		const wrapper = mount(SpotifyQueue, mountOptions)

		expect(wrapper.text()).toContain('Desktop')
		expect(wrapper.text()).toContain('Refresh device')
	})

	it('calls refresh device mutation', async () => {
		const mutations = mockApi(connectedCapabilities, [], connectedCapabilities.activeDevice)

		const wrapper = mount(SpotifyQueue, mountOptions)
		const refreshButton = wrapper
			.findAll('button')
			.find((button) => button.text().includes('Refresh device'))
		expect(refreshButton).toBeDefined()
		await refreshButton!.trigger('click')

		expect(mutations.refreshDevice).toHaveBeenCalledTimes(1)
	})

	it('renders queue rows with status badges', () => {
		mockApi(connectedCapabilities, [
			queuedRequest,
			{ ...queuedRequest, id: 'request-2', status: SpotifySongRequestStatus.Playing },
		])

		const wrapper = mount(SpotifyQueue, mountOptions)

		expect(wrapper.text()).toContain('Song')
		expect(wrapper.text()).toContain('Artist')
		expect(wrapper.text()).toContain('Viewer')
		expect(wrapper.text()).toContain('Queued')
		expect(wrapper.text()).toContain('Playing')
	})

	it('shows empty state when there are no requests', () => {
		mockApi(connectedCapabilities, [])

		const wrapper = mount(SpotifyQueue, mountOptions)

		expect(wrapper.text()).toContain('No active Spotify requests')
	})

	it('pauses the queue query when Spotify cannot be used', () => {
		let capturedPause: unknown
		api.useSongRequestsApi.mockReturnValue({
			useSongRequestQuery: () =>
				createSettingsQuery({ connected: false, hasPlaybackScope: false, canUseSpotify: false }),
			useSpotifyQueueQuery: (pause: unknown) => {
				capturedPause = pause
				return createQueueQuery([])
			},
			useSpotifySkipMutation: () => ({ executeMutation: vi.fn() }),
			useSpotifyCancelMutation: () => ({ executeMutation: vi.fn() }),
			useSpotifyRefreshDeviceMutation: () => ({ executeMutation: vi.fn() }),
		} as never)

		mount(SpotifyQueue, mountOptions)

		expect(capturedPause).toBeDefined()
	})
})
