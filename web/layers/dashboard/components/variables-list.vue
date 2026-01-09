<script setup lang="ts">
import { toast } from 'vue-sonner'

import { useVariablesApi } from '#layers/dashboard/api/variables'

defineProps<{
	popoverAlign?: 'start' | 'center' | 'end'
	popoverSide?: 'top' | 'right' | 'bottom' | 'left'
}>()

const { t } = useI18n()
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
	toast.success('Copied', {
		duration: 2500,
	})
	open.value = false
}
</script>

<template>
	<UiPopover v-model:open="open">
		<UiPopoverTrigger>
			<slot name="trigger" />
		</UiPopoverTrigger>
		<UiPopoverContent class="p-0 z-9999 max-w-100" :align="popoverAlign" :side="popoverSide">
			<UiCommand :reset-search-term-on-blur="false">
				<UiCommandInput class="h-9" :placeholder="t('sharedTexts.searchPlaceholder')" />
				<UiCommandEmpty> Not found </UiCommandEmpty>
				<UiCommandList>
					<UiCommandGroup>
						<UiCommandItem
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
									>
										{{ link.name }}
									</a>
								</div>
							</div>
						</UiCommandItem>
					</UiCommandGroup>
				</UiCommandList>
			</UiCommand>
		</UiPopoverContent>
	</UiPopover>
</template>
