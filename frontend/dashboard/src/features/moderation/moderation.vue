<script lang="ts" setup>
import { PlusIcon } from 'lucide-vue-next'
import { ref, toRaw } from 'vue'
import { useI18n } from 'vue-i18n'

import { useModerationManager, useUserAccessFlagChecker } from '@/api/index.js'
import DialogOrSheet from '@/components/dialog-or-sheet.vue'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogHeader,
} from '@/components/ui/dialog'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import Card from '@/features/moderation/ui/card.vue'
import { Icons, availableSettings, availableSettingsTypes , useEditableItem } from '@/features/moderation/ui/form/helpers.ts'
import Modal from '@/features/moderation/ui/form/modal.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const manager = useModerationManager()
const { data: settings } = manager.getAll({})

const { editableItem } = useEditableItem()

const settingsOpened = ref(false)
function showSettings(id: string) {
	const item = settings.value?.body.find(i => i.id === id)
	if (!item) return

	editableItem.value = toRaw(item)
	settingsOpened.value = true
}

const { t } = useI18n()

const canEditModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration)

async function createNewItem(itemType: string) {
	const defaultSettings = availableSettings.find(s => s.type === itemType)
	if (!defaultSettings) return
	editableItem.value = {
		data: structuredClone(defaultSettings),
	}
	settingsOpened.value = true
}
</script>

<template>
	<PageLayout>
		<template #title>
			{{ t('sidebar.moderation') }}
		</template>

		<template #content>
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
					@show-settings="showSettings(item.id)"
				/>
			</div>
		</template>

		<template #action>
			<DropdownMenu>
				<DropdownMenuTrigger as-child>
					<Button :disabled="!canEditModeration">
						<PlusIcon class="size-4 mr-2" />
						{{ t('sharedButtons.create') }}
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent>
					<DropdownMenuItem
						v-for="itemType of availableSettingsTypes"
						:key="itemType"
						@click="createNewItem(itemType)"
					>
						<div class="flex items-center gap-1">
							<component
								:is="Icons[itemType]"
								:size="20"
							/>
							<span>{{ t(`moderation.types.${itemType}.name`) }}</span>
						</div>
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</template>
	</PageLayout>

	<Dialog v-model:open="settingsOpened">
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
