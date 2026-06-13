<script lang="ts" setup>
import { ExternalLinkIcon, SettingsIcon } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import { CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'

import Card from './card.vue'
import Ban from './events/ban.vue'
import ChatClear from './events/chatClear.vue'
import Donate from './events/donate.vue'
import FirstUserMessage from './events/firstUserMessage.vue'
import Follow from './events/follow.vue'
import Raid from './events/raid.vue'
import RedemptionCreated from './events/redemptionCreated.vue'
import ReSubscribe from './events/resubscribe.vue'
import SubGift from './events/subgift.vue'
import Subscribe from './events/subscribe.vue'

import { useProfile } from '@/api/auth'
import UnbanRequestCreated from '@/components/dashboard/events/unban-request-created.vue'
import UnbanRequestResolved from '@/components/dashboard/events/unban-request-resolved.vue'
import { useEvents } from '@/features/dashboard/widgets/composables/events'
import { DashboardEventType } from '@/gql/graphql'

const props = defineProps<{
	popup?: boolean
}>()

const { t } = useI18n()
const { events, enabledEvents, enabledEventsOptions } = useEvents()

const { data } = useProfile()

function openPopup() {
	if (!data.value) return

	const height = 800
	const width = 500
	const top = Math.max(0, (screen.height - height) / 2)
	const left = Math.max(0, (screen.width - width) / 2)

	window.open(
		`${window.location.origin}/dashboard/popup/widgets/eventslist?apiKey=${data.value.apiKey}`,
		'_blank',
		`height=${height},width=${width},top=${top},left=${left},status=0,location=0,menubar=0,toolbar=0`,
	)
}
</script>

<template>
	<Card :popup="props.popup" class="flex flex-col">
		<template #header-extra>
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger as-child>
						<Button size="sm" variant="ghost" @click="openPopup">
							<ExternalLinkIcon class="h-4 w-4" />
						</Button>
					</TooltipTrigger>
					<TooltipContent>
						{{ t('sharedButtons.popout') }}
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>

			<Popover>
				<PopoverTrigger as-child>
					<Button variant="ghost" size="sm">
						<SettingsIcon class="h-4 w-4" />
					</Button>
				</PopoverTrigger>
				<PopoverContent class="w-64">
					<div class="space-y-3">
						<h4 class="font-medium text-sm">{{ t('dashboard.events.settings') || 'Event Settings' }}</h4>
						<div class="space-y-2">
							<div
								v-for="option in enabledEventsOptions"
								:key="option.value"
								class="flex items-center space-x-2"
							>
								<Checkbox
									:id="`event-${option.value}`"
									:checked="enabledEvents.includes(option.value)"
									@update:checked="(checked: boolean) => {
										if (checked) {
											enabledEvents.push(option.value)
										} else {
											const index = enabledEvents.indexOf(option.value)
											if (index > -1) enabledEvents.splice(index, 1)
										}
									}"
								/>
								<Label
									:for="`event-${option.value}`"
									class="text-sm font-normal cursor-pointer"
								>
									{{ option.label }}
								</Label>
							</div>
						</div>
					</div>
				</PopoverContent>
			</Popover>
		</template>
		<CardContent v-if="events.length" class="flex-1 overflow-hidden p-0">
			<ScrollArea class="h-full">
			<TransitionGroup name="list">
				<template v-for="(event) of events" :key="event.createdAt">
					<Follow
						v-if="event.type === DashboardEventType.Follow"
						:created-at="event.createdAt"
						:user-name="event.data.followUserName"
						:user-display-name="event.data.followUserDisplayName"
					/>
					<Raid
						v-if="event.type === DashboardEventType.Raided"
						:created-at="event.createdAt"
						:user-name="event.data.raidedFromUserName"
						:user-display-name="event.data.raidedFromDisplayName"
						:viewers="event.data.raidedViewersCount"
					/>
					<Donate
						v-if="event.type === DashboardEventType.Donation"
						:created-at="event.createdAt"
						:user-name="event.data.donationUserName"
						:amount="event.data.donationAmount"
						:message="event.data.donationMessage"
						:currency="event.data.donationCurrency"
					/>
					<Subscribe
						v-if="event.type === DashboardEventType.Subscribe"
						:created-at="event.createdAt"
						:user-name="event.data.subUserName"
						:user-display-name="event.data.subUserDisplayName"
						:level="event.data.subLevel"
					/>
					<ReSubscribe
						v-if="event.type === DashboardEventType.Resubscribe"
						:created-at="event.createdAt"
						:user-name="event.data.reSubUserName"
						:user-display-name="event.data.reSubUserDisplayName"
						:level="event.data.reSubLevel"
						:months="event.data.reSubMonths"
						:streak="event.data.reSubStreak"
					/>
					<SubGift
						v-if="event.type === DashboardEventType.Subgift"
						:created-at="event.createdAt"
						:user-name="event.data.subGiftTargetUserName"
						:user-display-name="event.data.subGiftUserDisplayName"
						:level="event.data.subGiftLevel"
						:target-user-name="event.data.subGiftTargetUserName"
						:target-user-display-name="event.data.subGiftTargetUserDisplayName"
					/>
					<FirstUserMessage
						v-if="event.type === DashboardEventType.FirstUserMessage"
						:created-at="event.createdAt"
						:user-name="event.data.firstUserMessageUserName"
						:user-display-name="event.data.firstUserMessageUserDisplayName"
						:message="event.data.firstUserMessageMessage"
					/>
					<ChatClear
						v-if="event.type === DashboardEventType.ChatClear"
						:created-at="event.createdAt"
					/>
					<RedemptionCreated
						v-if="event.type === DashboardEventType.RedemptionCreated"
						:created-at="event.createdAt"
						:title="event.data.redemptionTitle"
						:input="event.data.redemptionInput"
						:user-name="event.data.redemptionUserName"
						:user-display-name="event.data.redemptionUserDisplayName"
						:cost="event.data.redemptionCost"
					/>
					<Ban
						v-if="event.type === DashboardEventType.ChannelBan"
						:created-at="event.createdAt"
						:ends-in="event.data.banEndsInMinutes"
						:moderator-user-login="event.data.moderatorName"
						:moderator-user-name="event.data.moderatorDisplayName"
						:reason="event.data.banReason"
						:user-login="event.data.bannedUserLogin"
						:user-name="event.data.bannedUserName"
					/>
					<UnbanRequestCreated
						v-if="event.type === DashboardEventType.ChannelUnbanRequestCreate"
						:created-at="event.createdAt"
						:message="event.data.message"
						:user-login="event.data.userLogin"
						:user-name="event.data.userName"
					/>
					<UnbanRequestResolved
						v-if="event.type === DashboardEventType.ChannelUnbanRequestResolve"
						:created-at="event.createdAt"
						:message="event.data.message"
						:user-login="event.data.userLogin"
						:user-name="event.data.userName"
						:moderator-user-login="event.data.moderatorName"
						:moderator-user-name="event.data.moderatorDisplayName"
					/>
				</template>
			</TransitionGroup>
			</ScrollArea>
		</CardContent>
		<CardContent v-else class="flex items-center justify-center flex-1">
			<p class="text-4xl text-muted-foreground">
				No events
			</p>
		</CardContent>
	</Card>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
	transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
	opacity: 0;
	transform: translateY(30px);
}
</style>
