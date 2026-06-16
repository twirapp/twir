<script setup lang="ts">
import type { FunctionalComponent } from 'vue'

import { computed, ref } from 'vue'
import { useVariablesApi } from '~~/layers/dashboard/api/variables'

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
		<div class="group relative flex w-full flex-col">
			<component
				v-bind="$attrs"
				:is="inputType === 'textarea' ? Textarea : Input"
				v-model="text"
				class="input w-full pr-10"
				:maxlength="500"
			/>
			<div
				class="absolute top-1 right-1 flex gap-0.5"
				:class="{ 'opacity-100!': open }"
			>
				<PopoverTrigger as-child>
					<button class="hover:bg-secondary/80 rounded-md p-1">
						<Icon
							name="lucide:variable"
							class="size-4 opacity-50"
						/>
					</button>
				</PopoverTrigger>
				<slot name="additional-buttons" />
			</div>
		</div>
		<PopoverContent
			class="z-9999 max-w-[600px] p-0"
			:align="popoverAlign"
			:side="popoverSide"
		>
			<Command :reset-search-term-on-blur="false">
				<CommandInput
					class="h-9"
					:placeholder="t('sharedTexts.searchPlaceholder')"
				/>
				<CommandEmpty> Not found </CommandEmpty>
				<CommandList>
					<CommandGroup>
						<CommandItem
							v-for="option in selectVariables"
							:key="option.value"
							:value="option.value"
							@select="handleSelect(option.value)"
						>
							<div class="flex flex-col flex-wrap gap-0.5">
								<span>{{ option.label }}</span>
								<span
									v-if="option.description"
									class="text-xs"
									>{{ option.description }}</span
								>
								<div
									v-if="option.links"
									class="flex flex-wrap gap-4"
								>
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
