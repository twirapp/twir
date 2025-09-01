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
		stoppedAt
		keyword
		createdByUserId
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
	const useGiveawaysList = () =>
		useQuery({
			query: graphql(`
				query GiveawaysList {
					giveaways {
						...Giveaway
					}
				}
			`),
		})

	const useGiveaway = (giveawayId: MaybeRef<string | null>) => {
		return useQuery({
			query: graphql(`
				query GiveawayById($giveawayId: String!) {
					giveaway(giveawayId: $giveawayId) {
						...Giveaway
						participants {
							...GiveawayParticipant
						}
					}
				}
			`),
			get variables() {
				return { giveawayId: unref(giveawayId)! }
			},
			pause: true,
			requestPolicy: 'cache-and-network',
		})
	}

	// Mutations
	const useMutationCreateGiveaway = () =>
		useMutation(
			graphql(`
				mutation CreateGiveaway($opts: GiveawaysCreateInput!) {
					giveawaysCreate(opts: $opts) {
						...Giveaway
					}
				}
			`)
		)

	const useMutationUpdateGiveaway = () =>
		useMutation(
			graphql(`
				mutation UpdateGiveaway($id: String!, $opts: GiveawaysUpdateInput!) {
					giveawaysUpdate(id: $id, opts: $opts) {
						...Giveaway
					}
				}
			`)
		)

	const useMutationRemoveGiveaway = () =>
		useMutation(
			graphql(`
				mutation RemoveGiveaway($id: String!) {
					giveawaysRemove(id: $id) {
						...Giveaway
					}
				}
			`)
		)

	const useMutationStartGiveaway = () =>
		useMutation(
			graphql(`
				mutation StartGiveaway($id: String!) {
					giveawaysStart(id: $id) {
						...Giveaway
					}
				}
			`)
		)

	const useMutationStopGiveaway = () =>
		useMutation(
			graphql(`
				mutation StopGiveaway($id: String!) {
					giveawaysStop(id: $id) {
						...Giveaway
					}
				}
			`)
		)

	const useMutationChooseWinners = () =>
		useMutation(
			graphql(`
				mutation ChooseWinners($id: String!) {
					giveawaysChooseWinners(id: $id) {
						...GiveawayWinner
					}
				}
			`)
		)

	// Subscriptions
	const useSubscriptionGiveawayParticipants = (giveawayId: MaybeRef<string | null>) => {
		return useSubscription({
			query: graphql(`
				subscription SubscribeToGiveawayParticipants($giveawayId: String!) {
					giveawaysParticipants(giveawayId: $giveawayId) {
						...GiveawaySubscriptionParticipant
					}
				}
			`),
			get variables() {
				return { giveawayId: unref(giveawayId)! }
			},
			pause: true,
		})
	}

	return {
		// Queries
		useGiveawaysList,
		useGiveaway,

		// Mutations
		useMutationCreateGiveaway,
		useMutationUpdateGiveaway,
		useMutationRemoveGiveaway,
		useMutationStartGiveaway,
		useMutationStopGiveaway,
		useMutationChooseWinners,

		// Subscriptions
		useSubscriptionGiveawayParticipants,
	}
})
