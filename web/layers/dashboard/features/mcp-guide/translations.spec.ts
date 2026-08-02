import { describe, expect, it } from 'vitest'

import en from '../../../../locales/en.json'
import ru from '../../../../locales/ru.json'

describe('MCP guide translations', () => {
	it('uses locale files loaded by Nuxt i18n', () => {
		expect(en.mcpGuide.title).toBe('AI access')
		expect(ru.mcpGuide.title).toBe('Доступ для ИИ')
	})
})
