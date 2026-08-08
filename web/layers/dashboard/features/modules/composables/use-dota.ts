import { useForm } from 'vee-validate'
import { computed, ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import * as z from 'zod'

import { useDotaApi } from '~~/layers/dashboard/api/dota'

const chatEventSchema = z.object({
	enabled: z.boolean().default(false),
	template: z.string().default(''),
	cooldown: z.number().int().min(0).default(0),
})

export const dotaFormSchema = z.object({
	enabled: z.boolean().default(false),
	mmr: z.number().int().min(0).default(0),
	mmrDelta: z.number().int().min(1).max(100).default(25),
	predictionSettings: z.object({
		enabled: z.boolean().default(false),
		titleTemplate: z.string().min(1).max(35, 'Twitch prediction titles are limited to 45 characters'),
		windowSeconds: z.number().int().min(30).max(1800).default(300),
	}),
	chatEvents: z.object({
		matchStarted: chatEventSchema,
		matchEnded: chatEventSchema,
		roshanKilled: chatEventSchema,
		aegisPickup: chatEventSchema,
	}),
	commandsSettings: z.object({
		mmr: z.boolean().default(true),
		wl: z.boolean().default(true),
		lg: z.boolean().default(true),
		gm: z.boolean().default(true),
		np: z.boolean().default(true),
		wp: z.boolean().default(true),
	}),
})

export type DotaFormSchema = z.infer<typeof dotaFormSchema>

const defaultChatEvents: DotaFormSchema['chatEvents'] = {
	matchStarted: { enabled: false, template: '🎮 Match started, playing {hero}!', cooldown: 0 },
	matchEnded: { enabled: false, template: 'Game over. MMR: {mmr} ({wins}-{losses})', cooldown: 0 },
	roshanKilled: { enabled: false, template: '🐲 Roshan killed by {team} at {time}!', cooldown: 0 },
	aegisPickup: { enabled: false, template: '🛡️ Aegis picked up at {time}!', cooldown: 0 },
}

export function useDota() {
	const { t } = useI18n()
	const isLoading = ref(false)

	const dotaApi = useDotaApi()
	const { data, fetching } = dotaApi.useQueryDota()
	const updateMutation = dotaApi.useMutationDotaUpdate()
	const steamLinkMutation = dotaApi.useMutationDotaSteamLink()
	const steamUnlinkMutation = dotaApi.useMutationDotaSteamUnlink()
	const resetSessionMutation = dotaApi.useMutationDotaResetSession()
	const regenerateTokenMutation = dotaApi.useMutationDotaRegenerateGsiToken()

	const settings = computed(() => data.value?.dota)

	const form = useForm<DotaFormSchema>({
		validationSchema: dotaFormSchema,
		initialValues: {
			enabled: false,
			mmr: 0,
			mmrDelta: 25,
			predictionSettings: {
				enabled: false,
				titleTemplate: 'Win this game?',
				windowSeconds: 300,
			},
			chatEvents: defaultChatEvents,
			commandsSettings: {
				mmr: true,
				wl: true,
				lg: true,
				gm: true,
				np: true,
				wp: true,
			},
		},
		validateOnMount: false,
		keepValuesOnUnmount: true,
	})

	watch(data, (v) => {
		if (!v?.dota) return

		form.setValues({
			enabled: v.dota.enabled,
			mmr: v.dota.mmr,
			mmrDelta: v.dota.mmrDelta,
			predictionSettings: {
				enabled: v.dota.predictionSettings.enabled,
				titleTemplate: v.dota.predictionSettings.titleTemplate || 'Win this game?',
				windowSeconds: v.dota.predictionSettings.windowSeconds || 300,
			},
			chatEvents: {
				matchStarted: withChatEventDefaults(v.dota.chatEvents.matchStarted, defaultChatEvents.matchStarted),
				matchEnded: withChatEventDefaults(v.dota.chatEvents.matchEnded, defaultChatEvents.matchEnded),
				roshanKilled: withChatEventDefaults(v.dota.chatEvents.roshanKilled, defaultChatEvents.roshanKilled),
				aegisPickup: withChatEventDefaults(v.dota.chatEvents.aegisPickup, defaultChatEvents.aegisPickup),
			},
			commandsSettings: {
				mmr: v.dota.commandsSettings.mmr,
				wl: v.dota.commandsSettings.wl,
				lg: v.dota.commandsSettings.lg,
				gm: v.dota.commandsSettings.gm,
				np: v.dota.commandsSettings.np,
				wp: v.dota.commandsSettings.wp,
			},
		})
	})

	function withChatEventDefaults(
		event: { enabled: boolean, template: string, cooldown: number },
		fallback: { enabled: boolean, template: string, cooldown: number },
	) {
		return {
			enabled: event.enabled,
			template: event.template || fallback.template,
			cooldown: event.cooldown,
		}
	}

	const handleSubmit = form.handleSubmit(async (values) => {
		isLoading.value = true
		try {
			const result = await updateMutation.executeMutation({
				input: {
					enabled: values.enabled,
					mmr: values.mmr,
					mmrDelta: values.mmrDelta,
					predictionSettings: values.predictionSettings,
					chatEvents: values.chatEvents,
					commandsSettings: values.commandsSettings,
				},
			})

			if (result.error) {
				toast.error(result.error.message || 'Error updating dota settings')
				return
			}

			toast.success(t('sharedTexts.saved'))
		} catch (err) {
			console.error(err)
			toast.error('An error occurred')
		} finally {
			isLoading.value = false
		}
	})

	async function linkSteam(queryString: string) {
		const result = await steamLinkMutation.executeMutation({ queryString })
		if (result.error) {
			toast.error(result.error.message || 'Failed to link Steam account')
			return false
		}

		toast.success(t('modules.dota.steam.linked'))
		return true
	}

	async function unlinkSteam() {
		const result = await steamUnlinkMutation.executeMutation({})
		if (result.error) {
			toast.error(result.error.message || 'Failed to unlink Steam account')
			return
		}

		toast.success(t('modules.dota.steam.unlinked'))
	}

	async function resetSession() {
		const result = await resetSessionMutation.executeMutation({})
		if (result.error) {
			toast.error(result.error.message || 'Failed to reset session')
			return
		}

		toast.success(t('sharedTexts.saved'))
	}

	async function regenerateGsiToken() {
		const result = await regenerateTokenMutation.executeMutation({})
		if (result.error) {
			toast.error(result.error.message || 'Failed to regenerate token')
			return
		}

		toast.success(t('sharedTexts.saved'))
	}

	return {
		form,
		handleSubmit,
		isLoading,
		fetching,
		settings,
		linkSteam,
		unlinkSteam,
		resetSession,
		regenerateGsiToken,
	}
}
