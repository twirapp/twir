import { flushPromises, mount } from '@vue/test-utils'
import { ref } from 'vue'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import { SongRequestMode } from '~/gql/graphql.js'

import SettingsModal from './settings.vue'

const songRequestsApi = vi.hoisted(() => ({
	useSongRequestsApi: vi.fn(),
}))

const commandsApi = vi.hoisted(() => ({
	useCommandsApi: vi.fn(),
}))

vi.mock('~~/layers/dashboard/api/song-requests', () => songRequestsApi)
vi.mock('~~/layers/dashboard/api/commands/commands', () => commandsApi)
vi.mock('#build/fetch.mjs', () => ({ $fetch: vi.fn() }))

const i18n = createI18n<[typeof en], 'en'>({
	legacy: false,
	locale: 'en',
	messages: { en },
})

function createSettingsFixture(mode: SongRequestMode) {
	return {
		enabled: true,
		mode,
		acceptOnlyWhenOnline: true,
		maxRequests: 500,
		channelPointsRewardId: '',
		announcePlay: true,
		neededVotesForSkip: 30,
		user: {
			maxRequests: 20,
			minWatchTime: 0,
			minMessages: 0,
			minFollowTime: 0,
		},
		song: {
			minLength: 0,
			maxLength: 10,
			minViews: 50000,
			acceptedCategories: [],
		},
		denyList: {
			users: [],
			songs: [],
			channels: [],
			artistsNames: [],
			words: [],
		},
		translations: {
			nowPlaying: 'Now playing',
			notEnabled: 'Not enabled',
			noText: 'No text',
			acceptOnlyWhenOnline: 'Online only',
			user: {
				denied: 'user.denied',
				maxRequests: 'user.maxRequests',
				minMessages: 'user.minMessages',
				minWatched: 'user.minWatched',
				minFollow: 'user.minFollow',
			},
			song: {
				denied: 'song.denied',
				notFound: 'song.notFound',
				alreadyInQueue: 'song.alreadyInQueue',
				ageRestrictions: 'song.ageRestrictions',
				cannotGetInformation: 'song.cannotGetInformation',
				live: 'song.live',
				maxLength: 'song.maxLength',
				minLength: 'song.minLength',
				requestedMessage: 'song.requestedMessage',
				maximumOrdered: 'song.maximumOrdered',
				minViews: 'song.minViews',
			},
			channel: {
				denied: 'channel.denied',
			},
		},
		takeSongFromDonationMessages: false,
		playerNoCookieMode: false,
		volume: 30,
		channelApiKey: 'secret-api-key',
		spotifyCapabilities: {
			connected: true,
			hasPlaybackScope: true,
			canUseSpotify: true,
			activeDevice: null,
			selectedDevice: null,
		},
	}
}

function mountSettings(mode: SongRequestMode) {
	const executeMutation = vi.fn().mockResolvedValue({ error: undefined })
	const executeSearch = () => ({
		data: ref({ songRequestsSearchChannelOrVideo: { items: [] } }),
		fetching: ref(false),
		error: ref(undefined),
	})

	songRequestsApi.useSongRequestsApi.mockReturnValue({
		useSongRequestQuery: () => ({
			data: ref({ songRequests: createSettingsFixture(mode) }),
			fetching: ref(false),
			error: ref(undefined),
		}),
		useSongRequestMutation: () => ({ executeMutation }),
		useYoutubeVideoOrChannelSearch: executeSearch,
		useSpotifyQueueQuery: () => ({ data: ref(null), fetching: ref(false) }),
		useSpotifyTrackSearch: () => ({ data: ref(null), fetching: ref(false) }),
		useSpotifySelectDeviceMutation: () => ({ executeMutation: vi.fn() }),
		useSpotifyRefreshDeviceMutation: () => ({ executeMutation: vi.fn() }),
		useSpotifySkipMutation: () => ({ executeMutation: vi.fn() }),
		useSpotifyCancelMutation: () => ({ executeMutation: vi.fn() }),
	} as never)

	commandsApi.useCommandsApi.mockReturnValue({
		useQueryCommands: () => ({ data: ref({ commands: [] }) }),
	} as never)

	const wrapper = mount(SettingsModal, {
		attachTo: document.body,
		props: { open: true },
		global: {
			plugins: [i18n],
			stubs: {
				NuxtIcon: { props: ['name'], template: '<i :data-icon="name" />' },
				TwitchSearchUsers: true,
				CommandsList: true,
				RewardsSelector: true,
			},
		},
	})

	return { wrapper, executeMutation }
}

describe('SongRequests settings modal', () => {
	beforeEach(() => {
		songRequestsApi.useSongRequestsApi.mockReset()
		commandsApi.useCommandsApi.mockReset()
		document.body.replaceChildren()
	})

	it('saves the selected mode and strips capabilities and api key', async () => {
		const { executeMutation } = mountSettings(SongRequestMode.Spotify)
		await flushPromises()

		const saveButton = Array.from(document.body.querySelectorAll('button')).find((button) =>
			button.textContent?.includes('Save'),
		)
		expect(saveButton).toBeDefined()
		saveButton!.dispatchEvent(new MouseEvent('click', { bubbles: true }))
		await flushPromises()

		expect(executeMutation).toHaveBeenCalledTimes(1)
		const payload = executeMutation.mock.calls[0]![0] as { opts: Record<string, unknown> }
		expect(payload.opts['mode']).toBe(SongRequestMode.Spotify)
		expect(payload.opts).not.toHaveProperty('spotifyCapabilities')
		expect(payload.opts).not.toHaveProperty('channelApiKey')
	})

	it('keeps youtube mode in the payload when mode is youtube', async () => {
		const { executeMutation } = mountSettings(SongRequestMode.Youtube)
		await flushPromises()

		const saveButton = Array.from(document.body.querySelectorAll('button')).find((button) =>
			button.textContent?.includes('Save'),
		)
		expect(saveButton).toBeDefined()
		saveButton!.dispatchEvent(new MouseEvent('click', { bubbles: true }))
		await flushPromises()

		expect(executeMutation).toHaveBeenCalledTimes(1)
		const payload = executeMutation.mock.calls[0]![0] as { opts: Record<string, unknown> }
		expect(payload.opts['mode']).toBe(SongRequestMode.Youtube)
	})
})
