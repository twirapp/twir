<script setup lang="ts">
import { XIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import { cn } from '@/lib/utils'

interface Option {
	label: string
	value: string
}

interface Props {
	modelValue: string[]
	options: Option[]
	placeholder?: string
	disabled?: boolean
	class?: string
}

const props = withDefaults(defineProps<Props>(), {
	placeholder: 'Select items...',
	disabled: false,
})

const emit = defineEmits<{
	'update:modelValue': [value: string[]]
}>()

const open = ref(false)
const search = ref('')

const selectedValues = computed(() => new Set(props.modelValue))

const filteredOptions = computed(() => {
	if (!search.value) return props.options

	return props.options.filter(option =>
		option.label.toLowerCase().includes(search.value.toLowerCase()),
	)
})

const selectedLabels = computed(() => {
	return props.options
		.filter(option => selectedValues.value.has(option.value))
		.map(option => option.label)
})

function toggleOption(value: string) {
	const newSelectedValues = new Set(selectedValues.value)

	if (newSelectedValues.has(value)) {
		newSelectedValues.delete(value)
	} else {
		newSelectedValues.add(value)
	}

	emit('update:modelValue', Array.from(newSelectedValues))
}

function removeOption(value: string) {
	const newSelectedValues = new Set(selectedValues.value)
	newSelectedValues.delete(value)
	emit('update:modelValue', Array.from(newSelectedValues))
}

function clearOptions() {
	emit('update:modelValue', [])
}

// Reset search when popover closes
watch(open, (isOpen) => {
	if (!isOpen) {
		search.value = ''
	}
})
</script>

<template>
	<Popover v-model:open="open">
		<PopoverTrigger as-child :disabled="disabled">
			<Button
				variant="outline"
				role="combobox"
				:aria-expanded="open"
				:class="[
					cn('w-full justify-between h-fit', props.class),
				]"
				:disabled="disabled"
			>
				<div class="flex gap-1 flex-wrap">
					<Badge
						v-for="(value, i) in props.modelValue"
						:key="value"
						variant="default"
						class="mr-1 mb-1"
					>
						{{ props.options.find(option => option.value === value)?.label || value }}
						<button
							class="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
							@click.stop="removeOption(value)"
						>
							<XIcon class="h-3 w-3 text-background hover:text-background" />
						</button>
					</Badge>
					<span v-if="!props.modelValue.length" class="text-muted-foreground">
						{{ placeholder }}
					</span>
				</div>
			</Button>
		</PopoverTrigger>
		<PopoverContent class="w-full p-0">
			<Command>
				<CommandInput v-model="search" placeholder="Search options..." />
				<CommandList>
					<CommandEmpty>No options found.</CommandEmpty>
					<CommandGroup>
						<CommandItem
							v-for="option in filteredOptions"
							:key="option.value"
							:value="option.value"
							:data-selected="selectedValues.has(option.value)"
							@click="toggleOption(option.value)"
						>
							<span
								:class="cn(
									'mr-2 h-4 w-4 border border-primary rounded-sm flex items-center justify-center',
									selectedValues.has(option.value)
										? 'bg-primary text-primary-foreground'
										: 'opacity-50',
								)"
							>
								<svg
									v-if="selectedValues.has(option.value)"
									xmlns="http://www.w3.org/2000/svg"
									viewBox="0 0 24 24"
									width="16"
									height="16"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
								>
									<polyline points="20 6 9 17 4 12" />
								</svg>
							</span>
							{{ option.label }}
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
