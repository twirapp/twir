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
		<div v-if="msg.badges?.length || msg.sender" class="profile">
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
			<div v-if="msg.senderDisplayName" :style="{ color: msg.senderColor }">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}{{ msg.isItalic ? '' : ':' }}
			</div>
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
  display: flex;
  gap: 8px;
  justify-content: flex-start;
  align-items: flex-start;
	width: 100%;
}

.message .badges {
  display: flex;
  gap: 4px;
}

.message .badges .badge {
	height: 20px;
	width: 20px;
}

.message .profile {
	display: flex;
  flex-wrap: nowrap;
  gap: 4px;
  align-items: center;
}

.message .text {

}

.message .text .emote {
	height: 20px;
	width: 20px;
}
</style>
