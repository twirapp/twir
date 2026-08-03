export type McpClientId = 'claude' | 'codex' | 'opencode' | 'cursor'

export interface McpClientGuide {
	readonly id: McpClientId
	readonly name: string
	readonly icon: string
	readonly descriptionKey: string
	readonly stepKeys: readonly string[]
	readonly docsUrl: string
	readonly fileName?: string
	readonly config: string
	readonly authCommand?: string
}

export function createMcpClientGuides(endpoint: string): readonly [
	McpClientGuide,
	McpClientGuide,
	McpClientGuide,
	McpClientGuide,
] {
	return [
		{
			id: 'claude',
			name: 'Claude Code',
			icon: 'simple-icons:anthropic',
			descriptionKey: 'mcpGuide.clients.claude.description',
			stepKeys: [
				'mcpGuide.clients.claude.step1',
				'mcpGuide.clients.claude.step2',
			],
			docsUrl: 'https://code.claude.com/docs/en/mcp',
			config: `claude mcp add --transport http --scope user twir ${JSON.stringify(endpoint)}`,
			authCommand: '/mcp',
		},
		{
			id: 'codex',
			name: 'Codex',
			icon: 'simple-icons:openai',
			descriptionKey: 'mcpGuide.clients.codex.description',
			stepKeys: [
				'mcpGuide.clients.codex.step1',
				'mcpGuide.clients.codex.step2',
			],
			docsUrl: 'https://developers.openai.com/codex/mcp',
			fileName: '~/.codex/config.toml',
			config: `[mcp_servers.twir]\nurl = ${JSON.stringify(endpoint)}`,
			authCommand: 'codex mcp login twir',
		},
		{
			id: 'opencode',
			name: 'OpenCode',
			icon: 'lucide:terminal',
			descriptionKey: 'mcpGuide.clients.opencode.description',
			stepKeys: [
				'mcpGuide.clients.opencode.step1',
				'mcpGuide.clients.opencode.step2',
			],
			docsUrl: 'https://opencode.ai/docs/mcp-servers/',
			fileName: 'opencode.json',
			config: JSON.stringify({
				mcp: {
					twir: {
						type: 'remote',
						url: endpoint,
						oauth: {},
						enabled: true,
					},
				},
			}, null, 2),
			authCommand: 'opencode mcp auth twir',
		},
		{
			id: 'cursor',
			name: 'Cursor',
			icon: 'lucide:mouse-pointer-2',
			descriptionKey: 'mcpGuide.clients.cursor.description',
			stepKeys: [
				'mcpGuide.clients.cursor.step1',
				'mcpGuide.clients.cursor.step2',
			],
			docsUrl: 'https://cursor.com/docs/mcp',
			fileName: '~/.cursor/mcp.json',
			config: JSON.stringify({ mcpServers: { twir: { url: endpoint } } }, null, 2),
		},
	]
}

export interface McpOAuthFlowStep {
	readonly titleKey: string
	readonly descriptionKey: string
}

export interface McpOAuthEndpoint {
	readonly method: 'GET' | 'POST'
	readonly url: string
	readonly descriptionKey: string
}

export interface McpOAuthGuide {
	readonly steps: readonly McpOAuthFlowStep[]
	readonly endpoints: readonly McpOAuthEndpoint[]
}

export function createMcpOAuthGuide(origin: string): McpOAuthGuide {
	const api = `${origin}/api`

	return {
		steps: [
			{
				titleKey: 'mcpGuide.oauth.steps.discovery.title',
				descriptionKey: 'mcpGuide.oauth.steps.discovery.description',
			},
			{
				titleKey: 'mcpGuide.oauth.steps.register.title',
				descriptionKey: 'mcpGuide.oauth.steps.register.description',
			},
			{
				titleKey: 'mcpGuide.oauth.steps.authorize.title',
				descriptionKey: 'mcpGuide.oauth.steps.authorize.description',
			},
			{
				titleKey: 'mcpGuide.oauth.steps.token.title',
				descriptionKey: 'mcpGuide.oauth.steps.token.description',
			},
			{
				titleKey: 'mcpGuide.oauth.steps.refresh.title',
				descriptionKey: 'mcpGuide.oauth.steps.refresh.description',
			},
			{
				titleKey: 'mcpGuide.oauth.steps.revoke.title',
				descriptionKey: 'mcpGuide.oauth.steps.revoke.description',
			},
		],
		endpoints: [
			{
				method: 'GET',
				url: `${api}/.well-known/oauth-protected-resource`,
				descriptionKey: 'mcpGuide.oauth.endpoints.resourceMetadata',
			},
			{
				method: 'GET',
				url: `${api}/.well-known/oauth-authorization-server`,
				descriptionKey: 'mcpGuide.oauth.endpoints.serverMetadata',
			},
			{
				method: 'POST',
				url: `${api}/oauth/register`,
				descriptionKey: 'mcpGuide.oauth.endpoints.register',
			},
			{
				method: 'GET',
				url: `${api}/oauth/authorize`,
				descriptionKey: 'mcpGuide.oauth.endpoints.authorize',
			},
			{
				method: 'POST',
				url: `${api}/oauth/token`,
				descriptionKey: 'mcpGuide.oauth.endpoints.token',
			},
			{
				method: 'POST',
				url: `${api}/oauth/revoke`,
				descriptionKey: 'mcpGuide.oauth.endpoints.revoke',
			},
		],
	}
}
