<script setup lang="ts">
import { InfoIcon, LoaderCircleIcon } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import IconNightbot from '~/assets/integrations/nightbot.svg?use'
import OauthComponent from '#layers/dashboard/components/integrations/variants/oauth.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

import { useNightbotIntegration } from './composables/use-nightbot-integration.js'

import type { NightbotImportCommandsOutput, NightbotImportTimersOutput } from '~/gql/graphql.js'

const nightbotIntegration = useNightbotIntegration()
const { data: authLinkData } = nightbotIntegration.useAuthLink()
const { data: nightbotData, fetching: dataFetching, executeQuery } = nightbotIntegration.useData()

onMounted(() => {
	nightbotIntegration.nightbotBroadcaster.onmessage = async (event) => {
		if (event.data !== 'refresh') return

		await executeQuery({ requestPolicy: 'network-only' })
	}
})

const commandsResponse = ref<NightbotImportCommandsOutput | null>(null)
const commandsImporting = ref(false)
async function importCommands() {
	commandsImporting.value = true
	try {
		const result = await nightbotIntegration.importCommands.executeMutation({})
		if (result.data?.nightbotImportCommands) {
			commandsResponse.value = result.data.nightbotImportCommands
		}
	} finally {
		commandsImporting.value = false
	}
}

const timersResponse = ref<NightbotImportTimersOutput | null>(null)
const timersImporting = ref(false)
async function importTimers() {
	timersImporting.value = true
	try {
		const result = await nightbotIntegration.importTimers.executeMutation({})
		if (result.data?.nightbotImportTimers) {
			timersResponse.value = result.data.nightbotImportTimers
		}
	} finally {
		timersImporting.value = false
	}
}

async function logout() {
	await nightbotIntegration.logout.executeMutation({})
}

const isNightbotIntegrationEnabled = computed(() => {
	return !!nightbotData.value?.nightbotGetData?.userName
})

const userCanManageCommands = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageCommands)
const userCanManageTimers = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageTimers)
</script>

<template>
	<oauth-component
		title="Nightbot"
		:data="nightbotData?.nightbotGetData"
		:logout="logout"
		:authLink="authLinkData?.nightbotGetAuthLink"
		:icon="IconNightbot"
		:is-loading="dataFetching"
		with-settings
	>
		<template #description>
			<i18n-t keypath="integrations.nightbot.info" />
		</template>

		<template #settings>
			<div class="flex flex-col w-full gap-4">
				<UiCard class="flex flex-col flex-1">
					<UiCardHeader>
						<UiCardTitle>Commands</UiCardTitle>
					</UiCardHeader>
					<UiCardContent class="flex-1">
						<div v-if="commandsResponse">
							<p>Imported Count: {{ commandsResponse.importedCount }}</p>
							<p>Failed Count: {{ commandsResponse.failedCount }}</p>
							<p v-if="commandsResponse.failedCommandsNames.length > 0">Failed Commands:</p>
							<ul
								v-if="commandsResponse.failedCommandsNames.length > 0"
								class="overflow-y-scroll max-h-60"
							>
								<li v-for="name in commandsResponse.failedCommandsNames" :key="name">
									{{ name }}
								</li>
							</ul>
						</div>
						<UiAlert v-else>
							<InfoIcon class="size-4" />
							<UiAlertDescription>Waiting import...</UiAlertDescription>
						</UiAlert>
					</UiCardContent>
					<UiCardFooter>
						<UiButton
							class="w-full"
							:disabled="
								!isNightbotIntegrationEnabled || !userCanManageCommands || commandsImporting
							"
							@click="importCommands"
						>
							<LoaderCircleIcon v-if="commandsImporting" class="animate-spin size-4 mr-2" />
							Import
						</UiButton>
					</UiCardFooter>
				</UiCard>

				<UiCard class="flex flex-col flex-1">
					<UiCardHeader>
						<UiCardTitle>Timers</UiCardTitle>
					</UiCardHeader>
					<UiCardContent class="flex-1">
						<div v-if="timersResponse">
							<p>Imported Count: {{ timersResponse.importedCount }}</p>
							<p>Failed Count: {{ timersResponse.failedCount }}</p>
							<p v-if="timersResponse.failedTimersNames.length > 0">Failed Timers:</p>
							<ul
								v-if="timersResponse.failedTimersNames.length > 0"
								class="overflow-y-scroll max-h-60"
							>
								<li v-for="name in timersResponse.failedTimersNames" :key="name">
									{{ name }}
								</li>
							</ul>
						</div>
						<UiAlert v-else>
							<InfoIcon class="size-4" />
							<UiAlertDescription>Waiting import...</UiAlertDescription>
						</UiAlert>
					</UiCardContent>
					<UiCardFooter>
						<UiButton
							class="w-full"
							:disabled="!isNightbotIntegrationEnabled || !userCanManageTimers || timersImporting"
							@click="importTimers"
						>
							<LoaderCircleIcon v-if="timersImporting" class="animate-spin size-4 mr-2" />
							Import
						</UiButton>
					</UiCardFooter>
				</UiCard>
			</div>
		</template>
	</oauth-component>
</template>
