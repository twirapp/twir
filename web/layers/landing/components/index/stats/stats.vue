<script setup lang="ts">
import { LandingStatsStoreKey } from '~/stores/landing-stats'

const formatter = Intl.NumberFormat('en-US', {
	notation: 'compact',
	maximumFractionDigits: 1,
})

const statsStore = useLandingStatsStore()
await callOnce(LandingStatsStoreKey, () => statsStore.fetchLandingStats())
const { t } = useI18n()

function formatNumber(value?: number | bigint) {
	return formatter.format(value ?? 0)
}

const totalChannels = computed(
	() =>
		(statsStore.stats?.twitchChannels ?? 0) +
		(statsStore.stats?.kickChannels ?? 0) +
		(statsStore.stats?.vkChannels ?? 0) +
		(statsStore.stats?.youtubeChannels ?? 0),
)

const stats = computed(() => [
	{
		key: 'activeChannels',
		value: formatNumber(totalChannels.value),
		isChannels: true,
	},
	{
		key: 'createdCommands',
		value: formatNumber(statsStore.stats?.createdCommands),
	},
	{
		key: 'usersSeen',
		value: formatNumber(statsStore.stats?.viewers),
	},
	{
		key: 'messagesProcessed',
		value: formatNumber(statsStore.stats?.messages),
	},
	{
		key: 'emotesProcessed',
		value: formatNumber(statsStore.stats?.usedEmotes),
	},
	{
		key: 'commandsProcessed',
		value: formatNumber(statsStore.stats?.usedCommands),
	},
])
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
				{{ t(`landing.stats.${stat.key}`) }}
			</span>
				<div v-if="stat.isChannels" class="flex gap-4 mt-2 text-[#ADB0B8] text-sm md:text-base font-medium">
					<div class="flex items-center gap-1.5" :title="t('landing.stats.platformChannels.twitch')">
					<Icon name="simple-icons:twitch" class="w-4 h-4 text-[#9146FF]" />
					<span>{{ formatNumber(statsStore.stats?.twitchChannels) }}</span>
				</div>
					<div class="flex items-center gap-1.5" :title="t('landing.stats.platformChannels.kick')">
					<Icon name="simple-icons:kick" class="w-4 h-4 text-[#53FC18]" />
					<span>{{ formatNumber(statsStore.stats?.kickChannels) }}</span>
				</div>
					<div class="flex items-center gap-1.5" :title="t('landing.stats.platformChannels.vk')">
					<Icon name="simple-icons:vk" class="w-4 h-4 text-[#0077FF]" />
					<span>{{ formatNumber(statsStore.stats?.vkChannels) }}</span>
				</div>
					<div class="flex items-center gap-1.5" :title="t('landing.stats.platformChannels.youtube')">
					<Icon name="simple-icons:youtube" class="w-4 h-4 text-[#FF0000]" />
					<span>{{ formatNumber(statsStore.stats?.youtubeChannels) }}</span>
				</div>
			</div>
		</div>
	</section>
</template>

<style scoped></style>
