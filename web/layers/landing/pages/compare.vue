<script setup lang="ts">
import { DISCORD_INVITE_URL } from '@twir/brand'

import CompareTable from '../components/compare/compare-table.vue'
import { compareFaqKeys } from '../components/compare/compare.data'
import LandingUiButton from '../components/landing-ui-button.vue'

definePageMeta({
	layout: 'landing',
})

const { t } = useI18n()
const localePath = useLocalePath()

useAppSeo({
	titleKey: 'compare.meta.title',
	descriptionKey: 'compare.meta.description',
})

useSchemaOrg([
	defineWebPage({
		'@type': ['WebPage', 'FAQPage'],
		'name': () => t('compare.meta.title'),
		'description': () => t('compare.meta.description'),
	}),
	...compareFaqKeys.map((key) =>
		defineQuestion({
			name: () => t(`compare.faq.${key}.question`),
			acceptedAnswer: () => t(`compare.faq.${key}.answer`),
		})
	),
	defineBreadcrumb({
		itemListElement: [
			{ name: 'Twir', item: '/' },
			{ name: () => t('compare.breadcrumb'), item: '/compare' },
		],
	}),
])
</script>

<template>
	<div class="text-white">
		<section class="container mx-auto flex flex-col items-center px-5 pt-24 pb-12 text-center md:px-8">
			<span
				class="mb-4 rounded-full border border-[#72757d26] bg-[#1a1a22] px-4 py-1.5 text-sm font-medium text-[#B0ADFF]"
			>
				{{ t('compare.hero.badge') }}
			</span>
			<h1 class="max-w-3xl text-4xl font-bold tracking-tight sm:text-5xl">
				{{ t('compare.hero.title') }}
			</h1>
			<p class="mt-5 max-w-2xl text-lg leading-relaxed text-[#ADB0B8]">
				{{ t('compare.hero.subtitle') }}
			</p>
		</section>

		<section class="container mx-auto px-5 pb-16 md:px-8" aria-label="Comparison table">
			<CompareTable />
		</section>

		<section class="container mx-auto flex flex-col items-center px-5 pb-24 text-center md:px-8">
			<h2 class="text-3xl font-bold tracking-tight">
				{{ t('compare.cta.title') }}
			</h2>
			<p class="mt-4 max-w-2xl text-lg leading-relaxed text-[#ADB0B8]">
				{{ t('compare.cta.subtitle') }}
			</p>
			<div class="mt-8 flex flex-wrap items-center justify-center gap-4">
				<LandingUiButton variant="primary" :href="localePath('/login')">
					{{ t('compare.cta.getStarted') }}
				</LandingUiButton>
				<LandingUiButton variant="secondary" :href="localePath('/dashboard/import')">
					{{ t('compare.cta.import') }}
				</LandingUiButton>
			</div>
			<p class="mt-6 text-sm text-[#ADB0B8]">
				{{ t('compare.cta.questions') }}
				<a
					:href="DISCORD_INVITE_URL"
					target="_blank"
					rel="noopener noreferrer"
					class="text-[#B0ADFF] underline-offset-4 hover:underline"
				>
					Discord
				</a>
			</p>
		</section>

		<section
			class="container mx-auto max-w-3xl px-5 pb-24 md:px-8"
			aria-labelledby="compare-faq-heading"
		>
			<h2 id="compare-faq-heading" class="text-center text-3xl font-bold tracking-tight">
				{{ t('compare.faq.title') }}
			</h2>
			<UiAccordion type="single" collapsible class="mt-8">
				<UiAccordionItem
					v-for="key of compareFaqKeys"
					:key="key"
					:value="key"
					class="border-[#72757d26]"
				>
					<UiAccordionTrigger class="text-left text-base text-white hover:no-underline">
						{{ t(`compare.faq.${key}.question`) }}
					</UiAccordionTrigger>
					<UiAccordionContent class="text-base leading-relaxed text-[#ADB0B8]">
						{{ t(`compare.faq.${key}.answer`) }}
					</UiAccordionContent>
				</UiAccordionItem>
			</UiAccordion>
		</section>

		<section class="container mx-auto max-w-3xl px-5 pb-24 md:px-8">
			<p class="text-center text-xs leading-relaxed text-[#6b6e75]">
				{{ t('compare.disclaimer') }}
			</p>
		</section>
	</div>
</template>
