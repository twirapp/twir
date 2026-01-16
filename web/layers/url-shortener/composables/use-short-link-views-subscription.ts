import type { MaybeRef } from 'vue'

import { graphql } from '~/gql'

export const useShortLinkViewsSubscription = (shortLinkId: MaybeRef<string>) => {
	const { data, isPaused, pause, resume, error } = useSubscription({
		query: graphql(`
			subscription ShortLinkViewsUpdates($shortLinkId: String!) {
				shortLinkViewsUpdates(shortLinkId: $shortLinkId) {
					shortLinkId
					totalViews
					lastView {
						shortLinkId
						userId
						country
						city
						createdAt
					}
				}
			}
		`),
		variables: computed(() => ({
			shortLinkId: unref(shortLinkId),
		})),
	})

	const totalViews = computed(() => data.value?.shortLinkViewsUpdates?.totalViews ?? null)
	const lastView = computed(() => data.value?.shortLinkViewsUpdates?.lastView ?? null)

	return {
		totalViews,
		lastView,
		isPaused: computed(() => isPaused.value),
		pause,
		resume,
		error: computed(() => error.value),
	}
}
