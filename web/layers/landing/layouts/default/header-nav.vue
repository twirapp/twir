<script setup lang="ts">
import { headerLinks } from '~~/layers/landing/layouts/default/anchor-links'

const { t } = useI18n()
</script>

<template>
	<nav class="lg:inline-flex gap-x-2 hidden">
		<NuxtLink
			v-for="anchor of headerLinks"
			:key="anchor.href"
			:to="anchor.href"
			class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 whitespace-nowrap hover:text-[#D5D8DF] rounded-md focus-visible:text-[#D5D8DF] navigation-link"
			prefetch-on="interaction"
		>
			{{ t(anchor.label) }}
		</NuxtLink>
	</nav>
</template>
