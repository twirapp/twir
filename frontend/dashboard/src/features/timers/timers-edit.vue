<script setup lang="ts">
import { toTypedSchema } from "@vee-validate/zod";
import { BadgePlus, GripVertical, TrashIcon } from "lucide-vue-next";
import { FieldArray, useForm } from "vee-validate";
import { computed, onMounted, ref, toRaw } from "vue";
import { VueDraggable } from "vue-draggable-plus";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";

import { Button } from "@/components/ui/button";
import { Card, CardAction, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { Slider } from "@/components/ui/slider";
import Switch from "@/components/ui/switch/Switch.vue";
import VariableInput from "@/components/variable-input.vue";
import { formSchema, useTimersEdit } from "@/features/timers/composables/use-timers-edit.js";
import { TwitchAnnounceColor } from "@/gql/graphql.js";
import PageLayout from "@/layout/page-layout.vue";
import { Separator } from "@/components/ui/separator";

const route = useRoute();
const { t } = useI18n();
const { findTimer, submit } = useTimersEdit();

const loading = ref(true);

const { resetForm, handleSubmit, controlledValues, errors, setValues } = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		timeInterval: 1,
		messageInterval: 0,
		responses: [{ text: "", isAnnounce: false, count: 1 }],
		offlineEnabled: false,
		onlineEnabled: true,
	},
});

onMounted(async () => {
	resetForm();

	if (typeof route.params.id === "string") {
		const timer = await findTimer(route.params.id);
		if (timer) {
			setValues(toRaw(timer));
		}
	}

	loading.value = false;
});

const onSubmit = handleSubmit(submit);

const responsesHasError = computed(() => {
	return Object.keys(errors.value).some((key) => key.startsWith("responses"));
});
</script>

