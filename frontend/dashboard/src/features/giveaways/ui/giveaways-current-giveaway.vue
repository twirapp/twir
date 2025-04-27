<script setup lang="ts">
import { ArchiveIcon, BanIcon, PlayIcon, ShuffleIcon } from 'lucide-vue-next'
import { computed } from 'vue'

import { useProfile } from '@/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts'
import GiveawaysCurrentGiveawayParticipants from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue'
import GiveawaysCurrentGiveawayWinners from '@/features/giveaways/ui/giveaways-current-giveaway-winners.vue'

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
	return currentGiveaway.value?.startedAt && !currentGiveaway.value?.stoppedAt && !currentGiveaway.value?.endedAt
})

const isStopped = computed(() => {
	return currentGiveaway.value?.stoppedAt
})

const isEnded = computed(() => {
	return currentGiveaway.value?.endedAt
})

const isArchived = computed(() => {
	return currentGiveaway.value?.archivedAt
})
</script>

<template>
	<div class="flex flex-col gap-4 h-[calc(100dvh-200px)] p-4">
		<div v-if="winners.length > 0" class="flex-none">
			<GiveawaysCurrentGiveawayWinners />
		</div>

		<div class="flex flex-row flex-wrap-reverse gap-4 flex-1">
			<Card class="flex-1 h-full">
				<CardHeader class="flex-row items-center p-2 justify-between border-b border-border border-solid">
					<CardTitle class="flex flex-col">
						<span>
							{{ currentGiveaway?.keyword }}
						</span>
						<span class="text-xs font-normal text-muted-foreground">
							{{ participants.length }} participants
						</span>
					</CardTitle>

					<div class="flex flex-row gap-1">
						<!--						<Button size="sm" variant="outline" class="flex gap-2 items-center"> -->
						<!--							<Settings2Icon class="size-4" /> -->
						<!--						</Button> -->

						<Button
							size="sm"
							variant="secondary"
							class="flex gap-2 items-center"
							:disabled="!isActive || participants.length === 0"
							@click="handleChooseWinners"
						>
							<ShuffleIcon class="size-4" />
							Choose winners
						</Button>

						<Button
							v-if="!isActive && !isStopped && !isEnded && !isArchived"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStartGiveaway"
						>
							<PlayIcon class="size-4" />
							Start
						</Button>

						<Button
							v-if="isActive"
							size="sm"
							variant="outline"
							class="flex gap-2 items-center"
							@click="handleStopGiveaway"
						>
							<BanIcon class="size-4" />
							Stop
						</Button>

						<Button
							v-if="!isArchived"
							size="sm"
							variant="destructive"
							class="flex gap-2 items-center"
							@click="handleArchiveGiveaway"
						>
							<ArchiveIcon class="size-4" />
							Archive
						</Button>
					</div>
				</CardHeader>
				<CardContent class="h-full p-0">
					<div class="flex flex-row gap-4 h-full">
						<GiveawaysCurrentGiveawayParticipants />
					</div>
				</CardContent>
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
	</div>
</template>
