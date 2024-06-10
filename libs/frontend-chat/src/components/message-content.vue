<script setup lang="ts">
import { EmoteFlag, type MessageChunk } from '../types.js'

const props = defineProps<{
	chunks: MessageChunk[]
	isItalic?: boolean
	textShadowColor?: string
	textShadowSize?: number
	userColor: string
}>()

function getEmoteWidth(isGrowX: boolean, width?: number) {
	if (isGrowX) {
		return `${width ? `${width * 2}px` : 'auto'}`
	}

	if (width) {
		return `${width}px`
	}
}

const mappedChunks = props.chunks.reduce((acc, chunk) => {
	if (chunk.type === 'text') {
		const lastChunk = acc[acc.length - 1]
		if (lastChunk && lastChunk.type === 'text') {
			lastChunk.value += ` ${chunk.value}`
		} else {
			acc.push(chunk)
		}
	} else {
		acc.push(chunk)
	}

	return acc
}, [] as MessageChunk[])
</script>

<template>
	<span
		class="text"
		:style="{
			fontStyle: isItalic ? 'italic' : 'normal',
			color: isItalic ? userColor : 'inherit',
		}"
	>
		<template v-for="(chunk, _) of mappedChunks" :key="_">
			<div
				v-if="['emote', '3rd_party_emote'].includes(chunk.type)"
				class="emote"
			>
				<img
					:src="chunk.value"
					:class="{
						'emote-cursed': chunk.flags?.includes(EmoteFlag.Cursed),
						'flipX': chunk.flags?.includes(EmoteFlag.FlipX),
						'flipY': chunk.flags?.includes(EmoteFlag.FlipY),
					}"
					:style="{
						width: getEmoteWidth(chunk.flags?.includes(EmoteFlag.GrowX) ?? false, chunk.emoteWidth),
					}"
				/>

				<img
					v-for="(c, idx) of chunk.zeroWidthModifiers" :key="idx" class="emote-zerowidth"
					:src="c"
				/>
			</div>

			<template v-else-if="['text', 'emoji'].includes(chunk.type)">
				{{ chunk.value }}
			</template>
		</template>
	</span>
</template>

<style scoped>
.text {
	vertical-align: top;
}

.text > span {
	font-style: inherit;
	font-weight: inherit;
	font-style: inherit;
}

.emote img {
	height: 1.5em;
}

.text .emote {
	position: relative;
	display: inline-block;
	margin-left: 4px;
	margin-right: 4px;
	vertical-align: middle;
}

.text .emote .emote-zerowidth {
	top: 50%;
	left: 50%;
	bottom: 0;
	transform: translate(-50%, -50%);
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
