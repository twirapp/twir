<script setup lang="ts">
import { IconMessage } from '@tabler/icons-vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useChatOverlayManager } from '@/api/index.js';
import Card from '@/components/overlays/card.vue';

const chatManager = useChatOverlayManager();
const { data: settings, isError } = chatManager.getSettings();
const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="IconMessage"
		title="Chat"
		:description="t('overlays.chat.description')"
		overlay-path="chat"
		:copy-disabled="!settings || isError"
		@open-settings="router.push({ name: 'ChatOverlay' })"
	>
	</card>
</template>
