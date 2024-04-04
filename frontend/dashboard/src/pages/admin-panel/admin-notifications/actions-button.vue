<script setup lang="ts">
import { TrashIcon, PencilIcon } from 'lucide-vue-next';
import { ref } from 'vue';

import { Button } from '@/components/ui/button';
import DeleteConfirm from '@/components/ui/delete-confirm.vue';

const emits = defineEmits<{
	(event: 'delete'): void
	(event: 'edit'): void
}>();

const showDelete = ref(false);
</script>

<template>
	<div class="flex items-center gap-2">
		<Button variant="destructive" size="icon" @click="emits('delete')">
			<TrashIcon class="h-4 w-4" />
		</Button>
		<Button variant="secondary" size="icon" @click="showDelete = true">
			<PencilIcon class="h-4 w-4" />
		</Button>
	</div>

	<DeleteConfirm
		v-model:open="showDelete"
		@confirm="emits('edit')"
	/>
</template>
