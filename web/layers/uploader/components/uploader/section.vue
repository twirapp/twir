<script setup lang="ts">
import UploaderRecentUploads from './recent-uploads.vue'
</script>

<template>
	<div class="w-full">
		<div class="flex flex-col gap-3 items-center">
			<UploaderHeader />
			<UploaderDropzone />
			<UploaderRecentUploads />
			<UploaderApiGuide />
		</div>
	</div>
</template>
