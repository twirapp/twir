<script setup lang="ts">
import { DISCORD_INVITE_URL, GITHUB_REPOSITORY_URL } from '@twir/brand'
import { anchorLinks } from '~~/layers/landing/layouts/default/anchor-links'

import TwirLogo from '../../../../../libs/brand/src/logo.svg?component'
</script>

<template>
	<footer class="bg-[#19181a] px-5 pt-12">
		<div class="container mx-auto">
			<div class="flex flex-col items-center gap-6">
				<div class="flex items-center gap-3">
					<TwirLogo
						alt="Twir"
						class="h-10 w-10"
					/>
					<span class="text-2xl font-semibold tracking-wide text-white uppercase">Twir</span>
				</div>
				<nav
					class="flex items-center justify-center"
					aria-label="Footer Navigation"
				>
					<ul class="m-0 flex list-none flex-wrap gap-4 p-0">
						<NuxtLink
							v-for="anchor of anchorLinks"
							:key="anchor.href"
							:to="anchor.href"
							class="navigation-link rounded-md px-3 py-2 leading-6 font-medium text-[#ADB0B8] hover:text-[#D5D8DF] focus-visible:text-[#D5D8DF]"
							prefetch-on="interaction"
						>
							{{ anchor.label }}
						</NuxtLink>
					</ul>
				</nav>
			</div>
			<div
				class="mt-12 flex justify-between gap-4 border-t border-stone-600 p-8 text-stone-400 max-[720px]:flex-col max-[720px]:items-center"
			>
				<span>2023 - {{ new Date().getFullYear() }} Twir</span>
				<div class="flex gap-[14px]">
					<a
						:href="GITHUB_REPOSITORY_URL"
						target="_blank"
						rel="noopener noreferrer"
						aria-label="Github repository"
					>
						<SvgoSocialGithub
							:fontControlled="false"
							class="h-6 w-6"
							aria-hidden
						/>
					</a>
					<a
						:href="DISCORD_INVITE_URL"
						target="_blank"
						rel="noopener noreferrer"
						aria-label="Discord channel for asking a question and keeping up with updates"
					>
						<SvgoSocialDiscord
							:fontControlled="false"
							class="h-6 w-6"
							aria-hidden
						/>
					</a>
				</div>
			</div>
		</div>
	</footer>
</template>

<style scoped></style>
