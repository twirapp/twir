import { flushPromises, mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import en from '../../../../i18n/locales/en.json'
import McpConsentPage from './mcp-consent.vue'

vi.mock('#build/fetch.mjs', () => ({ $fetch: vi.fn() }))

const consentApi = vi.hoisted(() => ({
	getMcpConsent: vi.fn(),
	submitMcpConsent: vi.fn(),
}))

const route = vi.hoisted(() => {
	const query: { attempt?: string } = { attempt: 'opaque-attempt-123' }
	return { query }
})

vi.mock('./api.js', () => ({
	createMcpConsentApi: () => consentApi,
}))
vi.mock('#app/composables/router', () => ({ useRoute: () => route }))

const i18n = createI18n<[typeof en], 'en'>({
	legacy: false,
	locale: 'en',
	messages: { en },
})

function translate(key: string): string {
	return i18n.global.t(key)
}

const consentFixture = {
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
	],
	csrf_token: 'csrf-token',
} as const

function mountPage() {
	return mount(McpConsentPage, {
		global: {
			plugins: [i18n],
			stubs: {
				NuxtIcon: { props: ['name'], template: '<i :data-icon="name" />' },
				PageLayout: {
					template:
						'<div><slot name="title" /><slot name="title-footer" /><slot name="content" /></div>',
				},
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

function getGroup(wrapper: Awaited<ReturnType<typeof mountReadyPage>>, group: string) {
	return wrapper.get(`[data-test="scope-group"][data-group="${group}"]`)
}

async function clickSwitch(switchElement: { trigger: (event: string) => Promise<void> }) {
	await switchElement.trigger('click')
	await flushPromises()
}

describe('MCP consent page', () => {
	beforeEach(() => {
		consentApi.getMcpConsent.mockReset()
		consentApi.submitMcpConsent.mockReset()
		route.query = { attempt: 'opaque-attempt-123' }
	})

	it('defaults every requested group to read-only even when edit was requested', async () => {
		const wrapper = await mountReadyPage()

		const timers = getGroup(wrapper, 'timers')
		expect(timers.get('[data-test="read-switch"]').attributes('aria-checked')).toBe('true')
		expect(timers.get('[data-test="edit-switch"]').attributes('aria-checked')).toBe('false')
		expect(wrapper.text()).toContain('Commands')
		expect(wrapper.text()).toContain('Manage channel timers')
		expect(wrapper.find('[data-test="edit-warning"]').exists()).toBe(false)
	})

	it('forces read on when edit is enabled and shows the sensitive access warning', async () => {
		const wrapper = await mountReadyPage()
		const timers = getGroup(wrapper, 'timers')

		await clickSwitch(timers.get('[data-test="read-switch"]'))
		expect(timers.get('[data-test="read-switch"]').attributes('aria-checked')).toBe('false')

		await clickSwitch(timers.get('[data-test="edit-switch"]'))
		expect(timers.get('[data-test="read-switch"]').attributes('aria-checked')).toBe('true')
		expect(timers.get('[data-test="edit-switch"]').attributes('aria-checked')).toBe('true')
		expect(wrapper.find('[data-test="edit-warning"]').exists()).toBe(true)
	})

	it('forces edit off when read is disabled while keeping other groups untouched', async () => {
		const wrapper = await mountReadyPage()
		const timers = getGroup(wrapper, 'timers')

		await clickSwitch(timers.get('[data-test="edit-switch"]'))
		expect(timers.get('[data-test="edit-switch"]').attributes('aria-checked')).toBe('true')

		await clickSwitch(timers.get('[data-test="read-switch"]'))
		expect(timers.get('[data-test="read-switch"]').attributes('aria-checked')).toBe('false')
		expect(timers.get('[data-test="edit-switch"]').attributes('aria-checked')).toBe('false')
		expect(wrapper.find('[data-test="edit-warning"]').exists()).toBe(false)

		expect(getGroup(wrapper, 'commands').get('[data-test="read-switch"]').attributes('aria-checked')).toBe('true')
	})

	it('keeps read on when edit is disabled', async () => {
		const wrapper = await mountReadyPage()
		const timers = getGroup(wrapper, 'timers')

		await clickSwitch(timers.get('[data-test="edit-switch"]'))
		await clickSwitch(timers.get('[data-test="edit-switch"]'))

		expect(timers.get('[data-test="read-switch"]').attributes('aria-checked')).toBe('true')
		expect(timers.get('[data-test="edit-switch"]').attributes('aria-checked')).toBe('false')
	})

	it('blocks approve with a visible validation message when every group is deselected', async () => {
		const wrapper = await mountReadyPage()

		await clickSwitch(getGroup(wrapper, 'commands').get('[data-test="read-switch"]'))
		await clickSwitch(getGroup(wrapper, 'timers').get('[data-test="read-switch"]'))

		await wrapper.get('form').trigger('submit')
		await flushPromises()

		expect(consentApi.submitMcpConsent).not.toHaveBeenCalled()
		expect(wrapper.get('[data-test="selection-error"]').text()).toBe(
			translate('mcpConsent.scopes.emptySelection'),
		)
	})

	it('submits only the selected canonical scopes on approve and redirects', async () => {
		const replace = vi.spyOn(window.location, 'replace').mockImplementation(() => {})
		consentApi.submitMcpConsent.mockResolvedValue({
			kind: 'success',
			data: { redirectTo: 'https://agent.example/callback' },
		})
		const wrapper = await mountReadyPage()

		await clickSwitch(getGroup(wrapper, 'commands').get('[data-test="read-switch"]'))
		await clickSwitch(getGroup(wrapper, 'timers').get('[data-test="edit-switch"]'))
		await wrapper.get('form').trigger('submit')
		await flushPromises()

		expect(consentApi.submitMcpConsent).toHaveBeenCalledWith({
			attempt: 'opaque-attempt-123',
			csrf_token: 'csrf-token',
			channel_id: 'channel-1',
			decision: 'approve',
			approved_scopes: ['timers:read', 'timers:edit'],
		})
		expect(replace).toHaveBeenCalledWith('https://agent.example/callback')
	})

	it('ignores a second approve while the first submission is in flight', async () => {
		let resolveSubmission!: (value: unknown) => void
		consentApi.submitMcpConsent.mockReturnValue(
			new Promise((resolve) => {
				resolveSubmission = resolve
			}),
		)
		const wrapper = await mountReadyPage()

		await wrapper.get('form').trigger('submit')
		await wrapper.get('form').trigger('submit')
		expect(consentApi.submitMcpConsent).toHaveBeenCalledTimes(1)

		resolveSubmission({ kind: 'network' })
		await flushPromises()
	})

	it('posts a denial without approved scopes and redirects', async () => {
		const replace = vi.spyOn(window.location, 'replace').mockImplementation(() => {})
		consentApi.submitMcpConsent.mockResolvedValue({
			kind: 'success',
			data: { redirectTo: 'https://agent.example/cancelled' },
		})
		const wrapper = await mountReadyPage()

		const denyButton = wrapper.get('[data-test="deny-button"]')
		await denyButton.trigger('click')
		await flushPromises()

		expect(consentApi.submitMcpConsent).toHaveBeenCalledWith({
			attempt: 'opaque-attempt-123',
			csrf_token: 'csrf-token',
			channel_id: 'channel-1',
			decision: 'deny',
		})
		expect(replace).toHaveBeenCalledWith('https://agent.example/cancelled')
	})

	it('retries after a network failure through the retry button', async () => {
		consentApi.getMcpConsent
			.mockResolvedValueOnce({ kind: 'network' })
			.mockResolvedValueOnce({ kind: 'success', data: consentFixture })
		const wrapper = mountPage()
		await flushPromises()
		expect(wrapper.text()).toContain(translate('mcpConsent.errors.network.title'))

		const retryButton = wrapper.get('button[type="button"]')
		await retryButton.trigger('click')
		await flushPromises()

		expect(consentApi.getMcpConsent).toHaveBeenCalledTimes(2)
		expect(wrapper.text()).toContain('Desktop Agent')
	})

	it('shows the expired state when the attempt is missing', async () => {
		route.query = {}
		const wrapper = mountPage()
		await flushPromises()

		expect(consentApi.getMcpConsent).not.toHaveBeenCalled()
		expect(wrapper.text()).toContain(translate('mcpConsent.errors.expired.title'))
	})
})
