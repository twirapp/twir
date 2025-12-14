<script setup lang="ts">
import { ValorantStatsWidget } from '@twir/frontend-valorant-stats'
import { BanIcon } from 'lucide-vue-next'
import { Label } from 'radix-vue'

import { useValorantStats } from './composables/use-valorant-stats'

import { useValorantIntegration } from '@/features/integrations/composables/valorant/use-valorant-integration.ts'
import { Button } from '@/components/ui/button'
import { ColorPicker } from '@/components/ui/color-picker'
import InputWithIcon from '@/components/ui/InputWithIcon.vue'
import Separator from '@/components/ui/separator/Separator.vue'
import { SwitchToggle } from '@/components/ui/switch'
import PageLayout from '@/layout/page-layout.vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const { settings, copyOverlayLink } = useValorantStats()

const { isConfigured: isConnected, isDataFetching } = useValorantIntegration()
</script>

<template>
	<PageLayout cleanBody>
		<template #title> {{ t('overlays.valorant.title') }} </template>

		<template #content>
			<div
				v-if="isDataFetching"
				class="flex items-center justify-center p-20 text-2xl bg-yellow-900/30 m-40 rounded-xl"
			>
				Loading...
			</div>
			<div class="relative">
				<div
					v-if="!isConnected"
					class="flex flex-col gap-4 items-center justify-center p-4 text bg-black/70 rounded-xl absolute inset-0 z-50"
				>
					<BanIcon class="size-10" />
					<span class="text-2xl"> {{ t('overlays.valorant.valorantIntegration.connect') }} </span>
				</div>

				<div
					class="flex flex-col-reverse lg:flex-row"
					:class="{ 'blur-xs pointer-events-none': !isConnected }"
				>
					<div class="overflow-auto bg-card lg:max-w-[400px] w-full flex flex-col p-4 shadow-md">
						<div class="flex flex-col gap-2">
							<span class="text-xs mb-2">{{ t('overlays.valorant.settings.colors.title') }}</span>

							<div class="flex flex-col gap-2">
								<Label for="backgroundColor">
									{{ t('overlays.valorant.settings.colors.background') }}
								</Label>
								<InputWithIcon id="backgroundColor" v-model="settings.backgroundColor">
									<ColorPicker id="backgroundColor" v-model:modelValue="settings.backgroundColor" />
								</InputWithIcon>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="textColor">{{ t('overlays.valorant.settings.colors.text') }}</Label>
								<InputWithIcon id="textColor" v-model="settings.textColor">
									<ColorPicker id="textColor" v-model:modelValue="settings.textColor" />
								</InputWithIcon>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="primaryTextColor">
									{{ t('overlays.valorant.settings.colors.primaryText') }}
								</Label>
								<InputWithIcon id="primaryTextColor" v-model="settings.primaryTextColor">
									<ColorPicker
										id="primaryTextColor"
										v-model:modelValue="settings.primaryTextColor"
									/>
								</InputWithIcon>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="winColor">{{ t('overlays.valorant.settings.colors.win') }}</Label>
								<InputWithIcon id="winColor" v-model="settings.winColor">
									<ColorPicker id="winColor" v-model:modelValue="settings.winColor" />
								</InputWithIcon>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="loseColor">{{ t('overlays.valorant.settings.colors.lose') }}</Label>
								<InputWithIcon id="loseColor" v-model="settings.loseColor">
									<ColorPicker id="loseColor" v-model:modelValue="settings.loseColor" />
								</InputWithIcon>
							</div>

							<Separator class="my-2" />

							<span class="text-xs mb-2">{{ t('overlays.valorant.settings.general.title') }}</span>

							<div class="flex flex-col gap-2">
								<Label for="disabledBorder">
									{{ t('overlays.valorant.settings.general.border') }}
								</Label>
								<SwitchToggle
									id="disabledBorder"
									:modelValue="!settings.disabledBorder"
									@update:model-value="(v) => (settings.disabledBorder = !v)"
								/>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="disabledGlowEffect">
									{{ t('overlays.valorant.settings.general.glow') }}
								</Label>
								<SwitchToggle
									id="disabledGlowEffect"
									:modelValue="!settings.disabledGlowEffect"
									@update:model-value="(v) => (settings.disabledGlowEffect = !v)"
								/>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="disabledLeaderboardPlace">
									{{ t('overlays.valorant.settings.general.leaderboard') }}
								</Label>
								<SwitchToggle
									id="disabledLeaderboardPlace"
									:modelValue="!settings.disabledLeaderboardPlace"
									@update:model-value="(v) => (settings.disabledLeaderboardPlace = !v)"
								/>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="disabledWinLose">
									{{ t('overlays.valorant.settings.general.winLose') }}
								</Label>
								<SwitchToggle
									id="disabledWinLose"
									:modelValue="!settings.disabledWinLose"
									@update:model-value="(v) => (settings.disabledWinLose = !v)"
								/>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="disabledProgress">
									{{ t('overlays.valorant.settings.general.progress') }}
								</Label>
								<SwitchToggle
									id="disabledProgress"
									:modelValue="!settings.disabledProgress"
									@update:model-value="(v) => (settings.disabledProgress = !v)"
								/>
							</div>

							<div class="flex flex-col gap-2">
								<Label for="lastTwentyMatches">
									{{ t('overlays.valorant.settings.general.last20MatchesStats') }}
								</Label>
								<SwitchToggle
									id="lastTwentyMatches"
									:modelValue="!settings.disabledTwentyLastMatches"
									@update:model-value="(v) => (settings.disabledTwentyLastMatches = !v)"
								/>
							</div>
						</div>

						<Button class="mt-4" @click="copyOverlayLink">
							{{ t('overlays.generateObsLink') }}
						</Button>
					</div>

					<div class="flex min-h-[200px] w-full h-full items-center justify-center">
						<ValorantStatsWidget class="w-[50%]" :settings="settings" />
					</div>
				</div>
			</div>
		</template>
	</PageLayout>
</template>
