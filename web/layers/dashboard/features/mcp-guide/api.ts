import { mcpScopeCatalogResponseSchema, type McpRequestedScope } from '~/utils/mcp-scopes.js'

export type McpGuideScope = McpRequestedScope

export type McpGuideScopesResult =
	| { readonly kind: 'success'; readonly scopes: readonly McpGuideScope[] }
	| { readonly kind: 'error' }

export type McpGuideFetcher = <T>(request: string) => Promise<T>

export function createMcpGuideApi(fetcher: McpGuideFetcher) {
	async function getScopesCatalog(): Promise<McpGuideScopesResult> {
		try {
			const response = await fetcher<unknown>('/api/oauth/scopes')
			const parsed = mcpScopeCatalogResponseSchema.parse(response)

			return { kind: 'success', scopes: parsed.scopes }
		} catch {
			return { kind: 'error' }
		}
	}

	return { getScopesCatalog }
}
