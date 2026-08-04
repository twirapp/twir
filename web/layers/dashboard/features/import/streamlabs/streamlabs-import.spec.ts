import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { createImportsTestI18n } from '../components/test-i18n.js'
import StreamlabsCallback from './streamlabs-callback.vue'
import StreamlabsImport from './streamlabs-import.vue'

const logout = vi.hoisted(() => vi.fn())
const postCode = vi.hoisted(() => vi.fn())
const broadcastRefresh = vi.hoisted(() => vi.fn())
const route = vi.hoisted(() => ({ query: {} as Record<string, string> }))
const page = vi.hoisted(() => ({
	streamlabsData: {
		value: { enabled: true, userName: 'labs-user', avatar: '' } as {
			enabled: boolean
			userName: string | null
			avatar: string | null
		},
	},
	streamlabsAuthLink: { value: 'https://streamlabs.example/oauth' },
	fetching: { value: false },
}))

vi.mock('~~/layers/dashboard/api/integrations/integrations-page.js', () => ({ useIntegrationsPageData: () => page }))
vi.mock('~~/layers/dashboard/api/integrations/integrations.js', () => ({
	useIntegrations: () => ({
		streamlabsLogout: () => ({ executeMutation: logout }),
		streamlabsPostCode: () => ({ executeMutation: postCode }),
		broadcastRefresh,
	}),
}))
vi.mock('~~/layers/dashboard/api/auth.js', () => ({ useUserAccessFlagChecker: () => ({ value: true }) }))
vi.mock('vue-router', () => ({ useRoute: () => route }))
vi.mock('vue-sonner', () => ({ toast: { error: vi.fn() } }))

const i18n = createImportsTestI18n()

describe('Streamlabs import card', () => {
	beforeEach(() => {
		vi.clearAllMocks()
		route.query = {}
		page.streamlabsData.value = { enabled: true, userName: 'labs-user', avatar: '' }
	})

	it('shows donation connection and explicitly keeps import unavailable', () => {
		const wrapper = mount(StreamlabsImport, {
			global: {
				plugins: [i18n],
				stubs: {
					ImportProviderCard: {
						props: ['account', 'connected', 'importAvailable', 'donationDescription', 'unavailableDescription'],
						template: '<section><span>{{ account?.userName }}</span><span>connected={{ connected }}</span><span>{{ importAvailable }}</span><span>{{ donationDescription }}</span><span>{{ unavailableDescription }}</span><slot name="settings" /></section>',
					},
				},
			},
		})

		expect(wrapper.text()).toContain('labs-user')
		expect(wrapper.text()).toContain('connected=true')
		expect(wrapper.text()).toContain('false')
		expect(wrapper.text()).toContain('enables Streamlabs donation events')
		expect(wrapper.text()).toContain('does not expose Cloudbot commands or timers')
		expect(wrapper.find('[data-test="import-commands"]').exists()).toBe(false)
		expect(wrapper.find('[data-test="provider-settings"]').exists()).toBe(false)
	})

	it('keeps OAuth available while disconnected without claiming donations are active', () => {
		page.streamlabsData.value = { enabled: false, userName: 'retained-profile', avatar: 'avatar' }
		const wrapper = mount(StreamlabsImport, {
			global: {
				plugins: [i18n],
				stubs: {
					ImportProviderCard: {
						props: ['account', 'authLink', 'connected', 'donationDescription', 'importAvailable'],
						template: '<section><span>{{ account?.userName }}</span><span>{{ authLink }}</span><span>connected={{ connected }}</span><span>{{ donationDescription }}</span><span>{{ importAvailable }}</span></section>',
					},
				},
			},
		})

		expect(wrapper.text()).toContain('https://streamlabs.example/oauth')
		expect(wrapper.text()).toContain('retained-profile')
		expect(wrapper.text()).toContain('connected=false')
		expect(wrapper.text()).toContain('false')
		expect(wrapper.text()).not.toContain('enables Streamlabs donation events')
	})

	it('submits OAuth state and only closes after a successful callback', async () => {
		route.query = { code: 'code', state: 'state' }
		postCode.mockResolvedValue({ data: { streamlabsPostCode: true } })
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)
		mount(StreamlabsCallback, { global: { plugins: [i18n], stubs: { NuxtIcon: true } } })
		await flushPromises()
		expect(postCode).toHaveBeenCalledWith({ input: { code: 'code', state: 'state' } })
		expect(broadcastRefresh).toHaveBeenCalled()
		expect(close).toHaveBeenCalled()
	})

	it('keeps the Streamlabs callback open when GraphQL returns an error', async () => {
		route.query = { code: 'code', state: 'state' }
		postCode.mockResolvedValue({ error: new Error('rejected') })
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)
		const wrapper = mount(StreamlabsCallback, { global: { plugins: [i18n], stubs: { NuxtIcon: true } } })
		await flushPromises()
		expect(wrapper.text()).toContain('Could not connect Streamlabs')
		expect(broadcastRefresh).not.toHaveBeenCalled()
		expect(close).not.toHaveBeenCalled()
	})
})
