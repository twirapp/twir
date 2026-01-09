<script lang="ts" setup>
import { PlusIcon } from 'lucide-vue-next'
import { computed, h } from 'vue'

import { useRouter } from 'vue-router'

import type { PageLayoutTab } from '~/layout/page-layout.vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useModerationApi } from '~/features/moderation/composables/use-moderation-api.ts'
import { Icons } from '~/features/moderation/composables/use-moderation-form.ts'
import ModerationTabChatWall from '~/features/moderation/tabs/moderation-tab-chat-wall.vue'
import ModerationTabRules from '~/features/moderation/tabs/moderation-tab-rules.vue'
// oxlint-disable-next-line consistent-type-imports
import { ChannelRolePermissionEnum, ModerationSettingsType } from '~/gql/graphql'
import PageLayout from '~/layout/page-layout.vue'

const { t } = useI18n()
const router = useRouter()
const { user: profile } = storeToRefs(useDashboardAuth())
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
	router.push({ name: 'ModerationForm', query: { ruleType }, params: { id: 'new' } })
}
</script>

<template>
	<PageLayout active-tab="rules" :tabs="tabs">
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
			<UiDropdownMenu v-if="activeTab === 'rules'">
				<UiDropdownMenuTrigger as-child>
					<UiButton :disabled="isCreateDisabled">
						<PlusIcon class="size-4 mr-2" />
						{{ items.length >= maxModerationRules ? t('moderation.limitExceeded') : t('sharedButtons.create') }}
						({{ items.length }}/{{ maxModerationRules }})
					</UiButton>
				</UiDropdownMenuTrigger>
				<UiDropdownMenuContent>
					<UiDropdownMenuItem
						v-for="itemType of ModerationSettingsType"
						:key="itemType"
						@click="createNewRule(itemType)"
					>
						<div class="flex items-center gap-1">
							<component :is="Icons[itemType]" :size="20" />
							<span>{{ t(`moderation.types.${itemType}.name`) }}</span>
						</div>
					</UiDropdownMenuItem>
				</UiDropdownMenuContent>
			</UiDropdownMenu>
		</template>
	</PageLayout>
</template>
