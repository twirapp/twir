<script setup lang="ts">

import { ImportTimersResponse } from '@twir/api/messages/integrations_nightbot/integrations_nightbot';
import { NCard, NButton } from 'naive-ui';
import { computed, ref } from 'vue';


import { useTimersManager, useNightbotIntegration, useNightbotIntegrationImporter } from '@/api';
import { useUserAccessFlagChecker } from '@/api/auth';

const emits = defineEmits<{
	close: []
}>();

const nightbotIntegrationManager = useNightbotIntegration();
const { data } = nightbotIntegrationManager.useData();

const nightbotManager = useNightbotIntegrationImporter();
const nightbotTimersImporter = nightbotManager.useTimersImporter();

const timersManager = useTimersManager();
const { refetch: refetchTimers } = timersManager.getAll({});

const res = ref<ImportTimersResponse | null>(null);

async function onImportClick() {
	res.value = await nightbotTimersImporter.mutateAsync();
	await refetchTimers();
}

const isNightbotIntegrationEnabled = computed(() => {
	return !!data.value?.userName;
});

const userCanManageTimers = useUserAccessFlagChecker('MANAGE_TIMERS');
</script>

<template>
	<div class="flex flex-row justify-between h-full ">
		<n-card class="flex items-end">
			<template #cover>
				<img src="@/assets/integrations/nightbot.png" />
			</template>

			<div class="h-full w-full">
				<div v-if="res">
					<p>Imported Count: {{ res.importedCount }}</p>
					<p>Failed Count: {{ res.failedCount }}</p>
					<p v-if="res.failedTimersNames.length > 0">
						Failed Commands:
					</p>
					<ul v-if="res.failedTimersNames.length > 0" class="overflow-y-scroll max-h-60">
						<li v-for="name in res.failedTimersNames" :key="name">
							{{ name }}
						</li>
					</ul>
				</div>
			</div>

			<template #footer>
				<n-button secondary type="success" :disabled="!userCanManageTimers || !isNightbotIntegrationEnabled" :loading="nightbotTimersImporter.isLoading.value" @click="onImportClick">
					IMPORT
				</n-button>
			</template>
		</n-card>
	</div>
</template>

<style scoped>
.n-card {
  max-width: 300px;
}
</style>
