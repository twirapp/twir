import * as z from 'zod'

import { type McpConsentAttempt, mcpConsentAttemptSchema } from '~/utils/mcp-consent.js'
import {
	mcpApprovedScopesSchema,
	mcpRequestedScopesSchema,
} from '~/utils/mcp-scopes.js'

const mcpConsentSchema = z.object({
	client: z.object({
		id: z.string().min(1),
		name: z.string().min(1),
		uri: z.string().min(1).optional(),
	}),
	channel_id: z.string().min(1),
	requested_scopes: mcpRequestedScopesSchema,
	csrf_token: z.string().min(1),
})

const mcpConsentResponseSchema = z.object({
	redirect_to: z.string().min(1),
})

const mcpConsentDecisionSchema = z.discriminatedUnion('decision', [
	z.object({
		attempt: mcpConsentAttemptSchema,
		csrf_token: z.string().min(1),
		channel_id: z.string().min(1),
		decision: z.literal('approve'),
		approved_scopes: mcpApprovedScopesSchema,
	}),
	z.object({
		attempt: mcpConsentAttemptSchema,
		csrf_token: z.string().min(1),
		channel_id: z.string().min(1),
		decision: z.literal('deny'),
	}),
])

const fetchErrorSchema = z.object({
	status: z.number().int().optional(),
})

export type McpConsent = z.infer<typeof mcpConsentSchema>
export type McpConsentDecision = z.infer<typeof mcpConsentDecisionSchema>

export type McpConsentRequestResult<T> =
	| { readonly kind: 'success'; readonly data: T }
	| { readonly kind: 'expired' }
	| { readonly kind: 'permission' }
	| { readonly kind: 'network' }

export type McpConsentFetcher = <T>(
	request: string,
	options: {
		readonly query?: { readonly attempt: McpConsentAttempt }
		readonly method?: 'POST'
		readonly body?: McpConsentDecision
	},
) => Promise<T>

function requestFailure(error: Error): McpConsentRequestResult<never> {
	const result = fetchErrorSchema.safeParse(error)

	switch (result.success ? result.data.status : undefined) {
		case 401:
		case 403:
			return { kind: 'permission' }
		case 404:
		case 409:
		case 410:
			return { kind: 'expired' }
		default:
			return { kind: 'network' }
	}
}

export function createMcpConsentApi(fetcher: McpConsentFetcher) {
	async function getMcpConsent(
		attempt: McpConsentAttempt,
	): Promise<McpConsentRequestResult<McpConsent>> {
		try {
			const response = await fetcher<unknown>('/api/oauth/consent', {
				query: { attempt },
			})

			return { kind: 'success', data: mcpConsentSchema.parse(response) }
		} catch (error) {
			return requestFailure(
				error instanceof Error ? error : new Error('MCP consent request failed'),
			)
		}
	}

	async function submitMcpConsent(
		decision: McpConsentDecision,
	): Promise<McpConsentRequestResult<{ readonly redirectTo: string }>> {
		try {
			const response = await fetcher<unknown>('/api/oauth/consent', {
				method: 'POST',
				body: mcpConsentDecisionSchema.parse(decision),
			})
			const parsed = mcpConsentResponseSchema.parse(response)

			return { kind: 'success', data: { redirectTo: parsed.redirect_to } }
		} catch (error) {
			return requestFailure(
				error instanceof Error ? error : new Error('MCP consent request failed'),
			)
		}
	}

	return { getMcpConsent, submitMcpConsent }
}
