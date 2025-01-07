<script setup lang="ts">
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import TwirLogo from '@twir/brand/src/logo.svg?component'

import { anchorLinks } from '~/layers/landing/layouts/default/anchor-links'
</script>

<template>
	<footer class="pt-12 px-5 bg-[#19181a]">
		<div class="container mx-auto">
			<div class="flex flex-col gap-6 items-center">
				<div class="flex gap-3 items-center">
					<TwirLogo alt="Twir" class="w-10 h-10" />
					<span class="text-white text-2xl font-semibold uppercase tracking-wide">Twir</span>
				</div>
				<nav class="flex justify-center items-center" aria-label="Footer Navigation">
					<ul class="list-none p-0 m-0 flex flex-wrap gap-4">
						<li
							v-for="anchor of anchorLinks"
							:key="anchor.href"
						>
							<a
								:href="anchor.href"
								class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] transition-colors"
							>
								{{ anchor.label }}
							</a>
						</li>
					</ul>
				</nav>
				<nav class="flex justify-center items-center" aria-label="Footer Navigation">
					<ul class="list-none p-0 m-0 flex flex-wrap gap-4">
						<NuxtLink to="/terms" class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] transition-colors">
							Terms of service
						</NuxtLink>
					</ul>
				</nav>
			</div>
			<div class="flex justify-between p-8 text-stone-400 mt-12 border-t border-stone-600 gap-4 max-[720px]:flex-col max-[720px]:items-center">
				<span>&copy; 2023 - {{ new Date().getFullYear() }} TwirApp. All rights reserved.</span>
				<div class="flex gap-[14px]">
					<a
						:href="GITHUB_REPOSITORY_URL"
						target="_blank"
						rel="noopener noreferrer"
						aria-label="Github repository"
					>
						<SvgoSocialGithub :fontControlled="false" class="w-6 h-6" aria-hidden />
					</a>
					<a
						:href="DISCORD_INVITE_URL"
						target="_blank"
						rel="noopener noreferrer"
						aria-label="Discord channel for asking a question and keeping up with updates"
					>
						<SvgoSocialDiscord :fontControlled="false" class="w-6 h-6" aria-hidden />
					</a>
				</div>
			</div>
		</div>
	</footer>
</template>

<style scoped>

</style>
