<script setup lang="ts">
import { useField } from 'vee-validate'
import { useI18n } from 'vue-i18n'

import Card from '../ui/card.vue'

import type { ModerationSettingsType } from '@/gql/graphql.ts'

import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Dialog, DialogHeader } from '@/components/ui/dialog'
import { Separator } from '@/components/ui/separator'
import { useModerationApi } from '@/features/moderation/composables/use-moderation-api.ts'
import {
	isModerationEditModalOpened,
	setModerationEditModalState,
} from '@/features/moderation/composables/use-moderation-form.ts'
import Modal from '@/features/moderation/ui/form/modal.vue'

const { t } = useI18n()
const { items } = useModerationApi()

const { setValue: setId } = useField<string | undefined>('id')
const { value: currentEditType } = useField<ModerationSettingsType>('type')

function showForm(itemId: string) {
	setModerationEditModalState(true)
	setId(itemId)
}
</script>

<template>
	<div>
		<div v-if="!items.length">
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
				v-for="item of items"
				:key="item.id"
				:item="item"
				@show-settings="showForm(item.id)"
			/>
		</div>
	</div>

	<Dialog v-model:open="isModerationEditModalOpened">
		<DialogOrSheet>
			<DialogHeader>
				<div class="flex flex-col gap-[2px]">
					<span>
						{{ currentEditType ? t(`moderation.types.${currentEditType}.name`) : 'Edit' }}
					</span>
					<span class="text-xs">
						{{ currentEditType ? t(`moderation.types.${currentEditType}.description`) : '' }}
					</span>
				</div>
			</DialogHeader>

			<Separator />

			<Modal />
		</DialogOrSheet>
	</Dialog>
</template>
