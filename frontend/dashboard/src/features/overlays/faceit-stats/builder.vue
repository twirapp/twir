<script setup lang="ts">
import { FaceitStatsWidget } from '@twir/frontend-faceit-stats'
import { Label } from 'radix-vue'

import { useI18n } from 'vue-i18n'
import { useFaceitStats } from './composables/use-faceit-stats'
import DisplaySwitcher from './ui/display-switcher.vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import Separator from '@/components/ui/separator/Separator.vue'
import PageLayout from '@/layout/page-layout.vue'

const { settings, copyOverlayLink } = useFaceitStats()

const { t } = useI18n()
</script>

<template>
	<PageLayout cleanBody>
		<template #title>
			{{ t('overlays.faceit.title') }}
		</template>

		<template #content>
			<div class="flex flex-col-reverse lg:flex-row">
				<div class="overflow-auto bg-card lg:max-w-[400px] w-full flex flex-col p-4 shadow-md">
					<div class="flex flex-col gap-2">
						<span class="text-xs mb-2">{{ t('overlays.faceit.settings.general.title') }}</span>

						<div class="flex flex-col gap-2">
							<Label for="nickname">
								{{ t('overlays.faceit.settings.general.faceitNickname') }}
							</Label>
							<Input id="nickname" v-model:modelValue="settings.nickname" placeholder="s1mple" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="game">{{ t('overlays.faceit.settings.general.game') }}</Label>
							<Input id="game" v-model:modelValue="settings.game" placeholder="cs2" disabled />
						</div>

						<Separator class="my-2" />

						<span class="text-xs mb-2">{{ t('overlays.faceit.settings.appearance.title') }}</span>

						<div class="flex flex-col gap-2">
							<Label for="backgroundColor">
								{{ t('overlays.faceit.settings.appearance.background') }}
							</Label>
							<Input id="backgroundColor" v-model:modelValue="settings.bgColor" type="color" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="textColor">
								{{ t('overlays.faceit.settings.appearance.textColor') }}
							</Label>
							<Input id="textColor" v-model:modelValue="settings.textColor" type="color" />
						</div>

						<div class="flex flex-col gap-2">
							<Label for="borderRadius">
								{{ t('overlays.faceit.settings.appearance.borderRadius') }}
							</Label>
							<Input id="borderRadius" v-model:modelValue="settings.borderRadius" type="number" />
						</div>

						<Separator class="my-2" />

						<span class="text-xs mb-2">{{ t('overlays.faceit.settings.display.title') }}</span>

						<div class="flex flex-col gap-2">
							<Label for="showAvarageKdr">
								{{ t('overlays.faceit.settings.display.showAvarageKdr') }}
							</Label>
							<DisplaySwitcher
								id="showAvarageKdr"
								v-model:modelValue="settings.displayAvarageKdr"
							/>
						</div>

						<div class="flex flex-col gap-2">
							<Label for="showWorldRanking">
								{{ t('overlays.faceit.settings.display.worldRanking') }}
							</Label>
							<DisplaySwitcher
								id="showWorldRanking"
								v-model:modelValue="settings.displayWorldRanking"
							/>
						</div>

						<div class="flex flex-col gap-2">
							<Label for="last20MatchesStats">
								{{ t('overlays.faceit.settings.display.last20MatchesStats') }}
							</Label>
							<DisplaySwitcher
								id="lastTwentyatches"
								v-model:modelValue="settings.displayLastTwentyMatches"
							/>
						</div>
					</div>

					<Button class="mt-4" @click="copyOverlayLink">
						{{ t('overlays.faceit.settings.generateObsLink') }}
					</Button>
				</div>

				<div class="flex min-h-[200px] w-full h-full items-center justify-center">
					<FaceitStatsWidget :settings="settings" />
				</div>
			</div>
		</template>
	</PageLayout>
</template>
