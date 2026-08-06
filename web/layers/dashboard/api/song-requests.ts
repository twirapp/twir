import type { MaybeRef } from 'vue'

import { useQuery, useSubscription } from '@urql/vue'
import { createGlobalState } from '@vueuse/core'
import { unref } from 'vue'
import { useMutation } from '~~/layers/dashboard/composables/use-mutation'

import type { SongRequestsSearchChannelOrVideoOpts } from '~/gql/graphql.js'

import { graphql } from '~/gql/gql.js'

export const useSongRequestsApi = createGlobalState(() => {
	const cacheKey = 'songRequests'

	const useSongRequestQuery = () =>
		useQuery({
			query: graphql(`
				query SongRequests {
					songRequests {
						enabled
						mode
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
						channelApiKey
						volume
						spotifyCapabilities {
							connected
							hasPlaybackScope
							canUseSpotify
							activeDevice {
								id
								name
								type
								isActive
							}
							selectedDevice {
								id
								name
								type
								isActive
							}
						}
					}
				}
			`),
			context: {
				additionalTypenames: [cacheKey],
			},
			variables: {},
		})

	const useSongRequestMutation = () =>
		useMutation(
			graphql(`
				mutation UpdateSongRequests($opts: SongRequestsSettingsOpts!) {
					songRequestsUpdate(opts: $opts)
				}
			`),
			[cacheKey]
		)

	const useYoutubeVideoOrChannelSearch = (opts: MaybeRef<SongRequestsSearchChannelOrVideoOpts>) =>
		useQuery({
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

	const spotifyQueueCacheKey = 'spotifySongRequestsQueue'

	const useSpotifyQueueQuery = (pause?: MaybeRef<boolean>) =>
		useQuery({
			query: graphql(`
				query SpotifySongRequestsQueue {
					spotifySongRequestsQueue {
						currentDevice {
							id
							name
							type
							isActive
						}
						requests {
							id
							title
							artist
							album
							durationMs
							requesterName
							requesterDisplayName
							source
							queuePosition
							status
							createdAt
						}
					}
				}
			`),
			context: {
				additionalTypenames: [spotifyQueueCacheKey],
			},
			get pause() {
				return unref(pause) ?? false
			},
			variables: {},
		})

	const useSpotifyTrackSearch = (query: MaybeRef<string>, pause?: MaybeRef<boolean>) =>
		useQuery({
			query: graphql(`
				query SpotifySongRequestsSearch($query: String!, $limit: Int) {
					spotifySongRequestsSearch(query: $query, limit: $limit) {
						id
						title
						artist
						album
						durationMs
						imageUrl
					}
				}
			`),
			context: {},
			get variables() {
				return { query: unref(query), limit: 10 }
			},
			get pause() {
				const q = unref(query)
				return (unref(pause) ?? false) || !q || q.length < 2
			},
		})

	const useSpotifyQueueSubscription = (channelId: MaybeRef<string>, pause?: MaybeRef<boolean>) =>
		useSubscription({
			query: graphql(`
				subscription SpotifySongRequestsQueueUpdated($channelId: UUID!) {
					spotifySongRequestsQueueUpdated(channelId: $channelId) {
						currentDevice {
							id
							name
							type
							isActive
						}
						requests {
							id
							title
							artist
							album
							durationMs
							requesterName
							requesterDisplayName
							source
							queuePosition
							status
							createdAt
						}
					}
				}
			`),
			get variables() {
				return { channelId: unref(channelId) }
			},
			get pause() {
				return (unref(pause) ?? false) || !unref(channelId)
			},
		})

	const useSpotifySelectDeviceMutation = () =>
		useMutation(
			graphql(`
				mutation SpotifySongRequestSelectDevice($deviceId: String!) {
					spotifySongRequestSelectDevice(deviceId: $deviceId)
				}
			`),
			[cacheKey, spotifyQueueCacheKey]
		)

	const useSpotifyRefreshDeviceMutation = () =>
		useMutation(
			graphql(`
				mutation SpotifySongRequestRefreshDevice {
					spotifySongRequestRefreshDevice {
						id
						name
						type
						isActive
					}
				}
			`),
			[cacheKey, spotifyQueueCacheKey]
		)

	const useSpotifySkipMutation = () =>
		useMutation(
			graphql(`
				mutation SpotifySongRequestSkip($requestId: UUID!) {
					spotifySongRequestSkip(requestId: $requestId)
				}
			`),
			[spotifyQueueCacheKey]
		)

	const useSpotifyCancelMutation = () =>
		useMutation(
			graphql(`
				mutation SpotifySongRequestCancel($requestId: UUID!) {
					spotifySongRequestCancel(requestId: $requestId)
				}
			`),
			[spotifyQueueCacheKey]
		)

	return {
		useSongRequestQuery,
		useSongRequestMutation,
		useYoutubeVideoOrChannelSearch,
		useSpotifyQueueQuery,
		useSpotifyQueueSubscription,
		useSpotifyTrackSearch,
		useSpotifySelectDeviceMutation,
		useSpotifyRefreshDeviceMutation,
		useSpotifySkipMutation,
		useSpotifyCancelMutation,
	}
})
