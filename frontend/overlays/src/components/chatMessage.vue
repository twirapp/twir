<script lang="ts" setup>
import { normalizeDisplayName } from './chat_tmi_helpers.js';
import type { Settings } from '../sockets/chat.js';
import type { Message } from '../sockets/chat_tmi.js';

defineProps<{
	msg: Message,
	settings: Settings
}>();
</script>

<template>
	<div class="message">
		<div v-if="msg.badges" class="badges">
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
		<span class="text" :style="{ fontStyle: msg.isItalic ? 'italic' : 'normal' }">
			<template v-for="(chunk, _) of msg.chunks" :key="_">
				<img
					v-if="chunk.type === 'emote'"
					:src="`https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/3.0`"
					class="emote"
				/>

				<img
					v-else-if="chunk.type === '3rd_party_emote'"
					:src="chunk.value"
					class="emote"
				/>

				<template v-else-if="chunk.type === 'text'">
					{{ chunk.value }}
				</template>
				{{ ' ' }}
			</template>
		</span>
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
