import { useMutation, useQuery } from '@tanstack/vue-query'
import { useQuery as useGqlQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { computed, isRef, unref } from 'vue'

import type { ComputedRef, MaybeRef, Ref } from 'vue'

import { protectedApiClient } from '@/api/twirp.js'
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

export function useTwitchSearchCategories(query: string | Ref<string>) {
	return useQuery({
		queryKey: ['twitchSearchCategories', query || ''],
		queryFn: async () => {
			const input = isRef(query) ? query.value : query
			if (!input) return { categories: [] }

			const call = await protectedApiClient.twitchSearchCategories({ query: input })
			return call.response
		},
	})
}

export function useTwitchGetCategories(
	ids: MaybeRef<string[]> | ComputedRef<string[]>,
	options?: { keepPreviousData?: boolean }
) {
	return useQuery({
		queryKey: ['twitchGetCategories', ids || ''],
		queryFn: async () => {
			const input = isRef(ids) ? ids.value : ids
			if (!input) return { categories: [] }

			const call = await protectedApiClient.twitchGetCategories({ ids: input })
			return call.response
		},
		keepPreviousData: options?.keepPreviousData,
	})
}

export function twitchSetChannelInformationMutation() {
	return useMutation({
		mutationKey: ['twitchSetChannelInformation'],
		mutationFn: async (req: { categoryId: string; title: string }) => {
			await protectedApiClient.twitchSetChannelInformation(req)
		},
	})
}
