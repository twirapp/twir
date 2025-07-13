<script setup lang="ts">
import { BanIcon, PlayIcon, ShuffleIcon, TrophyIcon, UsersIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCurrentGiveawayParticipants from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue'
import GiveawaysCurrentGiveawayWinners from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners.vue'

const { t } = useI18n()

const {
	participants,
	currentGiveaway,
	currentGiveawayId,
	startGiveaway,
	stopGiveaway,
	chooseWinners,
	winners,
} = useGiveaways()

const { data: profile } = useProfile()

const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
		?.twitchProfile
})

const chatUrl = computed(() => {
	console.log(selectedDashboardTwitchUser.value)
	if (!selectedDashboardTwitchUser.value) return

	return `https://www.twitch.tv/embed/${selectedDashboardTwitchUser.value.login}/chat?parent=${window.location.host}&darkpopout`
})

// Tab state
const activeTab = ref('participants')

// Watch for winners and switch to winners tab when they are chosen
watch(winners, (newWinners) => {
	if (newWinners.length > 0) {
		activeTab.value = 'winners'
	}
})

async function handleStartGiveaway() {
	if (currentGiveawayId.value) {
		await startGiveaway(currentGiveawayId.value)
	}
}

async function handleStopGiveaway() {
	if (currentGiveawayId.value) {
		await stopGiveaway(currentGiveawayId.value)
	}
}

async function handleChooseWinners() {
	if (currentGiveawayId.value) {
		await chooseWinners(currentGiveawayId.value)
	}
}
</script>

<template>
	<div class="flex flex-row flex-wrap-reverse gap-4 h-[98dvh] p-4">
		<Card class="flex-1 h-full min-h-0">
			<Tabs v-model="activeTab" class="h-full flex flex-col">
				<CardHeader
					class="flex-row items-center p-2 justify-between border-b border-border border-solid"
				>
					<CardTitle class="flex items-center">
						<span>
							{{ currentGiveaway?.keyword }}
						</span>
					</CardTitle>
					<div class="ml-2 flex flex-row gap-1">
						<Button
							v-if="!currentGiveaway?.startedAt"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStartGiveaway"
						>
							<PlayIcon class="size-4" />
							{{ t('giveaways.start') }}
						</Button>

						<Button
							v-if="!currentGiveaway?.stoppedAt && currentGiveaway?.startedAt"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStopGiveaway"
						>
							<BanIcon class="size-4" />
							{{ t('giveaways.stop') }}
						</Button>

						<!--						<Button -->
						<!--							v-if="!archived" -->
						<!--							size="sm" -->
						<!--							variant="destructive" -->
						<!--							class="flex gap-2 items-center" -->
						<!--							@click="handleArchiveGiveaway" -->
						<!--						> -->
						<!--							<ArchiveIcon class="size-4" /> -->
						<!--							Archive -->
						<!--						</Button> -->
					</div>

					<div>
						<TabsList>
							<TabsTrigger value="participants" class="flex flex-row gap-2">
								<UsersIcon class="size-4 inline" />
								{{ t('giveaways.currentGiveaway.tabs.participants') }}
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ participants.length }}
								</span>
							</TabsTrigger>
							<TabsTrigger value="winners" class="flex flex-row gap-2">
								<TrophyIcon class="size-4 inline" />
								<span>
									{{ t('giveaways.currentGiveaway.tabs.winners') }}
								</span>
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ winners.length }}
								</span>
							</TabsTrigger>
						</TabsList>
					</div>
				</CardHeader>
				<CardContent class="h-[calc(100%-56px)] min-h-0 p-0">
					<TabsContent value="participants" class="h-full mt-0 border-none">
						<GiveawaysCurrentGiveawayParticipants />
					</TabsContent>

					<TabsContent value="winners" class="mt-0 h-full border-none">
						<div class="flex flex-col h-full">
							<div class="p-2 border-b flex justify-between flex-wrap gap-2 items-center">
								<span class="text-sm font-medium">{{
									t('giveaways.currentGiveaway.totalWinners', { count: winners.length })
								}}</span>
								<Button
									size="sm"
									variant="secondary"
									class="flex gap-2 items-center"
									:disabled="participants.length === 0 || winners.length === participants.length"
									@click="handleChooseWinners"
								>
									<ShuffleIcon class="size-4" />
									{{ t('giveaways.chooseWinner') }}
								</Button>
							</div>
							<GiveawaysCurrentGiveawayWinners />
						</div>
					</TabsContent>
				</CardContent>
			</Tabs>
		</Card>
		<Card class="flex-1 h-full">
			<CardContent class="p-0 h-full">
				<iframe v-if="chatUrl" :src="chatUrl" frameborder="0" class="w-full h-full" />
			</CardContent>
		</Card>
	</div>
</template>
