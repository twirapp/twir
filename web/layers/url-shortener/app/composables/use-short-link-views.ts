import type { MaybeRef } from 'vue'

import { graphql } from '~/gql'

export const useShortLinkViews = (
	shortLinkId: MaybeRef<string>,
	page: MaybeRef<number> = ref(0),
	perPage: MaybeRef<number> = ref(20)
) => {
	const { data, isPaused, pause, resume, error, executeQuery } = useQuery({
		query: graphql(`
			query ShortLinkViews($input: ShortLinkViewsInput!) {
				shortLinkViews(input: $input) {
					views {
						shortLinkId
						userId
						country
						city
						createdAt
					}
					total
				}
			}
		`),
		variables: computed(() => ({
			input: {
				shortLinkId: unref(shortLinkId),
				page: unref(page),
				perPage: unref(perPage),
			},
		})),
		pause: true,
	})

	const views = computed(() => data.value?.shortLinkViews?.views ?? [])
	const total = computed(() => data.value?.shortLinkViews?.total ?? 0)

	return {
		views,
		total,
		isPaused: computed(() => isPaused.value),
		pause,
		resume,
		error: computed(() => error.value),
		executeQuery,
	}
}
