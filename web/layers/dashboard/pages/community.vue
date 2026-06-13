<script setup lang='ts'>
import { computed } from 'vue'

import type { PageLayoutTab } from '~/layers/dashboard/layout/page-layout.vue'

import { useUserAccessFlagChecker } from '~/layers/dashboard/api/auth'
import CommunityChatMessages from '~/layers/dashboard/features/community-chat-messages/community-chat-messages.vue'
import CommunityEmotesStatistic
	from '~/layers/dashboard/features/community-emotes-statistic/community-emotes-statistic.vue'
import CommunityRewardsHistory from '~/layers/dashboard/features/community-rewards-history/community-rewards-history.vue'
import CommunityRoles from '~/layers/dashboard/features/community-roles/community-roles.vue'
import CommunityUsers from '~/layers/dashboard/features/community-users/community-users.vue'
import { ChannelRolePermissionEnum } from '~/app/gql/graphql'
import PageLayout from '~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { t } = useI18n()

const canViewRoles = useUserAccessFlagChecker(ChannelRolePermissionEnum.ViewRoles)

const tabs = computed<PageLayoutTab[]>(() => ([
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
]))
</script>

<template>
	<PageLayout :tabs="tabs" active-tab="users">
		<template #title>
			{{ t('community.title') }}
		</template>
	</PageLayout>
</template>
