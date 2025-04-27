import { createGlobalState } from '@vueuse/core'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import type { Giveaway, GiveawayParticipant, GiveawaySubscriptionParticipant, GiveawayWinner } from '@/api/giveaways.js'

import { useGiveawaysApi } from '@/api/giveaways.js'
import { useToast } from '@/components/ui/toast'

export const useGiveaways = createGlobalState(() => {
	const giveawaysApi = useGiveawaysApi()
	const { data: giveaways, fetching: giveawaysListFetching } = giveawaysApi.useGiveawaysList()
	const route = useRoute()
	const router = useRouter()
	const { toast } = useToast()

	// Current giveaway state
	const currentGiveawayId = ref<string | null>(null)
	const participants = ref<GiveawayParticipant[]>([])
	const winners = ref<GiveawayWinner[]>([])

	// Computed values
	const giveawaysList = computed<Giveaway[]>(() => {
		return giveaways.value?.giveaways as Giveaway[] ?? []
	})

	const activeGiveaways = computed(() => {
		return giveawaysList.value.filter(g => !g.archivedAt && !g.endedAt)
	})

	const archivedGiveaways = computed(() => {
		return giveawaysList.value.filter(g => g.archivedAt || g.endedAt)
	})

	const currentGiveaway = computed(() => {
		if (!currentGiveawayId.value) return null
		return giveawaysList.value.find(g => g.id === currentGiveawayId.value) || null
	})

	// Mutations
	const createGiveawayMutation = giveawaysApi.useMutationCreateGiveaway()
	const startGiveawayMutation = giveawaysApi.useMutationStartGiveaway()
	const stopGiveawayMutation = giveawaysApi.useMutationStopGiveaway()
	const archiveGiveawayMutation = giveawaysApi.useMutationArchiveGiveaway()
	const chooseWinnersMutation = giveawaysApi.useMutationChooseWinners()

	// Actions
	async function createGiveaway(keyword: string) {
		try {
			const result = await createGiveawayMutation.executeMutation({ opts: { keyword } })
			if (result.error) {
				throw new Error(result.error.message)
			}
			toast({
				title: 'Giveaway created',
				description: `Giveaway with keyword "${keyword}" has been created`,
			})
			return result.data?.giveawaysCreate
		} catch (error) {
			toast({
				variant: 'destructive',
				title: 'Error creating giveaway',
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

	async function archiveGiveaway(id: string) {
		try {
			const result = await archiveGiveawayMutation.executeMutation({ id })
			if (result.error) {
				throw new Error(result.error.message)
			}
			toast({
				title: 'Giveaway archived',
				description: 'The giveaway has been archived successfully',
			})
			return result.data?.giveawaysArchive
		} catch (error) {
			toast({
				variant: 'destructive',
				title: 'Error archiving giveaway',
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
			winners.value = result.data?.giveawaysChooseWinners || []
			toast({
				title: 'Winners chosen',
				description: `${winners.value.length} winners have been chosen`,
			})
			return winners.value
		} catch (error) {
			toast({
				variant: 'destructive',
				title: 'Error choosing winners',
				description: error instanceof Error ? error.message : 'Unknown error',
			})
			return []
		}
	}

	// Navigation
	function viewGiveaway(id: string) {
		router.push(`/dashboard/giveaways/view/${id}`)
	}

	// Load giveaway data (participants and winners)
	function loadParticipants(giveawayId: string) {
		if (!giveawayId) return

		currentGiveawayId.value = giveawayId

		// Load giveaway details to get winners
		const { data: giveawayData } = giveawaysApi.useGiveaway(giveawayId)
		watch(giveawayData, (newData) => {
			if (newData?.giveaway) {
				// Set winners from giveaway data
				if (newData.giveaway.winners && newData.giveaway.winners.length > 0) {
					winners.value = newData.giveaway.winners
				}
			}
		}, { immediate: true })

		// Load participants
		const { data } = giveawaysApi.useGiveawayParticipants(giveawayId)
		watch(data, (newData) => {
			if (newData?.giveawayParticipants) {
				participants.value = newData.giveawayParticipants
			}
		}, { immediate: true })

		// Subscribe to new participants
		const { data: subscriptionData } = giveawaysApi.useSubscriptionGiveawayParticipants(giveawayId)
		watch(subscriptionData, (newData) => {
			if (newData?.giveawaysParticipants) {
				const newParticipant = newData.giveawaysParticipants
				// Check if participant already exists
				const exists = participants.value.some(p => p.userId === newParticipant.userId)
				if (!exists) {
					// Add new participant
					participants.value.push({
						id: `temp-${Date.now()}`,
						giveawayId: newParticipant.giveawayId,
						userId: newParticipant.userId,
						displayName: newParticipant.userDisplayName,
						isWinner: newParticipant.isWinner,
					})
				}
			}
		})
	}

	return {
		// State
		giveawaysList,
		giveawaysListFetching,
		participants,
		winners,
		currentGiveawayId,
		currentGiveaway,
		activeGiveaways,
		archivedGiveaways,

		// Actions
		createGiveaway,
		startGiveaway,
		stopGiveaway,
		archiveGiveaway,
		chooseWinners,
		viewGiveaway,
		loadParticipants,
	}
})
