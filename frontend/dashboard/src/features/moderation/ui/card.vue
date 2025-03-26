<script setup lang="ts">
import { Settings, Trash2 } from 'lucide-vue-next'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import { Icons } from './form/helpers.ts'

import type { ItemWithId } from '@twir/api/messages/moderation/moderation'

import { useModerationManager, useUserAccessFlagChecker } from '@/api'
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
import { useToast } from '@/components/ui/toast/use-toast.ts'
import { ChannelRolePermissionEnum } from '@/gql/graphql.ts'

const props = defineProps<{
	item: ItemWithId
}>()

defineEmits<{
	showSettings: []
}>()

const manager = useModerationManager()
const patcher = manager.patch!
const deleter = manager.deleteOne

const patchExecuting = ref(false)
const showDeleteDialog = ref(false)

const { t } = useI18n()
const { toast } = useToast()

const userCanManageModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration)

async function switchState(id: string, v: boolean) {
	patchExecuting.value = true

	try {
		await patcher.mutateAsync({ id, enabled: v })
	} catch (error) {
		console.error(error)
	} finally {
		patchExecuting.value = false
	}
}

async function removeItem() {
	await deleter.mutateAsync({ id: props.item.id })
	showDeleteDialog.value = false
	toast({
		title: t('sharedTexts.deleted'),
		duration: 2000,
	})
}
</script>

<template>
	<Card
		:title="t(`moderation.types.${item.data!.type}.name`)"
		:icon="Icons[item.data!.type]"
		class="h-full"
	>
		<template #headerExtra>
			<Switch
				:disabled="!userCanManageModeration"
				:checked="item.data!.enabled"
				:aria-disabled="patchExecuting"
				@update:checked="(v) => switchState(item.id, v)"
			/>
		</template>

		<template #content>
			{{ t(`moderation.types.${item.data!.type}.description`) }}
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
