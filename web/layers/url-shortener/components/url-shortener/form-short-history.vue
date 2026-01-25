<script setup lang="ts">
import { useUrlShortener } from '#layers/url-shortener/composables/use-url-shortener'
import UrlShortenerHistoryCard from '#layers/url-shortener/components/url-shortener/short-history-card.vue'

const api = useUrlShortener()
const recentUrlsError = ref<string | null>(null)

onMounted(async () => {
	if (!import.meta.client) return

	const response = await api.refetchLatestShortenedUrls({
		page: 0,
		perPage: 3,
		sortBy: 'created_at',
	})
	if (response.error) {
		recentUrlsError.value = response.error
		return
	}
})
</script>

<template>
	<div class="flex flex-col w-full max-w-xl gap-3">
		<div class="flex items-center justify-between">
			<h3 class="text-sm font-semibold">Recent links</h3>
			<NuxtLink
				to="/url-shortener/profile"
				class="text-xs text-[hsl(240,11%,65%)] hover:text-white transition-colors"
			>
				View all
			</NuxtLink>
		</div>
		<p v-if="recentUrlsError" class="text-sm text-red-400">
			{{ recentUrlsError }}
		</p>
		<p v-else-if="api.latestShortenedUrls.length === 0" class="text-sm text-[hsl(240,11%,65%)]">
			No recent links yet. Create your first short link above.
		</p>
		<TransitionGroup v-else name="list" tag="div" class="flex flex-col gap-3">
			<div v-for="url in api.latestShortenedUrls" :key="url.id">
				<UrlShortenerHistoryCard :url="url" />
			</div>
		</TransitionGroup>
	</div>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
	transition: all 0.3s ease;
}

.list-enter-from {
	opacity: 0;
	transform: translateY(-20px);
}

.list-leave-to {
	opacity: 0;
	transform: translateY(20px);
}

.list-move {
	transition: transform 0.3s ease;
}
</style>
