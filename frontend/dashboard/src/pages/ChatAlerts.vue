<script setup lang="ts">
import {
	NCard,
	NTabs,
	NTabPane,
	NButton,
	NDynamicTags,
	NFormItem,
	useNotification,
} from 'naive-ui';
import { ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useChatAlertsApi } from '@/api/chat-alerts.js';
import { useUserAccessFlagChecker } from '@/api/index.js';
import Settings from '@/components/chatAlerts/settings.vue';

const formValue = ref({
	chatCleared: {
		enabled: false,
		messages: [],
		cooldown: 2,
	},
	cheers: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	donations: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	firstUserMessage: {
		enabled: false,
		messages: [],
		cooldown: 2,
	},
	followers: {
		enabled: false,
		messages: [],
		cooldown: 3,
	},
	raids: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	redemptions: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	streamOffline: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	streamOnline: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	subscribers: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	ban: {
		enabled: false,
		messages: [],
		cooldown: 2,
		ignoreTimeoutFrom: [],
	},
	unbanRequestCreate: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
	unbanRequestResolve: {
		enabled: false,
		messages: [],
		cooldown: 0,
	},
});

const chatAlertsApi = useChatAlertsApi();
const updateChatAlerts = chatAlertsApi.useMutationUpdateChatAlerts();

watch(chatAlertsApi.chatAlerts, (v) => {
	if (!v) return;

	const raw = toRaw(v);

	for (const key of Object.keys(raw)) {
		// eslint-disable-next-line @typescript-eslint/ban-ts-comment
		// @ts-ignore
		formValue.value[key] = raw[key];
	}
}, { immediate: true });

const message = useNotification();
const { t } = useI18n();

async function save() {
	const raw = toRaw(formValue.value);
	if (!raw) return;

	try {
		await updateChatAlerts.executeMutation({ input: raw });
		message.success({
			title: t('sharedTexts.saved'),
			duration: 2500,
		});
	} catch (error) {
		message.error({
			title: t('sharedTexts.errorOnSave'),
			duration: 2500,
		});
	}
}

const hasAccessToManageAlerts = useUserAccessFlagChecker('MANAGE_ALERTS');
</script>

