import { useSubscription } from '@urql/vue'
import { computed } from 'vue'

import { useProfile } from '~~/layers/dashboard/api/auth.js'
import { useIntegrationsPageData } from '~~/layers/dashboard/api/integrations/integrations-page.js'

import { graphql } from '~/gql/gql.js'

const defaultNowPlayingTrack = {
	imageUrl: 'https://i.scdn.co/image/ab67616d0000b273e7fbc0883149094912559f2c',
	artist: 'Slipknot',
	title: 'Psychosocial',
	progressMs: 70_000,
	durationMs: 238_000,
}

export function useNowPlayingPreviewTrack() {
	const profile = useProfile()
	const integrationsPage = useIntegrationsPageData()

	const isSomeSongIntegrationEnabled = computed(() => {
		return (
			integrationsPage.spotifyData.value?.userName ||
			integrationsPage.lastfmData.value?.userName ||
			integrationsPage.vkData.value?.userName
		)
	})

	const currentTrackPaused = computed(() => {
		return !isSomeSongIntegrationEnabled.value || !profile.data.value?.channelApiKey
	})

	const { data: currentTrackSub } = useSubscription({
		query: graphql(`
			subscription NowPlayingOverlayNowPlaying($apiKey: String!) {
				nowPlayingCurrentTrack(apiKey: $apiKey) {
					title
					artist
					imageUrl
					progressMs
					durationMs
				}
			}
		`),
		get variables() {
			return {
				apiKey: profile.data.value!.channelApiKey!,
			}
		},
		pause: currentTrackPaused,
	})

	const track = computed(() => {
		if (currentTrackSub.value?.nowPlayingCurrentTrack) {
			return {
				imageUrl: currentTrackSub.value.nowPlayingCurrentTrack.imageUrl,
				artist: currentTrackSub.value.nowPlayingCurrentTrack.artist,
				title: currentTrackSub.value.nowPlayingCurrentTrack.title,
				progressMs: currentTrackSub.value.nowPlayingCurrentTrack.progressMs,
				durationMs: currentTrackSub.value.nowPlayingCurrentTrack.durationMs,
			}
		}

		return defaultNowPlayingTrack
	})

	return { track, isSomeSongIntegrationEnabled }
}
