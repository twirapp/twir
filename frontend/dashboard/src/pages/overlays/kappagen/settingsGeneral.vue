<script lang="ts" setup>
import { EmojiStyle } from '@twir/api/messages/overlays_kappagen/overlays_kappagen'
import { NAlert, NButton, NDivider, NSelect, NSlider, NSwitch } from 'naive-ui'
import { h } from 'vue'
import { useI18n } from 'vue-i18n'

import { useKappagenFormSettings } from './store.js'

import type { SelectBaseOption } from 'naive-ui/es/select/src/interface'
import type { VNodeChild } from 'vue'

import CommandButton from '@/features/commands/ui/command-button.vue'

const { settings: formValue } = useKappagenFormSettings()
const { t } = useI18n()

const formatSizeValue = (v: number) => Number.parseInt(`${v}`.split('.')[1])

const emojiStylesOptions: SelectBaseOption[] = [
	{ label: 'Disabled', value: EmojiStyle.None },
	{ label: 'Twemoji', value: EmojiStyle.Twemoji },
	{ label: 'Openmoji', value: EmojiStyle.Openmoji },
	{ label: 'Noto', value: EmojiStyle.Noto },
	{ label: 'Blob', value: EmojiStyle.Blobmoji },
]

function renderEmojiLabel(option: SelectBaseOption): VNodeChild {
	const style = emojiStylesOptions.find(s => s.value === option.value)
	if (!style || style.value === EmojiStyle.None) return 'Disabled'

	const preview = `https://cdn.frankerfacez.com/static/emoji/images/${style.label?.toString().toLowerCase()}/1f609.png`

	return [
		h(
			'div',
			{ class: 'flex items-center gap-1' },
			{
				default: () => [
					h('span', undefined, { default: () => option.label }),
					h('img', { class: 'h-5 w-5', src: preview }),
				],
			},
		),
	]
}
</script>

<template>
	<div class="tab">
		<NAlert type="info" :show-icon="false" class="mt-1">
			{{ t('overlays.kappagen.info') }}
		</NAlert>
		<CommandButton name="kappagen" />

		<div class="switch">
			<NSwitch v-model:value="formValue.enableSpawn" />
			<span>{{ t('overlays.kappagen.settings.spawn') }}</span>
		</div>

		<NDivider />

		<div class="slider">
			{{ t('overlays.kappagen.settings.size') }}({{ formatSizeValue(formValue.size!.ratioNormal) }})
			<NSlider
				v-model:value="formValue.size!.ratioNormal"
				:format-tooltip="formatSizeValue"
				:step="0.01"
				:min="0.05"
				:max="0.15"
			/>
		</div>

		<div class="slider">
			{{ t('overlays.kappagen.settings.sizeSmall') }}({{
				formatSizeValue(formValue.size!.ratioSmall)
			}})
			<NSlider
				v-model:value="formValue.size!.ratioSmall"
				:format-tooltip="formatSizeValue"
				:step="0.01"
				:min="0.02"
				:max="0.07"
			/>
		</div>

		<div class="switchers">
			<div class="switch">
				<NSwitch v-model:value="formValue.emotes!.bttvEnabled" />
				<span>{{ t('overlays.kappagen.settings.emotes.bttvEnabled') }}</span>
			</div>

			<div class="switch">
				<NSwitch v-model:value="formValue.emotes!.ffzEnabled" />
				<span>{{ t('overlays.kappagen.settings.emotes.ffzEnabled') }}</span>
			</div>

			<div class="switch">
				<NSwitch v-model:value="formValue.emotes!.sevenTvEnabled" />
				<span>{{ t('overlays.kappagen.settings.emotes.seventvEnabled') }}</span>
			</div>

			<div class="switch">
				<span>{{ t('overlays.kappagen.settings.emotes.emojiStyle') }}</span>
				<NSelect
					v-model:value="formValue.emotes!.emojiStyle"
					:options="emojiStylesOptions"
					style="width: 40%;"
					size="tiny"
					:render-label="renderEmojiLabel"
				/>
			</div>
		</div>

		<NDivider />

		<div class="slider">
			{{ t('overlays.kappagen.settings.time') }}({{ formValue.emotes!.time }}s)
			<NSlider
				v-model:value="formValue.emotes!.time"
				:min="1"
				:max="15"
			/>
		</div>

		<div class="slider">
			{{ t('overlays.kappagen.settings.maxEmotes') }}({{ formValue.emotes!.max }})
			<NSlider
				v-model:value="formValue.emotes!.max"
				:min="0"
				:max="250"
			/>
		</div>

		<NDivider />

		<div class="switchers">
			<span>{{ t('overlays.kappagen.settings.animationsOnAppear') }}</span>

			<div class="switch">
				<NSwitch v-model:value="formValue.animation!.fadeIn" />
				<span>Fade</span>
			</div>

			<div class="switch">
				<NSwitch v-model:value="formValue.animation!.zoomIn" />
				<span>Zoom</span>
			</div>
		</div>

		<NDivider />

		<div class="switchers">
			<span>{{ t('overlays.kappagen.settings.animationsOnDisappear') }}</span>

			<div class="switch">
				<NSwitch v-model:value="formValue.animation!.fadeOut" />
				<span>Fade</span>
			</div>

			<div class="switch">
				<NSwitch v-model:value="formValue.animation!.zoomOut" />
				<span>Zoom</span>
			</div>
		</div>

		<NDivider />

		<div class="switch">
			<NSwitch v-model:value="formValue.enableRave" />
			<span>{{ t('overlays.kappagen.settings.rave') }}</span>
		</div>

		<NDivider />

		<div class="flex flex-col gap-1">
			<span>{{ t('overlays.kappagen.settings.excludedEmotes') }}</span>

			<NSelect
				v-model:value="formValue.excludedEmotes"
				filterable
				multiple
				tag
				:placeholder="t('overlays.kappagen.settings.excludedEmotes')"
				:show-arrow="false"
				:show="false"
			/>
			<NButton secondary type="error" @click="formValue.excludedEmotes = []">
				Clear
			</NButton>
		</div>
	</div>
</template>
