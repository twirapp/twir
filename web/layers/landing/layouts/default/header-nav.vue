<script setup lang="ts">
import { anchorLinks } from '~~/layers/landing/layouts/default/anchor-links'
</script>

<template>
	<nav class="md:inline-flex gap-x-2 hidden">
		<a
			v-for="anchor of anchorLinks"
			:key="anchor.href"
			:href="anchor.href"
			class="text-[#ADB0B8] font-medium px-3 py-2 leading-6 hover:text-[#D5D8DF] rounded-md focus-visible:text-[#D5D8DF] navigation-link"
		>
			{{ anchor.label }}
		</a>
	</nav>
</template>
