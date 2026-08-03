<script setup lang="ts">
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
</script>

<template>
	<div
		class="flex flex-col gap-4 rounded-lg border p-4 sm:flex-row sm:items-start sm:justify-between"
		data-test="scope-group"
		:data-group="scope.group"
	>
		<div class="flex min-w-0 flex-col gap-1">
			<span class="font-medium">{{ scope.name }}</span>
			<span class="text-muted-foreground text-sm">{{ scope.description }}</span>
		</div>

		<div class="flex shrink-0 flex-col gap-3">
			<div class="flex items-center justify-between gap-4 sm:justify-end">
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
				<div class="flex items-center justify-between gap-4 sm:justify-end">
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
				<p v-if="edit" :id="editHintId" class="text-muted-foreground text-xs sm:text-right">
					{{ t('mcpConsent.scopes.editIncludesRead') }}
				</p>
			</div>
		</div>
	</div>
</template>
