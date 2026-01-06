<script setup lang="ts">
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api/auth'
import BanSettings from '@/features/chat-alerts/ui/ban-settings.vue'
import ChatAlertsRewardsSettings from '@/features/chat-alerts/ui/chat-alerts-rewards-settings.vue'
import Settings from '@/features/chat-alerts/ui/settings.vue'
import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue'

const { t } = useI18n()
const { data: profile } = useProfile()

const maxChatAlertsMessages = computed(() => {
	const selectedDashboard = profile.value?.availableDashboards.find(
		(d) => d.id === profile.value?.selectedDashboardId
	)
	return selectedDashboard?.plan?.maxChatAlertsMessages ?? 20
})

const pageTabs = computed<PageLayoutTab[]>(() => [
	{
		name: 'followers',
		title: t('chatAlerts.labels.followers'),
		component: () =>
			h(Settings, {
				formKey: 'followers',
				title: t('chatAlerts.labels.followers'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText:
					'Yay, there is new follower, say hello to {user}! Total followers for current stream: {streamFollowers}',
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{user}, {followers}, {streamFollowers}' })}
		`,
			}),
	},
	{
		name: 'raids',
		title: t('chatAlerts.labels.raids'),
		component: () =>
			h(Settings, {
				formKey: 'raids',
				title: t('chatAlerts.labels.raids'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: '{user} raided us with {count} viewers PogChamp',
				count: {
					label: 'Viewers',
				},
				alertMessage: `
			${t('chatAlerts.randomMessageWithCount')}
			${t('chatAlerts.replacedInfo', { vars: '{user}, {count}' })}
		`,
			}),
	},
	{
		name: 'donations',
		title: t('chatAlerts.labels.donations'),
		component: () =>
			h(Settings, {
				formKey: 'donations',
				title: t('chatAlerts.labels.donations'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				count: {
					label: 'Amount',
				},
				defaultMessageText: '{user} just donated {count}{currency} and want to say us {message}',
				alertMessage: `
			${t('chatAlerts.randomMessageWithCount')}
			${t('chatAlerts.replacedInfo', { vars: '{user}, {count}, {currency}, {message}' })}
		`,
			}),
	},
	{
		name: 'subscriptions',
		title: t('chatAlerts.labels.subscriptions'),
		component: () =>
			h(Settings, {
				formKey: 'subscribers',
				title: t('chatAlerts.labels.subscriptions'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				count: {
					label: 'Months',
				},
				defaultMessageText: '{user} just subscribed {month} months in a row',
				alertMessage: `
			${t('chatAlerts.randomMessageWithCount')}
			1 month message will be used for new subscribers. ${t('chatAlerts.replacedInfo', { vars: '{user}, {month}' })}
		`,
			}),
	},
	{
		name: 'rewards',
		title: t('chatAlerts.labels.rewards'),
		component: () =>
			h(
				Settings,
				{
					formKey: 'redemptions',
					title: t('chatAlerts.labels.rewards'),
					minCooldown: 0,
					maxMessages: maxChatAlertsMessages.value,
					defaultMessageText: '{user} activated {reward} reward',
					alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{user}, {reward}' })}
		`,
				},
				{
					additionalSettings: () => h(ChatAlertsRewardsSettings),
				}
			),
	},
	{
		name: 'first-user-message',
		title: t('chatAlerts.labels.firstUserMessage'),
		component: () =>
			h(Settings, {
				formKey: 'firstUserMessage',
				title: t('chatAlerts.labels.firstUserMessage'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: '{user} new on the channel! Say hello.',
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{user}' })}
		`,
			}),
	},
	{
		name: 'stream-online',
		title: t('chatAlerts.labels.streamOnline'),
		component: () =>
			h(Settings, {
				formKey: 'streamOnline',
				title: t('chatAlerts.labels.streamOnline'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: "We're just online in {category} | {title}",
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{title}, {category}' })}
		`,
			}),
	},
	{
		name: 'stream-offline',
		title: t('chatAlerts.labels.streamOffline'),
		component: () =>
			h(Settings, {
				formKey: 'streamOffline',
				title: t('chatAlerts.labels.streamOffline'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: "We're now offline, stay in touch, follow socials.",
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
		`,
			}),
	},
	{
		name: 'chat-cleared',
		title: t('chatAlerts.labels.chatCleared'),
		component: () =>
			h(Settings, {
				formKey: 'chatCleared',
				title: t('chatAlerts.labels.chatCleared'),
				minCooldown: 2,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: 'Chat cleared, but who knows why? Kappa',
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
		`,
			}),
	},
	{
		name: 'ban',
		title: t('chatAlerts.labels.userBanned'),
		component: () =>
			h(
				Settings,
				{
					formKey: 'ban',
					title: t('chatAlerts.labels.userBanned'),
					minCooldown: 2,
					maxMessages: maxChatAlertsMessages.value,
					defaultMessageText:
						'How dare are you {userName}? Glad we have {moderatorName} to calm you down. Please sit {time} in prison for {reason}, and think about your behavior.',
					count: {
						label: t('chatAlerts.ban.countLabel'),
					},
					minCount: 0,
					alertMessage: `
			${t('chatAlerts.ban.alertInfo')}
			${t('chatAlerts.randomMessageWithCount')}
			${t('chatAlerts.replacedInfo', { vars: `{userName}, {moderatorName}, {time} - seconds or 'permanent', {reason}` })}
		`,
				},
				{
					additionalSettings: () => h(BanSettings),
				}
			),
	},
	{
		name: 'unban-request-create',
		title: t('chatAlerts.labels.channelUnbanRequestCreate'),
		component: () =>
			h(Settings, {
				formKey: 'unbanRequestCreate',
				title: t('chatAlerts.labels.channelUnbanRequestCreate'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: 'User {userName} requesting unban with message {message}',
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{userName}, {message}' })}
		`,
			}),
	},
	{
		name: 'unban-request-resolve',
		title: t('chatAlerts.labels.channelUnbanRequestResolve'),
		component: () =>
			h(Settings, {
				formKey: 'unbanRequestResolve',
				title: t('chatAlerts.labels.channelUnbanRequestResolve'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText:
					'User {userName} unban request resolved with message {message} by moderator {moderatorName}',
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{userName}, {moderatorName}, {message}' })}
		`,
			}),
	},
	{
		name: 'message-delete',
		title: t('chatAlerts.labels.messageDelete'),
		component: () =>
			h(Settings, {
				formKey: 'messageDelete',
				title: t('chatAlerts.labels.messageDelete'),
				minCooldown: 0,
				maxMessages: maxChatAlertsMessages.value,
				defaultMessageText: 'Message of user {userName} deleted',
				alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{userName}' })}
		`,
			}),
	},
])
</script>

<template>
	<PageLayout active-tab="followers" :tabs="pageTabs">
		<template #title>
			{{ t('sidebar.chatAlerts') }}
		</template>
	</PageLayout>
</template>

<style scoped>
/* TODO: webkit line clamp */
.tags :deep(.n-tag__content) {
	text-overflow: ellipsis;
	white-space: nowrap;
	overflow: hidden;
}

.tags :deep(.n-space),
.tags :deep(.n-tag) {
	width: 100%;
}
</style>
