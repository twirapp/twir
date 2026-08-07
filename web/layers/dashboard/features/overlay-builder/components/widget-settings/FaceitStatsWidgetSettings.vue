<script setup lang="ts">
import type { Settings } from '@twir/frontend-faceit-stats'
import { watch } from 'vue'

import { ColorPicker } from '@/components/ui/color-picker'
import InputWithIcon from '@/components/ui/InputWithIcon/InputWithIcon.vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { SwitchToggle } from '@/components/ui/switch'

import { useFaceitStats } from '~~/layers/dashboard/features/overlays/faceit-stats/composables/use-faceit-stats'

interface Props {
	readonly params?: Record<string, string>
}

const props = defineProps<Props>()
const emit = defineEmits<{
	'update-params': [params: Record<string, string>]
}>()

const { settings, setSettings } = useFaceitStats()
const { t } = useI18n()

function booleanParam(value: string | undefined, fallback: boolean) {
	return value === undefined ? fallback : value === 'true'
}

function paramsToSettings(params: Record<string, string> | undefined): Settings {
	return {
		nickname: params?.nickname ?? settings.value.nickname,
		game: 'cs2',
		bgColor: params?.bgColor ?? settings.value.bgColor,
		textColor: params?.textColor ?? settings.value.textColor,
		borderRadius: params?.borderRadius ? Number(params.borderRadius) : settings.value.borderRadius,
		displayAvarageKdr: booleanParam(params?.displayAvarageKdr, settings.value.displayAvarageKdr),
		displayWorldRanking: booleanParam(params?.displayWorldRanking, settings.value.displayWorldRanking),
		displayLastTwentyMatches: booleanParam(
			params?.displayLastTwentyMatches,
			settings.value.displayLastTwentyMatches,
		),
	}
}

function settingsToParams(value: Settings): Record<string, string> {
	return Object.fromEntries(Object.entries(value).map(([key, setting]) => [key, String(setting)]))
}

watch(
	() => props.params,
	(params) => setSettings(paramsToSettings(params)),
	{ immediate: true, deep: true },
)

watch(
	settings,
	(value) => emit('update-params', settingsToParams(value)),
	{ deep: true },
)
</script>

<template>
	<div class="flex flex-col gap-4">
		<p class="text-sm text-muted-foreground">{{ t('overlayBuilder.widgets.faceitStats.settingsHint') }}</p>
		<div class="flex flex-col gap-2">
			<Label for="faceit-nickname">{{ t('overlays.faceit.settings.general.faceitNickname') }}</Label>
			<Input id="faceit-nickname" v-model="settings.nickname" placeholder="s1mple" />
		</div>

		<div class="flex flex-col gap-2">
			<Label for="faceit-game">{{ t('overlays.faceit.settings.general.game') }}</Label>
			<Input id="faceit-game" :model-value="settings.game" disabled />
		</div>

		<div class="flex flex-col gap-2">
			<Label for="faceit-background">{{ t('overlays.faceit.settings.appearance.background') }}</Label>
			<InputWithIcon id="faceit-background" v-model="settings.bgColor">
				<ColorPicker v-model="settings.bgColor" />
			</InputWithIcon>
		</div>

		<div class="flex flex-col gap-2">
			<Label for="faceit-text">{{ t('overlays.faceit.settings.appearance.textColor') }}</Label>
			<InputWithIcon id="faceit-text" v-model="settings.textColor">
				<ColorPicker v-model="settings.textColor" />
			</InputWithIcon>
		</div>

		<div class="flex flex-col gap-2">
			<Label for="faceit-radius">{{ t('overlays.faceit.settings.appearance.borderRadius') }}</Label>
			<Input id="faceit-radius" v-model.number="settings.borderRadius" type="number" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="faceit-kdr">{{ t('overlays.faceit.settings.display.showAvarageKdr') }}</Label>
			<SwitchToggle id="faceit-kdr" v-model="settings.displayAvarageKdr" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="faceit-ranking">{{ t('overlays.faceit.settings.display.worldRanking') }}</Label>
			<SwitchToggle id="faceit-ranking" v-model="settings.displayWorldRanking" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="faceit-last-matches">{{ t('overlays.faceit.settings.display.last20MatchesStats') }}</Label>
			<SwitchToggle id="faceit-last-matches" v-model="settings.displayLastTwentyMatches" />
		</div>
	</div>
</template>
