<script setup lang="ts">
const { t } = useI18n()

const faqKeys = ['q1', 'q2', 'q3', 'q4', 'q5'] as const

useSchemaOrg(
	faqKeys.map((key) =>
		defineQuestion({
			name: () => t(`landing.faq.${key}.question`),
			acceptedAnswer: () => t(`landing.faq.${key}.answer`),
		})
	)
)
</script>

<template>
	<section id="faq" aria-labelledby="faq-heading">
		<div class="container mx-auto flex flex-col items-center px-5 pb-24 md:px-8 text-center">
			<h2
				id="faq-heading"
				class="font-semibold text-[#B0ADFF] text-base uppercase tracking-wider mb-3"
			>
				{{ t('landing.faq.badge') }}
			</h2>
			<p class="font-bold text-white text-4xl tracking-tight leading-tight">
				{{ t('landing.faq.title') }}
			</p>
			<UiAccordion type="single" collapsible class="mt-[64px] w-full max-w-3xl text-left">
				<UiAccordionItem
					v-for="key of faqKeys"
					:key="key"
					:value="key"
					class="border-[#72757d26]"
				>
					<UiAccordionTrigger class="text-left text-base text-white hover:no-underline">
						{{ t(`landing.faq.${key}.question`) }}
					</UiAccordionTrigger>
					<UiAccordionContent class="text-base leading-relaxed text-[#ADB0B8]">
						{{ t(`landing.faq.${key}.answer`) }}
					</UiAccordionContent>
				</UiAccordionItem>
			</UiAccordion>
		</div>
	</section>
</template>
