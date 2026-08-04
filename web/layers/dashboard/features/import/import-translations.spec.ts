import { describe, expect, it } from 'vitest'

import de from '../../../../i18n/locales/de.json'
import en from '../../../../i18n/locales/en.json'
import es from '../../../../i18n/locales/es.json'
import ja from '../../../../i18n/locales/ja.json'
import pt from '../../../../i18n/locales/pt.json'
import ru from '../../../../i18n/locales/ru.json'
import sk from '../../../../i18n/locales/sk.json'
import uk from '../../../../i18n/locales/uk.json'

function keys(value: object, prefix = ''): string[] {
	return Object.entries(value).flatMap(([key, nested]) => {
		const path = prefix ? `${prefix}.${key}` : key
		const isCompiledMessage = nested && typeof nested === 'object' && 'body' in nested && 'type' in nested
		return nested && typeof nested === 'object' && !isCompiledMessage ? keys(nested, path) : [path]
	})
}

describe('import translations', () => {
	it('keeps the imports namespace in parity across every locale', () => {
		const expected = keys(en.imports).sort()
		for (const locale of [de, es, ja, pt, ru, sk, uk]) {
			expect(keys(locale.imports).sort()).toEqual(expected)
		}
	})
})
