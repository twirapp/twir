import { useQuery } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'

import type { SongRequestsSearchChannelOrVideoOpts } from '@/gql/graphql'
import type { MaybeRef } from 'vue'

import { useMutation } from '@/composables/use-mutation'
import { graphql } from '@/gql'

export const useSongRequestsApi = createGlobalState(() => {
	const cacheKey = 'songRequests'

	const useSongRequestQuery = () => useQuery({
		query: graphql(`
			query SongRequests {
				songRequests {
					enabled
					acceptOnlyWhenOnline
					maxRequests
					channelPointsRewardId
					announcePlay
					neededVotesForSkip
					user {
						maxRequests
						minWatchTime
						minMessages
						minFollowTime
					}
					song {
						minLength
						maxLength
						minViews
						acceptedCategories
					}
					denyList {
						users
						songs
						channels
						artistsNames
						words
					}
					translations {
						nowPlaying
						notEnabled
						noText
						acceptOnlyWhenOnline
						user {
							denied
							maxRequests
							minMessages
							minWatched
							minFollow
						}
						song {
							denied
							notFound
							alreadyInQueue
							ageRestrictions
							cannotGetInformation
							live
							maxLength
							minLength
							requestedMessage
							maximumOrdered
							minViews
						}
						channel {
							denied
						}
					}
					takeSongFromDonationMessages
					playerNoCookieMode
				}
			}
		`),
		context: {
			additionalTypenames: [cacheKey],
		},
		variables: {},
	})

	const useSongRequestMutation = () => useMutation(
		graphql(`
			mutation UpdateSongRequests($opts: SongRequestsSettingsOpts!) {
				songRequestsUpdate(opts: $opts)
			}
		`),
		[cacheKey],
	)

	const useYoutubeVideoOrChannelSearch = (opts: MaybeRef<SongRequestsSearchChannelOrVideoOpts>) => useQuery({
		query: graphql(`
			query YoutubeVideoOrChannelSearch($opts: SongRequestsSearchChannelOrVideoOpts!) {
				songRequestsSearchChannelOrVideo(opts: $opts) {
					items {
						id
						title
						thumbnail
					}
				}
			}
		`),
		context: {},
		get variables() {
			return {
				opts: unref(opts),
			}
		},
	})

	return {
		useSongRequestQuery,
		useSongRequestMutation,
		useYoutubeVideoOrChannelSearch,
	}
})
