import { describe, expect, it } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import ru from '../../../../i18n/locales/ru.json'

describe('MCP guide translations', () => {
	it('keeps the OAuth client set aligned between supported locales', () => {
		expect(Object.keys(en.mcpGuide.clients)).toEqual(['claude', 'codex', 'opencode', 'cursor'])
		expect(Object.keys(ru.mcpGuide.clients)).toEqual(['claude', 'codex', 'opencode', 'cursor'])
		expect(en.mcpGuide.credentials).not.toHaveProperty('apiKey')
		expect(ru.mcpGuide.credentials).not.toHaveProperty('apiKey')
	})
})
