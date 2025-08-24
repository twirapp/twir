<script setup lang="ts">
import { ValorantStatsWidget } from '@twir/frontend-valorant-stats'
import { Label } from 'radix-vue'

import { useValorantStats } from './composables/use-valorant-stats'

import { Button } from '@/components/ui/button'
import { ColorPicker } from '@/components/ui/color-picker'
import Separator from '@/components/ui/separator/Separator.vue'
import { SwitchToggle } from '@/components/ui/switch'
import PageLayout from '@/layout/page-layout.vue'

const { settings, copyOverlayLink } = useValorantStats()
</script>

<template>
	<PageLayout cleanBody>
		<template #title>
			Valorant Stats Builder
		</template>

		<template #content>
			<div class="flex flex-col-reverse lg:flex-row">
				<div class="overflow-auto bg-card lg:max-w-[400px] w-full flex flex-col p-4 shadow-md">
					<div class="flex flex-col gap-2">
						<span class="text-xs mb-2">Colors</span>

						<div class="flex flex-col gap-2">
							<Label for="backgroundColor">Background color</Label>
							<ColorPicker v-model:modelValue="settings.backgroundColor" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="textColor">Text color</Label>
							<ColorPicker v-model:modelValue="settings.textColor" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="primaryTextColor">Primary text color</Label>
							<ColorPicker v-model:modelValue="settings.primaryTextColor" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="winColor">Win color</Label>
							<ColorPicker v-model:modelValue="settings.winColor" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="loseColor">Lose color</Label>
							<ColorPicker v-model:modelValue="settings.loseColor" />
						</div>

						<Separator class="my-2" />

						<span class="text-xs mb-2">General</span>

						<div class="flex flex-col gap-2">
							<Label for="disabledBackground">Background</Label>
							<SwitchToggle id="disabledBackground" :modelValue="!settings.disabledBackground" @update:model-value="v => settings.disabledBackground = !v" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="disabledBorder">Border</Label>
							<SwitchToggle id="disabledBorder" :modelValue="!settings.disabledBorder" @update:model-value="v => settings.disabledBorder = !v" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="disabledGlowEffect">Glow effect</Label>
							<SwitchToggle id="disabledGlowEffect" :modelValue="!settings.disabledGlowEffect" @update:model-value="v => settings.disabledGlowEffect = !v" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="disabledLeaderboardPlace">Leaderboard place</Label>
							<SwitchToggle id="disabledLeaderboardPlace" :modelValue="!settings.disabledLeaderboardPlace" @update:model-value="v => settings.disabledLeaderboardPlace = !v" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="disabledWinLose">Win / Lose</Label>
							<SwitchToggle id="disabledWinLose" :modelValue="!settings.disabledWinLose" @update:model-value="v => settings.disabledWinLose = !v" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="disabledProgress">Progress</Label>
							<SwitchToggle id="disabledProgress" :modelValue="!settings.disabledProgress" @update:model-value="v => settings.disabledProgress = !v" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="lastTwentyMatches">Show last 20 matches stats</Label>
							<SwitchToggle
								id="lastTwentyMatches"
								:modelValue="!settings.disabledTwentyLastMatches"
								@update:model-value="v => settings.disabledTwentyLastMatches = !v"
							/>
						</div>
					</div>

					<Button class="mt-4" @click="copyOverlayLink">
						Generate obs link
					</Button>
				</div>

				<div class="flex min-h-[200px] w-full h-full items-center justify-center">
					<ValorantStatsWidget class="w-[50%]" :settings="settings" />
				</div>
			</div>
		</template>
	</PageLayout>
</template>
