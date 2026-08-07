<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Dialog, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'
import { Input } from '@/components/ui/input'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

import { useSevenTvEmotes } from './composables/useSevenTvEmotes'
import type { SelectedSevenTvEmote } from './composables/useSevenTvEmotes'

const open = defineModel<boolean>('open', { default: false })

const emit = defineEmits<{
	select: [emote: SelectedSevenTvEmote]
}>()

const {
	activeTab,
	searchQuery,
	currentStatus,
	currentError,
	filteredEmotes,
	emptyMessage,
	tabValues,
	retryActiveTab,
	handleSelectEmote,
	getEmoteUrl,
} = useSevenTvEmotes(open, (emote) => emit('select', emote))
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger v-if="$slots.trigger" as-child><slot name="trigger" /></DialogTrigger>
		<DialogOrSheet class="max-h-[80dvh] max-w-3xl gap-0 overflow-hidden rounded-t-2xl p-0 sm:rounded-2xl">
			<div class="flex min-h-0 flex-col">
				<DialogHeader class="border-b px-6 py-4">
					<DialogTitle>Эмоции</DialogTitle>
					<DialogDescription class="sr-only">Выберите эмоцию из наборов 7TV.</DialogDescription>
				</DialogHeader>

				<Tabs v-model="activeTab" class="min-h-0 gap-0">
					<div class="sticky top-0 z-10 border-b bg-background/95 px-4 py-3 backdrop-blur">
						<TabsList class="grid h-9 w-full grid-cols-2 bg-muted/60">
							<TabsTrigger value="global">7TV</TabsTrigger>
							<TabsTrigger value="channel">7TV канал</TabsTrigger>
						</TabsList>
						<div class="relative mt-3">
							<Icon name="lucide:search" class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
							<Input v-model="searchQuery" class="pl-9" placeholder="Поиск по имени" aria-label="Поиск эмоций" @keydown.stop />
						</div>
					</div>

					<TabsContent v-for="tab in tabValues" :key="tab" :value="tab" class="max-h-[60dvh] overflow-y-auto px-4 pb-4 pt-3">
						<div v-if="currentStatus === 'loading'" class="grid grid-cols-4 gap-2 sm:grid-cols-6 md:grid-cols-8">
							<Skeleton v-for="index in 24" :key="index" class="aspect-square rounded-lg" />
						</div>
						<div v-else-if="currentStatus === 'error'" class="flex min-h-48 flex-col items-center justify-center gap-3 text-center">
							<Icon name="lucide:triangle-alert" class="size-8 text-destructive" />
							<p class="max-w-sm text-sm text-muted-foreground">{{ currentError }}</p>
							<Button variant="outline" size="sm" @click="retryActiveTab">Повторить</Button>
						</div>
						<div v-else-if="filteredEmotes.length === 0" class="flex min-h-48 flex-col items-center justify-center gap-2 text-center">
							<Icon name="lucide:smile-plus" class="size-8 text-muted-foreground" />
							<p class="max-w-sm text-sm text-muted-foreground">{{ emptyMessage }}</p>
						</div>
						<div v-else class="grid grid-cols-4 gap-2 sm:grid-cols-6 md:grid-cols-8">
							<button v-for="emote in filteredEmotes" :key="emote.id" type="button" class="group flex aspect-square min-w-0 items-center justify-center rounded-lg border border-transparent bg-muted/30 p-2 transition-colors hover:border-border hover:bg-accent focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 active:scale-[0.98]" :title="emote.name" :aria-label="emote.name" @click="handleSelectEmote(emote)">
								<img :src="getEmoteUrl(emote)" :alt="emote.name" width="64" height="64" loading="lazy" class="size-full max-h-16 max-w-16 object-contain transition-transform group-hover:scale-110" />
							</button>
						</div>
					</TabsContent>
				</Tabs>
			</div>
		</DialogOrSheet>
	</Dialog>
</template>
