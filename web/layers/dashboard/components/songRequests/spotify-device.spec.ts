import { mount } from '@vue/test-utils'
import { ref } from 'vue'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import en from '../../../../i18n/locales/en.json'

import SpotifyDevice from './spotify-device.vue'

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
		},
	},
}

const deviceFixture = { id: 'device-1', name: 'rach', type: 'Computer', isActive: true }

function mockApi(capabilities: Record<string, unknown>) {
	const refreshDevice = vi.fn().mockResolvedValue({ error: undefined })

	api.useSongRequestsApi.mockReturnValue({
		useSongRequestQuery: () => ({
			data: ref({ songRequests: { spotifyCapabilities: capabilities } }),
			fetching: ref(false),
			error: ref(undefined),
		}),
		useSpotifyRefreshDeviceMutation: () => ({ executeMutation: refreshDevice }),
	} as never)

	return { refreshDevice }
}

describe('SpotifyDevice', () => {
	beforeEach(() => {
		api.useSongRequestsApi.mockReset()
		document.body.replaceChildren()
	})

	it('renders nothing when spotify cannot be used', () => {
		mockApi({ connected: false, hasPlaybackScope: false, canUseSpotify: false })

		const wrapper = mount(SpotifyDevice, mountOptions)

		expect(wrapper.html()).toBe('<!--v-if-->')
	})

	it('shows the selected device name', () => {
		mockApi({
			connected: true,
			hasPlaybackScope: true,
			canUseSpotify: true,
			selectedDevice: deviceFixture,
			activeDevice: deviceFixture,
		})

		const wrapper = mount(SpotifyDevice, mountOptions)

		expect(wrapper.text()).toContain('rach')
		expect(wrapper.text()).toContain('Computer')
	})

	it('shows no-device hint when no device is cached', () => {
		mockApi({
			connected: true,
			hasPlaybackScope: true,
			canUseSpotify: true,
			selectedDevice: null,
			activeDevice: null,
		})

		const wrapper = mount(SpotifyDevice, mountOptions)

		expect(wrapper.text()).toContain('No active Spotify device')
	})

	it('calls refresh mutation on button click', async () => {
		const { refreshDevice } = mockApi({
			connected: true,
			hasPlaybackScope: true,
			canUseSpotify: true,
			selectedDevice: deviceFixture,
			activeDevice: deviceFixture,
		})

		const wrapper = mount(SpotifyDevice, mountOptions)
		await wrapper.get('button').trigger('click')

		expect(refreshDevice).toHaveBeenCalledTimes(1)
	})
})
