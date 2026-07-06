<script setup lang="ts">
import { useMutation, useQuery, useSubscription } from '@urql/vue'
import { useDebouncedSeek } from '~~/layers/dashboard/composables/useDebouncedSeek.js'
import { useDebouncedVolume } from '~~/layers/dashboard/composables/useDebouncedVolume.js'
import { graphql } from '~/gql/gql.js'

const route = useRoute()
const channelApiKey = computed(() => route.params.channelApiKey as string)

const channelQuery = useQuery({
	query: graphql(`
		query WidgetChannelByApiKey($apiKey: String!) {
			channelByApiKey(apiKey: $apiKey) {
				id
			}
		}
	`),
	variables: computed(() => ({ apiKey: channelApiKey.value })),
	pause: computed(() => !channelApiKey.value),
})

const channelId = computed(() => channelQuery.data.value?.channelByApiKey?.id ?? '')

const widgetDataQuery = useQuery({
	query: graphql(`
		query WidgetSongRequestData($channelId: UUID!) {
			songRequestWidgetData(channelId: $channelId) {
				playbackState {
					videoId
					title
					position
					isPlaying
					volume
					updatedAt
				}
				queue {
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
		}
	`),
	variables: computed(() => ({ channelId: channelId.value })),
	pause: computed(() => !channelId.value),
})

const playbackSub = useSubscription({
	query: graphql(`
		subscription WidgetSongRequestPlaybackState($channelId: UUID!) {
			songRequestPlaybackState(channelId: $channelId) {
				videoId
				title
				position
				isPlaying
				volume
				updatedAt
			}
		}
	`),
	variables: computed(() => ({ channelId: channelId.value })),
	pause: computed(() => !channelId.value),
})

const queueSub = useSubscription({
	query: graphql(`
		subscription WidgetSongRequestQueueUpdated($channelId: UUID!) {
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
	`),
	variables: computed(() => ({ channelId: channelId.value })),
	pause: computed(() => !channelId.value),
})

const playbackState = computed(() =>
	playbackSub.data.value?.songRequestPlaybackState
	?? widgetDataQuery.data.value?.songRequestWidgetData?.playbackState
	?? null,
)

const queue = computed(() =>
	queueSub.data.value?.songRequestQueueUpdated
	?? widgetDataQuery.data.value?.songRequestWidgetData?.queue
	?? [],
)

const { executeMutation: playMutation } = useMutation(graphql(`
	mutation WidgetSongRequestPlay($channelId: UUID!, $videoId: String!) {
		songRequestPlay(channelId: $channelId, videoId: $videoId)
	}
`))

const { executeMutation: pauseMutation } = useMutation(graphql(`
	mutation WidgetSongRequestPause($channelId: UUID!) {
		songRequestPause(channelId: $channelId)
	}
`))

const { executeMutation: skipMutation } = useMutation(graphql(`
	mutation WidgetSongRequestSkip($channelId: UUID!) {
		songRequestSkip(channelId: $channelId)
	}
`))

const { executeMutation: setVolumeMutation } = useMutation(graphql(`
	mutation WidgetSongRequestSetVolume($channelId: UUID!, $volume: Int!) {
		songRequestSetVolume(channelId: $channelId, volume: $volume)
	}
`))

const { executeMutation: updatePositionMutation } = useMutation(graphql(`
	mutation WidgetSongRequestUpdatePosition($channelId: UUID!, $position: Float!) {
		songRequestUpdatePosition(channelId: $channelId, position: $position)
	}
`))

const { executeMutation: deleteFromQueueMutation } = useMutation(graphql(`
	mutation WidgetSongRequestDeleteFromQueue($channelId: UUID!, $videoId: String!) {
		songRequestDeleteFromQueue(channelId: $channelId, videoId: $videoId)
	}
`))

const { executeMutation: clearQueueMutation } = useMutation(graphql(`
	mutation WidgetSongRequestClearQueue($channelId: UUID!) {
		songRequestClearQueue(channelId: $channelId)
	}
`))

const { localVolume, handleVolumeInput, syncFromServer } = useDebouncedVolume(() => channelId.value, setVolumeMutation)

const {
	localPosition,
	handleSeekInput,
	syncFromServer: syncSeekFromServer,
} = useDebouncedSeek(() => channelId.value, updatePositionMutation)

