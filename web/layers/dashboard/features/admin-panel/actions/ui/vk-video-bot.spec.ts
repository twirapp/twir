import { flushPromises, mount } from '@vue/test-utils'
import { nextTick, ref } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

const adminActions = vi.hoisted(() => ({
	vkVideoBotSetupBroadcastChannelName: 'vk_video_bot_setup',
	useMutationVKVideoBotSetupComplete: vi.fn(),
	useMutationVKVideoBotSetupLink: vi.fn(),
	useMutationVKVideoBotSetupStatus: vi.fn(),
}))

const dashboard = vi.hoisted(() => ({
	useBotJoinPart: vi.fn(),
	useBotStatuses: vi.fn(),
}))

const sonner = vi.hoisted(() => ({
	toast: {
		error: vi.fn(),
		success: vi.fn(),
	},
}))

const route = vi.hoisted(() => ({
	query: {},
}))

vi.mock('~~/layers/dashboard/api/admin/actions', () => adminActions)
vi.mock('~~/layers/dashboard/api/dashboard', () => dashboard)
vi.mock('vue-sonner', () => sonner)
vi.mock('vue-router', () => ({ useRoute: () => route }))

import HeaderBotStatus from '~~/layers/dashboard/layout/header/header-bot-status.vue'
import SetupVKVideoBot from './setup-vk-video-bot.vue'
import VKVideoBotCallback from './vk-video-bot-callback.vue'

const componentStubs = {
	ActionConfirm: {
		props: ['open'],
		emits: ['confirm'],
		template: '<button v-if="open" data-test="replace-confirm" @click="$emit(\'confirm\')">Confirm</button>',
	},
	Button: {
		props: ['disabled'],
		template: '<button :disabled="disabled"><slot /></button>',
	},
	Card: { template: '<section><slot /></section>' },
	CardContent: { template: '<div><slot /></div>' },
	CardDescription: { template: '<p><slot /></p>' },
	CardHeader: { template: '<header><slot /></header>' },
	CardTitle: { template: '<h2><slot /></h2>' },
	NuxtIcon: {
		props: ['name'],
		template: '<i :data-icon="name" />',
	},
	Alert: { template: '<div><slot /></div>' },
	AlertDescription: { template: '<p><slot /></p>' },
}

