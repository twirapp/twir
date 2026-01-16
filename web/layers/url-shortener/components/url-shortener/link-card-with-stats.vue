<script setup lang="ts">
import type { LinkOutputDto } from "@twir/api/openapi";
import { toast } from "vue-sonner";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import Button from "@/components/ui/button/Button.vue";
import UrlShortenerLinkChart from "./link-chart.vue";
import ViewsHistoryDialog from "./views-history-dialog.vue";
import TopCountriesDialog from "./top-countries-dialog.vue";
import { useMetaExtractor } from "../../composables/use-meta-extractor";
import { useShortLinkViewsSubscription } from "../../composables/use-short-link-views-subscription";

const props = defineProps<{
	link: LinkOutputDto;
}>();

const shortId = computed(() => props.link.id);
const { statistics, range, isDayRange, isLoading, refetch } = useShortLinkStatistics(shortId);

// Subscribe to real-time view updates
const { totalViews: liveViews, lastView } = useShortLinkViewsSubscription(shortId);

// Use live views if available, otherwise fall back to initial link views
const displayViews = computed(() => liveViews.value ?? props.link.views);

// Throttled refetch to update chart smoothly without overwhelming the backend
// Refetch at most once per 10 seconds
let lastRefetchTime = 0;
let pendingRefetch = false;
let refetchTimeout: ReturnType<typeof setTimeout> | null = null;
const REFETCH_THROTTLE = 10000; // 10 seconds

watch(liveViews, (newViews, oldViews) => {
	if (newViews !== null && oldViews !== null && newViews > oldViews) {
		const now = Date.now();
		const timeSinceLastRefetch = now - lastRefetchTime;

		if (timeSinceLastRefetch >= REFETCH_THROTTLE) {
			// Enough time has passed, refetch immediately
			refetch();
			lastRefetchTime = now;
			pendingRefetch = false;
		} else if (!pendingRefetch) {
			// Schedule a refetch for when the throttle period ends
			pendingRefetch = true;
			const delay = REFETCH_THROTTLE - timeSinceLastRefetch;

			if (refetchTimeout) {
				clearTimeout(refetchTimeout);
			}

			refetchTimeout = setTimeout(() => {
				refetch();
				lastRefetchTime = Date.now();
				pendingRefetch = false;
				refetchTimeout = null;
			}, delay);
		}
	}
});

// Cleanup timeout on unmount
onUnmounted(() => {
	if (refetchTimeout) {
		clearTimeout(refetchTimeout);
	}
});

// Views history dialog
const showViewsDialog = ref(false);

// Top countries dialog
const showTopCountriesDialog = ref(false);

const clipboardApi = useClipboard();

function copyShortUrl() {
	clipboardApi.copy(props.link.short_url);
	toast.success("Copied", {
		description: "Shortened URL copied to clipboard",
		duration: 2500,
	});
}

