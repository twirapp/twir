<script setup lang="ts">
import { GavelIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { ref, toRaw, watch } from 'vue'


import { formSchema, useVotebanForm } from './composables/use-voteban-form'

import Card from '#layers/dashboard/components/games/card.vue'










import CommandButton from '~/features/commands/ui/command-button.vue'
import { VoteBanGameVotingMode } from '~/gql/graphql'

const isDialogOpen = ref(false)
const { t } = useI18n()
const { initialValues, save, settings } = useVotebanForm()

const votebanForm = useForm({
	validationSchema: formSchema,
	initialValues,
	validateOnMount: false,
	keepValuesOnUnmount: true,
})

watch(
	() => [
    votebanForm.values.chatVotesWordsPositive,
    votebanForm.values.chatVotesWordsNegative,
  ],
  () => {
    votebanForm.validate()
  }
)

watch(
	settings,
	(newSettings) => {
		if (!newSettings) return
		votebanForm.setValues(toRaw(newSettings.gamesVoteban))
	},
	{ immediate: true }
)

const onSubmit = votebanForm.handleSubmit(async (values) => {
	await save(values)
	isDialogOpen.value = false
})

function resetSettings() {
	votebanForm.setValues(initialValues)
}
</script>

<template>
	<UiDialog v-model:open="isDialogOpen">
		<UiDialogTrigger asChild>
			<Card
				title="Voteban"
				:icon="GavelIcon"
				:icon-stroke="1"
				:description="t('games.voteban.description')"
			/>
		</UiDialogTrigger>

		<UiDialogContent class="sm:max-w-[625px] max-h-[80vh] overflow-y-auto">
			<UiDialogHeader>
				<UiDialogTitle>Voteban</UiDialogTitle>
			</UiDialogHeader>

			<form @submit="onSubmit">
				<div class="grid gap-4 py-4">
					<div class="flex items-center gap-6">
						<UiFormField v-slot="{ value, handleChange }" name="enabled">
							<UiFormItem class="flex flex-col items-center gap-1">
								<UiFormLabel>{{ t('sharedTexts.enabled') }}</UiFormLabel>
								<UiFormControl>
									<UiSwitch :model-value="value" @update:model-value="handleChange" />
								</UiFormControl>
							</UiFormItem>
						</UiFormField>

						<CommandButton name="voteban" />
					</div>

					<UiSeparator />

					<div class="space-y-4">
						<h4 class="font-medium">Messages</h4>

						<UiFormField v-slot="{ componentField }" name="initMessage">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.initialMessage') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="banMessage">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.banMessage') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="surviveMessage">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.surviveMessage') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
					</div>

					<UiSeparator />

					<div class="space-y-4">
						<h4 class="font-medium">
							{{ t('sharedTexts.settings') }}
						</h4>

						<UiFormField v-slot="{ componentField }" name="votingMode">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.voteMode') }}</UiFormLabel>
								<UiSelect v-bind="componentField">
									<UiFormControl>
										<UiSelectTrigger>
											<UiSelectValue />
										</UiSelectTrigger>
									</UiFormControl>
									<UiSelectContent>
										<UiSelectItem :value="VoteBanGameVotingMode.Chat"> Chat </UiSelectItem>
										<UiSelectItem :value="VoteBanGameVotingMode.Polls" disabled>
											Twitch polls (soon)
										</UiSelectItem>
									</UiSelectContent>
								</UiSelect>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ field }" name="chatVotesWordsPositive">
							<UiFormItem>
								<UiFormControl>
									<UiFormLabel>{{ t('games.voteban.wordsPositive') }}</UiFormLabel>
									<UiTagsInput
										:model-value="field.value"
										@update:model-value="field['onUpdate:modelValue']"
										:max="10"
									>
										<UiTagsInputItem v-for="item in field.value" :key="item" :value="item">
											<UiTagsInputItemText />
											<UiTagsInputItemDelete />
										</UiTagsInputItem>

										<UiTagsInputInput placeholder="Enter words..." />
									</UiTagsInput>
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ field }" name="chatVotesWordsNegative">
							<UiFormItem>
								<UiFormControl>
									<UiFormLabel>{{ t('games.voteban.wordsNegative') }}</UiFormLabel>
									<UiTagsInput
										:model-value="field.value"
										@update:model-value="field['onUpdate:modelValue']"
										:placeholder="t('games.voteban.wordsPositive')"
										:max="10"
									>
										<UiTagsInputItem v-for="item in field.value" :key="item" :value="item">
											<UiTagsInputItemText />
											<UiTagsInputItemDelete />
										</UiTagsInputItem>

										<UiTagsInputInput placeholder="Enter words..." />
									</UiTagsInput>
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="voteDuration">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.voteDuration') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="number" v-bind="componentField" min="1" max="86400" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="neededVotes">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.neededVotes') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="number" v-bind="componentField" min="1" max="999999" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="timeoutSeconds">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.banDuration') }}</UiFormLabel>
								<UiFormControl>
									<UiInput type="number" v-bind="componentField" min="1" max="86400" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
					</div>

					<UiSeparator />

					<div class="space-y-4">
						<h4 class="font-medium">Moderators</h4>

						<UiFormField v-slot="{ value, handleChange }" name="timeoutModerators">
							<UiFormItem class="flex items-center justify-between">
								<UiFormLabel>{{ t('games.voteban.timeoutModerators') }}</UiFormLabel>
								<UiFormControl>
									<UiSwitch :model-value="value" @update:model-value="handleChange" />
								</UiFormControl>
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="banMessageModerators">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.banMessageModerators') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>

						<UiFormField v-slot="{ componentField }" name="surviveMessageModerators">
							<UiFormItem>
								<UiFormLabel>{{ t('games.voteban.surviveMessageModerators') }}</UiFormLabel>
								<UiFormControl>
									<UiTextarea v-bind="componentField" />
								</UiFormControl>
								<UiFormMessage />
							</UiFormItem>
						</UiFormField>
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

						<UiButton type="submit">
							{{ t('sharedButtons.save') }}
						</UiButton>
					</div>
				</div>
			</form>
		</UiDialogContent>
	</UiDialog>
</template>
