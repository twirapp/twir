import { describe, expect, it } from 'vitest'

import { createMcpClientGuides } from './config.js'

describe('createMcpClientGuides', () => {
	it('creates URL-only OAuth configurations for every supported client', () => {
		const endpoint = 'https://twir.example/api/mcp'
		const [claude, codex, opencode, cursor] = createMcpClientGuides(endpoint)

		expect(claude).toMatchObject({
			id: 'claude',
			config: `claude mcp add --transport http --scope user twir ${JSON.stringify(endpoint)}`,
		})
		expect(codex).toMatchObject({
			id: 'codex',
			config: `[mcp_servers.twir]\nurl = ${JSON.stringify(endpoint)}`,
		})
		expect(JSON.parse(opencode.config)).toEqual({
			mcp: {
				twir: {
					type: 'remote',
					url: endpoint,
					oauth: {},
					enabled: true,
				},
			},
		})
		expect(JSON.parse(cursor.config)).toEqual({
			mcpServers: {
				twir: { url: endpoint },
			},
		})
	})
})