function removeProtocol(url: string) {
	return url.replace(/^https?:\/\//, "");
}

const displayShortUrl = computed(() => removeProtocol(props.link.short_url));
const displayUrl = computed(() => removeProtocol(props.link.url));

function formatViews(views: number) {
	const intl = new Intl.NumberFormat(import.meta.client ? navigator.language : "en-US", {
		notation: "compact",
		maximumFractionDigits: 1,
	});
	return intl.format(views);
}

const { extractMetaFromUrl, loading: metaLoading } = useMetaExtractor();
const hasLoaded = ref(false);
const metaData = ref<any>(null);

async function fetchMetadata() {
	if (metaLoading.value || hasLoaded.value) return;

	if (props.link.url) {
		hasLoaded.value = true;

		try {
			const meta = await extractMetaFromUrl(props.link.url);
			metaData.value = meta;
		} catch (err) {
			console.error("Failed to fetch metadata:", err);
			hasLoaded.value = false;
		}
	}
}

watch(
	() => props.link.url,
	() => {
		fetchMetadata();
	},
	{ immediate: true },
);
</script>

<template>
	<div
		class="flex flex-col bg-[hsl(240,11%,9%)] border border-[hsl(240,11%,18%)] w-full rounded-2xl p-4 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<!-- Header with link info and copy button -->
		<div class="flex justify-between items-start mb-4">
			<div class="flex items-center min-w-0 gap-x-3 flex-1">
				<div
					class="flex-none size-fit p-3 rounded-full font-semibold border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,20%)]"
				>
					<Icon v-if="!hasLoaded || !metaData" name="lucide:link" class="w-4 h-4" />
					<img v-else :src="metaData.favicon" class="w-4 h-4" />
				</div>
				<div class="overflow-hidden min-w-0 flex-1">
					<div class="flex items-center gap-2">
						<a
							:href="link.short_url"
							target="_blank"
							class="font-bold hover:text-[hsl(240,11%,85%)] transition-colors truncate"
						>
							{{ displayShortUrl }}
						</a>
						<button
							@click="copyShortUrl"
							class="flex-none p-1.5 rounded-lg border border-[hsl(240,11%,25%)] hover:border-[hsl(240,11%,40%)] bg-[hsl(240,11%,20%)] hover:bg-[hsl(240,11%,30%)] transition-colors"
							title="Copy short URL"
						>
							<Icon name="lucide:copy" class="w-3.5 h-3.5" />
						</button>
					</div>
					<span class="flex gap-1 items-center">
						<Icon
							name="lucide:corner-down-right"
							class="w-4 h-4 text-[hsl(240,11%,50%)] shrink-0"
						/>
						<a
							:href="link.url"
							target="_blank"
							class="text-sm font-medium text-[hsl(240,11%,50%)] hover:text-[hsl(240,11%,65%)] transition-colors truncate"
						>
							{{ displayUrl }}
						</a>
					</span>
				</div>
			</div>
			<div
				class="flex-none ml-3 flex gap-1.5 items-center justify-center text-sm px-2.5 py-1.5 rounded-lg font-semibold border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,20%)] transition-all"
				:class="{
					'ring-2 ring-green-500/50': liveViews !== null && liveViews !== props.link.views,
				}"
			>
				<Icon name="lucide:eye" class="w-4 h-4" />
				<span>{{ formatViews(displayViews) }}</span>
			</div>
		</div>

		<!-- Range selector and action buttons -->
		<div class="flex justify-between items-center mb-3">
			<div class="flex gap-2">
				<Button
					variant="outline"
					size="sm"
					@click="showViewsDialog = true"
					class="border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)]"
				>
					<Icon name="lucide:list" class="h-4 w-4 mr-2" />
					View History
				</Button>
				<Button
					variant="outline"
					size="sm"
					@click="showTopCountriesDialog = true"
					class="border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)]"
				>
					<Icon name="lucide:globe" class="h-4 w-4 mr-2" />
					Top Countries
				</Button>
			</div>
			<Select v-model="range">
				<SelectTrigger
					class="w-32 border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)] transition-colors"
				>
					<SelectValue placeholder="Range" />
				</SelectTrigger>
				<SelectContent>
					<SelectItem value="week">Week</SelectItem>
					<SelectItem value="month">Month</SelectItem>
					<SelectItem value="3months">3 Months</SelectItem>
				</SelectContent>
			</Select>
		</div>

		<!-- Chart section -->
		<div class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4">
			<div v-if="isLoading" class="flex items-center justify-center h-[200px]">
				<Icon name="lucide:loader-2" class="h-8 w-8 animate-spin text-[hsl(240,11%,50%)]" />
			</div>
			<div
				v-else-if="statistics.length === 0"
				class="flex items-center justify-center h-[200px] text-[hsl(240,11%,50%)]"
			>
				No data available
			</div>
			<UrlShortenerLinkChart v-else :is-day-range="isDayRange" :usages="statistics" />
		</div>

		<!-- Views History Dialog -->
		<ViewsHistoryDialog
			v-model:open="showViewsDialog"
			:short-link-id="link.id"
			:short-url="displayShortUrl"
		/>

		<!-- Top Countries Dialog -->
		<TopCountriesDialog
			v-model:open="showTopCountriesDialog"
			:short-link-id="link.id"
			:short-url="displayShortUrl"
		/>
	</div>
</template>
