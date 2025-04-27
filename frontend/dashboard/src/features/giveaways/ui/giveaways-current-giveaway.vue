<script setup lang="ts">
import { ArchiveIcon, BanIcon, PlayIcon, ShuffleIcon, TrophyIcon, UsersIcon } from 'lucide-vue-next'
import { computed, ref, watch } from 'vue'

import { useProfile } from '@/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCurrentGiveawayParticipants from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue'
import GiveawaysCurrentGiveawayWinners from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners.vue'

const {
	participants,
	currentGiveaway,
	currentGiveawayId,
	startGiveaway,
	stopGiveaway,
	archiveGiveaway,
	chooseWinners,
	winners,
} = useGiveaways()

const { data: profile } = useProfile()
const chatUrl = computed(() => {
	if (!profile?.value?.selectedDashboardTwitchUser) return

	const user = profile.value.selectedDashboardTwitchUser
	return `https://www.twitch.tv/embed/${user.login}/chat?parent=${window.location.host}&darkpopout`
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

async function handleArchiveGiveaway() {
	if (currentGiveawayId.value) {
		await archiveGiveaway(currentGiveawayId.value)
	}
}

async function handleChooseWinners() {
	if (currentGiveawayId.value) {
		await chooseWinners(currentGiveawayId.value)
	}
}

const isActive = computed(() => {
	return !!(currentGiveaway.value?.startedAt && !currentGiveaway.value?.stoppedAt && !currentGiveaway.value?.endedAt)
})

const stopped = computed(() => {
	return currentGiveaway.value?.stoppedAt
})

const ended = computed(() => {
	return currentGiveaway.value?.endedAt
})

const archived = computed(() => {
	return currentGiveaway.value?.archivedAt
})

const canBeRunned = computed(() => {
	return !currentGiveaway.value?.startedAt && !currentGiveaway.value?.stoppedAt && !currentGiveaway.value?.endedAt && !currentGiveaway.value?.archivedAt
})
</script>

<template>
	<div class="flex flex-row flex-wrap-reverse gap-4 h-[98dvh] p-4">
		<Card class="flex-1 h-full">
			<Tabs v-model="activeTab" class="h-full flex flex-col">
				<CardHeader class="flex-row items-center p-2 justify-between border-b border-border border-solid">
					<CardTitle class="flex items-center">
						<span>
							{{ currentGiveaway?.keyword }}
						</span>
					</CardTitle>
					<div class="ml-2 flex flex-row gap-1">
						<Button
							v-if="canBeRunned"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStartGiveaway"
						>
							<PlayIcon class="size-4" />
							Start
						</Button>

						<Button
							v-if="isActive && !stopped"
							size="sm"
							variant="outline"
							class="flex gap-2 items-center"
							@click="handleStopGiveaway"
						>
							<BanIcon class="size-4" />
							Stop
						</Button>

						<Button
							v-if="!archived"
							size="sm"
							variant="destructive"
							class="flex gap-2 items-center"
							@click="handleArchiveGiveaway"
						>
							<ArchiveIcon class="size-4" />
							Archive
						</Button>
					</div>

					<div>
						<TabsList>
							<TabsTrigger value="participants" class="flex flex-row gap-2">
								<UsersIcon class="size-4 inline" />
								Participants
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ participants.length }}
								</span>
							</TabsTrigger>
							<TabsTrigger value="winners" class="flex flex-row gap-2">
								<TrophyIcon class="size-4 inline" />
								<span>
									Winners
								</span>
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ winners.length }}
								</span>
							</TabsTrigger>
						</TabsList>
					</div>
				</CardHeader>
				<CardContent class="h-[calc(100%-56px)] p-0">
					<TabsContent value="participants" class="flex-1 mt-0 border-none">
						<GiveawaysCurrentGiveawayParticipants />
					</TabsContent>

					<TabsContent value="winners" class="flex-1 mt-0 border-none">
						<div class="flex flex-col h-full">
							<div class="p-2 border-b flex justify-between flex-wrap gap-2 items-center">
								<span class="text-sm font-medium">Total winners: {{ winners.length }}</span>
								<Button
									size="sm"
									variant="secondary"
									class="flex gap-2 items-center"
									:disabled="participants.length === 0 || winners.length == participants.length"
									@click="handleChooseWinners"
								>
									<ShuffleIcon class="size-4" />
									Choose winner
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
				<iframe
					v-if="chatUrl"
					:src="chatUrl"
					frameborder="0"
					class="w-full h-full"
				/>
			</CardContent>
		</Card>
	</div>
</template>
