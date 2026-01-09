import { useQuery } from '@urql/vue'

import { graphql } from '@/gql'
import { useMutation } from '~/composables/use-mutation'

export function useNowPlayingOverlayApi() {
	const cacheKey = ['nowPlayingOverlay']

	const useNowPlayingQuery = () =>
		useQuery({
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
		})

	const useNowPlayingCreate = () =>
		useMutation(
			graphql(`
			mutation NowPlayingOverlayCreate($input: NowPlayingOverlayMutateOpts!) {
				nowPlayingOverlayCreate(opts: $input)
			}
		`),
			cacheKey
		)

	const useNowPlayingUpdate = () =>
		useMutation(
			graphql(`
			mutation NowPlayingOverlayUpdate($id: String!, $input: NowPlayingOverlayMutateOpts!) {
				nowPlayingOverlayUpdate(id: $id, opts: $input)
			}
		`),
			cacheKey
		)

	const useNowPlayingDelete = () =>
		useMutation(
			graphql(`
			mutation NowPlayingOverlayDelete($id: String!) {
				nowPlayingOverlayDelete(id: $id)
			}
		`),
			cacheKey
		)

	return {
		useNowPlayingQuery,
		useNowPlayingCreate,
		useNowPlayingUpdate,
		useNowPlayingDelete,
	}
}
