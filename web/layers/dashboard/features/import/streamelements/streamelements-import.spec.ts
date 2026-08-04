import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { createImportsTestI18n } from '../components/test-i18n.js'
import StreamElementsCallback from './streamelements-callback.vue'
import StreamElementsImport from './streamelements-import.vue'

const api = vi.hoisted(() => ({
	importCommands: vi.fn(), importTimers: vi.fn(), logout: vi.fn(), postCode: vi.fn(), broadcastRefresh: vi.fn(),
}))
const page = vi.hoisted(() => ({
	streamelementsData: { value: { userName: 'elements-user', avatar: '' } as { userName: string; avatar: string } | null },
	streamelementsAuthLink: { value: 'https://streamelements.example/oauth' },
	fetching: { value: false },
}))
const route = vi.hoisted(() => ({ query: {} as Record<string, string> }))
const toast = vi.hoisted(() => ({ error: vi.fn() }))

vi.mock('./composables/use-streamelements-integration.js', () => ({
	useStreamElementsIntegration: () => ({
		importCommands: { executeMutation: api.importCommands },
		importTimers: { executeMutation: api.importTimers },
		logout: { executeMutation: api.logout },
		postCode: { executeMutation: api.postCode },
		broadcastRefresh: api.broadcastRefresh,
	}),
}))
vi.mock('~~/layers/dashboard/api/integrations/integrations-page.js', () => ({ useIntegrationsPageData: () => page }))
vi.mock('~~/layers/dashboard/api/auth.js', () => ({ useUserAccessFlagChecker: () => ({ value: true }) }))
vi.mock('vue-router', () => ({ useRoute: () => route }))
vi.mock('vue-sonner', () => ({ toast }))

const i18n = createImportsTestI18n()
const stubs = {
	ImportProviderCard: {
		props: ['account', 'authLink', 'donationDescription'],
		template: '<section><span>{{ account?.userName }}</span><span>{{ authLink }}</span><span>{{ donationDescription }}</span><slot name="settings" /></section>',
	},
	ImportSettings: {
		emits: ['importCommands', 'importTimers'],
		template: '<div><button data-test="commands" @click="$emit(\'importCommands\')">commands</button><button data-test="timers" @click="$emit(\'importTimers\')">timers</button></div>',
	},
	NuxtIcon: true,
}

describe('StreamElements import', () => {
	beforeEach(() => {
		vi.clearAllMocks()
		route.query = {}
		page.streamelementsData.value = { userName: 'elements-user', avatar: '' }
		api.importCommands.mockResolvedValue({ data: { streamelementsImportCommands: { importedCount: 1, failedCount: 0, failures: [] } } })
		api.importTimers.mockResolvedValue({ data: { streamelementsImportTimers: { importedCount: 1, failedCount: 0, failures: [] } } })
	})

	it('renders disconnected StreamElements without donation-enabled copy', () => {
		page.streamelementsData.value = null
		const wrapper = mount(StreamElementsImport, { global: { plugins: [i18n], stubs } })
		expect(wrapper.text()).not.toContain('elements-user')
		expect(wrapper.text()).toContain('https://streamelements.example/oauth')
		expect(wrapper.text()).not.toContain('enables StreamElements donation events')
	})

	it('uses one connection for donations and independent imports', async () => {
		const wrapper = mount(StreamElementsImport, { global: { plugins: [i18n], stubs } })
		expect(wrapper.text()).toContain('elements-user')
		expect(wrapper.text()).toContain('enables StreamElements donation events')
		await wrapper.get('[data-test="commands"]').trigger('click')
		await wrapper.get('[data-test="timers"]').trigger('click')
		expect(api.importCommands).toHaveBeenCalledWith({})
		expect(api.importTimers).toHaveBeenCalledWith({})
	})

	it('submits code and state and closes only after success', async () => {
		route.query = { code: 'code', state: 'state' }
		api.postCode.mockResolvedValue({ data: { streamelementsPostCode: true } })
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)
		mount(StreamElementsCallback, { global: { plugins: [i18n], stubs } })
		await flushPromises()
		expect(api.postCode).toHaveBeenCalledWith({ input: { code: 'code', state: 'state' } })
		expect(api.broadcastRefresh).toHaveBeenCalled()
		expect(close).toHaveBeenCalled()
	})

	it('does not broadcast or close when callback mutation fails', async () => {
		route.query = { code: 'code', state: 'state' }
		api.postCode.mockResolvedValue({ error: new Error('rejected') })
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)
		const wrapper = mount(StreamElementsCallback, { global: { plugins: [i18n], stubs } })
		await flushPromises()
		expect(wrapper.text()).toContain('Could not connect StreamElements')
		expect(api.broadcastRefresh).not.toHaveBeenCalled()
		expect(close).not.toHaveBeenCalled()
	})
})
