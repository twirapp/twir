<script setup lang="ts">
import MessageBadges from '../../components/message-badges.vue'
import MessageContent from '../../components/message-content.vue'
import { normalizeDisplayName } from '../../helpers'

import type { MessageComponentProps } from '../../types'

defineProps<MessageComponentProps>()
</script>

<template>
	<div class="message">
		<div class="profile">
			<MessageBadges :msg="msg" :settings="settings" />
			<div v-if="msg.sender" class="username">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}
			</div>
			<span v-if="msg.sender">
				{{ msg.isItalic ? '' : ':' }}
			</span>
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
	display: inline-flex;
	align-items: center;
	white-space: nowrap;
}

.profile {
	display: flex;
	align-items: center;
	margin-right: 4px;
}

.username {
	color: v-bind(userColor);
	font-weight: 700;
}
</style>