<template>
	<n-card
		:title="t('sidebar.chatAlerts')"
		segmented
	>
		<template #header-extra>
			<n-button secondary type="success" :disabled="!hasAccessToManageAlerts" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</template>

		<n-tabs type="line">
			<n-tab-pane name="followers" :tab="t('chatAlerts.labels.followers')">
				<Settings
					v-model:enabled="formValue.followers.enabled"
					v-model:messages="formValue.followers.messages"
					v-model:cooldown="formValue.followers.cooldown"
					:min-cooldown="0"
					:max-messages="20"
					default-message-text="Yay, there is new follower, say hello to {user}!"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{user}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="raids" :tab="t('chatAlerts.labels.raids')">
				<Settings
					v-model:enabled="formValue.raids.enabled"
					v-model:messages="formValue.raids.messages"
					v-model:cooldown="formValue.raids.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="{user} raided us with {count} viewers PogChamp"
					with-count
					count-label="Viewers"
					:alert-message="`
						${t('chatAlerts.randomMessageWithCount')}
						${t('chatAlerts.replacedInfo', { vars: '{user}, {count}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="donations" :tab="t('chatAlerts.labels.donations')">
				<Settings
					v-model:enabled="formValue.donations.enabled"
					v-model:messages="formValue.donations.messages"
					v-model:cooldown="formValue.donations.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					with-count
					count-label="Amount"
					default-message-text="{user} just donated {count}{currency} and want to say us {message}"
					:alert-message="`
						${t('chatAlerts.randomMessageWithCount')}
						${t('chatAlerts.replacedInfo', { vars: '{user}, {count}, {currency}, {message}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="subscriptions" :tab="t('chatAlerts.labels.subscriptions')">
				<Settings
					v-model:enabled="formValue.subscribers.enabled"
					v-model:messages="formValue.subscribers.messages"
					v-model:cooldown="formValue.subscribers.cooldown"
					:max-messages="500"
					:min-cooldown="0"
					with-count
					count-label="Months"
					default-message-text="{user} just subscribed {month} months in a row"
					:alert-message="`
						${t('chatAlerts.randomMessageWithCount')}
						1 month message will be used for new subscribers. ${t('chatAlerts.replacedInfo', { vars: '{user}, {month}'})}
					`"
				/>
			</n-tab-pane>

			<!-- <n-tab-pane name="cheers" tab="Cheers">
				<Settings
					v-model:enabled="formValue.cheers.enabled"
					v-model:messages="formValue.cheers.messages"
					:max-messages="500"
					with-count
					count-label="Months"
					default-message-text="{user} just donated {count}{currency} and want to say us {message}"
					:alert-message="`
						{user}, {month} – will be replaced with actual information.
					`"
				/>
			</n-tab-pane> -->

			<n-tab-pane name="rewards" :tab="t('chatAlerts.labels.rewards')">
				<Settings
					v-model:enabled="formValue.redemptions.enabled"
					v-model:messages="formValue.redemptions.messages"
					v-model:cooldown="formValue.redemptions.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="{user} activated {reward} reward"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{user}, {reward}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="firstUserMessage" :tab="t('chatAlerts.labels.firstUserMessage')">
				<Settings
					v-model:enabled="formValue.firstUserMessage.enabled"
					v-model:messages="formValue.firstUserMessage.messages"
					v-model:cooldown="formValue.firstUserMessage.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="{user} new on the channel! Say hello."
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{user}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="streamOnline" :tab="t('chatAlerts.labels.streamOnline')">
				<Settings
					v-model:enabled="formValue.streamOnline.enabled"
					v-model:messages="formValue.streamOnline.messages"
					v-model:cooldown="formValue.streamOnline.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="We're just online in {category} | {title}"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{title}, {category}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="streamOffline" :tab="t('chatAlerts.labels.streamOffline')">
				<Settings
					v-model:enabled="formValue.streamOffline.enabled"
					v-model:messages="formValue.streamOffline.messages"
					v-model:cooldown="formValue.streamOffline.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="We're now offline, stay in touch, follow socials."
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="chatCleared" :tab="t('chatAlerts.labels.chatCleared')">
				<Settings
					v-model:enabled="formValue.chatCleared.enabled"
					v-model:messages="formValue.chatCleared.messages"
					v-model:cooldown="formValue.chatCleared.cooldown"
					:max-messages="20"
					:min-cooldown="2"
					default-message-text="Chat cleared, but who knows why? Kappa"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="ban" :tab="t('chatAlerts.labels.userBanned')">
				<Settings
					v-model:enabled="formValue.ban.enabled"
					v-model:messages="formValue.ban.messages"
					v-model:cooldown="formValue.ban.cooldown"
					:max-messages="20"
					:min-cooldown="2"
					default-message-text="How dare are you {userName}? Glad we have {moderatorName} to calm you down. Please sit {time} in prison for {reason}, and think about your behavior."
					:count-label="t('chatAlerts.ban.countLabel')"
					with-count
					:min-count="0"
					:alert-message="`
						${t('chatAlerts.ban.alertInfo')}
						${t('chatAlerts.randomMessageWithCount')}
						${t('chatAlerts.replacedInfo', { vars: `{userName}, {moderatorName}, {time} - seconds or 'permanent', {reason}`})}
					`"
				>
					<template #header>
						<n-form-item
							:label="t('chatAlerts.ban.ignoreTimeoutFrom')" label-style="padding: 0;"
							class="tags"
						>
							<n-dynamic-tags
								v-model:value="formValue.ban.ignoreTimeoutFrom"
								:max="100"
							/>
						</n-form-item>
					</template>
				</Settings>
			</n-tab-pane>

			<n-tab-pane name="channelUnbanRequestCreate" :tab="t('chatAlerts.labels.channelUnbanRequestCreate')">
				<Settings
					v-model:enabled="formValue.unbanRequestCreate.enabled"
					v-model:messages="formValue.unbanRequestCreate.messages"
					v-model:cooldown="formValue.unbanRequestCreate.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="User {userName} requesting unban with message {message}"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{userName}, {message}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="channelUnbanRequestResolve" :tab="t('chatAlerts.labels.channelUnbanRequestResolve')">
				<Settings
					v-model:enabled="formValue.unbanRequestResolve.enabled"
					v-model:messages="formValue.unbanRequestResolve.messages"
					v-model:cooldown="formValue.unbanRequestResolve.cooldown"
					:max-messages="20"
					:min-cooldown="0"
					default-message-text="User {userName} unban request resolved with message {message} by moderator {moderatorName}"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{userName}, {moderatorName}, {message}'})}
					`"
				/>
			</n-tab-pane>
		</n-tabs>
	</n-card>
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
