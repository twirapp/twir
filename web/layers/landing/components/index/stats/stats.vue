<script setup lang="ts">
import { LandingStatsStoreKey, useLandingStatsStore } from '#layers/landing/stores/landing-stats'

const formatter = Intl.NumberFormat('en-US', {
	notation: 'compact',
	maximumFractionDigits: 1,
})

const { stats: statsData, fetchLandingStats } = useLandingStatsStore()
await callOnce(LandingStatsStoreKey, () => fetchLandingStats())

function formatNumber(value?: number | bigint) {
	return formatter.format(value ?? 0)
}

const stats = [
	{
		key: 'Channels',
		value: formatNumber(statsData?.channels),
	},
	{
		key: 'Created commands',
		value: formatNumber(statsData?.createdCommands),
	},
	{
		key: 'Viewers',
		value: formatNumber(statsData?.viewers),
	},
	{
		key: 'Messages',
		value: formatNumber(statsData?.messages),
	},
	{
		key: 'Used emotes',
		value: formatNumber(statsData?.usedEmotes),
	},
	{
		key: 'Used commands',
		value: formatNumber(statsData?.usedCommands),
	},
]
</script>

<template>
	<section id="stats" class="bg-[#17171A] px-5 py-6 gap-32 flex flex-wrap justify-center">
		<div
			v-for="stat of stats"
			:key="stat.key"
			class="inline-flex flex-col items-center justify-center"
		>
			<span
				class="font-semibold lg:text-6xl text-[min(40px,11vw)] text-white leading-[1.2] tracking-tight"
			>
				{{ stat.value }}
			</span>
			<span class="text-[#ADB0B8] lg:text-lg lg:mt-2 leading-normal whitespace-nowrap">
				{{ stat.key }}
			</span>
		</div>
	</section>
</template>

<style scoped></style>
