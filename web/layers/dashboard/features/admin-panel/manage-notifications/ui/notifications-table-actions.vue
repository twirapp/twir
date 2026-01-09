<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'



const emits = defineEmits<{
	(event: 'delete'): void
	(event: 'edit'): void
}>()

const showDelete = ref(false)
</script>

<template>
	<div class="flex items-center gap-2">
		<UiButton variant="secondary" size="icon" @click="emits('edit')">
			<PencilIcon class="h-4 w-4" />
		</UiButton>
		<UiButton variant="destructive" size="icon" @click="showDelete = true">
			<TrashIcon class="h-4 w-4" />
		</UiButton>
	</div>

	<UiActionConfirm
		v-model:open="showDelete"
		@confirm="emits('delete')"
	/>
</template>
