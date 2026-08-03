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
