<script setup lang="ts">
import { PencilIcon, TrashIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import ActionConfirm from '@/components/ui/action-confirm.vue'
import { Button } from '@/components/ui/button'

const emits = defineEmits<{
	(event: 'delete'): void
	(event: 'edit'): void
}>()

const showDelete = ref(false)
</script>

<template>
	<div class="flex items-center gap-2">
		<Button variant="secondary" size="icon" @click="emits('edit')">
			<PencilIcon class="h-4 w-4" />
		</Button>
		<Button variant="destructive" size="icon" @click="showDelete = true">
			<TrashIcon class="h-4 w-4" />
		</Button>
	</div>

	<ActionConfirm
		v-model:open="showDelete"
		@confirm="emits('delete')"
	/>
</template>
