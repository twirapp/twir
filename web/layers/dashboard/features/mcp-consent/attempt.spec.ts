import { describe, expect, it } from 'vitest'

import { getMcpConsentAuthorizePath, parseMcpConsentAttempt } from '~/utils/mcp-consent.js'

describe('MCP consent attempt', () => {
	it('preserves only a valid opaque attempt in the fixed authorize route', () => {
		const attempt = parseMcpConsentAttempt('opaque-attempt_123456')

		expect(attempt).not.toBeNull()
		if (attempt === null) throw new Error('expected valid attempt fixture')
		expect(getMcpConsentAuthorizePath(attempt)).toBe(
			'/dashboard/mcp/authorize?attempt=opaque-attempt_123456',
		)
	})

	it('rejects malformed and duplicate attempt values', () => {
		expect(parseMcpConsentAttempt('../outside')).toBeNull()
		expect(parseMcpConsentAttempt(['opaque-attempt-123', 'opaque-attempt-456'])).toBeNull()
		expect(parseMcpConsentAttempt('short')).toBeNull()
	})
})
