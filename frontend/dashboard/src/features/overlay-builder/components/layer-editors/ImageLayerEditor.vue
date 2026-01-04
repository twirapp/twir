<script setup lang="ts">
import { computed } from 'vue'
import { Image } from 'lucide-vue-next'

import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

import type { Layer } from '../../types'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

const imageUrl = computed({
	get: () => props.layer?.settings?.imageUrl ?? '',
	set: (_value: string) => {
		// Setter is not used directly anymore
	},
})

function handleUrlChange(event: Event) {
	const target = event.target as HTMLInputElement
	emit('update', {
		settings: {
			...props.layer.settings,
			imageUrl: target.value,
		},
	})
}

function setPlaceholder() {
	emit('update', {
		settings: {
			...props.layer.settings,
			imageUrl: 'https://via.placeholder.com/300x200',
		},
	})
}
</script>

<template>
	<div class="space-y-4">
		<div class="flex items-center gap-2 text-sm font-medium">
			<Image class="h-4 w-4" />
			<span>Image Settings</span>
		</div>

		<Separator />

		<div class="space-y-2">
			<Label for="image-url">Image URL</Label>
			<Input
				id="image-url"
				:model-value="imageUrl"
				type="url"
				placeholder="https://example.com/image.png"
				@input="handleUrlChange"
			/>
			<p class="text-xs text-muted-foreground">
				Enter a direct URL to an image (PNG, JPG, GIF, etc.)
			</p>
		</div>

		<div class="space-y-2">
			<Label>Preview</Label>
			<div class="border rounded-lg p-4 bg-slate-50 dark:bg-slate-900 min-h-[100px] flex items-center justify-center">
				<img
					v-if="imageUrl"
					:src="imageUrl"
					alt="Image preview"
					class="max-w-full max-h-[200px] object-contain"
					@error="() => {}"
				/>
				<p v-else class="text-sm text-muted-foreground">
					No image URL provided
				</p>
			</div>
		</div>

		<Button variant="outline" size="sm" class="w-full" @click="setPlaceholder">
			Use Placeholder Image
		</Button>

		<div class="p-3 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg">
			<p class="text-xs text-blue-900 dark:text-blue-100">
				<strong>Tip:</strong> You can use Twir variables in the URL, like
				<code class="px-1 py-0.5 bg-blue-100 dark:bg-blue-900 rounded">$(user.login)</code>
				or
				<code class="px-1 py-0.5 bg-blue-100 dark:bg-blue-900 rounded">$(stream.title)</code>
			</p>
		</div>
	</div>
</template>
