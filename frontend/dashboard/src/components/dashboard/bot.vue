<script setup lang='ts'>
import { NCard, NAlert, NSkeleton, NButton } from 'naive-ui';
import { onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { useBotInfo, useBotJoinPart } from '@/api/index.js';

const { t } = useI18n();

const { data, isLoading, isFetching, refetch } = useBotInfo();

const stateMutation = useBotJoinPart();
const isStateButtonDisabled = ref(true);
async function changeBotState() {
	isStateButtonDisabled.value = true;
	await stateMutation.mutate(data?.value?.enabled ? 'part' : 'join');
}

watch(isFetching, (value) => {
	if (!value) {
		isStateButtonDisabled.value = false;
		return;
	}
});

onMounted(() => {
	refetch();
});
</script>

<template>
	<n-card
		:title="t('dashboard.botManage.title')"
		:content-style="{ padding: isLoading ? '10px' : '0px' }"
		:segmented="{
			content: true,
			footer: 'soft'
		}"
	>
		<n-skeleton v-if="!data" :sharp="false" />

		<n-alert v-else :type="data!.isMod ? 'success' : 'error'" :bordered="false" class="bot-alert">
			<span
				v-if="data?.isMod"
				v-html="t('dashboard.botManage.moderator', { botName: data.botName})"
			/>
			<span
				v-else
				v-html="t('dashboard.botManage.notModerator', { botName: data.botName})"
			/>
		</n-alert>

		<template #action>
			<n-button
				:type="data?.enabled ? 'error' : 'success'"
				block
				:loading="isStateButtonDisabled"
				:disabled="isStateButtonDisabled"
				@click="changeBotState"
			>
				{{ t(`dashboard.botManage.${data?.enabled ? 'leave' : 'join'}`) }}
			</n-button>
		</template>
	</n-card>
</template>

<style scoped>
.bot-alert {
	border-radius: 0px;
}
</style>
