<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';
import type { GetTwirStreamersResponse_Streamer } from '@twir/grpc/generated/api/api/stats';
import { ref } from 'vue';

defineProps<{
	streamers: GetTwirStreamersResponse_Streamer[][]
}>();

const wrapper = ref<HTMLElement>(null);
</script>

<template>
	<div class="flicking-viewport">
		<div ref="wrapper" class="wrapper">
			<Flicking
				v-if="wrapper"
				ref="flickingRef"
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
					// 	parentEl: wrapper,
					// }),
					// new Pagination({
					// 	parentEl: wrapper,
					// 	type: 'scroll',
					// }),
				]"
			>
				<div v-for="(item, idx) in streamers" :key="idx" class="slider-review-card">
					<div v-for="(streamer, streamerIdx) of item" :key="streamerIdx" class="flex gap-3 items-center">
						<img :src="streamer.avatar" class="rounded-full w-10 h-10" draggable="false" />
						<a class="flex flex-col gap-1" :href="`https://twitch.tv/${streamer.userLogin}`" target="_blank">
							<span>{{ streamer.userDisplayName }}</span>
							<span class="text-xs uppercase">{{ streamer.followersCount }} followers</span>
						</a>
					</div>
				</div>
			</Flicking>
		</div>
	</div>
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

:deep(.flicking-pagination-bullet) {
	background-color: #fff !important;
}
</style>
