<script setup lang="ts">
import KappagenOverlay, { type KappagenEmoteConfig, type Emote } from 'kappagen';
import { computed, onMounted, reactive, ref, toRef } from 'vue';
import { useRoute } from 'vue-router';
import 'kappagen/style.css';

import { useThirdPartyEmotes } from '../../components/chat_tmi_emotes.js';
import { animations } from '../../components/kappagen_animations.js';

const kappagen = ref<InstanceType<typeof KappagenOverlay>>();
const route = useRoute();
const apiKey = route.params.apiKey as string;
console.log(apiKey);

const { emotes } = useThirdPartyEmotes(toRef('fukushine'), toRef('971211575'));

const kappagenEmotes = computed(() => {
	const emotesArray = Object.values(emotes.value);

	return emotesArray.filter(e => !e.isZeroWidth && !e.isModifier);
});

const emoteConfig = reactive<KappagenEmoteConfig>({
  max: 10,
  time: 10,
  queue: 100,
  cube: {
    speed: 10,
  },
});

const emote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
};

// на ивент\команду
const kappa = async () => {
	kappagen.value?.emote.addEmotes([emote]);
	kappagen.value?.emote.showEmotes();
};

const enabledAnimations = animations.filter(a => a.style !== 'Text');

// смайлики в чате
const spawn = async () => {
	const emotesValue = Object.values(kappagenEmotes.value);

	const randomEmotes = Array(50)
		.fill(null)
		.map(() => emotesValue[Math.floor(Math.random() * emotesValue.length)].urls.at(-1));
	const randomAnimation = enabledAnimations[Math.floor(Math.random() * enabledAnimations.length)];

	await kappagen.value!.kappagen.run(
		randomEmotes.map(url => ({ url: url! })),
		randomAnimation,
	);
};

onMounted(() => {
	kappagen.value?.init();
});
</script>

<template>
	<button @click="kappa">
		kappa
	</button>
	<button @click="spawn">
		spawn
	</button>
	<kappagen-overlay ref="kappagen" :emote-config="emoteConfig" />
</template>
