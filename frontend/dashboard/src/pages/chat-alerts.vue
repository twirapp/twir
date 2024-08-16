<script setup lang="ts">
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import BanSettings from '@/features/chat-alerts/ui/ban-settings.vue'
import Settings from '@/features/chat-alerts/ui/settings.vue'
import PageLayout, { type PageLayoutTab } from '@/layout/page-layout.vue'

const { t } = useI18n()

const pageTabs = computed<PageLayoutTab[]>(() => [
	{
		name: 'followers',
		title: t('chatAlerts.labels.followers'),
		component: () => h(Settings, {
			formKey: 'followers',
			minCooldown: 0,
			maxMessages: 20,
			defaultMessageText: 'Yay, there is new follower, say hello to {user}!',
			alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{user}' })}
		`,
		}),
	},
	{
		name: 'raids',
		title: t('chatAlerts.labels.raids'),
		component: () => h(Settings, {
			formKey: 'raids',
			minCooldown: 0,
			maxMessages: 20,
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
		component: () => h(Settings, {
			formKey: 'donations',
			minCooldown: 0,
			maxMessages: 20,
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
		component: () => h(Settings, {
			formKey: 'subscribers',
			minCooldown: 0,
			maxMessages: 500,
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
		component: () => h(Settings, {
			formKey: 'redemptions',
			minCooldown: 0,
			maxMessages: 20,
			defaultMessageText: '{user} activated {reward} reward',
			alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{user}, {reward}' })}
		`,
		}),
	},
	{
		name: 'first-user-message',
		title: t('chatAlerts.labels.firstUserMessage'),
		component: () => h(Settings, {
			formKey: 'firstUserMessage',
			minCooldown: 0,
			maxMessages: 20,
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
		component: () => h(Settings, {
			formKey: 'streamOnline',
			minCooldown: 0,
			maxMessages: 20,
			defaultMessageText: 'We\'re just online in {category} | {title}',
			alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{title}, {category}' })}
		`,
		}),
	},
	{
		name: 'stream-offline',
		title: t('chatAlerts.labels.streamOffline'),
		component: () => h(Settings, {
			formKey: 'streamOffline',
			minCooldown: 0,
			maxMessages: 20,
			defaultMessageText: 'We\'re now offline, stay in touch, follow socials.',
			alertMessage: `
			${t('chatAlerts.randomedMessage')}
		`,
		}),
	},
	{
		name: 'chat-cleared',
		title: t('chatAlerts.labels.chatCleared'),
		component: () => h(Settings, {
			formKey: 'chatCleared',
			minCooldown: 2,
			maxMessages: 20,
			defaultMessageText: 'Chat cleared, but who knows why? Kappa',
			alertMessage: `
			${t('chatAlerts.randomedMessage')}
		`,
		}),
	},
	{
		name: 'ban',
		title: t('chatAlerts.labels.userBanned'),
		component: () => h(Settings, {
			formKey: 'ban',
			minCooldown: 2,
			maxMessages: 20,
			defaultMessageText: 'How dare are you {userName}? Glad we have {moderatorName} to calm you down. Please sit {time} in prison for {reason}, and think about your behavior.',
			count: {
				label: t('chatAlerts.ban.countLabel'),
			},
			minCount: 0,
			alertMessage: `
			${t('chatAlerts.ban.alertInfo')}
			${t('chatAlerts.randomMessageWithCount')}
			${t('chatAlerts.replacedInfo', { vars: `{userName}, {moderatorName}, {time} - seconds or 'permanent', {reason}` })}
		`,
		}, {
			additionalSettings: () => h(BanSettings),
		}),
	},
	{
		name: 'unban-request-create',
		title: t('chatAlerts.labels.channelUnbanRequestCreate'),
		component: () => h(Settings, {
			formKey: 'unbanRequestCreate',
			minCooldown: 0,
			maxMessages: 20,
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
		component: () => h(Settings, {
			formKey: 'unbanRequestResolve',
			minCooldown: 0,
			maxMessages: 20,
			defaultMessageText: 'User {userName} unban request resolved with message {message} by moderator {moderatorName}',
			alertMessage: `
			${t('chatAlerts.randomedMessage')}
			${t('chatAlerts.replacedInfo', { vars: '{userName}, {moderatorName}, {message}' })}
		`,
		}),
	},
	{
		name: 'message-delete',
		title: t('chatAlerts.labels.messageDelete'),
		component: () => h(Settings, {
			formKey: 'messageDelete',
			minCooldown: 0,
			maxMessages: 20,
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

.tags :deep(.n-space), .tags :deep(.n-tag) {
	width: 100%;
}
</style>
