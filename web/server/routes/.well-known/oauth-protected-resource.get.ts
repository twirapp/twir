import {
	createMcpProtectedResourceMetadata,
	MCP_OAUTH_METADATA_RESPONSE_HEADERS,
} from '../../utils/mcp-oauth-metadata.js'

export default defineEventHandler((event) => {
	const { siteUrl } = useRuntimeConfig(event).public

	for (const [header, value] of Object.entries(MCP_OAUTH_METADATA_RESPONSE_HEADERS)) {
		setResponseHeader(event, header, value)
	}

	return createMcpProtectedResourceMetadata(siteUrl)
})
