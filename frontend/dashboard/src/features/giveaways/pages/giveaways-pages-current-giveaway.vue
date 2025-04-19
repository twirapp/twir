<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'
import { computed, ref } from 'vue'

import { useProfile } from '@/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

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

const participants = ref<number[]>(
	Array.from({ length: 100 }, (_, i) => i),
)

useIntervalFn(() => {
	participants.value.push(participants.value.length + 1)
}, 1000, { immediate: true })
</script>

<template>
	<div class="flex flex-row flex-wrap-reverse gap-4 h-[calc(100dvh-200px)]">
		<Card class="flex-1 h-full">
			<CardHeader class="p-2">
				<CardTitle>
					Manage {{ participants.length }}
				</CardTitle>
			</CardHeader>
			<CardContent>
				<div class="flex flex-row gap-4">
					<RecycleScroller
						v-slot="{ item }"
						class="scroller"
						:items="participants"
						:item-size="32"
					>
						<div>
							{{ item.name }}
						</div>
					</RecycleScroller>
					<div class="flex flex-1 flex-col h-full">
						asd
					</div>
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
				>
				</iframe>
			</CardContent>
		</Card>
	</div>
</template>
