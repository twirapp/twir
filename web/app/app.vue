<script lang="ts" setup>
const description = 'Powerful and useful Twitch bot that helps manage chat on big channels. Developed from streamers for streamers with love.'
const siteName = 'Twir'
const keywords = 'twitch, bot, chat, moderation, moderation bot, twitch bot, twitch chat, twitch moderation, twitch moderation bot'

useHead({
	title: siteName,
	titleTemplate: (title) => title === siteName ? title : `${title} - ${siteName}`,
	meta: [
		{ name: 'darkreader-lock' },
		{ name: 'author', content: 'Satont, me@satont.dev' },
		{ name: 'keywords', content: keywords },
		{ name: 'viewport', content: 'width=device-width, initial-scale=1' },
		{ name: 'description', content: description },
	],
	script: [
		!import.meta.dev
			? {
				'src': 'https://umami.twir.app/script.js',
				'async': true,
				'data-website-id': '77f327a4-1744-4fb5-9fea-97ffbbc76806',
			}
			: {},
	],
})

useSeoMeta({
	title: siteName,
	description,
	viewport: 'width=device-width, initial-scale=1',
	ogDescription: description,
	applicationName: siteName,
	ogSiteName: siteName,
	ogTitle: siteName,
	ogType: 'website',
	keywords: 'twitch, bot, chat, moderation, moderation bot, twitch bot, twitch chat, twitch moderation, twitch moderation bot',
})
</script>

<template>
	<NuxtLayout>
		<NuxtPage />
	</NuxtLayout>
</template>
