<script setup lang="ts">
import { computed, ref, watch, watchEffect } from 'vue'
import { useProfile } from '~~/layers/dashboard/api/auth'
import { useGiveaways } from '~~/layers/dashboard/features/giveaways/composables/giveaways-use-giveaways.js'
import GiveawaysCurrentGiveawayParticipants from '~~/layers/dashboard/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue'
import GiveawaysCurrentGiveawayWinners from '~~/layers/dashboard/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners.vue'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'

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
const requestUrl = useRequestURL()

const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
		?.twitchProfile
})

const chatUrl = computed(() => {
	if (!selectedDashboardTwitchUser.value) return

	return `https://www.twitch.tv/embed/${selectedDashboardTwitchUser.value.login}/chat?parent=${requestUrl.host}&darkpopout`
})

const isOnlineChatterGiveaway = computed(() => currentGiveaway.value?.type === 'ONLINE_CHATTERS')

// Tab state - для ONLINE_CHATTERS открываем сразу winners
const activeTab = ref('participants')

// Устанавливаем правильную вкладку при загрузке или смене гива
watchEffect(() => {
	if (isOnlineChatterGiveaway.value) {
		activeTab.value = 'winners'
	} else if (winners.value.length > 0) {
		activeTab.value = 'winners'
	} else {
		activeTab.value = 'participants'
	}
})

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
	<div class="flex h-[90dvh] flex-row flex-wrap-reverse gap-4 p-4">
		<Card class="h-full min-h-0 flex-1">
			<Tabs
				v-model="activeTab"
				class="flex h-full flex-col"
			>
				<CardHeader
					class="border-border flex flex-row items-center justify-between border-b border-solid"
				>
					<CardTitle class="flex flex-col gap-1">
						<div class="flex items-center gap-2 text-3xl">
							<span v-if="currentGiveaway?.type === 'KEYWORD'">
								{{ currentGiveaway?.keyword }}
							</span>
							<span v-else>
								{{ t('giveaways.typeOnlineChatters') }}
							</span>
						</div>
						<div class="text-muted-foreground flex flex-col flex-wrap gap-2 text-sm font-normal">
							<span>
								{{ t('giveaways.type') }}:
								{{
									currentGiveaway?.type === 'KEYWORD'
										? t('giveaways.typeKeyword')
										: t('giveaways.typeOnlineChatters')
								}}
							</span>
							<div class="flex flex-wrap gap-2">
								<Badge
									class="bg-blue-500 text-white"
									v-if="currentGiveaway?.minWatchedTime"
								>
									{{ t('giveaways.minWatchedTime') }}: {{ currentGiveaway.minWatchedTime }}m
								</Badge>
								<Badge
									class="bg-blue-500 text-white"
									v-if="currentGiveaway?.minMessages"
								>
									{{ t('giveaways.minMessages') }}: {{ currentGiveaway.minMessages }}
								</Badge>
								<Badge
									class="bg-blue-500 text-white"
									v-if="currentGiveaway?.minUsedChannelPoints"
									>{{ t('giveaways.minUsedChannelPoints') }}:
									{{ currentGiveaway.minUsedChannelPoints }}
								</Badge>
								<Badge
									class="bg-blue-500 text-white"
									v-if="currentGiveaway?.minFollowDuration"
								>
									{{ t('giveaways.minFollowDuration') }}: {{ currentGiveaway.minFollowDuration }}d
								</Badge>
								<Badge
									class="bg-blue-500 text-white"
									v-if="currentGiveaway?.requireSubscription"
									>{{ t('giveaways.requireSubscription') }}
								</Badge>
							</div>
						</div>
					</CardTitle>

					<div class="flex items-center justify-center gap-2 place-self-start">
						<div class="ml-2 flex flex-row gap-1">
							<!-- Для KEYWORD гивов показываем кнопку Start/Stop -->
							<Button
								v-if="!isOnlineChatterGiveaway && !currentGiveaway?.startedAt"
								size="sm"
								class="flex items-center gap-2"
								@click="handleStartGiveaway"
							>
								<Icon
									name="lucide:play"
									class="size-4"
								/>
								{{ t('giveaways.start') }}
							</Button>

							<Button
								v-if="!currentGiveaway?.stoppedAt && currentGiveaway?.startedAt"
								size="sm"
								class="flex items-center gap-2"
								@click="handleStopGiveaway"
							>
								<Icon
									name="lucide:ban"
									class="size-4"
								/>
								{{ t('giveaways.stop') }}
							</Button>

							<!-- Для ONLINE_CHATTERS показываем сразу кнопку выбора победителя -->
							<Button
								v-if="isOnlineChatterGiveaway && !currentGiveaway?.stoppedAt"
								size="sm"
								variant="secondary"
								class="flex items-center gap-2"
								@click="handleChooseWinners"
							>
								<Icon
									name="lucide:shuffle"
									class="size-4"
								/>
								{{ t('giveaways.chooseWinner') }}
							</Button>
						</div>
						<TabsList>
							<TabsTrigger
								value="participants"
								class="flex flex-row gap-2"
								:disabled="isOnlineChatterGiveaway"
							>
								<Icon
									name="lucide:users"
									class="inline size-4"
								/>
								{{ t('giveaways.currentGiveaway.tabs.participants') }}
								<span class="bg-primary text-primary-foreground ml-1 rounded-full px-2 text-xs">
									{{ participants.length }}
								</span>
							</TabsTrigger>
							<TabsTrigger
								value="winners"
								class="flex flex-row gap-2"
							>
								<Icon
									name="lucide:trophy"
									class="inline size-4"
								/>
								<span>
									{{ t('giveaways.currentGiveaway.tabs.winners') }}
								</span>
								<span class="bg-primary text-primary-foreground ml-1 rounded-full px-2 text-xs">
									{{ winners.length }}
								</span>
							</TabsTrigger>
						</TabsList>
					</div>
				</CardHeader>
				<CardContent class="h-[calc(100%-56px)] min-h-0 p-0">
					<TabsContent
						value="participants"
						class="mt-0 h-full border-none"
					>
						<GiveawaysCurrentGiveawayParticipants />
					</TabsContent>

					<TabsContent
						value="winners"
						class="mt-0 h-full border-none"
					>
						<div class="flex h-full flex-col">
							<div class="flex flex-wrap items-center justify-between gap-2 border-b p-2">
								<span class="text-sm font-medium">{{
									t('giveaways.currentGiveaway.totalWinners', { count: winners.length })
								}}</span>
								<Button
									size="sm"
									variant="secondary"
									class="flex items-center gap-2"
									:disabled="
										!isOnlineChatterGiveaway &&
										(participants.length === 0 || winners.length === participants.length)
									"
									@click="handleChooseWinners"
								>
									<Icon
										name="lucide:shuffle"
										class="size-4"
									/>
									{{ t('giveaways.chooseWinner') }}
								</Button>
							</div>
							<GiveawaysCurrentGiveawayWinners />
						</div>
					</TabsContent>
				</CardContent>
			</Tabs>
		</Card>
		<Card class="h-full flex-1">
			<CardContent class="h-full p-0">
				<iframe
					v-if="chatUrl"
					:src="chatUrl"
					frameborder="0"
					class="h-full w-full"
				/>
			</CardContent>
		</Card>
	</div>
</template>