describe('VK Video Live bot setup', () => {
	beforeEach(() => {
		adminActions.useMutationVKVideoBotSetupComplete.mockReset()
		adminActions.useMutationVKVideoBotSetupLink.mockReset()
		adminActions.useMutationVKVideoBotSetupStatus.mockReset()
		dashboard.useBotJoinPart.mockReset()
		dashboard.useBotStatuses.mockReset()
		sonner.toast.error.mockReset()
		sonner.toast.success.mockReset()
		route.query = {}
		vi.stubGlobal('BroadcastChannel', class {
			close = vi.fn()
			postMessage = vi.fn()
		})
	})

	afterEach(() => {
		vi.restoreAllMocks()
		vi.unstubAllGlobals()
	})

	it('opens the VK authorization popup when the global bot is disconnected', async () => {
		const status = { executeMutation: vi.fn().mockResolvedValue({ data: { vkVideoBotSetupStatus: false } }) }
		const setup = { executeMutation: vi.fn().mockResolvedValue({ data: { vkVideoBotSetupLink: 'https://vk.example/authorize' } }) }
		adminActions.useMutationVKVideoBotSetupStatus.mockReturnValue(status)
		adminActions.useMutationVKVideoBotSetupLink.mockReturnValue(setup)
		const open = vi.spyOn(window, 'open').mockReturnValue(null)

		const wrapper = mount(SetupVKVideoBot, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(wrapper.text()).toContain('Not connected')
		await wrapper.get('button').trigger('click')

		expect(setup.executeMutation).toHaveBeenCalledWith({})
		expect(open).toHaveBeenCalledWith('https://vk.example/authorize', 'vk-video-bot-setup', 'popup')
	})

	it('requires confirmation before replacing the connected global bot', async () => {
		const status = { executeMutation: vi.fn().mockResolvedValue({ data: { vkVideoBotSetupStatus: true } }) }
		const setup = { executeMutation: vi.fn().mockResolvedValue({ data: { vkVideoBotSetupLink: 'https://vk.example/authorize' } }) }
		adminActions.useMutationVKVideoBotSetupStatus.mockReturnValue(status)
		adminActions.useMutationVKVideoBotSetupLink.mockReturnValue(setup)
		const open = vi.spyOn(window, 'open').mockReturnValue(null)

		const wrapper = mount(SetupVKVideoBot, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(wrapper.get('button').text()).toContain('Replace')
		await wrapper.get('button').trigger('click')
		expect(setup.executeMutation).not.toHaveBeenCalled()

		await wrapper.get('[data-test="replace-confirm"]').trigger('click')
		expect(setup.executeMutation).toHaveBeenCalledWith({})
		expect(open).toHaveBeenCalledWith('https://vk.example/authorize', 'vk-video-bot-setup', 'popup')
	})

	it('completes the callback, broadcasts refresh, and closes the popup', async () => {
		const complete = { executeMutation: vi.fn().mockResolvedValue({ data: { vkVideoBotSetupComplete: true } }) }
		adminActions.useMutationVKVideoBotSetupComplete.mockReturnValue(complete)
		route.query = { code: 'authorization-code', state: 'opaque-state' }
		const postMessage = vi.fn()
		vi.stubGlobal('BroadcastChannel', class {
			postMessage = postMessage
		})
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)

		mount(VKVideoBotCallback, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(complete.executeMutation).toHaveBeenCalledWith({ code: 'authorization-code', state: 'opaque-state' })
		expect(postMessage).toHaveBeenCalledWith('refresh')
		expect(sonner.toast.success).toHaveBeenCalled()
		expect(close).toHaveBeenCalled()
	})

	it('shows a callback failure without completing setup', async () => {
		const complete = { executeMutation: vi.fn() }
		adminActions.useMutationVKVideoBotSetupComplete.mockReturnValue(complete)
		route.query = { error: 'access_denied' }

		const wrapper = mount(VKVideoBotCallback, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(complete.executeMutation).not.toHaveBeenCalled()
		expect(wrapper.text()).toContain('VK Video Live authorization was not completed')
		expect(sonner.toast.error).toHaveBeenCalled()
	})

	it('renders VK alongside Twitch and Kick and toggles only the selected VK binding', async () => {
		const botStatuses = ref<Array<{
			dashboardId: string
			platform: string
			channelName: string
			botName: string
			enabled: boolean
		}>>([])
		const executeSubscription = vi.fn()
		const executeMutation = vi.fn().mockResolvedValue({ data: { botJoinLeave: true } })
		dashboard.useBotStatuses.mockReturnValue({ botStatuses, executeSubscription })
		dashboard.useBotJoinPart.mockReturnValue({ executeMutation })

		const wrapper = mount(HeaderBotStatus, {
			global: {
				stubs: {
					...componentStubs,
					CircleSvg: true,
					DropdownMenu: { template: '<div><slot /></div>' },
					DropdownMenuContent: { template: '<div><slot /></div>' },
					DropdownMenuItem: { template: '<button><slot /></button>' },
					DropdownMenuLabel: { template: '<p><slot /></p>' },
					DropdownMenuSeparator: true,
					DropdownMenuTrigger: { template: '<div><slot /></div>' },
				},
			},
		})

		botStatuses.value = [
			{ dashboardId: 'dashboard-1', platform: 'twitch', channelName: 'twitch-channel', botName: 'Twitch bot', enabled: true },
			{ dashboardId: 'dashboard-1', platform: 'kick', channelName: 'kick-channel', botName: 'Kick bot', enabled: true },
			{ dashboardId: 'dashboard-1', platform: 'vk_video_live', channelName: 'vk-channel', botName: 'VK bot', enabled: true },
		]
		await nextTick()

		expect(wrapper.find('[data-icon="simple-icons:twitch"]').exists()).toBe(true)
		expect(wrapper.find('[data-icon="simple-icons:kick"]').exists()).toBe(true)
		expect(wrapper.find('[data-icon="simple-icons:vk"]').exists()).toBe(true)
		expect(wrapper.text()).toContain('VK Video Live')

		const vkButton = wrapper.findAll('button').find((button) => button.text().includes('VK Video Live'))
		expect(vkButton).toBeDefined()
		if (!vkButton) return
		await vkButton.trigger('click')

		expect(executeMutation).toHaveBeenCalledWith({
			action: 'LEAVE',
			dashboardId: 'dashboard-1',
			platform: 'vk_video_live',
		})
		expect(executeSubscription).toHaveBeenCalled()

		const twitchButton = wrapper.findAll('button').find((button) => button.text().includes('Twitch'))
		expect(twitchButton).toBeDefined()
		if (!twitchButton) return
		await twitchButton.trigger('click')

		expect(executeMutation).toHaveBeenLastCalledWith({
			action: 'LEAVE',
			dashboardId: 'dashboard-1',
			platform: 'twitch',
		})
	})
})
