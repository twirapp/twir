<script setup lang="ts">
import { NCard, NTabs, NTabPane, useMessage, NButton } from 'naive-ui';
import { ref, toRaw, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { type ChatAlertsSettings } from '@/api/modules/index.js';
import { useChatAlertsSettings, useChatAlertsSettingsUpdate } from '@/api/modules/index.js';
import Settings from '@/components/chatAlerts/settings.vue';


const formValue = ref<Required<ChatAlertsSettings>>({
	chatCleared: {
		enabled: true,
		messages: [],
	},
	cheers: {
		enabled: true,
		messages: [],
	},
	donations: {
		enabled: true,
		messages: [],
	},
	firstUserMessage: {
		enabled: false,
		messages: [],
	},
	followers: {
		enabled: true,
		messages: [],
	},
	raids: {
		enabled: true,
		messages: [],
	},
	redemptions: {
		enabled: true,
		messages: [],
	},
	streamOffline: {
		enabled: true,
		messages: [],
	},
	streamOnline: {
		enabled: true,
		messages: [],
	},
	subscribers: {
		enabled: true,
		messages: [],
	},
});

const { data: settings } = useChatAlertsSettings();
const updater = useChatAlertsSettingsUpdate();

watch(settings, (v) => {
	if (!v) return;

	const raw = toRaw(v);

	for (const key of Object.keys(raw)) {
		// eslint-disable-next-line @typescript-eslint/ban-ts-comment
		// @ts-ignore
		formValue.value[key] = raw[key];
	}
});

const message = useMessage();
const { t } = useI18n();

async function save() {
	const raw = toRaw(formValue.value);

	try {
		await updater.mutateAsync(raw);
		message.success(t('sharedTexts.saved'));
	} catch (error) {
		message.error(t('sharedTexts.errorOnSave'));
	}
}
</script>

<template>
	<n-card
		:title="t('sidebar.chatAlerts')"
		segmented
	>
		<template #header-extra>
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</template>

		<n-tabs type="line">
			<n-tab-pane name="followers" tab="Followers">
				<Settings
					v-model:enabled="formValue.followers.enabled"
					v-model:messages="formValue.followers.messages"
					:max-messages="20"
					default-message-text="Yay, there is new follower, say hello to {user}!"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{user}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="raids" tab="Raids">
				<Settings
					v-model:enabled="formValue.raids.enabled"
					v-model:messages="formValue.raids.messages"
					:max-messages="20"
					default-message-text="{user} raided us with {count} viewers PogChamp"
					with-count
					count-label="Viewers"
					:alert-message="`
						${t('chatAlerts.randomMessageWithCount')}
						${t('chatAlerts.replacedInfo', { vars: '{user}, {count}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="donations" tab="Donations">
				<Settings
					v-model:enabled="formValue.donations.enabled"
					v-model:messages="formValue.donations.messages"
					:max-messages="20"
					with-count
					count-label="Amount"
					default-message-text="{user} just donated {count}{currency} and want to say us {message}"
					:alert-message="`
						${t('chatAlerts.randomMessageWithCount')}
						${t('chatAlerts.replacedInfo', { vars: '{user}, {count}, {currency}, {message}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="subscriptions" tab="Subscriptions">
				<Settings
					v-model:enabled="formValue.subscribers.enabled"
					v-model:messages="formValue.subscribers.messages"
					:max-messages="500"
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

			<n-tab-pane name="rewards" tab="Rewards">
				<Settings
					v-model:enabled="formValue.redemptions.enabled"
					v-model:messages="formValue.redemptions.messages"
					:max-messages="20"
					default-message-text="{user} activated {reward} reward"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{user}, {reward}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="firstUserMessage" tab="First user message">
				<Settings
					v-model:enabled="formValue.firstUserMessage.enabled"
					v-model:messages="formValue.firstUserMessage.messages"
					:max-messages="20"
					default-message-text="{user} new on the channel! Say hello."
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{user}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="streamOnline" tab="Stream Online">
				<Settings
					v-model:enabled="formValue.streamOnline.enabled"
					v-model:messages="formValue.streamOnline.messages"
					:max-messages="20"
					default-message-text="We're just online in {category} | {title}"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
						${t('chatAlerts.replacedInfo', { vars: '{title}, {category}'})}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="streamOffline" tab="Stream Offline">
				<Settings
					v-model:enabled="formValue.streamOffline.enabled"
					v-model:messages="formValue.streamOffline.messages"
					:max-messages="20"
					default-message-text="We're now offline, stay in touch, follow socials."
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
					`"
				/>
			</n-tab-pane>

			<n-tab-pane name="chatClearted" tab="Chat cleared">
				<Settings
					v-model:enabled="formValue.chatCleared.enabled"
					v-model:messages="formValue.chatCleared.messages"
					:max-messages="20"
					default-message-text="Chat cleared, but who knows why? Kappa"
					:alert-message="`
						${t('chatAlerts.randomedMessage')}
					`"
				/>
			</n-tab-pane>
		</n-tabs>
	</n-card>
</template>