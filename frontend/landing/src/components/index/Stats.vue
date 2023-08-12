<script lang="ts" setup>
import { computed, onMounted } from 'vue';

import { unProtectedClient } from '../../api/twirp.js';

const { response } = await unProtectedClient.getStats({});

const stats = computed(() => {
	return {
		'Users': response.users ?? 0,
		'Channels': response.channels ?? 0,
		'Commands': response.commands ?? 0,
		'Messages': response.messages ?? 0,
	};
});

onMounted(() => {
	if (typeof window === 'undefined') return;
	const numberFormatter = new Intl.NumberFormat(navigator.language, { notation: 'compact' });
  const spans = document.querySelectorAll<HTMLElement>('[data-stats-count]');
  spans.forEach((s) => {
    s.innerText = numberFormatter.format(Number(s.dataset.statsCount ?? 0));
  });
});
</script>

<template>
	<div class="bg-[#17171A] px-[32px] py-[20px]">
		<div class="container mx-auto">
			<div class="flex justify-center flex-wrap gap-[32px]">
				<div v-for="([key, value]) of Object.entries(stats)" :key="key" class="flex flex-col items-center gap-[12px] min-w-[200px]">
					<span class="font-semibold text-6xl text-white stats-count" data-stats-count="{item.count}">{{ value }}</span>
					<span class="font-normal text-lg text-[#ADB0B8]">{{ key }}</span>
				</div>
			</div>
		</div>
	</div>
</template>
