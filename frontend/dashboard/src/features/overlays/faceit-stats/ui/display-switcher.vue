<script setup lang="ts">
import { CheckIcon, ChevronDownIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import { Button } from '@/components/ui/button'
import { Command, CommandItem, CommandList } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
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
	<Popover v-model:open="show">
		<PopoverTrigger as-child>
			<Button :id="id" variant="outline" class="flex justify-between" @click="show = true">
				{{ model ? 'Show' : 'Hide' }}
				<ChevronDownIcon class="size-4" />
			</Button>
		</PopoverTrigger>

		<PopoverContent class="p-1">
			<Command>
				<CommandList>
					<CommandItem value="false" @select="handleSelect">
						<CheckIcon v-if="!model" class="size-4 mr-2" /> Hide
					</CommandItem>
					<CommandItem value="true" @select="handleSelect">
						<CheckIcon v-if="model" class="size-4 mr-2" /> Show
					</CommandItem>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
