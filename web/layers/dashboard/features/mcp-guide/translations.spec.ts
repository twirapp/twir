import { describe, expect, it } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import ru from '../../../../i18n/locales/ru.json'

describe('MCP guide translations', () => {
	it('uses locale files loaded by Nuxt i18n', () => {
		expect(en.mcpGuide.title.loc.source).toBe('AI access')
		expect(ru.mcpGuide.title.loc.source).toBe('Доступ для ИИ')
	})
})
