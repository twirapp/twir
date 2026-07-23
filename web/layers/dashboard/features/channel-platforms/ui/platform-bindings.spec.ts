import { flushPromises, mount } from '@vue/test-utils'
import { ref } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const channelPlatforms = vi.hoisted(() => ({
	useChannelPlatforms: vi.fn(),
}))
const sonner = vi.hoisted(() => ({
	toast: {
		error: vi.fn(),
	},
}))

vi.mock('../composables/use-channel-platforms.js', () => channelPlatforms)
vi.mock('vue-sonner', () => sonner)

import PlatformBindings from './platform-bindings.vue'

describe('PlatformBindings', () => {
	beforeEach(() => {
		channelPlatforms.useChannelPlatforms.mockReset()
		sonner.toast.error.mockReset()
	})

	it('sends a disconnected platform connect intent through the composable', async () => {
		const connect = vi.fn().mockResolvedValue(undefined)
		channelPlatforms.useChannelPlatforms.mockReturnValue({
			cards: ref([
				{
					platform: 'TWITCH',
					presentation: { label: 'Twitch', icon: 'lucide:radio' },
					capabilities: [{ name: 'chat.write' }],
					binding: null,
				},
			]),
			fetching: ref(false),
			error: ref(null),
			connect,
			disconnect: vi.fn(),
			setEnabled: vi.fn(),
		})

		const wrapper = mount(PlatformBindings, {
			global: {
				stubs: {
					NuxtIcon: true,
				},
			},
		})

		await wrapper.get('button').trigger('click')
		expect(connect).toHaveBeenCalledWith('TWITCH')
	})

	it('shows a mutation error after a failed platform action', async () => {
		const connect = vi.fn().mockResolvedValue(new Error('Connection failed'))
		channelPlatforms.useChannelPlatforms.mockReturnValue({
			cards: ref([
				{
					platform: 'TWITCH',
					presentation: { label: 'Twitch', icon: 'lucide:radio' },
					capabilities: [],
					binding: null,
				},
			]),
			fetching: ref(false),
			error: ref(null),
			connect,
			disconnect: vi.fn(),
			setEnabled: vi.fn(),
		})

		const wrapper = mount(PlatformBindings, {
			global: {
				stubs: {
					NuxtIcon: true,
				},
			},
		})

		await wrapper.get('button').trigger('click')
		await flushPromises()

		expect(sonner.toast.error).toHaveBeenCalledWith('Unable to update platform binding')
	})
})
