<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';
import type { Profile } from '@twir/grpc/generated/api/api/auth';
import { type TwitchUser } from '@twir/grpc/generated/api/api/twitch';
import { onMounted, ref, computed } from 'vue';

import ReviewDialog from './ReviewDialog.vue';
import { reviews } from '../../data/home/reviews.js';

const plugins = [new AutoPlay({ stopOnHover: true })];

defineProps<{
	profile?: Profile
}>();

const twitchUsers = ref<TwitchUser[]>([]);

onMounted(async () => {
	const { browserUnProtectedClient } = await import('../../api/twirp-browser.js');

	const request = await browserUnProtectedClient.twitchGetUsers({
		ids: reviews.map((r) => r.id),
		names: [],
	});
	twitchUsers.value = request.response.users;
});

const mappedReviews = computed(() => {
	if (!twitchUsers.value.length) return reviews;

	return reviews.map((r) => {
		const twitchUser = twitchUsers.value.find((u) => u.id === r.id);

		return {
			...r,
			avatarUrl: twitchUser.profileImageUrl ?? r.avatarUrl,
			username: twitchUser.login ?? r.username,
		};
	});
});

</script>

<template>
	<section id="reviews">
		<div class="reviews-bg">
			<div class="container py-24">
				<div class="flex justify-between">
					<h3 class="text-white text-4xl font-bold">
						Reviews from streamers and viewers
					</h3>

					<ReviewDialog :profile="profile" />
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
					<div
						v-for="review of mappedReviews"
						:key="review.id"
						class="panel gap-5 p-6 rounded-[12px] bg-[#24242780] border-[#393A3E] inline-flex flex-col border select-none slider-review-card"
					>
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
		</div>
	</section>
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
