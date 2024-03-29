<script setup lang="ts">
import {
	IconChevronRight,
	IconLogin,
	IconLogout,
	IconRobotOff,
	IconRobot,
} from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { NButton, useThemeVars } from 'naive-ui';
import { storeToRefs } from 'pinia';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

import { useBotInfo, useBotJoinPart } from '@/api';
import { useSidebarCollapseStore } from '@/layout/use-sidebar-collapse';

const { data, refetch } = useBotInfo();
const stateMutation = useBotJoinPart();

async function changeBotState() {
	await stateMutation.mutateAsync(data?.value?.enabled ? 'part' : 'join');
}

onMounted(() => {
	refetch();
});

const isMod = computed(() => {
	return data.value?.isMod ?? false;
});

const isJoined = computed(() => {
	return data.value?.enabled ?? false;
});

const cardBackgroundColor = computed(() => {
	if (isMod.value && isJoined.value) {
		return 'rgba(42, 148, 125, 0.25)';
	}

	return 'rgba(208, 58, 82, 0.25)';
});

const { t } = useI18n();
const themeVars = useThemeVars();


const isCollapsed = useLocalStorage('twirIsBotStatusCollapsed', false);

const collapsedStore = useSidebarCollapseStore();
const { isCollapsed: isSidebarCollapsed } = storeToRefs(collapsedStore);
</script>

<template>
	<div
		class="bot-status-card"
	>
		<div
			class="header"
			:style="{
				justifyContent: isSidebarCollapsed ? 'center' : 'space-between',
				cursor: isSidebarCollapsed ? 'not-allowed' : 'pointer',
			}"
			@click="isCollapsed = !isCollapsed"
		>
			<div style="display: flex; gap: 4px; align-items: center;">
				<IconRobotOff v-if="!isMod || !isJoined" />
				<IconRobot v-else />
				<span v-if="!isSidebarCollapsed">{{ t('dashboard.botManage.title') }}</span>
			</div>

			<n-button v-if="!isSidebarCollapsed" text>
				<IconChevronRight
					:style="{
						transition: '0.2s transform ease',
						transform: `rotate(${isCollapsed ? 90 : -90}deg)`
					}"
				/>
			</n-button>
		</div>

		<div v-if="!isCollapsed && !isSidebarCollapsed" class="body">
			<span class="title">
				<template v-if="!isJoined">{{ t('dashboard.botManage.notEnabledTitle') }}</template>
				<template v-else-if="!isMod">
					<i18n-t
						keypath="dashboard.botManage.notModerator"
					>
						<b>/mod {{ data?.botName }}</b>
					</i18n-t>
				</template>
				<template v-else>
					<i18n-t
						keypath="dashboard.botManage.success"
					>
						<b>{{ data?.botName }}</b>
					</i18n-t>
				</template>
			</span>

			<n-button
				v-if="!isCollapsed"
				block
				secondary
				:loading="stateMutation.isLoading.value"
				size="small"
				@click="changeBotState"
			>
				<template #icon>
					<IconLogin v-if="!data?.enabled" />
					<IconLogout v-else />
				</template>
				{{ t(`dashboard.botManage.${data?.enabled ? 'leave' : 'join'}`) }}
			</n-button>
		</div>
	</div>
</template>

<style scoped>
.bot-status-card {
	background-color: v-bind(cardBackgroundColor);
	border-bottom: 1px solid v-bind('themeVars.borderColor');
}

.header {
	display: flex;
	align-items: center;
	padding: 4px;
	-webkit-user-select: none;
	-ms-user-select: none;
	user-select: none;
}

.body {
	padding: 8px;
}

.title {
	font-size: 13px;
}
</style>