const currentQueueItem = computed(() => {
	if (!playbackState.value?.videoId) return null
	return queue.value.find((q) => q.id === playbackState.value?.videoId) ?? null
})

const duration = computed(() => currentQueueItem.value?.durationSeconds ?? 0)

watch(() => playbackState.value?.volume, (v) => {
	if (v !== undefined) syncFromServer(v)
})

watch(() => playbackState.value?.position, (position) => {
	if (position !== undefined) syncSeekFromServer(position)
}, { immediate: true })

const hasPlayableVideo = computed(() => !!playbackState.value || queue.value.length > 0)

function handlePlayPause() {
	if (playbackState.value?.isPlaying) {
		pauseMutation({ channelId: channelId.value })
	} else {
		const videoId = playbackState.value?.videoId
			|| (queue.value[0]?.songLink?.match(/(?:v=|youtu\.be\/)([^&?/]+)/)?.[1] ?? '')
		if (videoId) {
			playMutation({ channelId: channelId.value, videoId })
		}
	}
}

function handleSkip() {
	skipMutation({ channelId: channelId.value })
}

function handleVolumeChange(e: Event) {
	const target = e.target as HTMLInputElement
	handleVolumeInput(Number(target.value))
}

function handleSeekChange(e: Event) {
	if (!playbackState.value) return

	const target = e.target as HTMLInputElement
	handleSeekInput(Number(target.value))
}

function handleDelete(videoId: string) {
	deleteFromQueueMutation({ channelId: channelId.value, videoId })
}

function handleClearQueue() {
	clearQueueMutation({ channelId: channelId.value })
}

function formatTime(seconds: number): string {
	const m = Math.floor(seconds / 60)
	const s = Math.floor(seconds % 60)
	return `${m}:${s.toString().padStart(2, '0')}`
}

function formatDuration(seconds: number): string {
	const m = Math.floor(seconds / 60)
	const s = seconds % 60
	return `${m}:${s.toString().padStart(2, '0')}`
}

const progressPercent = computed(() => {
	if (!playbackState.value || duration.value <= 0) return 0
	return Math.min((localPosition.value / duration.value) * 100, 100)
})
</script>

<template>
	<div class="widget-root">
		<div class="now-playing">
			<div class="track-title">
				{{ playbackState?.title ?? 'No track playing' }}
			</div>
			<div class="progress-bar">
				<div
					class="progress-fill"
					:style="{ width: `${progressPercent}%` }"
				/>
			</div>
			<div class="progress-time">
				{{ formatTime(localPosition) }}
			</div>
			<input
				type="range"
				min="0"
				:disabled="!playbackState"
				:max="duration || 1"
				:value="localPosition"
				class="seek-slider"
				@input="handleSeekChange"
			>
			<div class="controls">
				<button
					class="ctrl-btn"
					:disabled="!hasPlayableVideo"
					@click="handlePlayPause"
				>
					<Icon :name="playbackState?.isPlaying ? 'lucide:pause' : 'lucide:play'" class="icon" />
				</button>
				<button
					class="ctrl-btn"
					:disabled="!hasPlayableVideo"
					@click="handleSkip"
				>
					<Icon name="lucide:skip-forward" class="icon" />
				</button>
				<div class="volume-control">
					<Icon name="lucide:volume-2" class="icon-small" />
					<input
						type="range"
						min="0"
						max="100"
						:value="localVolume"
						class="volume-slider"
						@input="handleVolumeChange"
					>
				</div>
			</div>
		</div>

		<div class="queue-section">
			<div class="queue-header">
				<span class="queue-title">Queue ({{ queue.length }})</span>
				<button
					v-if="queue.length > 0"
					class="clear-btn"
					@click="handleClearQueue"
				>
					Clear All
				</button>
			</div>
			<div class="queue-list">
				<div
					v-for="(item, index) in queue"
					:key="item.id"
					class="queue-item"
				>
					<span class="queue-index">{{ index + 1 }}.</span>
					<div class="queue-item-info">
						<span class="queue-item-title">{{ item.title }}</span>
						<span class="queue-item-meta">
							{{ item.orderedByDisplayName }} &middot; {{ formatDuration(item.durationSeconds) }}
						</span>
					</div>
					<button
						class="delete-btn"
						@click="handleDelete(item.id)"
					>
						<Icon name="lucide:x" class="icon-tiny" />
					</button>
				</div>
				<div v-if="queue.length === 0" class="queue-empty">
					No songs in queue
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.widget-root {
	width: 100%;
	height: 100%;
	background: rgba(0, 0, 0, 0.9);
	border-radius: 8px;
	overflow: hidden;
	display: flex;
	flex-direction: column;
	font-family: system-ui, -apple-system, sans-serif;
	color: white;
}

