import { useSubscription } from '@urql/vue'
import { computed, type Ref } from 'vue'

import { useMutation } from '~~/layers/dashboard/composables/use-mutation.js'
import { graphql } from '~/gql/gql.js'

const SongRequestPlaybackStateSubscription = graphql(`
	subscription SongRequestPlaybackState($channelId: UUID!) {
		songRequestPlaybackState(channelId: $channelId) {
			videoId
			title
			position
			isPlaying
			volume
			updatedAt
		}
	}
`)

const SongRequestQueueUpdatedSubscription = graphql(`
	subscription SongRequestQueueUpdated($channelId: UUID!) {
		songRequestQueueUpdated(channelId: $channelId) {
			id
			title
			songLink
			durationSeconds
			orderedByName
			orderedByDisplayName
			queuePosition
			createdAt
		}
	}
`)

const SongRequestPlayMutation = graphql(`
	mutation SongRequestPlay($channelId: UUID!, $videoId: String!) {
		songRequestPlay(channelId: $channelId, videoId: $videoId)
	}
`)

const SongRequestPauseMutation = graphql(`
	mutation SongRequestPause($channelId: UUID!) {
		songRequestPause(channelId: $channelId)
	}
`)

const SongRequestSkipMutation = graphql(`
	mutation SongRequestSkip($channelId: UUID!) {
		songRequestSkip(channelId: $channelId)
	}
`)

const SongRequestSetVolumeMutation = graphql(`
	mutation SongRequestSetVolume($channelId: UUID!, $volume: Int!) {
		songRequestSetVolume(channelId: $channelId, volume: $volume)
	}
`)

const SongRequestReorderMutation = graphql(`
	mutation SongRequestReorder($channelId: UUID!, $videoIds: [String!]!) {
		songRequestReorder(channelId: $channelId, videoIds: $videoIds)
	}
`)

const SongRequestDeleteFromQueueMutation = graphql(`
	mutation SongRequestDeleteFromQueue($channelId: UUID!, $videoId: String!) {
		songRequestDeleteFromQueue(channelId: $channelId, videoId: $videoId)
	}
`)

const SongRequestClearQueueMutation = graphql(`
	mutation SongRequestClearQueue($channelId: UUID!) {
		songRequestClearQueue(channelId: $channelId)
	}
`)

export function useSongRequestGql(channelId: Ref<string>) {
	const playbackStateSub = useSubscription({
		query: SongRequestPlaybackStateSubscription,
		get variables() {
			return { channelId: channelId.value }
		},
	})

	const queueSub = useSubscription({
		query: SongRequestQueueUpdatedSubscription,
		get variables() {
			return { channelId: channelId.value }
		},
	})

	const playbackState = computed(() => playbackStateSub.data.value?.songRequestPlaybackState ?? null)
	const queue = computed(() => queueSub.data.value?.songRequestQueueUpdated ?? [])

	const { executeMutation: play } = useMutation(SongRequestPlayMutation)
	const { executeMutation: pause } = useMutation(SongRequestPauseMutation)
	const { executeMutation: skip } = useMutation(SongRequestSkipMutation)
	const { executeMutation: setVolume } = useMutation(SongRequestSetVolumeMutation)
	const { executeMutation: reorder } = useMutation(SongRequestReorderMutation)
	const { executeMutation: deleteFromQueue } = useMutation(SongRequestDeleteFromQueueMutation)
	const { executeMutation: clearQueue } = useMutation(SongRequestClearQueueMutation)

	return {
		playbackState,
		queue,
		play: (videoId: string) => play({ channelId: channelId.value, videoId }),
		pause: () => pause({ channelId: channelId.value }),
		skip: () => skip({ channelId: channelId.value }),
		setVolume: (volume: number) => setVolume({ channelId: channelId.value, volume }),
		reorder: (videoIds: string[]) => reorder({ channelId: channelId.value, videoIds }),
		deleteFromQueue: (videoId: string) => deleteFromQueue({ channelId: channelId.value, videoId }),
		clearQueue: () => clearQueue({ channelId: channelId.value }),
	}
}
