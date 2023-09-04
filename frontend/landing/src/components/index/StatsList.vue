<script lang="ts" setup>
import { onMounted, onUnmounted, ref } from 'vue';


const stats = ref({
	Users: 0,
	Channels: 0,
	Commands: 0,
	Messages: 0,
	'Used Emotes': 0,
});

async function fetchStats() {
	const { browserUnProtectedClient } = await import('../../api/twirp-browser.js');

	const req = await browserUnProtectedClient.getStats({});
	const res = await req.response;
	const formatter = Intl.NumberFormat('en-US', {
		notation: 'compact',
		maximumFractionDigits: 1,
	});

	stats.value.Users = formatter.format(res.users);
	stats.value.Channels = formatter.format(res.channels);
	stats.value.Commands = formatter.format(res.commands);
	stats.value.Messages = formatter.format(res.messages);
	stats.value['Used Emotes'] = formatter.format(res.usedEmotes);
}

let interval;

onMounted(async () => {
	if (typeof window === 'undefined') return;
	await fetchStats();
	interval = setInterval(fetchStats, 5 * 1000);
});

onUnmounted(() => {
	clearInterval(interval);
});
</script>

<template>
	<div class="bg-[#17171A] px-5 py-6 gap-32 flex flex-wrap justify-center">
		<div
			v-for="key of Object.keys(stats)"
			:key="key"
			class="inline-flex flex-col items-center justify-center"
		>
			<span
				class="font-semibold lg:text-6xl text-[min(40px,11vw)] text-white leading-[1.2] tracking-tight"
			>
				{{ stats[key] }}
			</span>
			<span class="text-[#ADB0B8] lg:text-lg lg:mt-2 leading-normal whitespace-nowrap">
				{{ key }}
			</span>
		</div>
	</div>
</template>