.now-playing {
	flex-shrink: 0;
	padding: 12px;
	border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.track-title {
	font-size: 14px;
	font-weight: 600;
	margin-bottom: 6px;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.progress-bar {
	width: 100%;
	height: 4px;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 2px;
	margin-bottom: 4px;
}

.progress-fill {
	height: 100%;
	background: #8b5cf6;
	border-radius: 2px;
	transition: width 1s linear;
}

.progress-time {
	font-size: 11px;
	color: rgba(255, 255, 255, 0.5);
	margin-bottom: 4px;
}

.seek-slider {
	width: 100%;
	height: 4px;
	-webkit-appearance: none;
	appearance: none;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 2px;
	outline: none;
	margin-bottom: 8px;
}

.seek-slider::-webkit-slider-thumb {
	-webkit-appearance: none;
	appearance: none;
	width: 12px;
	height: 12px;
	background: #8b5cf6;
	border-radius: 50%;
	cursor: pointer;
}

.seek-slider:disabled::-webkit-slider-thumb {
	cursor: default;
	opacity: 0.3;
}

.controls {
	display: flex;
	align-items: center;
	gap: 8px;
}

.ctrl-btn {
	background: rgba(255, 255, 255, 0.1);
	border: none;
	border-radius: 4px;
	padding: 6px;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: background 0.15s;
}

.ctrl-btn:hover:not(:disabled) {
	background: rgba(255, 255, 255, 0.2);
}

.ctrl-btn:disabled {
	opacity: 0.3;
	cursor: default;
}

.icon {
	width: 18px;
	height: 18px;
}

.icon-small {
	width: 14px;
	height: 14px;
	flex-shrink: 0;
}

.icon-tiny {
	width: 14px;
	height: 14px;
}

.volume-control {
	display: flex;
	align-items: center;
	gap: 6px;
	flex: 1;
}

.volume-slider {
	flex: 1;
	height: 4px;
	-webkit-appearance: none;
	appearance: none;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 2px;
	outline: none;
}

.volume-slider::-webkit-slider-thumb {
	-webkit-appearance: none;
	appearance: none;
	width: 12px;
	height: 12px;
	background: #8b5cf6;
	border-radius: 50%;
	cursor: pointer;
}

.queue-section {
	flex: 1;
	display: flex;
	flex-direction: column;
	min-height: 0;
}

.queue-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 10px 12px 6px;
	flex-shrink: 0;
}

.queue-title {
	font-size: 13px;
	font-weight: 600;
}

.clear-btn {
	background: none;
	border: 1px solid rgba(239, 68, 68, 0.5);
	border-radius: 4px;
	padding: 3px 8px;
	font-size: 11px;
	color: #ef4444;
	cursor: pointer;
	transition: all 0.15s;
}

.clear-btn:hover {
	background: rgba(239, 68, 68, 0.15);
	border-color: #ef4444;
}

.queue-list {
	flex: 1;
	overflow-y: auto;
	padding: 0 12px 12px;
}

.queue-item {
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 6px 0;
	border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.queue-index {
	font-size: 12px;
	color: rgba(255, 255, 255, 0.4);
	flex-shrink: 0;
	width: 18px;
	text-align: right;
}

.queue-item-info {
	flex: 1;
	min-width: 0;
}

.queue-item-title {
	display: block;
	font-size: 13px;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.queue-item-meta {
	font-size: 11px;
	color: rgba(255, 255, 255, 0.4);
}

.delete-btn {
	background: none;
	border: none;
	padding: 4px;
	cursor: pointer;
	color: rgba(255, 255, 255, 0.3);
	display: flex;
	align-items: center;
	transition: color 0.15s;
	flex-shrink: 0;
}

.delete-btn:hover {
	color: #ef4444;
}

.queue-empty {
	text-align: center;
	padding: 20px 0;
	font-size: 13px;
	color: rgba(255, 255, 255, 0.3);
}
</style>
