<script setup lang="ts">
import { AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Switch } from '@/components/ui/switch'
import type { McpRequestedScope } from '~/utils/mcp-scopes.js'

interface Props {
	scope: McpRequestedScope
	read: boolean
	edit: boolean
	disabled: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
	(e: 'update:read', value: boolean): void
	(e: 'update:edit', value: boolean): void
}>()

const { t } = useI18n()

const readSwitchId = computed(() => `mcp-scope-${props.scope.group}-read`)
const editSwitchId = computed(() => `mcp-scope-${props.scope.group}-edit`)
const editHintId = computed(() => `mcp-scope-${props.scope.group}-edit-hint`)
const hasEditAction = computed(() => props.scope.actions.includes('edit'))

const selectionSummary = computed(() => {
	if (props.edit) return `${t('mcpConsent.scopes.read')}, ${t('mcpConsent.scopes.edit')}`
	if (props.read) return t('mcpConsent.scopes.read')
	return t('mcpConsent.scopes.none')
})
</script>

<template>
	<AccordionItem
		v-slot="{ open }"
		:value="scope.group"
		data-test="scope-group"
		:data-group="scope.group"
	>
		<AccordionTrigger class="items-center gap-3 py-2.5 hover:no-underline">
			<span class="min-w-0 flex-1 truncate">{{ scope.name }}</span>
			<span
				class="shrink-0 text-xs font-normal"
				:class="read ? 'text-muted-foreground' : 'text-muted-foreground/60'"
				data-test="selection-summary"
			>
				{{ selectionSummary }}
			</span>
		</AccordionTrigger>

		<AccordionContent force-mount :class="open ? undefined : 'hidden'">
			<div class="flex flex-col gap-3">
				<p class="text-muted-foreground text-sm">{{ scope.description }}</p>

				<div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:gap-6">
					<div class="flex items-center justify-between gap-4 sm:justify-start">
						<label class="cursor-pointer text-sm font-medium" :for="readSwitchId">
							{{ t('mcpConsent.scopes.read') }}
						</label>
						<Switch
							:id="readSwitchId"
							data-test="read-switch"
							:model-value="read"
							:disabled="disabled"
							@update:model-value="emit('update:read', $event)"
						/>
					</div>

					<div v-if="hasEditAction" class="flex flex-col gap-1">
						<div class="flex items-center justify-between gap-4 sm:justify-start">
							<label class="cursor-pointer text-sm font-medium" :for="editSwitchId">
								{{ t('mcpConsent.scopes.edit') }}
							</label>
							<Switch
								:id="editSwitchId"
								data-test="edit-switch"
								:model-value="edit"
								:disabled="disabled"
								:aria-describedby="edit ? editHintId : undefined"
								@update:model-value="emit('update:edit', $event)"
							/>
						</div>
						<p v-if="edit" :id="editHintId" class="text-muted-foreground text-xs">
							{{ t('mcpConsent.scopes.editIncludesRead') }}
						</p>
					</div>
				</div>
			</div>
		</AccordionContent>
	</AccordionItem>
</template>
