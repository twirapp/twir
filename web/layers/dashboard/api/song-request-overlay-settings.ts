import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'

import { graphql } from '~/gql/gql.js'

const cacheKey = 'SongRequestOverlaySettings'

export const useSongRequestOverlaySettingsApi = createGlobalState(() => {
	const useSettingsQuery = () =>
		useQuery({
			query: graphql(`
				query SongRequestOverlaySettingsDashboard {
					songRequestOverlaySettings {
						style
						accentColor
						tickerBackgroundColor
						tickerTextColor
						tickerSpeed
						hideOnPause
					}
				}
			`),
			context: {
				additionalTypenames: [cacheKey],
			},
			variables: {},
		})

	const useUpdateMutation = () =>
		useMutation(
			graphql(`
				mutation SongRequestOverlaySettingsUpdate($opts: SongRequestOverlaySettingsUpdateInput!) {
					songRequestOverlaySettingsUpdate(opts: $opts) {
						style
						accentColor
						tickerBackgroundColor
						tickerTextColor
						tickerSpeed
						hideOnPause
					}
				}
			`),
			[cacheKey]
		)

	return {
		useSettingsQuery,
		useUpdateMutation,
	}
})
