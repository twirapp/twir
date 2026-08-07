<script setup lang="ts">
import type { Settings } from '@twir/frontend-valorant-stats'
import { watch } from 'vue'

import { ColorPicker } from '@/components/ui/color-picker'
import InputWithIcon from '@/components/ui/InputWithIcon/InputWithIcon.vue'
import { Label } from '@/components/ui/label'
import { SwitchToggle } from '@/components/ui/switch'

import { useValorantStats } from '~~/layers/dashboard/features/overlays/valorant-stats/composables/use-valorant-stats'

interface Props {
	readonly params?: Record<string, string>
}

const props = defineProps<Props>()
const emit = defineEmits<{
	'update-params': [params: Record<string, string>]
}>()

const { settings, setSettings } = useValorantStats()
const { t } = useI18n()

function booleanParam(value: string | undefined, fallback: boolean) {
	return value === undefined ? fallback : value === 'true'
}

function paramsToSettings(params: Record<string, string> | undefined): Required<Settings> {
	return {
		backgroundColor: params?.backgroundColor ?? settings.value.backgroundColor,
		textColor: params?.textColor ?? settings.value.textColor,
		primaryTextColor: params?.primaryTextColor ?? settings.value.primaryTextColor,
		winColor: params?.winColor ?? settings.value.winColor,
		loseColor: params?.loseColor ?? settings.value.loseColor,
		disabledPeakRR: booleanParam(params?.disabledPeakRR, settings.value.disabledPeakRR),
		disabledLeaderboardPlace: booleanParam(params?.disabledLeaderboardPlace, settings.value.disabledLeaderboardPlace),
		disabledPeakRankIcon: booleanParam(params?.disabledPeakRankIcon, settings.value.disabledPeakRankIcon),
		disabledBorder: booleanParam(params?.disabledBorder, settings.value.disabledBorder),
		disabledWinLose: booleanParam(params?.disabledWinLose, settings.value.disabledWinLose),
		disabledProgress: booleanParam(params?.disabledProgress, settings.value.disabledProgress),
		disabledGlowEffect: booleanParam(params?.disabledGlowEffect, settings.value.disabledGlowEffect),
		disabledTwentyLastMatches: booleanParam(
			params?.disabledTwentyLastMatches,
			settings.value.disabledTwentyLastMatches,
		),
	}
}

function settingsToParams(value: Required<Settings>): Record<string, string> {
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
		<p class="text-sm text-muted-foreground">{{ t('overlayBuilder.widgets.valorantStats.settingsHint') }}</p>
		<div class="flex flex-col gap-2">
			<Label for="valorant-background">{{ t('overlays.valorant.settings.colors.background') }}</Label>
			<InputWithIcon id="valorant-background" v-model="settings.backgroundColor">
				<ColorPicker v-model="settings.backgroundColor" />
			</InputWithIcon>
		</div>

		<div class="flex flex-col gap-2">
			<Label for="valorant-text">{{ t('overlays.valorant.settings.colors.text') }}</Label>
			<InputWithIcon id="valorant-text" v-model="settings.textColor">
				<ColorPicker v-model="settings.textColor" />
			</InputWithIcon>
		</div>

		<div class="flex flex-col gap-2">
			<Label for="valorant-primary-text">{{ t('overlays.valorant.settings.colors.primaryText') }}</Label>
			<InputWithIcon id="valorant-primary-text" v-model="settings.primaryTextColor">
				<ColorPicker v-model="settings.primaryTextColor" />
			</InputWithIcon>
		</div>

		<div class="flex flex-col gap-2">
			<Label for="valorant-win">{{ t('overlays.valorant.settings.colors.win') }}</Label>
			<InputWithIcon id="valorant-win" v-model="settings.winColor">
				<ColorPicker v-model="settings.winColor" />
			</InputWithIcon>
		</div>

		<div class="flex flex-col gap-2">
			<Label for="valorant-lose">{{ t('overlays.valorant.settings.colors.lose') }}</Label>
			<InputWithIcon id="valorant-lose" v-model="settings.loseColor">
				<ColorPicker v-model="settings.loseColor" />
			</InputWithIcon>
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="valorant-border">{{ t('overlays.valorant.settings.general.border') }}</Label>
			<SwitchToggle id="valorant-border" :model-value="!settings.disabledBorder" @update:model-value="settings.disabledBorder = !$event" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="valorant-glow">{{ t('overlays.valorant.settings.general.glow') }}</Label>
			<SwitchToggle id="valorant-glow" :model-value="!settings.disabledGlowEffect" @update:model-value="settings.disabledGlowEffect = !$event" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="valorant-leaderboard">{{ t('overlays.valorant.settings.general.leaderboard') }}</Label>
			<SwitchToggle id="valorant-leaderboard" :model-value="!settings.disabledLeaderboardPlace" @update:model-value="settings.disabledLeaderboardPlace = !$event" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="valorant-win-lose">{{ t('overlays.valorant.settings.general.winLose') }}</Label>
			<SwitchToggle id="valorant-win-lose" :model-value="!settings.disabledWinLose" @update:model-value="settings.disabledWinLose = !$event" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="valorant-progress">{{ t('overlays.valorant.settings.general.progress') }}</Label>
			<SwitchToggle id="valorant-progress" :model-value="!settings.disabledProgress" @update:model-value="settings.disabledProgress = !$event" />
		</div>

		<div class="flex items-center justify-between gap-4">
			<Label for="valorant-last-matches">{{ t('overlays.valorant.settings.general.last20MatchesStats') }}</Label>
			<SwitchToggle id="valorant-last-matches" :model-value="!settings.disabledTwentyLastMatches" @update:model-value="settings.disabledTwentyLastMatches = !$event" />
		</div>
	</div>
</template>
