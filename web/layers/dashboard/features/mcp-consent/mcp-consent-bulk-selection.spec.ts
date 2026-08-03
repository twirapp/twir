import { flushPromises, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import type { McpConsent } from './api.js'
import McpConsentPage from './mcp-consent.vue'

vi.mock('#build/fetch.mjs', () => ({ $fetch: vi.fn() }))

const consentApi = vi.hoisted(() => ({
	getMcpConsent: vi.fn(),
	submitMcpConsent: vi.fn(),
}))

const route = vi.hoisted(() => ({
	query: { attempt: 'opaque-attempt-123' },
}))

vi.mock('./api.js', () => ({
	createMcpConsentApi: () => consentApi,
}))
vi.mock('#app/composables/router', () => ({ useRoute: () => route }))

const i18n = createI18n<[typeof en], 'en'>({
	legacy: false,
	locale: 'en',
	messages: { en },
})

const consentFixture: McpConsent = {
	client: { id: 'client-1', name: 'Desktop Agent', uri: 'https://agent.example' },
	channel_id: 'channel-1',
	requested_scopes: [
		{
			group: 'commands',
			name: 'Commands',
			description: 'Manage channel commands',
			actions: ['read'],
		},
		{
			group: 'timers',
			name: 'Timers',
			description: 'Manage channel timers',
			actions: ['read', 'edit'],
		},
		{
			group: 'future_tools',
			name: 'Future tools',
			description: 'Manage future tools',
			actions: ['read', 'edit'],
		},
	],
	csrf_token: 'csrf-token',
}

function mountPage() {
	return mount(McpConsentPage, {
		global: {
			plugins: [i18n],
			stubs: {
				NuxtIcon: { props: ['name'], template: '<i :data-icon="name" />' },
			},
		},
	})
}

async function mountReadyPage() {
	consentApi.getMcpConsent.mockResolvedValue({ kind: 'success', data: consentFixture })
	const wrapper = mountPage()
	await flushPromises()
	return wrapper
}

function translate(key: string): string {
	return i18n.global.t(key)
}

function getGroup(wrapper: Awaited<ReturnType<typeof mountReadyPage>>, group: string) {
	return wrapper.get(`[data-test="scope-group"][data-group="${group}"]`)
}

async function clickBulkSelection(wrapper: Awaited<ReturnType<typeof mountReadyPage>>) {
	await wrapper.get('[data-test="bulk-selection-button"]').trigger('click')
	await flushPromises()
}

describe('MCP consent bulk selection', () => {
	beforeEach(() => {
		consentApi.getMcpConsent.mockReset()
		consentApi.submitMcpConsent.mockReset()
	})

	it('renders every backend scope read-only without submitting before approval', async () => {
		const wrapper = await mountReadyPage()
		const futureTools = getGroup(wrapper, 'future_tools')

		expect(futureTools.get('[data-test="read-switch"]').attributes('aria-checked')).toBe('true')
		expect(futureTools.get('[data-test="edit-switch"]').attributes('aria-checked')).toBe('false')
		expect(consentApi.submitMcpConsent).not.toHaveBeenCalled()
	})

	it('labels a partially selected request as Select all', async () => {
		const wrapper = await mountReadyPage()

		expect(wrapper.get('[data-test="bulk-selection-button"]').text()).toBe(
			translate('mcpConsent.scopes.selectAll'),
		)
	})

	it('selects every backend-provided action including a future group', async () => {
		consentApi.submitMcpConsent.mockResolvedValue({ kind: 'network' })
		const wrapper = await mountReadyPage()

		await clickBulkSelection(wrapper)
		await wrapper.get('form').trigger('submit')
		await flushPromises()

		expect(consentApi.submitMcpConsent).toHaveBeenCalledWith({
			attempt: 'opaque-attempt-123',
			csrf_token: 'csrf-token',
			channel_id: 'channel-1',
			decision: 'approve',
			approved_scopes: [
				'commands:read',
				'timers:read',
				'timers:edit',
				'future_tools:read',
				'future_tools:edit',
			],
		})
	})

	it('labels every selected backend action as Deselect all', async () => {
		const wrapper = await mountReadyPage()

		await clickBulkSelection(wrapper)

		expect(wrapper.get('[data-test="bulk-selection-button"]').text()).toBe(
			translate('mcpConsent.scopes.deselectAll'),
		)
	})

	it('clears every backend scope when deselecting all', async () => {
		const wrapper = await mountReadyPage()

		await clickBulkSelection(wrapper)
		await clickBulkSelection(wrapper)

		for (const group of ['commands', 'timers', 'future_tools']) {
			expect(getGroup(wrapper, group).get('[data-test="read-switch"]').attributes('aria-checked')).toBe(
				'false',
			)
		}
		expect(getGroup(wrapper, 'timers').get('[data-test="edit-switch"]').attributes('aria-checked')).toBe(
			'false',
		)
		expect(
			getGroup(wrapper, 'future_tools').get('[data-test="edit-switch"]').attributes('aria-checked'),
		).toBe('false')
	})

	it('blocks approval without a POST after deselecting all', async () => {
		const wrapper = await mountReadyPage()

		await clickBulkSelection(wrapper)
		await clickBulkSelection(wrapper)
		await wrapper.get('form').trigger('submit')
		await flushPromises()

		expect(consentApi.submitMcpConsent).not.toHaveBeenCalled()
		expect(wrapper.get('[data-test="selection-error"]').text()).toBe(
			translate('mcpConsent.scopes.emptySelection'),
		)
	})
})
