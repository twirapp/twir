<script setup lang="ts">
import { type SongRequestOverlayVisualProps, getSongRequestPlaybackMetrics } from './types'

const props = defineProps<SongRequestOverlayVisualProps>()

const viewportRef = useTemplateRef<HTMLElement>('viewportRef')
const measureRef = useTemplateRef<HTMLElement>('measureRef')
const firstTextRef = useTemplateRef<HTMLElement>('firstTextRef')
const secondTextRef = useTemplateRef<HTMLElement>('secondTextRef')
const isOverflowing = ref(false)
const travelDistance = ref(0)

let resizeObserver: ResizeObserver | null = null

const metrics = computed(() => getSongRequestPlaybackMetrics(props.position, props.duration))
const tickerLabel = computed(() =>
	props.requester ? `@${props.requester} • ${props.title}` : props.title
)
const animationDuration = computed(() => travelDistance.value / props.tickerSpeed)
const marqueeStyle = computed(() => ({
	'--ticker-travel-distance': `${travelDistance.value}px`,
	animationDuration: `${animationDuration.value}s`,
}))

function updateMeasurements() {
	const viewport = viewportRef.value
	const measure = measureRef.value
	if (!viewport || !measure) return

	isOverflowing.value = measure.scrollWidth > viewport.clientWidth

	if (!isOverflowing.value) {
		travelDistance.value = 0
		return
	}

	void nextTick(() => {
		const firstText = firstTextRef.value
		const secondText = secondTextRef.value
		if (!firstText || !secondText) return

		travelDistance.value = Math.max(0, secondText.offsetLeft - firstText.offsetLeft)
	})
}

function renderTickerContent() {
	updateMeasurements()
}

watch(tickerLabel, () => {
	void nextTick(updateMeasurements)
})

onMounted(() => {
	updateMeasurements()

	if (typeof ResizeObserver !== 'undefined') {
		resizeObserver = new ResizeObserver(updateMeasurements)
		if (viewportRef.value) resizeObserver.observe(viewportRef.value)
		if (measureRef.value) resizeObserver.observe(measureRef.value)
	}

	window.addEventListener('resize', renderTickerContent)
})

onUnmounted(() => {
	resizeObserver?.disconnect()
	window.removeEventListener('resize', renderTickerContent)
})
</script>

<template>
	<div
		class="absolute inset-x-0 bottom-0 overflow-hidden text-sm font-medium sm:text-base"
		:style="{ backgroundColor: tickerBackgroundColor, color: tickerTextColor }"
	>
		<div class="flex min-h-14 items-center gap-4 px-4 sm:px-6">
			<div
				ref="viewportRef"
				class="relative min-w-0 flex-1 overflow-hidden"
			>
				<span
					ref="measureRef"
					aria-hidden="true"
					class="invisible absolute whitespace-nowrap"
				>
					{{ tickerLabel }}
				</span>

				<div
					v-if="isOverflowing"
					class="song-request-ticker-track flex w-max gap-12 whitespace-nowrap motion-reduce:hidden motion-reduce:animate-none"
					:style="marqueeStyle"
				>
					<span ref="firstTextRef">
						<template v-if="requester">
							@{{ requester }}
							<span :style="{ color: accentColor }"> • </span>
						</template>
						{{ title }}
					</span>
					<span
						ref="secondTextRef"
						aria-hidden="true"
					>
						<template v-if="requester">
							@{{ requester }}
							<span :style="{ color: accentColor }"> • </span>
						</template>
						{{ title }}
					</span>
				</div>
				<p
					v-if="isOverflowing"
					class="hidden truncate whitespace-nowrap motion-reduce:block"
				>
					<template v-if="requester">
						@{{ requester }}
						<span :style="{ color: accentColor }"> • </span>
					</template>
					{{ title }}
				</p>
				<p
					v-else
					class="truncate whitespace-nowrap"
				>
					<template v-if="requester">
						@{{ requester }}
						<span :style="{ color: accentColor }"> • </span>
					</template>
					{{ title }}
				</p>
			</div>

			<p class="shrink-0 font-mono text-xs tabular-nums opacity-75 sm:text-sm">
				{{ metrics.formattedPosition }} / {{ metrics.formattedDuration }}
			</p>
		</div>
		<div class="h-0.5 bg-black/20">
			<div
				class="h-full transition-[width] duration-200 motion-reduce:transition-none"
				:style="{ width: `${metrics.progress}%`, backgroundColor: accentColor }"
			/>
		</div>
	</div>
</template>

<style scoped>
@keyframes song-request-ticker-marquee {
	to {
		transform: translateX(calc(-1 * var(--ticker-travel-distance)));
	}
}

.song-request-ticker-track {
	animation-name: song-request-ticker-marquee;
	animation-timing-function: linear;
	animation-iteration-count: infinite;
}
</style>
