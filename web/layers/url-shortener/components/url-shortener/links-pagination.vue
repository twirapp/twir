<script setup lang="ts">
interface Props {
	currentPage: number;
	totalPages: number;
	total: number;
}

const props = defineProps<Props>();

const emit = defineEmits<{
	(e: "goToPage", page: number): void;
	(e: "nextPage"): void;
	(e: "previousPage"): void;
}>();

// Calculate which page numbers to show
const visiblePages = computed(() => {
	const pages: (number | "ellipsis-start" | "ellipsis-end")[] = [];

	if (props.totalPages <= 7) {
		// Show all pages if 7 or less
		for (let i = 1; i <= props.totalPages; i++) {
			pages.push(i);
		}
	} else {
		// Always show first page
		pages.push(1);

		// Calculate range around current page
		const currentPageNumber = props.currentPage + 1;
		const start = Math.max(2, currentPageNumber - 1);
		const end = Math.min(props.totalPages - 1, currentPageNumber + 2);

		// Add ellipsis after first page if needed
		if (start > 2) {
			pages.push("ellipsis-start");
		}

		// Add pages around current page
		for (let i = start; i <= end; i++) {
			pages.push(i);
		}

		// Add ellipsis before last page if needed
		if (end < props.totalPages - 1) {
			pages.push("ellipsis-end");
		}

		// Always show last page
		pages.push(props.totalPages);
	}

	return pages;
});
</script>

<template>
	<div class="flex items-center justify-between border-t border-[hsl(240,11%,30%)] pt-4">
		<div class="text-sm text-[hsl(240,11%,65%)]">
			Page {{ currentPage + 1 }} of {{ totalPages }} ({{ total }} total links)
		</div>

		<div class="flex items-center gap-2">
			<button
				:disabled="currentPage === 0"
				@click="emit('previousPage')"
				class="px-3 py-2 rounded-lg border border-[hsl(240,11%,30%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				aria-label="Previous page"
			>
				<Icon name="lucide:chevron-left" class="w-4 h-4" />
			</button>

			<div class="flex items-center gap-1">
				<template v-for="(page, index) in visiblePages" :key="`page-${index}`">
					<button
						v-if="typeof page === 'number'"
						@click="emit('goToPage', page - 1)"
						:class="[
							'min-w-[40px] px-3 py-2 rounded-lg border transition-colors',
							currentPage === page - 1
								? 'border-[hsl(240,11%,45%)] bg-[hsl(240,11%,25%)] text-white'
								: 'border-[hsl(240,11%,30%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)]',
						]"
						:aria-label="`Go to page ${page}`"
						:aria-current="currentPage === page - 1 ? 'page' : undefined"
					>
						{{ page }}
					</button>
					<span v-else class="px-2 text-[hsl(240,11%,65%)]" aria-hidden="true"> ... </span>
				</template>
			</div>

			<button
				:disabled="currentPage >= totalPages - 1"
				@click="emit('nextPage')"
				class="px-3 py-2 rounded-lg border border-[hsl(240,11%,30%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,25%)] disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
				aria-label="Next page"
			>
				<Icon name="lucide:chevron-right" class="w-4 h-4" />
			</button>
		</div>
	</div>
</template>
