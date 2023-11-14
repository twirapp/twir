<script setup lang="ts">
import { type MessageChunk, EmoteFlag } from '../types.js';

defineProps<{
	chunks: MessageChunk[]
	isItalic?: boolean
}>();
</script>

<template>
	<span class="text" :style="{ fontStyle: isItalic ? 'italic' : 'normal' }">
		<template v-for="(chunk, _) of chunks" :key="_">
			<div
				v-if="['emote', '3rd_party_emote'].includes(chunk.type)"
				class="emote"
			>
				<img
					:src="chunk.type === 'emote'
						? `https://static-cdn.jtvnw.net/emoticons/v2/${chunk.value}/default/dark/1.0`
						: chunk.value
					"
					:class="{
						'emote-cursed': chunk.flags?.includes(EmoteFlag.Cursed),
						'flipX': chunk.flags?.includes(EmoteFlag.FlipX),
						'flipY': chunk.flags?.includes(EmoteFlag.FlipY),
					}"
					:style="{
						width: chunk.emoteWidth
							? `${chunk.flags?.includes(EmoteFlag.GrowX) ? chunk.emoteWidth * 2 : chunk.emoteWidth}px`
							: undefined,
						height: chunk.emoteHeight
							? `${chunk.emoteHeight}px`
							: undefined
					}"
				/>

				<span v-for="(c, idx) of chunk.zeroWidthModifiers" :key="idx" class="emote-zerowidth">
					<img
						:src="c"
						:style="{
							maxWidth: chunk.emoteWidth
								? `${chunk.flags?.includes(EmoteFlag.GrowX) ? chunk.emoteWidth * 2 : chunk.emoteWidth}px`
								: `32px`,
							maxHeight: chunk.emoteHeight
								? `${chunk.emoteHeight}px`
								: `32px`
						}"
					/>
				</span>
			</div>

			<template v-else-if="chunk.type === 'text'">
				{{ chunk.value }}
			</template>
			{{ ' ' }}
		</template>
	</span>
</template>

<style scoped>
.text {
	overflow-wrap: break-word;
}

.text .emote {
	max-height: 1em;
	position: relative;
	display: inline-block;
}

.text .emote .emote-zerowidth {
	top: 50%;
	left: 50%;
	bottom: 0;
	transform: translate(-50%,-50%);
	position: absolute;
}

.text .emote .emote-cursed {
	filter: grayscale(1) brightness(0.7) contrast(2.5);
}

.text .emote .flipX {
	transform: scaleX(-1);
}

.text .emote .flipY {
	transform: scaleY(-1);
}
</style>
