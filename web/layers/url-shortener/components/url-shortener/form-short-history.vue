<script setup lang="ts">
import { useUrlShortener } from '#layers/url-shortener/composables/use-url-shortener'

const api = useUrlShortener()
const recentUrlsError = ref()

onMounted(async () => {
	if (!import.meta.client) return

	const response = await api.refetchLatestShortenedUrls({ page: 0, perPage: 3 })
	if (response.error) {
		recentUrlsError.value = response.error
		return
	}
})
</script>

<template>
	<div class="flex flex-col w-full max-w-xl gap-3">
		<template v-if="recentUrlsError">
			{{ recentUrlsError }}
		</template>
		<template v-else>
			<TransitionGroup name="list" tag="div" class="flex flex-col gap-3">
				<div v-for="url in api.latestShortenedUrls" :key="url.id">
					<UrlShortenerCard :url="url" />
				</div>
			</TransitionGroup>
		</template>
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
