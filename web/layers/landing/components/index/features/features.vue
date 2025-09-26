<script setup lang="ts">
import { DISCORD_INVITE_URL } from '@twir/brand'

import { featuresData } from './features-data.js'
import { LandingStatsStoreKey } from '~/stores/landing-stats'

const statsStore = useLandingStatsStore()
await callOnce(LandingStatsStoreKey, () => statsStore.fetchLandingStats())

function prepareDescription(d: string): string {
	return d
		.replaceAll('{landingStatsCreatedHastebins}', (statsStore.stats?.hasteBins ?? 0).toString())
		.replaceAll('{landingStatsCreatedShortUrls}', (statsStore.stats?.shortUrls ?? 0).toString())
}
</script>

<template>
	<section id="features">
		<div
			class="container flex flex-col items-center py-24 md:px-8 px-5 text-center purple-gradient"
		>
			<h2 class="font-semibold text-[#B0ADFF] text-base uppercase tracking-wider mb-3">Features</h2>
			<div class="flex flex-col items-center gap-[20px]">
				<h3 class="font-bold text-white text-4xl tracking-tight leading-tight">
					Explore our bot's awesome features
				</h3>
				<p class="max-w-3xl text-[#ADB0B8] text-xl leading-normal">
					Uncover a treasure trove of features that'll take your stream to the next level. Our bot
					got the tools to transform your stream adventure!
				</p>
			</div>
			<ul
				class="mt-[64px] grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-x-[32px] gap-y-[64px]"
			>
				<!-- Full width feature -->
				<li
					v-for="feature of featuresData"
					:key="feature.title"
					class="flex flex-col items-center text-center"
					:class="[feature.fullWidth ? 'col-span-full mb-8' : 'col-span-1']"
				>
					<div class="bg-[#B1AEFF]/[.12] p-2.5 rounded-full flex ring-8 ring-[#B1AEFF]/5 mb-5">
						<component :is="feature.icon" class="w-6 h-6 text-[#B0ADFF]" />
					</div>
					<span class="text-xl text-white font-semibold mb-2">{{ feature.title }}</span>
					<p
						class="text-base leading-normal text-[#ADB0B8] font-normal"
						:class="{ 'max-w-2xl': feature.fullWidth }"
					>
						{{ prepareDescription(feature.description) }}
					</p>
				</li>
			</ul>

			<div class="flex flex-col items-center gap-2 pt-16">
				<p class="max-w-3xl text-[#ADB0B8] text-sm leading-normal italic">
					This is by no means an exhaustive list of the bot's capabilities; it will be continuously
					updated for you as we strive to bring forth new and innovative features.
				</p>
				<a
					class="flex gap-[8px] items-center pr-[10px] pl-2 py-1 rounded-full bg-[#1a1a22] hover:bg-[#272730] border border-[#72757d26] font-medium text-sm text-[#E3E6ED] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#72757d]/80 transition-shadow"
					:href="DISCORD_INVITE_URL"
					target="_blank"
				>
					<SvgoSocialDiscord :fontControlled="false" class="h-4 w-4" />
					Join our Discord server to stay up to date
				</a>
			</div>
		</div>
	</section>
</template>
