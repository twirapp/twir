<script setup lang="ts">
import { computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import { toast } from "vue-sonner";
import { useForm } from "vee-validate";
import { toTypedSchema } from "@vee-validate/zod";
import { z } from "zod";

import { useUserAccessFlagChecker } from "@/api/auth";
import { twitchSetChannelInformationMutation } from "@/api/twitch";
import TwitchCategorySelector from "@/components/twitch-category-selector.vue";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { FormField, FormItem, FormLabel, FormControl, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { ChannelRolePermissionEnum } from "@/gql/graphql";

const props = defineProps<{
	title?: string;
	categoryId?: string;
	categoryName?: string;
}>();

const open = defineModel<boolean>("open", { default: false });

const { t } = useI18n();

const formSchema = computed(() =>
	z.object({
		title: z.string().max(140, t("dashboard.statsWidgets.streamInfo.titleMaxLength")),
		categoryId: z.string().optional(),
	}),
);

const { handleSubmit, values, setValues, resetForm } = useForm({
	validationSchema: computed(() => toTypedSchema(formSchema.value)),
});

watch(
	() => props,
	(v) => {
		setValues({
			title: v.title ?? "",
			categoryId: v.categoryId ?? undefined,
		});
	},
	{ immediate: true, deep: true },
);

watch(
	() => open.value,
	(isOpen) => {
		if (isOpen) {
			resetForm({
				values: {
					title: props.title ?? "",
					categoryId: props.categoryId ?? undefined,
				},
			});
		}
	},
);

const titleLength = computed(() => values.title?.length ?? 0);

const informationUpdater = twitchSetChannelInformationMutation();

const onSubmit = handleSubmit(async (formValues) => {
	try {
		await informationUpdater.executeMutation({
			categoryId: formValues.categoryId,
			title: formValues.title,
		});
		toast.success(t("sharedTexts.saved"));
		open.value = false;
	} catch (error) {
		toast.error(error instanceof Error ? error.message : "Failed to save");
	}
});

const userCanEditTitle = useUserAccessFlagChecker(ChannelRolePermissionEnum.UpdateChannelTitle);
const userCanEditCategory = useUserAccessFlagChecker(
	ChannelRolePermissionEnum.UpdateChannelCategory,
);
</script>

<template>
	<Dialog v-model:open="open">
		<DialogContent class="sm:max-w-[500px]">
			<DialogHeader>
				<DialogTitle>{{ t("dashboard.statsWidgets.streamInfo.editStreamInfo") }}</DialogTitle>
				<DialogDescription>
					{{ t("dashboard.statsWidgets.streamInfo.editStreamInfoDescription") }}
				</DialogDescription>
			</DialogHeader>

			<form @submit="onSubmit">
				<div class="grid gap-4 py-4">
					<FormField v-slot="{ componentField }" name="title">
						<FormItem>
							<FormLabel>
								{{ t("dashboard.statsWidgets.streamInfo.title") }}
							</FormLabel>
							<FormControl>
								<div class="space-y-1">
									<Input
										id="title"
										v-bind="componentField"
										:disabled="!userCanEditTitle"
										:placeholder="t('dashboard.statsWidgets.streamInfo.title')"
									/>
									<div class="flex justify-end">
										<span
											class="text-xs"
											:class="titleLength > 140 ? 'text-destructive' : 'text-muted-foreground'"
										>
											{{ titleLength }} / 140
										</span>
									</div>
								</div>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="categoryId">
						<FormItem>
							<FormLabel for="category">
								{{ t("dashboard.statsWidgets.streamInfo.category") }}
							</FormLabel>
							<FormControl>
								<TwitchCategorySelector
									id="category"
									v-bind="componentField"
									:disabled="!userCanEditCategory"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<DialogFooter>
					<Button type="button" variant="outline" @click="open = false">
						{{ t("sharedButtons.cancel") }}
					</Button>
					<Button type="submit" :disabled="informationUpdater.fetching.value">
						{{
							informationUpdater.fetching.value
								? t("sharedButtons.saving")
								: t("sharedButtons.save")
						}}
					</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
