import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { describe, expect, it } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import type { McpGuideScope } from './api.js'
import McpGuideScopes from './mcp-guide-scopes.vue'

const i18n = createI18n<[typeof en], 'en'>({
	legacy: false,
	locale: 'en',
	messages: { en },
})

function translate(key: string): string {
	return i18n.global.t(key)
}

const scopesFixture: readonly McpGuideScope[] = [
	{
		group: 'commands',
		name: 'Commands',
		description: 'Manage channel commands',
		actions: ['read', 'edit'],
	},
	{
		group: 'timers',
		name: 'Timers',
		description: 'Manage channel timers',
		actions: ['read'],
	},
]

function mountScopes(state: 'loading' | 'ready' | 'error', scopes: readonly McpGuideScope[] = []) {
	return mount(McpGuideScopes, {
		props: { state, scopes },
		global: { plugins: [i18n] },
	})
}

describe('MCP guide scopes', () => {
	it('renders each catalog group with its name, description, and action badges', () => {
		const wrapper = mountScopes('ready', scopesFixture)

		const commands = wrapper.get('[data-test="scope-group"][data-group="commands"]')
		expect(commands.text()).toContain('Commands')
		expect(commands.text()).toContain('Manage channel commands')
		expect(commands.findAll('[data-test="scope-action"]').map((badge) => badge.text())).toEqual([
			'read',
			'edit',
		])

		const timers = wrapper.get('[data-test="scope-group"][data-group="timers"]')
		expect(timers.findAll('[data-test="scope-action"]')).toHaveLength(1)
	})

	it('always shows the static notes about edit access and legacy aliases', () => {
		const wrapper = mountScopes('ready', scopesFixture)

		expect(wrapper.text()).toContain(translate('mcpGuide.oauth.scopesEditIncludesRead'))
		expect(wrapper.text()).toContain(translate('mcpGuide.oauth.scopesLegacyAliases'))
		expect(wrapper.text()).toContain(translate('mcpGuide.oauth.tokenLifetimes'))
	})

	it('shows a loading note while the catalog is being fetched', () => {
		const wrapper = mountScopes('loading')

		expect(wrapper.get('[data-test="scopes-loading"]').text()).toBe(
			translate('mcpGuide.oauth.scopesLoading'),
		)
		expect(wrapper.find('[data-test="scope-group"]').exists()).toBe(false)
	})

	it('shows a graceful fallback note instead of the list when the catalog fails to load', () => {
		const wrapper = mountScopes('error')

		expect(wrapper.get('[data-test="scopes-unavailable"]').text()).toBe(
			translate('mcpGuide.oauth.scopesUnavailable'),
		)
		expect(wrapper.find('[data-test="scope-group"]').exists()).toBe(false)
		expect(wrapper.text()).toContain(translate('mcpGuide.oauth.scopesEditIncludesRead'))
	})
})
