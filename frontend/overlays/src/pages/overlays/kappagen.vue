<script setup lang="ts">
import KappagenOverlay, { type Emote } from 'kappagen';
import { onMounted, ref, toRef } from 'vue';
import { useRoute } from 'vue-router';
import 'kappagen/style.css';

import { useThirdPartyEmotes } from '../../components/chat_tmi_emotes.js';
import { animations } from '../../components/kappagen_animations.js';

const kappagen = ref<InstanceType<typeof KappagenOverlay>>();
const route = useRoute();
const apiKey = route.params.apiKey as string;

const { emotes } = useThirdPartyEmotes(toRef('fukushine'), toRef('971211575'));

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
	kappagen.value?.emote.addToShowList([emote]);
	kappagen.value?.emote.showEmotes();
};

const enabledAnimations = animations.filter(a => a.style !== 'Text');

// смайлики в чате
const spawn = async () => {
	const emotesValue = Object.values(emotes.value);

	const randomEmotes = Array(50)
		.fill(null)
		.map(() => emotesValue[Math.floor(Math.random()*emotesValue.length)].urls.at(-1));

	await kappagen.value?.kappagen.show(
    randomEmotes.map(url => ({ url: url! })),
		enabledAnimations[Math.floor(Math.random()*enabledAnimations.length)],
  );
};

onMounted(() => {
	kappagen.value?.startup();
});
</script>

<template>
	<button @click="kappa">
		kappa
	</button>
	<button @click="spawn">
		spawn
	</button>
	<kappagen-overlay
		ref="kappagen"
	/>
</template>
