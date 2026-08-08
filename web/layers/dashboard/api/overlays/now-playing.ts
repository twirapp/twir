import { useQuery } from '@urql/vue'
import type { Ref } from 'vue'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation'
import { graphql } from '~/gql/gql.js'

export function useNowPlayingOverlayApi() {
	const cacheKey = ['nowPlayingOverlay']

	const useNowPlayingQuery = (pause?: Ref<boolean>) => useQuery({
		query: graphql(`
			query NowPlayingOverlays {
				nowPlayingOverlays {
					id
					channelId
					preset
					hideTimeout
					fontWeight
					fontFamily
					backgroundColor
					showImage
				}
			}
		`),
		context: {
			additionalTypenames: cacheKey,
		},
		variables: {},
		pause,
	})

	const useNowPlayingCreate = () => useMutation(
		graphql(`
			mutation NowPlayingOverlayCreate($input: NowPlayingOverlayMutateOpts!) {
				nowPlayingOverlayCreate(opts: $input)
			}
		`),
		cacheKey,
	)

	const useNowPlayingUpdate = () => useMutation(
		graphql(`
			mutation NowPlayingOverlayUpdate($id: String!, $input: NowPlayingOverlayMutateOpts!) {
				nowPlayingOverlayUpdate(id: $id, opts: $input)
			}
		`),
		cacheKey,
	)

	const useNowPlayingDelete = () => useMutation(
		graphql(`
			mutation NowPlayingOverlayDelete($id: String!) {
				nowPlayingOverlayDelete(id: $id)
			}
		`),
		cacheKey,
	)

	return {
		useNowPlayingQuery,
		useNowPlayingCreate,
		useNowPlayingUpdate,
		useNowPlayingDelete,
	}
}
