import { describe, expect, it } from 'vitest'

import {
	buildMcpScopeToken,
	flattenMcpScopes,
	mcpApprovedScopesSchema,
	mcpRequestedScopesSchema,
	mcpRequestedScopeSchema,
	mcpScopeGroupSchema,
	mcpScopeTokenSchema,
	mcpScopeCatalogResponseSchema,
} from './mcp-scopes.ts'

function groupedScope(group: string, actions: ('read' | 'edit')[] = ['read', 'edit']) {
	return { group, name: `${group} name`, description: `${group} description`, actions }
}

describe('MCP scope model', () => {
	it('parses any valid scope group grammar and rejects malformed groups', () => {
		expect(mcpScopeGroupSchema.safeParse('commands').success).toBe(true)
		expect(mcpScopeGroupSchema.safeParse('future_scope').success).toBe(true)
		expect(mcpScopeGroupSchema.safeParse('').success).toBe(false)
		expect(mcpScopeGroupSchema.safeParse('Commands').success).toBe(false)
		expect(mcpScopeGroupSchema.safeParse('future-scope').success).toBe(false)
	})

	it('parses scope tokens and rejects unknown actions', () => {
		expect(mcpScopeTokenSchema.safeParse('commands:read').success).toBe(true)
		expect(mcpScopeTokenSchema.safeParse('future_scope:edit').success).toBe(true)
		expect(mcpScopeTokenSchema.safeParse('commands:write').success).toBe(false)
		expect(mcpScopeTokenSchema.safeParse('Commands:read').success).toBe(false)
	})

	it('flattens requested scopes in backend payload order', () => {
		const requestedScopes = [
			groupedScope('commands', ['read']),
			groupedScope('future_scope', ['read', 'edit']),
		]

		expect(flattenMcpScopes(requestedScopes)).toEqual([
			'commands:read',
			'future_scope:read',
			'future_scope:edit',
		])
		expect(buildMcpScopeToken('future_scope', 'edit')).toBe('future_scope:edit')
	})

	it('requires unique requested groups and edit to include read', () => {
		expect(
			mcpRequestedScopesSchema.safeParse([
				groupedScope('commands', ['read']),
				groupedScope('commands', ['edit']),
			]).success,
		).toBe(false)
		expect(mcpRequestedScopeSchema.safeParse(groupedScope('commands', ['edit'])).success).toBe(
			false,
		)
		expect(mcpRequestedScopeSchema.safeParse(groupedScope('commands', ['edit', 'read'])).success).toBe(
			false,
		)
		expect(mcpRequestedScopeSchema.safeParse(groupedScope('commands', ['read', 'edit'])).success).toBe(
			true,
		)
	})

	it('accepts backend-shaped catalogs with unknown groups and any count, while rejecting duplicates', () => {
		expect(
			mcpScopeCatalogResponseSchema.safeParse({
				scopes: [
					groupedScope('commands', ['read']),
					groupedScope('future_scope', ['read', 'edit']),
					groupedScope('another_scope', ['read']),
				],
			}).success,
		).toBe(true)
		expect(
			mcpScopeCatalogResponseSchema.safeParse({
				scopes: [groupedScope('commands', ['read']), groupedScope('commands', ['edit'])],
			}).success,
		).toBe(false)
		expect(mcpScopeCatalogResponseSchema.safeParse({ scopes: [] }).success).toBe(true)
	})

	it('rejects malformed approved scope tokens and empty approvals', () => {
		expect(mcpApprovedScopesSchema.safeParse([]).success).toBe(false)
		expect(mcpApprovedScopesSchema.safeParse(['commands:write']).success).toBe(false)
		expect(mcpApprovedScopesSchema.safeParse(['future_scope:read']).success).toBe(true)
	})
})
