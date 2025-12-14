<script setup lang="ts">
definePageMeta({
	layout: 'landing',
})
</script>

<template>
	<div class="h-full w-full">
		<div
			class="absolute inset-0 bg-[linear-gradient(to_right,hsl(240,11%,9%)_1px,transparent_1px),linear-gradient(to_bottom,hsl(240,11%,9%)_1px,transparent_1px)] bg-size-[36px_36px] mask-[linear-gradient(to_bottom,transparent_15%,black_100%)]"
		></div>
		<div
			class="container mx-auto flex flex-col items-center justify-center min-h-[calc(100vh-73px)] relative"
		>
			<UrlShortenerSection />
		</div>
	</div>
</template>
