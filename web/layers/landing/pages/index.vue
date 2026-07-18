<script setup lang="ts">
import Faq from '~~/layers/landing/components/index/faq/faq.vue'
import Features from '~~/layers/landing/components/index/features/features.vue'
import Hero from '~~/layers/landing/components/index/hero/hero.vue'
import Integrations from '~~/layers/landing/components/index/integrations/integrations.vue'
import Stats from '~~/layers/landing/components/index/stats/stats.vue'

definePageMeta({
	layout: 'landing',
})

useAppSeo({
	titleKey: 'landing.meta.title',
	descriptionKey: 'landing.meta.description',
})

useSchemaOrg([
	defineSoftwareApp({
		name: 'Twir',
		applicationCategory: 'MultimediaApplication',
		operatingSystem: 'Web',
		offers: defineOffer({
			price: 0,
			priceCurrency: 'USD',
		}),
	}),
])
</script>

<template>
	<div class="text-white">
		<Hero />
		<Stats />
		<Features />
		<Integrations />
		<Faq />
		<!--		<Team /> -->
	</div>
</template>
