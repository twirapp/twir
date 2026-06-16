<script lang="ts" setup>
import type { PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'

import { computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { useProfile, useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useModerationApi } from '~~/layers/dashboard/features/moderation/composables/use-moderation-api.js'
import { Icons } from '~~/layers/dashboard/features/moderation/composables/use-moderation-form.js'
import ModerationTabChatWall from '~~/layers/dashboard/features/moderation/tabs/moderation-tab-chat-wall.vue'
import ModerationTabRules from '~~/layers/dashboard/features/moderation/tabs/moderation-tab-rules.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

import { Button } from '@/components/ui/button'
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
// oxlint-disable-next-line consistent-type-imports
import { ChannelRolePermissionEnum, ModerationSettingsType } from '~/gql/graphql.js'

const { t } = useI18n()
const router = useRouter()
const { data: profile } = useProfile()
const { items } = useModerationApi()

const canEditModeration = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageModeration)

const maxModerationRules = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxModerationRules ?? 50
})

const isCreateDisabled = computed(() => {
	return items.value.length >= maxModerationRules.value || !canEditModeration.value
})

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

function createNewRule(ruleType: ModerationSettingsType) {
	router.push({ path: `/dashboard/moderation/new`, query: { ruleType } })
}
</script>

<template>
	<PageLayout
		active-tab="rules"
		:tabs="tabs"
	>
		<template #title>
			{{ t('sidebar.moderation') }}
		</template>

		<template #title-footer="{ activeTab }">
			<div
				v-if="activeTab === 'rules'"
				class="flex flex-col gap-0.5"
			>
				<span>{{ t('moderationRules.description.line1') }}</span>
				<span class="text-xs">{{ t('moderationRules.description.line2') }}</span>
			</div>
			<div
				v-if="activeTab === 'chat-wall'"
				class="flex flex-col gap-0.5"
			>
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
					<Button :disabled="isCreateDisabled">
						<Icon
							name="lucide:plus"
							class="mr-2 size-4"
						/>
						{{
							items.length >= maxModerationRules
								? t('moderation.limitExceeded')
								: t('sharedButtons.create')
						}}
						({{ items.length }}/{{ maxModerationRules }})
					</Button>
				</DropdownMenuTrigger>
				<DropdownMenuContent>
					<DropdownMenuItem
						v-for="itemType of ModerationSettingsType"
						:key="itemType"
						@click="createNewRule(itemType)"
					>
						<div class="flex items-center gap-1">
							<Icon
								:name="Icons[itemType]"
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
