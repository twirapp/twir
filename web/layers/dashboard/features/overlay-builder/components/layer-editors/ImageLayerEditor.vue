<script setup lang="ts">
import { computed } from 'vue'
import { Image } from 'lucide-vue-next'



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
	set: (value: string) => {
		emit('update', {
			settings: {
				...props.layer.settings,
				imageUrl: value,
			},
		})
	},
})

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

		<UiSeparator />

		<div class="space-y-2">
			<UiLabel for="image-url">Image URL</UiLabel>
			<UiInput
				id="image-url"
				v-model="imageUrl"
				type="url"
				placeholder="https://example.com/image.png"
				@keydown.stop
			/>
			<p class="text-xs text-muted-foreground">
				Enter a direct URL to an image (PNG, JPG, GIF, etc.)
			</p>
		</div>

		<UiButton variant="outline" size="sm" class="w-full" @click="setPlaceholder">
			Use Placeholder Image
		</UiButton>

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
