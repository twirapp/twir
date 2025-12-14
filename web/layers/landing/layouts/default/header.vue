<script setup lang="ts">
import Nav from './header-nav.vue'
import HeaderProfile from './header-profile.vue'

import TwirLogo from '~/components/twir-logo.vue'

onMounted(() => {
	if (import.meta.server) return

	const el = document.querySelector('header')
	const observer = new IntersectionObserver(
		([e]) => {
			e?.target.classList.toggle('sticky-header', e?.intersectionRatio < 1)
		},
		{ threshold: [1] }
	)

	observer.observe(el!)
})

const title = `Twir${import.meta.dev ? ' dev' : ''}`
</script>

<template>
	<header id="top" class="bg-[#09090B]/25 sm:px-8 px-4 sticky header -top-px z-50">
		<div class="lg:container lg:mx-auto">
			<div class="flex justify-between items-center py-4">
				<div class="flex items-center gap-1 divide-x-2">
					<NuxtLink to="/" class="flex items-center gap-3 cursor-pointer">
						<TwirLogo :src="TwirLogo" alt="Twir" class="w-9 h-9" />
						<span class="text-2xl font-semibold text-white">{{ title }}</span>
					</NuxtLink>
				</div>

				<Nav />

				<HeaderProfile />
			</div>
		</div>
	</header>
</template>

<style>
.sticky-header {
	background-color: #17171a !important;
}
</style>
