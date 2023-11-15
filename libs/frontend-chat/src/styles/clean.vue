<script lang="ts" setup>
import MessageContent from '../components/messageContent.vue';
import { normalizeDisplayName } from '../helpers.js';
import type { Settings, Message } from '../types.js';


defineProps<{
	msg: Message,
	settings: Settings
}>();
</script>

<template>
	<div class="message">
		<div v-if="settings.showBadges && msg.badges" class="badges">
			<template
				v-for="(badgeValue, badgeName) of msg.badges"
				:key="badgeName+badgeValue"
			>
				<img
					v-if="settings.channelBadges.get(`${badgeName}-${badgeValue}`)"
					:src="settings.channelBadges.get(`${badgeName}-${badgeValue}`)!.image_url_4x"
					class="badge"
				/>

				<img
					v-else-if="settings.globalBadges.get(badgeName)?.versions.length"
					:src="settings.globalBadges.get(badgeName)!.versions.at(-1)!.image_url_4x"
					class="badge"
				/>
			</template>
		</div>
		<div v-if="msg.sender" class="profile" :style="{ color: msg.senderColor }">
			{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}{{ msg.isItalic ? '' : ':' }}
		</div>
		<message-content
			:chunks="msg.chunks"
			:is-italic="msg.isItalic"
			:text-shadow-color="settings.textShadowColor"
			:text-shadow-size="settings.textShadowSize"
		/>
	</div>
</template>

<style scoped>
.message {
	width: 100%;
	line-height: 1em;
	margin-left: 0.2em;
}

.message .badges {
	display: inline-flex;
	gap: 4px;
	align-self: center;
	max-height: 0.8em;
	margin-right: 4px;
	transform: translateY(0.2em);
}

.message .badges .badge {
	height: 1em;
	width: 1em;
}

.message .profile {
	display: inline-flex;
}

.message .text {
	margin-left: 4px;
}

.message .text .emote {
	height: 1em;
	width: 1em;
}
</style>
