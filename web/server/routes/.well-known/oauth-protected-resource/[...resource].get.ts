import {
	createMcpProtectedResourceMetadata,
	isMcpProtectedResourceSuffix,
	MCP_OAUTH_METADATA_RESPONSE_HEADERS,
} from '../../../utils/mcp-oauth-metadata.js'

export default defineEventHandler((event) => {
	const resource = getRouterParam(event, 'resource')
	if (!isMcpProtectedResourceSuffix(resource)) {
		throw createError({
			statusCode: 404,
			statusMessage: 'Not Found',
		})
	}

	const { siteUrl } = useRuntimeConfig(event).public

	for (const [header, value] of Object.entries(MCP_OAUTH_METADATA_RESPONSE_HEADERS)) {
		setResponseHeader(event, header, value)
	}

	return createMcpProtectedResourceMetadata(siteUrl)
})
