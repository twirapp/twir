import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { createImportsTestI18n } from '../components/test-i18n.js'
import NightbotCallback from './nightbot-callback.vue'
import NightbotImport from './nightbot-import.vue'

const api = vi.hoisted(() => ({
	importCommands: vi.fn(),
	importTimers: vi.fn(),
	logout: vi.fn(),
	postCode: vi.fn(),
	broadcastRefresh: vi.fn(),
}))
const page = vi.hoisted(() => ({
	nightbotData: { value: { userName: 'caster', avatar: '' } as { userName: string; avatar: string } | null },
	nightbotAuthLink: { value: 'https://nightbot.example/oauth' },
	fetching: { value: false },
}))
const route = vi.hoisted(() => ({ query: {} as Record<string, string> }))
const toast = vi.hoisted(() => ({ error: vi.fn() }))

vi.mock('./composables/use-nightbot-integration.js', () => ({
	useNightbotIntegration: () => ({
		importCommands: { executeMutation: api.importCommands },
		importTimers: { executeMutation: api.importTimers },
		logout: { executeMutation: api.logout },
		postCode: { executeMutation: api.postCode },
		broadcastRefresh: api.broadcastRefresh,
	}),
}))
vi.mock('~~/layers/dashboard/api/integrations/integrations-page.js', () => ({
	useIntegrationsPageData: () => page,
}))
vi.mock('~~/layers/dashboard/api/auth.js', () => ({ useUserAccessFlagChecker: () => ({ value: true }) }))
vi.mock('vue-router', () => ({ useRoute: () => route }))
vi.mock('vue-sonner', () => ({ toast }))

const i18n = createImportsTestI18n()
const stubs = {
	ImportProviderCard: {
		props: ['account', 'authLink'],
		emits: ['logout'],
		template: '<section><span>{{ account?.userName }}</span><span>{{ authLink }}</span><button data-test="logout" @click="$emit(\'logout\')">logout</button><slot name="settings" /></section>',
	},
	ImportSettings: {
		props: ['commandsReport', 'timersReport'],
		emits: ['importCommands', 'importTimers'],
		template: '<div><button data-test="commands" @click="$emit(\'importCommands\')">commands</button><button data-test="timers" @click="$emit(\'importTimers\')">timers</button><span>{{ commandsReport?.importedCount }}</span><span>{{ timersReport?.importedCount }}</span></div>',
	},
	NuxtIcon: true,
}

describe('Nightbot import', () => {
	beforeEach(() => {
		vi.clearAllMocks()
		route.query = {}
		page.nightbotData.value = { userName: 'caster', avatar: '' }
		api.logout.mockResolvedValue({ data: { nightbotLogout: true } })
		api.importCommands.mockResolvedValue({ data: { nightbotImportCommands: { importedCount: 2, failedCount: 0, failures: [] } } })
		api.importTimers.mockResolvedValue({ data: { nightbotImportTimers: { importedCount: 3, failedCount: 0, failures: [] } } })
	})

	it('renders the disconnected state from the unified query', () => {
		page.nightbotData.value = null
		const wrapper = mount(NightbotImport, { global: { plugins: [i18n], stubs } })
		expect(wrapper.text()).not.toContain('caster')
		expect(wrapper.text()).toContain('https://nightbot.example/oauth')
	})

	it('uses unified provider data and keeps command and timer imports separate', async () => {
		const wrapper = mount(NightbotImport, { global: { plugins: [i18n], stubs } })
		expect(wrapper.text()).toContain('caster')
		expect(wrapper.text()).toContain('https://nightbot.example/oauth')

		await wrapper.get('[data-test="commands"]').trigger('click')
		await wrapper.get('[data-test="timers"]').trigger('click')
		await flushPromises()

		expect(api.importCommands).toHaveBeenCalledWith({})
		expect(api.importTimers).toHaveBeenCalledWith({})
		expect(wrapper.text()).toContain('2')
		expect(wrapper.text()).toContain('3')
	})

	it('submits code and state and only closes after success', async () => {
		route.query = { code: 'code', state: 'state' }
		api.postCode.mockResolvedValue({ data: { nightbotPostCode: true } })
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)

		mount(NightbotCallback, { global: { plugins: [i18n], stubs } })
		await flushPromises()

		expect(api.postCode).toHaveBeenCalledWith({ input: { code: 'code', state: 'state' } })
		expect(api.broadcastRefresh).toHaveBeenCalled()
		expect(close).toHaveBeenCalled()
	})

	it('keeps the callback open and visible on GraphQL error', async () => {
		route.query = { code: 'code', state: 'state' }
		api.postCode.mockResolvedValue({ error: new Error('rejected') })
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)

		const wrapper = mount(NightbotCallback, { global: { plugins: [i18n], stubs } })
		await flushPromises()

		expect(wrapper.text()).toContain('Could not connect Nightbot')
		expect(api.broadcastRefresh).not.toHaveBeenCalled()
		expect(close).not.toHaveBeenCalled()
	})
})
