import { flushPromises, mount } from '@vue/test-utils'
import { nextTick, ref } from 'vue'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import HeaderBotStatus from '~~/layers/dashboard/layout/header/header-bot-status.vue'
import SetupYouTubeBot from './setup-youtube-bot.vue'
import YouTubeBotCallback from './youtube-bot-callback.vue'

const adminActions = vi.hoisted(() => ({
	youtubeBotSetupBroadcastChannelName: 'youtube_bot_setup',
	useMutationYouTubeBotSetupComplete: vi.fn(),
	useMutationYouTubeBotSetupLink: vi.fn(),
	useMutationYouTubeBotSetupStatus: vi.fn(),
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

describe('YouTube bot setup', () => {
	beforeEach(() => {
		adminActions.useMutationYouTubeBotSetupComplete.mockReset()
		adminActions.useMutationYouTubeBotSetupLink.mockReset()
		adminActions.useMutationYouTubeBotSetupStatus.mockReset()
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

	it('opens the YouTube authorization popup when the global bot is disconnected', async () => {
		const status = { executeMutation: vi.fn().mockResolvedValue({ data: { youtubeBotSetupStatus: false } }) }
		const setup = { executeMutation: vi.fn().mockResolvedValue({ data: { youtubeBotSetupLink: 'https://youtube.example/authorize' } }) }
		adminActions.useMutationYouTubeBotSetupStatus.mockReturnValue(status)
		adminActions.useMutationYouTubeBotSetupLink.mockReturnValue(setup)
		const open = vi.spyOn(window, 'open').mockReturnValue(null)

		const wrapper = mount(SetupYouTubeBot, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(wrapper.text()).toContain('Not connected')
		await wrapper.get('button').trigger('click')

		expect(setup.executeMutation).toHaveBeenCalledWith({})
		expect(open).toHaveBeenCalledWith('https://youtube.example/authorize', 'youtube-bot-setup', 'popup')
	})

	it('requires confirmation before replacing the connected global bot', async () => {
		const status = { executeMutation: vi.fn().mockResolvedValue({ data: { youtubeBotSetupStatus: true } }) }
		const setup = { executeMutation: vi.fn().mockResolvedValue({ data: { youtubeBotSetupLink: 'https://youtube.example/authorize' } }) }
		adminActions.useMutationYouTubeBotSetupStatus.mockReturnValue(status)
		adminActions.useMutationYouTubeBotSetupLink.mockReturnValue(setup)
		const open = vi.spyOn(window, 'open').mockReturnValue(null)

		const wrapper = mount(SetupYouTubeBot, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(wrapper.get('button').text()).toContain('Replace')
		await wrapper.get('button').trigger('click')
		expect(setup.executeMutation).not.toHaveBeenCalled()

		await wrapper.get('[data-test="replace-confirm"]').trigger('click')
		expect(setup.executeMutation).toHaveBeenCalledWith({})
		expect(open).toHaveBeenCalledWith('https://youtube.example/authorize', 'youtube-bot-setup', 'popup')
	})

	it('completes the callback, broadcasts refresh, and closes the popup', async () => {
		const complete = { executeMutation: vi.fn().mockResolvedValue({ data: { youtubeBotSetupComplete: true } }) }
		adminActions.useMutationYouTubeBotSetupComplete.mockReturnValue(complete)
		route.query = { code: 'authorization-code', state: 'opaque-state' }
		const postMessage = vi.fn()
		vi.stubGlobal('BroadcastChannel', class {
			postMessage = postMessage
		})
		const close = vi.spyOn(window, 'close').mockImplementation(() => undefined)

		mount(YouTubeBotCallback, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(complete.executeMutation).toHaveBeenCalledWith({ code: 'authorization-code', state: 'opaque-state' })
		expect(postMessage).toHaveBeenCalledWith('refresh')
		expect(sonner.toast.success).toHaveBeenCalled()
		expect(close).toHaveBeenCalled()
	})

	it('shows a callback failure without completing setup', async () => {
		const complete = { executeMutation: vi.fn() }
		adminActions.useMutationYouTubeBotSetupComplete.mockReturnValue(complete)
		route.query = { error: 'access_denied' }

		const wrapper = mount(YouTubeBotCallback, { global: { stubs: componentStubs } })
		await flushPromises()

		expect(complete.executeMutation).not.toHaveBeenCalled()
		expect(wrapper.text()).toContain('YouTube authorization was not completed')
		expect(sonner.toast.error).toHaveBeenCalled()
	})

	it('renders YouTube alongside Twitch and Kick and toggles only the selected YouTube binding', async () => {
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
			{ dashboardId: 'dashboard-1', platform: 'youtube', channelName: 'youtube-channel', botName: 'YouTube bot', enabled: true },
		]
		await nextTick()

		expect(wrapper.find('[data-icon="simple-icons:twitch"]').exists()).toBe(true)
		expect(wrapper.find('[data-icon="simple-icons:kick"]').exists()).toBe(true)
		expect(wrapper.find('[data-icon="simple-icons:youtube"]').exists()).toBe(true)
		expect(wrapper.text()).toContain('YouTube')

		const youtubeButton = wrapper.findAll('button').find((button) => button.text().includes('YouTube'))
		expect(youtubeButton).toBeDefined()
		if (!youtubeButton) return
		await youtubeButton.trigger('click')

		expect(executeMutation).toHaveBeenCalledWith({
			action: 'LEAVE',
			dashboardId: 'dashboard-1',
			platform: 'youtube',
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
