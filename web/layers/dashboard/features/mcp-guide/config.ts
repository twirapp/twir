export type McpClientId = 'claude' | 'pi' | 'codex' | 'opencode'

export interface McpClientGuide {
	id: McpClientId
	name: string
	icon: string
	descriptionKey: string
	stepKeys: string[]
	docsUrl: string
	fileName?: string
	config: string
}

export function createMcpClientGuides(endpoint: string, apiKey: string): McpClientGuide[] {
	const server = {
		type: 'http',
		url: endpoint,
		headers: {
			'Api-Key': apiKey,
		},
	}

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
			config: `claude mcp add-json --scope user twir '${JSON.stringify(server)}'`,
		},
		{
			id: 'pi',
			name: 'Pi',
			icon: 'lucide:pi',
			descriptionKey: 'mcpGuide.clients.pi.description',
			stepKeys: [
				'mcpGuide.clients.pi.step1',
				'mcpGuide.clients.pi.step2',
				'mcpGuide.clients.pi.step3',
			],
			docsUrl: 'https://github.com/nicobailon/pi-mcp-adapter',
			fileName: '~/.config/mcp/mcp.json',
			config: JSON.stringify({ mcpServers: { twir: { url: endpoint, headers: server.headers } } }, null, 2),
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
			docsUrl: 'https://learn.chatgpt.com/docs/extend/mcp.md',
			fileName: '~/.codex/config.toml',
			config: `[mcp_servers.twir]\nurl = ${JSON.stringify(endpoint)}\nhttp_headers = { "Api-Key" = ${JSON.stringify(apiKey)} }`,
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
						oauth: false,
						headers: server.headers,
						enabled: true,
					},
				},
			}, null, 2),
		},
	]
}
