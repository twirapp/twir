<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { PlusIcon } from "lucide-vue-next";
import { useForm } from "vee-validate";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import { z } from "zod";

import type { Giveaway } from "@/api/giveaways.ts";

import { useUserAccessFlagChecker } from "@/api/auth";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { useGiveaways } from "@/features/giveaways/composables/giveaways-use-giveaways.ts";
import { ChannelRolePermissionEnum, GiveawayType } from "@/gql/graphql.ts";

const { t } = useI18n();
const open = ref(false);
const { createGiveaway, viewGiveaway } = useGiveaways();

const canManageGiveaways = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageGiveaways);

// Form validation schema
const formSchema = toTypedSchema(
	z
		.object({
			type: z.enum(["KEYWORD", "ONLINE_CHATTERS"]),
			keyword: z
				.string()
				.min(3, "Keyword must be at least 3 characters")
				.max(100, "Keyword must be at most 100 characters")
				.optional(),
			minWatchedTime: z.coerce.number().min(0).optional(),
			minMessages: z.coerce.number().min(0).optional(),
			minUsedChannelPoints: z.coerce.number().min(0).optional(),
			minFollowDuration: z.coerce.number().min(0).optional(),
			requireSubscription: z.boolean().optional(),
		})
		.refine(
			(data) => {
				if (data.type === "KEYWORD" && !data.keyword) {
					return false;
				}
				return true;
			},
			{
				message: "Keyword is required for KEYWORD type giveaways",
				path: ["keyword"],
			},
		),
);

// Form setup
const giveawayCreateForm = useForm({
	validationSchema: formSchema,
	initialValues: {
		type: "KEYWORD" as const,
		keyword: "",
		minWatchedTime: undefined,
		minMessages: undefined,
		minUsedChannelPoints: undefined,
		minFollowDuration: undefined,
		requireSubscription: false,
	},
	validateOnMount: false,
});

const selectedType = computed(() => giveawayCreateForm.values.type);

const handleSubmit = giveawayCreateForm.handleSubmit(async (values) => {
	try {
		const filters = {
			minWatchedTime: values.minWatchedTime,
			minMessages: values.minMessages,
			minUsedChannelPoints: values.minUsedChannelPoints,
			minFollowDuration: values.minFollowDuration,
			requireSubscription: values.requireSubscription,
		};

		// Remove undefined values
		const cleanFilters = Object.fromEntries(
			Object.entries(filters).filter(([_, v]) => v !== undefined),
		);

		const result = await createGiveaway({
			type: values.type as GiveawayType,
			keyword: values.keyword,
			filters: Object.keys(cleanFilters).length > 0 ? cleanFilters : undefined,
		});

		if (result) {
			giveawayCreateForm.resetForm();
			open.value = false;
			viewGiveaway((result as Giveaway).id);
		}
	} catch (error) {
		console.error(error);
	}
});
</script>

<template>
	<Dialog v-model:open="open">
		<DialogTrigger as-child>
			<Button size="sm" class="flex gap-2 items-center" :disabled="!canManageGiveaways">
				<PlusIcon class="size-4" />
				{{ t("giveaways.createNew") }}
			</Button>
		</DialogTrigger>

		<DialogContent class="sm:max-w-[525px] max-h-[90vh] overflow-y-auto">
			<DialogHeader>
				<DialogTitle>{{ t("giveaways.createDialog.title") }}</DialogTitle>
				<DialogDescription>
					{{ t("giveaways.createDialog.description") }}
				</DialogDescription>
			</DialogHeader>

			<form class="space-y-4" @submit.prevent="handleSubmit">
				<!-- Giveaway Type Selection -->
				<FormField v-slot="{ componentField }" name="type">
					<FormItem>
						<FormLabel>{{ t("giveaways.createDialog.typeLabel") }}</FormLabel>
						<Select v-bind="componentField">
							<FormControl>
								<SelectTrigger>
									<SelectValue placeholder="Select giveaway type" />
								</SelectTrigger>
							</FormControl>
							<SelectContent>
								<SelectItem value="KEYWORD">
									{{ t("giveaways.createDialog.typeKeyword") }}
								</SelectItem>
								<SelectItem value="ONLINE_CHATTERS">
									{{ t("giveaways.createDialog.typeOnlineChatters") }}
								</SelectItem>
							</SelectContent>
						</Select>
						<FormMessage />
					</FormItem>
				</FormField>

				<!-- Keyword Field (only for KEYWORD type) -->
				<FormField
					v-if="selectedType === 'KEYWORD'"
					v-slot="{ componentField, errorMessage }"
					name="keyword"
				>
					<FormItem>
						<FormLabel>{{ t("giveaways.createDialog.keywordLabel") }}</FormLabel>
						<FormControl>
							<Input
								:placeholder="t('giveaways.createDialog.keywordPlaceholder')"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage>{{ errorMessage }}</FormMessage>
					</FormItem>
				</FormField>

				<!-- Filters Section -->
				<div class="space-y-3 border-t pt-4">
					<h4 class="text-sm font-medium">{{ t("giveaways.createDialog.filtersTitle") }}</h4>

					<FormField v-slot="{ componentField }" name="minWatchedTime">
						<FormItem>
							<FormLabel>{{ t("giveaways.createDialog.minWatchedTime") }}</FormLabel>
							<FormControl>
								<Input
									type="number"
									:placeholder="t('giveaways.createDialog.minWatchedTimePlaceholder')"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="minMessages">
						<FormItem>
							<FormLabel>{{ t("giveaways.createDialog.minMessages") }}</FormLabel>
							<FormControl>
								<Input
									type="number"
									:placeholder="t('giveaways.createDialog.minMessagesPlaceholder')"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="minUsedChannelPoints">
						<FormItem>
							<FormLabel>{{ t("giveaways.createDialog.minUsedChannelPoints") }}</FormLabel>
							<FormControl>
								<Input
									type="number"
									:placeholder="t('giveaways.createDialog.minUsedChannelPointsPlaceholder')"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ componentField }" name="minFollowDuration">
						<FormItem>
							<FormLabel>{{ t("giveaways.createDialog.minFollowDuration") }}</FormLabel>
							<FormControl>
								<Input
									type="number"
									:placeholder="t('giveaways.createDialog.minFollowDurationPlaceholder')"
									v-bind="componentField"
								/>
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<FormField v-slot="{ value, handleChange }" name="requireSubscription">
						<FormItem class="flex items-center justify-between space-y-0 rounded-lg border p-3">
							<div class="space-y-0.5">
								<FormLabel>{{ t("giveaways.createDialog.requireSubscription") }}</FormLabel>
							</div>
							<FormControl>
								<Switch :model-value="value" @update:model-value="handleChange" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>
				</div>

				<DialogFooter>
					<Button type="button" variant="outline" @click="open = false">
						{{ t("giveaways.createDialog.cancel") }}
					</Button>
					<Button type="submit">
						{{ t("giveaways.createDialog.create") }}
					</Button>
				</DialogFooter>
			</form>
		</DialogContent>
	</Dialog>
</template>
