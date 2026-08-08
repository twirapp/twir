<script setup lang="ts">
import UploaderUploadPreviewCard from './upload-preview-card.vue'

const uploader = useUploader()
const recentUploadsError = ref<string | null>(null)

const response = await uploader.refetchUploads({
	page: 0,
	perPage: 3,
})
if (response.error) {
	recentUploadsError.value = response.error
}
</script>

<template>
	<div class="flex flex-col w-full max-w-xl gap-3">
		<div class="flex items-center justify-between">
			<h3 class="text-sm font-semibold text-[hsl(240,11%,90%)]">{{ $t('uploader.recent.title') }}</h3>
			<NuxtLink
				to="/uploader/profile"
				class="text-xs text-[hsl(240,11%,65%)] hover:text-white transition-colors"
			>
				{{ $t('uploader.recent.viewAll') }}
			</NuxtLink>
		</div>
		<p v-if="recentUploadsError" class="text-sm text-red-400">
			{{ recentUploadsError }}
		</p>
		<p v-else-if="uploader.latestUploads.length === 0" class="text-sm text-[hsl(240,11%,65%)]">
			{{ $t('uploader.recent.empty') }}
		</p>
		<TransitionGroup v-else name="list" tag="div" class="grid gap-3 sm:grid-cols-3">
			<div v-for="upload in uploader.latestUploads" :key="upload.id">
				<UploaderUploadPreviewCard :upload="upload" />
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
