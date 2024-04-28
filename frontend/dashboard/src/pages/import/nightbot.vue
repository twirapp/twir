<script setup lang="ts">
import type {
	ImportCommandsResponse,
	ImportTimersResponse,
} from '@twir/api/messages/integrations_nightbot/integrations_nightbot';
import { NButton, NCard, NAlert } from 'naive-ui';
import { computed, ref } from 'vue';

import { useUserAccessFlagChecker } from '@/api/auth';
import {
	useNightbotIntegration,
	useNightbotIntegrationImporter,
} from '@/api/index.js';
import IconNightbot from '@/assets/integrations/nightbot.svg?use';
import OauthComponent from '@/components/integrations/variants/oauth.vue';
import { ChannelRolePermissionEnum } from '@/gql/graphql';

const integrationManager = useNightbotIntegration();
const { data: authLink } = integrationManager.useAuthLink();
const { data } = integrationManager.useData();
const logout = integrationManager.useLogout();

const nightbotManager = useNightbotIntegrationImporter();
const nightbotCommandsImporter = nightbotManager.useCommandsImporter();
const nightbotTimersImporter = nightbotManager.useTimersImporter();

const commandsResponse = ref<ImportCommandsResponse | null>(null);
async function importCommands() {
	commandsResponse.value = await nightbotCommandsImporter.mutateAsync();
}

const timersResponse = ref<ImportTimersResponse | null>(null);
async function importTimers() {
	timersResponse.value = await nightbotTimersImporter.mutateAsync();
}

const isNightbotIntegrationEnabled = computed(() => {
	return !!data.value?.userName;
});

const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands);
const userCanManageTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageTimers);
</script>

<template>
	<oauth-component
		title="Nightbot"
		:data="data"
		:logout="() => logout.mutateAsync({})"
		:authLink="authLink?.link"
		:icon="IconNightbot"
		with-settings
	>
		<template #description>
			<i18n-t keypath="integrations.nightbot.info" />
		</template>

		<template #settings>
			<div class="flex w-full gap-4">
				<n-card class="flex h-full" title="Commands">
					<div class="h-full">
						<div v-if="commandsResponse">
							<p>Imported Count: {{ commandsResponse.importedCount }}</p>
							<p>Failed Count: {{ commandsResponse.failedCount }}</p>
							<p v-if="commandsResponse.failedCommandsNames.length > 0">
								Failed Commands:
							</p>
							<ul v-if="commandsResponse.failedCommandsNames.length > 0" class="overflow-y-scroll max-h-60">
								<li v-for="name in commandsResponse.failedCommandsNames" :key="name">
									{{ name }}
								</li>
							</ul>
						</div>
						<n-alert v-else type="info">
							Waiting import...
						</n-alert>
					</div>

					<template #footer>
						<n-button
							secondary
							type="success"
							:disabled="!isNightbotIntegrationEnabled || !userCanManageCommands"
							:loading="nightbotCommandsImporter.isLoading.value"
							block
							@click="importCommands"
						>
							IMPORT
						</n-button>
					</template>
				</n-card>

				<n-card class="flex h-full" title="Timers">
					<div class="h-full w-full">
						<div v-if="timersResponse">
							<p>Imported Count: {{ timersResponse.importedCount }}</p>
							<p>Failed Count: {{ timersResponse.failedCount }}</p>
							<p v-if="timersResponse.failedTimersNames.length > 0">
								Failed Timers:
							</p>
							<ul v-if="timersResponse.failedTimersNames.length > 0" class="overflow-y-scroll max-h-60">
								<li v-for="name in timersResponse.failedTimersNames" :key="name">
									{{ name }}
								</li>
							</ul>
						</div>
						<n-alert v-else type="info">
							Waiting import...
						</n-alert>
					</div>

					<template #footer>
						<n-button
							secondary
							type="success"
							:disabled="!isNightbotIntegrationEnabled || !userCanManageTimers"
							:loading="nightbotTimersImporter.isLoading.value"
							block
							@click="importTimers"
						>
							IMPORT
						</n-button>
					</template>
				</n-card>
			</div>
		</template>
	</oauth-component>
</template>
