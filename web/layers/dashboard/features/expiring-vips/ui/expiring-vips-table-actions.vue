<script setup lang="ts">
import { TrashIcon } from "lucide-vue-next";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { toast } from "vue-sonner";

import type { ScheduledVip } from "@/api/scheduled-vips.ts";

import { useScheduledVipsApi } from "@/api/scheduled-vips.ts";
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";

const props = defineProps<{
	scheduledVip: ScheduledVip;
}>();

const { t } = useI18n();
const api = useScheduledVipsApi();
const { executeMutation: deleteMutation, fetching: isDeleting } =
	api.useMutationRemoveScheduledVip();

const showDeleteDialog = ref(false);

const handleDelete = async () => {
	const result = await deleteMutation({
		id: props.scheduledVip.id,
		input: {
			keepVip: false,
		},
	});

	if (result.error) {
		toast.error(result.error.message);
		return;
	}

	toast.success(t("expiringVips.successDelete"));
	showDeleteDialog.value = false;
};
</script>

<template>
	<div class="flex items-center gap-2">
		<Button
			variant="destructive"
			size="icon"
			@click="showDeleteDialog = true"
			:disabled="isDeleting"
		>
			<TrashIcon class="h-4 w-4" />
		</Button>

		<AlertDialog v-model:open="showDeleteDialog">
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>{{ t("expiringVips.deleteDialog.title") }}</AlertDialogTitle>
					<AlertDialogDescription>
						{{
							t("expiringVips.deleteDialog.description", {
								user: scheduledVip.twitchProfile.displayName,
							})
						}}
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>{{ t("sharedButtons.cancel") }}</AlertDialogCancel>
					<AlertDialogAction @click="handleDelete" :disabled="isDeleting">
						{{ t("sharedButtons.delete") }}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	</div>
</template>
