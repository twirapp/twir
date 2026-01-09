<script setup lang="ts">
import { Bomb } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { onMounted, ref, toRaw, watch } from 'vue'


import { formSchema, useDuelForm } from './composables/use-duel-form'

import DialogOrSheet from '#layers/dashboard/components/dialog-or-sheet.vue'
import Card from '#layers/dashboard/components/games/card.vue'
import CommandButton from '~/features/commands/ui/command-button.vue'

const isDialogOpen = ref(false)
const { t } = useI18n()
const { initialValues, save, settings } = useDuelForm()

const duelForm = useForm({
	validationSchema: formSchema,
	initialValues,
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

watch(
	settings,
	(newSettings) => {
		if (!newSettings) return
		duelForm.setValues(toRaw(newSettings.gamesDuel))
	},
	{ immediate: true }
)

onMounted(() => {
	if (!settings.value) return
	duelForm.setValues(structuredClone(toRaw(settings.value.gamesDuel)))
})

const onSubmit = duelForm.handleSubmit(async (values) => {
	await save(values)
	isDialogOpen.value = false
})

function resetSettings() {
	duelForm.setValues(initialValues)
}
</script>

<template>
	<UiDialog v-model:open="isDialogOpen">
		<UiDialogTrigger asChild>
			<Card title="Duel" :icon="Bomb" :icon-stroke="1" :description="t('games.duel.description')" />
		</UiDialogTrigger>

		<DialogOrSheet class="sm:max-w-[625px]">
			<UiDialogHeader>
				<UiDialogTitle>{{ t('games.duel.title') }}</UiDialogTitle>
			</UiDialogHeader>

			<form>
				<div class="grid gap-4 py-4">
					<UiFormField v-slot="{ value, handleChange }" name="enabled">
						<UiFormItem
							class="flex flex-row items-center justify-between rounded-lg border p-4 space-y-0"
						>
							<UiFormLabel>Enabled</UiFormLabel>
							<UiFormControl>
								<UiSwitch :model-value="value" @update:model-value="handleChange" />
							</UiFormControl>
						</UiFormItem>
					</UiFormField>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.commands.title') }}
						</h4>
						<div class="flex gap-2">
							<CommandButton name="duel" :title="t('games.duel.commands.duel')" />
							<CommandButton name="duel accept" :title="t('games.duel.commands.accept')" />
							<CommandButton name="duel stats" :title="t('games.duel.commands.stats')" />
						</div>
					</div>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.cooldown.title') }}
						</h4>
						<div class="grid grid-cols-2 gap-4">
							<UiFormField v-slot="{ componentField }" name="userCooldown">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.cooldown.user') }}</UiFormLabel>
									<UiFormControl>
										<UiInput type="number" v-bind="componentField" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="globalCooldown">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.cooldown.global') }}</UiFormLabel>
									<UiFormControl>
										<UiInput type="number" v-bind="componentField" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>
						</div>
					</div>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.settings.title') }}
						</h4>
						<div class="grid grid-cols-2 gap-4">
							<UiFormField v-slot="{ componentField }" name="timeoutSeconds">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.settings.timeoutTime') }}</UiFormLabel>
									<UiFormControl>
										<UiInput type="number" v-bind="componentField" :max="84000" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="secondsToAccept">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.settings.secondsToAccept') }}</UiFormLabel>
									<UiFormControl>
										<UiInput type="number" v-bind="componentField" :max="3600" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="bothDiePercent">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.settings.bothDiePercent') }}</UiFormLabel>
									<UiFormControl>
										<UiInput type="number" v-bind="componentField" :max="100" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<!--							<FormField -->
							<!--								v-slot="{ componentField }" -->
							<!--								name="pointsPerWin" -->
							<!--							> -->
							<!--								<FormItem> -->
							<!--									<FormLabel>{{ t('games.duel.settings.pointsPerWin') }}</FormLabel> -->
							<!--									<FormControl> -->
							<!--										<Input type="number" v-bind="componentField" :max="999999" /> -->
							<!--									</FormControl> -->
							<!--									<FormMessage /> -->
							<!--								</FormItem> -->
							<!--							</FormField> -->

							<!--							<FormField -->
							<!--								v-slot="{ componentField }" -->
							<!--								name="pointsPerLose" -->
							<!--							> -->
							<!--								<FormItem> -->
							<!--									<FormLabel>{{ t('games.duel.settings.pointsPerLose') }}</FormLabel> -->
							<!--									<FormControl> -->
							<!--										<Input type="number" v-bind="componentField" :max="999999" /> -->
							<!--									</FormControl> -->
							<!--									<FormMessage /> -->
							<!--								</FormItem> -->
							<!--							</FormField> -->
						</div>
					</div>

					<div class="rounded-lg border p-4">
						<h4 class="mb-4 text-sm font-medium">
							{{ t('games.duel.messages.title') }}
						</h4>
						<div class="space-y-4">
							<UiFormField v-slot="{ componentField }" name="startMessage">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.messages.start.title') }}</UiFormLabel>
									<UiFormControl>
										<UiInput v-bind="componentField" />
									</UiFormControl>
									<p class="text-sm text-muted-foreground">
										{{ t('games.duel.messages.start.description') }}
									</p>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="resultMessage">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.messages.result.title') }}</UiFormLabel>
									<UiFormControl>
										<UiInput v-bind="componentField" />
									</UiFormControl>
									<p class="text-sm text-muted-foreground">
										{{ t('games.duel.messages.result.description') }}
									</p>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>

							<UiFormField v-slot="{ componentField }" name="bothDieMessage">
								<UiFormItem>
									<UiFormLabel>{{ t('games.duel.messages.bothDie.title') }}</UiFormLabel>
									<UiFormControl>
										<UiInput v-bind="componentField" />
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>
						</div>
					</div>

					<div class="flex justify-between">
						<UiAlertDialog>
							<UiAlertDialogTrigger asChild>
								<UiButton variant="outline">
									{{ t('sharedButtons.setDefaultSettings') }}
								</UiButton>
							</UiAlertDialogTrigger>
							<UiAlertDialogContent>
								<UiAlertDialogHeader>
									<UiAlertDialogTitle>{{ t('sharedTexts.areYouSure') }}</UiAlertDialogTitle>
									<UiAlertDialogDescription> Are you sure? </UiAlertDialogDescription>
								</UiAlertDialogHeader>
								<UiAlertDialogFooter>
									<UiAlertDialogCancel>
										{{ t('sharedButtons.cancel') }}
									</UiAlertDialogCancel>
									<UiAlertDialogAction @click="resetSettings">
										{{ t('sharedButtons.confirm') }}
									</UiAlertDialogAction>
								</UiAlertDialogFooter>
							</UiAlertDialogContent>
						</UiAlertDialog>

						<UiButton type="submit" @click="onSubmit">
							{{ t('sharedButtons.save') }}
						</UiButton>
					</div>
				</div>
			</form>
		</DialogOrSheet>
	</UiDialog>
</template>
