<script setup lang="ts">
import CinemaStyle from './CinemaStyle.vue'
import CompactStyle from './CompactStyle.vue'
import PillStyle from './PillStyle.vue'
import PortraitStyle from './PortraitStyle.vue'
import StudioStyle from './StudioStyle.vue'
import TickerStyle from './TickerStyle.vue'
import {
	type SongRequestOverlayProps,
	type SongRequestOverlayStyle,
	getYouTubeThumbnailUrl,
	resolveSongRequestOverlayProps,
} from './types'

const props = defineProps<SongRequestOverlayProps>()

defineSlots<{
	media?: (props: { thumbnailUrl: string }) => unknown
}>()

const styleComponents = {
	CINEMA: CinemaStyle,
	COMPACT: CompactStyle,
	TICKER: TickerStyle,
	STUDIO: StudioStyle,
	PORTRAIT: PortraitStyle,
	PILL: PillStyle,
} satisfies Record<SongRequestOverlayStyle, typeof CinemaStyle>

const resolvedProps = computed(() => resolveSongRequestOverlayProps(props))
const activeComponent = computed(() => styleComponents[resolvedProps.value.style])
const mediaClass = computed(() => ({
	'song-request-overlay-media--cinema': resolvedProps.value.style === 'CINEMA',
	'song-request-overlay-media--compact': resolvedProps.value.style === 'COMPACT',
	'song-request-overlay-media--audio': !['CINEMA', 'COMPACT'].includes(resolvedProps.value.style),
}))
const visualProps = computed(() => {
	const { style: _style, ...rest } = resolvedProps.value
	return rest
})
const thumbnailUrl = computed(() => getYouTubeThumbnailUrl(resolvedProps.value.videoId))
</script>

<template>
	<div class="relative size-full overflow-hidden">
		<div class="pointer-events-none absolute inset-0 z-10">
			<component
				:is="activeComponent"
				v-bind="visualProps"
			>
				<template
					v-if="$slots.media"
					#media
				>
					<div class="size-full" />
				</template>
			</component>
		</div>

		<div
			class="song-request-overlay-media absolute z-0 overflow-hidden"
			:class="mediaClass"
		>
			<slot
				name="media"
				:thumbnail-url="thumbnailUrl"
			/>
		</div>
	</div>
</template>

<style scoped>
.song-request-overlay-media--cinema,
.song-request-overlay-media--audio {
	inset: 0;
}

.song-request-overlay-media--compact {
	left: 0.5rem;
	bottom: 5.25rem;
	width: min(calc(100% - 1rem), 30rem);
	aspect-ratio: 16 / 9;
	border-radius: 0.75rem 0.75rem 0 0;
}

.song-request-overlay-media--audio {
	opacity: 0;
	pointer-events: none;
}

@media (min-width: 640px) {
	.song-request-overlay-media--compact {
		left: 1rem;
		bottom: 6.75rem;
		width: min(calc(100% - 2rem), 30rem);
	}
}
</style>
