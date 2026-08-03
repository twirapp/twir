import { describe, expect, it } from 'vitest'

import {
	MCP_OAUTH_METADATA_RESPONSE_HEADERS,
	createMcpAuthorizationServerMetadata,
	createMcpProtectedResourceMetadata,
	isMcpProtectedResourceSuffix,
} from './mcp-oauth-metadata.js'

describe('MCP OAuth metadata', () => {
	it('creates canonical discovery documents for the MCP resource', () => {
		const givenOrigin = 'https://twir.example/dashboard?source=metadata'

		const protectedResource = createMcpProtectedResourceMetadata(givenOrigin)
		const authorizationServer = createMcpAuthorizationServerMetadata(givenOrigin)

		expect(protectedResource).toEqual({
			resource: 'https://twir.example/api/mcp',
			authorization_servers: ['https://twir.example'],
			scopes_supported: ['read', 'write'],
			bearer_methods_supported: ['header'],
		})
		expect(authorizationServer).toEqual({
			issuer: 'https://twir.example',
			authorization_endpoint: 'https://twir.example/api/oauth/authorize',
			token_endpoint: 'https://twir.example/api/oauth/token',
			registration_endpoint: 'https://twir.example/api/oauth/register',
			revocation_endpoint: 'https://twir.example/api/oauth/revoke',
			scopes_supported: ['read', 'write'],
			response_types_supported: ['code'],
			grant_types_supported: ['authorization_code', 'refresh_token'],
			token_endpoint_auth_methods_supported: ['none'],
			code_challenge_methods_supported: ['S256'],
		})
		expect(MCP_OAUTH_METADATA_RESPONSE_HEADERS).toEqual({
			'Access-Control-Allow-Origin': '*',
			'Cache-Control': 'public, max-age=3600',
		})
	})

	it('accepts only the MCP protected-resource suffix', () => {
		expect(isMcpProtectedResourceSuffix('api/mcp')).toBe(true)
		expect(isMcpProtectedResourceSuffix('api/other')).toBe(false)
		expect(isMcpProtectedResourceSuffix(undefined)).toBe(false)
	})
})
