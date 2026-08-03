import * as z from 'zod'

export const MCP_SCOPE_ACTIONS = ['read', 'edit'] as const

export const mcpScopeGroupSchema = z.string().min(1).regex(/^[a-z][a-z0-9_]*$/)
export const mcpScopeActionSchema = z.enum(MCP_SCOPE_ACTIONS)
export const mcpScopeTokenSchema = z.templateLiteral([
	mcpScopeGroupSchema,
	z.literal(':'),
	mcpScopeActionSchema,
])

export type McpScopeGroup = z.infer<typeof mcpScopeGroupSchema>
export type McpScopeAction = z.infer<typeof mcpScopeActionSchema>
export type McpScopeToken = z.infer<typeof mcpScopeTokenSchema>

const scopeActionsSchema = z
	.array(mcpScopeActionSchema)
	.min(1)
	.max(MCP_SCOPE_ACTIONS.length)
	.superRefine((actions, context) => {
		if (new Set(actions).size !== actions.length) {
			context.addIssue({ code: 'custom', message: 'scope actions must be unique' })
		}

		if (actions.includes('edit') && !actions.includes('read')) {
			context.addIssue({ code: 'custom', message: 'edit scope requires read scope' })
		}

		if (actions.length === 2 && (actions[0] !== 'read' || actions[1] !== 'edit')) {
			context.addIssue({ code: 'custom', message: 'scope actions must use canonical order' })
		}
	})

export const mcpRequestedScopeSchema = z.object({
	group: mcpScopeGroupSchema,
	name: z.string().min(1),
	description: z.string().min(1),
	actions: scopeActionsSchema,
})

export type McpRequestedScope = z.infer<typeof mcpRequestedScopeSchema>

export const mcpRequestedScopesSchema = z
	.array(mcpRequestedScopeSchema)
	.min(1)
	.superRefine((scopes, context) => {
		const seenGroups = new Set<McpScopeGroup>()

		for (const [index, scope] of scopes.entries()) {
			if (seenGroups.has(scope.group)) {
				context.addIssue({
					code: 'custom',
					message: 'scope groups must be unique',
					path: [index, 'group'],
				})
			}
			seenGroups.add(scope.group)
		}
	})

export type McpRequestedScopes = z.infer<typeof mcpRequestedScopesSchema>

export const mcpApprovedScopesSchema = z
	.array(mcpScopeTokenSchema)
	.min(1)
	.superRefine((scopes, context) => {
		if (new Set(scopes).size !== scopes.length) {
			context.addIssue({ code: 'custom', message: 'approved scopes must be unique' })
		}
	})

export type McpApprovedScopes = z.infer<typeof mcpApprovedScopesSchema>

export function buildMcpScopeToken(group: McpScopeGroup, action: McpScopeAction): McpScopeToken {
	return mcpScopeTokenSchema.parse(`${group}:${action}`)
}

export function flattenMcpScopes(scopes: readonly McpRequestedScope[]): McpScopeToken[] {
	return scopes.flatMap(({ group, actions }) =>
		actions.map((action) => buildMcpScopeToken(group, action)),
	)
}
