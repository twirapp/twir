import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import PlatformBindingCard from './platform-binding-card.vue'

describe('PlatformBindingCard', () => {
	it('offers OAuth connection for a disconnected platform', () => {
		const wrapper = mount(PlatformBindingCard, {
			props: {
				platform: 'TWITCH',
				presentation: { label: 'Twitch', icon: 'lucide:radio' },
				capabilities: [{ name: 'chat.write' }],
				binding: null,
			},
		})

		expect(wrapper.get('button').text()).toContain('Connect')
	})
})
