<script setup lang="ts">
import { Variable } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import type { FunctionalComponent } from 'vue'

import { useVariablesApi } from '@/api/variables'
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

withDefaults(
	defineProps<{
		inputType?: 'text' | 'textarea'
		minRows?: number
		maxRows?: number
		popoverAlign?: 'start' | 'center' | 'end'
		popoverSide?: 'top' | 'right' | 'bottom' | 'left'
	}>(),
	{
		inputType: 'text',
	}
)

defineSlots<{
	'additional-buttons': FunctionalComponent
}>()

const text = defineModel<string | undefined>({ default: '' })
const { t } = useI18n()

const { allVariables } = useVariablesApi()

const selectVariables = computed(() => {
	return allVariables.value.map((variable) => ({
		label: `$(${variable.example})`,
		value: `$(${variable.example})`,
		description: variable.description,
		links: variable.links,
	}))
})

const open = ref(false)

function handleSelect(value: string) {
	text.value += ` ${value}`
}
</script>

<template>
	<Popover v-model:open="open">
		<div class="flex flex-col w-full group relative">
			<component
				v-bind="$attrs"
				:is="inputType === 'textarea' ? Textarea : Input"
				v-model="text"
				class="input pr-10 w-full"
				:maxlength="500"
			/>
			<div class="flex gap-0.5 absolute right-1 top-1" :class="{ 'opacity-100!': open }">
				<PopoverTrigger as-child>
					<button class="hover:bg-secondary/80 p-1 rounded-md">
						<Variable class="size-4 opacity-50" />
					</button>
				</PopoverTrigger>
				<slot name="additional-buttons" />
			</div>
		</div>
		<PopoverContent class="p-0 z-9999 max-w-[600px]" :align="popoverAlign" :side="popoverSide">
			<Command :reset-search-term-on-blur="false">
				<CommandInput class="h-9" :placeholder="t('sharedTexts.searchPlaceholder')" />
				<CommandEmpty> Not found </CommandEmpty>
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
								<div v-if="option.links" class="flex flex-wrap gap-4">
									<a
										v-for="link of option.links"
										:key="link.href"
										:href="link.href"
										target="_blank"
										class="text-xs underline"
										@click.stop
									>
										{{ link.name }}
									</a>
								</div>
							</div>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</PopoverContent>
	</Popover>
</template>
