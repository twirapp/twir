<script setup lang="ts">
import 'vue-sonner/style.css'
import { Toaster } from '@/components/ui/sonner'

const description =
	'Powerful and useful Twitch bot that helps manage chat on big channels. Developed from streamers for streamers with love.'
const siteName = `Twir${import.meta.dev ? ' dev' : ''}`
const keywords =
	'twitch, bot, chat, moderation, moderation bot, twitch bot, twitch chat, twitch moderation, twitch moderation bot'

const metaImg = '/meta.webp'

useColorMode()
useHead({
	title: siteName,
	titleTemplate: (title) => (title === siteName ? title : `${title} - ${siteName}`),
	meta: [
		{ name: 'darkreader-lock', content: '' },
		{ name: 'author', content: 'Satont, me@satont.dev' },
		{ name: 'keywords', content: keywords },
		{ name: 'viewport', content: 'width=device-width, initial-scale=1' },
		{ name: 'description', content: description },
	],
	link: [
		{ rel: 'dns-prefetch', href: 'https://static-cdn.jtvnw.net' },
		{ rel: 'dns-prefetch', href: 'https://avatars.githubusercontent.com' },
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
	ogImage: metaImg,
	twitterCard: 'summary_large_image',
	twitterTitle: siteName,
	twitterDescription: description,
	twitterImage: metaImg,
	keywords,
})
</script>

<template>
	<UiTooltipProvider :delay-duration="100">
		<NuxtLayout>
			<NuxtPage />
		</NuxtLayout>

		<Toaster />
	</UiTooltipProvider>
</template>
