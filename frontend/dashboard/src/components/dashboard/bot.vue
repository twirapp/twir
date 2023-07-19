<script setup lang='ts'>
import { NCard, NAlert, NSkeleton, NButton } from 'naive-ui';
import { onMounted, ref, watch } from 'vue';

import { useBotInfo, useBotJoinPart } from '@/api/index.js';

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
		title="Bot manage"
		:content-style="{ padding: isLoading ? '10px' : '0px' }"
		:segmented="{
			content: true,
			footer: 'soft'
		}"
	>
		<n-skeleton v-if="!data" :sharp="false" />

		<n-alert v-else :type="data!.isMod ? 'success' : 'error'" :bordered="false" class="bot-alert">
			<span v-if="data?.isMod">
				<b>{{ data!.botName }}</b> is a moderator.
			</span>
			<span v-else>
				We have found that the bot is not a moderator on this channel. Please, use <b>/mod {{ data!.botName }}</b>, or some of functionality may work incorrectly.
			</span>
		</n-alert>

		<template #action>
			<n-button
				:type="data?.enabled ? 'error' : 'success'"
				block
				:loading="isStateButtonDisabled"
				:disabled="isStateButtonDisabled"
				@click="changeBotState"
			>
				{{ data?.enabled ? 'Leave' : 'Join' }}
			</n-button>
		</template>
	</n-card>
</template>

<style scoped>
.bot-alert {
	border-radius: 0px;
}
</style>
