<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
</script>

<template>
	<div class="flex items-center gap-2">
		<UiButton disabled size="icon">
			<PencilIcon class="h-4 w-4" />
		</UiButton>
		<UiButton disabled variant="destructive" size="icon">
			<TrashIcon class="h-4 w-4" />
		</UiButton>
	</div>
</template>
