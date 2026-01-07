<script setup lang="ts">
import { LandingStatsStoreKey } from '~/stores/landing-stats'

const formatter = Intl.NumberFormat('en-US', {
	notation: 'compact',
	maximumFractionDigits: 1,
})

const statsStore = useLandingStatsStore()
await callOnce(LandingStatsStoreKey, () => statsStore.fetchLandingStats())

function formatNumber(value?: number | bigint) {
	return formatter.format(value ?? 0)
}

const stats = [
	{
		key: 'Active Channels',
		value: formatNumber(statsStore.stats?.channels),
	},
	{
		key: 'Created Commands',
		value: formatNumber(statsStore.stats?.createdCommands),
	},
	{
		key: 'Users Seen',
		value: formatNumber(statsStore.stats?.viewers),
	},
	{
		key: 'Messages Processed',
		value: formatNumber(statsStore.stats?.messages),
	},
	{
		key: 'Emotes Processed',
		value: formatNumber(statsStore.stats?.usedEmotes),
	},
	{
		key: 'Commands Processed',
		value: formatNumber(statsStore.stats?.usedCommands),
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
