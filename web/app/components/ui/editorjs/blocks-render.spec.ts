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

	it('lays images out in a wrapping row while text spans the full width', () => {
		const wrapper = mount(BlocksRender, {
			props: {
				data: {
					blocks: [
						{
							id: 'text-1',
							type: 'paragraph',
							data: { text: 'Release notes' },
						},
						{
							id: 'image-1',
							type: 'image',
							data: { url: 'https://cdn.example/one.png' },
						},
						{
							id: 'image-2',
							type: 'image',
							data: { url: 'https://cdn.example/two.png' },
						},
					],
				},
			},
		})

		expect(wrapper.get('div').classes()).toContain('flex-wrap')
		expect(wrapper.get('p').classes()).toContain('w-full')
		expect(wrapper.findAll('figure')).toHaveLength(2)
	})

	it('renders Discord inline formatting and an interactive spoiler', () => {
		const wrapper = mount(BlocksRender, {
			props: {
				data: {
					blocks: [
						{
							id: 'text-1',
							type: 'paragraph',
							data: {
								text: '<strong>bold</strong> <em>italic</em> <span class="discord-spoiler" tabindex="0">secret</span>',
							},
						},
						{
							id: 'quote-1',
							type: 'quote',
							data: { text: 'quoted text' },
						},
					],
				},
			},
		})

		expect(wrapper.get('strong').text()).toBe('bold')
		expect(wrapper.get('em').text()).toBe('italic')
		expect(wrapper.get('.discord-spoiler').attributes('tabindex')).toBe('0')
		expect(wrapper.get('blockquote.bq').text()).toBe('quoted text')
	})
})
