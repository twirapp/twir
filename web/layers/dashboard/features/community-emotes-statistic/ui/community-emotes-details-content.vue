<script setup lang="ts">
import { RadioGroupItem, RadioGroupRoot } from "reka-ui";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import { VisArea, VisAxis, VisLine, VisXYContainer } from "@unovis/vue";

import type { ChartConfig } from "@/components/ui/chart";

import {
	ChartContainer,
	ChartCrosshair,
	ChartTooltip,
	ChartTooltipContent,
	componentToString,
} from "@/components/ui/chart";
import {
	useCommunityEmotesDetails,
	useCommunityEmotesDetailsName,
} from "@/features/community-emotes-statistic/composables/use-community-emotes-details.js";
import { useTranslatedRanges } from "@/features/community-emotes-statistic/composables/use-translated-ranges.js";
import CommunityEmotesDetailsContentUsersHistory from "@/features/community-emotes-statistic/ui/community-emotes-details-content-users-history.vue";
import CommunityEmotesDetailsContentUsersTop from "@/features/community-emotes-statistic/ui/community-emotes-details-content-users-top.vue";

const { t } = useI18n();
const { ranges } = useTranslatedRanges();
const { details, range } = useCommunityEmotesDetails();
const { emoteName } = useCommunityEmotesDetailsName();

interface Data {
	timestamp: number;
	date: Date;
	count: number;
};

const chartData = computed<Data[]>(() => {
	if (!details.value?.emotesStatisticEmoteDetailedInformation?.graphicUsages) {
		return [];
	}

	return details.value.emotesStatisticEmoteDetailedInformation.graphicUsages.map(
		({ timestamp, count }) => ({
			timestamp,
			date: new Date(timestamp),
			count,
		}),
	);
});

const chartConfig = {
	count: {
		label: "Usage",
		theme: {
			light: "#10b981",
			dark: "#10b981",
		},
	},
} satisfies ChartConfig;

const svgDefs = `
  <linearGradient id="fillCount" x1="0" y1="0" x2="0" y2="1">
    <stop
      offset="5%"
      stop-color="var(--color-count)"
      stop-opacity="0.8"
    />
    <stop
      offset="95%"
      stop-color="var(--color-count)"
      stop-opacity="0.1"
    />
  </linearGradient>
`;

const tableTabs = [
	{ key: "top", text: t("community.emotesStatistic.details.usersTabs.top") },
	{ key: "history", text: t("community.emotesStatistic.details.usersTabs.history") },
];

const tableTab = ref<"top" | "history">("top");
</script>

<template>
	<div class="flex flex-col divide-y divide-white/10">
		<h1 class="text-4xl font-medium px-6 py-6">
			{{ emoteName }}
		</h1>
		<div class="flex flex-col gap-6 px-6 py-7">
			<div class="flex justify-between flex-wrap">
				<h1 class="text-2xl font-medium">
					{{ t("community.emotesStatistic.details.stats") }}
				</h1>
				<RadioGroupRoot
					v-model="range"
					class="inline-flex w-full rounded-[7px] bg-zinc-800 p-px md:w-auto"
				>
					<RadioGroupItem
						v-for="[key, text] of Object.entries(ranges)"
						:key="key"
						:value="key"
						class="h-8 flex-1 rounded-md px-3 text-[13px] text-white/75 transition-colors hover:bg-white/5 data-[active=true]:bg-white/20 data-[active=true]:text-white data-[active=true]:shadow-md md:flex-auto whitespace-nowrap"
					>
						{{ text }}
					</RadioGroupItem>
				</RadioGroupRoot>
			</div>

			<div class="relative h-[240px]">
				<ChartContainer
					v-if="chartData.length > 0"
					:config="chartConfig"
					class="w-full h-full"
					cursor
				>
					<VisXYContainer
						:data="chartData"
						:svg-defs="svgDefs"
						:margin="{ left: 0, right: 0, top: 10, bottom: 30 }"
						:y-domain="[0, undefined]"
					>
						<VisArea
							:x="(d: Data) => d.date"
							:y="(d: Data) => d.count"
							color="url(#fillCount)"
							:opacity="0.4"
						/>
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
							:num-ticks="6"
							:tick-format="
								(d: number) => {
									const date = new Date(d);
									return date.toLocaleDateString('en-US', {
										month: 'short',
										day: 'numeric',
									});
								}
							"
						/>
						<VisAxis type="y" :num-ticks="4" :tick-line="false" :domain-line="false" />
						<ChartTooltip />
						<ChartCrosshair
							:template="
								componentToString(chartConfig, ChartTooltipContent, {
									labelFormatter(d) {
										return new Date(d).toLocaleDateString('en-US', {
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
		</div>
		<div class="flex flex-col gap-6 px-6 py-7">
			<div class="flex justify-between flex-wrap">
				<h1 class="text-2xl font-medium">
					{{ t("community.emotesStatistic.details.users") }}
				</h1>
				<RadioGroupRoot
					v-model="tableTab"
					class="inline-flex w-full rounded-[7px] bg-zinc-800 p-px md:w-auto"
				>
					<RadioGroupItem
						v-for="tab of tableTabs"
						:key="tab.key"
						:value="tab.key"
						class="h-8 flex-1 rounded-md px-3 text-[13px] text-white/75 transition-colors hover:bg-white/5 data-[active=true]:bg-white/20 data-[active=true]:text-white data-[active=true]:shadow-md md:flex-auto whitespace-nowrap"
					>
						{{ tab.text }}
					</RadioGroupItem>
				</RadioGroupRoot>
			</div>
			<CommunityEmotesDetailsContentUsersTop v-if="tableTab === 'top'" />
			<CommunityEmotesDetailsContentUsersHistory v-else />
		</div>
	</div>
</template>
