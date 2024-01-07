<script setup lang="ts">
import { IconMessage } from '@tabler/icons-vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import Card from '@/components/overlays/card.vue';

const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="IconMessage"
		title="Chat"
		:description="t('overlays.chat.description')"
		overlay-path="chat"
		:show-copy="false"
		@open-settings="router.push({ name: 'ChatOverlay' })"
	>
	</card>
</template>
