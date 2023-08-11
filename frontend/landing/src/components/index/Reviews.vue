<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';
import { type TwitchUser } from '@twir/grpc/generated/api/api/twitch';
import { onMounted, ref, computed } from 'vue';

const plugins = [new AutoPlay({ stopOnHover: true })];

const reviews = [
{
    id: '139336353',
    username: '7ssk7',
		roles: ['Streamer', 'pro player'],
    comment: `I've been using Twir for a few years now. There are useful integrations with Volaroant, Spotify. I am pleased with its stability and functionality.`,
    avatarUrl:
      'https://static-cdn.jtvnw.net/jtv_user_pictures/66cb7060-1a8a-4fca-9ccd-f760b70af636-profile_image-70x70.png',
    rating: 5,
  },
  {
    id: '104435562',
    username: 'qrushcsgo',
		roles: ['Streamer', 'player'],
    comment: `Good, handy bot for streaming. Easy to set all the settings and everything works clearly. Recommended 👍`,
    avatarUrl:
      'https://static-cdn.jtvnw.net/jtv_user_pictures/a477bccc-9b23-44d7-a379-fe64f67898c3-profile_image-70x70.png',
    rating: 4,
  },
	{
		id: '48385787',
		username: 'promotive',
		roles: ['Streamer'],
		comment: 'Been using the bot since 2017, now evolved to a huge and feature-rich bot, with all the right features and needs for any streamer. It has everything you need to easily automate your chat, the functionality is very advanced.',
		avatarUrl: 'https://static-cdn.jtvnw.net/jtv_user_pictures/6808c622-2cf0-4319-a2ac-d91ae5212928-profile_image-70x70.png',
		rating: 5,
	},
	{
		id: '155644238',
		username: 'le_xot',
		roles: ['Streamer'],
		comment: 'Twir combines a simple and clear interface, extensive customization and integration options, making this bot an indispensable assistant on my broadcasts.',
		avatarUrl: 'https://static-cdn.jtvnw.net/jtv_user_pictures/423e40e6-9534-46ac-9ed8-5714657dd03b-profile_image-70x70.png',
		rating: 5,
	},
];

const twitchUsers = ref<TwitchUser[]>([]);

onMounted(async () => {
	if (typeof window === 'undefined') return;
  const { browserUnProtectedClient } = await import('../../api/twirp-browser.js');

	const request = await browserUnProtectedClient.twitchGetUsers({
		ids: reviews.map(r => r.id),
		names: [],
	});
	twitchUsers.value = request.response.users;
});

const mappedReviews = computed(() => {
	if (!twitchUsers.value.length) return reviews;

	return reviews.map(r => {
		const twitchUser = twitchUsers.value.find(u => u.id === r.id);

		return {
			...r,
			avatarUrl: twitchUser.profileImageUrl ?? r.avatarUrl,
			username: twitchUser.login ?? r.username,
		};
	});
});
</script>

<template>
	<div class="container py-24">
		<div class="flex justify-between">
			<h1 class="text-white text-4xl font-bold">
				Reviews from streamers and viewers
			</h1>

			<button class="inline-flex px-6 py-3 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4">
				Leave feedback
			</button>
		</div>
	</div>

	<div class="wrapper">
		<Flicking
			:options="{
				align: 'center',
				renderOnlyVisible: true,
				moveType: 'snap',
				circular: true,
				circularFallback: 'bound',
			}"
			:plugins="plugins"
			class="cursor-grab mb-20"
		>
			<div v-for="review of mappedReviews" :key="review.id" class="panel gap-5 p-6 rounded-[12px] bg-[#24242780] border-[#393A3E] inline-flex flex-col border select-none slider-review-card">
				<div class="flex justify-between">
					<div class="flex items-center gap-4">
						<img :src="review.avatarUrl" class="w-11 h-11 rounded-full" />
						<div class="flex flex-col">
							<a
								class="text-lg text-white font-semibold"
								:href="'https://twitch.tv/' + review.username"
								target="_blank"
							>
								@{{ review.username }}
							</a>
							<span class="text-sm text-stone-400">{{ review.roles.join(', ') }}</span>
						</div>
					</div>

				<!-- <div class="flex gap-1">
					<img v-for="(_, index) of Array(review.rating)" :key="index" :src="RatingStarFilled" />
					<img v-for="(_, index) of Array(5 - review.rating)" :key="index" :src="RatingStarFilled" />
				</div> -->
				</div>
				<span class="text-base text-stone-400 line-clamp-4 text-ellipsis" :alt="review.comment">
					{{ review.comment }}
				</span>
			</div>
		</Flicking>
	</div>
</template>

<style>
@import '@egjs/vue3-flicking/dist/flicking.css';

.slider-review-card {
  width: 380px;
  margin: 0 12px;
	height: 230px;
	opacity: 1 !important;

  @media screen and (max-width: 565.98px) {
    width: calc(100vw - 24px * 2);
  }
}
</style>
