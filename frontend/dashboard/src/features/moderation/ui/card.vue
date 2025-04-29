<script setup lang="ts">
import { Settings, Trash2 } from 'lucide-vue-next'
import { ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import type { EditableItem } from '@/features/moderation/composables/use-moderation-form.ts'

import { useUserAccessFlagChecker } from '@/api'
import Card from '@/components/card/card.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog'
import { Switch } from '@/components/ui/switch'
import { useModerationApi } from '@/features/moderation/composables/use-moderation-api.ts'
import { Icons } from '@/features/moderation/composables/use-moderation-form.ts'
import { ChannelRolePermissionEnum } from '@/gql/graphql.ts'

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
	const item = manager.items.value.find(i => i.id === id)
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
	<Card
		:title="t(`moderation.types.${item.type}.name`)"
		:icon="Icons[item.type]"
		class="h-full"
	>
		<template #headerExtra>
			<Switch
				v-if="item.id"
				:disabled="!userCanManageModeration"
				:checked="item.enabled"
				:aria-disabled="patchExecuting"
				@update:checked="(v) => switchState(item.id!, v)"
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
				<Button
					:disabled="!userCanManageModeration"
					variant="outline"
					@click="$emit('showSettings')"
				>
					{{ t('sharedButtons.settings') }}
					<Settings class="ml-2 h-4 w-4" />
				</Button>

				<Button
					:disabled="!userCanManageModeration"
					variant="destructive"
					@click="showDeleteDialog = true"
				>
					{{ t('sharedButtons.delete') }}
					<Trash2 class="ml-2 h-4 w-4" />
				</Button>
			</div>
		</template>
	</Card>

	<Dialog v-model:open="showDeleteDialog">
		<DialogContent>
			<DialogHeader>
				<DialogTitle>{{ t('deleteConfirmation.text') }}</DialogTitle>
			</DialogHeader>
			<DialogFooter>
				<Button
					variant="outline"
					@click="showDeleteDialog = false"
				>
					{{ t('deleteConfirmation.cancel') }}
				</Button>
				<Button
					variant="destructive"
					@click="removeItem"
				>
					{{ t('deleteConfirmation.confirm') }}
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
