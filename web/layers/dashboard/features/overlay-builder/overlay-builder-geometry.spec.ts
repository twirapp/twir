import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import { useOverlayBuilder } from './composables/useOverlayBuilder'
import { type Layer, createLayerSettings } from './types'

import { ChannelOverlayLayerType } from '~/gql/graphql.js'

const i18n = createI18n<[typeof en], 'en'>({
	legacy: false,
	locale: 'en',
	messages: { en },
})

function setupBuilder() {
	let api: ReturnType<typeof useOverlayBuilder> | undefined
	mount(
		defineComponent({
			setup() {
				api = useOverlayBuilder()
				return () => null
			},
		}),
		{ global: { plugins: [i18n] } }
	)
	if (!api) throw new Error('builder was not initialized')
	return api
}

function makeLayer(id: string, overrides: Partial<Layer> = {}): Layer {
	return {
		id,
		type: ChannelOverlayLayerType.Text,
		name: id,
		posX: 10,
		posY: 10,
		width: 100,
		height: 100,
		rotation: 0,
		opacity: 1,
		visible: true,
		locked: false,
		zIndex: 0,
		periodicallyRefetchData: false,
		settings: createLayerSettings(),
		...overrides,
	}
}

describe('useOverlayBuilder geometry rounding', () => {
	let builder: ReturnType<typeof useOverlayBuilder>

	beforeEach(() => {
		builder = setupBuilder()
		builder.loadProject({
			id: 'overlay-1',
			name: 'Test overlay',
			width: 1920,
			height: 1080,
			instaSave: false,
			layers: [],
		})
	})

	it('rounds single-layer center alignment to integers', () => {
		builder.project.layers.push(makeLayer('a', { width: 401, height: 201 }))
		builder.selectLayers(['a'])

		builder.alignLayers('center')
		builder.alignLayers('middle')

		const layer = builder.project.layers[0]
		expect(Number.isInteger(layer.posX)).toBe(true)
		expect(Number.isInteger(layer.posY)).toBe(true)
		expect(layer.posX).toBe(Math.round((1920 - 401) / 2))
		expect(layer.posY).toBe(Math.round((1080 - 201) / 2))
	})

	it('rounds horizontal distribution spacing to integers', () => {
		builder.project.layers.push(
			makeLayer('a', { posX: 0, width: 100 }),
			makeLayer('b', { posX: 50, width: 100 }),
			makeLayer('c', { posX: 400, width: 100 }),
		)
		builder.selectLayers(['a', 'b', 'c'])

		builder.distributeLayersHorizontally()

		const middle = builder.project.layers.find((l) => l.id === 'b')
		expect(middle && Number.isInteger(middle.posX)).toBe(true)
	})

	it('rounds guide snapping to canvas center for odd-sized layers', () => {
		builder.canvasState.showGuides = true
		const layer = makeLayer('a', { width: 401, height: 201 })
		// near canvas center: center is (960 - 200.5, 540 - 100.5)
		const snapped = builder.snapToGuides(layer, 760, 440)
		expect(Number.isInteger(snapped.x)).toBe(true)
		expect(Number.isInteger(snapped.y)).toBe(true)
	})
})
