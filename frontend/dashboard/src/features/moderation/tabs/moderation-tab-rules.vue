<script setup lang="ts">
import { useI18n } from 'vue-i18n'

import Card from '../ui/card.vue'

import { useModerationManager } from '@/api'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Dialog, DialogHeader } from '@/components/ui/dialog'
import { Separator } from '@/components/ui/separator'
import { useModerationRules } from '@/features/moderation/composables/use-moderation-rules.ts'
import { useEditableItem } from '@/features/moderation/ui/form/helpers.ts'
import Modal from '@/features/moderation/ui/form/modal.vue'

const { t } = useI18n()
const manager = useModerationManager()
const { data: settings } = manager.getAll({})

const rules = useModerationRules()
const { editableItem } = useEditableItem()
</script>

<template>
	<div>
		<div v-if="!settings?.body.length">
			<div class="flex flex-col gap-2 items-center justify-center h-full">
				<h2 class="text-2xl font-bold">
					No rules
				</h2>
				<p class="text-sm">
					Create a new rule to start moderating your chat.
				</p>
			</div>
		</div>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<Card
				v-for="item of settings?.body"
				:key="item.id"
				:item="item"
				@show-settings="rules.showSettings(item.id)"
			/>
		</div>
	</div>

	<Dialog v-model:open="rules.settingsOpened.value">
		<DialogOrSheet>
			<DialogHeader>
				<div class="flex flex-col gap-[2px]">
					<span>
						{{ editableItem?.data ? t(`moderation.types.${editableItem.data.type}.name`) : 'Edit' }}
					</span>
					<span class="text-xs">
						{{ editableItem?.data ? t(`moderation.types.${editableItem.data.type}.description`) : '' }}
					</span>
				</div>
			</DialogHeader>

			<Separator />

			<Modal v-if="editableItem" />
		</DialogOrSheet>
	</Dialog>
</template>
