<script setup lang="ts">
import { computed, ref } from 'vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'

import type { StorageEntry } from './composables/use-storage-api'

import { useStorageApi } from './composables/use-storage-api'
import StorageCreateDialog from './ui/storage-create-dialog.vue'
import StorageDeleteDialog from './ui/storage-delete-dialog.vue'
import StorageEditor from './ui/storage-editor.vue'

const storageApi = useStorageApi()
const { entries, isLoading, totalSize, refresh } = storageApi

const searchQuery = ref('')
const selectedKey = ref<string | null>(null)
const isCreateOpen = ref(false)
const isDeleteAllOpen = ref(false)

const filteredEntries = computed(() => {
	if (!searchQuery.value) return entries.value
	const q = searchQuery.value.toLowerCase()
	return entries.value.filter((e) => e.key.toLowerCase().includes(q))
})

const selectedEntry = computed(() => {
	if (!selectedKey.value) return null
	return entries.value.find((e) => e.key === selectedKey.value) ?? null
})

function getValueType(value: unknown): string {
	if (value === null || value === undefined) return 'null'
	if (Array.isArray(value)) return 'array'
	return typeof value
}

function getValuePreview(value: unknown): string {
	if (value === null || value === undefined) return 'null'
	if (typeof value === 'string') return value.length > 50 ? `${value.slice(0, 50)}...` : value
	if (typeof value === 'number' || typeof value === 'boolean') return String(value)
	if (Array.isArray(value)) return `Array(${value.length})`
	if (typeof value === 'object') return `Object(${Object.keys(value as object).length} keys)`
	return String(value)
}

function formatBytes(bytes: number): string {
	if (bytes === 0) return '0 B'
	const k = 1024
	const sizes = ['B', 'KB', 'MB', 'GB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))
	return `${Number.parseFloat((bytes / k ** i).toFixed(1))} ${sizes[i]}`
}

function selectKey(key: string) {
	selectedKey.value = key
}
</script>

<template>
	<div class="flex h-full flex-col gap-4">
		<div class="flex items-center justify-between">
			<p class="text-muted-foreground text-sm">
				Store data accessible in scripts via
				<code class="bg-muted rounded px-1 py-0.5 text-xs">twir.storage.get('key')</code>
			</p>
			<div class="flex gap-2">
				<Button
					variant="outline"
					:disabled="isLoading"
					@click="refresh"
				>
					<Icon
						name="lucide:refresh-cw"
						class="mr-2 h-4 w-4"
						:class="{ 'animate-spin': isLoading }"
					/>
					Refresh
				</Button>
				<Button
					variant="destructive"
					:disabled="entries.length === 0"
					@click="isDeleteAllOpen = true"
				>
					<Icon
						name="lucide:trash-2"
						class="mr-2 h-4 w-4"
					/>
					Clear All
				</Button>
				<Button @click="isCreateOpen = true">
					<Icon
						name="lucide:plus"
						class="mr-2 h-4 w-4"
					/>
					Add Entry
				</Button>
			</div>
		</div>

		<div
			v-if="isLoading"
			class="flex items-center justify-center py-12"
		>
			<Icon
				name="lucide:loader-2"
				class="text-muted-foreground h-6 w-6 animate-spin"
			/>
		</div>

		<div
			v-else
			class="flex min-h-[60dvh] overflow-hidden rounded-lg border"
		>
			<div class="flex w-64 flex-col border-r">
				<div class="border-b p-3">
					<Input
						v-model="searchQuery"
						placeholder="Search keys..."
						class="h-8"
					/>
				</div>
				<div class="flex-1 overflow-y-auto">
					<div
						v-if="filteredEntries.length === 0"
						class="text-muted-foreground p-4 text-center text-sm"
					>
						{{ entries.length === 0 ? 'No entries yet' : 'No matches' }}
					</div>
					<button
						v-for="entry in filteredEntries"
						:key="entry.key"
						class="hover:bg-accent flex w-full items-center justify-between gap-2 px-3 py-2 text-left transition-colors"
						:class="{ 'bg-accent': selectedKey === entry.key }"
						@click="selectKey(entry.key)"
					>
						<span class="truncate text-sm">{{ entry.key }}</span>
						<span
							class="bg-muted text-muted-foreground shrink-0 rounded-full px-1.5 py-0.5 text-[10px]"
						>
							{{ getValueType(entry.value) }}
						</span>
					</button>
				</div>
			</div>

			<div class="flex flex-1 flex-col">
				<div
					v-if="!selectedEntry"
					class="text-muted-foreground flex flex-1 items-center justify-center text-sm"
				>
					Select a key to view its value
				</div>
				<StorageEditor
					v-else
					:entry="selectedEntry"
					@close="selectedKey = null"
				/>
			</div>
		</div>

		<div class="text-muted-foreground flex items-center justify-between text-xs">
			<span>{{ entries.length }} entries</span>
			<span>{{ formatBytes(totalSize) }} / 30 MB</span>
		</div>

		<StorageCreateDialog v-model:open="isCreateOpen" />

		<StorageDeleteDialog
			v-model:open="isDeleteAllOpen"
			:delete-all="true"
		/>
	</div>
</template>
