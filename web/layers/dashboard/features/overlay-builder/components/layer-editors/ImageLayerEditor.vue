<script setup lang="ts">
import { toRef } from 'vue'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

import type { Layer } from '../../types'
import FieldInput from '../fields/FieldInput.vue'
import { useImageLayerEditor } from '../../composables/useImageLayerEditor'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()
const { t } = useI18n()

const emit = defineEmits<{
	update: [updates: Partial<Layer>]
}>()

const { imageUrl, setPlaceholder } = useImageLayerEditor(toRef(props, 'layer'), (updates) => {
	emit('update', { settings: { ...props.layer.settings, ...updates } })
})
</script>

<template>
	<div class="space-y-4">
		<div class="flex items-center gap-2 text-sm font-medium">
			<Icon name="lucide:image" class="h-4 w-4" />
			<span>{{ t('overlayBuilder.editors.image.title') }}</span>
		</div>

		<Separator />

		<FieldInput
			id="image-url"
			v-model="imageUrl"
			:label="t('overlayBuilder.editors.image.url')"
			type="url"
			placeholder="https://example.com/image.png"
			:description="t('overlayBuilder.editors.image.urlDescription')"
		/>

		<Button variant="outline" size="sm" class="w-full" @click="setPlaceholder">
			{{ t('overlayBuilder.editors.image.usePlaceholder') }}
		</Button>

		<div class="p-3 bg-blue-50 dark:bg-blue-950 border border-blue-200 dark:border-blue-800 rounded-lg">
			<p class="text-xs text-blue-900 dark:text-blue-100">
				<strong>{{ t('overlayBuilder.editors.image.tip') }}</strong> {{ t('overlayBuilder.editors.image.variablesTip') }}
				<code class="px-1 py-0.5 bg-blue-100 dark:bg-blue-900 rounded">$(user.login)</code>
				or
				<code class="px-1 py-0.5 bg-blue-100 dark:bg-blue-900 rounded">$(stream.title)</code>
			</p>
		</div>
	</div>
</template>
