<script setup lang="ts">
import { useRoute } from 'vue-router'

import { anchorLinks } from '~~/layers/landing/layouts/default/anchor-links'

const route = useRoute()
const isHomePage = computed(() => route.path === '/')
const isTermsPage = computed(() => route.path === '/terms')
</script>

<template>
	<nav class="md:inline-flex gap-x-2 hidden">
		<template v-if="isHomePage">
			<a
				v-for="anchor of anchorLinks"
				:key="anchor.href"
				:href="anchor.href"
				class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] rounded-md focus-visible:text-[#D5D8DF] navigation-link"
			>
				{{ anchor.label }}
			</a>
			<NuxtLink
				to="/h"
				class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] rounded-md focus-visible:text-[#D5D8DF] navigation-link"
				prefetch-on="interaction"
			>
				Hastebin
			</NuxtLink>
		</template>
		<NuxtLink
			v-if="!isTermsPage"
			to="/terms"
			class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] rounded-md focus-visible:text-[#D5D8DF] navigation-link"
			prefetch-on="interaction"
		>
			Terms of Service
		</NuxtLink>
	</nav>
</template>
