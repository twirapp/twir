import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import { afterEach, describe, expect, it } from 'vitest'

import PlatformBindingCard from './platform-binding-card.vue'

const connectedTwitchBinding = {
	id: 'binding-twitch',
	platform: 'TWITCH',
	userId: 'user-1',
	platformChannelId: 'channel-twitch',
	enabled: true,
	platformUserId: 'twitch-user',
	platformLogin: 'twitch_login',
	platformDisplayName: 'Twitch User',
	platformAvatar: 'https://example.com/twitch-avatar.png',
	capabilities: [{ name: 'chat.write' }],
}

const connectedKickBinding = {
	...connectedTwitchBinding,
	id: 'binding-kick',
	platform: 'KICK',
	platformChannelId: 'channel-kick',
	platformUserId: 'kick-user',
	platformLogin: 'kick_login',
	platformDisplayName: 'Kick User',
}

const mountOptions = {
	global: {
		stubs: {
			NuxtIcon: true,
		},
	},
}

describe('PlatformBindingCard', () => {
	afterEach(() => {
		document.body.replaceChildren()
	})

	it('emits a generic OAuth connection intent for a disconnected platform', async () => {
		const wrapper = mount(PlatformBindingCard, {
			...mountOptions,
			props: {
				platform: 'TWITCH',
				presentation: { label: 'Twitch', icon: 'lucide:radio' },
				capabilities: [{ name: 'chat.write' }],
				binding: null,
			},
		})

		expect(wrapper.get('button').text()).toContain('Connect')
		expect(wrapper.text()).toContain('chat.write')
		await wrapper.get('button').trigger('click')
		expect(wrapper.emitted('connect')).toEqual([['TWITCH']])
	})

	it('shows the connected account, capabilities, and a disconnect action', () => {
		const wrapper = mount(PlatformBindingCard, {
			...mountOptions,
			props: {
				platform: 'TWITCH',
				presentation: { label: 'Twitch', icon: 'lucide:radio' },
				capabilities: [{ name: 'chat.write' }],
				binding: connectedTwitchBinding,
			},
		})

		expect(wrapper.text()).toContain(connectedTwitchBinding.platformDisplayName)
		expect(wrapper.text()).toContain(connectedTwitchBinding.platformLogin)
		expect(wrapper.text()).toContain('chat.write')
		expect(wrapper.text()).toContain('Disconnect')
		expect(wrapper.get('img').attributes('src')).toBe(connectedTwitchBinding.platformAvatar)
	})

	it('confirms disconnect and emits the platform', async () => {
		const wrapper = mount(PlatformBindingCard, {
			...mountOptions,
			attachTo: document.body,
			props: {
				platform: 'TWITCH',
				presentation: { label: 'Twitch', icon: 'lucide:radio' },
				capabilities: [],
				binding: connectedTwitchBinding,
			},
		})

		const disconnectTrigger = wrapper
			.findAll('button')
			.find((button) => button.text().trim() === 'Disconnect')
		expect(disconnectTrigger).toBeDefined()
		await disconnectTrigger!.trigger('click')
		await nextTick()

		const disconnectActions = Array.from(document.body.querySelectorAll('button')).filter(
			(button) => button.textContent?.trim() === 'Disconnect',
		)
		expect(disconnectActions).toHaveLength(2)
		disconnectActions[1]?.click()
		await nextTick()

		expect(wrapper.emitted('disconnect')).toEqual([['TWITCH']])
	})

	it('keeps the generic enabled switch for a capability-less binding', async () => {
		const wrapper = mount(PlatformBindingCard, {
			...mountOptions,
			props: {
				platform: 'KICK',
				presentation: { label: 'Kick', icon: 'lucide:circle-play' },
				capabilities: [],
				binding: connectedKickBinding,
			},
		})

		expect(wrapper.text()).toContain('Enable bot')
		const enabledSwitch = wrapper.get('[role="switch"]')
		await enabledSwitch.trigger('click')
		expect(wrapper.emitted('setEnabled')).toEqual([['KICK', false]])
	})
})
