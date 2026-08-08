import { ref } from 'vue'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { type OverlayProject, createLayerSettings } from '../types'
import { useOverlaySave } from './useOverlaySave'

import { ChannelOverlayLayerType } from '~/gql/graphql.js'

const api = vi.hoisted(() => ({
	useChannelOverlayCreate: vi.fn(),
	useChannelOverlayUpdate: vi.fn(),
	useChannelOverlaysQuery: vi.fn(),
}))

const instantSave = vi.hoisted(() => ({
	useOverlayInstantSave: vi.fn(),
}))

vi.mock('~~/layers/dashboard/api/overlays/custom', () => api)
vi.mock('./useOverlayInstantSave', () => instantSave)
vi.mock('vue-sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }))

function makeProject(): OverlayProject {
	return {
		id: 'overlay-1',
		name: 'Overlay #1',
		width: 1920,
		height: 1080,
		instaSave: true,
		layers: [
			{
				id: 'layer-1',
				type: ChannelOverlayLayerType.Iframe,
				name: 'stream stats',
				posX: 641.5,
				posY: 36.4,
				width: 1277.6,
				height: 225.2,
				rotation: 10.6,
				opacity: 1,
				visible: true,
				locked: false,
				zIndex: 0,
				periodicallyRefetchData: false,
				settings: createLayerSettings(),
			},
		],
	}
}

describe('useOverlaySave', () => {
	const executeUpdate = vi.fn().mockResolvedValue({ data: {}, error: undefined })

	beforeEach(() => {
		vi.clearAllMocks()
		api.useChannelOverlayCreate.mockReturnValue({ executeMutation: vi.fn() })
		api.useChannelOverlayUpdate.mockReturnValue({ executeMutation: executeUpdate })
		api.useChannelOverlaysQuery.mockReturnValue({ executeQuery: vi.fn(), data: ref(undefined) })
		instantSave.useOverlayInstantSave.mockReturnValue({
			sendLayerPositions: vi.fn(),
			isEnabled: ref(false),
		})
	})

	it('rounds fractional layer geometry before sending the update mutation', async () => {
		const { saveOverlay } = useOverlaySave('overlay-1')

		await saveOverlay(makeProject())

		expect(executeUpdate).toHaveBeenCalledOnce()
		const variables = executeUpdate.mock.calls[0]?.[0]
		const layer = variables.input.layers[0]
		expect(layer.posX).toBe(642)
		expect(layer.posY).toBe(36)
		expect(layer.width).toBe(1278)
		expect(layer.height).toBe(225)
		expect(layer.rotation).toBe(11)
	})
})
