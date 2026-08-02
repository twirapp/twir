import { print } from 'graphql'
import { describe, expect, it, vi } from 'vitest'
import { IntegrationsPageQuery } from '~~/layers/dashboard/api/integrations/integrations-page.js'

vi.mock('~~/layers/dashboard/api/auth.js', () => ({
	useUserAccessFlagChecker: () => ({ value: false }),
}))

describe('imports page permissions', () => {
	it('loads provider data while protected OAuth links are skipped', () => {
		const source = print(IntegrationsPageQuery)

		expect(source).toContain('query IntegrationsPageData($canManageIntegrations: Boolean!)')
		expect(source).toMatch(/nightbotGetData\s*\{/)
		expect(source).toMatch(/streamelementsGetData\s*\{/)
		expect(source).toMatch(
			/nightbotGetAuthLink @include\(if: \$canManageIntegrations\)/
		)
		expect(source).toMatch(
			/streamelementsGetAuthorizationUrl @include\(if: \$canManageIntegrations\)/
		)

		for (const field of [
			'discordIntegrationAuthLink',
			'valorantAuthLink',
			'lastfmAuthLink',
			'donationAlertsAuthLink',
			'vkAuthLink',
			'faceitAuthLink',
			'streamlabsAuthLink',
		]) {
			expect(source).toContain(`${field} @include(if: $canManageIntegrations)`)
		}
	})
})
