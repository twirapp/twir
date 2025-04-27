import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type {
	ChannelGiveawaySubscriptionParticipant,
	GiveawayFragment,
	GiveawayParticipantFragment,
	GiveawayWinnerFragment,
} from '@/gql/graphql.ts'
import type { MaybeRef } from 'vue'

import { graphql } from '@/gql'

graphql(`
	fragment Giveaway on ChannelGiveaway {
		id
		channelId
		createdAt
		updatedAt
		startedAt
		endedAt
		stoppedAt
		keyword
		createdByUserId
		archivedAt
		winners {
			...GiveawayWinner
		}
	}

	fragment GiveawayWinner on ChannelGiveawayWinner {
		displayName
		userId
		userLogin
		twitchProfile {
			profileImageUrl
			displayName
			login
		}
	}

	fragment GiveawayParticipant on ChannelGiveawayParticipants {
		displayName
		userId
		isWinner
		id
		giveawayId
	}

	fragment GiveawaySubscriptionParticipant on ChannelGiveawaySubscriptionParticipant {
		userId
		userLogin
		userDisplayName
		isWinner
		giveawayId
	}
`)

export type Giveaway = GiveawayFragment
export type GiveawayParticipant = GiveawayParticipantFragment
export type GiveawayWinner = GiveawayWinnerFragment
export type GiveawaySubscriptionParticipant = ChannelGiveawaySubscriptionParticipant

export const useGiveawaysApi = createGlobalState(() => {
	// Queries
	const useGiveawaysList = () => useQuery({
		query: graphql(`
			query GiveawaysList {
				giveaways {
					...Giveaway
				}
			}
		`),
	})

	const useGiveaway = (giveawayId: MaybeRef<string | null>) => useQuery({
		query: graphql(`
			query GiveawayById($giveawayId: String!) {
				giveaway(giveawayId: $giveawayId) {
					...Giveaway
				}
			}
		`),
		get variables() {
			const id = unref(giveawayId)!
			return { giveawayId: id }
		},
		pause: !!unref(giveawayId),
	})

	const useGiveawayParticipants = (giveawayId: MaybeRef<string | null>) => useQuery({
		query: graphql(`
			query GetGiveawayParticipants($giveawayId: String!) {
				giveawayParticipants(giveawayId: $giveawayId) {
					...GiveawayParticipant
				}
			}
		`),
		get variables() {
			const id = unref(giveawayId)!
			return { giveawayId: id }
		},
		pause: !!unref(giveawayId),
	})

	// Mutations
	const useMutationCreateGiveaway = () => useMutation(graphql(`
		mutation CreateGiveaway($opts: GiveawaysCreateInput!) {
			giveawaysCreate(opts: $opts) {
				...Giveaway
			}
		}
	`))

	const useMutationUpdateGiveaway = () => useMutation(graphql(`
		mutation UpdateGiveaway($id: String!, $opts: GiveawaysUpdateInput!) {
			giveawaysUpdate(id: $id, opts: $opts) {
				...Giveaway
			}
		}
	`))

	const useMutationRemoveGiveaway = () => useMutation(graphql(`
		mutation RemoveGiveaway($id: String!) {
			giveawaysRemove(id: $id) {
				...Giveaway
			}
		}
	`))

	const useMutationStartGiveaway = () => useMutation(graphql(`
		mutation StartGiveaway($id: String!) {
			giveawaysStart(id: $id) {
				...Giveaway
			}
		}
	`))

	const useMutationStopGiveaway = () => useMutation(graphql(`
		mutation StopGiveaway($id: String!) {
			giveawaysStop(id: $id) {
				...Giveaway
			}
		}
	`))

	const useMutationArchiveGiveaway = () => useMutation(graphql(`
		mutation ArchiveGiveaway($id: String!) {
			giveawaysArchive(id: $id) {
				...Giveaway
			}
		}
	`))

	const useMutationChooseWinners = () => useMutation(graphql(`
		mutation ChooseWinners($id: String!) {
			giveawaysChooseWinners(id: $id) {
				...GiveawayWinner
			}
		}
	`))

	// Subscriptions
	const useSubscriptionGiveawayParticipants = (giveawayId: MaybeRef<string | null>) => useSubscription({
		query: graphql(`
			subscription SubscribeToGiveawayParticipants($giveawayId: String!) {
				giveawaysParticipants(giveawayId: $giveawayId) {
					...GiveawaySubscriptionParticipant
				}
			}
		`),
		get variables() {
			const id = unref(giveawayId)!
			return { giveawayId: id }
		},
		pause: !!unref(giveawayId),
	})

	return {
		// Queries
		useGiveawaysList,
		useGiveaway,
		useGiveawayParticipants,

		// Mutations
		useMutationCreateGiveaway,
		useMutationUpdateGiveaway,
		useMutationRemoveGiveaway,
		useMutationStartGiveaway,
		useMutationStopGiveaway,
		useMutationArchiveGiveaway,
		useMutationChooseWinners,

		// Subscriptions
		useSubscriptionGiveawayParticipants,
	}
})
