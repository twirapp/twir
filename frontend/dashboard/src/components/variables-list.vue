<script setup lang="ts">
import { useClipboard } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { useToast } from './ui/toast'

import { useVariablesApi } from '@/api/variables'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

defineProps<{
	popoverAlign?: 'start' | 'center' | 'end'
	popoverSide?: 'top' | 'right' | 'bottom' | 'left'
}>()

const { t } = useI18n()
const { toast } = useToast()
const clipboard = useClipboard()

const { builtInVariables } = useVariablesApi()

const open = ref(false)

const selectVariables = computed(() => {
	return builtInVariables.value.map((variable) => ({
		label: `$(${variable.example})`,
		value: `$(${variable.example})`,
		description: variable.description,
		links: variable.links,
	}))
})

function handleSelect(value: string) {
	clipboard.copy(value)
	toast({
		title: 'Copied',
		duration: 2500,
	})
	open.value = false
}
</script>

<template>
	<Popover v-model:open="open">
		<PopoverTrigger>
			<slot name="trigger" />
		</PopoverTrigger>
		<PopoverContent
			class="p-0 z-[9999] max-w-[400px]"
			:align="popoverAlign"
			:side="popoverSide"
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
								<div v-if="option.links" class="flex flex-wrap gap-4">
									<a v-for="link of option.links" :key="link.href" :href="link.href" target="_blank" class="text-xs underline">
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
