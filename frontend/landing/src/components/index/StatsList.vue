<template>
	<div class="bg-[#17171A] w-full px-5 py-6">
		<Flicking
			:plugins="plugins"
			:options="{
				panelsPerView,
				align: 'next',
				bound: true,
			}"
			class="flex w-full max-w-5xl mx-auto cursor-grab select-none"
		>
			<div v-for="item in stats" :key="item.key" class="inline-flex flex-col items-center justify-center w-full">
				<span class="font-semibold lg:text-6xl text-[min(40px,11vw)] text-white leading-[1.2] tracking-tight">
					{{ item.value }}
				</span>
				<span class="text-[#ADB0B8] lg:text-lg lg:mt-2 leading-normal whitespace-nowrap">
					{{ item.key }}
				</span>
			</div>
		</Flicking>
	</div>
</template>

<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';
import { type Response } from '@twir/grpc/generated/api/api/stats';
import { useWindowSize } from '@vueuse/core';
import { computed } from 'vue';

import { unProtectedClient } from '../../api/twirp.js';


let res: Response | undefined;

try {
  const { response } = await unProtectedClient.getStats({}, { timeout: 5000 });
  res = response;
} catch { /* ignore error */ }

const formatter = Intl.NumberFormat('en-US', {
  notation: 'compact',
  maximumFractionDigits: 1,
});

const stats = [
  {
    key: 'Users',
    value: formatter.format(res?.users ?? 0),
  },
  {
    key: 'Channels',
    value: formatter.format(res?.channels ?? 0),
  },
  {
    key: 'Commands',
    value:formatter.format(res?.commands ?? 0),
  },
  {
    key: 'Messages',
    value: formatter.format(res?.messages ?? 0),
  },
];

const { width: windowWidth } = useWindowSize();

const panelsPerView = computed(() => {
	if (windowWidth.value === Infinity) {
		return 4;
	} else if (windowWidth.value < 410) {
    return 1;
  } else if (windowWidth.value < 568) {
    return 2;
  } else if (windowWidth.value < 768) {
    return 3;
  } else {
    return 4;
  }
});

const plugins = [new AutoPlay()];
</script>

<style>
@import '@egjs/vue3-flicking/dist/flicking.css';
</style>
