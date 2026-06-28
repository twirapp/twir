<script setup lang="ts">
import { ref, watch } from 'vue'

import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'

import { useStorageApi } from '../composables/use-storage-api'

const open = defineModel<boolean>('open', { default: false })

const storageApi = useStorageApi()
const setMutation = storageApi.useMutationStorageSet()

const key = ref('')
const value = ref('')
const error = ref<string | null>(null)
const isSaving = ref(false)

watch(open, (isOpen) => {
	if (isOpen) {
		key.value = ''
		value.value = ''
		error.value = null
	}
})

async function create() {
	error.value = null

	if (!key.value.trim()) {
		error.value = 'Key is required'
		return
	}

	let parsed: unknown
	try {
		parsed = value.value ? JSON.parse(value.value) : null
	} catch {
		error.value = 'Invalid JSON value'
		return
	}

	try {
		isSaving.value = true
		await setMutation.executeMutation({
			key: key.value.trim(),
			value: parsed,
		})
		open.value = false
	} catch (err: any) {
		error.value = err.message || 'Failed to create entry'
	} finally {
		isSaving.value = false
	}
}
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[425px]">
			<DialogHeader>
				<DialogTitle>Add Storage Entry</DialogTitle>
				<DialogDescription>
					Create a new key-value entry in your channel storage.
				</DialogDescription>
			</DialogHeader>

			<form
				class="space-y-4"
				@submit.prevent="create"
			>
				<div class="space-y-2">
					<label class="text-sm font-medium">Key</label>
					<Input
						v-model="key"
						placeholder="myKey"
					/>
				</div>

				<div class="space-y-2">
					<label class="text-sm font-medium">Value (JSON)</label>
					<textarea
						v-model="value"
						placeholder='{"name": "value"} or "string" or 123'
						class="w-full h-32 font-mono text-sm p-3 bg-muted rounded-lg border resize-none focus:outline-none focus:ring-2 focus:ring-ring"
						spellcheck="false"
					/>
				</div>

				<div
					v-if="error"
					class="p-2 bg-destructive/10 text-destructive text-sm rounded"
				>
					{{ error }}
				</div>

				<DialogFooter>
					<Button
						type="button"
						variant="outline"
						@click="open = false"
					>
						Cancel
					</Button>
					<Button
						type="submit"
						:disabled="isSaving"
					>
						Create
					</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
