<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';
import { type TwitchUser } from '@twir/grpc/generated/api/api/twitch';
import { onMounted, ref, computed, watch } from 'vue';

const plugins = [new AutoPlay({ stopOnHover: true })];

const reviews = [
	{
		id: '139336353',
		username: '7ssk7',
		roles: ['Streamer', 'pro player'],
		comment: `I've been using Twir for a few years now. There are useful integrations with Valorant, Spotify. I am pleased with its stability and functionality.`,
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
		comment:
			'Started used Twir in 2017, now it\'s a robust, feature-rich bot for streamers — complete automation for advanced chat functionality.',
		avatarUrl:
			'https://static-cdn.jtvnw.net/jtv_user_pictures/6808c622-2cf0-4319-a2ac-d91ae5212928-profile_image-70x70.png',
		rating: 5,
	},
	{
		id: '155644238',
		username: 'le_xot',
		roles: ['Streamer'],
		comment:
			'Twir combines a simple and clear interface, extensive customization and integration options, making this bot an indispensable assistant on my broadcasts.',
		avatarUrl:
			'https://static-cdn.jtvnw.net/jtv_user_pictures/423e40e6-9534-46ac-9ed8-5714657dd03b-profile_image-70x70.png',
		rating: 5,
	},
	{
		id: '189703483',
		username: 'daetojekara',
		roles: ['Streamer', 'viewer'],
		comment:
			'Streamline broadcasts: bot controls, configures & monitors. Integrated 3rd-party tools enable dashboard-based basic setup.',
		avatarUrl:
			'https://static-cdn.jtvnw.net/jtv_user_pictures/b73f81e7-3fe1-415b-a543-4fe164d16e56-profile_image-70x70.png',
		rating: 5,
	},
];

const twitchUsers = ref<TwitchUser[]>([]);

onMounted(async () => {
	if (typeof window === 'undefined') return;
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

const dialog = ref<HTMLDialogElement>(null);

function showDialog() {
	dialog.value.showModal();
	const body = document.body;
	body.style.overflowY = 'hidden';
}

watch(dialog, (v) => {
	if (!v) return;

	const body = document.body;

	const closeListener = () => {
		body.style.overflowY = 'auto';
	};

	v.addEventListener('close', closeListener);

	return () => {
		v.removeEventListener('close', closeListener);
	};
});

const reviewForm = ref({
	author: '',
	message: '',
});

const feedBackError = ref<string | null>(null);

async function sendFeedBack() {
	const { author, message } = reviewForm.value;

	if (!author || !message) {
		feedBackError.value = 'Author and message required';
		return;
	}

	const req = await fetch('/a/feedback', {
		method: 'POST',
		body: JSON.stringify({
			author,
			message,
		}),
	});

	if (!req.ok) {
		const response = await req.json();
		if (response.error) {
			feedBackError.value = response.error;
		}
	} else {
		dialog.value?.close();
		reviewForm.value.author = '';
		reviewForm.value.message = '';
	}
}
</script>

<template>
	<div class="reviews-bg">
		<div class="container py-24">
			<div class="flex justify-between">
				<h3 class="text-white text-4xl font-bold">
					Reviews from streamers and viewers
				</h3>

				<button
					class="inline-flex px-6 py-3 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
					@click="showDialog"
				>
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

	<dialog id="feedback-dialog" ref="dialog">
		<div v-if="feedBackError" class="p-3 bg-red-500">
			{{ feedBackError }}
		</div>
		<form class="flex flex-col gap-2 p-5">
			<fieldset class="flex flex-col">
				<label for="author">Username</label>
				<input
					id="author"
					v-model="reviewForm.author"
					placeholder="Enter your twitch username"
					class="p-2 rounded-md shadow-sm outline-none"
					maxlength="25"
					autocomplete="off"
				/>
			</fieldset>
			<fieldset class="flex flex-col relative">
				<label for="message">Message</label>
				<textarea
					id="message"
					v-model="reviewForm.message"
					placeholder="Write your feedback"
					rows="6"
					class="p-2 rounded-md shadow-sm outline-none"
					maxlength="200"
					autocomplete="off"
				/>
				<span class="absolute px-2 py-1 text-xs text-white bg-[#5D58F5] rounded right-2 bottom-2">{{
					reviewForm.message.length
				}}/200</span>
			</fieldset>
			<div class="flex justify-end gap-2">
				<button
					value="cancel"
					formmethod="dialog"
					class="inline-flex px-3 py-1 font-semibold text-white rounded-lg bg-red-500 transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
				>
					Cancel
				</button>
				<button
					class="inline-flex px-3 py-1 font-semibold text-white rounded-lg bg-[#5D58F5] hover:bg-[#6964FF] transition-[background,box-shadow] text-lg focus-visible:outline-none focus-visible:ring-[#5D58F5]/50 focus-visible:ring-4"
					@click="
						(e) => {
							e.preventDefault();
							sendFeedBack();
						}
					"
				>
					Confirm
				</button>
			</div>
		</form>
	</dialog>
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

.reviews-bg {
	background: radial-gradient(ellipse at top, #270a3b, transparent),
	radial-gradient(ellipse at bottom, #000, transparent);
}

#feedback-dialog::backdrop {
	background-color: rgba(0, 0, 0, 9);
	opacity: 0.7;
}

#feedback-dialog {
	box-shadow: 0 4px 5px rgb(0 0 0 / 90%);
	border-radius: 5px;
	min-width: 500px;
}
</style>
