<script setup lang="ts">
// eslint-disable typescript-eslint(consistent-type-imports)
import { ShortUrlProfileParamsSortByEnum } from '@twir/api/openapi'
import UrlShortenerLinkCardWithStats from '../../components/url-shortener/link-card-with-stats.vue'
import UrlShortenerLinksPagination from '../../components/url-shortener/links-pagination.vue'

definePageMeta({
	layout: 'landing',
})

const urlShortener = useUrlShortener()
const pagination = useLinksPagination()

const { data: linksData, refresh } = await useAsyncData(
	async () => {
		const result = await urlShortener.refetchLatestShortenedUrls({
			page: pagination.currentPage.value,
			perPage: pagination.perPage.value,
			sortBy: pagination.sortBy.value,
		})
		return result.data
	},
	{
		watch: [pagination.currentPage, pagination.perPage, pagination.sortBy],
	}
)

const links = computed(() => linksData.value?.items ?? [])

// Update total when data changes
watch(
	() => linksData.value?.total,
	(newTotal) => {
		if (newTotal !== undefined) {
			pagination.setTotal(newTotal)
		}
	},
	{ immediate: true }
)

function handleLinkUpdated() {
	refresh()
}

function handleLinkDeleted() {
	refresh()
}

function handleSortChange(value: string) {
	pagination.setSortBy(value as ShortUrlProfileParamsSortByEnum)
}
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
						<Icon
							name="lucide:arrow-left"
							class="w-4 h-4"
						/>
						Back to Shortener
					</NuxtLink>

					<div class="flex items-center justify-between flex-wrap gap-4">
						<div>
							<h1 class="text-3xl font-bold">My Short Links</h1>
							<p class="text-[hsl(240,11%,65%)]">
								Manage your short links and view their statistics
							</p>
						</div>

						<!-- Sort selector -->
						<div class="flex items-center gap-2">
							<span class="text-sm text-[hsl(240,11%,65%)]">Sort by:</span>
							<select
								:value="pagination.sortBy.value"
								@change="handleSortChange(($event.target as HTMLSelectElement).value)"
								class="px-3 py-2 rounded-lg border border-[hsl(240,11%,30%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)] transition-colors text-sm"
							>
								<option :value="ShortUrlProfileParamsSortByEnum.Views">Most Views</option>
								<option :value="ShortUrlProfileParamsSortByEnum.CreatedAt">Most Recent</option>
							</select>
						</div>
					</div>
				</div>

				<div
					v-if="links.length === 0"
					class="text-center py-12"
				>
					<p class="text-[hsl(240,11%,65%)]">You don't have any short links yet</p>
					<NuxtLink
						to="/url-shortener"
						class="text-[hsl(240,11%,85%)] hover:text-white transition-colors underline"
					>
						Create your first link
					</NuxtLink>
				</div>

				<div
					v-else
					class="space-y-6"
				>
					<!-- Top Pagination -->
					<UrlShortenerLinksPagination
						v-if="pagination.totalPages.value > 1"
						:current-page="pagination.currentPage.value"
						:total-pages="pagination.totalPages.value"
						:total="pagination.total.value"
						@go-to-page="pagination.goToPage"
						@next-page="pagination.nextPage"
						@previous-page="pagination.previousPage"
					/>

					<!-- Links Grid -->
					<div class="grid gap-4">
						<UrlShortenerLinkCardWithStats
							v-for="link in links"
							:key="link.id"
							:link="link"
							@updated="handleLinkUpdated"
							@deleted="handleLinkDeleted"
						/>
					</div>

					<!-- Bottom Pagination -->
					<UrlShortenerLinksPagination
						v-if="pagination.totalPages.value > 1"
						:current-page="pagination.currentPage.value"
						:total-pages="pagination.totalPages.value"
						:total="pagination.total.value"
						@go-to-page="pagination.goToPage"
						@next-page="pagination.nextPage"
						@previous-page="pagination.previousPage"
					/>
				</div>
			</div>
		</div>
	</div>
</template>
