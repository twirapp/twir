<script lang="ts" setup>
import { PlusIcon } from 'lucide-vue-next'
import { h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { PageLayoutTab } from '@/layout/page-layout.vue'

import { useUserAccessFlagChecker } from '@/api/index.js'
import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useModerationRules } from '@/features/moderation/composables/use-moderation-rules.ts'
import ModerationTabChatWall from '@/features/moderation/tabs/moderation-tab-chat-wall.vue'
import ModerationTabRules from '@/features/moderation/tabs/moderation-tab-rules.vue'
import { Icons, availableSettings, availableSettingsTypes, useEditableItem } from '@/features/moderation/ui/form/helpers.ts'
import { ChannelRolePermissionEnum } from '@/gql/graphql'
import PageLayout from '@/layout/page-layout.vue'

const { editableItem } = useEditableItem()
const rules = useModerationRules()

const { t } = useI18n()

const canEditModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration)

async function createNewItem(itemType: string) {
	const defaultSettings = availableSettings.find(s => s.type === itemType)
	if (!defaultSettings) return
	editableItem.value = {
		data: structuredClone(defaultSettings),
	}
	rules.settingsOpened.value = true
}

const tabs: PageLayoutTab[] = [
	{
		name: 'rules',
		title: 'Rules',
		component: h(ModerationTabRules),
	},
	{
		name: 'chat-wall',
		title: 'Chat wall (firewall)',
		component: h(ModerationTabChatWall),
	},
]
</script>

<template>
	<PageLayout active-tab="rules" :tabs="tabs" sticky-header>
		<template #title>
			{{ t('sidebar.moderation') }}
		</template>

		<template #title-footer="{ activeTab }">
			<div v-if="activeTab === 'rules'" class="flex flex-col gap-0.5">
				<span>{{ t('moderationRules.description.line1') }}</span>
				<span class="text-xs">{{ t('moderationRules.description.line2') }}</span>
			</div>
			<div v-if="activeTab === 'chat-wall'" class="flex flex-col gap-0.5">
				<span>
					{{ t('chatWall.description.line1') }}
				</span>
				<span class="text-xs text-orange-500">
					{{ t('chatWall.description.line2') }}
				</span>
			</div>
		</template>

		<template #action="{ activeTab }">
			<DropdownMenu v-if="activeTab === 'rules'">
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
</template>
