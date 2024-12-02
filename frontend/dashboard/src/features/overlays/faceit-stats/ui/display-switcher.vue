<script setup lang="ts">
import { CheckIcon, ChevronDownIcon } from 'lucide-vue-next'
import { ref } from 'vue'

import type { SelectEvent } from 'radix-vue/dist/Combobox/ComboboxItem'
import type { AcceptableValue } from 'radix-vue/dist/Combobox/ComboboxRoot'

import { Button } from '@/components/ui/button'
import { Command, CommandItem, CommandList } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

interface Props {
	id: string
}

defineProps<Props>()

const model = defineModel<boolean>({ required: true })

const show = ref(false)

function handleSelect(event: SelectEvent<AcceptableValue>) {
	if (typeof event.detail.value !== 'boolean') {
		return
	}

	model.value = event.detail.value
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
			<Command v-model:modelValue="model">
				<CommandList>
					<CommandItem :value="false" @select="handleSelect">
						<CheckIcon v-if="!model" class="size-4 mr-2" />	Hide
					</CommandItem>
					<CommandItem :value="true" @select="handleSelect">
						<CheckIcon v-if="model" class="size-4 mr-2" /> Show
					</CommandItem>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
