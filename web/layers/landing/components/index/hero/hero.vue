<script setup lang="ts">
import { DISCORD_INVITE_URL } from '@twir/brand'
import KickIcon from '~~/layers/landing/components/kick-icon.vue'
import UiButton from '~~/layers/landing/components/landing-ui-button.vue'

import { UserStoreKey } from '~/stores/user'

import HeroChat from './hero-chat.vue'

const userStore = useAuth()
const localePath = useLocalePath()

await Promise.all([callOnce(UserStoreKey, () => userStore.getUserDataWithoutDashboards())])
</script>

<template>
	<section class="overflow-hidden px-5 md:px-9">
		<div
			class="relative container mx-auto py-20 before:pointer-events-none before:absolute before:-right-12 before:-bottom-8 before:-z-10 before:h-[607px] before:w-[950px] before:-rotate-30 before:rounded-[950px] before:bg-[radial-gradient(50%_50.00%_at_50%_50%,#181F4E_0%,rgba(9,9,11,0.00)_100%)] before:content-[''] lg:py-28"
		>
			<div class="flex flex-wrap items-center justify-between gap-[60px] md:flex-nowrap">
				<div class="flex w-full flex-col items-start">
					<a
						class="flex items-center gap-[8px] rounded-full border border-[#72757d26] bg-[#1a1a22] py-1 pr-[10px] pl-2 text-sm font-medium text-[#E3E6ED] transition-shadow hover:bg-[#272730] focus-visible:ring-2 focus-visible:ring-[#72757d]/80 focus-visible:outline-none"
						:href="DISCORD_INVITE_URL"
						target="_blank"
					>
						🚀 View latest updates
						<SvgoArrowRight
							:fontControlled="false"
							class="h-4 w-4 shrink-0 stroke-white/50 stroke-[1.5]"
						/>
					</a>
					<h1
						class="max-w-2xl pt-4 text-[min(48px,11vw)] leading-[1.2] font-bold tracking-tight text-white lg:text-[64px]"
					>
						Engage your audience like never before
					</h1>

					<p
						class="max-w-xl pt-6 text-[min(18px,5vw)] leading-normal text-[#ADB0B8] lg:text-[20px]"
					>
						Our bot is the ultimate all-in-one solution for streamers looking to take their channel
						to the next level.
					</p>

					<div class="inline-flex w-full flex-col gap-3 pt-[48px] lg:flex-row">
						<UiButton
							:href="localePath('/compare')"
							variant="secondary"
						>
							Learn more
						</UiButton>
						<UiButton
							v-if="userStore.userWithoutDashboards"
							href="/dashboard"
							variant="primary"
						>
							Dashboard
						</UiButton>
						<template v-else>
							<UiButton
								as="button"
								variant="primary"
								@click="userStore.login()"
							>
								Start with Twitch
							</UiButton>
							<button
								class="xs:py-4 inline-flex items-center justify-center gap-2 rounded-lg bg-[#27272a] px-7 py-3 text-center text-base font-semibold whitespace-nowrap text-white transition-[background,box-shadow] hover:bg-[#27272a]/80 focus-visible:ring-4 focus-visible:ring-[#53FC18]/50 focus-visible:outline-none sm:text-lg"
								@click="userStore.loginWithKick()"
							>
								Start with Kick
								<KickIcon class="text-[#53FC18]" />
							</button>
						</template>
					</div>
				</div>

				<HeroChat />
			</div>
		</div>
	</section>
</template>
