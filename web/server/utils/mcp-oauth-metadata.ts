export const MCP_OAUTH_METADATA_RESPONSE_HEADERS = {
	'Access-Control-Allow-Origin': '*',
	'Cache-Control': 'public, max-age=3600',
} as const

const MCP_RESOURCE_PATH = '/api/mcp'
const MCP_SCOPES = ['read', 'write'] as const

export function createMcpProtectedResourceMetadata(origin: string) {
	const canonicalOrigin = new URL(origin).origin

	return {
		resource: `${canonicalOrigin}${MCP_RESOURCE_PATH}`,
		authorization_servers: [canonicalOrigin],
		scopes_supported: MCP_SCOPES,
		bearer_methods_supported: ['header'],
	}
}

export function createMcpAuthorizationServerMetadata(origin: string) {
	const canonicalOrigin = new URL(origin).origin

	return {
		issuer: canonicalOrigin,
		authorization_endpoint: `${canonicalOrigin}/api/oauth/authorize`,
		token_endpoint: `${canonicalOrigin}/api/oauth/token`,
		registration_endpoint: `${canonicalOrigin}/api/oauth/register`,
		revocation_endpoint: `${canonicalOrigin}/api/oauth/revoke`,
		scopes_supported: MCP_SCOPES,
		response_types_supported: ['code'],
		grant_types_supported: ['authorization_code', 'refresh_token'],
		token_endpoint_auth_methods_supported: ['none'],
		code_challenge_methods_supported: ['S256'],
	}
}

export function isMcpProtectedResourceSuffix(resource: string | undefined): boolean {
	return resource === 'api/mcp'
}
