<script setup lang="ts">
import type { PageLayoutTab } from '~~/layers/dashboard/app/layout/page-layout.vue'

import { computed } from 'vue'
import { useUserAccessFlagChecker } from '~~/layers/dashboard/app/api/auth'
import CommunityChatMessages from '~~/layers/dashboard/app/features/community-chat-messages/community-chat-messages.vue'
import CommunityEmotesStatistic from '~~/layers/dashboard/app/features/community-emotes-statistic/community-emotes-statistic.vue'
import CommunityRewardsHistory from '~~/layers/dashboard/app/features/community-rewards-history/community-rewards-history.vue'
import CommunityRoles from '~~/layers/dashboard/app/features/community-roles/community-roles.vue'
import CommunityUsers from '~~/layers/dashboard/app/features/community-users/community-users.vue'
import PageLayout from '~~/layers/dashboard/app/layout/page-layout.vue'

import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })

const { t } = useI18n()

const canViewRoles = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewRoles)

const tabs = computed<PageLayoutTab[]>(() => [
	{
		title: 'Chat logs',
		component: CommunityChatMessages,
		name: 'chat-logs',
	},
	{
		title: t('community.users.title'),
		component: CommunityUsers,
		name: 'users',
	},
	{
		title: t('sidebar.roles'),
		component: CommunityRoles,
		name: 'permissions',
		disabled: !canViewRoles.value,
	},
	{
		title: t('community.emotesStatistic.title'),
		component: CommunityEmotesStatistic,
		name: 'emotes-stats',
	},
	{
		title: 'Rewards history',
		component: CommunityRewardsHistory,
		name: 'rewards-history',
	},
])
</script>

<template>
	<PageLayout
		:tabs="tabs"
		active-tab="users"
	>
		<template #title>
			{{ t('community.title') }}
		</template>
	</PageLayout>
</template>