<template>
	<form @submit="onSubmit">
		<PageLayout sticky-header show-back backRedirectTo="/dashboard/timers">
			<template #title>
				{{ route.params.id === "create" ? t("sharedTexts.create") : t("sharedTexts.edit") }}
			</template>

			<template #action>
				<Button type="submit">
					{{ t("sharedButtons.save") }}
				</Button>
			</template>

			<template #content>
				<div class="flex flex-col gap-4 max-w-4xl mx-auto" :class="{ 'blur-xs': loading }">
					<FormField v-slot="{ componentField }" name="name">
						<FormItem>
							<FormLabel>{{ t("sharedTexts.name") }}</FormLabel>
							<FormControl>
								<Input type="text" v-bind="componentField" />
							</FormControl>
							<FormMessage />
						</FormItem>
					</FormField>

					<Card class="p-0">
						<CardContent class="py-4 space-y-4">
							<div class="space-y-2">
								<h3 class="text-lg font-semibold">Timer Intervals</h3>
								<p class="text-sm text-muted-foreground">
									Both intervals work <strong>together (AND)</strong>. The timer will trigger when
									<strong>both</strong> the time interval <strong>AND</strong> message interval
									conditions are met.
								</p>
							</div>

							<FormField v-slot="{ componentField }" name="timeInterval">
								<FormItem>
									<FormLabel>{{ t("timers.table.columns.intervalInMinutes") }}</FormLabel>
									<FormControl>
										<div class="flex gap-6 flex-wrap">
											<Input type="number" v-bind="componentField" />
											<Slider
												:model-value="[componentField.modelValue]"
												:max="1000"
												:default-value="[0, 1000]"
												:min="0"
												:step="1"
												@update:model-value="
													(v) => {
														if (!v) return;
														componentField.onChange(v[0]);
													}
												"
											/>
										</div>
									</FormControl>
									<FormMessage />
									<FormDescription class="flex justify-end">
										<span>{{ componentField.modelValue }} minutes</span>
									</FormDescription>
								</FormItem>
							</FormField>

							<Separator class="my-4" />

							<FormField v-slot="{ componentField }" name="messageInterval">
								<FormItem>
									<FormLabel>{{ t("timers.table.columns.intervalInMessages") }}</FormLabel>
									<FormControl>
										<Input type="number" placeholder="0" v-bind="componentField" />
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<div class="mt-4 p-4 bg-muted rounded-lg space-y-3">
								<h4 class="text-sm font-semibold">Examples:</h4>
								<div class="space-y-2 text-sm">
									<div class="flex items-start gap-2">
										<span class="text-primary">•</span>
										<div>
											<strong>Time: {{ controlledValues.timeInterval }} min, Messages: 0</strong>
											<br />
											<span class="text-muted-foreground">
												Triggers every {{ controlledValues.timeInterval }} minutes (no message
												requirement)
											</span>
										</div>
									</div>
									<div class="flex items-start gap-2">
										<span class="text-primary">•</span>
										<div>
											<strong
												>Time: {{ controlledValues.timeInterval }} min, Messages:
												{{ controlledValues.messageInterval || 10 }}</strong
											>
											<br />
											<span class="text-muted-foreground">
												Triggers after {{ controlledValues.timeInterval }} minutes
												<strong>AND</strong> {{ controlledValues.messageInterval || 10 }} messages
												have been sent in chat
											</span>
										</div>
									</div>
									<div class="flex items-start gap-2">
										<span class="text-primary">•</span>
										<div>
											<strong>Time: 10 min, Messages: 50</strong>
											<br />
											<span class="text-muted-foreground">
												If 10 minutes pass but only 30 messages were sent, the timer will
												<strong>NOT</strong> trigger until 50 messages are reached
											</span>
										</div>
									</div>
								</div>
							</div>
						</CardContent>
					</Card>

					<Card class="p-0">
						<CardContent class="py-4 space-y-4">
							<div class="space-y-2">
								<h3 class="text-lg font-semibold">Delivery</h3>
								<p class="text-sm text-muted-foreground">
									Choose when the timer is allowed to send messages.
								</p>
							</div>

							<FormField v-slot="{ value, handleChange }" name="onlineEnabled">
								<FormItem class="space-y-2 rounded-lg border p-4">
									<div class="flex items-center justify-between gap-4">
										<div class="space-y-1">
											<FormLabel>Send while online</FormLabel>
											<FormDescription>
												Allow this timer to trigger when the channel is live.
											</FormDescription>
										</div>
										<FormControl>
											<Switch :model-value="value" @update:model-value="handleChange" />
										</FormControl>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>

							<FormField v-slot="{ value, handleChange }" name="offlineEnabled">
								<FormItem class="space-y-2 rounded-lg border p-4">
									<div class="flex items-center justify-between gap-4">
										<div class="space-y-1">
											<FormLabel>Send while offline</FormLabel>
											<FormDescription>
												Allow this timer to trigger when the channel is offline.
											</FormDescription>
										</div>
										<FormControl>
											<Switch :model-value="value" @update:model-value="handleChange" />
										</FormControl>
									</div>
									<FormMessage />
								</FormItem>
							</FormField>
						</CardContent>
					</Card>

					<Label :class="{ 'text-destructive': responsesHasError }">{{
						t("sharedTexts.responses")
					}}</Label>
					<span class="text-sm text-muted-foreground">
						Responses are sent in sequence: the first on the initial trigger, the second after
						<b>{{ controlledValues.timeInterval }}</b> minutes, etc., cycling back to the first
						after last are sent.
					</span>

					<FieldArray v-slot="{ fields, push, remove }" name="responses">
						<VueDraggable
							v-model="controlledValues.responses!"
							handle=".drag-handle"
							class="flex flex-col gap-2"
						>
							<div v-for="(field, index) in fields" :key="`responses-text-${field.key}`">
								<Card class="relative flex items-center p-0">
									<div
										class="absolute flex left-0 rounded-l-md h-full bg-accent w-4 cursor-move drag-handle z-10 border cursor-move"
									>
										<GripVertical class="my-auto size-6" />
									</div>
									<CardContent class="pt-2 w-full">
										<FormField v-slot="{ componentField }" :name="`responses[${index}].text`">
											<FormItem class="flex flex-col gap-4">
												<FormControl>
													<VariableInput
														input-type="textarea"
														class="relative p-2"
														:model-value="componentField.modelValue"
														:min-rows="1"
														:rows="1"
														popoverAlign="end"
														popoverSide="bottom"
														@update:model-value="componentField.onChange"
													/>

													<div class="flex flex-row flex-wrap gap-4">
														<div class="flex flex-col gap-2">
															<Label :for="`responses[${index}].count`">
																How many times send this message on trigger
															</Label>

															<Input
																:id="`responses[${index}].count`"
																v-model:modelValue="(field.value as any).count"
																type="number"
															/>
														</div>

														<div class="flex flex-col gap-2">
															<Label :for="`responses[${index}].isAnnounce`">
																Send as announcement
															</Label>
															<Checkbox
																:id="`responses[${index}].isAnnounce`"
																v-model:model-value="(field.value as any).isAnnounce"
															/>
														</div>

														<div class="flex flex-col gap-2">
															<Label :for="`responses[${index}].announceColor`">
																Announcement color
															</Label>
															<Select
																:id="`responses[${index}].announceColor`"
																v-model:modelValue="(field.value as any).announceColor"
																:default-value="TwitchAnnounceColor.Primary"
															>
																<SelectTrigger :disabled="!(field.value as any).isAnnounce">
																	<SelectValue placeholder="Select a color" />
																</SelectTrigger>
																<SelectContent>
																	<SelectGroup>
																		<SelectItem
																			v-for="color of TwitchAnnounceColor"
																			:key="color"
																			:value="color"
																		>
																			{{
																				color.at(0)!.toUpperCase() + color.slice(1).toLowerCase()
																			}}
																		</SelectItem>
																	</SelectGroup>
																</SelectContent>
															</Select>
														</div>
													</div>
												</FormControl>
												<FormMessage />
											</FormItem>
										</FormField>
										<CardAction class="w-full flex justify-end py-2">
											<Button
												variant="destructive"
												class="flex gap-2 place-self-end"
												@click="remove(index)"
											>
												<TrashIcon class="size-4" />
												Remove
											</Button>
										</CardAction>
									</CardContent>
								</Card>
							</div>
						</VueDraggable>
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="text-xs w-full flex gap-2 items-center mt-2"
							:disabled="(controlledValues.responses?.length ?? 0) >= 10"
							@click="push({ text: '', isAnnounce: false, count: 1 })"
						>
							<BadgePlus class="size-4" />
							Add response {{ controlledValues.responses?.length ?? 0 }} / 10
						</Button>
					</FieldArray>
				</div>
			</template>
		</PageLayout>
	</form>
</template>
