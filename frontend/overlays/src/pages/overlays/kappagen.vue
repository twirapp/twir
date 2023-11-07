<script setup lang="ts">
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';
import '../../components/kappagen/index.js';

import { kappagenAnimations } from '../../components/kappagen/animations.js';
import type { Emote } from '../../components/kappagen/index.js';
import Kappagen from '../../components/kappagen/kappagen.vue';

const route = useRoute();
const apiKey = route.params.apiKey as string;

const emote: Emote = {
	url: 'https://cdn.7tv.app/emote/6548b7074789656a7be787e1/4x.webp',
	zwe: [
		{
			url: 'https://cdn.7tv.app/emote/6128ed55a50c52b1429e09dc/4x.webp',
		},
	],
};

const emote2: Emote = {
	url: 'https://cdn.betterttv.net/emote/65425cc67080a9fc246a30d6/3x.webp',
};


// на ивент\команду
const kappa = () => {
	window.kappagen.show(
		[emote, emote2],
		// kappagenAnimations[window.random(kappagenAnimations.length)],
		kappagenAnimations.find(a => a.style === 'Fountain')!,
	);
};

// смайлики в чате
const spawn = () => {
	const randomCountEmotes = window.random(50) + 1;
	const emotes = new Array(randomCountEmotes).fill([emote, emote2]).flat();
	window.emote.addToShowList(emotes);
	window.emote.showEmotes();
};

onMounted(() => {
	window.startup();
});
</script>

<template>
	<button @click="kappa">
		kappa
	</button>
	<button @click="spawn">
		spawn
	</button>
	<Kappagen />
</template>
