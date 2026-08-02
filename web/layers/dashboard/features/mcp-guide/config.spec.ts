import { describe, expect, it } from 'vitest'

import { createMcpClientGuides } from './config.js'

describe('createMcpClientGuides', () => {
	it('creates ready-to-use configurations for every supported client', () => {
		const endpoint = 'https://twir.example/api/mcp'
		const apiKey = 'channel-api-key'
		const guides = createMcpClientGuides(endpoint, apiKey)

		expect(guides.map((guide) => guide.id)).toEqual(['claude', 'pi', 'codex', 'opencode'])
		for (const guide of guides) {
			expect(guide.config).toContain(endpoint)
			expect(guide.config).toContain(apiKey)
		}
		expect(guides.find((guide) => guide.id === 'opencode')?.config).toContain('"oauth": false')
		expect(guides.find((guide) => guide.id === 'codex')?.config).toContain('http_headers')
	})
})
