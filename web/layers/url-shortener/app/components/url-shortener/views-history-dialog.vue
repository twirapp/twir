<script setup lang="ts">
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow,
} from "@/components/ui/table";
import Button from "@/components/ui/button/Button.vue";
import { useShortLinkViews } from "../../composables/use-short-link-views";

const props = defineProps<{
	open: boolean;
	shortLinkId: string;
	shortUrl: string;
}>();

const emit = defineEmits<{
	(e: "update:open", value: boolean): void;
}>();

const page = ref(0);
const perPage = ref(20);

const { views, total, executeQuery } = useShortLinkViews(
	computed(() => props.shortLinkId),
	page,
	perPage,
);

// Fetch data when dialog opens
watch(
	() => props.open,
	(isOpen) => {
		if (isOpen) {
			page.value = 0;
			executeQuery();
		}
	},
);

const totalPages = computed(() => Math.ceil(total.value / perPage.value));
const hasNextPage = computed(() => page.value < totalPages.value - 1);
const hasPrevPage = computed(() => page.value > 0);

function nextPage() {
	if (hasNextPage.value) {
		page.value++;
		executeQuery();
	}
}

function prevPage() {
	if (hasPrevPage.value) {
		page.value--;
		executeQuery();
	}
}

function formatDate(dateString: string) {
	const date = new Date(dateString);
	return new Intl.DateTimeFormat(import.meta.client ? navigator.language : "en-US", {
		dateStyle: "short",
		timeStyle: "short",
	}).format(date);
}

function closeDialog() {
	emit("update:open", false);
}
</script>

<template>
	<Dialog :open="open" @update:open="closeDialog">
		<DialogContent class="max-w-3xl max-h-[80vh] flex flex-col">
			<DialogHeader>
				<DialogTitle>View History</DialogTitle>
				<DialogDescription> Detailed view history for {{ shortUrl }} </DialogDescription>
			</DialogHeader>

			<div class="flex-1 overflow-auto">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead>Date</TableHead>
							<TableHead>Country</TableHead>
							<TableHead>City</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						<TableRow v-if="views.length === 0">
							<TableCell colspan="3" class="text-center text-muted-foreground">
								No views yet
							</TableCell>
						</TableRow>
						<TableRow v-for="view in views" :key="`${view.createdAt}-${view.userId}`">
							<TableCell>{{ formatDate(view.createdAt) }}</TableCell>
							<TableCell>{{ view.country || "Unknown" }}</TableCell>
							<TableCell>{{ view.city || "Unknown" }}</TableCell>
						</TableRow>
					</TableBody>
				</Table>
			</div>

			<div class="flex items-center justify-between pt-4 border-t">
				<div class="text-sm text-muted-foreground">
					Page {{ page + 1 }} of {{ totalPages || 1 }} ({{ total }} total views)
				</div>
				<div class="flex gap-2">
					<Button variant="outline" size="sm" :disabled="!hasPrevPage" @click="prevPage">
						<Icon name="lucide:chevron-left" class="h-4 w-4 mr-1" />
						Previous
					</Button>
					<Button variant="outline" size="sm" :disabled="!hasNextPage" @click="nextPage">
						Next
						<Icon name="lucide:chevron-right" class="h-4 w-4 ml-1" />
					</Button>
				</div>
			</div>
		</DialogContent>
	</Dialog>
</template>
