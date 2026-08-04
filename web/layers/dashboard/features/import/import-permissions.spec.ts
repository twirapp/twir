import { print } from 'graphql'
import { describe, expect, it, vi } from 'vitest'
import { IntegrationsPageQuery } from '~~/layers/dashboard/api/integrations/integrations-page.js'

vi.mock('~~/layers/dashboard/api/auth.js', () => ({
	useUserAccessFlagChecker: () => ({ value: false }),
}))

describe('imports page permissions', () => {
	it('loads provider data while every manage-only root field is skipped', () => {
		const source = print(IntegrationsPageQuery)

		expect(source).toContain('query IntegrationsPageData($canManageIntegrations: Boolean!)')
		expect(source).toMatch(/nightbotGetData\s*\{/)
		expect(source).toMatch(/streamelementsGetData\s*\{/)
		const manageIntegrationsFields = [
			'discordIntegrationAuthLink',
			'valorantAuthLink',
			'lastfmAuthLink',
			'donationAlertsAuthLink',
			'donatello',
			'integrationsDonateStream',
			'vkAuthLink',
			'faceitAuthLink',
			'streamlabsAuthLink',
			'nightbotGetAuthLink',
			'streamelementsGetAuthorizationUrl',
		]

		for (const field of manageIntegrationsFields) {
			expect(source).toContain(`${field} @include(if: $canManageIntegrations)`)
		}
	})
})
