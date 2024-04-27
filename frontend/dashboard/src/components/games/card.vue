<script setup lang="ts">

import { IconSettings } from '@tabler/icons-vue';
import { NButton } from 'naive-ui';
import { FunctionalComponent } from 'vue';
import { useI18n } from 'vue-i18n';

import { useUserAccessFlagChecker } from '@/api';
import Card from '@/components/card/card.vue';
import { ChannelRolePermissionEnum } from '@/gql/graphql';

defineEmits<{
	openSettings: [];
}>();

withDefaults(defineProps<{
	description: string;
	title: string;
	icon: FunctionalComponent;
	iconStroke?: number
	showSettings?: boolean
	iconFill?: string
}>(), { showSettings: true });

const { t } = useI18n();

const userCanManageGames = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageGames);
</script>

<template>
	<card :title="title" :icon="icon" :icon-stroke="iconStroke" :icon-fill="iconFill">
		<template #content>
			<p>{{ description }}</p>
		</template>
		<template #footer>
			<n-button
				v-if="showSettings"
				:disabled="!userCanManageGames"
				secondary
				size="large"
				@click="$emit('openSettings')"
			>
				<div class="flex gap-[6px]">
					<span>{{ t('sharedButtons.settings') }}</span>
					<IconSettings />
				</div>
			</n-button>
		</template>
	</card>
</template>
