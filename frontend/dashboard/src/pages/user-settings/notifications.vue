<script setup lang="ts">
import { type FormInst, type SelectOption, NCard, NFormItem, NButton, NInput, NSpace, NSwitch, NSelect, NAvatar, NText, SelectRenderTag } from 'naive-ui';
import { ref } from 'vue';
import { computed, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { useStreamers } from '@/api/streamers';

const { t } = useI18n();

interface StreamerOption {
	label: string
	value: string
	followers: number
	avatar: string
}

const { data: streamers } = useStreamers();
const streamersOptions = computed<SelectOption[]>(() => {
	if (!streamers.value?.streamers) return [];
	return streamers.value.streamers.map((streamer) => ({
		label: streamer.userDisplayName,
		value: streamer.userId,
		followers: streamer.followersCount,
		avatar: streamer.avatar,
	}));
});

const renderLabel = (option: StreamerOption) => {
	return h('div', { class: 'flex items-center' },
		[
			h(NAvatar, {
				src: option.avatar,
				round: true,
				size: 'small',
			}),
			h('div', { class: 'ml-3 py-1' }, [
				h('div', null, [option.label]),
				h(NText, { depth: 3, tag: 'div' }, { default: () => `${option.followers} followers` }),
			]),
		],
	);
};

const renderSingleSelectTag: SelectRenderTag = ({ option }) => {
	return h(
		'div',
		{ class: 'flex items-center' },
		[
			h(NAvatar, {
				src: option.avatar as string,
				round: true,
				size: 24,
				class: 'mr-3',
			}),
			option.label as string,
		],
	);
};

type FormParams = {
	userId?: string;
	message: string;
	url?: string
};

const isUserMessage = ref(false);
const formRef = ref<FormInst | null>(null);

const formData = ref<FormParams>({
	userId: undefined,
	message: '',
	url: undefined,
});

async function sendNotification() {
	// formRef
	// TODO: api
}
</script>

<template>
	<div class="w-full flex flex-wrap gap-4">
		<n-card :title="t('userSettings.notifications.createNotification')" size="small" bordered>
			<n-form-item>
				<n-space align="center" item-style="display: flex;">
					<n-select
						v-model:value="formData.userId" :render-label="renderLabel"
						:render-tag="renderSingleSelectTag"
						:disabled="!isUserMessage" filterable placeholder="Please select a streamer"
						:options="streamersOptions"
					/>
					<n-switch v-model:value="isUserMessage" />
				</n-space>
			</n-form-item>

			<n-form-item :label="t('userSettings.notifications.messageLabel')">
				<n-input v-model:value="formData.message" type="textarea" placeholder="" :autosize="{ minRows: 3 }" />
			</n-form-item>

			<n-form-item :label="t('userSettings.notifications.urlLabel')">
				<n-input v-model:value="formData.url" type="text" placeholder="" />
			</n-form-item>

			<n-button @click="sendNotification">
				{{ t('sharedButtons.send') }}
			</n-button>
		</n-card>
	</div>
</template>
