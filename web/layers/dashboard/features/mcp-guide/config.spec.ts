import { describe, expect, it } from 'vitest'

import { createMcpClientGuides, createMcpOAuthGuide } from './config.js'

describe('createMcpClientGuides', () => {
	it('creates URL-only OAuth configurations for every supported client', () => {
		const endpoint = 'https://twir.example/api/mcp'
		const [claude, codex, opencode, cursor] = createMcpClientGuides(endpoint)

		expect(claude).toMatchObject({
			id: 'claude',
			config: `claude mcp add --transport http --scope user twir ${JSON.stringify(endpoint)}`,
			authCommand: '/mcp',
		})
		expect(codex).toMatchObject({
			id: 'codex',
			config: `[mcp_servers.twir]\nurl = ${JSON.stringify(endpoint)}`,
			authCommand: 'codex mcp login twir',
		})
		expect(opencode.config).toBe(
			`opencode mcp add twir --url ${endpoint}\nopencode mcp auth twir`,
		)
		expect(opencode.authCommand).toBe('opencode mcp auth twir')
		expect(opencode.fileName).toBeUndefined()
		expect(JSON.parse(cursor.config)).toEqual({
			mcpServers: {
				twir: { url: endpoint },
			},
		})
		expect(cursor.authCommand).toBeUndefined()
	})
})

describe('createMcpOAuthGuide', () => {
	it('documents the complete flow and endpoints for the site origin', () => {
		const guide = createMcpOAuthGuide('https://twir.example')

		expect(guide.steps).toHaveLength(6)
		expect(guide.steps.map((step) => step.titleKey)).toEqual([
			'mcpGuide.oauth.steps.discovery.title',
			'mcpGuide.oauth.steps.register.title',
			'mcpGuide.oauth.steps.authorize.title',
			'mcpGuide.oauth.steps.token.title',
			'mcpGuide.oauth.steps.refresh.title',
			'mcpGuide.oauth.steps.revoke.title',
		])
		expect(guide.endpoints).toEqual([
			{
				method: 'GET',
				url: 'https://twir.example/api/.well-known/oauth-protected-resource',
				descriptionKey: 'mcpGuide.oauth.endpoints.resourceMetadata',
			},
			{
				method: 'GET',
				url: 'https://twir.example/api/.well-known/oauth-authorization-server',
				descriptionKey: 'mcpGuide.oauth.endpoints.serverMetadata',
			},
			{
				method: 'POST',
				url: 'https://twir.example/api/oauth/register',
				descriptionKey: 'mcpGuide.oauth.endpoints.register',
			},
			{
				method: 'GET',
				url: 'https://twir.example/api/oauth/authorize',
				descriptionKey: 'mcpGuide.oauth.endpoints.authorize',
			},
			{
				method: 'POST',
				url: 'https://twir.example/api/oauth/token',
				descriptionKey: 'mcpGuide.oauth.endpoints.token',
			},
			{
				method: 'POST',
				url: 'https://twir.example/api/oauth/revoke',
				descriptionKey: 'mcpGuide.oauth.endpoints.revoke',
			},
		])
		expect(guide).not.toHaveProperty('scopes')
	})
})
