<script setup lang="ts">
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import Button from "@/components/ui/button/Button.vue";
import { useOapi } from "~/composables/use-oapi";
import { toast } from "vue-sonner";

const props = defineProps<{
	open: boolean;
	linkId: string;
	shortUrl: string;
}>();

const emit = defineEmits<{
	(e: "update:open", value: boolean): void;
	(e: "deleted"): void;
}>();

const api = useOapi();
const isDeleting = ref(false);

async function handleDelete() {
	isDeleting.value = true;

	try {
		await api.v1.shortUrlDelete(props.linkId);

		toast.success("Deleted", {
			description: "Short link deleted successfully",
		});

		emit("deleted");
		closeDialog();
	} catch (err: any) {
		console.error("Failed to delete link:", err);
		const errorMessage = err?.data?.detail || "Failed to delete link";
		toast.error("Error", {
			description: errorMessage,
		});
	} finally {
		isDeleting.value = false;
	}
}

function closeDialog() {
	emit("update:open", false);
}
</script>

<template>
	<Dialog :open="open" @update:open="closeDialog">
		<DialogContent class="max-w-md">
			<DialogHeader>
				<DialogTitle>Delete Short Link</DialogTitle>
				<DialogDescription>
					Are you sure you want to delete <strong class="font-semibold">{{ shortUrl }}</strong
					>? This action cannot be undone.
				</DialogDescription>
			</DialogHeader>

			<DialogFooter>
				<Button type="button" variant="outline" @click="closeDialog" :disabled="isDeleting">
					Cancel
				</Button>
				<Button type="button" variant="destructive" @click="handleDelete" :disabled="isDeleting">
					<Icon v-if="isDeleting" name="lucide:loader-2" class="h-4 w-4 mr-2 animate-spin" />
					Delete
				</Button>
			</DialogFooter>
		</DialogContent>
	</Dialog>
</template>
