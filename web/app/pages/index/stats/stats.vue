<script setup lang="ts">
const formatter = Intl.NumberFormat('en-US', {
	notation: 'compact',
	maximumFractionDigits: 1,
})

function formatNumber(value?: number | bigint) {
	return formatter.format(value ?? 0)
}

const stats = [
	{
		key: 'Channels',
		value: formatNumber(1),
	},
	{
		key: 'Created commands',
		value: formatNumber(1),
	},
	{
		key: 'Viewers',
		value: formatNumber(1),
	},
	{
		key: 'Messages',
		value: formatNumber(1),
	},
	{
		key: 'Used emotes',
		value: formatNumber(1),
	},
	{
		key: 'Used commands',
		value: formatNumber(1),
	},
]
</script>

<template>
	<section id="stats" class="bg-[#17171A] px-5 py-6 gap-32 flex flex-wrap justify-center">
		<div v-for="stat of stats" :key="stat.key" class="inline-flex flex-col items-center justify-center">
			<span class="font-semibold lg:text-6xl text-[min(40px,11vw)] text-white leading-[1.2] tracking-tight">
				{{ stat.value }}
			</span>
			<span class="text-[#ADB0B8] lg:text-lg lg:mt-2 leading-normal whitespace-nowrap">
				{{ stat.key }}
			</span>
		</div>
	</section>
</template>

<style scoped>

</style>
