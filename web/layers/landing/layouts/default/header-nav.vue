<script setup lang="ts">
import { useRoute } from 'vue-router'

import { anchorLinks } from '~~/layers/landing/layouts/default/anchor-links'

const route = useRoute()
</script>

<template>
	<nav class="md:inline-flex gap-x-2 hidden">
		<NuxtLink
			v-for="anchor of anchorLinks"
			:key="anchor.href"
			:to="anchor.href"
			class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] rounded-md focus-visible:text-[#D5D8DF] navigation-link"
			prefetch-on="interaction"
		>
			{{ anchor.label }}
		</NuxtLink>
	</nav>
</template>
