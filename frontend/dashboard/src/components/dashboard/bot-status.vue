<script setup lang="ts">
import { IconLogin } from '@tabler/icons-vue';
import { NAlert, NButton } from 'naive-ui';
import { computed, h, onMounted } from 'vue';
import { useI18n, I18nT } from 'vue-i18n';

import { useBotInfo, useBotJoinPart } from '@/api';

const { data, refetch } = useBotInfo();
const stateMutation = useBotJoinPart();

async function join() {
	await stateMutation.mutateAsync('join');
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

const haveIssues = computed(() => {
	return !isMod.value || !isJoined.value;
});

const { t } = useI18n();

const alertContent = computed(() => {
	if (!isMod.value) {
		return h(
			'span', {},
			h(
				I18nT,
				{ keypath: 'dashboard.botManage.notModerator' },
				{
					default: () => h('b', {}, `/mod ${data?.value?.botName}`),
				},
			),
		);
	}

	if (!isJoined.value) {
		return h('p', {}, t('dashboard.botManage.notEnabledTitle'));
	}

	return null;
});

</script>
<template>
	<div v-if="haveIssues" class="p-4" :title="t('dashboard.botManage.notEnabledTitle')">
		<n-alert type="error">
			<div class="flex flex-col">
				<component :is="alertContent" />
				<div>
					<n-button v-if="!isJoined" secondary @click="join">
						<IconLogin />
						{{ t(`dashboard.botManage.join`) }}
					</n-button>
				</div>
			</div>
		</n-alert>
	</div>
</template>

<style scoped>

</style>
