<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { InfoIcon } from "lucide-vue-next";
import { useForm } from "vee-validate";
import { computed, watch } from "vue";
import { useI18n } from "vue-i18n";
import { toast } from "vue-sonner";
import { z } from "zod";

import { useScheduledVipsApi } from "@/api/scheduled-vips.ts";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { DatePicker } from "@/components/ui/date-picker";
import {
	Dialog,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import TwitchUserSelect from "@/components/twitchUsers/twitch-user-select.vue";
import { ScheduledVipRemoveType } from "@/gql/graphql";
import DialogOrSheet from "@/components/dialog-or-sheet.vue";

const open = defineModel<boolean>({ default: false });

const { t } = useI18n();
const api = useScheduledVipsApi();
const { executeMutation: createMutation, fetching: isCreating } =
	api.useMutationCreateScheduledVip();

const formSchema = z
	.object({
		userID: z.string().min(1, t("expiringVips.form.errors.userRequired")).nullable(),
		removeType: z.nativeEnum(ScheduledVipRemoveType),
		removeAt: z.number().optional().nullable(),
	})
	.superRefine((data, ctx) => {
		if (!data.userID) {
			ctx.addIssue({
				code: z.ZodIssueCode.custom,
				message: t("expiringVips.form.errors.userRequired"),
				path: ["userID"],
			});
		}
		if (data.removeType === ScheduledVipRemoveType.Time && !data.removeAt) {
			ctx.addIssue({
				code: z.ZodIssueCode.custom,
				message: t("expiringVips.form.errors.dateRequired"),
				path: ["removeAt"],
			});
		}
	});

const { handleSubmit, resetForm, values, setFieldValue } = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		userID: null as string | null,
		removeType: ScheduledVipRemoveType.Time,
		removeAt: null,
	},
});

const showDatePicker = computed(() => {
	return values.removeType === ScheduledVipRemoveType.Time;
});

watch(
	() => values.removeType,
	(newType) => {
		if (newType === ScheduledVipRemoveType.StreamEnd) {
			setFieldValue("removeAt", null);
		}
	},
);

const onSubmit = handleSubmit(async (formValues) => {
	if (!formValues.userID) {
		toast.error(t("expiringVips.form.errors.userRequired"));
		return;
	}

	const result = await createMutation({
		input: {
			userID: formValues.userID,
			removeType: formValues.removeType,
			removeAt: formValues.removeAt,
		},
	});

	if (result.error) {
		toast.error(t("sharedTexts.error"), {
			description: result.error.message,
		});
		return;
	}

	toast.success(t("sharedTexts.success"), {
		description: t("expiringVips.form.successCreate"),
	});

	resetForm();
	open.value = false;
});
</script>

<template>
	<Dialog v-model:open="open">
		<DialogOrSheet>
			<DialogHeader>
				<DialogTitle>{{ t("expiringVips.form.createTitle") }}</DialogTitle>
				<DialogDescription>
					{{ t("expiringVips.form.createDescription") }}
				</DialogDescription>
			</DialogHeader>

			<form @submit="onSubmit" class="space-y-4">
				<Alert>
					<InfoIcon class="h-4 w-4" />
					<AlertDescription>
						{{ t("expiringVips.form.info") }}
					</AlertDescription>
				</Alert>

				<FormField v-slot="{ componentField }" name="userID">
					<FormItem>
						<FormLabel>{{ t("expiringVips.form.user") }}</FormLabel>
						<FormControl>
							<TwitchUserSelect
								v-model="componentField.modelValue"
								@update:model-value="componentField['onUpdate:modelValue']"
								:placeholder="t('expiringVips.form.userPlaceholder')"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-slot="{ componentField }" name="removeType">
					<FormItem>
						<FormLabel>{{ t("expiringVips.form.removeType") }}</FormLabel>
						<Select v-bind="componentField">
							<FormControl>
								<SelectTrigger>
									<SelectValue :placeholder="t('expiringVips.form.removeTypePlaceholder')" />
								</SelectTrigger>
							</FormControl>
							<SelectContent>
								<SelectGroup>
									<SelectItem :value="ScheduledVipRemoveType.Time">
										{{ t("expiringVips.form.removeTypeTime") }}
									</SelectItem>
									<SelectItem :value="ScheduledVipRemoveType.StreamEnd">
										{{ t("expiringVips.form.removeTypeStreamEnd") }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField v-if="showDatePicker" v-slot="{ field }" name="removeAt">
					<FormItem>
						<FormLabel>{{ t("expiringVips.form.removeAt") }}</FormLabel>
						<FormControl>
							<DatePicker
								:uid="field.name"
								auto-apply
								model-type="timestamp"
								dark
								:model-value="field.value"
								:min-date="new Date()"
								:config="{ closeOnAutoApply: true }"
								:placeholder="t('expiringVips.form.pickDate')"
								@update:model-value="field['onUpdate:modelValue']"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<DialogFooter>
					<Button type="submit" :disabled="isCreating">
						{{ t("sharedButtons.create") }}
					</Button>
				</DialogFooter>
			</form>
		</DialogOrSheet>
	</Dialog>
</template>
