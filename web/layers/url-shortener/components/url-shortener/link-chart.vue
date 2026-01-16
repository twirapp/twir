<script setup lang="ts">
import { computed } from "vue";
import {
	ChartContainer,
	ChartCrosshair,
	ChartTooltip,
	ChartTooltipContent,
	componentToString,
} from "@/components/ui/chart";
import { VisAxis, VisLine, VisXYContainer } from "@unovis/vue";
import type { ChartConfig } from "@/components/ui/chart";

const props = defineProps<{
	isDayRange: boolean;
	usages: {
		timestamp: number;
		count: number;
	}[];
}>();

interface Data {
	timestamp: number;
	date: Date;
	count: number;
}

const chartData = computed<Data[]>(() => {
	return props.usages.map(({ timestamp, count }) => ({
		timestamp,
		date: new Date(timestamp),
		count,
	}));
});

const chartConfig = {
	count: {
		label: "Views",
		theme: {
			light: "#10b981",
			dark: "#10b981",
		},
	},
} satisfies ChartConfig;
</script>

<template>
	<div class="w-full h-[200px] overflow-hidden">
		<ChartContainer :config="chartConfig" class="w-full h-full smooth-chart" cursor>
			<VisXYContainer
				:data="chartData"
				:margin="{ left: 0, right: 0, top: 5, bottom: 20 }"
				:y-domain="[0, undefined]"
				:duration="300"
				class="w-full"
			>
				<VisLine
					:x="(d: Data) => d.date"
					:y="(d: Data) => d.count"
					color="var(--color-count)"
					:line-width="2"
				/>
				<VisAxis
					type="x"
					:x="(d: Data) => d.date"
					:tick-line="false"
					:domain-line="false"
					:grid-line="false"
					:num-ticks="props.isDayRange ? 4 : 3"
					:tick-format="
						(d: number) => {
							const date = new Date(d);
							if (props.isDayRange) {
								return date.toLocaleTimeString('en-US', {
									hour: '2-digit',
									minute: '2-digit',
								});
							}
							return date.toLocaleDateString('en-US', {
								month: 'short',
								day: 'numeric',
							});
						}
					"
				/>
				<VisAxis
					type="y"
					:num-ticks="3"
					:tick-line="false"
					:domain-line="false"
					:tick-format="() => ''"
				/>
				<ChartTooltip />
				<ChartCrosshair
					:template="
						componentToString(chartConfig, ChartTooltipContent, {
							labelFormatter(d) {
								const date = new Date(d);
								if (props.isDayRange) {
									return date.toLocaleString('en-US', {
										month: 'short',
										day: 'numeric',
										hour: '2-digit',
										minute: '2-digit',
									});
								}
								return date.toLocaleDateString('en-US', {
									month: 'short',
									day: 'numeric',
									year: 'numeric',
								});
							},
						})
					"
					color="var(--color-count)"
				/>
			</VisXYContainer>
		</ChartContainer>
	</div>
</template>

<style scoped>
.smooth-chart :deep(path) {
	transition:
		d 300ms ease-in-out,
		opacity 300ms ease-in-out;
}

.smooth-chart :deep(circle) {
	transition:
		cx 300ms ease-in-out,
		cy 300ms ease-in-out,
		opacity 300ms ease-in-out;
}

.smooth-chart :deep(line) {
	transition:
		x1 300ms ease-in-out,
		x2 300ms ease-in-out,
		y1 300ms ease-in-out,
		y2 300ms ease-in-out,
		opacity 300ms ease-in-out;
}
</style>
