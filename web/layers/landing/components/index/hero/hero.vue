<script setup lang="ts">
import { DISCORD_INVITE_URL } from '@twir/brand'
import UiButton from '~~/layers/landing/components/landing-ui-button.vue'

import { UserStoreKey } from '~/stores/user'

import HeroChat from './hero-chat.vue'

const userStore = useAuth()
const localePath = useLocalePath()
const { t } = useI18n()

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
						{{ t('landing.hero.updates') }}
						<SvgoArrowRight
							:fontControlled="false"
							class="h-4 w-4 shrink-0 stroke-white/50 stroke-[1.5]"
						/>
					</a>
					<h1
						class="max-w-2xl pt-4 text-[min(48px,11vw)] leading-[1.2] font-bold tracking-tight text-white lg:text-[64px]"
					>
						{{ t('landing.hero.title') }}
					</h1>

					<p
						class="max-w-xl pt-6 text-[min(18px,5vw)] leading-normal text-[#ADB0B8] lg:text-[20px]"
					>
						{{ t('landing.hero.description') }}
					</p>

					<div class="flex gap-4 pt-8 text-sm font-medium text-[#ADB0B8] md:text-base">
						<div class="flex items-center gap-1.5">
							<Icon
								name="simple-icons:twitch"
								class="h-4 w-4 text-[#9146FF]"
							/>
							<span>{{ t('landing.hero.platforms.twitch') }}</span>
						</div>
						<div class="flex items-center gap-1.5">
							<Icon
								name="simple-icons:kick"
								class="h-4 w-4 text-[#53FC18]"
							/>
							<span>{{ t('landing.hero.platforms.kick') }}</span>
						</div>
						<div class="flex items-center gap-1.5">
							<Icon
								name="simple-icons:vk"
								class="h-4 w-4 text-[#0077FF]"
							/>
							<span>{{ t('landing.hero.platforms.vk') }}</span>
						</div>
						<div class="flex items-center gap-1.5">
							<Icon
								name="simple-icons:youtube"
								class="h-4 w-4 text-[#FF0000]"
							/>
							<span>{{ t('landing.hero.platforms.youtube') }}</span>
						</div>
					</div>

					<div class="inline-flex w-full flex-col gap-3 pt-[48px] lg:flex-row">
						<NuxtLink
							v-slot="{ navigate, href }"
							:to="localePath('/compare')"
							custom
						>
							<UiButton
								:href="href!"
								variant="secondary"
								@click="navigate"
							>
								{{ t('landing.hero.learnMore') }}
							</UiButton>
						</NuxtLink>

						<NuxtLink
							v-if="userStore.userWithoutDashboards"
							v-slot="{ navigate, href }"
							to="/dashboard"
							custom
						>
							<UiButton
								:href="href!"
								variant="primary"
								@click="navigate"
							>
								{{ t('landing.hero.dashboard') }}
							</UiButton>
						</NuxtLink>

						<LoginDropdown
							v-else
							variant="hero"
						/>
					</div>
				</div>

				<HeroChat />
			</div>
		</div>
	</section>
</template>
