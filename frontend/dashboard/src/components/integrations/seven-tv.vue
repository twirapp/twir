<script setup lang="ts">
import { IconCircleCheck, IconPlugOff } from '@tabler/icons-vue';
import type {
	UpdateDataRequest,
} from '@twir/api/messages/integrations_seventv/integrations_seventv';
import {
	NTag,
	NTimeline,
	NTimelineItem,
	NA,
	NSpin,
	useThemeVars,
	NAlert,
	NForm,
	NFormItem,
	NSpace,
	NText,
} from 'naive-ui';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import SevenTvButtonEditors from './seventv-button-editors.png';

import SevenTVSvg from '@/assets/integrations/seventv.svg?use';
import { useSevenTv } from '@/components/integrations/use-seven-tv';
import WithSettings from '@/components/integrations/variants/withSettings.vue';
import RewardsSelector from '@/components/rewardsSelector.vue';

const themeVars = useThemeVars();
const { t } = useI18n();

const sevenTvStore = useSevenTv();
const { isNotRegistered, data: sevenTvData, sevenTvProfileLink } = storeToRefs(sevenTvStore);

const form = ref<UpdateDataRequest>({});

const router = useRouter();

function goToEvents() {
	router.push('/dashboard/events/custom');
}

async function saveSettings() {
	await sevenTvStore.save(form.value);
}
</script>

<template>
	<with-settings
		title="7TV"
		:save="saveSettings"
		:icon="SevenTVSvg"
		icon-width="48px"
	>
		<template #description>
			{{ t('integrations.sevenTv.description') }}
		</template>

		<template #settings>
			<template v-if="isNotRegistered">
				<n-alert type="error">
					<i18n-t keypath="integrations.sevenTv.notRegistered">
						<n-a href="https://7tv.io" target="_blank">
							7tv.app
						</n-a>
					</i18n-t>
				</n-alert>
			</template>

			<template v-else>
				<n-spin :show="!sevenTvData?.isEditor">
					<n-form>
						<n-alert type="info" style="margin-bottom: 10px;">
							<i18n-t keypath="integrations.sevenTv.alert">
								<n-a @click="goToEvents">
									{{ t('sidebar.events').toLocaleLowerCase() }}
								</n-a>
							</i18n-t>
						</n-alert>

						<n-form-item :label="t('integrations.sevenTv.rewardForAddEmote')">
							<n-space vertical>
								<rewards-selector v-model="form.rewardIdForAddEmote" clearable />
								<n-text :depth="3" style="font-size: 12px">
									{{ t('integrations.sevenTv.rewardSelectorDescription') }}
								</n-text>
							</n-space>
						</n-form-item>

						<n-form-item :label="t('integrations.sevenTv.rewardForRemoveEmote')">
							<n-space vertical>
								<rewards-selector v-model="form.rewardIdForRemoveEmote" clearable />
								<n-text :depth="3" style="font-size: 12px">
									{{ t('integrations.sevenTv.rewardSelectorDescription') }}
								</n-text>
							</n-space>
						</n-form-item>
					</n-form>

					<template #description>
						<n-timeline>
							<n-timeline-item>
								<i18n-t keypath="integrations.sevenTv.connectSteps.step1">
									<n-a :href="sevenTvProfileLink" target="_blank">
										7tv
									</n-a>
								</i18n-t>
							</n-timeline-item>
							<n-timeline-item>
								<div style="display: flex; flex-direction: column">
									<span>{{ t('integrations.sevenTv.connectSteps.step2') }}</span>
									<img :src="SevenTvButtonEditors" height="50" width="100" />
								</div>
							</n-timeline-item>
							<n-timeline-item>
								<i18n-t keypath="integrations.sevenTv.connectSteps.step3">
									<b :style="{color: themeVars.successColor}">{{ sevenTvData?.botSeventvProfile?.username }}</b>
								</i18n-t>
							</n-timeline-item>
						</n-timeline>
					</template>
				</n-spin>
			</template>
		</template>

		<template #additionalFooter>
			<n-tag
				style="padding: 20px;"
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
			</n-tag>
		</template>
	</with-settings>
</template>

<style scoped>

</style>
