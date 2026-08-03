import { describe, expect, it } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import ru from '../../../../i18n/locales/ru.json'

describe('MCP consent translations', () => {
	it('keeps consent states and scope chrome aligned between supported locales', () => {
		expect(Object.keys(en.mcpConsent)).toEqual(Object.keys(ru.mcpConsent))
		expect(Object.keys(en.mcpConsent.errors)).toEqual(['expired', 'permission', 'network'])
		expect(Object.keys(ru.mcpConsent.errors)).toEqual(['expired', 'permission', 'network'])
		expect(Object.keys(en.mcpConsent.scopes)).toEqual([
			'title',
			'description',
			'read',
			'edit',
			'editIncludesRead',
			'emptySelection',
		])
		expect(Object.keys(ru.mcpConsent.scopes)).toEqual(Object.keys(en.mcpConsent.scopes))
		expect(Object.keys(en.mcpConsent.writeWarning)).toEqual(['title', 'description'])
		expect(Object.keys(ru.mcpConsent.writeWarning)).toEqual(['title', 'description'])
	})
})
