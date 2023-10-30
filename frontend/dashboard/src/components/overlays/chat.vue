<script setup lang="ts">
import { IconMessage } from '@tabler/icons-vue';
import { NModal } from 'naive-ui';
import { ref } from 'vue';

import Settings from './chat/settings.vue';

import { useChatOverlayManager } from '@/api/index.js';
import Card from '@/components/overlays/card.vue';

const isModalOpened = ref(false);
const chatManager = useChatOverlayManager();
const { data: settings, isError } = chatManager.getSettings();
</script>

<template>
	<card
		:icon="IconMessage"
		title="Chat"
		description="chat"
		overlay-path="chat"
		:copy-disabled="!settings || isError"
		@open-settings="isModalOpened = true"
	>
	</card>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Chat"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<Settings />
	</n-modal>
</template>
