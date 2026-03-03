<script setup lang="ts">
import { toast } from 'vue-sonner';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from '@/components/ui/dialog';
import Button from '@/components/ui/button/Button.vue';

const props = defineProps<{
	pasteId: string;
}>();

const emit = defineEmits<{
	(e: "deleted"): void;
}>();

const open = defineModel({ default: false })

const api = useOapi();
const isDeleting = ref(false);

async function handleDelete() {
	isDeleting.value = true;
	try {
		const result = await api.v1.pastebinDelete(props.pasteId);
		if (result.error) {
			toast.error("Error", {
				description: "Failed to delete paste",
			});
			return;
		}

		toast.success("Deleted", {
			description: "Paste has been deleted successfully",
		});

		open.value = false
		emit("deleted");
	} catch (error) {
		toast.error("Error", {
			description: "Failed to delete paste",
		});
	} finally {
		isDeleting.value = false;
	}
}

function handleCancel() {
	open.value = false
}
</script>

<template>
	<Dialog :open="open" @update:open="(val) => open = val">
		<DialogContent class="sm:max-w-md bg-[hsl(240,11%,9%)] border-[hsl(240,11%,18%)]">
			<DialogHeader>
				<DialogTitle>Delete Paste</DialogTitle>
				<DialogDescription>
					Are you sure you want to delete this paste? This action cannot be undone.
				</DialogDescription>
			</DialogHeader>
			<div class="flex justify-end gap-3 mt-4">
				<Button
					variant="outline"
					@click="handleCancel"
					:disabled="isDeleting"
					class="border-[hsl(240,11%,25%)] bg-[hsl(240,11%,15%)] hover:bg-[hsl(240,11%,20%)]"
				>
					Cancel
				</Button>
				<Button variant="destructive" @click="handleDelete" :disabled="isDeleting">
					<Icon v-if="isDeleting" name="lucide:loader-2" class="w-4 h-4 mr-2 animate-spin" />
					<Icon v-else name="lucide:trash-2" class="w-4 h-4 mr-2" />
					Delete
				</Button>
			</div>
		</DialogContent>
	</Dialog>
</template>
