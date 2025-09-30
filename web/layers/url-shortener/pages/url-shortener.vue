<script setup lang="ts">
import UrlShortenerForm from '~~/layers/url-shortener/components/url-shortener-form.vue'
import Particles from '~/components/ui/bits/backgrounds/Particles/Particles.vue'

definePageMeta({
	layout: 'landing',
})
</script>

<template>
	<div class="h-full w-full">
		<div class="container flex items-center justify-center h-screen relative">
			<UrlShortenerForm />
		</div>
		<Particles
			:particle-count="200"
			:particle-spread="10"
			:speed="0.05"
			:particle-colors="['#ffffff']"
			:move-particles-on-hover="false"
			:particle-hover-factor="1"
			:alpha-particles="true"
			:particle-base-size="100"
			:size-randomness="1"
			:camera-distance="20"
			:disable-rotation="true"
			class="w-full h-full pointer-events-none"
		/>
	</div>
</template>
