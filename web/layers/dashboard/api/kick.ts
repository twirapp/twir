import type { MaybeRef } from 'vue'

import { useQuery } from '@urql/vue'
import { computed, unref } from 'vue'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'

import type { KickSearchCategoriesQuery } from '~/gql/graphql.js'
import { graphql } from '~/gql/gql.js'

export function useKickSearchCategories(query: MaybeRef<string>) {
	const pause = computed(() => unref(query).trim().length < 3)

	const gqlQuery = useQuery({
		query: graphql(`
			query KickSearchCategories($query: String!) {
				kickSearchCategories(query: $query) {
					categories {
						id
						name
						thumbnail
					}
				}
			}
		`),
		variables: computed(() => ({ query: unref(query).trim() })),
		pause,
	})

	return {
		data: computed<KickSearchCategoriesQuery['kickSearchCategories']['categories']>(
			() => gqlQuery.data.value?.kickSearchCategories.categories ?? []
		),
		isLoading: gqlQuery.fetching,
		error: gqlQuery.error,
	}
}

export function channelSetStreamInformationMutation() {
	return useMutation(
		graphql(`
			mutation ChannelSetStreamInformation($platform: Platform!, $title: String, $categoryId: String) {
				channelSetStreamInformation(platform: $platform, title: $title, categoryId: $categoryId)
			}
		`)
	)
}
