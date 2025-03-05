<script setup lang="ts">
import Nav from './header-nav.vue'
import HeaderProfile from './header-profile.vue'

import TwirLogo from '~/components/twir-logo.vue'

onMounted(() => {
	if (import.meta.server) return

	const el = document.querySelector('header')
	const observer = new IntersectionObserver(
		([e]) => {
			e.target.classList.toggle('sticky-header', e.intersectionRatio < 1)
		},
		{ threshold: [1] },
	)

	observer.observe(el!)
})
</script>

<template>
	<header
		id="top"
		class="bg-[#09090B]/25 sm:px-8 px-4 sticky header top-[-1px] z-50"
	>
		<div class="lg:container">
			<div
				class="flex justify-between items-center border-b border-[#9DA2BA]/30 py-4"
			>
				<div class="flex items-center gap-[4px] divide-x-2">
					<NuxtLink to="/" class="flex items-center gap-[12px] cursor-pointer">
						<TwirLogo :src="TwirLogo" alt="Twir" class="w-9 h-9" />
						<span class="text-2xl font-semibold text-white">Twir</span>
					</NuxtLink>

					<div class="bg-red inline-block h-full"></div>

					<span class="pl-2 text-sm text-white/70 hidden sm:inline-block">
						Hosted by
						<a
							href="https://avoro.eu/cp/aff.php?aff=82"
							target="_blank"
							class="underline"
						>avoro.eu</a>
					</span>
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
