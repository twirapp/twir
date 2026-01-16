<script setup lang="ts">
import UrlShortenerLinkCardWithStats from "../../components/url-shortener/link-card-with-stats.vue";

definePageMeta({
	layout: "landing",
});

const urlShortener = useUrlShortener();

const { data: linksData, refresh } = await useAsyncData(async () => {
	const result = await urlShortener.refetchLatestShortenedUrls({ page: 0, perPage: 10 });
	return result.data;
});

const links = computed(() => linksData.value?.items ?? []);
</script>

<template>
	<div class="h-full w-full">
		<div
			class="absolute inset-0 bg-[linear-gradient(to_right,hsl(240,11%,9%)_1px,transparent_1px),linear-gradient(to_bottom,hsl(240,11%,9%)_1px,transparent_1px)] bg-size-[36px_36px] mask-[linear-gradient(to_bottom,transparent_15%,black_100%)]"
		></div>
		<div class="container mx-auto py-8 relative min-h-[calc(100vh-73px)]">
			<div class="max-w-6xl mx-auto space-y-8">
				<div class="space-y-4">
					<NuxtLink
						to="/url-shortener"
						class="inline-flex items-center gap-2 px-4 py-2 rounded-lg border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] transition-colors text-sm font-medium"
					>
						<Icon name="lucide:arrow-left" class="w-4 h-4" />
						Back to Shortener
					</NuxtLink>
					<h1 class="text-3xl font-bold">My Short Links</h1>
					<p class="text-[hsl(240,11%,65%)]">Manage your short links and view their statistics</p>
				</div>

				<div v-if="links.length === 0" class="text-center py-12">
					<p class="text-[hsl(240,11%,65%)]">You don't have any short links yet</p>
					<NuxtLink
						to="/url-shortener"
						class="text-[hsl(240,11%,85%)] hover:text-white transition-colors underline"
					>
						Create your first link
					</NuxtLink>
				</div>

				<div v-else class="grid gap-4">
					<UrlShortenerLinkCardWithStats v-for="link in links" :key="link.id" :link="link" />
				</div>
			</div>
		</div>
	</div>
</template>
