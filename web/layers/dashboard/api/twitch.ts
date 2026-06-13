import { useQuery as useGqlQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed, unref } from 'vue'
import { useMutation } from '@/composables/use-mutation.ts'

import type { MaybeRef, Ref } from 'vue'

import { graphql } from '@/gql/gql.js'
import { TwitchGetUsersQuery, TwitchSearchChannelsQuery } from '@/gql/graphql.ts'

type TwitchIn = MaybeRef<string | string[] | null>
export function useTwitchGetUsers(opts: { ids?: TwitchIn; names?: TwitchIn }) {
	const queryIds = computed(() => {
		const rawIds = unref(opts.ids) ?? []
		let ids: string[] = Array.isArray(rawIds) ? rawIds : [rawIds]
		return ids.filter((n) => n !== '')
	})

	const queryNames = computed(() => {
		const rawNames = unref(opts.names) ?? []
		let names: string[] = Array.isArray(rawNames) ? rawNames : [rawNames]
		return names.filter((n) => n !== '')
	})

	const pause = computed(() => {
		return queryIds.value.length === 0 && queryNames.value.length === 0
	})

	const query = useGqlQuery({
		query: graphql(`
			query TwitchGetUsers($ids: [ID!], $names: [String!]) {
				twitchGetUsers(ids: $ids, names: $names) {
					id
					login
					displayName
					profileImageUrl
					description
				}
			}
		`),
		variables: computed(() => ({
			ids: queryIds.value.length > 0 ? queryIds.value : undefined,
			names: queryNames.value.length > 0 ? queryNames.value : undefined,
		})),
		pause,
	})

	return {
		data: computed<TwitchGetUsersQuery['twitchGetUsers']>(() => {
			if (!query.data.value) {
				return []
			}

			return query.data.value.twitchGetUsers
		}),
		isLoading: query.fetching,
		error: query.error,
	}
}

export function useTwitchSearchChannels(params: Ref<{ query: string; twirOnly?: boolean }>) {
	const pause = computed(() => !params.value.query)

	const query = useGqlQuery({
		query: graphql(`
			query TwitchSearchChannels($query: String!, $twirOnly: Boolean) {
				twitchSearchChannels(query: $query, twirOnly: $twirOnly) {
					channels {
						id
						login
						displayName
						profileImageUrl
						title
						gameName
						gameId
						isLive
					}
				}
			}
		`),
		variables: params,
		pause,
	})

	return {
		data: computed<TwitchSearchChannelsQuery['twitchSearchChannels']['channels']>(() => {
			if (!query.data.value) {
				return []
			}

			return query.data.value.twitchSearchChannels?.channels ?? []
		}),
		isLoading: query.fetching,
		error: query.error,
	}
}

export const useTwitchRewardsNew = createGlobalState(() =>
	useGqlQuery({
		query: graphql(`
			query GetChannelRewards {
				twitchRewards {
					id
					title
					cost
					imageUrls
					backgroundColor
					enabled
					usedTimes
					userInputRequired
				}
			}
		`),
	})
)

export function useTwitchSearchCategories(query: MaybeRef<string>) {
	const pause = computed(() => {
		const q = unref(query)
		return !q || q.trim() === ''
	})

	const gqlQuery = useGqlQuery({
		query: graphql(`
			query TwitchSearchCategories($query: String!) {
				twitchSearchCategories(query: $query) {
					categories {
						id
						name
						boxArtUrl
					}
				}
			}
		`),
		variables: computed(() => ({
			query: unref(query) || '',
		})),
		pause,
	})

	return {
		data: computed(() => {
			if (!gqlQuery.data.value) return []
			return gqlQuery.data.value.twitchSearchCategories.categories
		}),
		isLoading: gqlQuery.fetching,
		error: gqlQuery.error,
	}
}

export function useTwitchGetCategories(ids: MaybeRef<string[]>) {
	const pause = computed(() => {
		const idsValue = unref(ids)
		return !idsValue || idsValue.length === 0
	})

	const gqlQuery = useGqlQuery({
		query: graphql(`
			query TwitchGetCategories($ids: [ID!]!) {
				twitchGetCategories(ids: $ids) {
					categories {
						id
						name
						boxArtUrl
					}
				}
			}
		`),
		variables: computed(() => ({
			ids: unref(ids) || [],
		})),
		pause,
	})

	return {
		data: computed(() => {
			if (!gqlQuery.data.value) return []
			return gqlQuery.data.value.twitchGetCategories.categories
		}),
		isLoading: gqlQuery.fetching,
		error: gqlQuery.error,
	}
}

export function twitchSetChannelInformationMutation() {
	return useMutation(
		graphql(`
			mutation TwitchSetChannelInformation(
				$title: String
				$categoryId: String
			) {
					twitchSetChannelInformation(title: $title, categoryId: $categoryId)
			}
		`),
	)
}
