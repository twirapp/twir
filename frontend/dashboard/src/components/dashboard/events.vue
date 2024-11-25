<script lang="ts" setup>
import { IconExternalLink, IconSettings } from '@tabler/icons-vue'
import { NButton, NPopselect, NScrollbar, NText, NTooltip } from 'naive-ui'
import { useI18n } from 'vue-i18n'

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

import { useProfile } from '@/api/index.js'
import UnbanRequestCreated from '@/components/dashboard/events/unban-request-created.vue'
import UnbanRequestResolved from '@/components/dashboard/events/unban-request-resolved.vue'
import { useEvents } from '@/features/dashboard/widgets/composables/events'
import { DashboardEventType } from '@/gql/graphql'

const props = defineProps<{
	popup?: boolean
}>()

const { t } = useI18n()
const { events, fetching, enabledEvents, enabledEventsOptions } = useEvents()

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
	<Card :content-style="{ padding: fetching ? '10px' : '0px', height: '80%' }" :popup="props.popup">
		<template #header-extra>
			<NTooltip trigger="hover" placement="bottom">
				<template #trigger>
					<NButton size="small" text @click="openPopup">
						<IconExternalLink />
					</NButton>
				</template>

				{{ t('sharedButtons.popout') }}
			</NTooltip>

			<NPopselect
				v-model:value="enabledEvents" multiple :options="enabledEventsOptions"
				trigger="click"
			>
				<NButton text>
					<IconSettings />
				</NButton>
			</NPopselect>
		</template>
		<NScrollbar v-if="events.length" trigger="none">
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
		</NScrollbar>
		<div v-else class="flex items-center justify-center h-full">
			<NText class="text-4xl">
				No events
			</NText>
		</div>
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
