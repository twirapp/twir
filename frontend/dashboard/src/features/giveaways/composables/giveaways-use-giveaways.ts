import { createGlobalState } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import type {
	Giveaway,
	GiveawayParticipant,
	GiveawaySubscriptionParticipant,
	GiveawayWinner,
} from '@/api/giveaways.js'

import {
	useGiveawaysApi,
} from '@/api/giveaways.js'
import { useToast } from '@/components/ui/toast'

export const useGiveaways = createGlobalState(() => {
	const giveawaysApi = useGiveawaysApi()
	const { data: giveaways, fetching: giveawaysListFetching } = giveawaysApi.useGiveawaysList()
	const router = useRouter()
	const { toast } = useToast()
	const { t } = useI18n()

	// Current giveaway state
	const currentGiveawayId = ref<string | null>(null)
	const participants = ref<GiveawayParticipant[]>([])
	const winners = ref<GiveawayWinner[]>([])

	// Computed values
	const giveawaysList = computed<Giveaway[]>(() => {
		return giveaways.value?.giveaways as Giveaway[] ?? []
	})

	const activeGiveaways = computed<Giveaway[]>(() => {
		return giveawaysList.value.filter(g => !g.stoppedAt)
	})

	const archivedGiveaways = computed<Giveaway[]>(() => {
		return giveawaysList.value.filter(g => g.stoppedAt)
	})

	const currentGiveaway = computed(() => {
		if (!currentGiveawayId.value) return null
		return giveawaysList.value.find(g => g.id === currentGiveawayId.value) as Giveaway
	})

	// Mutations
	const createGiveawayMutation = giveawaysApi.useMutationCreateGiveaway()
	const startGiveawayMutation = giveawaysApi.useMutationStartGiveaway()
	const stopGiveawayMutation = giveawaysApi.useMutationStopGiveaway()
	const chooseWinnersMutation = giveawaysApi.useMutationChooseWinners()

	// Actions
	async function createGiveaway(keyword: string) {
		try {
			const result = await createGiveawayMutation.executeMutation({ opts: { keyword } })
			if (result.error) {
				throw new Error(result.error.message)
			}
			toast({
				title: t('giveaways.notifications.created'),
				description: t('giveaways.notifications.createdDescription', { keyword }),
			})
			return result.data?.giveawaysCreate
		} catch (error) {
			toast({
				variant: 'destructive',
				title: t('giveaways.notifications.error'),
				description: error instanceof Error ? error.message : 'Unknown error',
			})
			return null
		}
	}

	async function startGiveaway(id: string) {
		try {
			const result = await startGiveawayMutation.executeMutation({ id })
			if (result.error) {
				throw new Error(result.error.message)
			}
			toast({
				title: 'Giveaway started',
				description: 'The giveaway has been started successfully',
			})
			return result.data?.giveawaysStart
		} catch (error) {
			toast({
				variant: 'destructive',
				title: 'Error starting giveaway',
				description: error instanceof Error ? error.message : 'Unknown error',
			})
			return null
		}
	}

	async function stopGiveaway(id: string) {
		try {
			const result = await stopGiveawayMutation.executeMutation({ id })
			if (result.error) {
				throw new Error(result.error.message)
			}
			toast({
				title: 'Giveaway stopped',
				description: 'The giveaway has been stopped successfully',
			})
			return result.data?.giveawaysStop
		} catch (error) {
			toast({
				variant: 'destructive',
				title: 'Error stopping giveaway',
				description: error instanceof Error ? error.message : 'Unknown error',
			})
			return null
		}
	}

	async function chooseWinners(id: string) {
		try {
			const result = await chooseWinnersMutation.executeMutation({ id })
			if (result.error) {
				throw new Error(result.error.message)
			}
			winners.value.push(...result.data?.giveawaysChooseWinners as GiveawayWinner[] || [])
			toast({
				title: t('giveaways.notifications.winnersChosen'),
				description: t('giveaways.notifications.winnersChosenDescription', { count: winners.value.length }),
			})
			return winners.value
		} catch (error) {
			toast({
				variant: 'destructive',
				title: t('giveaways.notifications.errorChoosingWinners'),
				description: error instanceof Error ? error.message : 'Unknown error',
			})
			return []
		}
	}

	// Navigation
	function viewGiveaway(id: string) {
		router.push(`/dashboard/giveaways/view/${id}`)
	}

	// Use the API with reactive giveaway ID directly
	const { data: giveawayData } = giveawaysApi.useGiveaway(currentGiveawayId)
	const { data: participantsData } = giveawaysApi.useGiveawayParticipants(currentGiveawayId)
	const { data: participantsSubscriptionData } = giveawaysApi.useSubscriptionGiveawayParticipants(currentGiveawayId)

	// Watch for giveaway data changes
	watch(giveawayData, (giveawayData) => {
		if (!giveawayData?.giveaway) return
		const g = giveawayData.giveaway as unknown as Giveaway
		if (g?.winners && g.winners.length > 0) {
			winners.value = g.winners as GiveawayWinner[]
		}
	}, { immediate: true })

	// Watch for participants data changes
	watch(participantsData, (participantsData) => {
		const newParticipants = participantsData?.giveawayParticipants

		if (newParticipants) {
			participants.value = newParticipants as GiveawayParticipant[]
		}
	}, { immediate: true })

	// Watch for new participants from subscription
	watch(participantsSubscriptionData, (data) => {
		if (!data) return
		const newParticipant = data.giveawaysParticipants
		const participant = newParticipant as unknown as GiveawaySubscriptionParticipant

		const exists = participants.value.some(p => p.userId === participant.userId)
		if (!exists) {
			// Add new participant
			participants.value.push({
				id: `temp-${Date.now()}`,
				giveawayId: participant.giveawayId,
				userId: participant.userId,
				displayName: participant.userDisplayName,
				isWinner: participant.isWinner,
			})
		}
	})

	// Function to set the current giveaway ID
	function loadParticipants(giveawayId: string) {
		if (!giveawayId) return
		currentGiveawayId.value = giveawayId
	}

	const participantsWithFixedWinners = computed(() => {
		return participants.value.map((p) => {
			const isWinner = winners.value.some(w => w.userId === p.userId)
			return {
				...p,
				isWinner,
			}
		})
	})

	return {
		// State
		giveawaysList,
		giveawaysListFetching,
		participants: participantsWithFixedWinners,
		winners,
		currentGiveawayId,
		currentGiveaway,
		activeGiveaways,
		archivedGiveaways,

		// Actions
		createGiveaway,
		startGiveaway,
		stopGiveaway,
		chooseWinners,
		viewGiveaway,
		loadParticipants,
	}
})
