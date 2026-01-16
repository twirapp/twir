<script setup lang="ts">
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { useOapi } from "~/composables/use-oapi";

const props = defineProps<{
	open: boolean;
	shortLinkId: string;
	shortUrl: string;
}>();

const emit = defineEmits<{
	(e: "update:open", value: boolean): void;
}>();

const api = useOapi();
const countries = ref<Array<{ country: string; count: number }>>([]);
const isLoading = ref(false);

// Fetch data when dialog opens
watch(
	() => props.open,
	async (isOpen) => {
		if (isOpen) {
			isLoading.value = true;
			try {
				const response = await api.v1.shortUrlGetTopCountries(props.shortLinkId, {
					limit: 10,
				});
				countries.value = response.data.data ?? [];
			} catch (err) {
				console.error("Failed to fetch top countries:", err);
				countries.value = [];
			} finally {
				isLoading.value = false;
			}
		}
	},
);

function formatNumber(num: number) {
	return new Intl.NumberFormat(import.meta.client ? navigator.language : "en-US").format(num);
}

function getCountryEmoji(countryCode: string) {
	if (!countryCode || countryCode.length !== 2) return "🌍";
	const codePoints = countryCode
		.toUpperCase()
		.split("")
		.map((char) => 127397 + char.charCodeAt(0));
	return String.fromCodePoint(...codePoints);
}

function closeDialog() {
	emit("update:open", false);
}
</script>

<template>
	<Dialog :open="open" @update:open="closeDialog">
		<DialogContent class="max-w-md">
			<DialogHeader>
				<DialogTitle>Top Countries</DialogTitle>
				<DialogDescription> Most popular countries viewing {{ shortUrl }} </DialogDescription>
			</DialogHeader>

			<div v-if="isLoading" class="flex items-center justify-center py-8">
				<Icon name="lucide:loader-2" class="h-8 w-8 animate-spin text-muted-foreground" />
			</div>

			<div v-else-if="countries.length === 0" class="py-8 text-center text-muted-foreground">
				No country data available
			</div>

			<div v-else class="space-y-2">
				<div
					v-for="(item, index) in countries"
					:key="item.country"
					class="flex items-center justify-between p-3 rounded-lg border hover:bg-accent/50 transition-colors"
				>
					<div class="flex items-center gap-3">
						<span class="text-2xl">{{ getCountryEmoji(item.country) }}</span>
						<div class="flex flex-col">
							<span class="font-medium">{{ item.country || "Unknown" }}</span>
							<span class="text-xs text-muted-foreground">Rank #{{ index + 1 }}</span>
						</div>
					</div>
					<div class="flex items-center gap-2">
						<span class="text-sm font-semibold">{{ formatNumber(item.count) }}</span>
						<span class="text-xs text-muted-foreground">views</span>
					</div>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
