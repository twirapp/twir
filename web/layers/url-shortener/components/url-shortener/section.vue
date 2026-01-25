<script setup lang="ts">
import FormShortHistory from '#layers/url-shortener/components/url-shortener/form-short-history.vue'
</script>

<template>
	<div class="w-full">
		<div class="flex flex-col gap-3 items-center">
			<UrlShortenerHeader />
			<UrlShortenerForm />
			<FormShortHistory />
			<UrlShortenerCustomDomainPanel />
		</div>
	</div>
</template>
