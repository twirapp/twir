import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import BlocksRender from './blocks-render.vue'

describe('BlocksRender', () => {
	it('renders a constrained image that opens in a new tab', () => {
		const wrapper = mount(BlocksRender, {
			props: {
				data: {
					blocks: [
						{
							id: 'image-1',
							type: 'image',
							data: {
								url: 'https://cdn.example/image.png',
								caption: 'image.png',
							},
						},
					],
				},
			},
		})

		const link = wrapper.get('a')
		expect(link.attributes('href')).toBe('https://cdn.example/image.png')
		expect(link.attributes('target')).toBe('_blank')
		expect(link.attributes('rel')).toBe('noopener noreferrer')

		const image = wrapper.get('img')
		expect(image.attributes('alt')).toBe('image.png')
		expect(image.attributes('loading')).toBe('lazy')
		expect(image.classes()).toContain('max-h-96')
		expect(wrapper.get('figcaption').text()).toBe('image.png')
	})
})
