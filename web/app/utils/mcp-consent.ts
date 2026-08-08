import * as z from 'zod'

export const mcpConsentAttemptSchema = z
	.string()
	.min(16)
	.max(512)
	.regex(/^[A-Za-z0-9_-]+$/)
	.brand<'McpConsentAttempt'>()

export type McpConsentAttempt = z.infer<typeof mcpConsentAttemptSchema>

export function parseMcpConsentAttempt(value: unknown): McpConsentAttempt | null {
	const result = mcpConsentAttemptSchema.safeParse(value)
	return result.success ? result.data : null
}

export function getMcpConsentAuthorizePath(attempt: McpConsentAttempt): string {
	return `/dashboard/mcp/authorize?${new URLSearchParams({ attempt }).toString()}`
}
