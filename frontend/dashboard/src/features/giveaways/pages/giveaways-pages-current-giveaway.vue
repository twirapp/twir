<script setup lang="ts">
import { BanIcon, PlayIcon, Settings2Icon, ShuffleIcon } from 'lucide-vue-next'
import { computed, ref } from 'vue'

import { useProfile } from '@/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useGiveaways } from '@/features/giveaways/composables/giveaways-use-api.ts'
import GiveawaysCurrentGiveawayParticipants from '@/features/giveaways/ui/giveaways-current-giveaway/giveaways-current-giveaway-participants.vue'

const { participants } = useGiveaways()

const { data: profile } = useProfile()
const chatUrl = computed(() => {
	if (!profile?.value?.selectedDashboardTwitchUser) return

	const user = profile.value.selectedDashboardTwitchUser

	const url = `https://www.twitch.tv/embed/${user.login}/chat?parent=${window.location.host}&darkpopout`

	// if (chatTheme.value === 'dark') {
	// 	url += '&darkpopout'
	// }
	//
	// if (openFrankerFaceZ.value) {
	// 	url += '&ffz-settings'
	// }

	return url
})

const currentState = ref('created')
async function setNewState(state: string) {
	currentState.value = state
}
</script>

<template>
	<div class="flex flex-row flex-wrap-reverse gap-4 h-[calc(100dvh-200px)]">
		<Card class="flex-1 h-full">
			<CardHeader class="flex-row items-center p-2 justify-between border-b border-border border-solid">
				<CardTitle class="flex flex-col">
					<span>
						#giveaway
					</span>
					<span class="text-xs font-normal text-muted-foreground">
						{{ participants.length }} participants
					</span>
				</CardTitle>

				<div class="flex flex-row gap-1">
					<Button size="sm" variant="outline" class="flex gap-2 items-center">
						<Settings2Icon class="size-4" />
					</Button>
					<Button size="sm" variant="secondary" class="flex gap-2 items-center">
						<ShuffleIcon class="size-4" />
						Choose winners
					</Button>
					<Button v-if="currentState === 'created'" size="sm" class="flex gap-2 items-center" @click="setNewState('started')">
						<PlayIcon class="size-4" />
						Start
					</Button>
					<Button v-if="currentState === 'started'" size="sm" variant="outline" class="flex gap-2 items-center">
						<BanIcon class="size-4" />
						Finish
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
</template>
