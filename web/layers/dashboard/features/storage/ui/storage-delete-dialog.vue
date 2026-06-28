<script setup lang="ts">
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'

import type { StorageEntry } from '../composables/use-storage-api'

import { useStorageApi } from '../composables/use-storage-api'

const props = defineProps<{
	entry?: StorageEntry | null
	deleteAll?: boolean
}>()

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
	deleted: []
}>()

const storageApi = useStorageApi()
const deleteMutation = storageApi.useMutationStorageDelete()
const deleteAllMutation = storageApi.useMutationStorageDeleteAll()

async function confirmDelete() {
	if (props.deleteAll) {
		await deleteAllMutation.executeMutation({})
	} else if (props.entry) {
		await deleteMutation.executeMutation({ key: props.entry.key })
	}
	emit('deleted')
	open.value = false
}
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>
					{{ deleteAll ? 'Clear All Storage' : 'Delete Entry' }}
				</DialogTitle>
				<DialogDescription>
					<template v-if="deleteAll">
						This will permanently delete all storage entries for this channel. This action cannot be undone.
					</template>
					<template v-else>
						Are you sure you want to delete
						<code class="bg-muted px-1 py-0.5 rounded text-xs">{{ entry?.key }}</code>?
						This cannot be undone.
					</template>
				</DialogDescription>
			</DialogHeader>

			<DialogFooter>
				<Button
					variant="outline"
					@click="open = false"
				>
					Cancel
				</Button>
				<Button
					variant="destructive"
					@click="confirmDelete"
				>
					{{ deleteAll ? 'Clear All' : 'Delete' }}
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
