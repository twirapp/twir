<script setup lang="ts">
import { IconCircleCheck, IconPlugOff } from '@tabler/icons-vue';
import { NTag, NTimeline, NTimelineItem, NA, NSpin, useThemeVars } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import SevenTvButtonEditors from './seventv-button-editors.png';

import { useSevenTvIntegration } from '@/api/integrations/seventv';
import SevenTVSvg from '@/assets/integrations/seventv.svg?use';
import WithSettings from '@/components/integrations/variants/withSettings.vue';

const themeVars = useThemeVars();
const { t } = useI18n();
const manager = useSevenTvIntegration();
const { data: sevenTvData } = manager.useData();
</script>

<template>
	<with-settings
		title="7TV"
		:icon="SevenTVSvg"
		icon-width="48px"
	>
		<template #description>
			{{ t('integrations.sevenTv.description') }}
		</template>

		<template #settings>
			<n-timeline>
				<n-timeline-item>
					Go to
					<n-a href="https://7tv.app" target="_blank">
						7tv.app
					</n-a>
				</n-timeline-item>
				<n-timeline-item>Open your profile</n-timeline-item>
				<n-timeline-item>
					<div style="display: flex; flex-direction: column">
						<span>Click on "Add editors" button</span>
						<img :src="SevenTvButtonEditors" height="50" width="100" />
					</div>
				</n-timeline-item>
				<n-timeline-item>
					Add <b :style="{color: themeVars.successColor}">{{ sevenTvData?.botUsername }}</b> as an
					editor
				</n-timeline-item>
				<n-timeline-item :type="sevenTvData?.isEditor ? 'success' : 'error'">
					<template v-if="sevenTvData?.isEditor">
						Done
					</template>
					<template v-else>
						<n-spin size="small" />
						Waiting until you add the bot as an editor
					</template>
				</n-timeline-item>
			</n-timeline>
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
