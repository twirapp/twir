<script setup lang="ts">
import { Settings, Trash2 } from 'lucide-vue-next'
import { ref, toRaw } from 'vue'


import type { EditableItem } from '~/features/moderation/composables/use-moderation-form.ts'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import Card from '#layers/dashboard/components/card/card.vue'



import { useModerationApi } from '~/features/moderation/composables/use-moderation-api.ts'
import { Icons } from '~/features/moderation/composables/use-moderation-form.ts'
import { ChannelRolePermissionEnum } from '~/gql/graphql.ts'

const props = defineProps<{
	item: EditableItem
}>()

defineEmits<{
	showSettings: []
}>()

const manager = useModerationApi()

const patchExecuting = ref(false)
const showDeleteDialog = ref(false)

const { t } = useI18n()

const userCanManageModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration)

async function switchState(id: string, v: boolean) {
	patchExecuting.value = true
	const item = manager.items.value.find((i) => i.id === id)
	if (!item) return

	const data = structuredClone(toRaw(item))

	delete data.id
	delete data.createdAt
	delete data.updatedAt

	try {
		await manager.update(id, {
			...data,
			enabled: v,
		})
	} catch (error) {
		console.error(error)
	} finally {
		patchExecuting.value = false
	}
}

async function removeItem() {
	if (!props.item.id) return

	await manager.remove(props.item.id)
	showDeleteDialog.value = false
}
</script>

<template>
	<Card :title="t(`moderation.types.${item.type}.name`)" :icon="Icons[item.type]" class="h-full">
		<template #headerExtra>
			<UiSwitch
				v-if="item.id"
				:disabled="!userCanManageModeration"
				:model-value="item.enabled"
				:aria-disabled="patchExecuting"
				@update:model-value="(v: boolean) => switchState(item.id!, v)"
			/>
		</template>

		<template #content>
			<div class="flex flex-col gap-2">
				<span v-if="item.name" class="text-xl text-white">{{ item.name }}</span>
				{{ t(`moderation.types.${item.type}.description`) }}
			</div>
		</template>

		<template #footer>
			<div class="flex gap-2">
				<UiButton
					:disabled="!userCanManageModeration"
					variant="outline"
					@click="$emit('showSettings')"
				>
					{{ t('sharedButtons.settings') }}
					<Settings class="ml-2 h-4 w-4" />
				</UiButton>

				<UiButton
					:disabled="!userCanManageModeration"
					variant="destructive"
					@click="showDeleteDialog = true"
				>
					{{ t('sharedButtons.delete') }}
					<Trash2 class="ml-2 h-4 w-4" />
				</UiButton>
			</div>
		</template>
	</Card>

	<UiDialog v-model:open="showDeleteDialog">
		<UiDialogContent>
			<UiDialogHeader>
				<UiDialogTitle>{{ t('deleteConfirmation.text') }}</UiDialogTitle>
			</UiDialogHeader>
			<UiDialogFooter>
				<UiButton variant="outline" @click="showDeleteDialog = false">
					{{ t('deleteConfirmation.cancel') }}
				</UiButton>
				<UiButton variant="destructive" @click="removeItem">
					{{ t('deleteConfirmation.confirm') }}
				</UiButton>
			</UiDialogFooter>
		</UiDialogContent>
	</UiDialog>
</template>
