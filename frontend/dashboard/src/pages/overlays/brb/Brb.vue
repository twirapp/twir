<script setup lang="ts">
import { FontSelector } from '@twir/fontsource'
import {
	NAlert,
	NButton,
	NColorPicker,
	NDivider,
	NInput,
	NInputNumber,
	NSwitch,
	useNotification,
	useThemeVars,
} from 'naive-ui'
import { computed, ref, toRaw, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import type { Settings } from '@twir/api/messages/overlays_be_right_back/overlays_be_right_back'
import type { Font } from '@twir/fontsource'

import { useBeRightBackOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js'
import CommandButton from '@/features/commands/ui/command-button.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const themeVars = useThemeVars()
const { t } = useI18n()

const { data: profile } = useProfile()

const defaultSettings = {
	backgroundColor: 'rgba(9, 8, 8, 0.50)',
	fontColor: '#fff',
	fontFamily: 'inter',
	fontSize: 100,
	text: 'AFK FOR',
	late: {
		text: 'LATE FOR',
		displayBrbTime: true,
		enabled: true,
	},
	opacity: 50,
}

const formValue = ref<Settings>(structuredClone(defaultSettings))

const manager = useBeRightBackOverlayManager()
const {
	data: settings,
	isError: isSettingsError,
	isLoading: isSettingsLoading,
} = manager.getSettings()
const updater = manager.updateSettings()

watch(settings, (v) => {
	if (!v) return
	formValue.value = toRaw(v)
}, { immediate: true })

const brbIframeRef = ref<HTMLIFrameElement | null>(null)
const brbIframeUrl = computed(() => {
	if (!profile.value) return null

	return `${window.location.origin}/overlays/${profile.value.apiKey}/brb`
})

function sendIframeMessage(key: string, data?: any) {
	if (!brbIframeRef.value) return
	const win = brbIframeRef.value

	win.contentWindow?.postMessage(JSON.stringify({
		key,
		data: toRaw(data),
	}))
}

function sendSettings() {
	sendIframeMessage('settings', {
		...toRaw(formValue.value),
		channelName: profile.value?.login,
		channelId: profile.value?.id,
	})
}

watch(brbIframeRef, (v) => {
	if (!v) return

	v.contentWindow?.addEventListener('message', (e) => {
		const parsed = JSON.parse(e.data)
		if (parsed.key !== 'getSettings') return

		sendSettings()
	})
})

watch(() => formValue, () => {
	if (!brbIframeRef.value) return

	sendSettings()
}, { deep: true })

const { copyOverlayLink } = useCopyOverlayLink('brb')

const message = useNotification()

async function save() {
	await updater.mutateAsync(formValue.value)

	message.success({
		title: t('sharedTexts.saved'),
		duration: 5000,
	})
}

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays
})

function setDefaultSettings() {
	formValue.value = structuredClone(defaultSettings)
}

// TODO: fontWeight and fontStyle should be a select
const fontData = ref<Font | null>(null)
watch(() => fontData.value, (font) => {
	if (!font) return
	formValue.value.fontFamily = font.id
})
</script>

<template>
	<div class="page">
		<div class="card">
			<div class="card-header">
				<NButton
					secondary
					type="error"
					@click="setDefaultSettings"
				>
					{{ t('sharedButtons.setDefaultSettings') }}
				</NButton>
				<NButton
					secondary
					type="info"
					:disabled="isSettingsError || isSettingsLoading || !canCopyLink"
					@click="copyOverlayLink()"
				>
					{{ t('overlays.copyOverlayLink') }}
				</NButton>
				<NButton
					secondary
					type="success"
					@click="save"
				>
					{{ t('sharedButtons.save') }}
				</NButton>
			</div>

			<div class="card-body">
				<div class="card-body-column">
					<NDivider class="!my-0">
						{{ t('overlays.brb.settings.main.label') }}
					</NDivider>

					<div class="item">
						<div class="flex flex-col gap-1">
							<CommandButton
								name="brb"
								:title="t('overlays.brb.settings.main.startCommand.description')"
							/>
							<NAlert type="info" :show-icon="false">
								<span v-html="t('overlays.brb.settings.main.startCommand.example')" />
							</NAlert>
						</div>

						<CommandButton
							name="brbstop"
							:title="t('overlays.brb.settings.main.stopCommand.description')"
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.text') }}</span>
						<NInput v-model:value="formValue.text" :maxlength="500" />
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.background') }}</span>
						<NColorPicker
							v-model:value="formValue.backgroundColor" :modes="['rgb']"
							show-preview
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.color') }}</span>
						<NColorPicker
							v-model:value="formValue.fontColor" :modes="['hex', 'rgb']"
							:show-alpha="false"
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.family') }}</span>
						<FontSelector
							v-model:font="fontData"
							:font-family="formValue.fontFamily"
							font-style="normal"
							:font-weight="400"
						/>
					</div>

					<div class="item">
						<span>{{ t('overlays.brb.settings.main.font.size') }}</span>
						<NInputNumber v-model:value="formValue.fontSize" :min="1" :max="500" />
					</div>
				</div>

				<div class="card-body-column">
					<NDivider class="!my-0">
						{{ t('overlays.brb.settings.late.label') }}
					</NDivider>

					<div class="item">
						<span>{{ t('overlays.brb.settings.late.text') }}</span>
						<NInput v-model:value="formValue.late!.text" :maxlength="500" />
					</div>

					<div class="flex gap-2">
						<NSwitch v-model:value="formValue.late!.enabled" />
						<span>{{ t('sharedTexts.enabled') }}</span>
					</div>

					<div class="flex gap-2">
						<NSwitch v-model:value="formValue.late!.displayBrbTime" />
						<span>{{ t('overlays.brb.settings.late.displayBrb') }}</span>
					</div>
				</div>
			</div>
		</div>

		<div>
			<iframe
				v-if="brbIframeUrl"
				ref="brbIframeRef"
				:src="brbIframeUrl"
				class="iframe"
				border="0"
			/>
			<div class="absolute top-9 right-10 font-medium">
				<div class="flex gap-2">
					<NButton secondary size="small" type="warning" @click="sendIframeMessage('stop')">
						{{ t('overlays.brb.preview.stop') }}
					</NButton>
					<NButton
						secondary
						size="small"
						type="success"
						@click="() => {
							sendSettings();
							sendIframeMessage('start', { minutes: 0.1 })
						}"
					>
						{{ t('overlays.brb.preview.start') }}
					</NButton>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
@import '../styles.css';

.card {
	background-color: v-bind('themeVars.cardColor');
}

.iframe {
	border: 1px solid v-bind('themeVars.borderColor');
}
</style>
