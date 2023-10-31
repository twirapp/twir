<script lang="ts" setup>
import { normalizeDisplayName } from '../helpers.js';
import type { Settings, Message } from '../types.js';

defineProps<{
	msg: Message,
	settings: Settings
}>();
</script>

<template>
	<div class="message">
		<div class="profile">
			<div v-if="msg.sender" :style="{ color: msg.senderColor }">
				{{ normalizeDisplayName(msg.sender!, msg.senderDisplayName!) }}
			</div>
			<div v-if="settings.showBadges" class="badges">
				<span v-if="settings.showAnnounceBadge && msg.isAnnounce" class="text-badge">Announce</span>
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
		</div>
		<span class="text" :style="{ fontStyle: msg.isItalic ? 'italic' : 'normal' }">
			<template v-for="(chunk, _) of msg.chunks" :key="_">
				<img
					v-if="chunk.type === 'emote'"
					:src="`https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/1.0`"
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
	padding-top: 4px;
	padding-bottom: 4px;
	padding-right: 8px;
	padding-left: 8px;
	font-size: 11px;
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
}
</style>
