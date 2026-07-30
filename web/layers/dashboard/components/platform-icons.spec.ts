import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import PlatformIcons from './platform-icons.vue'

function mountIcons(platforms: string[]) {
	return mount(PlatformIcons, {
		props: { platforms },
		global: {
			stubs: {
				NuxtIcon: {
					props: ['name'],
					template: '<i :data-icon="name" />',
				},
			},
		},
	})
}

describe('PlatformIcons', () => {
	it('renders nothing when the platforms list is empty', () => {
		// Given an entry available on every platform
		const wrapper = mountIcons([])

		// Then no container and no icons are rendered
		expect(wrapper.find('div').exists()).toBe(false)
		expect(wrapper.findAll('i')).toHaveLength(0)
	})

	it('renders one brand icon per known platform', () => {
		// Given an entry restricted to Twitch and Kick
		const wrapper = mountIcons(['twitch', 'kick'])

		// Then two icons render with the brand glyphs and tooltip labels
		const icons = wrapper.findAll('i')
		expect(icons).toHaveLength(2)
		expect(icons[0]!.attributes('data-icon')).toBe('simple-icons:twitch')
		expect(icons[0]!.attributes('title')).toBe('Twitch')
		expect(icons[1]!.attributes('data-icon')).toBe('simple-icons:kick')
		expect(icons[1]!.attributes('title')).toBe('Kick')
	})

	it('silently ignores unknown platform strings', () => {
		// Given a mix of a known and an unknown platform
		const wrapper = mountIcons(['twitch', 'myspace'])

		// Then only the known platform renders an icon
		const icons = wrapper.findAll('i')
		expect(icons).toHaveLength(1)
		expect(icons[0]!.attributes('data-icon')).toBe('simple-icons:twitch')
	})
})
