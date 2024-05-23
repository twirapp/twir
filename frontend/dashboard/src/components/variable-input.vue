<script setup lang="ts">
import { Variable } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { FunctionalComponent } from 'vue'

import { useVariablesApi } from '@/api/variables'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Textarea } from '@/components/ui/textarea'

withDefaults(defineProps<{
	inputType?: 'text' | 'textarea'
	minRows?: number
	maxRows?: number
}>(), {
	inputType: 'text',
})

defineSlots<{
	underSelect: FunctionalComponent
}>()

const text = defineModel<string>({ default: '' })
const { t } = useI18n()

const { allVariables } = useVariablesApi()

const selectVariables = computed(() => {
	return allVariables.value.map((variable) => ({
		label: `$(${variable.example})`,
		value: `$(${variable.example})`,
		description: variable.description,
	}))
})

const open = ref(false)

function handleSelect(value: string) {
	text.value += ` ${value}`
}
</script>

<template>
	<Popover v-model:open="open">
		<component :is="inputType === 'textarea' ? Textarea : Input" v-model="text" :maxlength="500" />
		<PopoverTrigger as-child>
			<Button class="ml-2" variant="ghost" size="icon">
				<Variable class="size-auto opacity-50" />
			</Button>
		</PopoverTrigger>
		<PopoverContent
			align="start"
			class="p-0 z-[9999] max-w-[400px]"
		>
			<Command
				:reset-search-term-on-blur="false"
			>
				<CommandInput class="h-9" :placeholder="t('sharedTexts.searchPlaceholder')" />
				<CommandEmpty>
					Not found
				</CommandEmpty>
				<CommandList>
					<CommandGroup>
						<CommandItem
							v-for="option in selectVariables"
							:key="option.value"
							:value="option.value"
							@select="handleSelect(option.value)"
						>
							<div class="flex flex-wrap flex-col gap-0.5">
								<span>{{ option.label }}</span>
								<span v-if="option.description" class="text-xs">{{ option.description }}</span>
							</div>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
