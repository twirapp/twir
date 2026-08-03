import * as z from 'zod'

const mcpGuideScopeSchema = z.object({
	group: z.string().min(1),
	name: z.string().min(1),
	description: z.string().min(1),
	actions: z.array(z.string().min(1)).min(1),
})

const mcpGuideScopesResponseSchema = z.object({
	scopes: z.array(mcpGuideScopeSchema),
})

export type McpGuideScope = z.infer<typeof mcpGuideScopeSchema>

export type McpGuideScopesResult =
	| { readonly kind: 'success'; readonly scopes: readonly McpGuideScope[] }
	| { readonly kind: 'error' }

export type McpGuideFetcher = <T>(request: string) => Promise<T>

export function createMcpGuideApi(fetcher: McpGuideFetcher) {
	async function getScopesCatalog(): Promise<McpGuideScopesResult> {
		try {
			const response = await fetcher<unknown>('/api/oauth/scopes')
			const parsed = mcpGuideScopesResponseSchema.parse(response)

			return { kind: 'success', scopes: parsed.scopes }
		} catch {
			return { kind: 'error' }
		}
	}

	return { getScopesCatalog }
}
