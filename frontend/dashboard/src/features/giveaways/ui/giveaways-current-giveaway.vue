<script setup lang="ts">
import { BanIcon, PlayIcon, ShuffleIcon, TrophyIcon, UsersIcon } from 'lucide-vue-next';
import { computed, ref, watch, watchEffect } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile } from '@/api/auth';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-giveaways.ts';
import GiveawaysCurrentGiveawayParticipants
	from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue';
import GiveawaysCurrentGiveawayWinners
	from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-winners.vue';
import { Badge } from '@/components/ui/badge';

const { t } = useI18n();

const {
	participants,
	currentGiveaway,
	currentGiveawayId,
	startGiveaway,
	stopGiveaway,
	chooseWinners,
	winners,
} = useGiveaways();

const { data: profile } = useProfile();

const selectedDashboardTwitchUser = computed(() => {
	return profile.value?.availableDashboards.find((d) => d.id === profile.value?.selectedDashboardId)
		?.twitchProfile;
});

const chatUrl = computed(() => {
	if (!selectedDashboardTwitchUser.value) return;

	return `https://www.twitch.tv/embed/${selectedDashboardTwitchUser.value.login}/chat?parent=${window.location.host}&darkpopout`;
});

const isOnlineChatterGiveaway = computed(() => currentGiveaway.value?.type === "ONLINE_CHATTERS");

// Tab state - для ONLINE_CHATTERS открываем сразу winners
const activeTab = ref("participants");

// Устанавливаем правильную вкладку при загрузке или смене гива
watchEffect(() => {
	if (isOnlineChatterGiveaway.value) {
		activeTab.value = "winners";
	} else if (winners.value.length > 0) {
		activeTab.value = "winners";
	} else {
		activeTab.value = "participants";
	}
});

// Watch for winners and switch to winners tab when they are chosen
watch(winners, (newWinners) => {
	if (newWinners.length > 0) {
		activeTab.value = "winners";
	}
});

async function handleStartGiveaway() {
	if (currentGiveawayId.value) {
		await startGiveaway(currentGiveawayId.value);
	}
}

async function handleStopGiveaway() {
	if (currentGiveawayId.value) {
		await stopGiveaway(currentGiveawayId.value);
	}
}

async function handleChooseWinners() {
	if (currentGiveawayId.value) {
		await chooseWinners(currentGiveawayId.value);
	}
}
</script>

<template>
	<div class="flex flex-row flex-wrap-reverse gap-4 h-[90dvh] p-4">
		<Card class="flex-1 h-full min-h-0">
			<Tabs v-model="activeTab" class="h-full flex flex-col">
				<CardHeader
					class="flex flex-row items-center justify-between border-b border-border border-solid"
				>
					<CardTitle class="flex flex-col gap-1">
						<div class="flex items-center gap-2 text-3xl">
							<span v-if="currentGiveaway?.type === 'KEYWORD'">
								{{ currentGiveaway?.keyword }}
							</span>
							<span v-else>
								{{ t("giveaways.typeOnlineChatters") }}
							</span>
						</div>
						<div class="text-sm font-normal text-muted-foreground flex flex-col flex-wrap gap-2">
							<span>
								{{ t("giveaways.type") }}:
								{{
									currentGiveaway?.type === "KEYWORD"
										? t("giveaways.typeKeyword")
										: t("giveaways.typeOnlineChatters")
								}}
							</span>
							<div class="flex flex-wrap gap-2">
								<Badge class="bg-blue-500 text-white" v-if="currentGiveaway?.minWatchedTime">
									{{ t("giveaways.minWatchedTime") }}: {{ currentGiveaway.minWatchedTime }}m
								</Badge>
								<Badge class="bg-blue-500 text-white" v-if="currentGiveaway?.minMessages">
									{{ t("giveaways.minMessages") }}: {{ currentGiveaway.minMessages }}
								</Badge>
								<Badge class="bg-blue-500 text-white" v-if="currentGiveaway?.minUsedChannelPoints"
									>{{ t("giveaways.minUsedChannelPoints") }}:
									{{ currentGiveaway.minUsedChannelPoints }}
								</Badge>
								<Badge class="bg-blue-500 text-white" v-if="currentGiveaway?.minFollowDuration">
									{{ t("giveaways.minFollowDuration") }}:
									{{ currentGiveaway.minFollowDuration }}d
								</Badge>
								<Badge class="bg-blue-500 text-white" v-if="currentGiveaway?.requireSubscription"
								>{{ t("giveaways.requireSubscription") }}
								</Badge>
							</div>
						</div>
					</CardTitle>

					<div class="flex gap-2 place-self-start items-center justify-center">
					<div class="ml-2 flex flex-row gap-1">
						<!-- Для KEYWORD гивов показываем кнопку Start/Stop -->
						<Button
							v-if="!isOnlineChatterGiveaway && !currentGiveaway?.startedAt"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStartGiveaway"
						>
							<PlayIcon class="size-4" />
							{{ t("giveaways.start") }}
						</Button>

						<Button
							v-if="!currentGiveaway?.stoppedAt && currentGiveaway?.startedAt"
							size="sm"
							class="flex gap-2 items-center"
							@click="handleStopGiveaway"
						>
							<BanIcon class="size-4" />
							{{ t("giveaways.stop") }}
						</Button>

						<!-- Для ONLINE_CHATTERS показываем сразу кнопку выбора победителя -->
						<Button
							v-if="isOnlineChatterGiveaway && !currentGiveaway?.stoppedAt"
							size="sm"
							variant="secondary"
							class="flex gap-2 items-center"
							@click="handleChooseWinners"
						>
							<ShuffleIcon class="size-4" />
							{{ t("giveaways.chooseWinner") }}
						</Button>
					</div>
						<TabsList>
							<TabsTrigger
								value="participants"
								class="flex flex-row gap-2"
								:disabled="isOnlineChatterGiveaway"
							>
								<UsersIcon class="size-4 inline" />
								{{ t("giveaways.currentGiveaway.tabs.participants") }}
								<span class="ml-1 rounded-full bg-primary text-primary-foreground text-xs px-2">
									{{ participants.length }}
								</span>
							</TabsTrigger>
							<TabsTrigger value="winners" class="flex flex-row gap-2">
								<TrophyIcon class="size-4 inline" />
								<span>
									{{ t("giveaways.currentGiveaway.tabs.winners") }}
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
									t("giveaways.currentGiveaway.totalWinners", { count: winners.length })
								}}</span>
								<Button
									size="sm"
									variant="secondary"
									class="flex gap-2 items-center"
									:disabled="
										!isOnlineChatterGiveaway &&
										(participants.length === 0 || winners.length === participants.length)
									"
									@click="handleChooseWinners"
								>
									<ShuffleIcon class="size-4" />
									{{ t("giveaways.chooseWinner") }}
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
