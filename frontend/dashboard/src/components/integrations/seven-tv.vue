<script setup lang="ts">
import { IconCircleCheck, IconPlugOff } from '@tabler/icons-vue'
import {
	NA,
	NAlert,
	NForm,
	NFormItem,
	NSpace,
	NSpin,
	NSwitch,
	NTabPane,
	NTabs,
	NTag,
	NText,
	NTimeline,
	NTimelineItem,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import SevenTvButtonEditors from './seventv-button-editors.png'

import type {
	UpdateDataRequest,
} from '@twir/api/messages/integrations_seventv/integrations_seventv'

import { useCommandsApi } from '@/api/commands/commands'
import SevenTVSvg from '@/assets/integrations/seventv.svg?use'
import { useSevenTv } from '@/components/integrations/use-seven-tv'
import WithSettings from '@/components/integrations/variants/withSettings.vue'
import RewardsSelector from '@/components/rewardsSelector.vue'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import CommandsList from '@/features/commands/components/list.vue'

const { t } = useI18n()

const { notification } = useNaiveDiscrete()
const {
	isNotRegistered,
	data: sevenTvData,
	sevenTvProfileLink,
	save,
} = useSevenTv()

const form = ref<UpdateDataRequest>({
	deleteEmotesOnlyAddedByApp: true,
})
watch(sevenTvData, (data) => {
	if (!data) return
	form.value = {
		rewardIdForAddEmote: data.rewardIdForAddEmote,
		rewardIdForRemoveEmote: data.rewardIdForRemoveEmote,
		deleteEmotesOnlyAddedByApp: data.deleteEmotesOnlyAddedByApp,
	}
})

const router = useRouter()

function goToEvents() {
	router.push('/dashboard/events/custom')
}

const isSameRewardsChoosed = computed(() => {
	if (!form.value.rewardIdForAddEmote || !form.value.rewardIdForRemoveEmote) return false

	return form.value.rewardIdForAddEmote === form.value.rewardIdForRemoveEmote
})

async function saveSettings() {
	try {
		await save(form.value)
		notification.success({ title: t('sharedTexts.saved'), duration: 2500 })
	} catch (err) {
		notification.error({ title: t('sharedTexts.errorOnSave'), duration: 2500 })
	}
}

const commandsManager = useCommandsApi()
const { data: commandsData } = commandsManager.useQueryCommands()
const commands = computed(() => {
	return commandsData.value?.commands.filter((c) => c.module === '7tv') ?? []
})
</script>

<template>
	<WithSettings
		title="7TV"
		:save="saveSettings"
		:icon="SevenTVSvg"
		icon-width="48px"
		:save-disabled="isSameRewardsChoosed"
		modal-content-style="padding: 0px"
	>
		<template #description>
			{{ t('integrations.sevenTv.description') }}
		</template>

		<template #settings>
			<template v-if="isNotRegistered || !sevenTvData?.emoteSetId">
				<NAlert v-if="isNotRegistered" type="error">
					<i18n-t keypath="integrations.sevenTv.notRegistered">
						<NA href="https://7tv.app" target="_blank">
							7tv.app
						</NA>
					</i18n-t>
				</NAlert>
				<NAlert v-else type="error">
					Emote set not created on 7tv, please create at least one set on
					<NA :href="sevenTvProfileLink" target="_blank">
						7tv
					</NA>
				</NAlert>
			</template>

			<template v-else>
				<NTabs
					default-value="settings"
					justify-content="space-evenly"
					type="line"
					pane-style="padding: 10px;"
				>
					<NTabPane name="settings" :tab="t('overlays.tts.tabs.general')">
						<NSpin :show="!sevenTvData?.isEditor">
							<NForm>
								<NFormItem :label="t('integrations.sevenTv.rewardForAddEmote')">
									<NSpace vertical>
										<RewardsSelector v-model="form.rewardIdForAddEmote" only-with-input clearable />
										<NText :depth="3" class="text-xs">
											{{ t('integrations.sevenTv.rewardSelectorDescription') }}
										</NText>
									</NSpace>
								</NFormItem>

								<NFormItem :label="t('integrations.sevenTv.rewardForRemoveEmote')">
									<NSpace vertical>
										<RewardsSelector v-model="form.rewardIdForRemoveEmote" only-with-input clearable />
										<NText :depth="3" class="text-xs">
											{{ t('integrations.sevenTv.rewardSelectorDescription') }}
										</NText>

										<div class="flex gap-1">
											<span>{{ t('integrations.sevenTv.deleteOnlyAddedByApp') }}</span>
											<NSwitch v-model:value="form.deleteEmotesOnlyAddedByApp" />
										</div>
									</NSpace>
								</NFormItem>
							</NForm>

							<div class="flex flex-col gap-1">
								<NAlert v-if="isSameRewardsChoosed" type="error">
									{{ t('integrations.sevenTv.errorSameReward') }}
								</NAlert>

								<NAlert type="info" class="mb-2.5">
									<i18n-t keypath="integrations.sevenTv.alert">
										<NA @click="goToEvents">
											{{ t('sidebar.events').toLocaleLowerCase() }}
										</NA>
									</i18n-t>
								</NAlert>
							</div>

							<template #description>
								<NTimeline>
									<NTimelineItem>
										<i18n-t keypath="integrations.sevenTv.connectSteps.step1">
											<NA :href="sevenTvProfileLink" target="_blank">
												7tv
											</NA>
										</i18n-t>
									</NTimelineItem>
									<NTimelineItem>
										<div class="flex flex-col">
											<span>{{ t('integrations.sevenTv.connectSteps.step2') }}</span>
											<img :src="SevenTvButtonEditors" height="50" width="100" />
										</div>
									</NTimelineItem>
									<NTimelineItem>
										<i18n-t keypath="integrations.sevenTv.connectSteps.step3">
											<b classs="text-[color:var(--n-color-target)]">{{ sevenTvData?.botSeventvProfile?.username }}</b>
										</i18n-t>
									</NTimelineItem>
								</NTimeline>
							</template>
						</NSpin>
					</NTabPane>

					<NTabPane name="commands" :tab="t('sidebar.commands.label')">
						<CommandsList :commands show-background />
					</NTabPane>
				</NTabs>
			</template>
		</template>

		<template #additionalFooter>
			<NTag
				class="!p-5"
				:bordered="false"
				:type="sevenTvData?.isEditor ? 'success' : 'error'"
			>
				<template #icon>
					<IconCircleCheck v-if="sevenTvData?.isEditor" />
					<IconPlugOff v-else />
				</template>
				<template v-if="sevenTvData?.isEditor">
					{{ t('integrations.sevenTv.connected') }}
				</template>
				<template v-else>
					{{ t('integrations.sevenTv.notConnected') }}
				</template>
			</NTag>
		</template>
	</WithSettings>
</template>

<style>
:deep(.n-card__content) {
	padding: 0px;
}
</style>
