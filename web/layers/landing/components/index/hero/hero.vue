<script setup lang="ts">
import { DISCORD_INVITE_URL } from '@twir/brand'

import HeroChat from './hero-chat.vue'

import { useAuthLink, useProfile } from '~/layers/landing/api/user'
import UiButton from '~/layers/landing/components/ui-button.vue'

const { data: profile } = await useProfile()

const pageUrl = useRequestURL()

const redirectUrl = computed(() => {
	return pageUrl.origin
})

const { data: authLinkData } = await useAuthLink(redirectUrl)

const isLogged = computed(() => {
	return !!profile.value
})

const startButtonHref = computed(() => {
	return isLogged.value ? '/dashboard' : authLinkData.value?.authLink
})

const startButtonText = computed(() => {
	return isLogged.value ? 'Dashboard' : 'Start for free'
})
</script>

<template>
	<section class="px-5 md:px-9 overflow-hidden">
		<div
			class="container py-20 lg:py-28 relative before:content-[''] before:absolute before:w-[950px] before:h-[607px] before:-right-12 before:-rotate-[30deg] before:rounded-[950px] before:pointer-events-none before:bg-[radial-gradient(50%_50.00%_at_50%_50%,_#181F4E_0%,_rgba(9,9,11,0.00)_100%)] before:-z-10 before:-bottom-8"
		>
			<div class="flex justify-between items-center md:flex-nowrap flex-wrap gap-[60px]">
				<div class="flex flex-col items-start w-full">
					<a
						class="flex gap-[8px] items-center pr-[10px] pl-2 py-1 rounded-full bg-[#1a1a22] hover:bg-[#272730] border border-[#72757d26] font-medium text-sm text-[#E3E6ED] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#72757d]/80 transition-shadow"
						:href="DISCORD_INVITE_URL"
						target="_blank"
					>
						🚀 View latest updates
						<SvgoArrowRight :fontControlled="false" class="h-4 w-4 stroke-white/50 stroke-[1.5] flex-shrink-0" />
					</a>
					<h1
						class="pt-4 lg:text-[64px] text-[min(48px,11vw)] font-bold text-white tracking-tight leading-[1.2] max-w-2xl"
					>
						Engage your audience like never before
					</h1>

					<p class="pt-6 max-w-xl text-[#ADB0B8] lg:text-[20px] text-[min(18px,5vw)] leading-normal">
						Our Twitch bot is the ultimate all-in-one solution for streamers looking to take their
						channel to the next level.
					</p>

					<div class="pt-[48px] w-full inline-flex flex-col lg:flex-row gap-3">
						<UiButton href="#" variant="secondary">
							Learn more
						</UiButton>
						<UiButton :href="startButtonHref" variant="primary">
							{{ startButtonText }}
						</UiButton>
					</div>
				</div>

				<HeroChat />
			</div>
		</div>
	</section>
</template>
