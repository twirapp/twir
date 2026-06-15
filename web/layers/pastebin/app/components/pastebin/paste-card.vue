<script setup lang="ts">
import type { PasteBinOutputDto } from "@twir/api/openapi";
import { toast } from "vue-sonner";
import Button from "@/components/ui/button/Button.vue";
import DeletePasteDialog from "./delete-paste-dialog.vue";

const props = defineProps<{
	paste: PasteBinOutputDto;
}>();

const emit = defineEmits<{
	(e: "deleted"): void;
}>();

const requestUrl = useRequestURL();
const clipboardApi = useClipboard();
const showDeleteDialog = ref(false);
const { detectLanguage, highlight } = useHighlight();

const pasteUrl = computed(() => `${requestUrl.origin}/h/${props.paste.id}`);

function copyPasteUrl() {
	clipboardApi.copy(pasteUrl.value);
	toast.success("Copied", {
		description: "Paste URL copied to clipboard",
		duration: 2500,
	});
}

function formatDate(date: string) {
	return new Date(date).toLocaleString();
}

const previewContent = computed(() => {
	const maxLength = 200;
	if (props.paste.content.length > maxLength) {
		return `${props.paste.content.slice(0, maxLength)}...`;
	}
	return props.paste.content;
});

const highlightedPreview = computed(() => {
	const lang = detectLanguage(previewContent.value);
	return highlight(previewContent.value, lang);
});

const hasExpiration = computed(() => !!props.paste.expire_at);

const isExpired = computed(() => {
	if (!props.paste.expire_at) return false;
	return new Date(props.paste.expire_at) < new Date();
});
</script>

<template>
	<div
		class="flex flex-col bg-[hsl(240,11%,9%)] border border-[hsl(240,11%,18%)] w-full max-w-full rounded-2xl p-4 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
		:class="{ 'opacity-60': isExpired }"
	>
		<!-- Header with paste info and actions -->
		<div class="flex justify-between items-start mb-4">
			<div class="flex items-center min-w-0 gap-x-3 flex-1">
				<div
					class="flex-none size-fit p-3 rounded-full font-semibold border border-[hsl(240,11%,25%)] bg-[hsl(240,11%,20%)]"
				>
					<Icon name="lucide:file-text" class="w-4 h-4" />
				</div>
				<div class="min-w-0 flex-1">
					<div class="flex items-center gap-2 overflow-hidden min-w-0">
						<a
							:href="pasteUrl"
							target="_blank"
							class="font-bold hover:text-[hsl(240,11%,85%)] transition-colors truncate"
						>
							{{ paste.id }}
						</a>
						<button
							@click="copyPasteUrl"
							class="flex-none p-1.5 rounded-lg border border-[hsl(240,11%,25%)] hover:border-[hsl(240,11%,40%)] bg-[hsl(240,11%,20%)] hover:bg-[hsl(240,11%,30%)] transition-colors"
							title="Copy paste URL"
						>
							<Icon name="lucide:copy" class="w-3.5 h-3.5" />
						</button>
						<button
							@click="showDeleteDialog = true"
							class="flex-none p-1.5 rounded-lg border border-red-900/50 hover:border-red-700 bg-red-950/30 hover:bg-red-950/50 text-red-400 hover:text-red-300 transition-colors"
							title="Delete paste"
						>
							<Icon name="lucide:trash-2" class="w-3.5 h-3.5" />
						</button>
					</div>
					<div class="flex gap-3 items-center text-sm text-[hsl(240,11%,50%)] mt-1">
						<span class="flex gap-1 items-center">
							<Icon name="lucide:calendar" class="w-3.5 h-3.5 shrink-0" />
							{{ formatDate(paste.created_at) }}
						</span>
						<span v-if="hasExpiration" class="flex gap-1 items-center">
							<Icon
								:name="isExpired ? 'lucide:alert-circle' : 'lucide:clock'"
								class="w-3.5 h-3.5 shrink-0"
								:class="{ 'text-red-400': isExpired }"
							/>
							<span :class="{ 'text-red-400': isExpired }">
								{{ isExpired ? "Expired" : "Expires" }}:
								{{ paste.expire_at ? formatDate(paste.expire_at) : "" }}
							</span>
						</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Content preview -->
		<div
			class="rounded-xl border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,12%)] p-4 overflow-hidden"
		>
			<div class="mb-2 text-xs text-[hsl(240,11%,50%)] font-semibold uppercase tracking-wider">
				Content Preview
			</div>
			<pre
				class="text-sm whitespace-pre-wrap break-words font-mono overflow-hidden"
			><code v-html="highlightedPreview"></code></pre>
			<div v-if="paste.content.length > 200" class="mt-3">
				<a
					:href="pasteUrl"
					target="_blank"
					class="inline-flex items-center gap-1 text-sm text-[hsl(240,11%,85%)] hover:text-white transition-colors"
				>
					View full content
					<Icon name="lucide:arrow-right" class="w-3.5 h-3.5" />
				</a>
			</div>
		</div>

		<!-- Delete Dialog -->
		<ClientOnly>
			<DeletePasteDialog
				v-model:open="showDeleteDialog"
				:paste-id="paste.id"
				@deleted="emit('deleted')"
			/>
		</ClientOnly>
	</div>
</template>

<style scoped>
:deep(code) {
	font-family: "JetBrains Mono", monospace;
}
</style>
