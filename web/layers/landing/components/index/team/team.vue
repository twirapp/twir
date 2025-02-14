<script setup lang="ts">
import { developers, socialLinkLabels } from './team-data.js'
</script>

<template>
	<section id="team" class="container py-[72px] team">
		<div class="flex flex-col gap-3 p-5">
			<h2 class="font-semibold text-[#B0ADFF] text-base uppercase text-center">
				Team
			</h2>
			<h3 class="text-white text-4xl font-bold text-center">
				They're doing their best for you
			</h3>
		</div>

		<div class="pt-16 flex items-center justify-items-center flex-wrap gap-8 px-4">
			<div v-for="developer of developers" :key="developer.username" class="flex gap-6">
				<img
					:src="developer.avatarUrl"
					class="size-20 rounded-full"
					:alt="`Avatar of ${developer.username} from GitHub`"
					loading="lazy"
				/>
				<div class="flex flex-col">
					<span class="font-semibold text-white text-lg">{{ developer.username }}</span>
					<span class="font-normal mt-1 text-[#ADB0B8] text-base">{{ developer.description }}</span>
					<div class="flex mt-4 gap-4">
						<a v-for="(value, name) of developer.social" :key="name" :href="value" target="_blank" :aria-label="socialLinkLabels[name]">
							<SvgoSocialTwitch v-if="name === 'twitch'" :fontControlled="false" :filled="true" class="w-5 h-5 text-[#ADB0B8]" />
							<SvgoSocialTelegram v-else-if="name === 'telegram'" :fontControlled="false" :filled="true" class="w-5 h-5 text-[#ADB0B8]" />
							<SvgoSocialGithub v-else-if="name === 'github'" :fontControlled="false" :filled="true" class="w-5 h-5 text-[#ADB0B8]" />
							<SvgoSocialInstagram v-else-if="name === 'instagram'" :fontControlled="false" :filled="true" class="w-5 h-5 text-[#ADB0B8]" />
						</a>
					</div>
				</div>
			</div>
		</div>
	</section>
</template>

<style scoped>
.team {
	background-image: radial-gradient(50% 50% at 50% 50%, #181f4e 0%, rgba(9, 9, 11, 0) 100%);
}
</style>
