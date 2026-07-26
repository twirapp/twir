import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import { Platform } from '~/gql/graphql.js'

import PlatformSelector from './platform-selector.vue'

const vkPlatformId = Platform.VkVideoLive.toLowerCase()

function mountSelector(props: { modelValue?: string[]; exclude?: string[] } = {}) {
	return mount(PlatformSelector, {
		props: {
			modelValue: props.modelValue ?? [],
			...(props.exclude ? { exclude: props.exclude } : {}),
		},
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

describe('PlatformSelector', () => {
	it('shows the VK Video Live option by default with the VK brand icon', () => {
		// Given the default option set
		const wrapper = mountSelector()

		// When the selector renders
		const buttons = wrapper.findAll('button')

		// Then VK Video Live is offered with the simple-icons VK logo
		const vkButton = buttons.find((button) => button.text().includes('VK Video Live'))
		expect(vkButton, 'expected a VK Video Live option').toBeDefined()
		expect(vkButton!.get('i').attributes('data-icon')).toBe('simple-icons:vk')
	})

	it('omits the VK Video Live option when the exclude prop filters it out', () => {
		// Given a selector excluding the VK platform
		const wrapper = mountSelector({ exclude: [vkPlatformId] })

		// When the selector renders
		const labels = wrapper.findAll('button').map((button) => button.text())

		// Then VK is omitted while the other platforms stay
		expect(labels.some((label) => label.includes('VK Video Live'))).toBe(false)
		expect(labels.some((label) => label.includes('Twitch'))).toBe(true)
		expect(labels.some((label) => label.includes('Kick'))).toBe(true)
	})

	it('emits the generated VK platform value when the VK option is selected', async () => {
		// Given an empty selection
		const wrapper = mountSelector()

		// When the VK Video Live option is clicked
		const vkButton = wrapper.findAll('button').find((button) => button.text().includes('VK Video Live'))!
		await vkButton.trigger('click')

		// Then the model update carries the generated VK platform constant, lowercased for the wire
		const emitted = wrapper.emitted('update:modelValue')
		expect(emitted).toBeDefined()
		expect(emitted![0]).toEqual([[vkPlatformId]])
	})

	it('marks the VK Video Live option active when the model already contains it', () => {
		// Given a model that already includes the VK platform
		const wrapper = mountSelector({ modelValue: [vkPlatformId] })

		// When the selector renders
		const vkButton = wrapper.findAll('button').find((button) => button.text().includes('VK Video Live'))!

		// Then the VK option renders in the active state
		expect(vkButton.attributes('data-active')).toBe('true')
	})
})
