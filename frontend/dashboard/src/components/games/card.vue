<script setup lang="ts">

import { IconSettings } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api';
import Card from '@/components/card/card.vue';

defineEmits<{
	openSettings: [];
}>();

withDefaults(defineProps<{
	description: string;
	title: string;
	icon: FunctionalComponent;
	showSettings?: boolean
}>(), { showSettings: true });

const { t } = useI18n();

const userCanManageGames = useUserAccessFlagChecker('MANAGE_GAMES');
</script>

<template>
	<card :title="title" :icon="icon">
		<template #content>
			<p>Ask the magic 8ball a question and it will answer you.</p>
		</template>
		<template #footer>
			<n-button
				v-if="showSettings"
				:disabled="!userCanManageGames"
				secondary
				size="large"
				@click="$emit('openSettings')"
			>
				<span>{{ t('sharedButtons.settings') }}</span>
				<IconSettings />
			</n-button>
		</template>
	</card>
</template>

<style scoped>

</style>
