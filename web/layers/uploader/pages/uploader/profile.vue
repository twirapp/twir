<script setup lang="ts">
import UploaderUploadPreviewCard from '../../components/uploader/upload-preview-card.vue'
import UploaderUploadsPagination from '../../components/uploader/uploads-pagination.vue'

definePageMeta({
	layout: 'landing',
	sitemap: false,
	robots: {
		noindex: true,
		nofollow: true,
	},
})

useSeoMeta({
	robots: 'noindex, nofollow',
})

const uploader = useUploader()
const pagination = useUploadsPagination()

const { data: uploadsData, refresh } = await useAsyncData(
	'uploader-files',
	async () => {
		const result = await uploader.refetchUploads({
			page: pagination.currentPage.value,
			perPage: pagination.perPage.value,
		})
		return result.data
	},
	{
		watch: [pagination.currentPage, pagination.perPage],
	}
)

const uploads = computed(() => uploadsData.value?.files ?? [])

// Update total when data changes
watch(
	() => uploadsData.value?.total,
	(newTotal) => {
		if (newTotal !== undefined) {
			pagination.setTotal(newTotal)
		}
	},
	{ immediate: true }
)

function handleUploadDeleted() {
	refresh()
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
						to="/uploader"
						class="inline-flex items-center gap-2 px-4 py-2 rounded-lg border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] transition-colors text-sm font-medium"
					>
						<Icon
							name="lucide:arrow-left"
							class="w-4 h-4"
						/>
						{{ $t('uploader.profile.back') }}
					</NuxtLink>

					<div>
						<h1 class="text-3xl font-bold">{{ $t('uploader.profile.title') }}</h1>
						<p class="text-[hsl(240,11%,65%)]">
							{{ $t('uploader.profile.description') }}
						</p>
					</div>
				</div>

				<div
					v-if="uploads.length === 0"
					class="text-center py-12"
				>
					<p class="text-[hsl(240,11%,65%)]">{{ $t('uploader.profile.empty') }}</p>
					<NuxtLink
						to="/uploader"
						class="text-[hsl(240,11%,85%)] hover:text-white transition-colors underline"
					>
						{{ $t('uploader.profile.createFirst') }}
					</NuxtLink>
				</div>

				<div
					v-else
					class="space-y-6"
				>
					<!-- Top Pagination -->
					<UploaderUploadsPagination
						v-if="pagination.totalPages.value > 1"
						:current-page="pagination.currentPage.value"
						:total-pages="pagination.totalPages.value"
						:total="pagination.total.value"
						@go-to-page="pagination.goToPage"
						@next-page="pagination.nextPage"
						@previous-page="pagination.previousPage"
					/>

					<!-- Uploads Grid -->
					<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
						<UploaderUploadPreviewCard
							v-for="upload in uploads"
							:key="upload.id"
							:upload="upload"
							@deleted="handleUploadDeleted"
						/>
					</div>

					<!-- Bottom Pagination -->
					<UploaderUploadsPagination
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
