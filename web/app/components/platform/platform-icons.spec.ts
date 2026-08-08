import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import { Platform } from '~/gql/graphql.js'

import PlatformIcons from './platform-icons.vue'

function mountIcons(platforms: Platform[]) {
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

	it('renders one brand icon per platform', () => {
		// Given an entry restricted to Twitch and Kick
		const wrapper = mountIcons([Platform.Twitch, Platform.Kick])

		// Then two icons render with the brand glyphs and tooltip labels
		const icons = wrapper.findAll('i')
		expect(icons).toHaveLength(2)
		expect(icons[0]!.attributes('data-icon')).toBe('simple-icons:twitch')
		expect(icons[0]!.attributes('title')).toBe('Twitch')
		expect(icons[1]!.attributes('data-icon')).toBe('simple-icons:kick')
		expect(icons[1]!.attributes('title')).toBe('Kick')
	})

	it('supports every platform of the enum', () => {
		// Given an entry restricted to all four platforms
		const wrapper = mountIcons([
			Platform.Twitch,
			Platform.Kick,
			Platform.VkVideoLive,
			Platform.Youtube,
		])

		// Then all four brand icons render
		expect(wrapper.findAll('i')).toHaveLength(4)
	})
})
