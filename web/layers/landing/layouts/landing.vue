<script lang="ts" setup>
import Footer from './default/footer.vue'
import Header from './default/header.vue'
</script>

<template>
	<span
		class="absolute content-[''] h-[482px] rounded-full pointer-events-none purple-gradient -z-20 -top-[220px] right-0 left-0 mx-auto"
	></span>
	<Header />

	<slot />

	<Footer />
</template>

<style>
html {
	scroll-behavior: smooth;
}

body {
	@apply bg-[#09090B] overflow-x-hidden
}
</style>
