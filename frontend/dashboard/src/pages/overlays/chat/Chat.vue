<script setup lang="ts">
import {
	ChatBox,
	type Message,
	type Settings as ChatBoxSettings,
	BadgeVersion,
} from '@twir/frontend-chat';
import { useIntervalFn } from '@vueuse/core';
import {
	NTabs,
	NTabPane,
	NAlert,
	NScrollbar,
	useThemeVars,
	NText,
} from 'naive-ui';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useChatOverlayForm } from './components/form.js';
import Form from './components/Form.vue';
import { globalBadges } from './constants.js';
import * as faker from './faker.js';

import {
	useChatOverlayManager, useUserAccessFlagChecker,
} from '@/api/index.js';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js';


const themeVars = useThemeVars();
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const chatManager = useChatOverlayManager();
const creator = chatManager.useCreate();
const deleter = chatManager.useDelete();

const { t } = useI18n();
const { dialog } = useNaiveDiscrete();

const {
	data: entities,
} = chatManager.useGetAll();

const globalBadgesObject = Object.fromEntries(globalBadges);

const messagesMock = ref<Message[]>([]);

useIntervalFn(() => {
	if (!formValue.value) return;

	const internalId = crypto.randomUUID();

	messagesMock.value.push({
		sender: faker.firstName(),
		chunks: [{
			type: 'text',
			value: formValue.value.direction === 'left' || formValue.value.direction === 'right'
				? faker.loremWithLen(3)
				: faker.lorem(),
		}],
		createdAt: new Date(),
		internalId,
		isAnnounce: faker.boolean(),
		isItalic: false,
		type: 'message',
		senderColor: faker.rgb(),
		announceColor: '',
		badges: {
			[faker.randomObjectKey(globalBadgesObject)]: '1',
		},
		id: crypto.randomUUID(),
		senderDisplayName: faker.firstName(),
	});

	if (formValue.value.messageHideTimeout != 0) {
		setTimeout(() => {
			messagesMock.value = messagesMock.value.filter(m => m.internalId != internalId);
		}, formValue.value.messageHideTimeout * 1000);
	}

	if (messagesMock.value.length >= 20) {
		messagesMock.value = messagesMock.value.slice(1);
	}
}, 1 * 1000);

const openedTab = ref<string>();

const { data: formValue, $setData, $getDefaultSettings } = useChatOverlayForm();

function resetTab() {
	if (!entities.value?.settings.at(0)) {
		openedTab.value = undefined;
		return;
	}

	openedTab.value = entities.value.settings.at(0)!.id;
}

watch(entities, () => {
	resetTab();
}, { immediate: true });

watch(openedTab, (v) => {
	const entity = entities.value?.settings.find(s => s.id === v);
	if (!entity) return;

	$setData(entity);
});

watch(openedTab, () => {
	messagesMock.value = [];
});

const chatBoxSettings = computed<ChatBoxSettings>(() => {
	return {
		channelId: '',
		channelName: '',
		channelDisplayName: '',
		globalBadges,
		channelBadges: new Map<string, BadgeVersion>(),
		...formValue.value,
	};
});

async function handleClose(id: string) {
	dialog.create({
		title: 'Delete preset',
		content: 'Are you sure you want to delete this preset?',
		positiveText: 'Delete',
		negativeText: 'Cancel',
		showIcon: false,
		onPositiveClick: async () => {
			const entity = entities.value?.settings.find(s => s.id === id);
			if (!entity?.id) return;

			await deleter.mutateAsync(entity.id);
			resetTab();
		},
	});
}

async function handleAdd() {
	await creator.mutateAsync($getDefaultSettings());
}

const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.settings.length ?? 0) < 5;
});
</script>

<template>
	<div class="page">
		<div class="chatBox" style="width: 70%">
			<ChatBox
				v-if="openedTab"
				class="chatBox"
				:messages="messagesMock"
				:settings="chatBoxSettings"
			/>
			<div v-else style="display: flex; justify-content: center; align-items: center; height: 100%;">
				<n-text style="font-size: 1rem;">
					Preview of chat will be here when you select some preset
				</n-text>
			</div>
		</div>
		<div style="width: 30%">
			<n-tabs
				v-model:value="openedTab"
				type="card"
				:closable="userCanEditOverlays"
				:addable="addable"
				tab-style="min-width: 80px;"
				@close="handleClose"
				@add="handleAdd"
			>
				<template #prefix>
					{{ t('overlays.chat.presets') }}
				</template>
				<template v-if="entities?.settings.length">
					<n-tab-pane
						v-for="(entity, entityIndex) in entities?.settings"
						:key="entity.id"
						:tab="`#${entityIndex+1}`"
						:name="entity.id!"
					>
						<n-scrollbar style="max-height: 75vh;" trigger="none">
							<Form />
						</n-scrollbar>
					</n-tab-pane>
				</template>
			</n-tabs>
			<n-alert v-if="!entities?.settings.length" type="info" style="margin-top: 8px;">
				Create new overlay for edit settings
			</n-alert>
		</div>
	</div>
</template>

<style scoped>
@import '../styles.css';

:deep(.chat) {
	height: 80dvh;
}

.chatBox {
	background-color: v-bind('themeVars.cardColor');
	border-radius: 8px;
	height: 80dvh;
}
</style>
