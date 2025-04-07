<script lang="ts" setup>
import MessageContent from '../../components/message-content.vue'
import { useMappedBadges } from '../../composables/mapped-badges'
import { normalizeDisplayName } from '../../helpers.js'

import type { MessageComponentProps } from '../../types.js'

defineProps<MessageComponentProps>()

const { globalMappedBadges, channelMappedBadges } = useMappedBadges()
</script>

<template>
	<div
		class="message"
		:style="{
			border: msg.isAnnounce && settings.showAnnounceBadge ? '2px solid #9146ff' : undefined,
		}"
	>
		<div class="profile">
			<div v-if="msg.sender" class="username">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}
			</div>
			<div v-if="settings.showBadges" class="badges">
				<template
					v-for="(badgeValue, badgeName) of msg.badges"
					:key="badgeName + badgeValue"
				>
					<img
						v-if="channelMappedBadges[badgeName]?.versions[badgeValue]"
						:src="channelMappedBadges[badgeName]!.versions[badgeValue].image_url_4x"
						class="badge"
					/>

					<img
						v-if="globalMappedBadges[badgeName]?.versions[badgeValue]"
						:src="globalMappedBadges[badgeName]!.versions[badgeValue].image_url_4x"
						class="badge"
					/>
				</template>
			</div>
		</div>
		<MessageContent
			:chunks="msg.chunks"
			:is-italic="msg.isItalic"
			:text-shadow-color="settings.textShadowColor"
			:text-shadow-size="settings.textShadowSize"
			:user-color="userColor"
		/>
	</div>
</template>

<style scoped>
.message {
	display: flex;
	flex-direction: column;
	padding: 0.5em;
	gap: 0.2em;
	border-radius: 8px;
	background-color: #252525;
}

.message .badges {
	display: inline-flex;
	gap: 4px;
	align-items: center;
}

.message .badges .text-badge {
	padding: 4px;
	font-size: inherit;
	background-color: #6d6767;
	border-radius: 4px;
	text-transform: uppercase;
}

.message .badges .badge {
	height: 1em;
	width: 1em;
}

.message .profile {
	display: flex;
	justify-content: space-between;
	gap: 4px;
	align-items: center;
}

.username {
	color: v-bind(userColor);
	font-weight: 700;
}
</style>
