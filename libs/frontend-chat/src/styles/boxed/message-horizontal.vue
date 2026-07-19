<script lang="ts" setup>
import MessageBadges from '../../components/message-badges.vue'
import MessageContent from '../../components/message-content.vue'
import { normalizeDisplayName } from '../../helpers.js'

import type { MessageComponentProps } from '../../types.js'

defineProps<MessageComponentProps>()
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
			<MessageBadges :msg="msg" :settings="settings" />
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
	white-space: nowrap;
	justify-content: center;
	--message-badges-margin: 0;
}

.message .profile {
	display: flex;
	justify-content: space-between;
	gap: 4px;
}

.username {
	color: v-bind(userColor);
	font-weight: 700;
}
</style>
