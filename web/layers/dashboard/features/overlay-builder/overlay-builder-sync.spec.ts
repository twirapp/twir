import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import { useOverlayBuilder } from './composables/useOverlayBuilder'
import { type Layer, type OverlayProject, createLayerSettings } from './types'

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

function seedProject(builder: ReturnType<typeof useOverlayBuilder>, layers: Layer[]) {
	builder.loadProject({
		id: 'overlay-1',
		name: 'Test overlay',
		width: 1920,
		height: 1080,
		instaSave: false,
		layers,
	})
}

describe('useOverlayBuilder remote sync', () => {
	let builder: ReturnType<typeof useOverlayBuilder>

	beforeEach(() => {
		builder = setupBuilder()
		seedProject(builder, [makeLayer('a'), makeLayer('b')])
	})

	it('applies remote layer add and ignores duplicates by updating instead', () => {
		builder.applyRemoteLayerAdd(makeLayer('c'))
		expect(builder.project.layers.map((l) => l.id)).toEqual(['a', 'b', 'c'])

		builder.applyRemoteLayerAdd(makeLayer('c', { name: 'renamed' }))
		expect(builder.project.layers).toHaveLength(3)
		expect(builder.project.layers.find((l) => l.id === 'c')?.name).toBe('renamed')
	})

	it('removes remote layer and prunes it from selection', () => {
		builder.selectLayers(['a', 'b'])
		builder.applyRemoteLayerRemove('a')
		expect(builder.project.layers.map((l) => l.id)).toEqual(['b'])
		expect(builder.canvasState.selectedLayerIds).toEqual(['b'])
	})

	it('reindexes zIndex after remote add and remove', () => {
		builder.applyRemoteLayerAdd(makeLayer('c'))
		builder.applyRemoteLayerRemove('a')
		expect(builder.project.layers.map((l) => [l.id, l.zIndex])).toEqual([
			['b', 0],
			['c', 1],
		])
	})

	it('applies remote update for untouched layer geometry', () => {
		builder.applyRemoteLayerUpdate(makeLayer('a', { posX: 500, name: 'peer edit' }))
		const layer = builder.project.layers.find((l) => l.id === 'a')
		expect(layer?.posX).toBe(500)
		expect(layer?.name).toBe('peer edit')
	})

	it('keeps local geometry on remote update while layer is being dragged locally', () => {
		builder.updateLayer('a', { posX: 300 })
		builder.applyRemoteLayerUpdate(makeLayer('a', { posX: 900, posY: 900, name: 'peer edit' }))
		const layer = builder.project.layers.find((l) => l.id === 'a')
		expect(layer?.posX).toBe(300)
		expect(layer?.posY).toBe(10)
		expect(layer?.name).toBe('peer edit')
	})

	it('skips remote positions for locally dragged layers and applies to the rest', () => {
		builder.updateLayer('a', { posX: 111 })
		builder.applyRemoteLayerPositions([
			{ id: 'a', posX: 900, posY: 900, rotation: 45, width: 200, height: 200, visible: true, opacity: 1 },
			{ id: 'b', posX: 700, posY: 700, rotation: 0, width: 100, height: 100, visible: false, opacity: 0.5 },
		])
		const a = builder.project.layers.find((l) => l.id === 'a')
		const b = builder.project.layers.find((l) => l.id === 'b')
		expect(a?.posX).toBe(111)
		expect(b?.posX).toBe(700)
		expect(b?.visible).toBe(false)
		expect(b?.opacity).toBe(0.5)
	})

	it('reorders layers by remote id list and appends unknown ids at the end', () => {
		builder.applyRemoteLayersReorder(['b', 'a', 'ghost'])
		expect(builder.project.layers.map((l) => l.id)).toEqual(['b', 'a'])
		expect(builder.project.layers.map((l) => l.zIndex)).toEqual([0, 1])
	})

	it('replaces the whole project and keeps still-existing selection', () => {
		builder.selectLayers(['a'])
		const replacement: OverlayProject = {
			id: 'overlay-1',
			name: 'Peer project',
			width: 1280,
			height: 720,
			instaSave: true,
			layers: [makeLayer('a'), makeLayer('z')],
		}
		builder.applyRemoteProject(replacement)
		expect(builder.project.name).toBe('Peer project')
		expect(builder.project.width).toBe(1280)
		expect(builder.project.instaSave).toBe(true)
		expect(builder.project.layers.map((l) => l.id)).toEqual(['a', 'z'])
		expect(builder.canvasState.selectedLayerIds).toEqual(['a'])
	})

	it('does not clobber remote ops when undoing a local change', async () => {
		builder.updateLayer('a', { posX: 42 })
		builder.applyRemoteLayerAdd(makeLayer('c'))
		builder.updateLayer('b', { posX: 77 })

		builder.undo()
		await flushPromises()

		const ids = builder.project.layers.map((l) => l.id)
		expect(ids).toContain('c')
		expect(builder.project.layers.find((l) => l.id === 'b')?.posX).toBe(10)
	})

	it('returns created layers from duplicate and paste for broadcasting', () => {
		const duplicated = builder.duplicateLayers(['a'])
		expect(duplicated).toHaveLength(1)
		expect(duplicated[0]?.id).not.toBe('a')

		builder.copyToClipboard()
		const pasted = builder.pasteFromClipboard()
		expect(pasted).toHaveLength(1)
		expect(builder.project.layers).toHaveLength(4)
	})
})
