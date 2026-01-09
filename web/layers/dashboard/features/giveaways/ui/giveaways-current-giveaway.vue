<script setup lang="ts">
import { BanIcon, PlayIcon, ShuffleIcon, TrophyIcon, UsersIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'


// ...existing code...
import { useGiveaways } from '~/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCurrentGiveawayParticipants from '~/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue'
import GiveawaysCurrentGiveawayWinners from '~/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners.vue'

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

const { user: profile } = storeToRefs(useDashboardAuth())

const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
		?.twitchProfile
})

const chatUrl = computed(() => {
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
		<UiCard class="flex-1 h-full min-h-0">
			<UiTabs v-model="activeTab" class="h-full flex flex-col">
				<UiCardHeader
					class="flex flex-row items-center justify-between border-b border-border border-solid"
				>
					<UiCardTitle class="flex items-center text-3xl">
						<span>
							{{ currentGiveaway?.keyword }}
						</span>
					</UiCardTitle>
					<div class="ml-2 flex flex-row gap-1">
						<UiButton
							v-if="!currentGiveaway?.startedAt"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStartGiveaway"
						>
							<PlayIcon class="size-4" />
							{{ t('giveaways.start') }}
						</UiButton>

						<UiButton
							v-if="!currentGiveaway?.stoppedAt && currentGiveaway?.startedAt"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStopGiveaway"
						>
							<BanIcon class="size-4" />
							{{ t('giveaways.stop') }}
						</UiButton>

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
						<UiTabsList>
							<UiTabsTrigger value="participants" class="flex flex-row gap-2">
								<UsersIcon class="size-4 inline" />
								{{ t('giveaways.currentGiveaway.tabs.participants') }}
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ participants.length }}
								</span>
							</UiTabsTrigger>
							<UiTabsTrigger value="winners" class="flex flex-row gap-2">
								<TrophyIcon class="size-4 inline" />
								<span>
									{{ t('giveaways.currentGiveaway.tabs.winners') }}
								</span>
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ winners.length }}
								</span>
							</UiTabsTrigger>
						</UiTabsList>
					</div>
				</UiCardHeader>
				<UiCardContent class="h-[calc(100%-56px)] min-h-0 p-0">
					<UiTabsContent value="participants" class="h-full mt-0 border-none">
						<GiveawaysCurrentGiveawayParticipants />
					</UiTabsContent>

					<UiTabsContent value="winners" class="mt-0 h-full border-none">
						<div class="flex flex-col h-full">
							<div class="p-2 border-b flex justify-between flex-wrap gap-2 items-center">
								<span class="text-sm font-medium">{{
									t('giveaways.currentGiveaway.totalWinners', { count: winners.length })
								}}</span>
								<UiButton
									size="sm"
									variant="secondary"
									class="flex gap-2 items-center"
									:disabled="participants.length === 0 || winners.length === participants.length"
									@click="handleChooseWinners"
								>
									<ShuffleIcon class="size-4" />
									{{ t('giveaways.chooseWinner') }}
								</UiButton>
							</div>
							<GiveawaysCurrentGiveawayWinners />
						</div>
					</UiTabsContent>
				</UiCardContent>
			</UiTabs>
		</UiCard>
		<UiCard class="flex-1 h-full">
			<UiCardContent class="p-0 h-full">
				<iframe v-if="chatUrl" :src="chatUrl" frameborder="0" class="w-full h-full" />
			</UiCardContent>
		</UiCard>
	</div>
</template>
