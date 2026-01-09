<script setup lang="ts">
import { CheckIcon, ChevronDownIcon } from 'lucide-vue-next'
import { ref } from 'vue'




import type { AcceptableValue } from 'reka-ui'

interface Props {
	id: string
}

defineProps<Props>()

const model = defineModel<boolean>({ required: true })

const show = ref(false)

function handleSelect(
	event: CustomEvent<{
		originalEvent: PointerEvent
		value?: AcceptableValue
	}>
) {
	if (typeof event.detail.value !== 'string') {
		return
	}

	model.value = event.detail.value === 'true'
	show.value = false
}
</script>

<template>
	<UiPopover v-model:open="show">
		<UiPopoverTrigger as-child>
			<UiButton :id="id" variant="outline" class="flex justify-between" @click="show = true">
				{{ model ? 'Show' : 'Hide' }}
				<ChevronDownIcon class="size-4" />
			</UiButton>
		</UiPopoverTrigger>

		<UiPopoverContent class="p-1">
			<UiCommand>
				<UiCommandList>
					<UiCommandItem value="false" @select="handleSelect">
						<CheckIcon v-if="!model" class="size-4 mr-2" /> Hide
					</UiCommandItem>
					<UiCommandItem value="true" @select="handleSelect">
						<CheckIcon v-if="model" class="size-4 mr-2" /> Show
					</UiCommandItem>
				</UiCommandList>
			</UiCommand>
		</UiPopoverContent>
	</UiPopover>
</template>
