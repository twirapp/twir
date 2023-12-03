<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';
import type { Profile } from '@twir/grpc/generated/api/api/auth';
import type { GetTwirStreamersResponse_Streamer } from '@twir/grpc/generated/api/api/stats';
import { onMounted, ref } from 'vue';

import ReviewDialog from './ReviewDialog.vue';
import { browserUnProtectedClient } from '../../api/twirp-browser';
import { chunk } from '../../utils/chunk.js';

const flickingWrapper = ref<HTMLElement>(null);

defineProps<{
  profile?: Profile
}>();

const streamersWithFollowers = ref<GetTwirStreamersResponse_Streamer[][]>([]);

onMounted(async () => {
	try {
		const twitchStreamersReq = await browserUnProtectedClient.getStatsTwirStreamers({});
		const sortedStreamers = twitchStreamersReq.response.streamers.sort((a, b) => a.followersCount - b.followersCount);
		// streamersWithFollowers.value = chunk(twitchStreamersReq.response.streamers, 3);
		streamersWithFollowers.value = chunk(Array.from({ length: 500 }).map(() => sortedStreamers.at(0)!), 3);
	} catch (e) {
		console.error('cannot get twitch streamers', e);
	}
});

</script>

<template>
	<section id="reviews">
		<div class="reviews-bg">
			<div class="container py-24">
				<div class="flex justify-between">
					<h3 class="text-white text-4xl font-bold">
						Streamers use twir
					</h3>

					<ReviewDialog :profile="profile" />
				</div>
			</div>
			<div ref="flickingWrapper" class="flicking-viewport">
				<div class="wrapper">
					<Flicking
						v-if="streamersWithFollowers?.length"
						:options="{
							renderOnlyVisible: true,
							moveType: 'snap',
							circularFallback: 'linear',
							circular: false,
							align: 'prev'
						}"
						:plugins="[
							new AutoPlay({ stopOnHover: true, duration: 2500 }),
							// new Arrow({
							// 	parentEl: flickingWrapper,
							// }),
							// new Pagination({
							// 	type: 'bullet',
							// 	parentEl: flickingWrapper,
							// }),
						]"
					>
						<div v-for="(item, idx) in streamersWithFollowers" :key="idx" class="slider-review-card">
							<div v-for="(streamer, streamerIdx) of item" :key="streamerIdx" class="flex gap-3 items-center">
								<img :src="streamer.avatar" class="rounded-full w-10 h-10" draggable="false" />
								<a class="flex flex-col gap-1" :href="`https://twitch.tv/${streamer.userLogin}`" target="_blank">
									<span>{{ streamer.userDisplayName }}</span>
									<span class="text-xs uppercase">{{ streamer.followersCount }} followers</span>
								</a>
							</div>
						</div>
					</Flicking>
					<!-- <span class="flicking-arrow-prev is-outside"></span>
					<span class="flicking-arrow-next is-outside"></span>
					<div class="flicking-pagination"></div> -->
				</div>
			</div>
		</div>
	</section>
</template>

<style>
@import "@egjs/vue3-flicking/dist/flicking.css";
@import "@egjs/flicking-plugins/dist/flicking-plugins.css";
@import "@egjs/flicking-plugins/dist/arrow.css";
@import "@egjs/flicking-plugins/dist/pagination.css";

.wrapper {
	display: flex;
	width: 50dvw;
	align-items: center;
	justify-content: center;
	padding-bottom: 40px;
	margin: auto;
	height: 300px;
}

.slider-review-card {
  width: auto;
  margin: 0 12px;
  height: auto;
  opacity: 1 !important;

  @media screen and (max-width: 565.98px) {
    width: calc(100vw - 24px * 2);
  }

	-webkit-user-select: none; /* Safari */
  -ms-user-select: none; /* IE 10 and IE 11 */
  user-select: none;

	@apply gap-5 p-6 rounded-[12px] bg-[#24242780] border-[#393A3E] inline-flex flex-col border select-none
}
</style>
