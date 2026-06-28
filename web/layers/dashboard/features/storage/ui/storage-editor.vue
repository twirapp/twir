<script setup lang="ts">
import { ref, watch } from 'vue'

import { Button } from '@/components/ui/button'

import type { StorageEntry } from '../composables/use-storage-api'

import { useStorageApi } from '../composables/use-storage-api'
import StorageDeleteDialog from './storage-delete-dialog.vue'

const props = defineProps<{
	entry: StorageEntry
}>()

const emit = defineEmits<{
	close: []
}>()

const storageApi = useStorageApi()
const setMutation = storageApi.useMutationStorageSet()

const editValue = ref('')
const isDeleteOpen = ref(false)
const isSaving = ref(false)
const error = ref<string | null>(null)

watch(
	() => props.entry,
	(entry) => {
		editValue.value = JSON.stringify(entry.value, null, 2)
		error.value = null
	},
	{ immediate: true },
)

function getValueType(value: unknown): string {
	if (value === null || value === undefined) return 'null'
	if (Array.isArray(value)) return 'array'
	return typeof value
}

function formatBytes(bytes: number): string {
	if (bytes === 0) return '0 B'
	const k = 1024
	const sizes = ['B', 'KB', 'MB', 'GB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))
	return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

async function save() {
	error.value = null
	try {
		const parsed = JSON.parse(editValue.value)
		isSaving.value = true
		await setMutation.executeMutation({
			key: props.entry.key,
			value: parsed,
		})
	} catch (err: any) {
		if (err instanceof SyntaxError) {
			error.value = 'Invalid JSON'
		} else {
			error.value = err.message || 'Failed to save'
		}
	} finally {
		isSaving.value = false
	}
}
</script>

<template>
	<div class="flex flex-col h-full">
		<div class="flex items-center justify-between p-4 border-b">
			<div class="flex items-center gap-3">
				<Button
					variant="ghost"
					size="icon"
					@click="emit('close')"
				>
					<Icon name="lucide:arrow-left" class="h-4 w-4" />
				</Button>
				<div>
					<h3 class="font-medium">{{ entry.key }}</h3>
					<p class="text-xs text-muted-foreground">
						{{ getValueType(entry.value) }}
						&middot;
						{{ formatBytes(JSON.stringify(entry.value).length) }}
					</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Button
					variant="destructive"
					size="sm"
					@click="isDeleteOpen = true"
				>
					<Icon name="lucide:trash" class="mr-2 h-4 w-4" />
					Delete
				</Button>
				<Button
					size="sm"
					:disabled="isSaving"
					@click="save"
				>
					<Icon
						v-if="isSaving"
						name="lucide:loader-2"
						class="mr-2 h-4 w-4 animate-spin"
					/>
					<Icon
						v-else
						name="lucide:save"
						class="mr-2 h-4 w-4"
					/>
					Save
				</Button>
			</div>
		</div>

		<div
			v-if="error"
			class="mx-4 mt-2 p-2 bg-destructive/10 text-destructive text-sm rounded"
		>
			{{ error }}
		</div>

		<div class="flex-1 p-4">
			<textarea
				v-model="editValue"
				class="w-full h-full min-h-[300px] font-mono text-sm p-3 bg-muted rounded-lg border resize-none focus:outline-none focus:ring-2 focus:ring-ring"
				spellcheck="false"
			/>
		</div>

		<StorageDeleteDialog
			v-model:open="isDeleteOpen"
			:entry="entry"
			@deleted="emit('close')"
		/>
	</div>
</template>
