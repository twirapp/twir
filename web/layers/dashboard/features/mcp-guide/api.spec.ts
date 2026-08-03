import { describe, expect, it, vi } from 'vitest'

import { createMcpGuideApi } from './api.js'

const catalogFixture = {
	scopes: [
		{
			group: 'commands',
			name: 'Commands',
			description: 'Manage channel commands',
			actions: ['read', 'edit'],
		},
		{
			group: 'timers',
			name: 'Timers',
			description: 'Manage channel timers',
			actions: ['read'],
		},
	],
} as const

describe('MCP guide API', () => {
	it('requests the backend scope catalog and parses each group', async () => {
		const fetchMock = vi.fn().mockResolvedValue(catalogFixture)
		const api = createMcpGuideApi(fetchMock)

		const result = await api.getScopesCatalog()

		expect(fetchMock).toHaveBeenCalledWith('/api/oauth/scopes')
		expect(result).toEqual({ kind: 'success', scopes: catalogFixture.scopes })
	})

	it('reports an error when the catalog request fails', async () => {
		const api = createMcpGuideApi(vi.fn().mockRejectedValue(new Error('network down')))

		expect(await api.getScopesCatalog()).toEqual({ kind: 'error' })
	})

	it('rejects a malformed catalog instead of rendering it', async () => {
		const api = createMcpGuideApi(
			vi.fn().mockResolvedValue({ scopes: [{ group: 'commands' }] }),
		)

		expect(await api.getScopesCatalog()).toEqual({ kind: 'error' })
	})
})
